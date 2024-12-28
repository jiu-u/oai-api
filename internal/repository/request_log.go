package repository

import (
	"context"
	"github.com/jiu-u/oai-api/internal/model"
	"time"
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
	InsertOne(ctx context.Context, log *model.RequestLog) error
	GetStatistics(ctx context.Context, query *RequestLogStatisticsQuery) ([]UserCallCount, error)
	GetNewRequestLogs(ctx context.Context, limit int) ([]*model.RequestLog, error)
	QueryItemByApiKey(ctx context.Context, apiKey string) (*model.RequestLog, error)
}

func NewRequestLogRepository(repo *Repository) RequestLogRepository {
	return &requestLogRepository{repo}
}

type requestLogRepository struct {
	*Repository
}

func (r *requestLogRepository) QueryItemByApiKey(ctx context.Context, apiKey string) (*model.RequestLog, error) {
	var log model.RequestLog
	err := r.DB(ctx).Where("api_key = ?", apiKey).First(&log).Error
	return &log, err
}

func (r *requestLogRepository) GetNewRequestLogs(ctx context.Context, limit int) ([]*model.RequestLog, error) {
	var logs []*model.RequestLog
	err := r.DB(ctx).Limit(limit).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *requestLogRepository) InsertOne(ctx context.Context, log *model.RequestLog) error {
	return r.DB(ctx).Create(log).Error
}

func (r *requestLogRepository) GetStatistics(ctx context.Context, query *RequestLogStatisticsQuery) ([]UserCallCount, error) {
	// 解析时间
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", query.StartTime, time.Local)
	if err != nil {
		panic("invalid start time format")
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", query.EndTime, time.Local)
	if err != nil {
		panic("invalid end time format")
	}
	// 查询用户调用排行榜单
	var userCallCounts []UserCallCount
	err = r.DB(ctx).Model(&model.RequestLog{}).
		Select("user_id, username, count(*) as call_count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("user_id, username").
		Order("call_count DESC").
		Find(&userCallCounts).Error
	return userCallCounts, err
}
