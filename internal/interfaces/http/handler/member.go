package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/member"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// MemberHandler 会员处理器
type MemberHandler struct {
	memberService *member.MemberService
	usageService  *member.UsageService
	storeService  *store.StoreService
	userService   *user.UserService
	logger        logger.Logger
}

// NewMemberHandler 创建会员处理器
func NewMemberHandler(
	memberService *member.MemberService,
	usageService *member.UsageService,
	storeService *store.StoreService,
	userService *user.UserService,
	log logger.Logger,
) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
		usageService:  usageService,
		storeService:  storeService,
		userService:   userService,
		logger:        log,
	}
}

// CreateMember 创建会员
// @Summary 创建会员
// @Description 创建会员（管理员、店长），包含套餐信息
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMemberRequest true "创建会员请求"
// @Success 200 {object} member.Member
// @Router /members [post]
func (h *MemberHandler) CreateMember(c *gin.Context) {
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
	default:
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建会员请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 权限检查：店长只能创建自己门店的会员
	if userRole == user.RoleStoreManager {
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

		if req.StoreID != nil && *req.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权创建其他门店的会员")
			return
		}

		// 店长创建会员时，自动设置为自己的门店
		if req.StoreID == nil {
			req.StoreID = &storeID
		}
	}

	// 验证门店是否存在
	if req.StoreID != nil {
		ctx := c.Request.Context()
		storeInfo, err := h.storeService.GetByID(ctx, *req.StoreID)
		if err != nil {
			h.logger.Error("查找门店失败",
				logger.NewField("request_id", requestID),
				logger.NewField("store_id", *req.StoreID),
				logger.NewField("error", err.Error()),
			)
			utils.Error(c, err)
			return
		}
		if storeInfo == nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "门店不存在")
			return
		}
	} else {
		utils.ErrorWithCode(c, http.StatusBadRequest, "门店ID不能为空")
		return
	}

	// 转换为领域实体
	newMember := &member.Member{
		Name:            req.Name,
		Phone:           req.Phone,
		PackageName:     req.PackageName,
		ServiceType:     req.ServiceType,
		Price:           req.Price,
		ValidityDuration: req.ValidityDuration,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
		StoreID:         *req.StoreID,
		PurchaseAmount:  req.PurchaseAmount,
		Status:          req.Status,
		Description:     req.Description,
		CreatedBy:       currentUserID,
	}

	ctx := c.Request.Context()
	if err := h.memberService.CreateMember(ctx, newMember); err != nil {
		h.logger.Error("创建会员失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建会员成功",
		logger.NewField("request_id", requestID),
		logger.NewField("member_id", newMember.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, newMember)
}

// GetMember 获取会员详情
// @Summary 获取会员详情
// @Description 根据会员ID获取会员详细信息
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {object} member.Member
// @Router /members/{id} [get]
func (h *MemberHandler) GetMember(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的会员ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	ctx := c.Request.Context()
	m, err := h.memberService.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取会员失败",
			logger.NewField("request_id", requestID),
			logger.NewField("member_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if m == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "会员不存在")
		return
	}

	// 权限检查：店长只能查看自己门店的会员
	if userRole == user.RoleStoreManager {
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

		if m.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权查看该会员")
			return
		}
	}

	h.logger.Info("获取会员成功",
		logger.NewField("request_id", requestID),
		logger.NewField("member_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, m)
}

// GetMemberList 获取会员列表
// @Summary 获取会员列表
// @Description 获取会员列表，支持多条件筛选和分页
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "会员姓名（模糊搜索）"
// @Param phone query string false "手机号（模糊搜索）"
// @Param store_id query int false "门店ID"
// @Param status query string false "状态" Enums(active, expired, inactive)
// @Param service_type query string false "服务类型" Enums(nail, eyelash, combo)
// @Param package_name query string false "套餐名称（模糊搜索）"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} utils.PaginationResponse
// @Router /members [get]
func (h *MemberHandler) GetMemberList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	// 获取查询参数
	name := c.Query("name")
	phone := c.Query("phone")
	status := c.Query("status")
	serviceType := c.Query("service_type")
	packageName := c.Query("package_name")

	var storeID *uint
	if userRole == user.RoleStoreManager {
		// 店长只能查询自己门店的会员
		storeIDInterface, exists := c.Get("store_id")
		if exists {
			switch v := storeIDInterface.(type) {
			case uint:
				storeID = &v
			case int:
				id := uint(v)
				storeID = &id
			case *uint:
				storeID = v
			}
		}
	} else {
		// 管理员可以查询所有门店或指定门店
		storeIDStr := c.Query("store_id")
		if storeIDStr != "" {
			id, err := strconv.ParseUint(storeIDStr, 10, 32)
			if err == nil {
				idUint := uint(id)
				storeID = &idUint
			}
		}
	}

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

	ctx := c.Request.Context()
	members, total, err := h.memberService.GetList(ctx, name, phone, storeID, status, serviceType, packageName, page, pageSize)
	if err != nil {
		h.logger.Error("获取会员列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取会员列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(members)),
		logger.NewField("total", total),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithPagination(c, members, page, pageSize, total)
}

// GetMemberByPhone 根据手机号查询会员列表
// @Summary 根据手机号查询会员
// @Description 根据手机号查询会员列表（同一手机号可能有多个会员记录）
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param phone path string true "手机号"
// @Success 200 {array} member.Member
// @Router /members/phone/{phone} [get]
func (h *MemberHandler) GetMemberByPhone(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	phone := c.Param("phone")
	if phone == "" {
		utils.ErrorWithCode(c, http.StatusBadRequest, "手机号不能为空")
		return
	}

	ctx := c.Request.Context()
	members, err := h.memberService.GetMemberByPhone(ctx, phone)
	if err != nil {
		h.logger.Error("根据手机号查询会员失败",
			logger.NewField("request_id", requestID),
			logger.NewField("phone", phone),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("根据手机号查询会员成功",
		logger.NewField("request_id", requestID),
		logger.NewField("phone", phone),
		logger.NewField("count", len(members)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, members)
}

// UpdateMember 更新会员信息
// @Summary 更新会员信息
// @Description 更新会员信息（管理员、店长），可更新套餐信息
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Param request body UpdateMemberRequest true "更新会员请求"
// @Success 200 {object} member.Member
// @Router /members/{id} [put]
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的会员ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	// 权限检查：店长只能更新自己门店的会员
	if userRole == user.RoleStoreManager {
		ctx := c.Request.Context()
		m, err := h.memberService.GetByID(ctx, uint(id))
		if err != nil {
			utils.Error(c, err)
			return
		}
		if m == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "会员不存在")
			return
		}

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

		if m.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权更新该会员")
			return
		}
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新会员请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 转换为领域实体
	updateMember := &member.Member{
		ID:              uint(id),
		Name:            req.Name,
		Phone:           req.Phone,
		PackageName:     req.PackageName,
		ServiceType:     req.ServiceType,
		Price:           req.Price,
		ValidityDuration: req.ValidityDuration,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
		PurchaseAmount:  req.PurchaseAmount,
		Status:          req.Status,
		Description:     req.Description,
	}

	ctx := c.Request.Context()
	if err := h.memberService.UpdateMember(ctx, updateMember); err != nil {
		h.logger.Error("更新会员失败",
			logger.NewField("request_id", requestID),
			logger.NewField("member_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	// 重新获取更新后的会员信息
	updatedMember, err := h.memberService.GetByID(ctx, uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新会员成功",
		logger.NewField("request_id", requestID),
		logger.NewField("member_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, updatedMember)
}

// UpdateMemberStatus 更新会员状态
// @Summary 更新会员状态
// @Description 更新会员状态（管理员、店长）
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Param request body UpdateMemberStatusRequest true "更新状态请求"
// @Success 200 {object} member.Member
// @Router /members/{id}/status [put]
func (h *MemberHandler) UpdateMemberStatus(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的会员ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	// 权限检查：店长只能更新自己门店的会员状态
	if userRole == user.RoleStoreManager {
		ctx := c.Request.Context()
		m, err := h.memberService.GetByID(ctx, uint(id))
		if err != nil {
			utils.Error(c, err)
			return
		}
		if m == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "会员不存在")
			return
		}

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

		if m.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权更新该会员状态")
			return
		}
	}

	var req UpdateMemberStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新会员状态请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.memberService.UpdateMemberStatus(ctx, uint(id), req.Status); err != nil {
		h.logger.Error("更新会员状态失败",
			logger.NewField("request_id", requestID),
			logger.NewField("member_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	// 重新获取更新后的会员信息
	updatedMember, err := h.memberService.GetByID(ctx, uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新会员状态成功",
		logger.NewField("request_id", requestID),
		logger.NewField("member_id", id),
		logger.NewField("status", req.Status),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, updatedMember)
}

// CreateUsage 记录使用
// @Summary 记录使用
// @Description 记录会员使用（管理员、店长、美甲师），记录使用历史并递增已使用次数
// @Tags 使用记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Param request body CreateUsageRequest true "创建使用记录请求"
// @Success 200 {object} member.MemberUsage
// @Router /members/{id}/usages [post]
func (h *MemberHandler) CreateUsage(c *gin.Context) {
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
	default:
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的用户ID类型")
		return
	}

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	memberIDStr := c.Param("id")
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的会员ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", memberIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	var req CreateUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建使用记录请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 权限检查：店长和美甲师只能记录自己门店的使用
	if userRole == user.RoleStoreManager || userRole == user.RoleTechnician {
		ctx := c.Request.Context()
		m, err := h.memberService.GetByID(ctx, uint(memberID))
		if err != nil {
			utils.Error(c, err)
			return
		}
		if m == nil {
			utils.ErrorWithCode(c, http.StatusNotFound, "会员不存在")
			return
		}

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

		if m.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权记录其他门店的会员使用")
			return
		}

		// 店长和美甲师创建使用记录时，自动设置为自己的门店
		if req.StoreID == nil {
			req.StoreID = &storeID
		} else if *req.StoreID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权记录其他门店的会员使用")
			return
		}
	}

	// 验证门店是否存在并获取门店名称
	var storeName string
	if req.StoreID != nil {
		ctx := c.Request.Context()
		storeInfo, err := h.storeService.GetByID(ctx, *req.StoreID)
		if err != nil {
			h.logger.Error("查找门店失败",
				logger.NewField("request_id", requestID),
				logger.NewField("store_id", *req.StoreID),
				logger.NewField("error", err.Error()),
			)
			utils.Error(c, err)
			return
		}
		if storeInfo == nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "门店不存在")
			return
		}
		storeName = storeInfo.Name
	} else {
		utils.ErrorWithCode(c, http.StatusBadRequest, "门店ID不能为空")
		return
	}

	// 验证美甲师是否存在并获取美甲师姓名
	var technicianName string
	if req.TechnicianID != nil {
		ctx := c.Request.Context()
		technician, err := h.userService.GetByID(ctx, *req.TechnicianID)
		if err != nil {
			h.logger.Error("查找美甲师失败",
				logger.NewField("request_id", requestID),
				logger.NewField("technician_id", *req.TechnicianID),
				logger.NewField("error", err.Error()),
			)
			utils.Error(c, err)
			return
		}
		if technician == nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "美甲师不存在")
			return
		}
		technicianName = technician.Username
	}

	// 获取会员信息以填充套餐名称
	ctx := c.Request.Context()
	memberInfo, err := h.memberService.GetByID(ctx, uint(memberID))
	if err != nil {
		utils.Error(c, err)
		return
	}
	if memberInfo == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "会员不存在")
		return
	}

	// 转换为领域实体
	newUsage := &member.MemberUsage{
		MemberID:       uint(memberID),
		PackageName:    memberInfo.PackageName,
		ServiceItem:    req.ServiceItem,
		StoreID:        *req.StoreID,
		StoreName:      storeName,
		TechnicianID:   req.TechnicianID,
		TechnicianName: technicianName,
		UsageDate:      req.UsageDate,
		Remark:         req.Remark,
		CreatedBy:      currentUserID,
	}

	if err := h.usageService.CreateUsage(ctx, newUsage); err != nil {
		h.logger.Error("创建使用记录失败",
			logger.NewField("request_id", requestID),
			logger.NewField("member_id", memberID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建使用记录成功",
		logger.NewField("request_id", requestID),
		logger.NewField("usage_id", newUsage.ID),
		logger.NewField("member_id", memberID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, newUsage)
}

// GetUsageList 获取使用记录列表
// @Summary 获取使用记录列表
// @Description 获取使用记录列表（不分页，管理员、店长、美甲师），支持筛选
// @Tags 使用记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param member_id query int false "会员ID"
// @Param store_id query int false "门店ID"
// @Param technician_id query int false "美甲师ID"
// @Success 200 {array} member.MemberUsage
// @Router /usages [get]
func (h *MemberHandler) GetUsageList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	var memberID *uint
	var storeID *uint
	var technicianID *uint

	// 权限检查：店长和美甲师只能查询自己门店的使用记录
	if userRole == user.RoleStoreManager || userRole == user.RoleTechnician {
		storeIDInterface, exists := c.Get("store_id")
		if exists {
			switch v := storeIDInterface.(type) {
			case uint:
				storeID = &v
			case int:
				id := uint(v)
				storeID = &id
			case *uint:
				storeID = v
			}
		}

		// 美甲师只能查询自己的使用记录
		if userRole == user.RoleTechnician {
			userIDInterface, exists := c.Get("user_id")
			if exists {
				switch v := userIDInterface.(type) {
				case uint:
					technicianID = &v
				case int:
					id := uint(v)
					technicianID = &id
				}
			}
		}
	} else {
		// 管理员可以查询所有使用记录或指定筛选条件
		memberIDStr := c.Query("member_id")
		if memberIDStr != "" {
			id, err := strconv.ParseUint(memberIDStr, 10, 32)
			if err == nil {
				idUint := uint(id)
				memberID = &idUint
			}
		}

		storeIDStr := c.Query("store_id")
		if storeIDStr != "" {
			id, err := strconv.ParseUint(storeIDStr, 10, 32)
			if err == nil {
				idUint := uint(id)
				storeID = &idUint
			}
		}

		technicianIDStr := c.Query("technician_id")
		if technicianIDStr != "" {
			id, err := strconv.ParseUint(technicianIDStr, 10, 32)
			if err == nil {
				idUint := uint(id)
				technicianID = &idUint
			}
		}
	}

	ctx := c.Request.Context()
	usages, err := h.usageService.GetUsageList(ctx, memberID, storeID, technicianID)
	if err != nil {
		h.logger.Error("获取使用记录列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取使用记录列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(usages)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, usages)
}

// GetUsageByMemberID 获取会员使用记录列表
// @Summary 获取会员使用记录列表
// @Description 根据会员ID获取使用记录列表（不分页）
// @Tags 使用记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {array} member.MemberUsage
// @Router /members/{id}/usages [get]
func (h *MemberHandler) GetUsageByMemberID(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	memberIDStr := c.Param("id")
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的会员ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", memberIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	ctx := c.Request.Context()
	usages, err := h.usageService.GetUsageByMemberID(ctx, uint(memberID))
	if err != nil {
		h.logger.Error("获取会员使用记录列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("member_id", memberID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取会员使用记录列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("member_id", memberID),
		logger.NewField("count", len(usages)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, usages)
}

// GetUsage 获取使用记录详情
// @Summary 获取使用记录详情
// @Description 根据使用记录ID获取详细信息
// @Tags 使用记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "使用记录ID"
// @Success 200 {object} member.MemberUsage
// @Router /usages/{id} [get]
func (h *MemberHandler) GetUsage(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的使用记录ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的使用记录ID")
		return
	}

	ctx := c.Request.Context()
	usage, err := h.usageService.GetUsageByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取使用记录失败",
			logger.NewField("request_id", requestID),
			logger.NewField("usage_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if usage == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "使用记录不存在")
		return
	}

	h.logger.Info("获取使用记录成功",
		logger.NewField("request_id", requestID),
		logger.NewField("usage_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, usage)
}

// DeleteUsage 删除使用记录
// @Summary 删除使用记录
// @Description 删除使用记录（仅管理员），删除时需回退会员的已使用次数
// @Tags 使用记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "使用记录ID"
// @Success 200 {string} string "删除成功"
// @Router /usages/{id} [delete]
func (h *MemberHandler) DeleteUsage(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	// 权限检查：仅管理员可删除
	if userRole != user.RoleAdmin {
		utils.ErrorWithCode(c, http.StatusForbidden, "无权删除使用记录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的使用记录ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的使用记录ID")
		return
	}

	ctx := c.Request.Context()
	if err := h.usageService.DeleteUsage(ctx, uint(id)); err != nil {
		h.logger.Error("删除使用记录失败",
			logger.NewField("request_id", requestID),
			logger.NewField("usage_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("删除使用记录成功",
		logger.NewField("request_id", requestID),
		logger.NewField("usage_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, "删除成功")
}

// 请求结构体定义

// CreateMemberRequest 创建会员请求
type CreateMemberRequest struct {
	Name            string     `json:"name" binding:"required"`              // 会员姓名
	Phone           string     `json:"phone" binding:"required"`            // 手机号
	PackageName     string     `json:"package_name" binding:"required"`      // 套餐名称
	ServiceType     string     `json:"service_type" binding:"required"`     // 服务类型
	Price           float64    `json:"price" binding:"required,min=0"`      // 套餐价格
	ValidityDuration int       `json:"validity_duration"`                    // 固定时长天数
	ValidFrom       time.Time  `json:"valid_from"`                           // 有效期开始时间
	ValidTo         time.Time  `json:"valid_to"`                              // 有效期结束时间
	StoreID         *uint      `json:"store_id" binding:"required"`          // 购买门店ID
	PurchaseAmount  float64    `json:"purchase_amount" binding:"min=0"`       // 购买金额
	Status          string     `json:"status"`                              // 状态
	Description     string     `json:"description"`                          // 套餐描述/备注
}

// UpdateMemberRequest 更新会员请求
type UpdateMemberRequest struct {
	Name            string     `json:"name"`                                // 会员姓名
	Phone           string     `json:"phone"`                               // 手机号
	PackageName     string     `json:"package_name"`                        // 套餐名称
	ServiceType     string     `json:"service_type"`                        // 服务类型
	Price           float64    `json:"price"`                               // 套餐价格
	ValidityDuration int       `json:"validity_duration"`                    // 固定时长天数
	ValidFrom       time.Time  `json:"valid_from"`                           // 有效期开始时间
	ValidTo         time.Time  `json:"valid_to"`                             // 有效期结束时间
	PurchaseAmount  float64    `json:"purchase_amount"`                      // 购买金额
	Status          string     `json:"status"`                              // 状态
	Description     string     `json:"description"`                         // 套餐描述/备注
}

// UpdateMemberStatusRequest 更新会员状态请求
type UpdateMemberStatusRequest struct {
	Status string `json:"status" binding:"required"` // 状态
}

// CreateUsageRequest 创建使用记录请求
type CreateUsageRequest struct {
	ServiceItem   string     `json:"service_item" binding:"required"`  // 服务项目
	StoreID       *uint      `json:"store_id" binding:"required"`      // 使用门店ID
	TechnicianID  *uint      `json:"technician_id"`                    // 美甲师ID
	UsageDate     time.Time  `json:"usage_date"`                        // 使用日期
	Remark        string     `json:"remark"`                            // 备注
}

