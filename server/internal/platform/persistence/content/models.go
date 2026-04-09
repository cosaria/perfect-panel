package content

import "time"

type Announcement struct {
	ID        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(255);not null;default:''"`
	Content   string    `gorm:"type:text"`
	Show      *bool     `gorm:"type:tinyint(1);not null;default:0"`
	Pinned    *bool     `gorm:"type:tinyint(1);not null;default:0"`
	Popup     *bool     `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (Announcement) TableName() string {
	return "announcements"
}

type Document struct {
	ID        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(255);not null;default:''"`
	Content   string    `gorm:"type:text"`
	Tags      string    `gorm:"type:varchar(255);not null;default:''"`
	Show      *bool     `gorm:"type:tinyint(1);not null;default:1"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (Document) TableName() string {
	return "documents"
}

type Ticket struct {
	ID          int64     `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null;default:''"`
	Description string    `gorm:"type:text"`
	UserID      int64     `gorm:"type:bigint;not null;default:0;index"`
	Status      uint8     `gorm:"type:tinyint(1);not null;default:1;index"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func (Ticket) TableName() string {
	return "tickets"
}

type TicketMessage struct {
	ID        int64     `gorm:"primaryKey"`
	TicketID  int64     `gorm:"type:bigint;not null;default:0;index"`
	From      string    `gorm:"type:varchar(255);not null;default:''"`
	Type      uint8     `gorm:"type:tinyint(1);not null;default:1"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"<-:create"`
}

func (TicketMessage) TableName() string {
	return "ticket_messages"
}

type AnnouncementFilter struct {
	Show   *bool
	Pinned *bool
	Popup  *bool
	Search string
}

type TicketDetail struct {
	Ticket
	Messages []TicketMessage
}
