package document

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type customDocumentLogicModel interface {
	QueryDocumentDetail(ctx context.Context, id int64) (*Document, error)
	QueryDocumentList(ctx context.Context, page, size int, tag string, search string) (int64, []*Document, error)
	GetDocumentListByAll(ctx context.Context) (int64, []*Document, error)
}

// NewModel returns a model for the database table.
func NewModel(conn *gorm.DB, c *redis.Client) Model {
	return &customDocumentModel{
		defaultDocumentModel: newDocumentModel(conn, c),
	}
}

// QueryDocumentDetail queries the details of a document.
func (m *customDocumentModel) QueryDocumentDetail(ctx context.Context, id int64) (*Document, error) {
	if m.content.Available() {
		data, err := m.content.FindDocument(ctx, id)
		if err != nil {
			return nil, err
		}
		return contentDocumentToLegacy(data), nil
	}
	var data Document
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Document{}).Preload("Group").Where("id = ?", id).Find(v).Error
	})
	return &data, err
}

// QueryDocumentList queries a list of documents.
func (m *customDocumentModel) QueryDocumentList(ctx context.Context, page, size int, tag string, search string) (int64, []*Document, error) {
	if m.content.Available() {
		total, list, err := m.content.ListDocuments(ctx, page, size, tag, search)
		if err != nil {
			return 0, nil, err
		}
		resp := make([]*Document, 0, len(list))
		for _, item := range list {
			resp = append(resp, contentDocumentToLegacy(item))
		}
		return total, resp, nil
	}
	var data []*Document
	var total int64
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		db := conn.Model(&Document{})
		if tag != "" {
			db = db.Where("FIND_IN_SET(?, tags)", tag)
		}
		if search != "" {
			db = db.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
		}
		return db.Count(&total).Offset((page - 1) * size).Limit(size).Find(v).Error
	})
	return total, data, err
}

// GetDocumentListByAll queries a list of documents.
func (m *customDocumentModel) GetDocumentListByAll(ctx context.Context) (int64, []*Document, error) {
	if m.content.Available() {
		total, list, err := m.content.ListVisibleDocuments(ctx)
		if err != nil {
			return 0, nil, err
		}
		resp := make([]*Document, 0, len(list))
		for _, item := range list {
			resp = append(resp, contentDocumentToLegacy(item))
		}
		return total, resp, nil
	}
	var data []*Document
	var total int64
	show := true
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Document{}).Where("`show` = ?", &show).Count(&total).Find(v).Error
	})
	return total, data, err
}
