package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/user"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *user.UserService
	logger  logger.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(service *user.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  log,
	}
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} user.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的用户ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	ctx := c.Request.Context()
	u, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, u)
}

// UpdateWorkStatus 更新美甲师工作状态
// @Summary 更新美甲师工作状态
// @Description 更新美甲师的工作状态（在岗/休息/离岗）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UpdateWorkStatusRequest true "工作状态请求"
// @Success 200 {object} map[string]string
// @Router /users/{id}/work-status [put]
func (h *UserHandler) UpdateWorkStatus(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的用户ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req UpdateWorkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新工作状态请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.UpdateWorkStatus(ctx, uint(id), req.WorkStatus); err != nil {
		h.logger.Error("更新工作状态失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", id),
			logger.NewField("work_status", req.WorkStatus),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新工作状态成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", id),
		logger.NewField("work_status", req.WorkStatus),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "更新工作状态成功", nil)
}

// UpdateWorkStatusRequest 更新工作状态请求
type UpdateWorkStatusRequest struct {
	WorkStatus string `json:"work_status" binding:"required,oneof=working rest offline" example:"working"`
}

// GetUsersByStore 根据门店获取用户列表
// @Summary 根据门店获取用户列表
// @Description 根据门店ID获取该门店的用户列表（店长和美甲师）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id path int true "门店ID"
// @Param role query string false "角色筛选" Enums(store_manager, technician)
// @Success 200 {object} []user.User
// @Router /stores/{store_id}/users [get]
func (h *UserHandler) GetUsersByStore(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	storeIDStr := c.Param("store_id")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", storeIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	role := c.Query("role")

	ctx := c.Request.Context()
	users, err := h.service.GetByStoreID(ctx, uint(storeID), role)
	if err != nil {
		h.logger.Error("获取门店用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取门店用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", storeID),
		logger.NewField("count", len(users)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, users)
}

