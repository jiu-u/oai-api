package oauth2

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AuthService interface {
	LinuxDoAuthHandle(ctx context.Context, req *apiV1.LinuxDoAuthRequest) (string, error)
	GetSessionUser(ctx context.Context, sessionId string) (*apiV1.AuthResponse, error)
}

func NewService(
	svc *service.Service,
	userRepo repository.UserRepository,
	LinuxDoOauth2 *LinuxDoOauth,
) AuthService {
	return &authService{
		Service:  svc,
		userRepo: userRepo,
		LinuxDo:  LinuxDoOauth2,
	}
}

type authService struct {
	*service.Service
	userRepo repository.UserRepository
	LinuxDo  *LinuxDoOauth
}

func (a *authService) LinuxDoAuthHandle(ctx context.Context, req *apiV1.LinuxDoAuthRequest) (string, error) {
	// 请求AccessToken
	// 请求用户信息
	userInfo, err := a.LinuxDo.CallBackHandle(ctx, req.Code)
	if err != nil {
		return "", err
	}
	var id uint64
	role := "user"
	// 存储用户信息
	user, err := a.userRepo.FindOneByLinuxDoId(ctx, uint64(userInfo.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			id = a.Sid.GenUint64()
			user = &model.User{
				Username:        userInfo.Username,
				UserEmail:       userInfo.Email,
				LinuxDoId:       uint64(userInfo.Id),
				LinuxDoUsername: userInfo.Username,
				LinuxDoLevel:    userInfo.TrustLevel,
			}
			user.Id = id
			err = a.userRepo.InsertOne(ctx, user)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		id = user.Id
		role = user.Role
	}
	// 生成token
	accessToken, err := a.Jwt.GenAccessToken(id, role)
	if err != nil {
		return "", err
	}
	refreshToken, err := a.Jwt.GenRefreshToken(id, role)
	if err != nil {
		return "", err
	}
	resp := new(apiV1.AuthResponse)
	resp.UserId = strconv.FormatUint(id, 10)
	resp.Success = true
	resp.RefreshToken = refreshToken
	resp.AccessToken = accessToken
	jsonBytes, err := sonic.Marshal(resp)
	if err != nil {
		return "", err
	}
	uuid := shortuuid.New()
	a.Cache.Set("session_"+uuid, string(jsonBytes), time.Minute*3)
	return uuid, nil
}

func (a *authService) GetSessionUser(ctx context.Context, sessionId string) (*apiV1.AuthResponse, error) {
	v, exist := a.Cache.Get("session_" + sessionId)
	if v == nil || !exist {
		return nil, errors.New("session not found")
	}
	defer a.Cache.Delete("session_" + sessionId)
	resp := new(apiV1.AuthResponse)
	err := sonic.Unmarshal([]byte(v.(string)), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
