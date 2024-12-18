package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        uint64         `gorm:"primaryKey;autoIncrement:false;comment:主键ID"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

type Model struct {
	BaseModel
	ProviderId    uint64 `gorm:"unique_index:provider_model_key;comment:提供者ID"`
	ModelKey      string `gorm:"unique_index:provider_model_key;size:100;comment:模型key"`
	ChatStatus    bool   `gorm:"default:true;comment:聊天是否可用"`
	TTSStatus     bool   `gorm:"default:true;comment:TTS是否可用"`
	TransStatus   bool   `gorm:"default:true;comment:转写是否可用"`
	IMageStatus   bool   `gorm:"default:true;comment:图片是否可用"`
	LastCheckTime int64
}
