package tunnel

import "sync"

type Registry struct {
	mu      sync.RWMutex
	tunnels map[string]*Tunnel
}

func NewRegistry() *Registry {
	return &Registry{
		tunnels: make(map[string]*Tunnel),
	}
}

// Shu yerdan boshlanadi ↓

func (r *Registry) Add(t *Tunnel) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tunnels[t.ID.String()] = t
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tunnels, id)
}

func (r *Registry) Get(id string) (*Tunnel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tunnels[id]
	return t, ok
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tunnels)
}
