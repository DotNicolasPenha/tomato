package bases

import (
	"net/http"

	"com.dotvinci.tm/internal/core/loader"
	"com.dotvinci.tm/internal/tmd/tapi/router/declarator"
)

type BaseContext struct {
	Route    declarator.TapiRoute
	Manifest loader.Manifest
	Server   *http.Server
	Mux      *http.ServeMux
	Writter  http.ResponseWriter
	Request  *http.Request
}
type Base interface {
	NameBase() string
	Exec(ctx *BaseContext) error
}

var Bases []Base

func RegistryBase(base Base) {
	Bases = append(Bases, base)
}
func Find(baseName string) Base {
	for _, b := range Bases {
		if b.NameBase() == baseName {
			return b
		}
	}
	return nil
}
