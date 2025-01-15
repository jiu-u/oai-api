package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type ApiKey struct {
	Id        uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	UserId    uint64                `gorm:"index;comment:用户id" json:"userId"`
	Content   string                `gorm:"size:64;uniqueIndex:idx_api_key_content;comment:api key" json:"content"`
	CreatedAt time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;uniqueIndex:idx_api_key_content;comment:删除时间" json:"deletedAt" `
}
