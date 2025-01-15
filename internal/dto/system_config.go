package dto

type EmailConfig struct {
	Id       uint64 `json:"id"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OAuthConfig struct {
	Provider     string   `json:"provider"` // prefix linuxDo_ github_
	ClientId     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Scopes       []string `json:"scopes"`
	AuthURL      string   `json:"authURL"`
	TokenURL     string   `json:"tokenURL"`
	UserURL      string   `json:"userURL"`
}
type OAuthProviderType = string

const (
	LinuxDoOAuthType OAuthProviderType = "linux_do"
	GithubOAuthType                    = "github"
)

type LinuxDoOAuthConfig struct {
	Id           uint64 `json:"id"`
	ClientId     string `json:"clientId" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
}

type GithubOAuthConfig = LinuxDoOAuthConfig

type ModelConfig struct {
	Id           uint64              `json:"id"`
	ModelMapping map[string][]string `json:"modelMapping"`
	CheckList    []string            `json:"checkList"`
}

type RegisterConfig struct {
	Id                      uint64 `json:"id"`
	AllowRegister           bool   `json:"allowRegister"`
	AllowRegisterByPassword bool   `json:"allowRegisterByPassword"`
	AllowLoginByPassword    bool   `json:"allowLoginByPassword"`
	AllowEmailValid         bool   `json:"allEmailValid"`
	AllowLinuxDoLogin       bool   `json:"allowLinuxDoLogin"`
	AllowGithubLogin        bool   `json:"allowGithubLogin"`
}
