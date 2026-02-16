package declarator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/common/reader"
)

type TapiRoute struct {
	Path                   string                         `json:"path"`
	Base                   string                         `json:"base"`
	BaseConfigs            map[string]any                 `json:"base-configs"`
	Method                 string                         `json:"method"`
	Request_RequiredFormat TapiRouteRequestRequiredFormat `json:"request-requiredFormat"`
}
type TapiRouteRequestRequiredFormat struct {
	Body_json  *map[string]TapiRouteBodyJsonPropertie `json:"body-json"`
	Headers    *TapiRouteHeaders                      `json:"headers"`
	Querys     *map[string]any                        `json:"querys"`
	PathParams *map[string]any                        `json:"path-params"`
}
type TapiRouteHeaders struct {
	ContentLength *int            `json:"content-length"`
	Authorization *string         `json:"authorization"`
	Cookies       *map[string]any `json:"cookies"`
}
type TapiRouteBodyJsonPropertie struct {
	Required *bool   `json:"required"`
	Type     *string `json:"type"`
	Max      *int    `json:"max"`
	Min      *int    `json:"min"`
}

func DeclareRoutes(path string) (map[string]TapiRoute, error) {
	var routesDeclared = map[string]TapiRoute{}
	walkerRoutes(path, routesDeclared)
	return routesDeclared, nil
}
func walkerRoutes(path string, routes map[string]TapiRoute) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("cannot read dir %s", path))
	}
	for _, dirent := range dir {
		fullpath := filepath.Join(path, dirent.Name())
		if dirent.Type().IsDir() {
			walkerRoutes(fullpath, routes)
			continue
		}
		if !dirent.IsDir() && dirent.Name() == "index.json" {
			route := loadRoute(fullpath)
			routes[route.Path] = route
		}
	}
	return nil
}
func loadRoute(path string) TapiRoute {
	routeJson := reader.Json[TapiRoute](path)
	if routeJson.Base == "" {
		logger.Error(
			fmt.Sprintf("in route %s base is undefined", path),
		)
	}
	if routeJson.Path == "" {
		logger.Error(
			fmt.Sprintf("in route %s path is undefined", path),
		)
	}
	allowed := map[string]struct{}{
		"get":    {},
		"post":   {},
		"delete": {},
		"put":    {},
	}

	method := strings.ToLower(routeJson.Method)

	if method == "" {
		logger.Error(
			fmt.Sprintf("in route %s method is undefined", path),
		)
	}

	if _, ok := allowed[method]; !ok {
		logger.Error(fmt.Sprintf("invalid HTTP method: %s", method))
	}
	if routeJson.Request_RequiredFormat.Body_json != nil {
		var errs []string
		for n, v := range *routeJson.Request_RequiredFormat.Body_json {
			if v.Required == nil {
				err := fmt.Sprintf("The propertie %s in route of path %s 'required' is empty", n, path)
				errs = append(errs, err)
				logger.Error(err)
			}
			if v.Type == nil {
				err := fmt.Sprintf("The propertie %s in route of path %s 'type' is empty", n, path)
				errs = append(errs, err)
				logger.Error(err)
			}
		}
		if len(errs) > 0 {
			logger.Fatal(fmt.Sprintf("connot load the route of path %s because the properties of 'request-requiredFormat' has empty fields", path))
		}
	}
	if method == "" || routeJson.Path == "" || routeJson.Base == "" {
		logger.Fatal(fmt.Sprintf("The route cannot be loaded because it has empty fields. path: %s", path))
	}
	return TapiRoute{
		Path:                   routeJson.Path,
		Base:                   routeJson.Base,
		Method:                 routeJson.Method,
		BaseConfigs:            routeJson.BaseConfigs,
		Request_RequiredFormat: routeJson.Request_RequiredFormat,
	}
}
