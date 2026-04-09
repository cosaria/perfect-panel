package identity

import "time"

type VerificationToken struct {
	ID         int64      `gorm:"primaryKey"`
	UserID     *int64     `gorm:"index:idx_verification_tokens_user_id;comment:User ID"`
	AuthType   string     `gorm:"type:varchar(100);not null;index:idx_verification_tokens_auth_type;comment:Auth Type"`
	Identifier string     `gorm:"type:varchar(255);not null;comment:Identity Identifier"`
	Token      string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_verification_tokens_token;comment:Verification Token"`
	Purpose    string     `gorm:"type:varchar(100);default:'';comment:Verification Purpose"`
	Channel    string     `gorm:"type:varchar(100);default:'';comment:Delivery Channel"`
	ExpiresAt  *time.Time `gorm:"comment:Expire Time"`
	ConsumedAt *time.Time `gorm:"comment:Consumed Time"`
	CreatedAt  time.Time  `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt  time.Time  `gorm:"comment:Update Time"`
}

func (VerificationToken) TableName() string {
	return "verification_tokens"
}

type VerificationDelivery struct {
	ID         int64      `gorm:"primaryKey"`
	TokenID    *int64     `gorm:"index:idx_verification_deliveries_token_id;comment:Verification Token ID"`
	AuthType   string     `gorm:"type:varchar(100);not null;comment:Auth Type"`
	Identifier string     `gorm:"type:varchar(255);not null;comment:Identity Identifier"`
	Channel    string     `gorm:"type:varchar(100);default:'';comment:Delivery Channel"`
	Status     string     `gorm:"type:varchar(50);default:'';comment:Delivery Status"`
	Message    string     `gorm:"type:text;comment:Delivery Payload"`
	SentAt     *time.Time `gorm:"comment:Sent Time"`
	CreatedAt  time.Time  `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt  time.Time  `gorm:"comment:Update Time"`
}

func (VerificationDelivery) TableName() string {
	return "verification_deliveries"
}
