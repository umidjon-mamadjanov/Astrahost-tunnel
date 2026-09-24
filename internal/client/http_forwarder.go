package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/astrahost/astrahost-tunnel/protocol"
)

type HTTPForwarder struct {
	Client *http.Client
}

func NewHTTPForwarder() *HTTPForwarder {
	return &HTTPForwarder{
		Client: &http.Client{},
	}
}

func (f *HTTPForwarder) Forward(
	req protocol.HTTPRequest,
) (protocol.HTTPResponse, error) {
	body := bytes.NewReader(req.Body)

	httpReq, err := http.NewRequest(
		req.Method,
		req.URL,
		body,
	)
	if err != nil {
		return protocol.HTTPResponse{}, fmt.Errorf(
			"create HTTP request: %w",
			err,
		)
	}

	for _, header := range req.Headers {
		httpReq.Header.Add(header.Name, header.Value)
	}

	resp, err := f.Client.Do(httpReq)
	if err != nil {
		return protocol.HTTPResponse{}, fmt.Errorf(
			"forward HTTP request: %w",
			err,
		)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return protocol.HTTPResponse{}, fmt.Errorf(
			"read HTTP response body: %w",
			err,
		)
	}

	headers := make([]protocol.HTTPHeader, 0, len(resp.Header))

	for name, values := range resp.Header {
		for _, value := range values {
			headers = append(headers, protocol.HTTPHeader{
				Name:  name,
				Value: value,
			})
		}
	}

	return protocol.HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       responseBody,
	}, nil
}
