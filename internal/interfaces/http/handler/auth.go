package handler

import (
	"member-pre/internal/domain/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	service *auth.AuthService
	logger  logger.Logger
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(service *auth.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  log,
	}
}

// Login 登录
// @Summary 用户登录
// @Description 用户登录接口，支持用户名或手机号登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "登录请求"
// @Success 200 {object} domain.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("登录请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	h.logger.Info("收到登录请求",
		logger.NewField("request_id", requestID),
		logger.NewField("username", req.Username),
	)

	ctx := c.Request.Context()
	resp, err := h.service.Login(ctx, &req)
	if err != nil {
		h.logger.Error("登录失败",
			logger.NewField("request_id", requestID),
			logger.NewField("username", req.Username),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("登录成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", resp.User.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, resp)
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.GetCurrentUserResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 从 context 中获取用户ID（通常由认证中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("获取当前用户失败：未授权",
			logger.NewField("request_id", requestID),
		)
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 转换用户ID类型
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case int:
		uid = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			h.logger.Warn("获取当前用户失败：无效的用户ID",
				logger.NewField("request_id", requestID),
				logger.NewField("user_id", v),
				logger.NewField("error", err.Error()),
			)
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
			return
		}
		uid = uint(id)
	default:
		h.logger.Warn("获取当前用户失败：无效的用户ID类型",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id_type", v),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.GetCurrentUser(ctx, uid)
	if err != nil {
		h.logger.Error("获取当前用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", uid),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取当前用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", uid),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, &auth.GetCurrentUserResponse{
		User: user,
	})
}

// Logout 登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 从请求头获取 token
	token := c.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		// 移除 "Bearer " 前缀
		token = token[7:]
	}

	h.logger.Info("收到登出请求",
		logger.NewField("request_id", requestID),
	)

	ctx := c.Request.Context()
	if err := h.service.Logout(ctx, token); err != nil {
		h.logger.Error("登出失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("登出成功",
		logger.NewField("request_id", requestID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "登出成功", nil)
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
