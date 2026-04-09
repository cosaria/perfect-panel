package node_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestUsageIngestDeduplicatesByIdempotencyKey(t *testing.T) {
	t.Parallel()

	db := openUsageIngestDB(t)
	repo := node.NewUsageIngestRepository(db)
	ctx := context.Background()

	first, err := repo.Ingest(ctx, &node.UsageIngestInput{
		ServerID:       11,
		Protocol:       "vless",
		IdempotencyKey: "server:11:vless:hash_1",
		AuthStatus:     "verified",
		RawPayload:     `{"traffic":[{"uid":1,"upload":10,"download":20}]}`,
		LogCount:       1,
	})
	if err != nil {
		t.Fatalf("ingest first usage report: %v", err)
	}
	if !first.Accepted {
		t.Fatalf("expected first usage report to be accepted, got %+v", first)
	}

	second, err := repo.Ingest(ctx, &node.UsageIngestInput{
		ServerID:       11,
		Protocol:       "vless",
		IdempotencyKey: "server:11:vless:hash_1",
		AuthStatus:     "verified",
		RawPayload:     `{"traffic":[{"uid":1,"upload":10,"download":20}]}`,
		LogCount:       1,
	})
	if err != nil {
		t.Fatalf("ingest duplicate usage report: %v", err)
	}
	if second.Accepted {
		t.Fatalf("expected duplicate usage report to be rejected, got %+v", second)
	}

	if err := repo.MarkProcessed(ctx, first.ReportID, "processed"); err != nil {
		t.Fatalf("mark usage report processed: %v", err)
	}

	third, err := repo.Ingest(ctx, &node.UsageIngestInput{
		ServerID:       11,
		Protocol:       "vless",
		IdempotencyKey: "server:11:vless:hash_1",
		AuthStatus:     "verified",
		RawPayload:     `{"traffic":[{"uid":1,"upload":10,"download":20}]}`,
		LogCount:       1,
	})
	if err != nil {
		t.Fatalf("ingest processed duplicate usage report: %v", err)
	}
	if third.Accepted {
		t.Fatalf("expected processed duplicate usage report to remain deduplicated, got %+v", third)
	}
	if third.ProcessingState != "processed" {
		t.Fatalf("expected duplicate usage report to expose processed state, got %+v", third.ProcessingState)
	}
}

func openUsageIngestDB(t *testing.T) *gorm.DB {
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
