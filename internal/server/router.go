package server

import "net/http"

func (a *App) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.health)

	mux.HandleFunc("/connect", a.ws.Handle)

	return mux
}
