package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"com.dotvinci.tm/internal/domain/schema"
	"com.dotvinci.tm/internal/tmd/tapi/bases"
)

type SqlCrudBase struct{}

type CrudConfig struct {
	Operation string `json:"operation"`
	Entity    string `json:"entity"`
	Driver    string `json:"driver"`
	DSN       string `json:"dsn"`
}

var dbPool = map[string]*sql.DB{}
var dbMu sync.Mutex

func (SqlCrudBase) NameBase() string {
	return "sql-crud"
}

func (SqlCrudBase) Exec(ctx *bases.BaseContext) error {
	cfg, err := readConfig(ctx.Route.BaseConfigs)
	if err != nil {
		return err
	}
	entity, err := schema.MustEntity(cfg.Entity)
	if err != nil {
		return err
	}
	db, err := getDB(cfg.Driver, cfg.DSN)
	if err != nil {
		return err
	}

	switch strings.ToLower(cfg.Operation) {
	case "create":
		return create(ctx, db, entity, cfg.Driver)
	case "read", "get":
		return readByID(ctx, db, entity, cfg.Driver)
	case "update":
		return update(ctx, db, entity, cfg.Driver)
	case "delete":
		return deleteByID(ctx, db, entity, cfg.Driver)
	case "list":
		return list(ctx, db, entity)
	default:
		return fmt.Errorf("unsupported CRUD operation: %s", cfg.Operation)
	}
}

func readConfig(in map[string]any) (CrudConfig, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return CrudConfig{}, err
	}
	var cfg CrudConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return CrudConfig{}, err
	}
	if cfg.Operation == "" || cfg.Entity == "" || cfg.Driver == "" || cfg.DSN == "" {
		return CrudConfig{}, fmt.Errorf("sql-crud requires operation, entity, driver and dsn in base-configs")
	}
	return cfg, nil
}

func getDB(driver string, dsn string) (*sql.DB, error) {
	key := fmt.Sprintf("%s|%s", driver, dsn)
	dbMu.Lock()
	defer dbMu.Unlock()
	if dbPool[key] != nil {
		return dbPool[key], nil
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	dbPool[key] = db
	return db, nil
}

func create(ctx *bases.BaseContext, db *sql.DB, entity schema.Schema, driver string) error {
	var payload map[string]any
	if err := json.NewDecoder(ctx.Request.Body).Decode(&payload); err != nil {
		return err
	}
	if errs := schema.ValidateObject(payload, entity); len(errs) > 0 {
		ctx.Writter.WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(ctx.Writter).Encode(map[string]any{"errors": errs})
	}

	keys := sortedFieldsForWrite(entity, payload, true)
	cols := make([]string, 0, len(keys))
	holders := make([]string, 0, len(keys))
	values := make([]any, 0, len(keys))
	for i, k := range keys {
		cols = append(cols, k)
		holders = append(holders, placeholder(driver, i+1))
		values = append(values, payload[k])
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", entity.Table, strings.Join(cols, ","), strings.Join(holders, ","))
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	ctx.Writter.WriteHeader(http.StatusCreated)
	return json.NewEncoder(ctx.Writter).Encode(map[string]any{"status": "created"})
}

func readByID(ctx *bases.BaseContext, db *sql.DB, entity schema.Schema, driver string) error {
	id := ctx.Request.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing query param 'id'")
	}
	pk := primaryField(entity)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s LIMIT 1", entity.Table, pk, placeholder(driver, 1))
	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	result, err := rowsToMaps(rows)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		ctx.Writter.WriteHeader(http.StatusNotFound)
		return json.NewEncoder(ctx.Writter).Encode(map[string]any{"error": "not found"})
	}
	return json.NewEncoder(ctx.Writter).Encode(result[0])
}

func update(ctx *bases.BaseContext, db *sql.DB, entity schema.Schema, driver string) error {
	id := ctx.Request.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing query param 'id'")
	}
	var payload map[string]any
	if err := json.NewDecoder(ctx.Request.Body).Decode(&payload); err != nil {
		return err
	}
	keys := sortedFieldsForWrite(entity, payload, false)
	sets := make([]string, 0, len(keys))
	values := make([]any, 0, len(keys)+1)
	for i, k := range keys {
		sets = append(sets, fmt.Sprintf("%s = %s", k, placeholder(driver, i+1)))
		values = append(values, payload[k])
	}
	pk := primaryField(entity)
	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s", entity.Table, strings.Join(sets, ","), pk, placeholder(driver, len(values)))
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}
	return json.NewEncoder(ctx.Writter).Encode(map[string]any{"status": "updated"})
}

func deleteByID(ctx *bases.BaseContext, db *sql.DB, entity schema.Schema, driver string) error {
	id := ctx.Request.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing query param 'id'")
	}
	pk := primaryField(entity)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = %s", entity.Table, pk, placeholder(driver, 1))
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return json.NewEncoder(ctx.Writter).Encode(map[string]any{"status": "deleted"})
}

func list(ctx *bases.BaseContext, db *sql.DB, entity schema.Schema) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", entity.Table))
	if err != nil {
		return err
	}
	defer rows.Close()
	result, err := rowsToMaps(rows)
	if err != nil {
		return err
	}
	return json.NewEncoder(ctx.Writter).Encode(result)
}

func rowsToMaps(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		item := map[string]any{}
		for i, c := range cols {
			if b, ok := vals[i].([]byte); ok {
				item[c] = string(b)
			} else {
				item[c] = vals[i]
			}
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func primaryField(entity schema.Schema) string {
	for name, f := range entity.Fields {
		if f.PrimaryKey != nil && *f.PrimaryKey {
			return name
		}
	}
	return "id"
}

func sortedFieldsForWrite(entity schema.Schema, payload map[string]any, skipAuto bool) []string {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		if _, ok := entity.Fields[k]; !ok {
			continue
		}
		if skipAuto {
			if field := entity.Fields[k]; field.AutoIncrement != nil && *field.AutoIncrement {
				continue
			}
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func placeholder(driver string, index int) string {
	d := strings.ToLower(driver)
	if strings.Contains(d, "postgres") || d == "pgx" || d == "pq" {
		return fmt.Sprintf("$%d", index)
	}
	return "?"
}
