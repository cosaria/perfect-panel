package system

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ExternalTrustEvent struct {
	ID              int64      `gorm:"primaryKey"`
	EntryPoint      string     `gorm:"type:varchar(100);not null;index:idx_external_trust_entrypoint"`
	Credential      string     `gorm:"type:varchar(255);default:''"`
	IdempotencyKey  string     `gorm:"type:varchar(255);default:'';index:idx_external_trust_idempotency"`
	AuthStatus      string     `gorm:"type:varchar(50);not null;default:''"`
	ProcessingState string     `gorm:"type:varchar(50);not null;default:''"`
	FailureReason   string     `gorm:"type:text;default:''"`
	RawPayload      string     `gorm:"type:text;default:''"`
	ProcessedAt     *time.Time `gorm:"default:null"`
	CreatedAt       time.Time  `gorm:"<-:create"`
	UpdatedAt       time.Time
}

func (ExternalTrustEvent) TableName() string {
	return "external_trust_events"
}

type ExternalTrustRepository struct {
	db *gorm.DB
}

func NewExternalTrustRepository(db *gorm.DB) *ExternalTrustRepository {
	return &ExternalTrustRepository{db: db}
}

func (r *ExternalTrustRepository) Record(ctx context.Context, data *ExternalTrustEvent, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Create(data).Error
}

func (r *ExternalTrustRepository) MarkProcessed(ctx context.Context, id int64, state string, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil || id == 0 {
		return nil
	}
	now := time.Now().UTC()
	return db.Model(&ExternalTrustEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processing_state": state,
			"processed_at":     &now,
		}).Error
}

func (r *ExternalTrustRepository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].WithContext(ctx)
	}
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx)
}
