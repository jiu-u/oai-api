package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type User struct {
	Id          uint64                `gorm:"primaryKey;autoIncrement:false;comment:主键ID" json:"id"`
	Username    string                `gorm:"type:varchar(255);uniqueIndex:idx_user_username_deleted;index;not null;comment:用户名(唯一，不可为空)" json:"username"`
	Email       *string               `gorm:"type:varchar(255);uniqueIndex:idx_user_email_deleted;index;comment:邮箱(唯一，可为空)" json:"email"`
	Phone       *string               `gorm:"type:varchar(20);uniqueIndex:idx_user_phone_deleted;index;comment:手机号(唯一，可为空)" json:"phone"`
	Password    string                `gorm:"type:varchar(255);comment:密码(可为空)" json:"password"`
	Avatar      string                `gorm:"type:varchar(255);comment:头像(可为空)" json:"avatar"`
	Role        string                `gorm:"type:varchar(64);default:user;comment:角色(可为空)" json:"role"`
	Status      int8                  `gorm:"default:1;index;comment:状态,1启用,2禁用" json:"status"`
	Nickname    string                `gorm:"type:varchar(255);comment:昵称(可为空)" json:"nickname"`
	Level       int                   `json:"level"`
	LastLoginAt time.Time             `json:"lastLoginAt"`
	LastLoginIP string                `gorm:"type:varchar(39)" json:"lastLoginIP"`
	CreatedAt   time.Time             `gorm:"index;comment:创建时间" json:"createdAt"`
	UpdatedAt   time.Time             `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt   soft_delete.DeletedAt `gorm:"uniqueIndex:idx_user_username_deleted,idx_user_email_deleted,idx_user_phone_deleted;index;comment:删除时间" json:"deletedAt"`
}
