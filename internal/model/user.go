package model

type User struct {
	BaseModel
	Username        string `gorm:"size:255;index;comment:用户名"`
	Password        string `gorm:"size:255;comment:用户密码"`
	UserEmail       string `gorm:"size:256;index;comment:用户邮箱"`
	Role            string `gorm:"index;default:user;comment:角色"`
	LinuxDoId       uint64 `gorm:"unique;index;comment:linuxDo用户id"`
	LinuxDoUsername string `gorm:"unique;index;size:256;comment:linuxDo用户名"`
	LinuxDoLevel    int    `gorm:"index;comment:linuxDo用户权限等级"`
	Status          int8   `gorm:"default:1;comment:状态,1启用,2禁用"`
}
