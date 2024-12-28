package service

import (
	"context"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/array"
)

type RequestLogReq struct {
	Model  string
	Ip     string
	Status int8
	Key    string
}

type RequestLogService interface {
	CreateRequestLog(ctx context.Context, req *RequestLogReq) error
	GetStatisticsData(ctx context.Context, req *apiV1.RequestLogRanking) (*apiV1.UseCallCountResponse, error)
	GetRealTimeData(ctx context.Context, n int) (*apiV1.RequestLogRealTimeResponse, error)
}

func NewRequestLogService(
	s *Service,
	userRepo repository.UserRepository,
	repo repository.RequestLogRepository,
	apiKeyRepo repository.ApiKeyRepository,
) RequestLogService {
	return &requestLogService{
		Service:    s,
		userRepo:   userRepo,
		repo:       repo,
		apiKeyRepo: apiKeyRepo,
	}
}

type requestLogService struct {
	*Service
	userRepo   repository.UserRepository
	repo       repository.RequestLogRepository
	apiKeyRepo repository.ApiKeyRepository
}

func (r *requestLogService) CreateRequestLog(ctx context.Context, req *RequestLogReq) error {
	apiKeyItem, err := r.apiKeyRepo.QueryItemByApiKey(ctx, req.Key)
	if err != nil {
		return err
	}
	userId := apiKeyItem.UserId
	id := r.Sid.GenUint64()
	user, err := r.userRepo.FindOne(ctx, userId)
	if err != nil {
		return err
	}
	reqLog := &model.RequestLog{
		Model:    req.Model,
		UserId:   userId,
		Username: user.Username,
		Email:    user.UserEmail,
		Ip:       req.Ip,
		Status:   req.Status,
	}
	reqLog.Id = id
	return r.repo.InsertOne(ctx, reqLog)
}

func (r *requestLogService) GetStatisticsData(ctx context.Context, req *apiV1.RequestLogRanking) (*apiV1.UseCallCountResponse, error) {
	resp := new(apiV1.UseCallCountResponse)
	query := new(repository.RequestLogStatisticsQuery)
	query.StartTime = req.StartTime
	query.EndTime = req.EndTime

	list, err := r.repo.GetStatistics(ctx, query)
	if err != nil {
		return nil, err
	}
	resp.Data = list
	return resp, nil
}

func (r *requestLogService) GetRealTimeData(ctx context.Context, n int) (*apiV1.RequestLogRealTimeResponse, error) {
	resp := new(apiV1.RequestLogRealTimeResponse)
	items, err := r.repo.GetNewRequestLogs(ctx, n)
	if err != nil {
		return nil, err
	}
	list := array.Map(items, func(item *model.RequestLog) apiV1.RequestLogRealTimeItem {
		return apiV1.RequestLogRealTimeItem{
			UserId:    item.UserId,
			Username:  item.Username,
			Email:     item.Email,
			Status:    item.Status,
			Model:     item.Model,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	})
	resp.Data = list
	return resp, nil
}
