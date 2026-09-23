package connection

import "sync"

type Registry struct {
	mu      sync.RWMutex
	tunnels map[string]*Session
}

func NewRegistry() *Registry {
	return &Registry{
		tunnels: make(map[string]*Session),
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

	delete(r.tunnels, tunnelID)
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
