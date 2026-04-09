package identity

import "time"

type SecurityEvent struct {
	ID         int64     `gorm:"primaryKey"`
	UserID     *int64    `gorm:"index:idx_security_events_user_id;comment:User ID"`
	EventType  string    `gorm:"type:varchar(100);not null;index:idx_security_events_event_type;comment:Event Type"`
	Identifier string    `gorm:"type:varchar(255);default:'';comment:Identity Identifier"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(255);default:'';comment:IP Address"`
	UserAgent  string    `gorm:"type:text;comment:User Agent"`
	Payload    string    `gorm:"type:text;comment:Event Payload"`
	CreatedAt  time.Time `gorm:"<-:create;comment:Creation Time"`
}

func (SecurityEvent) TableName() string {
	return "security_events"
}
