package service

import (
	"context"
	"errors"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/dto"
	"github.com/jiu-u/oai-api/internal/model"
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/internal/service/oauth2"
	"github.com/jiu-u/oai-api/pkg/datautils"
	"github.com/jiu-u/oai-api/pkg/encrypte"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/vaild"
	"github.com/lithammer/shortuuid/v4"
	"strconv"
	"strings"
	"time"
)

// todo oauth2统一封装抽象
// ! 需要记录那些信息，clientId、clientSecret、authRUL、tokenURL、userURL。。。
// ! 数据库存储那些信息

type AuthService interface {
	NewAccessToken(ctx context.Context, userId uint64) (*apiV1.AccessTokenResponse, error)
	UserLogin(ctx context.Context, req *apiV1.UserLoginReq) (*apiV1.AuthResponse, error)
	UserRegister(ctx context.Context, req *apiV1.UserRegisterReq) (*apiV1.AuthResponse, error)
	GetAuthResponseBySessionId(ctx context.Context, sessionId string) (*apiV1.AuthResponse, error)
	LinuxDoCallBack(ctx context.Context, req *apiV1.OAuthCbRequest) (sessionId string, err error)
	GetGitHubRedirectURL(ctx context.Context, state string) (string, error)
	GitHubCallBack(ctx context.Context, req *apiV1.OAuthCbRequest) (sessionId string, err error)
}

func NewAuthService(
	s *Service,
	userRepo repository.UserRepository,
	systemConfigSvc SystemConfigService,
	linuxDoAuth *oauth2.LinuxDoOauthService,
	githubAuth *oauth2.GitHubOauthService,
	oauth2Repo repository.UserAuthProviderRepository,
) AuthService {
	return &authService{
		Service:         s,
		userRepo:        userRepo,
		systemConfigSvc: systemConfigSvc,
		linuxDoAuth:     linuxDoAuth,
		githubAuth:      githubAuth,
		oauth2Repo:      oauth2Repo,
	}
}

type authService struct {
	*Service
	userRepo        repository.UserRepository
	systemConfigSvc SystemConfigService
	linuxDoAuth     *oauth2.LinuxDoOauthService
	githubAuth      *oauth2.GitHubOauthService
	oauth2Repo      repository.UserAuthProviderRepository
}

func (s *authService) GenAuthResponse(ctx context.Context, userId uint64, role string) (*apiV1.AuthResponse, error) {
	resp := new(apiV1.AuthResponse)
	accessToken, err := s.Jwt.GenAccessToken(userId, role)
	if err != nil {
		return nil, err
	}
	resp.AccessToken = accessToken
	refreshToken, err := s.Jwt.GenRefreshToken(userId, role)
	if err != nil {
		return nil, err
	}
	resp.RefreshToken = refreshToken
	resp.ExpiredAt = time.Now().Add(jwt.AccessTokenDuration).Unix()
	resp.TokenType = "Bearer"
	resp.UserId = strconv.FormatUint(userId, 10)
	resp.Success = true
	resp.Role = role
	return resp, nil
}

func (s *authService) NewAccessToken(ctx context.Context, userId uint64) (*apiV1.AccessTokenResponse, error) {
	user, err := s.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	resp := new(apiV1.AccessTokenResponse)
	accessToken, err := s.Jwt.GenAccessToken(user.Id, user.Role)
	if err != nil {
		return nil, err
	}
	resp.AccessToken = accessToken
	return resp, nil
}

