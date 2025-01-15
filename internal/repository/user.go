package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/pkg/vaild"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
	FindUserById(ctx context.Context, id uint64) (*model.User, error)
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

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	if !vaild.IsValidUsername(user.Username) {
		return errors.New("invalid username")
	}
	return r.DB(ctx).Create(user).Error
}

func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepo) FindOneForUpdate(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).Set("gorm:query_option", "FOR UPDATE").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with Id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}

func (r *userRepo) FindUserById(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.DB(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with Id %d not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	return &user, nil
}

func (r *userRepo) UpdateOne(ctx context.Context, user *model.User) error {
	if !vaild.IsValidUsername(user.Username) {
		return errors.New("invalid username")
	}
	return r.DB(ctx).Updates(user).Error
}
