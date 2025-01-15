package service

import (
	"context"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/array"
	"strconv"
)

type RequestLogReq struct {
	Model        string
	Ip           string
	Status       int8
	Key          string
	RetryTimes   int
	ChannelNames string
	ChannelIds   string
}

type RequestLogService interface {
	CreateRequestLog(ctx context.Context, req *RequestLogReq) error
	GetRequestLogs(ctx context.Context, req *apiV1.RequestLogsQuery) (*apiV1.RequestLogsResponse, error)
	GetRequestLogsModelRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) (*apiV1.RequestLogsModelRankingResponse, error)
	GetRequestLogsUserRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) (*apiV1.RequestLogsUserRankingResponse, error)
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
	user, err := r.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return err
	}
	reqLog := &model.RequestLog{
		Model:            req.Model,
		UserId:           userId,
		Username:         user.Username,
		Ip:               req.Ip,
		Status:           req.Status,
		Key:              req.Key,
		RetryTimes:       req.RetryTimes,
		ChannelNameTrace: req.ChannelNames,
		ChannelIdTrace:   req.ChannelIds,
	}
	if user.Email != nil {
		reqLog.Email = *user.Email
	}
	reqLog.Id = id
	return r.repo.CreateRequestLog(ctx, reqLog)
}

func (r *requestLogService) GetRequestLogs(ctx context.Context, req *apiV1.RequestLogsQuery) (*apiV1.RequestLogsResponse, error) {
	list, total, err := r.repo.FindRequestLogs(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := &apiV1.RequestLogsResponse{
		List:     nil,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	temp := array.Map(list, func(item *model.RequestLog) apiV1.RequestLogItem {
		return apiV1.RequestLogItem{
			Id:               strconv.FormatUint(item.Id, 10),
			Model:            item.Model,
			UserId:           strconv.FormatUint(item.UserId, 10),
			Username:         item.Username,
			Email:            item.Email,
			Ip:               item.Ip,
			Status:           item.Status,
			RetryTimes:       item.RetryTimes,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			ChannelNameTrace: item.ChannelNameTrace,
			ChannelIdTrace:   item.ChannelIdTrace,
		}
	})
	resp.List = temp
	return resp, nil
}

func (r *requestLogService) GetRequestLogsModelRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) (*apiV1.RequestLogsModelRankingResponse, error) {
	resp := new(apiV1.RequestLogsModelRankingResponse)
	list, err := r.repo.FindRequestLogsModelRanking(ctx, req)
	if err != nil {
		return nil, err
	}
	temp := array.Map(list, func(item *apiV1.RequestLogsModelRanking) apiV1.RequestLogsModelRanking {
		return *item
	})
	resp.List = temp
	return resp, nil
}

func (r *requestLogService) GetRequestLogsUserRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) (*apiV1.RequestLogsUserRankingResponse, error) {
	resp := new(apiV1.RequestLogsUserRankingResponse)
	list, err := r.repo.FindRequestLogsUserRanking(ctx, req)
	if err != nil {
		return nil, err
	}
	temp := array.Map(list, func(item *apiV1.RequestLogsUserRanking) apiV1.RequestLogsUserRanking {
		return *item
	})
	resp.List = temp
	return resp, nil
}
