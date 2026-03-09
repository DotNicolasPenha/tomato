package envx

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const envPrefix = "@env:"

func Resolve(v any) (any, error) {
	return resolveValue(v)
}

func resolveValue(v any) (any, error) {
	if v == nil {
		return nil, nil
	}

	switch value := v.(type) {
	case string:
		if strings.HasPrefix(value, envPrefix) {
			key := strings.TrimSpace(strings.TrimPrefix(value, envPrefix))
			if key == "" {
				return nil, fmt.Errorf("invalid @env usage: missing key")
			}
			envValue, ok := os.LookupEnv(key)
			if !ok {
				return nil, fmt.Errorf("environment variable '%s' not found", key)
			}
			return envValue, nil
		}
		return value, nil
	case []any:
		out := make([]any, 0, len(value))
		for _, item := range value {
			resolved, err := resolveValue(item)
			if err != nil {
				return nil, err
			}
			out = append(out, resolved)
		}
		return out, nil
	case map[string]any:
		out := map[string]any{}
		for k, item := range value {
			resolved, err := resolveValue(item)
			if err != nil {
				return nil, err
			}
			out[k] = resolved
		}
		return out, nil
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Map {
			iter := rv.MapRange()
			out := map[string]any{}
			for iter.Next() {
				key := fmt.Sprintf("%v", iter.Key().Interface())
				resolved, err := resolveValue(iter.Value().Interface())
				if err != nil {
					return nil, err
				}
				out[key] = resolved
			}
			return out, nil
		}
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			out := make([]any, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				resolved, err := resolveValue(rv.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				out = append(out, resolved)
			}
			return out, nil
		}
		return v, nil
	}
}
