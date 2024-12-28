package model

type RequestLog struct {
	BaseModel
	Model    string `gorm:"size:100;comment:模型key"`
	UserId   uint64 `gorm:"index;comment:用户id"`
	Username string `gorm:"size:255;comment:用户名"`
	Ip       string `gorm:"size:255;comment:ip"`
	Email    string `gorm:"size:255;comment:用户邮箱"`
	Status   int8   `gorm:"default:1;comment:状态"`
}
