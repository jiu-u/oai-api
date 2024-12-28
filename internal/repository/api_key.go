package repository

import (
	"context"
	"errors"
	"github.com/jiu-u/oai-api/internal/model"
)

type ApiKeyRepository interface {
	InsertOne(ctx context.Context, apiKey *model.ApiKey) error
	DeleteOne(ctx context.Context, id uint64) error
	DeleteKey(ctx context.Context, key string) error
	DeleteKeyByUserId(ctx context.Context, userId uint64) error
	IsExist(ctx context.Context, apiKey string) (bool, error)
	QueryItemByApiKey(ctx context.Context, apiKey string) (*model.ApiKey, error)
	GetUserApiKey(ctx context.Context, userId uint64) (*model.ApiKey, error)
}

func NewApiKeyRepository(r *Repository) ApiKeyRepository {
	return &apiKeyRepo{r}
}

type apiKeyRepo struct {
	*Repository
}

func (r *apiKeyRepo) GetUserApiKey(ctx context.Context, userId uint64) (*model.ApiKey, error) {
	var key model.ApiKey
	err := r.DB(ctx).Where("user_id = ?", userId).First(&key).Error
	return &key, err
}

func (r *apiKeyRepo) QueryItemByApiKey(ctx context.Context, apiKey string) (*model.ApiKey, error) {
	var key model.ApiKey
	err := r.DB(ctx).Where("content = ?", apiKey).First(&key).Error
	return &key, err
}

func (r *apiKeyRepo) DeleteKeyByUserId(ctx context.Context, userId uint64) error {
	return r.DB(ctx).Where("user_id = ?", userId).Delete(&model.ApiKey{}).Error
}

// InsertOne 插入一条 ApiKey 记录
func (r *apiKeyRepo) InsertOne(ctx context.Context, apiKey *model.ApiKey) error {
	if apiKey == nil {
		return errors.New("apiKey is nil")
	}
	return r.DB(ctx).Create(apiKey).Error
}

// DeleteOne 根据 ID 删除一条 ApiKey 记录
func (r *apiKeyRepo) DeleteOne(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("id is invalid")
	}
	return r.DB(ctx).Delete(&model.ApiKey{}, id).Error
}

// DeleteKey 根据 Content 删除一条 ApiKey 记录
func (r *apiKeyRepo) DeleteKey(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	return r.DB(ctx).Where("content = ?", key).Delete(&model.ApiKey{}).Error
}

// IsExist 检查指定的 ApiKey 是否存在
func (r *apiKeyRepo) IsExist(ctx context.Context, apiKey string) (bool, error) {
	if apiKey == "" {
		return false, errors.New("apiKey is empty")
	}
	var count int64
	err := r.DB(ctx).Model(&model.ApiKey{}).Where("content = ?", apiKey).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
