package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ChannelModelRepository interface {
	CreateChannelModel(ctx context.Context, channelModel *model.ChannelModel) error
	CreateChannelModelBatch(ctx context.Context, channelModels []*model.ChannelModel) error
	CreateChannelModelIfNotExists(ctx context.Context, channel *model.ChannelModel) error
	ExistsChannelModel(ctx context.Context, channelModel *model.ChannelModel) (*model.ChannelModel, error)

	FindChannelModelById(ctx context.Context, id uint64) (*model.ChannelModel, error)
	FindChannelModelByIdForUpdate(ctx context.Context, id uint64) (*model.ChannelModel, error)
	FindChannelModelByIdForShare(ctx context.Context, id uint64) (*model.ChannelModel, error)
	FindUsefulChannelModels(ctx context.Context, modelIds []string) ([]*model.ChannelModel, error)
	FindCheckChannelModels(ctx context.Context, modelIds []string) ([]*model.ChannelModel, error)
	FindAllChannelModels(ctx context.Context) ([]*model.ChannelModel, error)
	FindAllChannelModelIds(ctx context.Context) ([]string, error)

	InCrChannelModelWeight(ctx context.Context, id uint64) error
	DecrChannelModelWeight(ctx context.Context, id uint64) error
	RestoreChannelModel(ctx context.Context) error
	UpdateChannelModel(ctx context.Context, channelModel *model.ChannelModel) error
	ResetChannelModels(ctx context.Context, channelId uint64, channelModels []*model.ChannelModel) error
	UpdateChannelModelsStatus(ctx context.Context, channelId uint64, status int8) error

	DeleteChannelModelByID(ctx context.Context, id uint64) error
	DeleteChannelModelByChannelId(ctx context.Context, channelId uint64) error
	PermanentlyDeleteChannelModel(ctx context.Context, channelModel *model.ChannelModel) error
}

func NewChannelModelRepository(repo *Repository) ChannelModelRepository {
	return &channelModelRepository{
		Repository:    repo,
		MaxErrorCount: 6,
		MaxWeight:     20,
	}
}

type channelModelRepository struct {
	*Repository
	MaxErrorCount int
	MaxWeight     int
}

func (r *channelModelRepository) PermanentlyDeleteChannelModel(ctx context.Context, channelModel *model.ChannelModel) error {
	if channelModel.Id > 0 {
		return r.DB(ctx).Unscoped().Where("id = ?", channelModel.Id).Delete(&model.ChannelModel{}).Error
	}
	return r.DB(ctx).Unscoped().Where("model_key = ? and channel_id = ?", channelModel.ModelKey, channelModel.ChannelId).Delete(&model.ChannelModel{}).Error
}

func (r *channelModelRepository) ExistsChannelModel(ctx context.Context, channelModel *model.ChannelModel) (*model.ChannelModel, error) {
	var temp model.ChannelModel
	var err error
	err = r.DB(ctx).Model(&model.ChannelModel{}).Where("model_key = ? and channel_id = ?", channelModel.ModelKey, channelModel.ChannelId).First(&temp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching channel model: %w", err)
	}
	return &temp, nil
}
func (r *channelModelRepository) ExistsChannelModels(ctx context.Context, channelModel *model.ChannelModel) (bool, error) {
	var count int64
	err := r.DB(ctx).Model(&model.ChannelModel{}).Where("model_key = ? and channel_id = ?", channelModel.ModelKey, channelModel.ChannelId).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error counting channel model: %w", err)
	}
	return count > 0, nil
}

func (r *channelModelRepository) UpdateChannelModelsStatus(ctx context.Context, channelId uint64, status int8) error {
	if status < 0 || status > 2 {
		return errors.New("invalid status")
	}
	return r.DB(ctx).Model(&model.ChannelModel{}).Where("channel_id = ?", channelId).Update("soft_limit", status).Error
}

func (r *channelModelRepository) FindCheckChannelModels(ctx context.Context, modelIds []string) ([]*model.ChannelModel, error) {
	var list []*model.ChannelModel
	err := r.DB(ctx).Where("model_key in (?) and hard_limit = 1", modelIds).Find(&list).Error
	return list, err
}

func (r *channelModelRepository) RestoreChannelModel(ctx context.Context) error {
	t := time.Now().Add(-2 * time.Hour)
	err := r.DB(ctx).Model(&model.ChannelModel{}).
		Where("last_check_time < ? and soft_limit = 2 and hard_limit = 1", t).
		Updates(map[string]any{
			"error_count":     3,
			"weight":          10,
			"last_check_time": time.Now(),
		}).Error
	return err
}

func (r *channelModelRepository) UpdateChannelModel(ctx context.Context, channelModel *model.ChannelModel) error {
	return r.DB(ctx).Updates(channelModel).Error
}

func (r *channelModelRepository) ResetChannelModels(ctx context.Context, channelId uint64, channelModels []*model.ChannelModel) error {
	// 删除所有的channelModel
	err := r.DB(ctx).Unscoped().Where("channel_id = ?", channelId).Delete(&model.ChannelModel{}).Error
	if err != nil {
		return err
	}
	// 插入新的channelModel
	err = r.DB(ctx).Create(channelModels).Error
	return err
}

func (r *channelModelRepository) DeleteChannelModelByChannelId(ctx context.Context, channelId uint64) error {
	return r.DB(ctx).Where("channel_id = ?", channelId).Delete(&model.ChannelModel{}).Error
}

