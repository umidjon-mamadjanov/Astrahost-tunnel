package engine

import (
	"testing"

	"github.com/astrahost/astrahost-tunnel/internal/core/dispatcher"
)

func TestNewEngine(t *testing.T) {

	d := dispatcher.New()

	e := New(d)

	if e == nil {
		t.Fatal("engine is nil")
	}

	if e.dispatcher == nil {
		t.Fatal("dispatcher is nil")
	}
}
