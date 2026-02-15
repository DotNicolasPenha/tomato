package renderer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/tmd/tapi/bases"
	"com.dotvinci.tm/internal/tmd/tapi/router/declarator"
)

// took the list of json (route) and render to http
func RenderRoutes(routes map[string]declarator.TapiRoute, ctx distros.DistroExecContext) {
	allowed := map[string]struct{}{
		"get":    {},
		"post":   {},
		"delete": {},
		"put":    {},
		"GET":    {},
		"POST":   {},
		"DELETE": {},
		"PUT":    {},
	}
	routesLoadedCount := 0
	for _, route := range routes {
		ctx.Mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			if _, ok := allowed[r.Method]; !ok {
				err := fmt.Sprintf("invalid HTTP method: %s", r.Method)
				logger.Error(err)
				w.WriteHeader(400)
				w.Write([]byte(err))
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				err := "error to decode json body"
				logger.Error(err)
				w.WriteHeader(500)
				w.Write([]byte(err))
			}
			if route.Request_RequiredFormat != nil {
				if route.Request_RequiredFormat.Body_json != nil {
					ValidateBody(body, route.Request_RequiredFormat.Body_json)
				}
			}
			base := bases.Find(route.Base)
			if base == nil {
				err := fmt.Sprintf("base not found: %s", route.Base)
				logger.Error(err)
				w.WriteHeader(500)
				w.Write([]byte(err))
			}
			baseCtx := bases.BaseContext{
				Route:    route,
				Manifest: ctx.Manifest,
				Server:   ctx.Server,
				Mux:      ctx.Mux,
				Writter:  w,
				Request:  r,
			}
			err := base.Exec(&baseCtx)
			if err != nil {
				err := fmt.Sprintf("base %s exec error: %s", route.Base, err)
				logger.Error(err)
				w.WriteHeader(500)
				w.Write([]byte(err))
			}
			routesLoadedCount++
		})
	}
	logger.Info(fmt.Sprintf("(%d) routes loaded", routesLoadedCount))
}
func ValidateBody(
	body map[string]any,
	schema *map[string]declarator.TapiRouteBodyJsonPropertie,
) []error {

	var errs []error

	for field, rule := range *schema {
		value, exists := body[field]
		if rule.Required && !exists {
			errs = append(errs, fmt.Errorf("field '%s' is required", field))
			continue
		}
		if !exists {
			continue
		}
		if !validateType(value, rule.Type) {
			errs = append(errs, fmt.Errorf(
				"field '%s' must be of type '%s'", field, rule.Type,
			))
			continue
		}

		if err := validateMinMax(field, value, rule); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func validateMinMax(
	field string,
	value any,
	rule declarator.TapiRouteBodyJsonPropertie,
) error {

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
		n := int(v)
		if rule.Min != nil && n < *rule.Min {
			return fmt.Errorf("field '%s' < min (%d)", field, *rule.Min)
		}
		if rule.Max != nil && n > *rule.Max {
			return fmt.Errorf("field '%s' > max (%d)", field, *rule.Max)
		}
	}

	return nil
}
func validateType(v any, expected string) bool {
	switch expected {
	case "string":
		_, ok := v.(string)
		return ok
	case "number":
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
