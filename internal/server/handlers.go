package server

import (
	"net/http"
)

func (a *App) health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}

func (a *App) connect(w http.ResponseWriter, r *http.Request) {

	http.Error(
		w,
		"WebSocket not implemented",
		http.StatusNotImplemented,
	)
}
