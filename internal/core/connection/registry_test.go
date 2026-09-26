package connection

import "testing"

func TestRegistryAllocateSubdomain(t *testing.T) {
	registry := NewRegistry()

	session1 := &Session{ID: "session-1"}
	session2 := &Session{ID: "session-2"}
	session3 := &Session{ID: "session-3"}

	subdomain1 := registry.AllocateSubdomain("custom", session1)
	subdomain2 := registry.AllocateSubdomain("custom", session2)
	subdomain3 := registry.AllocateSubdomain("custom", session3)

	if subdomain1 != "custom" {
		t.Fatalf("expected custom, got %s", subdomain1)
	}

	if subdomain2 != "custom01" {
		t.Fatalf("expected custom01, got %s", subdomain2)
	}

	if subdomain3 != "custom02" {
		t.Fatalf("expected custom02, got %s", subdomain3)
	}

	if got, ok := registry.GetBySubdomain("custom"); !ok || got != session1 {
		t.Fatal("custom is not registered correctly")
	}

	if got, ok := registry.GetBySubdomain("custom01"); !ok || got != session2 {
		t.Fatal("custom01 is not registered correctly")
	}

	if got, ok := registry.GetBySubdomain("custom02"); !ok || got != session3 {
		t.Fatal("custom02 is not registered correctly")
	}
}

func TestRegistryUnregisterReleasesSubdomain(t *testing.T) {
	registry := NewRegistry()

	session1 := &Session{ID: "session-1"}
	tunnelID := "tunnel-1"

	subdomain := registry.AllocateSubdomain("custom", session1)
	session1.SetSubdomain(subdomain)

	registry.Register(tunnelID, session1)

	registry.Unregister(tunnelID)

	if _, ok := registry.Get(tunnelID); ok {
		t.Fatal("tunnel should be unregistered")
	}

	if _, ok := registry.GetBySubdomain("custom"); ok {
		t.Fatal("subdomain should be released")
	}

	session2 := &Session{ID: "session-2"}

	got := registry.AllocateSubdomain("custom", session2)

	if got != "custom" {
		t.Fatalf("expected custom to be reusable, got %s", got)
	}
}
