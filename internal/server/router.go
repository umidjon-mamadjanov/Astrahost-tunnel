package server

import (
	"net/http"
	"strings"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.health)
	mux.HandleFunc("/connect", a.ws.Handle)
	mux.HandleFunc("/tunnel/", a.tunnel)
	mux.HandleFunc("/", a.publicTunnel)

	return mux
}

func (a *App) publicTunnel(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if host == "" {
		http.Error(w, "host is required", http.StatusBadRequest)
		return
	}

	host = strings.Split(host, ":")[0]

	baseDomain := strings.TrimSuffix(
		strings.ToLower(a.cfg.BaseDomain),
		".",
	)

	if baseDomain == "" {
		http.Error(w, "public routing is not configured", http.StatusNotFound)
		return
	}

	suffix := "." + baseDomain

	if !strings.HasSuffix(strings.ToLower(host), suffix) {
		http.Error(w, "unknown host", http.StatusNotFound)
		return
	}

	subdomain := strings.TrimSuffix(
		strings.ToLower(host[:len(host)-len(suffix)]),
		".",
	)

	if subdomain == "" {
		http.Error(w, "subdomain is required", http.StatusBadRequest)
		return
	}

	session, ok := a.registry.GetBySubdomain(subdomain)
	if !ok {
		http.Error(w, "tunnel not found", http.StatusNotFound)
		return
	}

	if session.GetState() != connection.StateReady {
		http.Error(w, "tunnel is not ready", http.StatusServiceUnavailable)
		return
	}

	a.proxyHTTP(w, r, session, r.URL.Path)
}
