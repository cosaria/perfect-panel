package worker

import "github.com/perfect-panel/server/worker/spec"

const (
	ForthwithSendEmail     = spec.ForthwithSendEmail
	EmailTypeVerify        = spec.EmailTypeVerify
	EmailTypeMaintenance   = spec.EmailTypeMaintenance
	EmailTypeExpiration    = spec.EmailTypeExpiration
	EmailTypeTrafficExceed = spec.EmailTypeTrafficExceed
	EmailTypeCustom        = spec.EmailTypeCustom
)

type SendEmailPayload = spec.SendEmailPayload
