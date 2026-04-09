package schema

import "time"

type Registry struct {
	ID        string    `gorm:"primaryKey;type:varchar(64);not null;comment:Revision ID"`
	Source    string    `gorm:"type:varchar(32);not null;default:'embedded';comment:Revision Source"`
	State     string    `gorm:"type:varchar(16);not null;default:'applied';comment:Revision State"`
	Checksum  string    `gorm:"type:varchar(128);default:'';comment:Revision Checksum"`
	AppliedAt time.Time `gorm:"not null;comment:Applied At"`
	CreatedAt time.Time `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt time.Time `gorm:"comment:Update Time"`
}

func (Registry) TableName() string {
	return "schema_registry"
}
