package repository

import (
	"context"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
)

type ModelRepo interface {
	Insert(ctx context.Context, model *model.Model) error
	FindUsefulModels(ctx context.Context, models []string, statusName string) ([]*ModelSearchResultItem, error)
	FindCheckModels(ctx context.Context, models []string) ([]*ModelSearchResultItem, error)
	ExistsHashId(ctx context.Context, modelId string, ProviderId uint64) (bool, error)
	UpdateOne(ctx context.Context, model *model.Model) error
	UpdateStatus(ctx context.Context, id uint64, statusName string, status int) error
	GetAllModelIds(ctx context.Context) ([]string, error)
}

func NewModelRepo(repo *Repository) ModelRepo {
	return &modelRepo{repo}
}

type modelRepo struct {
	*Repository
}

func (r *modelRepo) FindCheckModels(ctx context.Context, models []string) ([]*ModelSearchResultItem, error) {
	var result []*ModelSearchResultItem
	// 查询所有模型并过滤
	err := r.DB(ctx).Model(&model.Model{}).
		Where("model_key IN (?)  ", models). // 查找有效的模型
		Select("id,model_key, provider_id,weight").
		Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询有效模型时出错: %w", err)
	}

	return result, nil
}

func (r *modelRepo) GetAllModelIds(ctx context.Context) ([]string, error) {
	var modelKeys []string
	// 使用 DISTINCT 进行去重查询
	result := r.DB(ctx).Model(&model.Model{}).
		Select("DISTINCT(model_key)"). // 选择不重复的 ModelKey 字段
		Find(&modelKeys)
	// 检查查询是否有错误
	if result.Error != nil {
		return nil, result.Error
	}
	return modelKeys, nil
}

type ModelSearchResultItem struct {
	Id         uint64
	Model      string `gorm:"column:model_key"`
	ProviderId uint64
	Weight     int
}

func (r *modelRepo) Insert(ctx context.Context, model *model.Model) error {
	// 使用 GORM 插入新数据
	if err := r.DB(ctx).Create(&model).Error; err != nil {
		fmt.Println("插入模型时出错", err)
		// 如果遇到唯一性冲突（ProviderId 和 ModelKey 组合唯一），返回错误
		//if strings.Contains(err.Error(), "unique_index:provider_model_key") {
		//	return fmt.Errorf("模型已存在，ProviderId 和 ModelKey 必须唯一")
		//}
		return fmt.Errorf("插入模型时出错: %w", err)
	}
	return nil
}

func (r *modelRepo) FindUsefulModels(ctx context.Context, models []string, statusName string) ([]*ModelSearchResultItem, error) {
	var result []*ModelSearchResultItem
	condition := fmt.Sprintf("%s = ?", statusName)

	// 动态构建查询条件，根据 statusName 来决定查询哪个状态字段
	//switch statusName {
	//case "ChatStatus":
	//	condition = "chat_status = ?"
	//case "FimStatus":
	//	condition = "fim_status = ?"
	//case "EmbeddingsStatus":
	//	condition = "embeddings_status = ?"
	//case "SpeechStatus":
	//	condition = "speech_status = ?"
	//case "TranscriptionStatus":
	//	condition = "transcription_status = ?"
	//case "TranslationStatus":
	//	condition = "translation_status = ?"
	//case "ImageGenStatus":
	//	condition = "image_gen_status = ?"
	//case "ImageEditStatus":
	//	condition = "image_edit_status = ?"
	//case "ImageVariationStatus":
	//	condition = "image_variation_status = ?"
	//default:
	//	return nil, fmt.Errorf("未知的状态名称: %s", statusName)
	//}

	// 查询所有模型并过滤
	err := r.DB(ctx).Model(&model.Model{}).
		Where("model_key IN (?) AND "+condition, models, 1). // 查找有效的模型
		Select("id,model_key, provider_id,weight").
		Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("查询有效模型时出错: %w", err)
	}

	return result, nil
}

func (r *modelRepo) ExistsHashId(ctx context.Context, modelId string, ProviderId uint64) (bool, error) {
	var count int64

	// 查找是否已存在满足条件的记录
	err := r.DB(ctx).Model(&model.Model{}).
		Where("model_key = ? AND provider_id = ?", modelId, ProviderId).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查模型是否存在时出错: %w", err)
	}

	// 如果 count > 0，表示存在此记录
	return count > 0, nil
}

func (r *modelRepo) UpdateOne(ctx context.Context, model *model.Model) error {
	return r.DB(ctx).Updates(model).Error
}

func (r *modelRepo) UpdateStatus(ctx context.Context, id uint64, statusName string, status int) error {
	return r.DB(ctx).Model(&model.Model{}).Where("id = ?", id).Update(statusName, status).Error
}
