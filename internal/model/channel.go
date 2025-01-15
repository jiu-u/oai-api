package model

import (
	"fmt"
	"github.com/jiu-u/oai-api/pkg/encrypte"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Channel struct {
	Id        uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	Name      string                `gorm:"size:100;not null;comment:渠道名称" json:"name"`
	Type      string                `gorm:"size:50;not null;comment:渠道类型"`
	EndPoint  string                `gorm:"size:255;not null;comment:基础URL"`
	Balance   float64               `gorm:"comment:余额"`
	APIKey    string                `gorm:"size:255;comment:访问令牌"`
	HashId    string                `gorm:"size:64;uniqueIndex:idx_channel_hash_id;comment:哈希ID" json:"hashId"`
	Status    int8                  `gorm:"default:1;comment:状态，1启用，2禁用"`
	Models    []ChannelModel        `gorm:"foreignKey:ChannelId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;uniqueIndex:idx_channel_hash_id;comment:删除时间" json:"deletedAt" `
}

func (c *Channel) GenerateHashId() {
	c.HashId = encrypte.Sha256Encode(fmt.Sprintf("%s%s%s", c.Type, c.EndPoint, c.APIKey))
}

func (c *Channel) AfterUpdate(tx *gorm.DB) (err error) {
	err = tx.Exec("update channels c set hash_id = SHA2(CONCAT(c.type,c.end_point,c.end_point),256) where id = ?", c.Id).Error
	return err
}
