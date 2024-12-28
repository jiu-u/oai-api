package v1

type AuthResponse struct {
	UserId       string `json:"userId"`
	Success      bool   `json:"success"`
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type LinuxDoAuthRequest struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" form:"refreshToken"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
