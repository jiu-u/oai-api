package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type BaseModel struct {
	Id        uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	CreatedAt time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt" `
}
