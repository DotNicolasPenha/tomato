package lta

import (
	"fmt"
	"net/http"
	"time"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/core/loader"
)

type Lta struct {
	Manifest *loader.Manifest
	Server   *http.Server
	Mux      *http.ServeMux
	Distros  []distros.Distro
}

func (lta *Lta) Init() error {
	start := time.Now()
	portStr := fmt.Sprintf(":%d", *lta.Manifest.Port)
	lta.Server = &http.Server{
		Addr:    portStr,
		Handler: lta.Mux,
	}
	elapsed := time.Since(start)
	logger.Ok(
		fmt.Sprintf(
			"%s:lta '%s' initialized in port %s",
			elapsed, *lta.Manifest.NameApplication, portStr),
	)
	err := lta.Server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}
func (lta *Lta) PlugDistro(distro distros.Distro) {
	lta.Distros = append(lta.Distros, distro)
	logger.Ok(
		fmt.Sprintf(
			"distro '%s' was pluged in lta '%s'",
			distro.NameDistro(), *lta.Manifest.NameApplication,
		),
	)
}
func (lta *Lta) ExecuteDistro(distroName string) error {
	for _, distro := range lta.Distros {
		if distroName == distro.NameDistro() {
			distroExecContext := distros.DistroExecContext{
				Manifest: *lta.Manifest,
				Server:   lta.Server,
				Mux:      lta.Mux,
			}
			err := distro.Exec(distroExecContext)
			logger.Ok(fmt.Sprintf("distro '%s' was executed", distro.NameDistro()))
			return err
		}
	}
	return fmt.Errorf("distro '%s' not found in lta '%s'.", distroName, *lta.Manifest.NameApplication)
}
