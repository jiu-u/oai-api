package v1

type AuthLoginReq struct {
	Username string
}

type VerificationEmailReq struct {
	Email string `json:"email" form:"email" binding:"required,email"`
	Code  string `json:"code" form:"code" binding:"required"`
}

type UserLoginReq struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=100"` // 用户名: 必填，长度3-20，字母和数字
	Password string `json:"password" form:"password" binding:"required,min=8,max=32"`  // 密码: 必填，最小长度8
}

type UserRegisterReq struct {
	Email            string `json:"email" form:"email"`
	Username         string `json:"username" form:"username" binding:"required,min=3,max=32"`
	Password         string `json:"password" form:"password" binding:"required,min=8,max=32"`
	VerificationCode string `json:"verificationCode" form:"verificationCode"`
}

type AuthResponse struct {
	UserId       string `json:"userId"`
	Success      bool   `json:"success"`
	Role         string `json:"role"`
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiredAt    int64  `json:"expiredAt"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" form:"refreshToken" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
