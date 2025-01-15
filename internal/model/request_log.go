package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type RequestLog struct {
	Id               uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	Model            string                `gorm:"type:varchar(128);comment:模型key"`
	UserId           uint64                `gorm:"index;comment:用户id"`
	Username         string                `gorm:"type:varchar(255);comment:用户名"`
	Key              string                `gorm:"type:varchar(255);comment:api key"`
	Ip               string                `gorm:"size:255;comment:ip"`
	Email            string                `gorm:"size:255;comment:用户邮箱"`
	Status           int8                  `gorm:"default:1;comment:状态"`
	RetryTimes       int                   `gorm:"default:0;comment:重试次数"`
	ChannelNameTrace string                `gorm:"size:255;comment:链式名称"`
	ChannelIdTrace   string                `gorm:"size:255;comment:链式id"`
	CreatedAt        time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt        time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt        soft_delete.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt" `
}
