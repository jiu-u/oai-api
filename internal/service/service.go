package service

import (
	"context"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/cache"
	"github.com/jiu-u/oai-api/pkg/constant"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/log"
	"github.com/jiu-u/oai-api/pkg/sid"
)

type Service struct {
	Sid    *sid.Sid
	Tm     repository.Transaction
	Logger *log.Logger
	Jwt    *jwt.JWT
	Cache  *cache.Cache
}

func NewService(
	sid *sid.Sid,
	tm repository.Transaction,
	logger *log.Logger,
	jwt *jwt.JWT,
	cache *cache.Cache,
) *Service {
	return &Service{
		Sid:    sid,
		Tm:     tm,
		Logger: logger,
		Jwt:    jwt,
		Cache:  cache,
	}
}

func GetClientIp(ctx context.Context) string {
	value := ctx.Value(constant.ClientIPKey)
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}
