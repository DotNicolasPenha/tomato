package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/common/reader"
)

type Field struct {
	Type          string `json:"type"`
	Required      *bool  `json:"required"`
	Max           *int   `json:"max"`
	Min           *int   `json:"min"`
	PrimaryKey    *bool  `json:"primaryKey"`
	AutoIncrement *bool  `json:"autoIncrement"`
}

type Schema struct {
	Name   string           `json:"name"`
	Table  string           `json:"table"`
	Fields map[string]Field `json:"fields"`
	Kind   string           `json:"-"`
}

var registry = map[string]Schema{}

func LoadDomain(cwd string) {
	registry = map[string]Schema{}
	loadKind(cwd, "entity", "@domain", "entitys")
	loadKind(cwd, "dto", "@domain", "dtos")
}

func loadKind(cwd string, kind string, parts ...string) {
	dir := filepath.Join(append([]string{cwd}, parts...)...)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Fatal(fmt.Sprintf("cannot read domain dir %s: %s", dir, err))
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		fullpath := filepath.Join(dir, entry.Name())
		s := reader.Json[Schema](fullpath)
		if s.Name == "" {
			logger.Fatal(fmt.Sprintf("schema name is required at path: %s", fullpath))
		}
		if kind == "entity" && s.Table == "" {
			s.Table = strings.ToLower(s.Name)
		}
		if len(s.Fields) == 0 {
			logger.Fatal(fmt.Sprintf("schema fields is required at path: %s", fullpath))
		}
		s.Kind = kind
		registry[typeKey(kind, s.Name)] = *s
		logger.Info(fmt.Sprintf("%s schema loaded: %s", kind, s.Name))
	}
}

func typeKey(kind string, name string) string {
	return fmt.Sprintf("@%s:%s", kind, name)
}

func Find(typeRef string) (Schema, bool) {
	s, ok := registry[typeRef]
	return s, ok
}

func MustEntity(name string) (Schema, error) {
	typeRef := fmt.Sprintf("@entity:%s", name)
	s, ok := Find(typeRef)
	if !ok {
		return Schema{}, fmt.Errorf("entity '%s' not found in @domain/entitys", name)
	}
	return s, nil
}

func ValidateObject(payload map[string]any, schema Schema) []error {
	var errs []error
	for fieldName, rule := range schema.Fields {
		value, exists := payload[fieldName]
		required := rule.Required != nil && *rule.Required
		if required && !exists {
			errs = append(errs, fmt.Errorf("field '%s' is required", fieldName))
			continue
		}
		if !exists {
			continue
		}
		if !validateType(value, rule.Type) {
			errs = append(errs, fmt.Errorf("field '%s' must be of type '%s'", fieldName, rule.Type))
			continue
		}
		if err := validateMinMax(fieldName, value, rule); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func validateType(v any, expected string) bool {
	switch expected {
	case "string":
		_, ok := v.(string)
		return ok
	case "number", "int", "float":
		_, ok := v.(float64)
		return ok
	case "bool":
		_, ok := v.(bool)
		return ok
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	default:
		return false
	}
}

func validateMinMax(field string, value any, rule Field) error {
	switch v := value.(type) {
	case string:
		l := len(v)
		if rule.Min != nil && l < *rule.Min {
			return fmt.Errorf("field '%s' length < min (%d)", field, *rule.Min)
		}
		if rule.Max != nil && l > *rule.Max {
			return fmt.Errorf("field '%s' length > max (%d)", field, *rule.Max)
		}
	case float64:
		if rule.Min != nil && int(v) < *rule.Min {
			return fmt.Errorf("field '%s' < min (%d)", field, *rule.Min)
		}
		if rule.Max != nil && int(v) > *rule.Max {
			return fmt.Errorf("field '%s' > max (%d)", field, *rule.Max)
		}
	}
	return nil
}
