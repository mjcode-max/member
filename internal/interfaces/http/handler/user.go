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

	// 获取当前用户信息
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	var currentUserID uint
	switch v := userIDInterface.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
			return
		}
		currentUserID = uint(id)
	default:
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的用户ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	ctx := c.Request.Context()
	u, err := h.service.GetByID(ctx, uint(targetID))
	if err != nil {
		h.logger.Error("获取用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", targetID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if u == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 权限检查
	if userRole == user.RoleAdmin {
		// 后台可以查看所有用户
	} else if userRole == user.RoleStoreManager {
		// 店长只能查看自己门店的用户
		storeIDInterface, exists := c.Get("store_id")
		if !exists {
			utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
			return
		}

		var storeID uint
		switch v := storeIDInterface.(type) {
		case uint:
			storeID = v
		case int:
			storeID = uint(v)
		case *uint:
			if v == nil {
				utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
				return
			}
			storeID = *v
		default:
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID类型")
			return
		}

		if u.StoreID == nil || *u.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权查看该用户")
			return
		}
	} else if userRole == user.RoleTechnician || userRole == user.RoleCustomer {
		// 美甲师和顾客只能查看自己的信息
		if currentUserID != uint(targetID) {
			utils.ErrorWithCode(c, http.StatusForbidden, "只能查看自己的信息")
			return
		}
	}

	h.logger.Info("获取用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", targetID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, u)
}

