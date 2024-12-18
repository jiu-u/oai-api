package model

type Provider struct {
	BaseModel
	Name    string `gorm:"size:100;not null;comment:提供者名称"`
	Type    string `gorm:"size:50;not null;comment:提供者类型"`
	BaseUrl string `gorm:"size:255;not null;comment:基础URL"`
	Token   string `gorm:"size:255;comment:访问令牌"`     // 访问令牌，最大长度 255
	HashId  string `gorm:"size:64;unique;comment:哈希ID"` // 哈希ID，最大长度 64，唯一
	Status  uint   `gorm:"default:1;comment:状态"`        // 状态，默认为 1
}
