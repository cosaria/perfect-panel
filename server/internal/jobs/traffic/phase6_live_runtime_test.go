package traffic

import (
	"context"
	"testing"

	"github.com/perfect-panel/server/models/node"
)

func TestResolveNodeMultiplierManagerFallsBackToLoader(t *testing.T) {
	t.Parallel()

	expected := &node.Manager{}
	loads := 0
	deps := Deps{
		LoadNodeMultiplierManager: func(context.Context) (*node.Manager, error) {
			loads++
			return expected, nil
		},
	}

	got, err := deps.ResolveNodeMultiplierManager(context.Background())
	if err != nil {
		t.Fatalf("expected loader fallback to succeed, got error %v", err)
	}
	if got != expected {
		t.Fatal("expected loader fallback to provide node multiplier manager")
	}
	if loads != 1 {
		t.Fatalf("expected loader to run once, got %d", loads)
	}
}

func TestResolveNodeMultiplierManagerPrefersCurrentManager(t *testing.T) {
	t.Parallel()

	current := &node.Manager{}
	loads := 0
	deps := Deps{
		NodeMultiplierManager: func() *node.Manager {
			return current
		},
		LoadNodeMultiplierManager: func(context.Context) (*node.Manager, error) {
			loads++
			return &node.Manager{}, nil
		},
	}

	got, err := deps.ResolveNodeMultiplierManager(context.Background())
	if err != nil {
		t.Fatalf("expected current manager path to succeed, got error %v", err)
	}
	if got != current {
		t.Fatal("expected current manager to win over fallback loader")
	}
	if loads != 0 {
		t.Fatalf("expected loader to stay unused when manager is already live, got %d", loads)
	}
}
