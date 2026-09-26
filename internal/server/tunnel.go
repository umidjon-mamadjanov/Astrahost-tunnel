package server

import (
	"net/http"
	"strings"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

func (a *App) tunnel(w http.ResponseWriter, r *http.Request) {
	// Expected:
	// /tunnel/{tunnel_id}/path

	path := strings.TrimPrefix(r.URL.Path, "/tunnel/")

	if path == "" {
		http.Error(
			w,
			"tunnel ID is required",
			http.StatusBadRequest,
		)
		return
	}

	parts := strings.SplitN(path, "/", 2)

	tunnelID := parts[0]

	if tunnelID == "" {
		http.Error(
			w,
			"tunnel ID is required",
			http.StatusBadRequest,
		)
		return
	}

	targetPath := "/"

	if len(parts) == 2 && parts[1] != "" {
		targetPath += parts[1]
	}

	session, ok := a.registry.Get(tunnelID)
	if !ok {
		http.Error(
			w,
			"tunnel not found",
			http.StatusNotFound,
		)
		return
	}

	if session.GetState() != connection.StateReady {
		http.Error(
			w,
			"tunnel is not ready",
			http.StatusServiceUnavailable,
		)
		return
	}

	a.proxyHTTP(w, r, session, targetPath)
}
