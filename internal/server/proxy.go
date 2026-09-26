package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

const tunnelRequestTimeout = 30 * time.Second

func (a *App) proxyHTTP(
	w http.ResponseWriter,
	r *http.Request,
	session *connection.Session,
	targetPath string,
) {
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

	if targetPath == "" {
		targetPath = "/"
	}

	targetURL := "http://" +
		localHost +
		":" +
		strconv.Itoa(int(localPort)) +
		targetPath

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			"failed to read request body",
			http.StatusBadRequest,
		)
		return
	}

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

	responseCh := session.PendingRequests.Add(requestID)

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
		"HTTP_REQUEST sent | Subdomain: %s | Request: %s | %s %s",
		session.GetSubdomain(),
		requestID,
		r.Method,
		targetURL,
	)

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

	response, err := protocol.DecodeHTTPResponse(responseData)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("decode HTTP_RESPONSE: %v", err),
			http.StatusBadGateway,
		)
		return
	}

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
}
