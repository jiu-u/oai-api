package repository

import (
	"context"
	"encoding/json"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/model"
)

type SystemRepository interface {
	SetEmailConfig(ctx context.Context, cfg *dto.EmailConfig) error
	GetEmailConfig(ctx context.Context) (*dto.EmailConfig, error)
	IsEmailServiceAvailable(ctx context.Context) (bool, error)
	SetLinuxDoOAuthConfig(ctx context.Context, cfg *dto.LinuxDoOAuthConfig) error
	GetLinuxDoOAuthConfig(ctx context.Context) (*dto.LinuxDoOAuthConfig, error)
	IsLinuxDoOAuthAvailable(ctx context.Context) (bool, error)
	SetGithubOAuthConfig(ctx context.Context, cfg *dto.GithubOAuthConfig) error
	GetGithubOAuthConfig(ctx context.Context) (*dto.GithubOAuthConfig, error)
	IsGithubOAuthAvailable(ctx context.Context) (bool, error)
	SetModelConfig(ctx context.Context, cfg *dto.ModelConfig) error
	GetModelConfig(ctx context.Context) (*dto.ModelConfig, error)
	SetRegisterConfig(ctx context.Context, cfg *dto.RegisterConfig) error
	GetRegisterConfig(ctx context.Context) (*dto.RegisterConfig, error)
}

func NewSystemRepository(r *Repository) SystemRepository {
	return &systemRepository{r}
}

type systemRepository struct {
	*Repository
}

func (r *systemRepository) SetEmailConfig(ctx context.Context, cfg *dto.EmailConfig) error {
	var err error
	cfg2, err := r.GetEmailConfig(ctx)
	// 已经存在
	if err == nil {
		cfg.Id = cfg2.Id
		err = r.UpdateEmailConfig(ctx, cfg)
		return err
	}

	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	kv := &model.SystemConfig{
		KeyName:     "smtp_service",
		Value:       string(jsonStr),
		ConfigType:  "email",
		Description: "email_smtp服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&model.SystemConfig{}).Create(kv).Error
	return err
}

func (r *systemRepository) UpdateEmailConfig(ctx context.Context, cfg *dto.EmailConfig) error {
	var err error
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	kv := &model.SystemConfig{
		KeyName:     "smtp_service",
		Value:       string(jsonStr),
		ConfigType:  "email",
		Description: "email_smtp服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&kv).Updates(&kv).Error
	return err
}

func (r *systemRepository) GetEmailConfig(ctx context.Context) (*dto.EmailConfig, error) {
	var err error
	var systemConfig model.SystemConfig

	err = r.DB(ctx).Model(&systemConfig).Where("config_type = ? and key_name=?", "email", "smtp_service").First(&systemConfig).Error
	if err != nil {
		return nil, err
	}
	var emailCfg dto.EmailConfig
	err = json.Unmarshal([]byte(systemConfig.Value), &emailCfg)
	return &emailCfg, err
}

func (r *systemRepository) IsEmailServiceAvailable(ctx context.Context) (bool, error) {
	_, err := r.GetEmailConfig(ctx)
	return err == nil, err
}

func (r *systemRepository) SetLinuxDoOAuthConfig(ctx context.Context, cfg *dto.LinuxDoOAuthConfig) error {
	var err error
	cfg2, err := r.GetLinuxDoOAuthConfig(ctx)
	if err == nil {
		cfg.Id = cfg2.Id
		err = r.UpdateLinuxDoOAuthConfig(ctx, cfg)
		return err
	}
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     dto.LinuxDoOAuthType,
		Value:       string(jsonStr),
		ConfigType:  "oauth2",
		Description: "linux do oauth2 服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&model.SystemConfig{}).Create(kv).Error
	return err
}

func (r *systemRepository) UpdateLinuxDoOAuthConfig(ctx context.Context, cfg *dto.LinuxDoOAuthConfig) error {
	var err error
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     dto.LinuxDoOAuthType,
		Value:       string(jsonStr),
		ConfigType:  "oauth2",
		Description: "linux do oauth2 服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&kv).Updates(&kv).Error
	return err
}

func (r *systemRepository) GetLinuxDoOAuthConfig(ctx context.Context) (*dto.LinuxDoOAuthConfig, error) {
	var err error
	var systemConfig model.SystemConfig
	err = r.DB(ctx).Model(&systemConfig).Where("config_type = ? and key_name=?", "oauth2", dto.LinuxDoOAuthType).First(&systemConfig).Error
	if err != nil {
		return nil, err
	}
	var linuxDoCfg dto.LinuxDoOAuthConfig
	err = json.Unmarshal([]byte(systemConfig.Value), &linuxDoCfg)
	return &linuxDoCfg, err
}

func (r *systemRepository) IsLinuxDoOAuthAvailable(ctx context.Context) (bool, error) {
	_, err := r.GetLinuxDoOAuthConfig(ctx)
	return err == nil, err
}

func (r *systemRepository) SetGithubOAuthConfig(ctx context.Context, cfg *dto.GithubOAuthConfig) error {
	var err error
	cfg2, err := r.GetGithubOAuthConfig(ctx)
	if err == nil {
		cfg.Id = cfg2.Id
		err = r.UpdateGithubOAuthConfig(ctx, cfg)
		return err
	}
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     dto.GithubOAuthType,
		Value:       string(jsonStr),
		ConfigType:  "oauth2",
		Description: "github oauth2 服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&model.SystemConfig{}).Create(kv).Error
	return err
}

func (r *systemRepository) UpdateGithubOAuthConfig(ctx context.Context, cfg *dto.GithubOAuthConfig) error {
	var err error
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     dto.GithubOAuthType,
		Value:       string(jsonStr),
		ConfigType:  "oauth2",
		Description: "github oauth2 服务",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&kv).Updates(&kv).Error
	return err
}

func (r *systemRepository) GetGithubOAuthConfig(ctx context.Context) (*dto.GithubOAuthConfig, error) {
	var err error
	var systemConfig model.SystemConfig
	err = r.DB(ctx).Model(&systemConfig).Where("config_type = ? and key_name=?", "oauth2", dto.GithubOAuthType).First(&systemConfig).Error
	if err != nil {
		return nil, err
	}
	var linuxDoCfg dto.GithubOAuthConfig
	err = json.Unmarshal([]byte(systemConfig.Value), &linuxDoCfg)
	return &linuxDoCfg, err
}

func (r *systemRepository) IsGithubOAuthAvailable(ctx context.Context) (bool, error) {
	_, err := r.GetGithubOAuthConfig(ctx)
	return err == nil, err
}

func (r *systemRepository) SetModelConfig(ctx context.Context, cfg *dto.ModelConfig) error {
	var err error
	cfg2, err := r.GetModelConfig(ctx)
	if err == nil {
		cfg.Id = cfg2.Id
		err = r.UpdateModelConfig(ctx, cfg)
		return err
	}
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     "model_config",
		Value:       string(jsonStr),
		ConfigType:  "model",
		Description: "model",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&model.SystemConfig{}).Create(kv).Error
	return err
}

