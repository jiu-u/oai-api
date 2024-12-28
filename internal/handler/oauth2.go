package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	apiV1 "github.com/jiu-u/oai-api/api/v1"
	"github.com/jiu-u/oai-api/internal/service/oauth2"
	"github.com/jiu-u/oai-api/pkg/config"
	"github.com/lithammer/shortuuid/v4"
	"net/http"
)

const RedirectURI = "/#/login?sessionId="

type OAuth2Handler struct {
	*Handler
	Oauth2Svc oauth2.AuthService
	clientId  string
}

func NewOAuth2Handler(h *Handler, Oauth2Svc oauth2.AuthService, conf *config.Config) *OAuth2Handler {
	return &OAuth2Handler{
		Handler:   h,
		Oauth2Svc: Oauth2Svc,
		clientId:  conf.Oauth.LinuxDo.ClientId,
	}
}

func (h *OAuth2Handler) LinuxDoLogin(ctx *gin.Context) {
	url := "%"
	proto := "http"
	if ctx.Request.TLS != nil {
		proto = "https"
	}
	if ctx.Request.URL != nil {
		url = proto + "://" + ctx.Request.URL.Host + oauth2.LinuxDoRedirectURI
	}
	session := sessions.Default(ctx)
	oauthState := shortuuid.New()
	session.Set("oauth_state", oauthState)
	session.Save()
	authorizationURL := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		oauth2.LinuxDoAuthorizationEndpoint, h.clientId, url, oauthState)
	ctx.Redirect(http.StatusFound, authorizationURL)
}

func (h *OAuth2Handler) LinuxDoCallback(ctx *gin.Context) {
	req := new(apiV1.LinuxDoAuthRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(ctx)
	oauthState := session.Get("oauth_state")
	if req.State != oauthState {
		ctx.JSON(400, gin.H{"error": "invalid state"})
		return
	}
	sessionId, err := h.Oauth2Svc.LinuxDoAuthHandle(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	url := RedirectURI + sessionId
	ctx.Redirect(http.StatusFound, url)
}

func (h *OAuth2Handler) GetUserInfo(ctx *gin.Context) {
	sessionId := ctx.Query("sessionId")
	resp, err := h.Oauth2Svc.GetSessionUser(ctx, sessionId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)
}
