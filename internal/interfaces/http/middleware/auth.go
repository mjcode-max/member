package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/auth"
	"member-pre/pkg/errors"
	"member-pre/pkg/utils"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService *auth.Service
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := extractToken(c)
		if token == "" {
			utils.Error(c, errors.ErrUnauthorized("请先登录"))
			c.Abort()
			return
		}

		// 验证token
		user, err := m.authService.ValidateToken(token)
		if err != nil {
			utils.Error(c, errors.ErrUnauthorized("token无效或已过期"))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行认证
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 检查角色
		role, exists := c.Get("role")
		if !exists {
			utils.Error(c, errors.ErrForbidden("权限不足"))
			c.Abort()
			return
		}

		roleStr := role.(string)
		hasRole := false
		for _, r := range roles {
			if roleStr == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.Error(c, errors.ErrForbidden("权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken 从请求中提取token
func extractToken(c *gin.Context) string {
	// 从Header中获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 从Query参数中获取
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 从Cookie中获取
	token, _ = c.Cookie("token")
	return token
}
