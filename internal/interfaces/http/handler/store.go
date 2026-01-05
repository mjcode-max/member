package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// StoreHandler 门店处理器
type StoreHandler struct {
	service *store.StoreService
	logger  logger.Logger
}

// NewStoreHandler 创建门店处理器
func NewStoreHandler(service *store.StoreService, log logger.Logger) *StoreHandler {
	return &StoreHandler{
		service: service,
		logger:  log,
	}
}

// GetStore 获取门店详情
// @Summary 获取门店详情
// @Description 根据门店ID获取门店详细信息
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Success 200 {object} store.Store
// @Router /stores/{id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	// 获取当前用户角色
	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	ctx := c.Request.Context()
	s, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取门店失败",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if s == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "门店不存在")
		return
	}

	// 权限检查：店长只能查看自己的门店
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

		if s.ID != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权查看该门店")
			return
		}
	}

	h.logger.Info("获取门店成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, s)
}

// GetStoreList 获取门店列表
// @Summary 获取门店列表
// @Description 获取门店列表，支持按状态、名称筛选和分页
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "状态筛选" Enums(operating, closed, shutdown)
// @Param name query string false "门店名称（模糊搜索）"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} utils.PaginationResponse
// @Router /stores [get]
func (h *StoreHandler) GetStoreList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 获取查询参数
	status := c.Query("status")
	name := c.Query("name")
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
	stores, total, err := h.service.GetList(ctx, status, name, page, pageSize)
	if err != nil {
		h.logger.Error("获取门店列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取门店列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(stores)),
		logger.NewField("total", total),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithPagination(c, stores, page, pageSize, total)
}

// CreateStore 创建门店
// @Summary 创建门店
// @Description 创建新门店（仅后台）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateStoreRequest true "创建门店请求"
// @Success 200 {object} store.Store
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建门店请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 转换为领域实体
	newStore := &store.Store{
		Name:              req.Name,
		Address:           req.Address,
		Phone:             req.Phone,
		ContactPerson:     req.ContactPerson,
		Status:            req.Status,
		BusinessHoursStart: req.BusinessHoursStart,
		BusinessHoursEnd:   req.BusinessHoursEnd,
		DepositAmount:     req.DepositAmount,
	}

	ctx := c.Request.Context()
	if err := h.service.Create(ctx, newStore); err != nil {
		h.logger.Error("创建门店失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建门店成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", newStore.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, newStore)
}

// UpdateStore 更新门店
// @Summary 更新门店
// @Description 更新门店信息
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Param request body UpdateStoreRequest true "更新门店请求"
// @Success 200 {object} store.Store
// @Router /stores/{id} [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 获取当前用户角色
	userRoleInterface, _ := c.Get("user_role")
	userRole, _ := userRoleInterface.(string)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	// 权限检查：店长只能更新自己的门店
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

		if uint(id) != storeID {
			utils.ErrorWithCode(c, http.StatusForbidden, "无权更新该门店")
			return
		}
	}

	var req UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新门店请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 转换为领域实体
	updateStore := &store.Store{
		ID:                uint(id),
		Name:              req.Name,
		Address:           req.Address,
		Phone:             req.Phone,
		ContactPerson:     req.ContactPerson,
		Status:            req.Status,
		BusinessHoursStart: req.BusinessHoursStart,
		BusinessHoursEnd:   req.BusinessHoursEnd,
		DepositAmount:     req.DepositAmount,
	}

	ctx := c.Request.Context()
	if err := h.service.Update(ctx, updateStore); err != nil {
		h.logger.Error("更新门店失败",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	// 重新获取门店信息
	updatedStore, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新门店成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, updatedStore)
}

// DeleteStore 删除门店
// @Summary 删除门店
// @Description 删除门店（软删除，仅后台）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Success 200 {object} map[string]string
// @Router /stores/{id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	ctx := c.Request.Context()
	if err := h.service.Delete(ctx, uint(id)); err != nil {
		h.logger.Error("删除门店失败",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("删除门店成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "删除门店成功", nil)
}

// UpdateStoreStatus 更新门店状态
// @Summary 更新门店状态
// @Description 更新门店状态（营业中/停业/关闭）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Param request body UpdateStoreStatusRequest true "更新门店状态请求"
// @Success 200 {object} map[string]string
// @Router /stores/{id}/status [put]
func (h *StoreHandler) UpdateStoreStatus(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	var req UpdateStoreStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新门店状态请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.UpdateStatus(ctx, uint(id), req.Status); err != nil {
		h.logger.Error("更新门店状态失败",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", id),
			logger.NewField("status", req.Status),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新门店状态成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", id),
		logger.NewField("status", req.Status),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "更新门店状态成功", nil)
}

// CreateStoreRequest 创建门店请求
type CreateStoreRequest struct {
	Name              string  `json:"name" binding:"required"`                    // 门店名称
	Address           string  `json:"address"`                                   // 门店地址
	Phone             string  `json:"phone"`                                     // 联系电话
	ContactPerson     string  `json:"contact_person"`                            // 联系人
	Status            string  `json:"status"`                                     // 状态: operating, closed, shutdown
	BusinessHoursStart string `json:"business_hours_start"`                       // 营业开始时间 (HH:MM格式)
	BusinessHoursEnd   string `json:"business_hours_end"`                         // 营业结束时间 (HH:MM格式)
	DepositAmount     float64 `json:"deposit_amount"`                            // 押金金额
}

// UpdateStoreRequest 更新门店请求
type UpdateStoreRequest struct {
	Name              string  `json:"name"`                // 门店名称
	Address           string  `json:"address"`             // 门店地址
	Phone             string  `json:"phone"`               // 联系电话
	ContactPerson     string  `json:"contact_person"`      // 联系人
	Status            string  `json:"status"`              // 状态
	BusinessHoursStart string `json:"business_hours_start"` // 营业开始时间
	BusinessHoursEnd   string `json:"business_hours_end"`   // 营业结束时间
	DepositAmount     float64 `json:"deposit_amount"`      // 押金金额
}

// UpdateStoreStatusRequest 更新门店状态请求
type UpdateStoreStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=operating closed shutdown"` // 状态
}

