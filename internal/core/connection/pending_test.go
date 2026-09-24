package connection

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPendingRequestsResolve(t *testing.T) {
	pending := NewPendingRequests()

	requestID := uuid.New()
	response := []byte("HTTP response")

	ch := pending.Add(requestID)

	if pending.Count() != 1 {
		t.Fatalf("expected 1 pending request, got %d", pending.Count())
	}

	if !pending.Resolve(requestID, response) {
		t.Fatal("expected request to be resolved")
	}

	select {
	case got := <-ch:
		if string(got) != string(response) {
			t.Fatalf(
				"unexpected response: got %q, want %q",
				got,
				response,
			)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for response")
	}

	if pending.Count() != 0 {
		t.Fatalf("expected 0 pending requests, got %d", pending.Count())
	}
}

func TestPendingRequestsRemove(t *testing.T) {
	pending := NewPendingRequests()

	requestID := uuid.New()
	ch := pending.Add(requestID)

	if !pending.Remove(requestID) {
		t.Fatal("expected request to be removed")
	}

	if pending.Count() != 0 {
		t.Fatalf("expected 0 pending requests, got %d", pending.Count())
	}

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("channel was not closed")
	}
}

func TestPendingRequestsUnknownRequest(t *testing.T) {
	pending := NewPendingRequests()

	requestID := uuid.New()

	if pending.Resolve(requestID, []byte("response")) {
		t.Fatal("expected unknown request to return false")
	}

	if pending.Remove(requestID) {
		t.Fatal("expected unknown request removal to return false")
	}
}
