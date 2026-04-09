package worker

import "github.com/perfect-panel/server/internal/jobs/spec"

const (
	DeferCloseOrder        = spec.DeferCloseOrder
	ForthwithActivateOrder = spec.ForthwithActivateOrder
)

type DeferCloseOrderPayload = spec.DeferCloseOrderPayload
type ForthwithActivateOrderPayload = spec.ForthwithActivateOrderPayload
