package model

type TaskType = string

const ()

type AsyncTask struct {
	BaseModel
	UserId  uint64 `gorm:"index;not null" json:"userId"`
	Type    string `gorm:"type:varchar(50);index;not null" json:"type"`
	Content string `gorm:"type:varchar(255);index;not null" json:"content"`
	Status  int8   `gorm:"default:1;index;comment:状态,1未执行，2正在执行，3执行完成，4执行失败" json:"status"`
}
