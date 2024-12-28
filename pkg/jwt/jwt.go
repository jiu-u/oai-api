package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jiu-u/oai-api/pkg/config"
	"strings"
	"time"
)

type TokenScopeType string

const (
	ACCESS  TokenScopeType = "access"
	REFRESH TokenScopeType = "refresh"
	//TOKEN_TYPE string         = "Bearer "
)

type JWT struct {
	key []byte
}

type MyCustomClaims struct {
	Role       string
	UserId     uint64
	Username   string
	TokenScope TokenScopeType
	jwt.RegisteredClaims
}

func NewJwt(conf *config.Config) *JWT {
	return &JWT{key: []byte(conf.Security.Jwt.Key)}
}

func (j *JWT) GenToken(userId uint64, role string, tokenScope TokenScopeType, expiresAt time.Time) (string, error) {
	beforeNow := time.Now().Add(-1 * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId:     userId,
		Role:       role,
		TokenScope: tokenScope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			//IssuedAt:  jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(beforeNow),
			NotBefore: jwt.NewNumericDate(beforeNow),
			Issuer:    "",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	})
	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string, prefix string) (*MyCustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, prefix)
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (j *JWT) GenAccessToken(userId uint64, role string) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 15)
	return j.GenToken(userId, role, ACCESS, expiresAt)
}
func (j *JWT) GenRefreshToken(userId uint64, role string) (string, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	return j.GenToken(userId, role, REFRESH, expiresAt)
}

func (j *JWT) ParseRefreshToken(tokenString string, prefix string) (*MyCustomClaims, error) {
	claims, err := j.ParseToken(tokenString, prefix)
	if err != nil {
		return nil, err
	}
	if claims.TokenScope != REFRESH {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
func (j *JWT) ParseAccessToken(tokenString string, prefix string) (*MyCustomClaims, error) {
	claims, err := j.ParseToken(tokenString, prefix)
	if err != nil {
		return nil, err
	}
	if claims.TokenScope != ACCESS {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}
