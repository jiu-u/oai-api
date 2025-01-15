package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type UserAuthProvider struct {
	Id             uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	UserId         uint64                `gorm:"index;not null" json:"userId"`
	Provider       string                `gorm:"type:varchar(50);index;not null" json:"provider"`
	ProviderUserId string                `gorm:"type:varchar(255);index;not null" json:"providerUserId"`
	AccessToken    string                `gorm:"type:varchar(255)" json:"accessToken"`
	RefreshToken   string                `gorm:"type:varchar(255)" json:"refreshToken"`
	TokenExpireAt  time.Time             `gorm:"not null" json:"tokenExpireAt"`
	Scope          string                `gorm:"type:varchar(255)" json:"scope"`
	ProviderEmail  string                `gorm:"type:varchar(128)" json:"providerEmail"`
	ProviderName   string                `gorm:"type:varchar(32)" json:"providerName"`
	ProviderAvatar string                `gorm:"type:varchar(255)" json:"providerAvatar"`
	UniqueProvider string                `gorm:"type:varchar(50);uniqueIndex:idx_user_provider_deleted"`
	UniqueUserId   uint64                `gorm:"uniqueIndex:idx_user_provider_deleted"`
	CreatedAt      time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt      time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt      soft_delete.DeletedAt `gorm:"uniqueIndex:idx_user_provider_deleted;index;comment:删除时间" json:"deletedAt"`

	User User `gorm:"foreignKey:UserId;references:Id" json:"user"`
}
