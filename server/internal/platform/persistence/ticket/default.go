package ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/cache"
	"github.com/perfect-panel/server/internal/platform/persistence/content"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customTicketModel)(nil)
var (
	cacheTicketIdPrefix = "cache:ticket:id:"
)

type (
	Model interface {
		ticketModel
		customTicketLogicModel
	}
	ticketModel interface {
		Insert(ctx context.Context, data *Ticket) error
		FindOne(ctx context.Context, id int64) (*Ticket, error)
		Update(ctx context.Context, data *Ticket) error
		Delete(ctx context.Context, id int64) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customTicketModel struct {
		*defaultTicketModel
	}
	defaultTicketModel struct {
		cache.CachedConn
		table   string
		content *content.Repository
	}
)

func newTicketModel(db *gorm.DB, c *redis.Client) *defaultTicketModel {
	return &defaultTicketModel{
		CachedConn: cache.NewConn(db, c),
		table:      "`ticket`",
		content:    content.NewRepository(db),
	}
}

//nolint:unused
func (m *defaultTicketModel) batchGetCacheKeys(Tickets ...*Ticket) []string {
	var keys []string
	for _, ticket := range Tickets {
		keys = append(keys, m.getCacheKeys(ticket)...)
	}
	return keys

}
func (m *defaultTicketModel) getCacheKeys(data *Ticket) []string {
	if data == nil {
		return []string{}
	}
	ticketIdKey := fmt.Sprintf("%s%v", cacheTicketIdPrefix, data.Id)
	cacheKeys := []string{
		ticketIdKey,
	}
	return cacheKeys
}

func (m *defaultTicketModel) Insert(ctx context.Context, data *Ticket) error {
	if m.content.Available() {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.content.UpsertTicket(ctx, legacyTicketToContent(data))
		}, m.getCacheKeys(data)...)
	}
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		return conn.Create(&data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultTicketModel) FindOne(ctx context.Context, id int64) (*Ticket, error) {
	if m.content.Available() {
		data, err := m.content.FindTicket(ctx, id)
		if err != nil {
			return nil, err
		}
		return contentTicketToLegacy(data), nil
	}
	TicketIdKey := fmt.Sprintf("%s%v", cacheTicketIdPrefix, id)
	var resp Ticket
	err := m.QueryCtx(ctx, &resp, TicketIdKey, func(conn *gorm.DB, v interface{}) error {

		return conn.Model(&Ticket{}).Where("`id` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultTicketModel) Update(ctx context.Context, data *Ticket) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if m.content.Available() {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.content.UpsertTicket(ctx, legacyTicketToContent(data))
		}, m.getCacheKeys(old)...)
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		return db.Save(data).Error
	}, m.getCacheKeys(old)...)
	return err
}

func (m *defaultTicketModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if m.content.Available() {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.content.DeleteTicket(ctx, id)
		}, m.getCacheKeys(data)...)
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		return db.Delete(&Ticket{}, id).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultTicketModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.TransactCtx(ctx, fn)
}

func legacyTicketToContent(data *Ticket) *content.Ticket {
	if data == nil {
		return nil
	}
	return &content.Ticket{
		ID:          data.Id,
		Title:       data.Title,
		Description: data.Description,
		UserID:      data.UserId,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func contentTicketToLegacy(data *content.Ticket) *Ticket {
	if data == nil {
		return nil
	}
	return &Ticket{
		Id:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		UserId:      data.UserID,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func legacyFollowToContent(data *Follow) *content.TicketMessage {
	if data == nil {
		return nil
	}
	return &content.TicketMessage{
		ID:        data.Id,
		TicketID:  data.TicketId,
		From:      data.From,
		Type:      data.Type,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}
}

func contentMessageToLegacy(data content.TicketMessage) Follow {
	return Follow{
		Id:        data.ID,
		TicketId:  data.TicketID,
		From:      data.From,
		Type:      data.Type,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}
}
