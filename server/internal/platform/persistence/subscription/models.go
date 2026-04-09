package subscription

import "time"

type Subscription struct {
	ID          int64     `gorm:"primaryKey"`
	UserID      int64     `gorm:"index;not null;default:0"`
	OrderID     int64     `gorm:"index;not null;default:0"`
	SubscribeID int64     `gorm:"index;not null;default:0"`
	StartTime   time.Time `gorm:"not null"`
	ExpireTime  time.Time `gorm:"not null"`
	FinishedAt  *time.Time
	Traffic     int64     `gorm:"not null;default:0"`
	Download    int64     `gorm:"not null;default:0"`
	Upload      int64     `gorm:"not null;default:0"`
	Status      uint8     `gorm:"type:tinyint(1);not null;default:0"`
	Note        string    `gorm:"type:varchar(500);default:''"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type SubscriptionPeriod struct {
	ID             int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"index;not null"`
	StartTime      time.Time `gorm:"not null"`
	ExpireTime     time.Time `gorm:"not null"`
	FinishedAt     *time.Time
	Status         uint8     `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func (SubscriptionPeriod) TableName() string {
	return "subscription_periods"
}

type SubscriptionToken struct {
	ID             int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"not null;index"`
	Token          string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	UUID           string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsPrimary      bool      `gorm:"not null;default:true"`
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func (SubscriptionToken) TableName() string {
	return "subscription_tokens"
}

type SubscriptionUsageSnapshot struct {
	ID             int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"not null;index"`
	Traffic        int64     `gorm:"not null;default:0"`
	Download       int64     `gorm:"not null;default:0"`
	Upload         int64     `gorm:"not null;default:0"`
	CapturedAt     time.Time `gorm:"not null"`
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func (SubscriptionUsageSnapshot) TableName() string {
	return "subscription_usage_snapshots"
}

type SubscriptionEvent struct {
	ID             int64     `gorm:"primaryKey"`
	SubscriptionID int64     `gorm:"not null;index"`
	EventType      string    `gorm:"type:varchar(100);not null;default:''"`
	Payload        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func (SubscriptionEvent) TableName() string {
	return "subscription_events"
}

type Record struct {
	ID          int64
	UserID      int64
	OrderID     int64
	SubscribeID int64
	StartTime   time.Time
	ExpireTime  time.Time
	FinishedAt  *time.Time
	Traffic     int64
	Download    int64
	Upload      int64
	Token       string
	UUID        string
	Status      uint8
	Note        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
