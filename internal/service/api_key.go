package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	v1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"strconv"
	"strings"
)

type ApiKeyService interface {
	CreateApiKey(ctx context.Context, req *v1.CreateApiKeyRequest) (*v1.CreateApiKeyResponse, error)
	ResetApiKey(ctx context.Context, req *v1.ResetApiKeyRequest) (*v1.ResetApiKeyResponse, error)
	IsActiveApiKey(ctx context.Context, key string) bool
	GetUserApiKey(ctx context.Context, userId uint64) (*model.ApiKey, error)
}

func NewApiKeyService(
	s *Service,
	userRepo repository.UserRepository,
	apiKeyRepo repository.ApiKeyRepository,
) ApiKeyService {
	return &apiKeyService{
		Service:    s,
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

type apiKeyService struct {
	*Service
	userRepo   repository.UserRepository
	apiKeyRepo repository.ApiKeyRepository
}

func (s *apiKeyService) GetUserApiKey(ctx context.Context, userId uint64) (*model.ApiKey, error) {
	return s.apiKeyRepo.GetUserApiKey(ctx, userId)
}

func (s *apiKeyService) IsActiveApiKey(ctx context.Context, key string) bool {
	exist, err := s.apiKeyRepo.IsExist(ctx, key)
	return err == nil && exist
}

func (s *apiKeyService) CreateApiKey(ctx context.Context, req *v1.CreateApiKeyRequest) (*v1.CreateApiKeyResponse, error) {
	userId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}
	apiKey := &model.ApiKey{
		UserId:  userId,
		Content: GenerateOpenAIKey(),
	}
	apiKey.Id = s.Sid.GenUint64()
	err = s.apiKeyRepo.InsertOne(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	resp := &v1.CreateApiKeyResponse{
		ApiKey: apiKey.Content,
	}
	return resp, nil
}

func (s *apiKeyService) ResetApiKey(ctx context.Context, req *v1.ResetApiKeyRequest) (*v1.ResetApiKeyResponse, error) {
	userId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return nil, err
	}
	apiKey := new(model.ApiKey)
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.FindOneForUpdate(ctx, userId)
		if err != nil {
			return err
		}
		if user.Status != 1 {
			return errors.New("用户已被禁用")
		}
		err = s.apiKeyRepo.DeleteKeyByUserId(ctx, userId)
		if err != nil {
			return err
		}
		apiKey = &model.ApiKey{
			UserId:  userId,
			Content: GenerateOpenAIKey(),
		}
		apiKey.Id = s.Sid.GenUint64()
		err = s.apiKeyRepo.InsertOne(ctx, apiKey)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp := &v1.ResetApiKeyResponse{
		ApiKey: apiKey.Content,
	}
	return resp, nil
}

// GenerateOpenAIKey 生成一个类似 OpenAI API Content 的随机字符串
func GenerateOpenAIKey() string {
	// 生成 32 字节的随机数据
	randomBytes := make([]byte, 50)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // 处理随机数生成错误
	}

	// 将随机数据编码为 Base64
	key := base64.StdEncoding.EncodeToString(randomBytes)

	// 移除 Base64 中的特殊字符（如 '+' 和 '/'），并截取前 48 个字符
	key = strings.ReplaceAll(key, "+", "")
	key = strings.ReplaceAll(key, "/", "")
	key = key[:48]

	// 添加 "sk-" 前缀
	return "sk-" + key
}
