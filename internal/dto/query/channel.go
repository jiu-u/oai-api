package query

type ChannelQueryRequest struct {
	Name     string `json:"name" form:"name"`
	Type     string `json:"type" form:"type"`
	Status   int8   `json:"status" form:"status"`
	Page     int    `json:"page" form:"page" default:"1" binding:"required,min=1"`
	PageSize int    `json:"pageSize" form:"pageSize" default:"10" binding:"required,min=1"`
}
