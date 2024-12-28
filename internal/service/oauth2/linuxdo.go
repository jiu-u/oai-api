package oauth2

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-api/pkg/config"
	"io"
	"net/http"
)

const (
	LinuxDoRedirectURI           = "/api/oauth2/linuxdo/callback"
	LinuxDoAuthorizationEndpoint = "https://connect.linux.do/oauth2/authorize"
	LinuxDoTokenEndpoint         = "https://connect.linux.do/oauth2/token"
	LinuxDoUserEndpoint          = "https://connect.linux.do/api/user"
)

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

type LinuxDoOauth struct {
	ClientId              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
}

func NewLinuxDoService(conf *config.Config) *LinuxDoOauth {
	return &LinuxDoOauth{
		ClientId:              conf.Oauth.LinuxDo.ClientId,
		ClientSecret:          conf.Oauth.LinuxDo.ClientSecret,
		AuthorizationEndpoint: LinuxDoAuthorizationEndpoint,
		TokenEndpoint:         LinuxDoTokenEndpoint,
		UserEndpoint:          LinuxDoUserEndpoint,
	}
}

func (o *LinuxDoOauth) CallBackHandle(ctx context.Context, code string) (*LinuxDoUser, error) {
	redirectUrl := LinuxDoRedirectURI
	token, err := o.GetAccessToken(ctx, code, redirectUrl)
	if err != nil {
		return nil, err
	}
	return o.GetUserInfo(ctx, token)
}

func (o *LinuxDoOauth) GetAccessToken(ctx context.Context, code, redirectUrl string) (string, error) {
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, redirectUrl)
	req, err := http.NewRequest("POST", o.TokenEndpoint, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(o.ClientId, o.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, err := io.ReadAll(resp.Body)
		if err != nil {
			detail = []byte("no detail")
		}
		return "", fmt.Errorf("failed to fetch access token: %d description: %s", resp.StatusCode, detail)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse map[string]interface{}
	if err := sonic.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse["access_token"].(string), nil
}

func (o *LinuxDoOauth) GetUserInfo(ctx context.Context, token string) (*LinuxDoUser, error) {
	req, err := http.NewRequest("GET", o.UserEndpoint, nil)
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
