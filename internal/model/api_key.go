package model

type ApiKey struct {
	BaseModel
	UserId  uint64 `gorm:"index;comment:用户id"`
	Content string `gorm:";size:64;index;unique;comment:api key"`
}
