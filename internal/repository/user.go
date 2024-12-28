package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertOne(ctx context.Context, user *model.User) error
	FindOneByLinuxDoId(ctx context.Context, linuxDoId uint64) (*model.User, error)
	FindOne(ctx context.Context, id uint64) (*model.User, error)
	FindOneForUpdate(ctx context.Context, id uint64) (*model.User, error)
	UpdateOne(ctx context.Context, user *model.User) error
	//FindAll(ctx context.Context) ([]*model.User, error)
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepo{repo}
}

type userRepo struct {
	*Repository
}

func (r *userRepo) FindOneForUpdate(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).Set("gorm:query_option", "FOR UPDATE").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}

func (r *userRepo) FindOne(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}

func (r *userRepo) UpdateOne(ctx context.Context, user *model.User) error {
	return r.DB(ctx).Updates(user).Error
}

func (r *userRepo) InsertOne(ctx context.Context, user *model.User) error {
	return r.DB(ctx).Create(user).Error
}

func (r *userRepo) FindOneByLinuxDoId(ctx context.Context, linuxDoId uint64) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).Where("linux_do_id = ?", linuxDoId).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}
