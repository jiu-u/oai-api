package v1

type RequestLogItem struct {
	Id               string `json:"id"`
	Model            string `json:"model"`
	UserId           string `json:"userId"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Ip               string `json:"ip"`
	Status           int8   `json:"status"`
	RetryTimes       int    `json:"retryTimes"`
	CreatedAt        string `json:"createdAt"`
	ChannelNameTrace string `json:"channelNameTrace"`
	ChannelIdTrace   string `json:"channelIdTrace"`
}

type RequestLogsQuery struct {
	StartTime string `form:"startTime" binding:"required"`
	EndTime   string `form:"endTime"  binding:"required"`
	Page      int    `form:"page"  binding:"required"`
	PageSize  int    `form:"pageSize"  binding:"required"`
	UserId    string `form:"userId"`
}

type RequestLogsResponse struct {
	List     []RequestLogItem `json:"list"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

type RequestLogsRankingRequest struct {
	StartTime string `json:"startTime" form:"startTime" binding:"required"`
	EndTime   string `json:"endTime" form:"endTime"  binding:"required"`
	Limit     int    `json:"limit" form:"limit"  binding:"required"`
	UserId    string `json:"userId" form:"userId"`
}

type RequestLogsUserRanking struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	CallCount int    `json:"callCount"`
}

type RequestLogsUserRankingResponse struct {
	List []RequestLogsUserRanking `json:"list"`
}

type RequestLogsModelRanking struct {
	Model     string `json:"model"`
	CallCount int    `json:"callCount"`
}

type RequestLogsModelRankingResponse struct {
	List []RequestLogsModelRanking `json:"list"`
}
