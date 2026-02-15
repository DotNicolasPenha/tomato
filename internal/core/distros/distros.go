package distros

import (
	// "net/http"

	"com.dotvinci.tm/internal/core/lta"
)

type Distro interface {
	NameDistro() string
	Exec(lta *lta.Lta) error
}

// type TomatoAPIs struct{}

// func (tmsAPIS *TomatoAPIs) NameDistro() string {
// 	return "tomato apis"
// }
// func (tmAPIS *TomatoAPIs) Exec(lta *lta.Lta) error {
// 	lta.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("test"))
// 	})
// 	return nil
// }
