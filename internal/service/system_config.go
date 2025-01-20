package service

import (
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/repository"
	"golang.org/x/net/context"
	"sync/atomic"
)

type SystemConfigService interface {
	SetEmailConfig(ctx context.Context, cfg *apiV1.EmailConfig) error
	GetEmailConfig(ctx context.Context) (*apiV1.EmailConfig, error)
	IsEmailServiceAvailable(ctx context.Context) (bool, error)
	SetLinuxDoOAuthConfig(ctx context.Context, cfg *apiV1.LinuxDoOAuthConfig) error
	GetLinuxDoOAuthConfig(ctx context.Context) (*apiV1.LinuxDoOAuthConfig, error)
	IsLinuxDoOAuthAvailable(ctx context.Context) (bool, error)
	SetGithubOAuthConfig(ctx context.Context, cfg *apiV1.GithubOAuthConfig) error
	GetGithubOAuthConfig(ctx context.Context) (*apiV1.GithubOAuthConfig, error)
	IsGithubOAuthAvailable(ctx context.Context) (bool, error)
	SetRegisterConfig(ctx context.Context, cfg *dto.RegisterConfig) error
	GetRegisterConfig(ctx context.Context) (*dto.RegisterConfig, error)
	InitSystemConfig(ctx context.Context) error
	GetModelConfig(ctx context.Context) (*dto.ModelConfig, error)
	SetModelConfig(ctx context.Context, cfg *dto.ModelConfig) error
}

func NewSystemConfigService(s *Service, repo repository.SystemRepository) SystemConfigService {
	return &systemConfigService{
		Service:           s,
		repo:              repo,
		initSystemCfgFlag: 0,
	}
}

type systemConfigService struct {
	*Service
	repo              repository.SystemRepository
	initSystemCfgFlag int32
}

func (s *systemConfigService) SetEmailConfig(ctx context.Context, cfg *apiV1.EmailConfig) error {

	cfg.Id = s.Sid.GenUint64()
	err := s.Tm.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.SetEmailConfig(ctx, cfg)
		return err
	})
	return err
}

func (s *systemConfigService) GetEmailConfig(ctx context.Context) (*apiV1.EmailConfig, error) {
	resp, err := s.repo.GetEmailConfig(ctx)
	return resp, err
}

func (s *systemConfigService) IsEmailServiceAvailable(ctx context.Context) (bool, error) {
	resp, err := s.repo.IsEmailServiceAvailable(ctx)
	return resp, err
}

func (s *systemConfigService) SetLinuxDoOAuthConfig(ctx context.Context, cfg *apiV1.LinuxDoOAuthConfig) error {
	id := s.Sid.GenUint64()
	cfg.Id = id
	err := s.repo.SetLinuxDoOAuthConfig(ctx, cfg)
	return err
}

func (s *systemConfigService) GetLinuxDoOAuthConfig(ctx context.Context) (*apiV1.LinuxDoOAuthConfig, error) {
	return s.repo.GetLinuxDoOAuthConfig(ctx)
}

func (s *systemConfigService) IsLinuxDoOAuthAvailable(ctx context.Context) (bool, error) {
	return s.repo.IsLinuxDoOAuthAvailable(ctx)
}

func (s *systemConfigService) SetGithubOAuthConfig(ctx context.Context, cfg *apiV1.GithubOAuthConfig) error {
	cfg.Id = s.Sid.GenUint64()
	err := s.repo.SetGithubOAuthConfig(ctx, cfg)
	return err
}

func (s *systemConfigService) GetGithubOAuthConfig(ctx context.Context) (*apiV1.GithubOAuthConfig, error) {
	return s.repo.GetGithubOAuthConfig(ctx)
}

func (s *systemConfigService) IsGithubOAuthAvailable(ctx context.Context) (bool, error) {
	return s.repo.IsGithubOAuthAvailable(ctx)
}

func (s *systemConfigService) SetRegisterConfig(ctx context.Context, cfg *dto.RegisterConfig) error {
	var err error
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.SetRegisterConfig(ctx, cfg)
		return err
	})
	return err
}

func (s *systemConfigService) GetRegisterConfig(ctx context.Context) (*dto.RegisterConfig, error) {
	err := s.InitSystemConfig(ctx)
	resp, err := s.repo.GetRegisterConfig(ctx)
	return resp, err
}

func (s *systemConfigService) InitSystemConfig(ctx context.Context) error {
	newValue := 1
	if atomic.LoadInt32(&s.initSystemCfgFlag) == 1 {
		return nil
	}
	if !atomic.CompareAndSwapInt32(&s.initSystemCfgFlag, 0, int32(newValue)) {
		return nil
	}
	var err error
	cfg, err := s.repo.GetRegisterConfig(ctx)
	if err == nil {
		// 配置已存在
		return nil
	}
	cfg = &dto.RegisterConfig{
		AllowRegister:           true,
		AllowRegisterByPassword: true,
		AllowLoginByPassword:    true,
		AllowEmailValid:         false,
		AllowLinuxDoLogin:       false,
		AllowGithubLogin:        false,
	}
	cfg.Id = s.Sid.GenUint64()
	err = s.repo.SetRegisterConfig(ctx, cfg)
	return err
}

func (s *systemConfigService) GetModelConfig(ctx context.Context) (*dto.ModelConfig, error) {
	resp, err := s.repo.GetModelConfig(ctx)
	return resp, err
}

func (s *systemConfigService) SetModelConfig(ctx context.Context, cfg *dto.ModelConfig) error {
	err := s.Tm.Transaction(ctx, func(ctx context.Context) error {
		cfg.Id = s.Sid.GenUint64()
		err := s.repo.SetModelConfig(ctx, cfg)
		return err
	})
	return err
}
