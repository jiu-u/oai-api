package model

type Channel struct {
	BaseModel
	Name     string  `json:"name"`
	Type     string  `gorm:"size:50;not null;comment:提供者类型"`
	EndPoint string  `gorm:"size:255;not null;comment:基础URL"`
	Balance  float64 `gorm:"comment:余额"`
	APIKey   string  `gorm:"size:255;comment:访问令牌"`        // 访问令牌，最大长度 255
	HashId   string  `gorm:"size:64;unique;comment:哈希ID"`  // 哈希ID，最大长度 64，唯一
	Status   int8    `gorm:"default:1;comment:状态，1启用，2禁用"` // 状态，默认为 1
}
