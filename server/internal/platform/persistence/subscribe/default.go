package subscribe

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/perfect-panel/server/internal/platform/cache"
	"github.com/perfect-panel/server/internal/platform/persistence/catalog"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customSubscribeModel)(nil)
var (
	cacheSubscribeIdPrefix = "cache:subscribe:id:"
)

type (
	Model interface {
		subscribeModel
		customSubscribeLogicModel
	}
	subscribeModel interface {
		Insert(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error
		FindOne(ctx context.Context, id int64) (*Subscribe, error)
		Update(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error
		Delete(ctx context.Context, id int64, tx ...*gorm.DB) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customSubscribeModel struct {
		*defaultSubscribeModel
	}
	defaultSubscribeModel struct {
		cache.CachedConn
		db             *gorm.DB
		table          string
		catalogRepo    *catalog.Repository
		assignmentRepo *node.AssignmentRepository
	}
)

func newSubscribeModel(db *gorm.DB, c *redis.Client) *defaultSubscribeModel {
	return &defaultSubscribeModel{
		CachedConn:     cache.NewConn(db, c),
		db:             db,
		table:          "`subscribe`",
		catalogRepo:    catalog.NewRepository(db),
		assignmentRepo: node.NewAssignmentRepository(db),
	}
}

//nolint:unused
func (m *defaultSubscribeModel) batchGetCacheKeys(Subscribes ...*Subscribe) []string {
	var keys []string
	for _, subscribe := range Subscribes {
		keys = append(keys, m.getCacheKeys(subscribe)...)
	}
	return keys

}
func (m *defaultSubscribeModel) getCacheKeys(data *Subscribe) []string {
	if data == nil {
		return []string{}
	}
	var keys []string
	if m.catalogRepo.Available() {
		nodeIds, err := m.catalogRepo.ResolveNodeIDsForSubscribe(context.Background(), data.Id)
		if err == nil && len(nodeIds) > 0 {
			var nodes []*node.Node
			err = m.QueryNoCacheCtx(context.Background(), &nodes, func(conn *gorm.DB, v interface{}) error {
				return conn.Model(&node.Node{}).Where("id IN ?", nodeIds).Find(&nodes).Error
			})
			if err == nil {
				for _, n := range nodes {
					keys = append(keys, fmt.Sprintf("%s%d", node.ServerUserListCacheKey, n.ServerId))
				}
			}
		}
		return append(tool.RemoveDuplicateElements(keys...), fmt.Sprintf("%s%v", cacheSubscribeIdPrefix, data.Id))
	}
	if data.Nodes != "" {
		var nodes []*node.Node
		ids := strings.Split(data.Nodes, ",")

		err := m.QueryNoCacheCtx(context.Background(), &nodes, func(conn *gorm.DB, v interface{}) error {
			return conn.Model(&node.Node{}).Where("id IN (?)", tool.StringSliceToInt64Slice(ids)).Find(&nodes).Error
		})
		if err == nil {
			for _, n := range nodes {
				keys = append(keys, fmt.Sprintf("%s%d", node.ServerUserListCacheKey, n.ServerId))
			}
		}
	}
	if data.NodeTags != "" {
		var nodes []*node.Node
		tags := tool.RemoveDuplicateElements(strings.Split(data.NodeTags, ",")...)
		err := m.QueryNoCacheCtx(context.Background(), &nodes, func(conn *gorm.DB, v interface{}) error {
			return conn.Model(&node.Node{}).Scopes(InSet("tags", tags)).Find(&nodes).Error
		})
		if err == nil {
			for _, n := range nodes {
				keys = append(keys, fmt.Sprintf("%s%d", node.ServerUserListCacheKey, n.ServerId))
			}
		}
	}

	return append(tool.RemoveDuplicateElements(keys...), fmt.Sprintf("%s%v", cacheSubscribeIdPrefix, data.Id))
}

func (m *defaultSubscribeModel) Insert(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error {
	if err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		return m.withWriteConn(ctx, conn, tx, func(db *gorm.DB) error {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
			return m.syncRelations(ctx, db, data)
		})
	}); err != nil {
		return err
	}
	return m.invalidateCacheKeys(ctx, m.getCacheKeys(data))
}

func (m *defaultSubscribeModel) FindOne(ctx context.Context, id int64) (*Subscribe, error) {
	SubscribeIdKey := fmt.Sprintf("%s%v", cacheSubscribeIdPrefix, id)
	var resp Subscribe
	err := m.QueryCtx(ctx, &resp, SubscribeIdKey, func(conn *gorm.DB, v interface{}) error {
		if err := conn.Model(&Subscribe{}).Where("`id` = ?", id).First(&resp).Error; err != nil {
			return err
		}
		return m.hydrateSelectors(ctx, conn, &resp)
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultSubscribeModel) Update(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	oldKeys := m.getCacheKeys(old)
	if err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		return m.withWriteConn(ctx, conn, tx, func(db *gorm.DB) error {
			if err := db.Save(data).Error; err != nil {
				return err
			}
			return m.syncRelations(ctx, db, data)
		})
	}); err != nil {
		return err
	}
	return m.invalidateCacheKeys(ctx, append(oldKeys, m.getCacheKeys(data)...))
}

func (m *defaultSubscribeModel) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	oldKeys := m.getCacheKeys(data)
	if err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		return m.withWriteConn(ctx, conn, tx, func(db *gorm.DB) error {
			if err := db.Delete(&Subscribe{}, id).Error; err != nil {
				return err
			}
			if m.catalogRepo.Available(db) {
				if err := m.catalogRepo.DeleteSubscribeSelectors(ctx, id, db); err != nil {
					return err
				}
				if err := m.assignmentRepo.RefreshAssignmentsForSubscribe(ctx, id, db); err != nil {
					return err
				}
			}
			return nil
		})
	}); err != nil {
		return err
	}
	return m.invalidateCacheKeys(ctx, oldKeys)
}

func (m *defaultSubscribeModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.TransactCtx(ctx, fn)
}

func (m *defaultSubscribeModel) withWriteConn(ctx context.Context, conn *gorm.DB, tx []*gorm.DB, fn func(db *gorm.DB) error) error {
	if len(tx) > 0 && tx[0] != nil {
		return fn(tx[0].WithContext(ctx))
	}
	return conn.WithContext(ctx).Transaction(fn)
}

func (m *defaultSubscribeModel) syncRelations(ctx context.Context, db *gorm.DB, data *Subscribe) error {
	if data == nil || !m.catalogRepo.Available(db) {
		return nil
	}
	if err := m.catalogRepo.SyncSubscribeSelectors(ctx, data.Id, tool.StringToInt64Slice(data.Nodes), splitCSV(data.NodeTags), db); err != nil {
		return err
	}
	return m.assignmentRepo.RefreshAssignmentsForSubscribe(ctx, data.Id, db)
}

func (m *defaultSubscribeModel) hydrateSelectors(ctx context.Context, db *gorm.DB, data *Subscribe) error {
	if data == nil || !(data.Nodes == "" && data.NodeTags == "") || !m.catalogRepo.Available(db) {
		return nil
	}
	snapshot, err := m.catalogRepo.LoadSelectorSnapshot(ctx, data.Id, db)
	if err != nil {
		return err
	}
	data.Nodes = tool.Int64SliceToString(snapshot.NodeIds)
	data.NodeTags = strings.Join(snapshot.Tags, ",")
	return nil
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	return tool.RemoveDuplicateElements(strings.Split(raw, ",")...)
}

func (m *defaultSubscribeModel) invalidateCacheKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return m.DelCacheCtx(ctx, tool.RemoveDuplicateElements(keys...)...)
}
