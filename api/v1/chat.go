package v1

type OnlyModelChatRequest struct {
	Model string `json:"model" binding:"required"`
}
