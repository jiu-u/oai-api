package oauth2

import (
	"context"
	"time"
)

type AuthCbReq struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type Record struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpireAt     time.Time `json:"expireAt"`
}

type Provider interface {
	GetRedirectURL(ctx context.Context, state string) (string, error)
	CallBackHandle(ctx context.Context, req *AuthCbReq) (*Record, error)
}
