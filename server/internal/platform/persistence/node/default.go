package node

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/persistence/catalog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customServerModel)(nil)

//goland:noinspection GoNameStartsWithPackageName
type (
	Model interface {
		serverModel
		NodeModel
		customCacheLogicModel
		customServerLogicModel
	}
	serverModel interface {
		InsertServer(ctx context.Context, data *Server, tx ...*gorm.DB) error
		FindOneServer(ctx context.Context, id int64) (*Server, error)
		UpdateServer(ctx context.Context, data *Server, tx ...*gorm.DB) error
		DeleteServer(ctx context.Context, id int64, tx ...*gorm.DB) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
		QueryServerList(ctx context.Context, ids []int64) (servers []*Server, err error)
	}

	NodeModel interface {
		InsertNode(ctx context.Context, data *Node, tx ...*gorm.DB) error
		FindOneNode(ctx context.Context, id int64) (*Node, error)
		UpdateNode(ctx context.Context, data *Node, tx ...*gorm.DB) error
		DeleteNode(ctx context.Context, id int64, tx ...*gorm.DB) error
	}

	customServerModel struct {
		*defaultServerModel
	}
	defaultServerModel struct {
		*gorm.DB
		Cache          *redis.Client
		catalogRepo    *catalog.Repository
		assignmentRepo *AssignmentRepository
	}
)

func newServerModel(db *gorm.DB, cache *redis.Client) *defaultServerModel {
	return &defaultServerModel{
		DB:             db,
		Cache:          cache,
		catalogRepo:    catalog.NewRepository(db),
		assignmentRepo: NewAssignmentRepository(db),
	}
}

// NewModel returns a model for the database table.
func NewModel(conn *gorm.DB, cache *redis.Client) Model {
	return &customServerModel{
		defaultServerModel: newServerModel(conn, cache),
	}
}

func (m *defaultServerModel) InsertServer(ctx context.Context, data *Server, tx ...*gorm.DB) error {
	db := m.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.WithContext(ctx).Create(data).Error
}

func (m *defaultServerModel) FindOneServer(ctx context.Context, id int64) (*Server, error) {
	var server Server
	err := m.WithContext(ctx).Model(&Server{}).Where("id = ?", id).First(&server).Error
	return &server, err
}

func (m *defaultServerModel) UpdateServer(ctx context.Context, data *Server, tx ...*gorm.DB) error {
	_, err := m.FindOneServer(ctx, data.Id)
	if err != nil {
		return err
	}

	db := m.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.WithContext(ctx).Where("`id` = ?", data.Id).Save(data).Error

}

func (m *defaultServerModel) DeleteServer(ctx context.Context, id int64, tx ...*gorm.DB) error {
	db := m.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.WithContext(ctx).Where("`id` = ?", id).Delete(&Server{}).Error
}

func (m *defaultServerModel) InsertNode(ctx context.Context, data *Node, tx ...*gorm.DB) error {
	return m.withWriteConn(ctx, tx, func(db *gorm.DB) error {
		if err := db.Create(data).Error; err != nil {
			return err
		}
		if m.catalogRepo.Available(db) {
			if err := m.catalogRepo.SyncNodeTagMemberships(ctx, data.Id, splitCSV(data.Tags), db); err != nil {
				return err
			}
			if err := m.assignmentRepo.RefreshAssignmentsForNode(ctx, data.Id, db); err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *defaultServerModel) FindOneNode(ctx context.Context, id int64) (*Node, error) {
	var node Node
	err := m.WithContext(ctx).Model(&Node{}).Where("id = ?", id).First(&node).Error
	return &node, err
}

func (m *defaultServerModel) UpdateNode(ctx context.Context, data *Node, tx ...*gorm.DB) error {
	_, err := m.FindOneNode(ctx, data.Id)
	if err != nil {
		return err
	}

	return m.withWriteConn(ctx, tx, func(db *gorm.DB) error {
		affectedSubscribeIds, err := m.collectAffectedSubscribeIds(ctx, db, data.Id)
		if err != nil {
			return err
		}
		if err := db.Where("`id` = ?", data.Id).Save(data).Error; err != nil {
			return err
		}
		if m.catalogRepo.Available(db) {
			if err := m.catalogRepo.SyncNodeTagMemberships(ctx, data.Id, splitCSV(data.Tags), db); err != nil {
				return err
			}
			currentSubscribeIds, err := m.collectAffectedSubscribeIds(ctx, db, data.Id)
			if err != nil {
				return err
			}
			if err := m.refreshAssignmentsForSubscribeIds(ctx, db, mergeInt64(affectedSubscribeIds, currentSubscribeIds)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *defaultServerModel) DeleteNode(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return m.withWriteConn(ctx, tx, func(db *gorm.DB) error {
		affectedSubscribeIds, err := m.collectAffectedSubscribeIds(ctx, db, id)
		if err != nil {
			return err
		}
		if err := db.Where("node_id = ?", id).Delete(&SubscriptionNodeAssignment{}).Error; err != nil {
			return err
		}
		if m.catalogRepo.Available(db) {
			if err := m.catalogRepo.DeleteNodeMemberships(ctx, id, db); err != nil {
				return err
			}
		}
		if err := db.Where("`id` = ?", id).Delete(&Node{}).Error; err != nil {
			return err
		}
		return m.refreshAssignmentsForSubscribeIds(ctx, db, affectedSubscribeIds)
	})
}

func (m *defaultServerModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.WithContext(ctx).Transaction(fn)
}

func (m *defaultServerModel) withWriteConn(ctx context.Context, tx []*gorm.DB, fn func(db *gorm.DB) error) error {
	if len(tx) > 0 && tx[0] != nil {
		return fn(tx[0].WithContext(ctx))
	}
	return m.WithContext(ctx).Transaction(fn)
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		result = append(result, part)
	}
	return result
}

func (m *defaultServerModel) collectAffectedSubscribeIds(ctx context.Context, db *gorm.DB, nodeId int64) ([]int64, error) {
	if !m.catalogRepo.Available(db) {
		return nil, nil
	}
	return m.catalogRepo.ListSubscribeIDsForNode(ctx, nodeId, db)
}

func (m *defaultServerModel) refreshAssignmentsForSubscribeIds(ctx context.Context, db *gorm.DB, subscribeIds []int64) error {
	if !m.assignmentRepo.Available(db) {
		return nil
	}
	for _, subscribeId := range subscribeIds {
		if err := m.assignmentRepo.RefreshAssignmentsForSubscribe(ctx, subscribeId, db); err != nil {
			return err
		}
	}
	return nil
}

func mergeInt64(left, right []int64) []int64 {
	seen := make(map[int64]struct{}, len(left)+len(right))
	result := make([]int64, 0, len(left)+len(right))
	for _, item := range left {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	for _, item := range right {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}
