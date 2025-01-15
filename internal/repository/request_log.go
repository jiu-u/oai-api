package repository

import (
	"context"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
)

type RequestLogStatisticsQuery struct {
	StartTime string
	EndTime   string
}

type UserCallCount struct {
	UserId    uint64 `json:"userId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CallCount int    `json:"callCount"`
}

type RequestLogRepository interface {
	CreateRequestLog(ctx context.Context, log *model.RequestLog) error
	FindRequestLogs(ctx context.Context, req *apiV1.RequestLogsQuery) ([]*model.RequestLog, int64, error)
	FindRequestLogsModelRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) ([]*apiV1.RequestLogsModelRanking, error)
	FindRequestLogsUserRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) ([]*apiV1.RequestLogsUserRanking, error)
}

func NewRequestLogRepository(repo *Repository) RequestLogRepository {
	return &requestLogRepository{repo}
}

type requestLogRepository struct {
	*Repository
}

func (r *requestLogRepository) CreateRequestLog(ctx context.Context, log *model.RequestLog) error {
	return r.DB(ctx).Create(log).Error
}

func (r *requestLogRepository) FindRequestLogs(ctx context.Context, req *apiV1.RequestLogsQuery) ([]*model.RequestLog, int64, error) {
	var logs []*model.RequestLog
	var err error
	query := r.DB(ctx).Model(&model.RequestLog{})
	if req.StartTime != "" && req.EndTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	}
	if req.UserId != "" {
		query = query.Where("user_id = ?", req.UserId)
	}
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if req.Page > 0 && req.PageSize > 0 {
		query = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("id desc")
	}
	err = query.Find(&logs).Error
	return logs, total, err
}

func (r *requestLogRepository) FindRequestLogsModelRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) ([]*apiV1.RequestLogsModelRanking, error) {
	var logs []*apiV1.RequestLogsModelRanking
	q := r.DB(ctx).Model(&model.RequestLog{})
	if req.StartTime != "" && req.EndTime != "" {
		q = q.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	}
	q = q.Select("model, count(*) as call_count").Group("model").Order("call_count DESC")
	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}
	err := q.Find(&logs).Error
	return logs, err
}

func (r *requestLogRepository) FindRequestLogsUserRanking(ctx context.Context, req *apiV1.RequestLogsRankingRequest) ([]*apiV1.RequestLogsUserRanking, error) {
	var logs []*apiV1.RequestLogsUserRanking
	q := r.DB(ctx).Model(&model.RequestLog{})
	if req.StartTime != "" && req.EndTime != "" {
		q = q.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	}
	q = q.Select("user_id, username, count(*) as call_count").Group("user_id, username").Order("call_count DESC")
	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}
	err := q.Find(&logs).Error
	return logs, err
}
