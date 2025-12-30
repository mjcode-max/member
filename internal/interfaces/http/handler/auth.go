package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/auth"
	"member-pre/pkg/errors"
	"member-pre/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, errors.ErrInvalidParams("参数错误: "+err.Error()))
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		utils.Error(c, errors.ErrUnauthorized(err.Error()))
		return
	}

	utils.Success(c, resp)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取token
	token := extractToken(c)
	if token == "" {
		utils.SuccessWithMessage(c, "登出成功", nil)
		return
	}

	if err := h.authService.Logout(token); err != nil {
		utils.Error(c, errors.ErrInternal("登出失败"))
		return
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.Error(c, errors.ErrUnauthorized("未登录"))
		return
	}

	utils.Success(c, user)
}

// extractToken 从请求中提取token（与middleware中的相同逻辑）
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return c.Query("token")
}
