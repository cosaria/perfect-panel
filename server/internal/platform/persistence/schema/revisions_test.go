package schema_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"gorm.io/gorm"
)

var registerTailRevisionOnce sync.Once
var registerStateTransitionRevisionOnce sync.Once
var errTestStateTransition = errors.New("state transition revision failed")
var testStateTransitionShouldFail bool

type testTailRevision struct{}
type testStateTransitionRevision struct{}

func (testTailRevision) Name() string {
	return schema.RevisionName(2, "tail")
}

func (testTailRevision) Up(*gorm.DB) error {
	return nil
}

func (testStateTransitionRevision) Name() string {
	return schema.RevisionName(3, "state_transition")
}

func (testStateTransitionRevision) Up(*gorm.DB) error {
	if testStateTransitionShouldFail {
		return errTestStateTransition
	}
	return nil
}

func registerTestStateTransitionRevision() {
	registerStateTransitionRevisionOnce.Do(func() {
		schemarevisions.RegisterEmbedded()
		schema.RegisterRevision(testStateTransitionRevision{})
	})
}

func registerTestTailRevision() {
	registerTailRevisionOnce.Do(func() {
		registerTestStateTransitionRevision()
		schemarevisions.RegisterEmbedded()
		schema.RegisterRevision(testTailRevision{})
	})
}

func TestRegisteredRevisionsIncludeBaseline(t *testing.T) {
	schemarevisions.RegisterEmbedded()
	revisions := schema.RegisteredRevisions()
	if len(revisions) == 0 {
		t.Fatal("expected at least one registered revision")
	}
	if revisions[0].Name() != schema.BaselineRevisionName {
		t.Fatalf("expected baseline revision first, got %s", revisions[0].Name())
	}
}

func TestApplyRevisionsSkipsAlreadyAppliedBaseline(t *testing.T) {
	registerTestTailRevision()
	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap db: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	var count int64
	if err := db.Model(&schema.Registry{}).Count(&count).Error; err != nil {
		t.Fatalf("count registry rows: %v", err)
	}
	if count != int64(len(schema.RegisteredRevisions())) {
		t.Fatalf("expected %d registered revisions, got %d", len(schema.RegisteredRevisions()), count)
	}

	var baselineCount int64
	if err := db.Model(&schema.Registry{}).Where("id = ?", schema.BaselineRevisionName).Count(&baselineCount).Error; err != nil {
		t.Fatalf("count baseline registry rows: %v", err)
	}
	if baselineCount != 1 {
		t.Fatalf("expected one baseline revision row, got %d", baselineCount)
	}
}

func TestApplyRevisionsRejectsOutOfOrderHistory(t *testing.T) {
	registerTestTailRevision()
	db := testSchemaDB(t)
	if err := db.AutoMigrate(&schema.Registry{}); err != nil {
		t.Fatalf("create schema registry table: %v", err)
	}
	if err := db.Create(&schema.Registry{
		ID:        schema.RevisionName(2, "tail"),
		Source:    schema.DefaultRevisionSource,
		AppliedAt: time.Now().UTC(),
	}).Error; err != nil {
		t.Fatalf("insert out-of-order registry row: %v", err)
	}

	err := schema.ApplyRevisions(db, schema.DefaultRevisionSource)
	if !errors.Is(err, schema.ErrRevisionOrderConflict) {
		t.Fatalf("expected out-of-order revision history to fail with ErrRevisionOrderConflict, got %v", err)
	}
}

func TestApplyRevisionsRejectsPendingRevision(t *testing.T) {
	registerTestTailRevision()
	db := testSchemaDB(t)
	if err := db.AutoMigrate(&schema.Registry{}); err != nil {
		t.Fatalf("create schema registry table: %v", err)
	}
	if err := db.Create(&schema.Registry{
		ID:        schema.RevisionName(2, "tail"),
		Source:    schema.DefaultRevisionSource,
		State:     schema.RevisionStatePending,
		AppliedAt: time.Now().UTC(),
	}).Error; err != nil {
		t.Fatalf("insert pending registry row: %v", err)
	}

	err := schema.ApplyRevisions(db, schema.DefaultRevisionSource)
	if !errors.Is(err, schema.ErrRevisionPending) {
		t.Fatalf("expected pending revision history to fail with ErrRevisionPending, got %v", err)
	}
}

func TestResetRejectsMultipleRegisteredRevisions(t *testing.T) {
	registerTestTailRevision()
	db := testSchemaDB(t)
	err := schema.Reset(db, schema.DefaultRevisionSource)
	if !errors.Is(err, schema.ErrResetRequiresManualCleanup) {
		t.Fatalf("expected reset to reject multiple revisions, got %v", err)
	}
}

func TestApplyRevisionsMarksSuccessfulRevisionApplied(t *testing.T) {
	registerTestTailRevision()
	testStateTransitionShouldFail = false
	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap db: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}

	var entry schema.Registry
	if err := db.Where("id = ?", schema.RevisionName(3, "state_transition")).Take(&entry).Error; err != nil {
		t.Fatalf("load state transition registry row: %v", err)
	}
	if entry.State != schema.RevisionStateApplied {
		t.Fatalf("expected state transition revision to be applied, got %s", entry.State)
	}
}

func TestApplyRevisionsLeavesFailedRevisionPending(t *testing.T) {
	registerTestTailRevision()
	testStateTransitionShouldFail = true
	t.Cleanup(func() {
		testStateTransitionShouldFail = false
	})

	db := testSchemaDB(t)
	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap db: %v", err)
	}
	err := schema.ApplyRevisions(db, schema.DefaultRevisionSource)
	if !errors.Is(err, errTestStateTransition) {
		t.Fatalf("expected failing revision error, got %v", err)
	}

	var entry schema.Registry
	if err := db.Where("id = ?", schema.RevisionName(3, "state_transition")).Take(&entry).Error; err != nil {
		t.Fatalf("load failed revision registry row: %v", err)
	}
	if entry.State != schema.RevisionStatePending {
		t.Fatalf("expected failed revision to remain pending, got %s", entry.State)
	}

	testStateTransitionShouldFail = false
	err = schema.ApplyRevisions(db, schema.DefaultRevisionSource)
	if !errors.Is(err, schema.ErrRevisionPending) {
		t.Fatalf("expected rerun to reject pending revision, got %v", err)
	}
}
