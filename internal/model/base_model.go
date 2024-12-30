package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        uint64         `gorm:"primaryKey;autoIncrement:false;comment:主键ID"`
	CreatedAt time.Time      `gorm:"index;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}
