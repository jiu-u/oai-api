package v1

import "github.com/jiu-u/oai-api/internal/repository"

type RequestLogRanking struct {
	StartTime string `json:"startTime" form:"startTime" binding:"required"`
	EndTime   string `json:"endTime" form:"endTime"  binding:"required"`
}

type UseCallCountResponse struct {
	Data []UserCallCount `json:"data"`
}

type UserCallCount = repository.UserCallCount

type RequestLogRealTimeResponse struct {
	Data []RequestLogRealTimeItem `json:"data"`
}
type RequestLogRealTimeItem struct {
	UserId    uint64 `json:"userId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int8   `json:"status"`
	Model     string `json:"model"`
	CreatedAt string `json:"createdAt"`
}
