package service

import (
	"context"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/repository"
)

type AuthService interface {
	NewAccessToken(ctx context.Context, userId uint64) (*apiV1.AccessTokenResponse, error)
	//BanUser(ctx context.Context, userId uint64) error
}

func NewAuthService(s *Service, userRepo repository.UserRepository) AuthService {
	return &authService{
		Service:  s,
		userRepo: userRepo,
	}
}

type authService struct {
	*Service
	userRepo repository.UserRepository
}

func (s *authService) NewAccessToken(ctx context.Context, userId uint64) (*apiV1.AccessTokenResponse, error) {
	user, err := s.userRepo.FindOne(ctx, userId)
	if err != nil {
		return nil, err
	}
	resp := new(apiV1.AccessTokenResponse)
	accessToken, err := s.Jwt.GenAccessToken(user.Id, user.Role)
	if err != nil {
		return nil, err
	}
	resp.AccessToken = accessToken
	return resp, nil
}
