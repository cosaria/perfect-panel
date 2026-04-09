package identity

import (
	"context"

	"gorm.io/gorm"
)

const (
	identitySystemRevisionName = "0002_identity_system"
	schemaRegistryTable        = "schema_registry"
	revisionStateApplied       = "applied"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	if !r.revisionApplied(db) {
		return false
	}
	return db.Migrator().HasTable(&User{}) &&
		db.Migrator().HasTable(&AuthIdentity{}) &&
		db.Migrator().HasTable(&UserDevice{})
}

func (r *Repository) FindUserByID(ctx context.Context, id int64, tx ...*gorm.DB) (*User, error) {
	var data User
	err := r.conn(ctx, tx...).Unscoped().
		Preload("AuthIdentities").
		Preload("UserDevices").
		Where("id = ?", id).
		First(&data).Error
	return &data, err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string, tx ...*gorm.DB) (*User, error) {
	identity, err := r.FindAuthIdentityByOpenID(ctx, "email", email, tx...)
	if err != nil {
		return nil, err
	}
	return r.FindUserByID(ctx, identity.UserID, tx...)
}

func (r *Repository) FindUserByReferCode(ctx context.Context, referCode string, tx ...*gorm.DB) (*User, error) {
	var data User
	err := r.conn(ctx, tx...).Unscoped().
		Preload("AuthIdentities").
		Preload("UserDevices").
		Where("refer_code = ?", referCode).
		First(&data).Error
	return &data, err
}

func (r *Repository) InsertUser(ctx context.Context, data *User, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Create(data).Error
}

func (r *Repository) UpdateUser(ctx context.Context, data *User, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Save(data).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Delete(&User{}, id).Error
}

func (r *Repository) FindAdminUsers(ctx context.Context, tx ...*gorm.DB) ([]*User, error) {
	var data []*User
	err := r.conn(ctx, tx...).Preload("AuthIdentities").Where("is_admin = ?", true).Find(&data).Error
	return data, err
}

func (r *Repository) FindUserAuthIdentities(ctx context.Context, userID int64, tx ...*gorm.DB) ([]*AuthIdentity, error) {
	var data []*AuthIdentity
	err := r.conn(ctx, tx...).Where("user_id = ?", userID).Find(&data).Error
	return data, err
}

func (r *Repository) InsertUserAuthIdentity(ctx context.Context, data *AuthIdentity, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Create(data).Error
}

func (r *Repository) UpdateUserAuthIdentity(ctx context.Context, data *AuthIdentity, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Where("user_id = ? AND auth_type = ?", data.UserID, data.AuthType).Save(data).Error
}

func (r *Repository) DeleteUserAuthIdentity(ctx context.Context, userID int64, platform string, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Where("user_id = ? AND auth_type = ?", userID, platform).Delete(&AuthIdentity{}).Error
}

func (r *Repository) FindAuthIdentityByOpenID(ctx context.Context, method, openID string, tx ...*gorm.DB) (*AuthIdentity, error) {
	var data AuthIdentity
	err := r.conn(ctx, tx...).Where("auth_type = ? AND auth_identifier = ?", method, openID).First(&data).Error
	return &data, err
}

func (r *Repository) FindAuthIdentityByUserID(ctx context.Context, method string, userID int64, tx ...*gorm.DB) (*AuthIdentity, error) {
	var data AuthIdentity
	err := r.conn(ctx, tx...).Where("auth_type = ? AND user_id = ?", method, userID).First(&data).Error
	return &data, err
}

func (r *Repository) FindAuthIdentityByPlatform(ctx context.Context, userID int64, platform string, tx ...*gorm.DB) (*AuthIdentity, error) {
	var data AuthIdentity
	err := r.conn(ctx, tx...).Where("user_id = ? AND auth_type = ?", userID, platform).First(&data).Error
	return &data, err
}

func (r *Repository) FindUserDevice(ctx context.Context, id int64, tx ...*gorm.DB) (*UserDevice, error) {
	var data UserDevice
	err := r.conn(ctx, tx...).Where("id = ?", id).First(&data).Error
	return &data, err
}

func (r *Repository) FindUserDeviceByIdentifier(ctx context.Context, identifier string, tx ...*gorm.DB) (*UserDevice, error) {
	var data UserDevice
	err := r.conn(ctx, tx...).Where("identifier = ?", identifier).First(&data).Error
	return &data, err
}

func (r *Repository) InsertUserDevice(ctx context.Context, data *UserDevice, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Create(data).Error
}

func (r *Repository) UpdateUserDevice(ctx context.Context, data *UserDevice, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Save(data).Error
}

func (r *Repository) DeleteUserDevice(ctx context.Context, id int64, tx ...*gorm.DB) error {
	return r.conn(ctx, tx...).Delete(&UserDevice{}, id).Error
}

func (r *Repository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		if ctx != nil {
			return tx[0].WithContext(ctx)
		}
		return tx[0]
	}
	if r.db == nil {
		return nil
	}
	if ctx != nil {
		return r.db.WithContext(ctx)
	}
	return r.db
}

func (r *Repository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(schemaRegistryTable) {
		return false
	}

	var count int64
	if err := db.Table(schemaRegistryTable).
		Where("id = ? AND state = ?", identitySystemRevisionName, revisionStateApplied).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
