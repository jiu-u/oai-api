package dto

type ModelCheckResult struct {
	ChannelId          string `json:"channelId"`
	ModelName          string `json:"model"`
	ConnectionDuration int64  `json:"connectionDuration"`
	TotalDuration      int64  `json:"totalDuration"`
	Status             int8   `json:"status"`
}
