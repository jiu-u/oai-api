package model

import "time"

type UserAuthProvider struct {
	BaseModel
	UserId         uint64    `gorm:"index;not null" json:"userId"`
	Provider       string    `gorm:"type:varchar(50);index;not null" json:"provider"`
	ProviderUserId string    `gorm:"type:varchar(255);index;not null" json:"providerUserId"`
	AccessToken    string    `gorm:"type:varchar(255)" json:"accessToken"`
	RefreshToken   string    `gorm:"type:varchar(255)" json:"refreshToken"`
	TokenExpireAt  time.Time `gorm:"not null" json:"tokenExpireAt"`
	Scope          string    `gorm:"type:varchar(255)" json:"scope"`
	ProviderEmail  string    `gorm:"type:varchar(128)" json:"providerEmail"`
	ProviderName   string    `gorm:"type:varchar(32)" json:"providerName"`
	ProviderAvatar string    `gorm:"type:varchar(255)" json:"providerAvatar"`
	UniqueProvider string    `gorm:"type:varchar(50);uniqueIndex:idx_user_provider"`
	UniqueUserId   uint64    `gorm:"uniqueIndex:idx_user_provider"`
	User           User      `gorm:"foreignKey:UserId;references:Id" json:"user"`
}
