package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// SystemConfig 系统配置表结构体
type SystemConfig struct {
	Id          uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	KeyName     string                `gorm:"type:varchar(255);not null;uniqueIndex:idx_type_key" json:"key"`
	Value       string                `gorm:"type:text;not null" json:"value"`
	ConfigType  string                `gorm:"type:varchar(50);not null;uniqueIndex:idx_type_key" json:"configType"`
	Description string                `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt   soft_delete.DeletedAt `gorm:"index;uniqueIndex:idx_type_key;comment:删除时间" json:"deletedAt" `
}
