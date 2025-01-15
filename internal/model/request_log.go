package model

type RequestLog struct {
	BaseModel
	Model            string `gorm:"type:varchar(128);comment:模型key"`
	UserId           uint64 `gorm:"index;comment:用户id"`
	Username         string `gorm:"type:varchar(255);comment:用户名"`
	Key              string `gorm:"type:varchar(255);comment:api key"`
	Ip               string `gorm:"size:255;comment:ip"`
	Email            string `gorm:"size:255;comment:用户邮箱"`
	Status           int8   `gorm:"default:1;comment:状态"`
	RetryTimes       int    `gorm:"default:0;comment:重试次数"`
	ChannelNameTrace string `gorm:"size:255;comment:链式名称"`
	ChannelIdTrace   string `gorm:"size:255;comment:链式id"`
}