func (r *systemRepository) UpdateModelConfig(ctx context.Context, cfg *dto.ModelConfig) error {
	var err error
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     "model_config",
		Value:       string(jsonStr),
		ConfigType:  "model",
		Description: "model",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&kv).Updates(&kv).Error
	return err
}

func (r *systemRepository) GetModelConfig(ctx context.Context) (*dto.ModelConfig, error) {
	var err error
	var systemConfig model.SystemConfig
	err = r.DB(ctx).Model(&systemConfig).Where("config_type = ? and key_name=?", "model", "model_config").First(&systemConfig).Error
	if err != nil {
		return nil, err
	}
	var modelCfg dto.ModelConfig
	err = json.Unmarshal([]byte(systemConfig.Value), &modelCfg)
	return &modelCfg, err
}

func (r *systemRepository) SetRegisterConfig(ctx context.Context, cfg *dto.RegisterConfig) error {
	var err error
	cfg2, err := r.GetRegisterConfig(ctx)
	if err == nil {
		cfg.Id = cfg2.Id
		err = r.UpdateRegisterConfig(ctx, cfg)
		return err
	}
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     "register_config",
		Value:       string(jsonStr),
		ConfigType:  "register",
		Description: "register",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&model.SystemConfig{}).Create(kv).Error
	return err
}

func (r *systemRepository) UpdateRegisterConfig(ctx context.Context, cfg *dto.RegisterConfig) error {
	var err error
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	kv := &model.SystemConfig{
		KeyName:     "register_config",
		Value:       string(jsonStr),
		ConfigType:  "register",
		Description: "register",
	}
	kv.Id = cfg.Id
	err = r.DB(ctx).Model(&kv).Updates(&kv).Error
	return err
}

func (r *systemRepository) GetRegisterConfig(ctx context.Context) (*dto.RegisterConfig, error) {
	var err error
	var systemConfig model.SystemConfig
	err = r.DB(ctx).Model(&systemConfig).Where("config_type = ? and key_name=?", "register", "register_config").First(&systemConfig).Error
	if err != nil {
		return nil, err
	}
	var registerCfg dto.RegisterConfig
	err = json.Unmarshal([]byte(systemConfig.Value), &registerCfg)
	return &registerCfg, err
}
