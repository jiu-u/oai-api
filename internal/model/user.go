package model

import "time"

type User struct {
	BaseModel
	Username    string  `gorm:"type:varchar(255);unique;index;not null;comment:用户名(唯一，不可为空)" json:"username"`
	Email       *string `gorm:"type:varchar(255);unique;index;comment:邮箱(唯一，可为空)" json:"email"`
	Phone       *string `gorm:"type:varchar(20);unique;index;comment:手机号(唯一，可为空)" json:"phone"`
	Password    string  `gorm:"type:varchar(255);comment:密码(可为空)" json:"password"`
	Avatar      string  `gorm:"type:varchar(255);comment:头像(可为空)" json:"avatar"`
	Role        string  `gorm:"type:varchar(64);default:user;comment:角色(可为空)" json:"role"`
	Status      int8    `gorm:"default:1;index;comment:状态,1启用,2禁用" json:"status"`
	Nickname    string  `gorm:"type:varchar(255);comment:昵称(可为空)" json:"nickname"`
	Level       int
	LastLoginAt time.Time `json:"lastLoginAt"`
	LastLoginIP string    `gorm:"type:varchar(39)" json:"lastLoginIP"`
}
