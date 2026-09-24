package connection

import (
	"sync"

	"github.com/google/uuid"
)

type PendingRequests struct {
	mu       sync.RWMutex
	requests map[uuid.UUID]chan []byte
}

func NewPendingRequests() *PendingRequests {
	return &PendingRequests{
		requests: make(map[uuid.UUID]chan []byte),
	}
}

func (p *PendingRequests) Add(
	requestID uuid.UUID,
) <-chan []byte {
	ch := make(chan []byte, 1)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.requests[requestID] = ch

	return ch
}

func (p *PendingRequests) Resolve(
	requestID uuid.UUID,
	data []byte,
) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch, ok := p.requests[requestID]
	if !ok {
		return false
	}

	delete(p.requests, requestID)
	ch <- data
	close(ch)

	return true
}

func (p *PendingRequests) Remove(
	requestID uuid.UUID,
) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch, ok := p.requests[requestID]
	if !ok {
		return false
	}

	delete(p.requests, requestID)
	close(ch)

	return true
}

func (p *PendingRequests) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.requests)
}
