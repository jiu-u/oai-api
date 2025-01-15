package repository

import (
	"context"
	"errors"
	"github.com/jiu-u/oai-api/internal/model"
	"gorm.io/gorm"
)

type UserAuthProviderRepository interface {
	CreateUserAuthProvider(ctx context.Context, p *model.UserAuthProvider) error
	FindUserAuthProvider(ctx context.Context, provider, providerUserId string) (*model.UserAuthProvider, error)
	UpdateUserAuthProviderById(ctx context.Context, p *model.UserAuthProvider) error
}

func NewUserAuthProviderRepository(r *Repository) UserAuthProviderRepository {
	return &userAuthRepository{
		Repository: r,
	}
}

type userAuthRepository struct {
	*Repository
}

func (r *userAuthRepository) CreateUserAuthProvider(ctx context.Context, p *model.UserAuthProvider) error {
	err := r.DB(ctx).Model(&model.UserAuthProvider{}).Create(p).Error
	return err
}

//func (r *userAuthRepository) FindUserAuthProvider(ctx context.Context, provider, providerUserId string) (*model.UserAuthProvider, error) {
//	var userAuthProvider model.UserAuthProvider
//	err := r.DB(ctx).Model(&model.UserAuthProvider{}).
//		Where("provider = ? AND provider_user_id = ?", provider, providerUserId).
//		First(&userAuthProvider).Error
//	if err != nil {
//		return nil, err
//	}
//	return &userAuthProvider, nil
//}
//
//func (r *userAuthRepository) UpdateUserAuthProviderById(ctx context.Context, p *model.UserAuthProvider) error {
//	err := r.DB(ctx).Model(&model.UserAuthProvider{}).
//		Where("provider = ? AND provider_user_id = ?", p.Provider, p.ProviderUserId).
//		Updates(p).Error
//	return err
//}

func (r *userAuthRepository) FindUserAuthProvider(ctx context.Context, provider, providerUserId string) (*model.UserAuthProvider, error) {
	var userAuthProvider model.UserAuthProvider
	err := r.DB(ctx).Where("provider = ? AND provider_user_id = ?", provider, providerUserId).First(&userAuthProvider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示未找到记录
		}
		return nil, err // 其他错误返回
	}
	return &userAuthProvider, nil
}

func (r *userAuthRepository) UpdateUserAuthProviderById(ctx context.Context, p *model.UserAuthProvider) error {
	err := r.DB(ctx).Model(&model.UserAuthProvider{}).
		//Where("provider = ? AND provider_user_id = ?", p.Provider, p.ProviderUserId).
		Where("id = ?", p.Id).
		Updates(map[string]interface{}{
			"access_token":    p.AccessToken,
			"refresh_token":   p.RefreshToken,
			"token_expire_at": p.TokenExpireAt,
			"scope":           p.Scope,
			"provider_email":  p.ProviderEmail,
			"provider_name":   p.ProviderName,
			"provider_avatar": p.ProviderAvatar,
		}).Error
	return err
}
