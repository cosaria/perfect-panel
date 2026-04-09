package content

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	contentRevisionName  = "0006_content_cleanup"
	contentRegistryTable = "schema_registry"
	contentAppliedState  = "applied"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil || !r.revisionApplied(db) {
		return false
	}
	return r.Installed(db)
}

func (r *Repository) Installed(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	return db.Migrator().HasTable(&Announcement{}) &&
		db.Migrator().HasTable(&Document{}) &&
		db.Migrator().HasTable(&Ticket{}) &&
		db.Migrator().HasTable(&TicketMessage{})
}

func (r *Repository) FindAnnouncement(ctx context.Context, id int64, tx ...*gorm.DB) (*Announcement, error) {
	var data Announcement
	err := r.conn(ctx, tx...).Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) ListAnnouncements(ctx context.Context, page, size int, filter AnnouncementFilter, tx ...*gorm.DB) (int64, []*Announcement, error) {
	query := r.conn(ctx, tx...).Model(&Announcement{})
	if filter.Show != nil {
		query = query.Where("`show` = ?", *filter.Show)
	}
	if filter.Pinned != nil {
		query = query.Where("`pinned` = ?", *filter.Pinned)
	}
	if filter.Popup != nil {
		query = query.Where("`popup` = ?", *filter.Popup)
	}
	if filter.Search != "" {
		query = query.Where("`title` LIKE ? OR `content` LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var list []*Announcement
	if err := query.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (r *Repository) UpsertAnnouncement(ctx context.Context, data *Announcement, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	return r.conn(ctx, tx...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(data).Error
}

func (r *Repository) DeleteAnnouncement(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Delete(&Announcement{}, id).Error
}

func (r *Repository) FindDocument(ctx context.Context, id int64, tx ...*gorm.DB) (*Document, error) {
	var data Document
	err := r.conn(ctx, tx...).Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) ListDocuments(ctx context.Context, page, size int, tag string, search string, tx ...*gorm.DB) (int64, []*Document, error) {
	query := r.conn(ctx, tx...).Model(&Document{})
	if tag != "" {
		query = applyDocumentTagFilter(query, tag)
	}
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var list []*Document
	if err := query.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (r *Repository) ListVisibleDocuments(ctx context.Context, tx ...*gorm.DB) (int64, []*Document, error) {
	show := true
	query := r.conn(ctx, tx...).Model(&Document{}).Where("`show` = ?", &show)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var list []*Document
	if err := query.Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (r *Repository) UpsertDocument(ctx context.Context, data *Document, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	return r.conn(ctx, tx...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(data).Error
}

func (r *Repository) DeleteDocument(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Delete(&Document{}, id).Error
}

func (r *Repository) FindTicket(ctx context.Context, id int64, tx ...*gorm.DB) (*Ticket, error) {
	var data Ticket
	err := r.conn(ctx, tx...).Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Repository) FindTicketDetail(ctx context.Context, id int64, tx ...*gorm.DB) (*TicketDetail, error) {
	var ticket Ticket
	if err := r.conn(ctx, tx...).Where("id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	var messages []TicketMessage
	if err := r.conn(ctx, tx...).Where("ticket_id = ?", id).Order("id asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return &TicketDetail{Ticket: ticket, Messages: messages}, nil
}

func (r *Repository) ListTickets(ctx context.Context, page, size int, userID int64, status *uint8, search string, tx ...*gorm.DB) (int64, []*Ticket, error) {
	query := r.conn(ctx, tx...).Model(&Ticket{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	} else {
		query = query.Where("status != ?", 4)
	}
	if search != "" {
		query = query.Where("title like ? or description like ?", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var list []*Ticket
	if err := query.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (r *Repository) UpsertTicket(ctx context.Context, data *Ticket, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	return r.conn(ctx, tx...).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(data).Error
}

func (r *Repository) DeleteTicket(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Delete(&Ticket{}, id).Error
}

func (r *Repository) InsertTicketMessage(ctx context.Context, data *TicketMessage, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	return r.conn(ctx, tx...).Create(data).Error
}

func (r *Repository) UpdateTicketStatus(ctx context.Context, id, userID int64, status uint8, tx ...*gorm.DB) error {
	query := r.conn(ctx, tx...).Model(&Ticket{}).Where("id = ?", id)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	return query.Update("status", status).Error
}

func (r *Repository) CountPendingTickets(ctx context.Context, tx ...*gorm.DB) (int64, error) {
	var total int64
	err := r.conn(ctx, tx...).Model(&Ticket{}).Where("status = ?", 1).Count(&total).Error
	return total, err
}

func (r *Repository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].WithContext(ctx)
	}
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(contentRegistryTable) {
		return false
	}
	var total int64
	if err := db.Table(contentRegistryTable).
		Where("id = ? AND state = ?", contentRevisionName, contentAppliedState).
		Count(&total).Error; err != nil {
		return false
	}
	return total > 0
}

func applyDocumentTagFilter(query *gorm.DB, tag string) *gorm.DB {
	if query == nil || tag == "" {
		return query
	}
	if query.Dialector != nil && query.Dialector.Name() == "sqlite" {
		return query.Where("instr(',' || tags || ',', ',' || ? || ',') > 0", tag)
	}
	return query.Where("FIND_IN_SET(?, tags)", tag)
}
