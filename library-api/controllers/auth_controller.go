package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/example/library-api/config"
	"github.com/example/library-api/database"
	"github.com/example/library-api/middleware"
	"github.com/example/library-api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func init() {
	// 修复：添加配置加载错误处理
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin 重定向到Google登录页面
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleLoginCallback 处理Google登录回调
func GoogleLoginCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	// 修复：正确初始化OAuth2服务
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	reqRes, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	// 修复：解析用户信息
	userInfoBytes, err := io.ReadAll(reqRes.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info: " + err.Error()})
		return
	}
	// 定义临时的用户信息结构体，用于解析Google返回的用户信息
	type UserInfo struct {
		Id    string `json:"sub"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var userInfo UserInfo
	if err := json.Unmarshal(userInfoBytes, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info: " + err.Error()})
		return
	}

	// 查找或创建用户
	var user models.User
	result := database.DB.Where("google_id = ? OR email = ?", userInfo.Id, userInfo.Email).First(&user)

	if result.Error != nil {
		// 创建新用户
		user = models.User{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			GoogleID: userInfo.Id,
			Role:     models.RoleUser,
		}
		// 修复：添加数据库错误处理
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}
	}
	// 加载配置
	cfg, err := config.LoadConfig()

	// 生成JWT令牌
	expirationTime := time.Now().Add(time.Duration(cfg.JWTExpiryHours) * time.Hour)
	jwtToken, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      jwtToken,
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"expires_at": expirationTime,
	})
}
