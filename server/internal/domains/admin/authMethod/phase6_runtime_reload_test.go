package authMethod

import (
	"context"
	"testing"
)

func TestPhase6UpdateGlobalUsesReloadCapabilities(t *testing.T) {
	t.Parallel()

	var emailReloads, mobileReloads, deviceReloads int

	logic := NewUpdateAuthMethodConfigLogic(context.Background(), Deps{
		ReloadEmail: func() {
			emailReloads++
		},
		ReloadMobile: func() {
			mobileReloads++
		},
		ReloadDevice: func() {
			deviceReloads++
		},
	})

	logic.UpdateGlobal("email")
	logic.UpdateGlobal("mobile")
	logic.UpdateGlobal("device")
	logic.UpdateGlobal("apple")

	if emailReloads != 1 {
		t.Fatalf("expected email reload once, got %d", emailReloads)
	}
	if mobileReloads != 1 {
		t.Fatalf("expected mobile reload once, got %d", mobileReloads)
	}
	if deviceReloads != 1 {
		t.Fatalf("expected device reload once, got %d", deviceReloads)
	}
}
