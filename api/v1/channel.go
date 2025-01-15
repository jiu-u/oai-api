package v1

import (
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/dto/query"
)

type ChannelResponse struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Balance  float64  `json:"balance"`
	EndPoint string   `json:"endPoint"`
	APIKey   string   `json:"apiKey"`
	Models   []string `json:"models"`
	Status   int8     `json:"status"`
}

type CreateChannelRequest struct {
	Name     string   `json:"name"`
	Type     string   `json:"type" binding:"required"`
	EndPoint string   `json:"endPoint" binding:"required"`
	APIKey   string   `json:"apiKey" binding:"required"`
	Weight   int      `json:"weight" default:"10"`
	Models   []string `json:"models"`
}

type ChannelQueryRequest = query.ChannelQueryRequest
type ChannelListResponse struct {
	Total    int64             `json:"total"`
	Page     int64             `json:"page"`
	PageSize int64             `json:"pageSize"`
	List     []ChannelResponse `json:"list"`
}

type UpdateChannelRequest struct {
	Name     string   `json:"name"`
	Type     string   `json:"type" `
	EndPoint string   `json:"endPoint"`
	APIKey   string   `json:"apiKey"`
	Models   []string `json:"models"`
	Status   int8     `json:"status"`
}

type ChannelModelTestResponse = dto.ModelCheckResult

type CheckModelRequest struct {
	ModelName string `json:"model" binding:"required"`
}
