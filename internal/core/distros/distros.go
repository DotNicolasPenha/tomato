package distros

import (
	"fmt"
	"net/http"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/loader"
)

type DistroExecContext struct {
	Manifest loader.Manifest
	Server   *http.Server
	Mux      *http.ServeMux
}

type Distro interface {
	NameDistro() string
	Exec(ctx DistroExecContext) error
}

var Registry = map[string]Distro{}

func Register(d Distro) {
	name := d.NameDistro()
	if name == "" {
		logger.Fatal("distro without name")
	}
	Registry[name] = d
	logger.Ok(fmt.Sprintf("distro '%s' registered", name))
}
func All() map[string]Distro {
	return Registry
}
func Find(distroname string) (Distro, error) {
	var distro = Registry[distroname]
	if distro == nil {
		return nil, fmt.Errorf("distro '%s' not found", distroname)
	}
	return distro, nil
}
