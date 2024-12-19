package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        uint64         `gorm:"primaryKey;autoIncrement:false;comment:主键ID"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// 定义常量表示状态
const (
	StatusEnabled  = 1  // 启用
	StatusDisabled = 0  // 禁用
	StatusUnknown  = -1 // 未初始化状态
)

// Model 结构体
type Model struct {
	BaseModel
	ProviderId           uint64 `gorm:"unique_index:provider_model_key;comment:提供者ID"`
	ModelKey             string `gorm:"unique_index:provider_model_key;size:100;comment:模型key"`
	ChatStatus           int8   `gorm:"default:1;index;comment:聊天是否可用"`      // 1: 启用，0: 禁用，-1: 未初始化
	FimStatus            int8   `gorm:"default:1;index;comment:FIM自动补全是否可用"` // 1: 启用，0: 禁用，-1: 未初始化
	EmbeddingsStatus     int8   `gorm:"default:1;index;comment:嵌入是否可用"`      // 1: 启用，0: 禁用，-1: 未初始化
	SpeechStatus         int8   `gorm:"default:1;index;comment:语音是否可用"`      // 1: 启用，0: 禁用，-1: 未初始化
	TranscriptionStatus  int8   `gorm:"default:1;index;comment:转写是否可用"`      // 1: 启用，0: 禁用，-1: 未初始化
	TranslationStatus    int8   `gorm:"default:1;index;comment:翻译是否可用"`      // 1: 启用，0: 禁用，-1: 未初始化
	ImageGenStatus       int8   `gorm:"default:1;index;comment:图片生成是否可用"`    // 1: 启用，0: 禁用，-1: 未初始化
	ImageEditStatus      int8   `gorm:"default:1;index;comment:图片编辑是否可用"`    // 1: 启用，0: 禁用，-1: 未初始化
	ImageVariationStatus int8   `gorm:"default:1;index;comment:图片变换是否可用"`    // 1: 启用，0: 禁用，-1: 未初始化
	Weight               int    `gorm:"default:1;comment:权重"`
	LastCheckTime        uint64 `gorm:"comment:最后一次检查时间"`
}

const (
	ChatStatus           = "chat_status"
	FimStatus            = "fim_status"
	EmbeddingsStatus     = "embeddings_status"
	SpeechStatus         = "speech_status"
	TranscriptionStatus  = "transcription_status"
	TranslationStatus    = "translation_status"
	ImageGenStatus       = "image_gen_status"
	ImageEditStatus      = "image_edit_status"
	ImageVariationStatus = "image_variation_status"
)
