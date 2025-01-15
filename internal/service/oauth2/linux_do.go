package oauth2

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-api/internal/repository"
	stdoauth2 "golang.org/x/oauth2"
	"io"
	"net/http"
)

const (
	LinuxDoRedirectURI           = "/v1/oauth2/linux_do/callback"
	LinuxDoAuthorizationEndpoint = "https://connect.linux.do/oauth2/authorize"
	LinuxDoTokenEndpoint         = "https://connect.linux.do/oauth2/token"
	LinuxDoUserEndpoint          = "https://connect.linux.do/api/user"
)

type LinuxDoOauthService struct {
	oauth2Config *stdoauth2.Config
	repo         repository.SystemRepository
}

func NewLinuxDoAuthService(repo repository.SystemRepository) *LinuxDoOauthService {
	return &LinuxDoOauthService{
		repo: repo,
		oauth2Config: &stdoauth2.Config{
			ClientID:     "",
			ClientSecret: "",
			Endpoint: stdoauth2.Endpoint{
				AuthURL:   LinuxDoAuthorizationEndpoint,
				TokenURL:  LinuxDoTokenEndpoint,
				AuthStyle: 0,
			},
			RedirectURL: LinuxDoRedirectURI,
			Scopes:      nil,
		},
	}
}

func (s *LinuxDoOauthService) CallBackHandle(ctx context.Context, code string) (*LinuxDoUser, error) {
	// 检查state
	// 查看是否可用
	cfg, err := s.repo.GetLinuxDoOAuthConfig(ctx)
	if err != nil {
		return nil, errors.New("该登录方式已被关闭")
	}
	s.oauth2Config.ClientID = cfg.ClientId
	s.oauth2Config.ClientSecret = cfg.ClientSecret
	redirectUrl := LinuxDoRedirectURI
	token, err := s.GetAccessToken(ctx, code, redirectUrl)
	if err != nil {
		return nil, err
	}
	return s.GetUserInfo(ctx, token)
}

func (s *LinuxDoOauthService) GetAccessToken(ctx context.Context, code, redirectUrl string) (string, error) {
	// 获取配置
	token, err := s.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, err
	//data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, redirectUrl)
	//req, err := http.NewRequest("POST", o.TokenEndpoint, bytes.NewBuffer([]byte(data)))
	//if err != nil {
	//	return "", err
	//}
	//
	//req.SetBasicAuth(o.ClientId, o.ClientSecret)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Accept", "application/json")
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	detail, err := io.ReadAll(resp.Body)
	//	if err != nil {
	//		detail = []byte("no detail")
	//	}
	//	return "", fmt.Errorf("failed to fetch access token: %d description: %s", resp.StatusCode, detail)
	//}
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//var tokenResponse map[string]interface{}
	//if err := sonic.Unmarshal(body, &tokenResponse); err != nil {
	//	return "", err
	//}
	//
	//return tokenResponse["access_token"].(string), nil
}

type LinuxDoUser struct {
	Active         bool
	ApiKey         string
	AvatarTemplate string
	AvatarUrl      string
	Email          string
	ExternalIds    interface{}
	Id             int
	Login          string
	Name           string
	Silenced       bool
	Sub            string
	TrustLevel     int
	Username       string
}

func (s *LinuxDoOauthService) GetUserInfo(ctx context.Context, token string) (*LinuxDoUser, error) {
	req, err := http.NewRequest("GET", LinuxDoUserEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo *LinuxDoUser
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
