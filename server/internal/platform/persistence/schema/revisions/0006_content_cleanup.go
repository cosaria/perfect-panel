package revisions

import (
	legacyannouncement "github.com/perfect-panel/server/internal/platform/persistence/announcement"
	"github.com/perfect-panel/server/internal/platform/persistence/content"
	legacydocument "github.com/perfect-panel/server/internal/platform/persistence/document"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	legacyticket "github.com/perfect-panel/server/internal/platform/persistence/ticket"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contentCleanupRevision struct{}

func (contentCleanupRevision) Name() string {
	return schema.RevisionName(6, "content_cleanup")
}

func (contentCleanupRevision) Up(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&content.Announcement{},
		&content.Document{},
		&content.Ticket{},
		&content.TicketMessage{},
	); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := backfillAnnouncements(tx); err != nil {
			return err
		}
		if err := backfillDocuments(tx); err != nil {
			return err
		}
		if err := backfillTickets(tx); err != nil {
			return err
		}
		return nil
	})
}

func backfillAnnouncements(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&legacyannouncement.Announcement{}) {
		return nil
	}
	var rows []legacyannouncement.Announcement
	if err := tx.Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		data := content.Announcement{
			ID:        row.Id,
			Title:     row.Title,
			Content:   row.Content,
			Show:      row.Show,
			Pinned:    row.Pinned,
			Popup:     row.Popup,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillDocuments(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&legacydocument.Document{}) {
		return nil
	}
	var rows []legacydocument.Document
	if err := tx.Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		data := content.Document{
			ID:        row.Id,
			Title:     row.Title,
			Content:   row.Content,
			Tags:      row.Tags,
			Show:      row.Show,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillTickets(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&legacyticket.Ticket{}) {
		var tickets []legacyticket.Ticket
		if err := tx.Find(&tickets).Error; err != nil {
			return err
		}
		for _, row := range tickets {
			data := content.Ticket{
				ID:          row.Id,
				Title:       row.Title,
				Description: row.Description,
				UserID:      row.UserId,
				Status:      row.Status,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&data).Error; err != nil {
				return err
			}
		}
	}

	if !tx.Migrator().HasTable(&legacyticket.Follow{}) {
		return nil
	}
	var follows []legacyticket.Follow
	if err := tx.Find(&follows).Error; err != nil {
		return err
	}
	for _, row := range follows {
		data := content.TicketMessage{
			ID:        row.Id,
			TicketID:  row.TicketId,
			From:      row.From,
			Type:      row.Type,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}
