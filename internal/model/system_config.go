package model

// SystemConfig 系统配置表结构体
type SystemConfig struct {
	BaseModel
	KeyName     string `gorm:"type:varchar(255);not null;uniqueIndex:idx_type_key" json:"key"`
	Value       string `gorm:"type:text;not null" json:"value"`
	ConfigType  string `gorm:"type:varchar(50);not null;uniqueIndex:idx_type_key" json:"configType"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}
