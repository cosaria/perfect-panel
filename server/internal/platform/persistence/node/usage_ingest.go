package node

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	usageIngestRegistryTable = "schema_registry"
	usageIngestAppliedState  = "applied"
	usageIngestRevisionID    = "0005_async_trust_and_usage"
)

type NodeUsageReport struct {
	ID              int64      `gorm:"primaryKey"`
	ServerID        int64      `gorm:"index;not null"`
	Protocol        string     `gorm:"type:varchar(100);not null;default:''"`
	IdempotencyKey  string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_node_usage_reports_idempotency"`
	AuthStatus      string     `gorm:"type:varchar(50);not null;default:''"`
	ProcessingState string     `gorm:"type:varchar(50);not null;default:''"`
	RawPayload      string     `gorm:"type:text"`
	LogCount        int        `gorm:"not null;default:0"`
	ProcessedAt     *time.Time `gorm:"default:null"`
	CreatedAt       time.Time  `gorm:"<-:create"`
	UpdatedAt       time.Time
}

func (NodeUsageReport) TableName() string {
	return "node_usage_reports"
}

type UsageIngestInput struct {
	ServerID       int64
	Protocol       string
	IdempotencyKey string
	AuthStatus     string
	RawPayload     string
	LogCount       int
}

type UsageIngestResult struct {
	Accepted        bool
	ReportID        int64
	ProcessingState string
}

type UsageIngestRepository struct {
	db *gorm.DB
}

func NewUsageIngestRepository(db *gorm.DB) *UsageIngestRepository {
	return &UsageIngestRepository{db: db}
}

func (r *UsageIngestRepository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil || !r.revisionApplied(db) {
		return false
	}
	return db.Migrator().HasTable(&NodeUsageReport{})
}

func (r *UsageIngestRepository) Ingest(ctx context.Context, input *UsageIngestInput, tx ...*gorm.DB) (*UsageIngestResult, error) {
	if input == nil {
		return nil, nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil, nil
	}
	report := NodeUsageReport{
		ServerID:        input.ServerID,
		Protocol:        input.Protocol,
		IdempotencyKey:  input.IdempotencyKey,
		AuthStatus:      input.AuthStatus,
		ProcessingState: "received",
		RawPayload:      input.RawPayload,
		LogCount:        input.LogCount,
	}
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idempotency_key"}},
		DoNothing: true,
	}).Create(&report)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &UsageIngestResult{
			Accepted:        true,
			ReportID:        report.ID,
			ProcessingState: report.ProcessingState,
		}, nil
	}
	var existing NodeUsageReport
	if err := db.Where("idempotency_key = ?", input.IdempotencyKey).First(&existing).Error; err != nil {
		return nil, err
	}
	return &UsageIngestResult{
		Accepted:        false,
		ReportID:        existing.ID,
		ProcessingState: existing.ProcessingState,
	}, nil
}

func (r *UsageIngestRepository) MarkProcessed(ctx context.Context, reportID int64, state string, tx ...*gorm.DB) error {
	if reportID == 0 {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	now := time.Now().UTC()
	return db.Model(&NodeUsageReport{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"processing_state": state,
			"processed_at":     &now,
		}).Error
}

func (r *UsageIngestRepository) ClaimPending(ctx context.Context, reportID int64, tx ...*gorm.DB) (*UsageIngestResult, error) {
	if reportID == 0 {
		return nil, nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil, nil
	}
	result := db.Model(&NodeUsageReport{}).
		Where("id = ? AND processing_state = ?", reportID, "received").
		Update("processing_state", "processing")
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &UsageIngestResult{
			Accepted:        true,
			ReportID:        reportID,
			ProcessingState: "processing",
		}, nil
	}
	var existing NodeUsageReport
	if err := db.Where("id = ?", reportID).First(&existing).Error; err != nil {
		return nil, err
	}
	return &UsageIngestResult{
		Accepted:        false,
		ReportID:        existing.ID,
		ProcessingState: existing.ProcessingState,
	}, nil
}

func (r *UsageIngestRepository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].WithContext(ctx)
	}
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx)
}

func (r *UsageIngestRepository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(usageIngestRegistryTable) {
		return false
	}
	var total int64
	if err := db.Table(usageIngestRegistryTable).
		Where("id = ? AND state = ?", usageIngestRevisionID, usageIngestAppliedState).
		Count(&total).Error; err != nil {
		return false
	}
	return total > 0
}
