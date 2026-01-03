package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"member-pre/internal/domain/auth"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
	"net/http"
	"time"
)

// RequestIDMiddleware 请求ID追踪中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 如果没有，生成一个新的请求ID
			requestID = uuid.New().String()
		}

		// 将请求ID存储到context中
		c.Set("request_id", requestID)

		// 将请求ID添加到响应头
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 获取请求ID
			requestID := getRequestID(c)

			// 记录错误日志
			log.Error("请求处理失败",
				logger.NewField("request_id", requestID),
				logger.NewField("method", c.Request.Method),
				logger.NewField("path", c.Request.URL.Path),
				logger.NewField("error", err.Error()),
			)

			// 处理应用错误
			if appErr, ok := errors.AsAppError(err); ok {
				utils.ErrorWithCode(c, appErr.HTTPStatus(), appErr.Message)
				return
			}

			// 处理未知错误
			utils.ErrorWithCode(c, http.StatusInternalServerError, "内部服务器错误")
		}
	}
}

// LoggerMiddleware 日志中间件（替换Gin默认的Logger）
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := getRequestID(c)

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		// 记录请求日志
		log.Info("HTTP请求",
			logger.NewField("request_id", requestID),
			logger.NewField("method", c.Request.Method),
			logger.NewField("path", c.Request.URL.Path),
			logger.NewField("status_code", statusCode),
			logger.NewField("duration_ms", duration.Milliseconds()),
			logger.NewField("client_ip", c.ClientIP()),
			logger.NewField("user_agent", c.Request.UserAgent()),
		)
	}
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(authService *auth.AuthService, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := getRequestID(c)

		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("未提供认证令牌",
				logger.NewField("request_id", requestID),
				logger.NewField("path", c.Request.URL.Path),
			)
			utils.ErrorWithCode(c, http.StatusUnauthorized, "未提供认证令牌")
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// 如果没有Bearer前缀，尝试直接使用
			tokenString = authHeader
		}

		// 验证token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Warn("令牌验证失败",
				logger.NewField("request_id", requestID),
				logger.NewField("path", c.Request.URL.Path),
				logger.NewField("error", err.Error()),
			)
			// 使用统一的错误处理
			utils.Error(c, err)
			c.Abort()
			return
		}

		// 将用户信息存储到context中
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		if claims.StoreID != nil {
			c.Set("store_id", *claims.StoreID)
		}

		log.Debug("认证成功",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", claims.UserID),
			logger.NewField("role", claims.Role),
		)

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorWithCode(c, http.StatusUnauthorized, "未授权访问")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			utils.ErrorWithCode(c, http.StatusForbidden, "无效的用户角色")
			c.Abort()
			return
		}

		// 如果没有指定允许的角色，则允许所有已认证用户
		if len(allowedRoles) == 0 {
			c.Next()
			return
		}

		// 检查角色是否在允许列表中
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorWithCode(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// getRequestID 从context中获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return "unknown"
}
