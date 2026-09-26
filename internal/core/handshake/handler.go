package handshake

import (
	"fmt"
	"strings"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

type Config struct {
	BaseDomain string
	Scheme     string
}

type Handler struct {
	registry *connection.Registry
	cfg      Config
}

func New(registry *connection.Registry, cfg Config) *Handler {
	return &Handler{
		registry: registry,
		cfg:      cfg,
	}
}

func (h *Handler) HandleConnect(
	session *connection.Session,
	req protocol.ConnectRequest,
) (*protocol.ConnectResponse, error) {
	if req.ProtocolVersion != protocol.Version {
		return nil, fmt.Errorf("unsupported protocol version")
	}

	if req.LocalHost == "" {
		return nil, fmt.Errorf("local host is required")
	}

	if req.LocalPort == 0 {
		return nil, fmt.Errorf("local port is required")
	}

	session.ProtocolVersion = req.ProtocolVersion
	session.ClientVersion = req.ClientVersion
	session.Platform = req.Platform
	session.Architecture = req.Architecture

	session.SetLocalTarget(
		req.LocalHost,
		req.LocalPort,
	)

	tunnelID := uuid.NewString()
	session.SetTunnelID(tunnelID)

	subdomain := h.registry.AllocateSubdomain(
		req.TunnelName,
		session,
	)

	session.SetSubdomain(subdomain)
	session.SetState(connection.StateReady)

	h.registry.Register(tunnelID, session)

	publicURL := ""

	if h.cfg.BaseDomain != "" && subdomain != "" {
		scheme := h.cfg.Scheme
		if scheme == "" {
			scheme = "https"
		}

		baseDomain := strings.TrimSuffix(h.cfg.BaseDomain, ".")
		publicURL = scheme + "://" + subdomain + "." + baseDomain
	}

	return &protocol.ConnectResponse{
		ProtocolVersion:   protocol.Version,
		SessionID:         session.ID,
		TunnelID:          tunnelID,
		HeartbeatInterval: 20,
		PublicURL:         publicURL,
	}, nil
}
