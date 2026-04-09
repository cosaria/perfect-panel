package identity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                    int64          `gorm:"primaryKey"`
	Password              string         `gorm:"type:varchar(100);not null;comment:User Password"`
	Algo                  string         `gorm:"type:varchar(20);default:'default';comment:Encryption Algorithm"`
	Salt                  string         `gorm:"type:varchar(20);default:null;comment:Password Salt"`
	Avatar                string         `gorm:"type:MEDIUMTEXT;comment:User Avatar"`
	Balance               int64          `gorm:"default:0;comment:User Balance"`
	ReferCode             string         `gorm:"type:varchar(20);default:'';comment:Referral Code"`
	RefererID             int64          `gorm:"index:idx_users_referer_id;comment:Referrer ID"`
	Commission            int64          `gorm:"default:0;comment:Commission"`
	ReferralPercentage    uint8          `gorm:"default:0;comment:Referral Percentage"`
	OnlyFirstPurchase     *bool          `gorm:"default:true;not null;comment:Only First Purchase Referral"`
	GiftAmount            int64          `gorm:"default:0;comment:User Gift Amount"`
	Enable                *bool          `gorm:"default:true;not null;comment:Is Account Enabled"`
	IsAdmin               *bool          `gorm:"default:false;not null;comment:Is Admin"`
	EnableBalanceNotify   *bool          `gorm:"default:false;not null;comment:Enable Balance Change Notifications"`
	EnableLoginNotify     *bool          `gorm:"default:false;not null;comment:Enable Login Notifications"`
	EnableSubscribeNotify *bool          `gorm:"default:false;not null;comment:Enable Subscription Notifications"`
	EnableTradeNotify     *bool          `gorm:"default:false;not null;comment:Enable Trade Notifications"`
	Rules                 string         `gorm:"type:TEXT;comment:User Rules"`
	AuthIdentities        []AuthIdentity `gorm:"foreignKey:UserID;references:ID"`
	UserDevices           []UserDevice   `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt             time.Time      `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt             time.Time      `gorm:"comment:Update Time"`
	DeletedAt             gorm.DeletedAt `gorm:"index;comment:Deletion Time"`
}

func (User) TableName() string {
	return "users"
}

type AuthIdentity struct {
	ID             int64     `gorm:"primaryKey"`
	UserID         int64     `gorm:"index:idx_user_auth_identities_user_id;not null;comment:User ID"`
	AuthType       string    `gorm:"type:varchar(255);not null;index:idx_user_auth_identities_type_identifier,priority:1;comment:Auth Type"`
	AuthIdentifier string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_auth_identities_type_identifier,priority:2;comment:Auth Identifier"`
	Verified       bool      `gorm:"default:false;not null;comment:Is Verified"`
	CreatedAt      time.Time `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt      time.Time `gorm:"comment:Update Time"`
}

func (AuthIdentity) TableName() string {
	return "user_auth_identities"
}

type UserSession struct {
	ID         int64      `gorm:"primaryKey"`
	UserID     int64      `gorm:"index:idx_user_sessions_user_id;not null;comment:User ID"`
	SessionID  string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_sessions_session_id;comment:Session Identifier"`
	LoginType  string     `gorm:"type:varchar(100);default:'';not null;comment:Login Type"`
	IPAddress  string     `gorm:"type:varchar(255);default:'';comment:IP Address"`
	UserAgent  string     `gorm:"type:text;comment:User Agent"`
	ExpiresAt  *time.Time `gorm:"comment:Session Expire Time"`
	RevokedAt  *time.Time `gorm:"comment:Session Revoked Time"`
	LastSeenAt *time.Time `gorm:"comment:Last Seen Time"`
	CreatedAt  time.Time  `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt  time.Time  `gorm:"comment:Update Time"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

type UserDevice struct {
	ID         int64     `gorm:"primaryKey"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(255);not null;comment:Device IP"`
	UserID     int64     `gorm:"index:idx_user_devices_user_id;not null;comment:User ID"`
	UserAgent  string    `gorm:"default:null;comment:User Agent"`
	Identifier string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_devices_identifier;comment:Device Identifier"`
	Online     bool      `gorm:"default:false;not null;comment:Online"`
	Enabled    bool      `gorm:"default:true;not null;comment:Enabled"`
	CreatedAt  time.Time `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt  time.Time `gorm:"comment:Update Time"`
}

func (UserDevice) TableName() string {
	return "user_devices"
}
