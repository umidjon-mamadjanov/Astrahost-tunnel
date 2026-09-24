package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrahost/astrahost-tunnel/protocol"
)

func TestHTTPForwarder(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"unexpected method: got %s, want POST",
					r.Method,
				)
			}

			if r.Header.Get("X-Test") != "astra" {
				t.Errorf("missing X-Test header")
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("Hello from test server"))
		}),
	)
	defer server.Close()

	forwarder := NewHTTPForwarder()

	response, err := forwarder.Forward(
		protocol.HTTPRequest{
			Method: http.MethodPost,
			URL:    server.URL,
			Headers: []protocol.HTTPHeader{
				{
					Name:  "X-Test",
					Value: "astra",
				},
			},
			Body: []byte("request body"),
		},
	)
	if err != nil {
		t.Fatalf("forward request: %v", err)
	}

	if response.StatusCode != http.StatusCreated {
		t.Fatalf(
			"unexpected status: got %d, want %d",
			response.StatusCode,
			http.StatusCreated,
		)
	}

	if string(response.Body) != "Hello from test server" {
		t.Fatalf(
			"unexpected body: got %q",
			response.Body,
		)
	}
}

func TestHTTPForwarderInvalidURL(t *testing.T) {
	forwarder := NewHTTPForwarder()

	_, err := forwarder.Forward(
		protocol.HTTPRequest{
			Method: http.MethodGet,
			URL:    "://invalid-url",
		},
	)

	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
