package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type TaskType = string

const ()

type AsyncTask struct {
	Id        uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	UserId    uint64                `gorm:"index;not null" json:"userId"`
	Type      string                `gorm:"type:varchar(50);index;not null" json:"type"`
	Content   string                `gorm:"type:varchar(255);index;not null" json:"content"`
	Status    int8                  `gorm:"default:1;index;comment:状态,1未执行，2正在执行，3执行完成，4执行失败" json:"status"`
	CreatedAt time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt" `
}
