package handler

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service"
	"github.com/jiu-u/oai-api/internal/service/oauth2"
	"github.com/jiu-u/oai-api/pkg/jwt"
	"github.com/jiu-u/oai-api/pkg/vaild"
	"github.com/lithammer/shortuuid/v4"
	"net/http"
	"os"
)

var RedirectURI string = "http://localhost:5173/#/login?sessionId="

func getRedirectURI() string {
	url := os.Getenv("OAI_REDIRECT_URI")
	if url == "" {
		return RedirectURI
	}
	return url
}

type AuthHandler struct {
	*Handler
	jwt                 *jwt.JWT
	svc                 service.AuthService
	systemConfigService service.SystemConfigService
}

func NewAuthHandler(
	handler *Handler,
	jwt *jwt.JWT,
	svc service.AuthService,
	systemConfigService service.SystemConfigService,
) *AuthHandler {
	return &AuthHandler{
		Handler:             handler,
		jwt:                 jwt,
		svc:                 svc,
		systemConfigService: systemConfigService,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	req := new(apiV1.UserLoginReq)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.UserLogin(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	req := new(apiV1.UserRegisterReq)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	if req.Username == "" && req.Email == "" {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "username or email is required")
		return
	}
	if req.Email != "" && !vaild.IsValidEmail(req.Email) {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "invalid email")
		return
	}
	resp, err := h.svc.UserRegister(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *AuthHandler) GetNewAccessToken(ctx *gin.Context) {
	req := new(apiV1.AccessTokenRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	claim, err := h.jwt.ParseRefreshToken(req.RefreshToken, "Bearer ")
	if err != nil {
		apiV1.HandleError(ctx, 403, errors.New("invalid refresh token"), err.Error())
		return
	}
	resp, err := h.svc.NewAccessToken(ctx, claim.UserId)
	if err != nil {
		apiV1.HandleError(ctx, 500, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *AuthHandler) LoginBySessionId(ctx *gin.Context) {
	sessionId := ctx.Query("sessionId")
	fmt.Println(sessionId)
	resp, err := h.svc.GetAuthResponseBySessionId(ctx, sessionId)
	if err != nil {
		apiV1.HandleError(ctx, 401, err, err.Error())
		return
	}
	apiV1.HandleSuccess(ctx, resp)
}

func (h *AuthHandler) LinuxDoLogin(c *gin.Context) {
	cfg, err := h.systemConfigService.GetLinuxDoOAuthConfig(c)
	if err != nil {
		apiV1.HandleError(c, 403, err, "未查找到相关配置")
		return
	}
	cfg2, err := h.systemConfigService.GetRegisterConfig(c)
	if err != nil {
		apiV1.HandleError(c, 500, err, fmt.Errorf("系统错误❌:%w", err).Error())
		return
	}
	if !cfg2.AllowLinuxDoLogin {
		apiV1.HandleError(c, 403, errors.New("该登录方式已被禁用"), "该登录方式已被禁用")
		return
	}
	session := sessions.Default(c)
	oauthState := shortuuid.New()
	session.Set("oauth_state", oauthState)
	err = session.Save()
	if err != nil {
		apiV1.HandleError(c, 500, err, "未查找到相关配置")
		return
	}
	authorizationURL := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		oauth2.LinuxDoAuthorizationEndpoint, cfg.ClientId, "", oauthState)
	c.Redirect(http.StatusFound, authorizationURL)
}

func (h *AuthHandler) LinuxDoCallBack(ctx *gin.Context) {
	req := new(apiV1.OAuthCbRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	session := sessions.Default(ctx)
	oauthState := session.Get("oauth_state")
	if req.State != oauthState {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "invalid state")
		return
	}
	sessionId, err := h.svc.LinuxDoCallBack(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	redirectURI := getRedirectURI()
	url := redirectURI + sessionId
	ctx.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) GithubLogin(c *gin.Context) {
	_, err := h.systemConfigService.GetGithubOAuthConfig(c)
	if err != nil {
		apiV1.HandleError(c, 403, err, "未查找到相关配置")
		return
	}
	cfg2, err := h.systemConfigService.GetRegisterConfig(c)
	if err != nil {
		apiV1.HandleError(c, 500, err, fmt.Errorf("系统错误❌:%w", err).Error())
		return
	}
	if !cfg2.AllowGithubLogin {
		apiV1.HandleError(c, 403, errors.New("该登录方式已被禁用"), "该登录方式已被禁用")
		return
	}
	session := sessions.Default(c)
	oauthState := shortuuid.New()
	session.Set("oauth_state", oauthState)
	err = session.Save()
	if err != nil {
		apiV1.HandleError(c, 500, err, "未查找到相关配置")
		return
	}
	url, err := h.svc.GetGitHubRedirectURL(c, oauthState)
	if err != nil {
		apiV1.HandleError(c, 403, err, "未查找到相关配置")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GithubCallBack(ctx *gin.Context) {
	req := new(apiV1.OAuthCbRequest)
	if err := ctx.ShouldBind(req); err != nil {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, err.Error())
		return
	}
	session := sessions.Default(ctx)
	oauthState := session.Get("oauth_state")
	if req.State != oauthState {
		apiV1.HandleError(ctx, 400, apiV1.ErrBadRequest, "invalid state")
		return
	}
	sessionId, err := h.svc.GitHubCallBack(ctx, req)
	if err != nil {
		apiV1.HandleError(ctx, 0, err, err.Error())
		return
	}
	redirectURI := getRedirectURI()
	url := redirectURI + sessionId
	ctx.Redirect(http.StatusFound, url)
}
