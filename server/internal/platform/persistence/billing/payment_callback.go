package billing

import "context"

func (r *Repository) RecordPaymentCallback(ctx context.Context, data *PaymentCallback) error {
	if data == nil {
		return nil
	}
	return r.conn(ctx).Create(data).Error
}
