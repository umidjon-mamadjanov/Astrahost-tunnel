package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

const tunnelRequestTimeout = 30 * time.Second

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

	// ------------------------------------------------------------
	// Find tunnel
	// ------------------------------------------------------------

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

	if session.PendingRequests == nil {
		http.Error(
			w,
			"pending request manager is not initialized",
			http.StatusInternalServerError,
		)
		return
	}

	// ------------------------------------------------------------
	// Get client's local target
	// ------------------------------------------------------------

	localHost, localPort := session.GetLocalTarget()

	if localHost == "" {
		http.Error(
			w,
			"local host is not configured",
			http.StatusInternalServerError,
		)
		return
	}

	if localPort == 0 {
		http.Error(
			w,
			"local port is not configured",
			http.StatusInternalServerError,
		)
		return
	}

	targetURL := "http://" +
		localHost +
		":" +
		strconv.Itoa(int(localPort)) +
		targetPath

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// ------------------------------------------------------------
	// Read request body
	// ------------------------------------------------------------

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			"failed to read request body",
			http.StatusBadRequest,
		)
		return
	}

	// ------------------------------------------------------------
	// Convert headers
	// ------------------------------------------------------------

	headers := make([]protocol.HTTPHeader, 0, len(r.Header))

	for name, values := range r.Header {
		for _, value := range values {
			headers = append(
				headers,
				protocol.HTTPHeader{
					Name:  name,
					Value: value,
				},
			)
		}
	}

	// ------------------------------------------------------------
	// Create request ID
	// ------------------------------------------------------------

	requestID := uuid.New()

	httpRequest := protocol.HTTPRequest{
		Method:  r.Method,
		URL:     targetURL,
		Headers: headers,
		Body:    body,
	}

	requestPacket, err := protocol.NewHTTPRequestPacket(
		requestID,
		httpRequest,
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("create HTTP_REQUEST: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	// Register pending request BEFORE sending the packet.
	responseCh := session.PendingRequests.Add(requestID)

	// ------------------------------------------------------------
	// Send HTTP_REQUEST to client
	// ------------------------------------------------------------

	if err := a.ws.SendPacket(session, requestPacket); err != nil {
		session.PendingRequests.Remove(requestID)

		http.Error(
			w,
			fmt.Sprintf("send HTTP_REQUEST: %v", err),
			http.StatusBadGateway,
		)
		return
	}

	log.Printf(
		"HTTP_REQUEST sent | Tunnel: %s | Request: %s | %s %s",
		tunnelID,
		requestID,
		r.Method,
		targetURL,
	)

	// ------------------------------------------------------------
	// Wait for HTTP_RESPONSE
	// ------------------------------------------------------------

	timer := time.NewTimer(tunnelRequestTimeout)
	defer timer.Stop()

	var responseData []byte

	select {
	case responseData = <-responseCh:

	case <-timer.C:
		session.PendingRequests.Remove(requestID)

		http.Error(
			w,
			"tunnel request timeout",
			http.StatusGatewayTimeout,
		)
		return

	case <-r.Context().Done():
		session.PendingRequests.Remove(requestID)
		return
	}

	// ------------------------------------------------------------
	// Decode HTTP_RESPONSE
	// ------------------------------------------------------------

	response, err := protocol.DecodeHTTPResponse(responseData)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("decode HTTP_RESPONSE: %v", err),
			http.StatusBadGateway,
		)
		return
	}

	// ------------------------------------------------------------
	// Write response
	// ------------------------------------------------------------

	for _, header := range response.Headers {
		w.Header().Add(
			header.Name,
			header.Value,
		)
	}

	w.WriteHeader(response.StatusCode)

	if len(response.Body) > 0 {
		_, _ = w.Write(response.Body)
	}

	log.Printf(
		"HTTP_RESPONSE received | Tunnel: %s | Request: %s | Status: %d",
		tunnelID,
		requestID,
		response.StatusCode,
	)
}
