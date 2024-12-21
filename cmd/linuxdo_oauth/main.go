package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// OAuth2 参数
const (
	ClientID              = "hi3geJYfTotoiR5S62u3rh4W5tSeC5UG" // 修改为你的 Client Id
	ClientSecret          = "VMPBVoAfOB5ojkGXRDEtzvDhRLENHpaN" // 修改为你的 Client Secret
	RedirectURI           = "/oauth2/callback"                 // 修改为你的回调地址
	AuthorizationEndpoint = "https://connect.linux.do/oauth2/authorize"
	TokenEndpoint         = "https://connect.linux.do/oauth2/token"
	UserEndpoint          = "https://connect.linux.do/api/user"
)

func main() {
	r := gin.Default()

	// 设置 session 中间件
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	// 路由: 初始化 OAuth2 认证
	r.GET("/oauth2/initiate", func(c *gin.Context) {
		base := c.Request.URL
		fmt.Println("base", c.Request.URL, c.Request.Host, base.Path, base.Host, c.Request.Proto)
		session := sessions.Default(c)
		oauthState := generateRandomString(16)
		fmt.Println("state", oauthState)
		session.Set("oauth_state", oauthState)
		session.Save()

		authorizationURL := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
			AuthorizationEndpoint, ClientID, RedirectURI, oauthState)
		c.Redirect(http.StatusFound, authorizationURL)
	})

	// 路由: 处理 OAuth2 回调
	r.GET("/oauth2/callback", func(c *gin.Context) {
		session := sessions.Default(c)
		code := c.Query("code")
		state := c.Query("state")
		oauthState := session.Get("oauth_state")

		if oauthState == nil || state != oauthState {
			c.String(http.StatusUnauthorized, "State value does not match")
			return
		}

		// 请求 Access Token
		token, err := fetchAccessToken(code)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch access token: %v", err)
			return
		}

		// 请求用户信息
		userInfo, err := fetchUserInfo(token)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch user info: %v", err)
			return
		}

		c.JSON(http.StatusOK, userInfo)
	})

	r.Run(":8181")
}

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// 获取 Access Token
func fetchAccessToken(code string) (string, error) {
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, RedirectURI)
	req, err := http.NewRequest("POST", TokenEndpoint, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(ClientID, ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch access token: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse["access_token"].(string), nil
}

// 获取用户信息
func fetchUserInfo(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", UserEndpoint, nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