func (s *authService) UserLogin(ctx context.Context, req *apiV1.UserLoginReq) (*apiV1.AuthResponse, error) {
	var err error
	cfg, err := s.systemConfigSvc.GetRegisterConfig(ctx)
	if err != nil {
		// 应该永远不会到这里
		// 配置不存在
		return nil, errors.New("配置不存在")
	}
	if !cfg.AllowLoginByPassword {
		return nil, errors.New("该登录方式已被禁用")
	}
	isEmail := strings.Contains(req.Username, "@")
	var user *model.User
	if isEmail {
		user, err = s.userRepo.FindUserByEmail(ctx, req.Username)
	} else {
		user, err = s.userRepo.FindUserByUsername(ctx, req.Username)
	}
	if err != nil {
		return nil, errors.New("用户不存在或数据库错误")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	// 校验秘密
	err = encrypte.VerifyPassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("密码错误")
	}
	resp, err := s.GenAuthResponse(ctx, user.Id, user.Role)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *authService) UserRegister(ctx context.Context, req *apiV1.UserRegisterReq) (*apiV1.AuthResponse, error) {
	if !vaild.IsValidUsername(req.Username) {
		return nil, errors.New("invalid username")
	}
	var err error
	cfg, err := s.systemConfigSvc.GetRegisterConfig(ctx)
	if err != nil {
		// 应该永远不会到这里
		// 配置不存在
		return nil, errors.New("配置不存在")
	}
	if !cfg.AllowRegister {
		return nil, errors.New("注册已关闭")
	}
	if !cfg.AllowRegisterByPassword {
		return nil, errors.New("该注册方式已被禁用")
	}
	// 检查用户名和邮箱
	var user *model.User
	if req.Email != "" {
		user, err = s.userRepo.FindUserByEmail(ctx, req.Email)
		if err == nil {
			return nil, errors.New("该邮箱已被注册")
		}
	}
	user, err = s.userRepo.FindUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("该用户名已被注册")
	}
	// 是否开启了邮箱验证
	if cfg.AllowEmailValid {
		if req.Email == "" {
			return nil, errors.New("邮箱不能为空")
		}
		// 比较验证码是否一致
		code, exist := s.Cache.Get("email_" + req.Email)
		if !exist || code == nil {
			return nil, errors.New("验证码已过期")
		}
		if code.(string) != req.VerificationCode {
			return nil, errors.New("验证码错误")
		}
		// 清除缓存
		s.Cache.Delete("email_" + req.Email)
	}
	password, err := encrypte.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("服务器内部出错")
	}
	user = &model.User{
		Username:    req.Username,
		Password:    password,
		Role:        "user",
		LastLoginAt: time.Now(),
		LastLoginIP: GetClientIp(ctx),
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	user.Id = s.Sid.GenUint64()
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return s.GenAuthResponse(ctx, user.Id, user.Role)
}

func (s *authService) GetAuthResponseBySessionId(ctx context.Context, sessionId string) (*apiV1.AuthResponse, error) {
	value, ok := s.Cache.Get("session_" + sessionId)
	if !ok {
		return nil, errors.New("无效的对话")
	}
	userId, err := strconv.ParseUint(value.(string), 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("not found")
	}
	s.Cache.Delete("session_" + sessionId)
	return s.GenAuthResponse(ctx, userId, user.Role)

}

