package v1

type OAuthCbRequest struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}
