package service

import (
	"context"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/repository"
	"strconv"
)

type UserService interface {
	BanUser(ctx context.Context, userId uint64) error
	GetUserInfo(ctx context.Context, userId uint64) (*apiV1.UserInfo, error)
}

func NewUserService(s *Service, userRepo repository.UserRepository, apikeyRepo repository.ApiKeyRepository) UserService {
	return &userService{
		Service:    s,
		userRepo:   userRepo,
		apikeyRepo: apikeyRepo,
	}
}

type userService struct {
	*Service
	userRepo   repository.UserRepository
	apikeyRepo repository.ApiKeyRepository
}

func (s *userService) GetUserInfo(ctx context.Context, userId uint64) (*apiV1.UserInfo, error) {
	user, err := s.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &apiV1.UserInfo{
		Id:       strconv.FormatUint(user.Id, 10),
		Username: user.Username,
		Email:    *user.Email,
		//LinuxDoId:       strconv.FormatUint(user.LinuxDoId, 10),
		//LinuxDoUsername: user.LinuxDoUsername,
	}, nil
}

func (s *userService) BanUser(ctx context.Context, userId uint64) error {
	return s.Tm.Transaction(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.FindOneForUpdate(ctx, userId)
		if err != nil {
			return err
		}
		user.Status = 2
		err = s.userRepo.UpdateOne(ctx, user)
		if err != nil {
			return err
		}
		return s.apikeyRepo.DeleteKeyByUserId(ctx, userId)
	})
}
