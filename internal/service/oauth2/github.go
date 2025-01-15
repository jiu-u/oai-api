package oauth2

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jiu-u/oai-api/internal/repository"
	stdoauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"net/http"
)

const (
	GithubUserEndpoint = "https://api.github.com/user"
)

type GitHubOauthService struct {
	oauth2Config *stdoauth2.Config
	repo         repository.SystemRepository
}

func NewGithubAuthService(repo repository.SystemRepository) *GitHubOauthService {
	return &GitHubOauthService{
		repo: repo,
		oauth2Config: &stdoauth2.Config{
			ClientID:     "",
			ClientSecret: "",
			Scopes:       []string{"user:email", "user:username"},
			Endpoint:     github.Endpoint,
			RedirectURL:  "",
		},
	}
}

func (s *GitHubOauthService) GetRedirectURL(ctx context.Context, state string) (string, error) {
	cfg, err := s.repo.GetGithubOAuthConfig(ctx)
	if err != nil {
		return "", errors.New("该登录方式已被关闭")
	}
	s.oauth2Config.ClientID = cfg.ClientId
	s.oauth2Config.ClientSecret = cfg.ClientSecret
	url := s.oauth2Config.AuthCodeURL(state, stdoauth2.AccessTypeOnline)
	return url, nil
}

func (s *GitHubOauthService) CallBackHandle(ctx context.Context, code string) (*GitHubUser, error) {
	// 检查state
	// 查看是否可用
	cfg, err := s.repo.GetGithubOAuthConfig(ctx)
	if err != nil {
		return nil, errors.New("该登录方式已被关闭")
	}
	s.oauth2Config.ClientID = cfg.ClientId
	s.oauth2Config.ClientSecret = cfg.ClientSecret
	redirectUrl := ""
	token, err := s.GetAccessToken(ctx, code, redirectUrl)
	if err != nil {
		return nil, err
	}
	return s.GetUserInfo(ctx, token)
}

func (s *GitHubOauthService) GetAccessToken(ctx context.Context, code, redirectUrl string) (string, error) {
	// 获取配置
	token, err := s.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, err
}

type GitHubResp struct {
	User GitHubUser `json:"user"`
}

type GitHubUser struct {
	AvatarURL         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	Company           string `json:"company"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsURL         string `json:"events_url"`
	Followers         int    `json:"followers"`
	FollowersURL      string `json:"followers_url"`
	Following         int    `json:"following"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	Hireable          bool   `json:"hireable"`
	HTMLURL           string `json:"html_url"`
	Id                int    `json:"id"`
	Location          string `json:"location"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	NodeID            string `json:"node_id"`
	NotificationEmail string `json:"notification_email"`
	OrganizationsURL  string `json:"organizations_url"`
	PublicGists       int    `json:"public_gists"`
	PublicRepos       int    `json:"public_repos"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	TwitterUsername   string `json:"twitter_username"`
	Type              string `json:"type"`
	UpdatedAt         string `json:"updated_at"`
	URL               string `json:"url"`
	UserViewType      string `json:"user_view_type"`
}

func (s *GitHubOauthService) GetUserInfo(ctx context.Context, token string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", GithubUserEndpoint, nil)
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

	var userInfo *GitHubResp
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo.User, nil
}