func (r *channelModelRepository) FindUsefulChannelModels(ctx context.Context, modelIds []string) ([]*model.ChannelModel, error) {
	var list []*model.ChannelModel
	err := r.DB(ctx).Model(&model.ChannelModel{}).Where("model_key in (?) And soft_limit = 1 and hard_limit = 1", modelIds).Find(&list).Error
	return list, err
}

func (r *channelModelRepository) FindAllChannelModels(ctx context.Context) ([]*model.ChannelModel, error) {
	var list []*model.ChannelModel
	err := r.DB(ctx).Model(&model.ChannelModel{}).Find(&list).Error
	return list, err
}

func (r *channelModelRepository) FindAllChannelModelIds(ctx context.Context) ([]string, error) {
	var modelKeys []string
	// 使用 DISTINCT 进行去重查询
	result := r.DB(ctx).Model(&model.ChannelModel{}).
		Select("DISTINCT(model_key)"). // 选择不重复的 ModelKey 字段
		Where("soft_limit = 1 and hard_limit = 1").
		Find(&modelKeys)
	// 检查查询是否有错误
	if result.Error != nil {
		return nil, result.Error
	}
	return modelKeys, nil
}

func (r *channelModelRepository) DeleteChannelModelByID(ctx context.Context, id uint64) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.ChannelModel{}).Error
}

func (r *channelModelRepository) DeleteChannelModelById(ctx context.Context, channelId uint64) error {
	return r.DB(ctx).Where("channel_id = ?", channelId).Delete(&model.ChannelModel{}).Error
}

func (r *channelModelRepository) CreateChannelModel(ctx context.Context, channelModel *model.ChannelModel) error {
	result, err := r.ExistsChannelModel(ctx, channelModel)
	if err != nil {
		return err
	}
	if result != nil {
		return errors.New("channel model already exists")
	}
	// 硬删除
	err = r.PermanentlyDeleteChannelModel(ctx, channelModel)
	if err != nil {
		return err
	}
	return r.DB(ctx).Create(channelModel).Error
}

func (r *channelModelRepository) CreateChannelModelBatch(ctx context.Context, channelModels []*model.ChannelModel) error {
	return r.DB(ctx).CreateInBatches(channelModels, 100).Error
}

func (r *channelModelRepository) CreateChannelModelIfNotExists(ctx context.Context, channelModel *model.ChannelModel) error {
	row := r.DB(ctx).Where(&model.ChannelModel{
		ChannelId: channelModel.ChannelId,
		ModelKey:  channelModel.ModelKey,
	}).First(&model.ChannelModel{})
	if row.Error != nil {
		if errors.Is(row.Error, gorm.ErrRecordNotFound) {
			return r.DB(ctx).Create(channelModel).Error
		}
		return row.Error
	}
	return nil
}

func (r *channelModelRepository) FindChannelModelById(ctx context.Context, id uint64) (*model.ChannelModel, error) {
	var channelModel model.ChannelModel
	err := r.DB(ctx).First(&channelModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel channelModel with Id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel channelModel: %w", err)
	}
	return &channelModel, nil

}

func (r *channelModelRepository) FindChannelModelByIdForUpdate(ctx context.Context, id uint64) (*model.ChannelModel, error) {
	var channelModel model.ChannelModel
	err := r.DB(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&channelModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel channelModel with Id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel channelModel: %w", err)
	}
	return &channelModel, nil
}

func (r *channelModelRepository) FindChannelModelByIdForShare(ctx context.Context, id uint64) (*model.ChannelModel, error) {
	var channelModel model.ChannelModel
	err := r.DB(ctx).Clauses(clause.Locking{Strength: "SHARE"}).First(&channelModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel channelModel with Id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel channelModel: %w", err)
	}
	return &channelModel, nil
}

func (r *channelModelRepository) InCrChannelModelWeight(ctx context.Context, id uint64) error {
	// times++
	// errorCount-- if errorCount > 0
	// weight++ if weight < MaxWeight
	// lastCheckTime
	err := r.DB(ctx).
		Model(&model.ChannelModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_count":     gorm.Expr("total_count + ?", 1),
			"error_count":     gorm.Expr("CASE WHEN error_count > 0 THEN error_count - 1 ELSE error_count END"),
			"weight":          gorm.Expr("CASE WHEN weight < ? THEN weight + 1 ELSE weight END", r.MaxWeight),
			"last_check_time": time.Now(),
			"soft_limit":      1,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *channelModelRepository) DecrChannelModelWeight(ctx context.Context, id uint64) error {
	err := r.DB(ctx).
		Model(&model.ChannelModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_count":     gorm.Expr("total_count + ?", 1),
			"error_count":     gorm.Expr("CASE WHEN error_count < ? THEN error_count + 1 ELSE error_count END", r.MaxErrorCount),
			"weight":          gorm.Expr("CASE WHEN weight > 0 THEN weight - 1 ELSE weight END"),
			"last_check_time": time.Now(),
			"soft_limit":      gorm.Expr("CASE WHEN error_count >= ? THEN 2 ELSE 1 END", r.MaxErrorCount-1),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *channelModelRepository) IncrChannelModelCount(ctx context.Context, id uint64) error {
	err := r.DB(ctx).
		Model(&model.ChannelModel{}).Where("id = ?", id).
		Update("count", gorm.Expr("count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *channelModelRepository) IncrChannelModelErrorCount(ctx context.Context, id uint64) error {
	err := r.DB(ctx).
		Model(&model.ChannelModel{}).Where("id = ?", id).
		Update("error_count", gorm.Expr("error_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}