// UpdateWorkStatus 更新美甲师工作状态
// @Summary 更新美甲师工作状态
// @Description 更新美甲师的工作状态（在岗/休息/离岗）。美甲师可以更新自己的状态，店长可以更新自己门店员工的状态，后台可以更新任何美甲师的状态
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

	// 获取当前用户信息
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	var currentUserID uint
	switch v := userIDInterface.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
			return
		}
		currentUserID = uint(id)
	default:
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
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

	// 权限检查
	if userRole == user.RoleTechnician {
		// 美甲师只能更新自己的状态
		if currentUserID != uint(targetID) {
			utils.ErrorWithCode(c, http.StatusForbidden, "只能更新自己的工作状态")
			return
		}
	} else if userRole == user.RoleStoreManager {
		// 店长只能更新自己门店员工的状态
		storeIDInterface, exists := c.Get("store_id")
		if !exists {
			utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
			return
		}

		var storeID uint
		switch v := storeIDInterface.(type) {
		case uint:
			storeID = v
		case int:
			storeID = uint(v)
		case *uint:
			if v == nil {
				utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
				return
			}
			storeID = *v
		default:
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID类型")
			return
		}

		// 验证目标用户是否属于当前门店
		targetUser, err := h.service.GetByID(ctx, uint(targetID))
		if err != nil {
			utils.Error(c, err)
			return
		}
		if targetUser == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "用户不存在")
			return
		}
		if targetUser.StoreID == nil || *targetUser.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权操作该员工")
			return
		}
	} else if userRole != user.RoleAdmin {
		// 只有后台、店长、美甲师可以更新工作状态
		utils.ErrorWithCode(c, http.StatusForbidden, "权限不足")
		return
	}

	// 更新工作状态
	if err := h.service.UpdateWorkStatus(ctx, uint(targetID), req.WorkStatus); err != nil {
		h.logger.Error("更新工作状态失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", targetID),
			logger.NewField("work_status", req.WorkStatus),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新工作状态成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", targetID),
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

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表，支持按角色、状态、门店筛选和分页
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role query string false "角色筛选" Enums(admin, store_manager, technician, customer)
// @Param status query string false "状态筛选" Enums(active, inactive)
// @Param store_id query int false "门店ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} utils.PaginationResponse
// @Router /users [get]
func (h *UserHandler) GetUserList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 获取当前用户角色
	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	// 获取查询参数
	role := c.Query("role")
	status := c.Query("status")
	storeIDStr := c.Query("store_id")
	username := c.Query("username")
	phone := c.Query("phone")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var storeID *uint
	// 店长只能查看自己门店的员工
	if userRole == user.RoleStoreManager {
		storeIDInterface, exists := c.Get("store_id")
		if !exists {
			utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
			return
		}
		switch v := storeIDInterface.(type) {
		case uint:
			storeID = &v
		case int:
			uid := uint(v)
			storeID = &uid
		case *uint:
			storeID = v
		}
		// 店长只能查看美甲师
		if role == "" {
			role = user.RoleTechnician
		} else if role != user.RoleTechnician {
			utils.ErrorWithCode(c, http.StatusForbidden, "只能查看美甲师员工")
			return
		}
	} else if storeIDStr != "" {
		// 后台可以指定门店ID
		id, err := strconv.ParseUint(storeIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			storeID = &uid
		}
	}

	ctx := c.Request.Context()
	users, total, err := h.service.GetList(ctx, role, status, storeID, username, phone, page, pageSize)
	if err != nil {
		h.logger.Error("获取用户列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取用户列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(users)),
		logger.NewField("total", total),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithPagination(c, users, page, pageSize, total)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户（店长、美甲师等）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "创建用户请求"
// @Success 200 {object} user.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 获取当前用户角色
	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建用户请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 店长只能创建美甲师
	if userRole == user.RoleStoreManager {
		if req.Role != "" && req.Role != user.RoleTechnician {
			utils.ErrorWithCode(c, http.StatusForbidden, "只能创建美甲师员工")
			return
		}
		req.Role = user.RoleTechnician // 强制设置为美甲师

		// 从context获取当前用户的store_id
		storeIDInterface, exists := c.Get("store_id")
		if !exists {
			utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
			return
		}
		switch v := storeIDInterface.(type) {
		case uint:
			req.StoreID = &v
		case int:
			uid := uint(v)
			req.StoreID = &uid
		case *uint:
			req.StoreID = v
		}
	}

	// 验证必填字段
	if req.Role == user.RoleCustomer {
		// 顾客必须提供手机号
		if req.Phone == "" {
			utils.ErrorWithCode(c, http.StatusBadRequest, "顾客必须提供手机号")
			return
		}
	} else {
		// 员工必须提供用户名和密码
		if req.Username == "" {
			utils.ErrorWithCode(c, http.StatusBadRequest, "员工必须提供用户名")
			return
		}
		if req.Password == "" {
			utils.ErrorWithCode(c, http.StatusBadRequest, "员工必须提供密码")
			return
		}
	}

	// 转换为领域实体
	newUser := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role,
		Status:   req.Status,
		StoreID:  req.StoreID,
	}

	ctx := c.Request.Context()
	if err := h.service.Create(ctx, newUser); err != nil {
		h.logger.Error("创建用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", newUser.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, newUser)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UpdateUserRequest true "更新用户请求"
// @Success 200 {object} user.User
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 获取当前用户信息
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	var currentUserID uint
	switch v := userIDInterface.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
			return
		}
		currentUserID = uint(id)
	default:
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的用户ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 权限检查
	if userRole == user.RoleAdmin {
		// 后台可以更新所有用户
	} else if userRole == user.RoleStoreManager {
		// 店长只能更新自己门店的用户
		storeIDInterface, exists := c.Get("store_id")
		if !exists {
			utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
			return
		}

		var storeID uint
		switch v := storeIDInterface.(type) {
		case uint:
			storeID = v
		case int:
			storeID = uint(v)
		case *uint:
			if v == nil {
				utils.ErrorWithCode(c, http.StatusForbidden, "当前用户未关联门店")
				return
			}
			storeID = *v
		default:
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID类型")
			return
		}

		// 验证目标用户是否属于当前门店
		targetUser, err := h.service.GetByID(c.Request.Context(), uint(targetID))
		if err != nil {
			utils.Error(c, err)
			return
		}
		if targetUser == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "用户不存在")
			return
		}
		if targetUser.StoreID == nil || *targetUser.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权操作该用户")
			return
		}
	} else if userRole == user.RoleTechnician || userRole == user.RoleCustomer {
		// 美甲师和顾客只能更新自己的信息
		if currentUserID != uint(targetID) {
			utils.ErrorWithCode(c, http.StatusForbidden, "只能更新自己的信息")
			return
		}
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新用户请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 美甲师和顾客不能修改角色、状态等敏感字段
	if userRole == user.RoleTechnician || userRole == user.RoleCustomer {
		// 获取当前用户信息，确保不能修改角色和状态
		currentUser, err := h.service.GetByID(c.Request.Context(), currentUserID)
		if err != nil {
			utils.Error(c, err)
			return
		}
		if currentUser == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "用户不存在")
			return
		}

		if userRole == user.RoleTechnician {
			// 美甲师不能修改角色、状态、门店ID
			if req.Role != "" && req.Role != currentUser.Role {
				utils.ErrorWithCode(c, http.StatusForbidden, "不能修改角色")
				return
			}
			if req.Status != "" && req.Status != currentUser.Status {
				utils.ErrorWithCode(c, http.StatusForbidden, "不能修改状态")
				return
			}
			if req.StoreID != nil {
				utils.ErrorWithCode(c, http.StatusForbidden, "不能修改门店")
				return
			}
		} else if userRole == user.RoleCustomer {
			// 顾客不能修改角色、状态
			if req.Role != "" {
				utils.ErrorWithCode(c, http.StatusForbidden, "不能修改角色")
				return
			}
			if req.Status != "" {
				utils.ErrorWithCode(c, http.StatusForbidden, "不能修改状态")
				return
			}
		}
	}

	// 转换为领域实体
	updateUser := &user.User{
		ID:         uint(targetID),
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   req.Password,
		Role:       req.Role,
		Status:     req.Status,
		StoreID:    req.StoreID,
		WorkStatus: req.WorkStatus,
	}

	ctx := c.Request.Context()
	if err := h.service.Update(ctx, updateUser); err != nil {
		h.logger.Error("更新用户失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", targetID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	// 重新获取用户信息
	updatedUser, err := h.service.GetByID(ctx, uint(targetID))
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新用户成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", targetID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, updatedUser)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 启用或禁用用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UpdateUserStatusRequest true "更新用户状态请求"
// @Success 200 {object} map[string]string
// @Router /users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
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

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新用户状态请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.UpdateStatus(ctx, uint(id), req.Status); err != nil {
		h.logger.Error("更新用户状态失败",
			logger.NewField("request_id", requestID),
			logger.NewField("user_id", id),
			logger.NewField("status", req.Status),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新用户状态成功",
		logger.NewField("request_id", requestID),
		logger.NewField("user_id", id),
		logger.NewField("status", req.Status),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "更新用户状态成功", nil)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username"`                                                              // 用户名（员工必填，顾客可选）
	Email    string `json:"email"`                                                                 // 邮箱
	Phone    string `json:"phone"`                                                                 // 手机号（顾客必填）
	Password string `json:"password"`                                                              // 密码（员工必填，顾客不需要）
	Role     string `json:"role" binding:"required,oneof=admin store_manager technician customer"` // 角色
	Status   string `json:"status"`                                                                // 状态: active, inactive
	StoreID  *uint  `json:"store_id"`                                                              // 门店ID（店长和美甲师必填）
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username   string  `json:"username"`    // 用户名
	Email      string  `json:"email"`       // 邮箱
	Phone      string  `json:"phone"`       // 手机号
	Password   string  `json:"password"`    // 密码（如果提供则更新）
	Role       string  `json:"role"`        // 角色
	Status     string  `json:"status"`      // 状态
	StoreID    *uint   `json:"store_id"`    // 门店ID
	WorkStatus *string `json:"work_status"` // 工作状态（美甲师）
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"` // 状态
}
