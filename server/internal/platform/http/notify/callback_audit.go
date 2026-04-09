package notify

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
)

func callbackHash(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])
}

func recordExternalTrust(ctx context.Context, deps Deps, event *system.ExternalTrustEvent) {
	if deps.DB == nil || event == nil {
		return
	}
	_ = system.NewExternalTrustRepository(deps.DB).Record(ctx, event)
}

func recordPaymentCallbackAttempt(ctx context.Context, deps Deps, paymentID int64, callbackType string, idempotencyKey string, rawPayload string) (*billing.PaymentCallbackDecision, error) {
	if deps.DB == nil {
		return &billing.PaymentCallbackDecision{
			Accepted:        true,
			CallbackID:      idempotencyKey,
			ProcessingState: "received",
		}, nil
	}
	return billing.NewRepository(deps.DB).RecordPaymentCallbackAttempt(ctx, &billing.PaymentCallbackAttempt{
		PaymentID:       paymentID,
		CallbackType:    callbackType,
		IdempotencyKey:  idempotencyKey,
		RawPayload:      rawPayload,
		ProcessingState: "received",
	})
}

func markPaymentCallbackProcessed(ctx context.Context, deps Deps, callbackID string, state string) {
	if deps.DB == nil || callbackID == "" {
		return
	}
	_ = billing.NewRepository(deps.DB).MarkPaymentCallbackProcessed(ctx, callbackID, state)
}
