package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type ChannelModel struct {
	Id            uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	ChannelId     uint64                `gorm:"uniqueIndex:channel_model_key;comment:渠道ID"`
	ModelKey      string                `gorm:"uniqueIndex:channel_model_key;size:100;comment:模型key"`
	SoftLimit     int8                  `gorm:"default:1;index;comment:软限制,1启用,2禁用"`
	HardLimit     int8                  `gorm:"default:1;index;comment:硬限制,1启用,2禁用"`
	Weight        int                   `gorm:"default:1;comment:权重"`
	LastCheckTime time.Time             `gorm:"comment:最后一次检查时间"`
	ErrorCount    int32                 `gorm:"default:0;comment:错误次数"`
	TotalCount    int64                 `gorm:"default:0;comment:总次数"`
	CreatedAt     time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt     time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt     soft_delete.DeletedAt `gorm:"index;uniqueIndex:channel_model_key;comment:删除时间" json:"deletedAt" `
}
