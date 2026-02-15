package lta

import (
	"net/http"

	"com.dotvinci.tm/internal/core/loader"
)

type Lta struct {
	Manifest loader.Manifest
	Server   *http.Server
	Mux      *http.ServeMux
}

func (lta *Lta) Init() {
	lta.Server = &http.Server{
		Addr:    string(rune(*lta.Manifest.Port)),
		Handler: lta.Mux,
	}
	lta.Server.ListenAndServe()
}
