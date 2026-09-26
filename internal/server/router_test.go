package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

func TestPublicTunnelUnknownSubdomain(t *testing.T) {
	app := &App{
		cfg: Config{
			BaseDomain: "astrahost.bond",
			Scheme:     "https",
		},
		registry: connection.NewRegistry(),
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"http://unknown.astrahost.bond/",
		nil,
	)
	req.Host = "unknown.astrahost.bond"

	rec := httptest.NewRecorder()

	app.publicTunnel(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}
