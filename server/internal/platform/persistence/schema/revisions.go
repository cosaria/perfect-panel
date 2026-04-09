package schema

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUnknownRevisionSource = errors.New("unknown revision source")
	ErrRevisionNotFound      = errors.New("revision not found")
	ErrRevisionOrderConflict = errors.New("revision history is out of order")
	ErrRevisionPending       = errors.New("revision is pending manual recovery")
)

type Revision interface {
	Name() string
	Up(*gorm.DB) error
}

const (
	RevisionStatePending = "pending"
	RevisionStateApplied = "applied"
)

var registry = struct {
	mu         sync.RWMutex
	revisions  []Revision
	registered map[string]struct{}
}{
	registered: map[string]struct{}{},
}

func RegisterRevision(revision Revision) {
	if revision == nil {
		return
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	name := revision.Name()
	if name == "" {
		return
	}
	if _, exists := registry.registered[name]; exists {
		return
	}

	registry.revisions = append(registry.revisions, revision)
	registry.registered[name] = struct{}{}
	sort.Slice(registry.revisions, func(i, j int) bool {
		return registry.revisions[i].Name() < registry.revisions[j].Name()
	})
}

func RegisteredRevisions() []Revision {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	result := make([]Revision, len(registry.revisions))
	copy(result, registry.revisions)
	return result
}

func ValidateRevisionSource(source string) error {
	switch NormalizeRevisionSource(source) {
	case DefaultRevisionSource:
		return nil
	default:
		return ErrUnknownRevisionSource
	}
}

func ApplyRevisions(db *gorm.DB, source string) error {
	if err := ValidateRevisionSource(source); err != nil {
		return err
	}
	if err := ensureRegistry(db); err != nil {
		return err
	}
	revisions := RegisteredRevisions()
	appliedState := make([]string, len(revisions))
	missingSeen := false
	for i, revision := range revisions {
		state, err := revisionState(db, revision.Name())
		if err != nil {
			return err
		}
		appliedState[i] = state
		switch state {
		case RevisionStatePending:
			return fmt.Errorf("%w: %s", ErrRevisionPending, revision.Name())
		case RevisionStateApplied:
			if missingSeen {
				return fmt.Errorf("%w: %s", ErrRevisionOrderConflict, revision.Name())
			}
			continue
		case "":
			missingSeen = true
		default:
			return fmt.Errorf("unknown revision state %q for %s", state, revision.Name())
		}
	}
	for i, revision := range revisions {
		if appliedState[i] == RevisionStateApplied {
			continue
		}
		if err := applyRevision(db, revision, source); err != nil {
			return err
		}
	}
	return nil
}

func ensureRegistry(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	return db.AutoMigrate(&Registry{})
}

func appliedRevision(db *gorm.DB, name string) (bool, error) {
	state, err := revisionState(db, name)
	if err != nil {
		return false, err
	}
	return state == RevisionStateApplied, nil
}

func revisionState(db *gorm.DB, name string) (string, error) {
	entry, err := revisionEntry(db, name)
	if err != nil {
		return "", err
	}
	if entry == nil {
		return "", nil
	}
	if entry.State == "" {
		return RevisionStateApplied, nil
	}
	return entry.State, nil
}

func revisionEntry(db *gorm.DB, name string) (*Registry, error) {
	var entry Registry
	err := db.Where("id = ?", name).Take(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func createPendingRevision(db *gorm.DB, name, source string) error {
	entry := Registry{
		ID:        name,
		Source:    NormalizeRevisionSource(source),
		State:     RevisionStatePending,
		AppliedAt: time.Now().UTC(),
	}
	return db.Create(&entry).Error
}

func markRevisionApplied(db *gorm.DB, name string) error {
	return db.Model(&Registry{}).Where("id = ?", name).Updates(map[string]any{
		"state":      RevisionStateApplied,
		"applied_at": time.Now().UTC(),
	}).Error
}

func applyRevision(db *gorm.DB, revision Revision, source string) error {
	if revision == nil {
		return nil
	}
	name := revision.Name()
	if name == "" {
		return nil
	}
	state, err := revisionState(db, name)
	if err != nil {
		return err
	}
	switch state {
	case RevisionStateApplied:
		return nil
	case RevisionStatePending:
		return fmt.Errorf("%w: %s", ErrRevisionPending, name)
	}
	if err := createPendingRevision(db, name, source); err != nil {
		return fmt.Errorf("create pending revision %s: %w", name, err)
	}
	if err := revision.Up(db); err != nil {
		return err
	}
	if err := markRevisionApplied(db, name); err != nil {
		return fmt.Errorf("record revision %s: %w", name, err)
	}
	return nil
}

func findRevision(name string) (Revision, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	for _, revision := range registry.revisions {
		if revision.Name() == name {
			return revision, true
		}
	}
	return nil, false
}
