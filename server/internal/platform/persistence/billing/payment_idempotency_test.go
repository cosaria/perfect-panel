package billing_test

import (
	"context"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/billing"
)

func TestPaymentCallbackDeduplicatesByIdempotencyKey(t *testing.T) {
	t.Parallel()

	db := openBillingCompatDB(t)
	repo := billing.NewRepository(db)
	ctx := context.Background()

	first, err := repo.RecordPaymentCallbackAttempt(ctx, &billing.PaymentCallbackAttempt{
		PaymentID:       41,
		CallbackType:    "stripe",
		IdempotencyKey:  "stripe:evt_1",
		RawPayload:      `{"id":"evt_1"}`,
		AuthStatus:      "verified",
		ProcessingState: "received",
	})
	if err != nil {
		t.Fatalf("record first callback attempt: %v", err)
	}
	if !first.Accepted {
		t.Fatalf("expected first callback attempt to be accepted, got %+v", first)
	}

	second, err := repo.RecordPaymentCallbackAttempt(ctx, &billing.PaymentCallbackAttempt{
		PaymentID:       41,
		CallbackType:    "stripe",
		IdempotencyKey:  "stripe:evt_1",
		RawPayload:      `{"id":"evt_1"}`,
		AuthStatus:      "verified",
		ProcessingState: "received",
	})
	if err != nil {
		t.Fatalf("record duplicate callback attempt: %v", err)
	}
	if second.Accepted {
		t.Fatalf("expected duplicate callback attempt to be rejected as duplicate, got %+v", second)
	}

	if err := repo.MarkPaymentCallbackProcessed(ctx, first.CallbackID, "processed"); err != nil {
		t.Fatalf("mark callback processed: %v", err)
	}

	third, err := repo.RecordPaymentCallbackAttempt(ctx, &billing.PaymentCallbackAttempt{
		PaymentID:       41,
		CallbackType:    "stripe",
		IdempotencyKey:  "stripe:evt_1",
		RawPayload:      `{"id":"evt_1"}`,
		AuthStatus:      "verified",
		ProcessingState: "received",
	})
	if err != nil {
		t.Fatalf("record processed duplicate callback attempt: %v", err)
	}
	if third.Accepted {
		t.Fatalf("expected processed duplicate callback attempt to stay deduplicated, got %+v", third)
	}
	if third.ProcessingState != "processed" {
		t.Fatalf("expected duplicate callback attempt to expose processed state, got %+v", third.ProcessingState)
	}
}
