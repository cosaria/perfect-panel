package worker

import "github.com/perfect-panel/server/internal/jobs/spec"

const (
	ForthwithSendEmail     = spec.ForthwithSendEmail
	EmailTypeVerify        = spec.EmailTypeVerify
	EmailTypeMaintenance   = spec.EmailTypeMaintenance
	EmailTypeExpiration    = spec.EmailTypeExpiration
	EmailTypeTrafficExceed = spec.EmailTypeTrafficExceed
	EmailTypeCustom        = spec.EmailTypeCustom
)

type SendEmailPayload = spec.SendEmailPayload