func (s *authService) LinuxDoCallBack(ctx context.Context, req *apiV1.OAuthCbRequest) (sessionId string, err error) {
	userInfo, err := s.linuxDoAuth.CallBackHandle(ctx, req.Code)
	if err != nil {
		return "", err
	}
	var userId uint64
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		// 查找存储记录
		record, err := s.oauth2Repo.FindUserAuthProvider(ctx, dto.LinuxDoOAuthType, strconv.Itoa(userInfo.Id))
		if err != nil {
			return err
		}
		if record == nil {
			// 记录不存在
			u := &model.User{
				Username:    "linux_do_" + datautils.SecureRandomString(5),
				Email:       nil,
				Phone:       nil,
				Role:        "user",
				Nickname:    "linux_do_" + datautils.SecureRandomString(5),
				Level:       userInfo.TrustLevel,
				LastLoginAt: time.Now(),
				LastLoginIP: GetClientIp(ctx),
			}
			u.Id = s.Sid.GenUint64()
			err = s.userRepo.CreateUser(ctx, u)
			if err != nil {
				return err
			}
			userId = u.Id
			p := &model.UserAuthProvider{
				UserId:         u.Id,
				Provider:       dto.LinuxDoOAuthType,
				ProviderUserId: strconv.Itoa(userInfo.Id),
				AccessToken:    "",
				RefreshToken:   "",
				TokenExpireAt:  time.Now(),
				Scope:          "",
				ProviderEmail:  userInfo.Email,
				ProviderName:   userInfo.Username,
				ProviderAvatar: "",
				UniqueProvider: dto.LinuxDoOAuthType,
				UniqueUserId:   userId,
			}
			p.Id = s.Sid.GenUint64()
			err = s.oauth2Repo.CreateUserAuthProvider(ctx, p)
			if err != nil {
				return err
			}
		} else {
			userId = record.UserId
			p := &model.UserAuthProvider{
				Provider:       dto.LinuxDoOAuthType,
				ProviderUserId: strconv.Itoa(userInfo.Id),
				AccessToken:    "",
				RefreshToken:   "",
				TokenExpireAt:  time.Now(),
				Scope:          "",
				ProviderEmail:  userInfo.Email,
				ProviderName:   userInfo.Username,
				ProviderAvatar: "",
				UniqueProvider: dto.LinuxDoOAuthType,
			}
			p.Id = record.Id
			err = s.oauth2Repo.UpdateUserAuthProviderById(ctx, p)
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	uuid := shortuuid.New()
	s.Cache.Set("session_"+uuid, strconv.FormatUint(userId, 10), time.Minute*10)
	return uuid, nil
}

func (s *authService) GetGitHubRedirectURL(ctx context.Context, state string) (string, error) {
	return s.githubAuth.GetRedirectURL(ctx, state)
}

func (s *authService) GitHubCallBack(ctx context.Context, req *apiV1.OAuthCbRequest) (sessionId string, err error) {
	userInfo, err := s.githubAuth.CallBackHandle(ctx, req.Code)
	if err != nil {
		return "", err
	}
	var userId uint64
	err = s.Tm.Transaction(ctx, func(ctx context.Context) error {
		// 查找存储记录
		record, err := s.oauth2Repo.FindUserAuthProvider(ctx, dto.LinuxDoOAuthType, strconv.Itoa(userInfo.Id))
		if err != nil {
			return err
		}
		if record == nil {
			// 记录不存在
			u := &model.User{
				Username:    "github_" + datautils.SecureRandomString(5),
				Email:       nil,
				Phone:       nil,
				Role:        "user",
				Nickname:    "github_" + datautils.SecureRandomString(5),
				Level:       1,
				LastLoginAt: time.Now(),
				LastLoginIP: GetClientIp(ctx),
			}
			u.Id = s.Sid.GenUint64()
			err = s.userRepo.CreateUser(ctx, u)
			if err != nil {
				return err
			}
			userId = u.Id

			p := &model.UserAuthProvider{
				UserId:         u.Id,
				Provider:       dto.LinuxDoOAuthType,
				ProviderUserId: strconv.Itoa(userInfo.Id),
				AccessToken:    "",
				RefreshToken:   "",
				TokenExpireAt:  time.Now(),
				Scope:          "",
				ProviderEmail:  userInfo.Email,
				ProviderName:   userInfo.Login,
				ProviderAvatar: "",
				UniqueProvider: dto.LinuxDoOAuthType,
				UniqueUserId:   userId,
			}
			p.Id = s.Sid.GenUint64()
			err = s.oauth2Repo.CreateUserAuthProvider(ctx, p)
			if err != nil {
				return err
			}
		} else {
			userId = record.UserId
			p := &model.UserAuthProvider{
				Provider:       dto.LinuxDoOAuthType,
				ProviderUserId: strconv.Itoa(userInfo.Id),
				AccessToken:    "",
				RefreshToken:   "",
				TokenExpireAt:  time.Now(),
				Scope:          "",
				ProviderEmail:  userInfo.Email,
				ProviderName:   userInfo.Name,
				ProviderAvatar: "",
				UniqueProvider: dto.LinuxDoOAuthType,
			}
			p.Id = record.Id
			err = s.oauth2Repo.UpdateUserAuthProviderById(ctx, p)
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	uuid := shortuuid.New()
	s.Cache.Set("session_"+uuid, strconv.FormatUint(userId, 10), time.Minute*10)
	return uuid, nil
}
