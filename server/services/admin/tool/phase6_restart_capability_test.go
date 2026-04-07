package tool

import (
	"context"
	"testing"
	"time"
)

func TestPhase6RestartSystemUsesExplicitRestartCapability(t *testing.T) {
	t.Parallel()

	restarted := make(chan struct{}, 1)
	logic := NewRestartSystemLogic(context.Background(), Deps{
		Restart: func() error {
			restarted <- struct{}{}
			return nil
		},
	})

	if err := logic.RestartSystem(); err != nil {
		t.Fatalf("restart system: %v", err)
	}

	select {
	case <-restarted:
	case <-time.After(time.Second):
		t.Fatal("expected restart capability to be invoked")
	}
}
