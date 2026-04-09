package content_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	modelannouncement "github.com/perfect-panel/server/internal/platform/persistence/announcement"
	modelcontent "github.com/perfect-panel/server/internal/platform/persistence/content"
	modeldocument "github.com/perfect-panel/server/internal/platform/persistence/document"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	modelticket "github.com/perfect-panel/server/internal/platform/persistence/ticket"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestAnnouncementAndDocumentCompatibilityPreferNormalizedContentSchema(t *testing.T) {
	t.Parallel()

	db := openContentCompatDB(t)
	rds := openContentCompatRedis(t)
	ctx := context.Background()

	show := true
	pinned := true
	if err := db.Create(&modelcontent.Announcement{
		ID:      11,
		Title:   "Normalized Announcement",
		Content: "hello",
		Show:    &show,
		Pinned:  &pinned,
	}).Error; err != nil {
		t.Fatalf("create normalized announcement: %v", err)
	}
	if err := db.Create(&modelcontent.Document{
		ID:      21,
		Title:   "Normalized Document",
		Content: "content",
		Tags:    "guide,install",
		Show:    &show,
	}).Error; err != nil {
		t.Fatalf("create normalized document: %v", err)
	}

	announcementModel := modelannouncement.NewModel(db, rds)
	totalAnnouncements, announcements, err := announcementModel.GetAnnouncementListByPage(ctx, 1, 10, modelannouncement.Filter{Show: &show})
	if err != nil {
		t.Fatalf("query announcements: %v", err)
	}
	if totalAnnouncements != 1 || len(announcements) != 1 || announcements[0].Title != "Normalized Announcement" {
		t.Fatalf("expected normalized announcement through compatibility facade, got total=%d list=%+v", totalAnnouncements, announcements)
	}

	documentModel := modeldocument.NewModel(db, rds)
	totalDocuments, documents, err := documentModel.QueryDocumentList(ctx, 1, 10, "guide", "Normalized")
	if err != nil {
		t.Fatalf("query documents: %v", err)
	}
	if totalDocuments != 1 || len(documents) != 1 || documents[0].Title != "Normalized Document" {
		t.Fatalf("expected normalized document through compatibility facade, got total=%d list=%+v", totalDocuments, documents)
	}
}

func TestTicketCompatibilityReadsNormalizedContentSchema(t *testing.T) {
	t.Parallel()

	db := openContentCompatDB(t)
	rds := openContentCompatRedis(t)
	ctx := context.Background()

	if err := db.Create(&modelcontent.Ticket{
		ID:          31,
		Title:       "Need help",
		Description: "normalized ticket",
		UserID:      41,
		Status:      modelticket.Pending,
	}).Error; err != nil {
		t.Fatalf("create normalized ticket: %v", err)
	}
	if err := db.Create(&modelcontent.TicketMessage{
		ID:       32,
		TicketID: 31,
		From:     "user",
		Type:     1,
		Content:  "first reply",
	}).Error; err != nil {
		t.Fatalf("create normalized ticket message: %v", err)
	}

	ticketModel := modelticket.NewModel(db, rds)
	detail, err := ticketModel.QueryTicketDetail(ctx, 31)
	if err != nil {
		t.Fatalf("query ticket detail: %v", err)
	}
	if detail.Id != 31 || len(detail.Follows) != 1 || detail.Follows[0].Content != "first reply" {
		t.Fatalf("expected normalized ticket detail through compatibility facade, got %+v", detail)
	}

	total, tickets, err := ticketModel.QueryTicketList(ctx, 1, 10, 41, nil, "Need")
	if err != nil {
		t.Fatalf("query ticket list: %v", err)
	}
	if total != 1 || len(tickets) != 1 || tickets[0].Id != 31 {
		t.Fatalf("expected normalized ticket list through compatibility facade, got total=%d list=%+v", total, tickets)
	}
}

func openContentCompatDB(t *testing.T) *gorm.DB {
	t.Helper()

	schemarevisions.RegisterEmbedded()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqliteDriver.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap schema: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}
	return db
}

func openContentCompatRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}
