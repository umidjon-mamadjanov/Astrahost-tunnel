package connection

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu         sync.RWMutex
	tunnels    map[string]*Session
	subdomains map[string]*Session
}

func NewRegistry() *Registry {
	return &Registry{
		tunnels:    make(map[string]*Session),
		subdomains: make(map[string]*Session),
	}
}

func (r *Registry) Register(tunnelID string, session *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tunnels[tunnelID] = session
}

func (r *Registry) Unregister(tunnelID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.tunnels[tunnelID]
	if !exists {
		return
	}

	delete(r.tunnels, tunnelID)

	subdomain := session.GetSubdomain()
	if subdomain != "" {
		delete(r.subdomains, subdomain)
	}
}

func (r *Registry) Get(tunnelID string) (*Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.tunnels[tunnelID]

	return session, ok
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tunnels)
}

func (r *Registry) AllocateSubdomain(
	requested string,
	session *Session,
) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.subdomains[requested]; !exists {
		r.subdomains[requested] = session
		return requested
	}

	for i := 1; ; i++ {
		candidate := requested + fmt.Sprintf("%02d", i)

		if _, exists := r.subdomains[candidate]; !exists {
			r.subdomains[candidate] = session
			return candidate
		}
	}
}

func (r *Registry) UnregisterSubdomain(subdomain string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.subdomains, subdomain)
}

func (r *Registry) GetBySubdomain(
	subdomain string,
) (*Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.subdomains[subdomain]
	return session, ok
}
