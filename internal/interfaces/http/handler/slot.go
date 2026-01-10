package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/slot"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// SlotTemplateHandler 时段模板处理器
type SlotTemplateHandler struct {
	service *slot.TemplateService
	logger  logger.Logger
}

// NewSlotTemplateHandler 创建时段模板处理器
func NewSlotTemplateHandler(service *slot.TemplateService, log logger.Logger) *SlotTemplateHandler {
	return &SlotTemplateHandler{
		service: service,
		logger:  log,
	}
}

// GetTemplate 获取时段模板详情
// @Summary 获取时段模板详情
// @Description 根据模板ID获取时段模板详细信息
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} slot.Template
// @Router /slot-templates/{id} [get]
func (h *SlotTemplateHandler) GetTemplate(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的模板ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	ctx := c.Request.Context()
	template, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取时段模板失败",
			logger.NewField("request_id", requestID),
			logger.NewField("template_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if template == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "时段模板不存在")
		return
	}

	h.logger.Info("获取时段模板成功",
		logger.NewField("request_id", requestID),
		logger.NewField("template_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, template)
}

// GetTemplateList 获取时段模板列表
// @Summary 获取时段模板列表
// @Description 根据门店ID获取时段模板列表
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id query int true "门店ID"
// @Success 200 {array} slot.Template
// @Router /slot-templates [get]
func (h *SlotTemplateHandler) GetTemplateList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	storeIDStr := c.Query("store_id")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", storeIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	ctx := c.Request.Context()
	templates, err := h.service.GetByStoreID(ctx, uint(storeID))
	if err != nil {
		h.logger.Error("获取时段模板列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取时段模板列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(templates)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, templates)
}

// CreateTemplate 创建时段模板
// @Summary 创建时段模板
// @Description 创建新的时段模板
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTemplateRequest true "创建时段模板请求"
// @Success 200 {object} slot.Template
// @Router /slot-templates [post]
func (h *SlotTemplateHandler) CreateTemplate(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建时段模板请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 转换为领域实体
	newTemplate := &slot.Template{
		StoreID:      req.StoreID,
		Name:         req.Name,
		Status:       req.Status,
		WeekdayRules: req.WeekdayRules,
	}

	ctx := c.Request.Context()
	if err := h.service.Create(ctx, newTemplate); err != nil {
		h.logger.Error("创建时段模板失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建时段模板成功",
		logger.NewField("request_id", requestID),
		logger.NewField("template_id", newTemplate.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, newTemplate)
}

// UpdateTemplate 更新时段模板
// @Summary 更新时段模板
// @Description 更新时段模板信息
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param request body UpdateTemplateRequest true "更新时段模板请求"
// @Success 200 {object} slot.Template
// @Router /slot-templates/{id} [put]
func (h *SlotTemplateHandler) UpdateTemplate(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的模板ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新时段模板请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 转换为领域实体
	updateTemplate := &slot.Template{
		ID:           uint(id),
		Name:         req.Name,
		Status:       req.Status,
		WeekdayRules: req.WeekdayRules,
	}

	ctx := c.Request.Context()
	if err := h.service.Update(ctx, updateTemplate); err != nil {
		h.logger.Error("更新时段模板失败",
			logger.NewField("request_id", requestID),
			logger.NewField("template_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	// 重新获取模板信息
	updatedTemplate, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.logger.Info("更新时段模板成功",
		logger.NewField("request_id", requestID),
		logger.NewField("template_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, updatedTemplate)
}

// DeleteTemplate 删除时段模板
// @Summary 删除时段模板
// @Description 删除时段模板
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} map[string]string
// @Router /slot-templates/{id} [delete]
func (h *SlotTemplateHandler) DeleteTemplate(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的模板ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	ctx := c.Request.Context()
	if err := h.service.Delete(ctx, uint(id)); err != nil {
		h.logger.Error("删除时段模板失败",
			logger.NewField("request_id", requestID),
			logger.NewField("template_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("删除时段模板成功",
		logger.NewField("request_id", requestID),
		logger.NewField("template_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "删除时段模板成功", nil)
}

// SlotHandler 时段处理器
type SlotHandler struct {
	service *slot.SlotService
	logger  logger.Logger
}

// NewSlotHandler 创建时段处理器
func NewSlotHandler(service *slot.SlotService, log logger.Logger) *SlotHandler {
	return &SlotHandler{
		service: service,
		logger:  log,
	}
}

// GetSlot 获取时段详情
// @Summary 获取时段详情
// @Description 根据时段ID获取时段详细信息
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "时段ID"
// @Success 200 {object} slot.Slot
// @Router /slots/{id} [get]
func (h *SlotHandler) GetSlot(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的时段ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的时段ID")
		return
	}

	ctx := c.Request.Context()
	slot, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("slot_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if slot == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "时段不存在")
		return
	}

	h.logger.Info("获取时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("slot_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, slot)
}

// GetSlotList 获取时段列表
// @Summary 获取时段列表
// @Description 根据门店ID和日期获取时段列表
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id query int true "门店ID"
// @Param date query string true "日期 (YYYY-MM-DD格式)"
// @Success 200 {array} slot.Slot
// @Router /slots [get]
func (h *SlotHandler) GetSlotList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	storeIDStr := c.Query("store_id")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", storeIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.logger.Warn("无效的日期格式",
			logger.NewField("request_id", requestID),
			logger.NewField("date", dateStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的日期格式，应为YYYY-MM-DD")
		return
	}

	ctx := c.Request.Context()
	slots, err := h.service.GetByStoreIDAndDate(ctx, uint(storeID), date)
	if err != nil {
		h.logger.Error("获取时段列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取时段列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(slots)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, slots)
}

// GetAvailableSlotList 获取可用时段列表
// @Summary 获取可用时段列表
// @Description 根据门店ID和日期获取可用时段列表
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id query int true "门店ID"
// @Param date query string true "日期 (YYYY-MM-DD格式)"
// @Success 200 {array} slot.Slot
// @Router /slots/available [get]
func (h *SlotHandler) GetAvailableSlotList(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	storeIDStr := c.Query("store_id")
	storeID, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的门店ID",
			logger.NewField("request_id", requestID),
			logger.NewField("store_id", storeIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.logger.Warn("无效的日期格式",
			logger.NewField("request_id", requestID),
			logger.NewField("date", dateStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的日期格式，应为YYYY-MM-DD")
		return
	}

	ctx := c.Request.Context()
	slots, err := h.service.GetAvailableByStoreIDAndDate(ctx, uint(storeID), date)
	if err != nil {
		h.logger.Error("获取可用时段列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取可用时段列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("count", len(slots)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, slots)
}

// GenerateSlots 生成时段
// @Summary 生成时段
// @Description 根据模板生成时段
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateSlotsRequest true "生成时段请求"
// @Success 200 {object} map[string]string
// @Router /slots/generate [post]
func (h *SlotHandler) GenerateSlots(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req GenerateSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("生成时段请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的开始日期格式，应为YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的结束日期格式，应为YYYY-MM-DD")
		return
	}

	ctx := c.Request.Context()
	if err := h.service.GenerateSlots(ctx, req.StoreID, startDate, endDate, req.TechnicianCount); err != nil {
		h.logger.Error("生成时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("生成时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "生成时段成功", nil)
}

// LockSlot 锁定时段
// @Summary 锁定时段
// @Description 锁定时段（预约时调用，原子操作）
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LockSlotRequest true "锁定时段请求"
// @Success 200 {object} map[string]string
// @Router /slots/lock [post]
func (h *SlotHandler) LockSlot(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req LockSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("锁定时段请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.LockSlot(ctx, req.SlotID, req.Count); err != nil {
		h.logger.Error("锁定时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("slot_id", req.SlotID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("锁定时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("slot_id", req.SlotID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "锁定时段成功", nil)
}

// UnlockSlot 解锁时段
// @Summary 解锁时段
// @Description 解锁时段（取消预约时调用，原子操作）
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UnlockSlotRequest true "解锁时段请求"
// @Success 200 {object} map[string]string
// @Router /slots/unlock [post]
func (h *SlotHandler) UnlockSlot(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req UnlockSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("解锁时段请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.UnlockSlot(ctx, req.SlotID, req.Count); err != nil {
		h.logger.Error("解锁时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("slot_id", req.SlotID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("解锁时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("slot_id", req.SlotID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "解锁时段成功", nil)
}

// BookSlot 预约时段
// @Summary 预约时段
// @Description 预约时段（从锁定转为已预约，原子操作）
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BookSlotRequest true "预约时段请求"
// @Success 200 {object} map[string]string
// @Router /slots/book [post]
func (h *SlotHandler) BookSlot(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req BookSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("预约时段请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.BookSlot(ctx, req.SlotID, req.Count); err != nil {
		h.logger.Error("预约时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("slot_id", req.SlotID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("预约时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("slot_id", req.SlotID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "预约时段成功", nil)
}

// ReleaseSlot 释放时段
// @Summary 释放时段
// @Description 释放时段（取消或完成时调用，原子操作）
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ReleaseSlotRequest true "释放时段请求"
// @Success 200 {object} map[string]string
// @Router /slots/release [post]
func (h *SlotHandler) ReleaseSlot(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req ReleaseSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("释放时段请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	if err := h.service.ReleaseSlot(ctx, req.SlotID, req.Count); err != nil {
		h.logger.Error("释放时段失败",
			logger.NewField("request_id", requestID),
			logger.NewField("slot_id", req.SlotID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("释放时段成功",
		logger.NewField("request_id", requestID),
		logger.NewField("slot_id", req.SlotID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "释放时段成功", nil)
}

// RecalculateCapacity 重新计算时段容量
// @Summary 重新计算时段容量
// @Description 重新计算时段容量（当美甲师状态变更时调用）
// @Tags 时段管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RecalculateCapacityRequest true "重新计算容量请求"
// @Success 200 {object} map[string]string
// @Router /slots/recalculate-capacity [post]
func (h *SlotHandler) RecalculateCapacity(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req RecalculateCapacityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("重新计算容量请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的日期格式，应为YYYY-MM-DD")
		return
	}

	ctx := c.Request.Context()
	if err := h.service.RecalculateCapacity(ctx, req.StoreID, fromDate, req.TechnicianCount); err != nil {
		h.logger.Error("重新计算容量失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("重新计算容量成功",
		logger.NewField("request_id", requestID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.SuccessWithMessage(c, "重新计算容量成功", nil)
}

// ==================== 请求结构体 ====================

// CreateTemplateRequest 创建时段模板请求
type CreateTemplateRequest struct {
	StoreID      uint                `json:"store_id" binding:"required"`
	Name         string              `json:"name" binding:"required"`
	Status       string              `json:"status"`
	WeekdayRules []slot.WeekdayRule `json:"weekday_rules" binding:"required"`
}

// UpdateTemplateRequest 更新时段模板请求
type UpdateTemplateRequest struct {
	Name         string              `json:"name"`
	Status       string              `json:"status"`
	WeekdayRules []slot.WeekdayRule `json:"weekday_rules"`
}

// GenerateSlotsRequest 生成时段请求
type GenerateSlotsRequest struct {
	StoreID          uint   `json:"store_id" binding:"required"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
	TechnicianCount  int    `json:"technician_count" binding:"required,min=1"`
}

// LockSlotRequest 锁定时段请求
type LockSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	Count  int  `json:"count" binding:"required,min=1"`
}

// UnlockSlotRequest 解锁时段请求
type UnlockSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	Count  int  `json:"count" binding:"required,min=1"`
}

// BookSlotRequest 预约时段请求
type BookSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	Count  int  `json:"count" binding:"required,min=1"`
}

// ReleaseSlotRequest 释放时段请求
type ReleaseSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	Count  int  `json:"count" binding:"required,min=1"`
}

// RecalculateCapacityRequest 重新计算容量请求
type RecalculateCapacityRequest struct {
	StoreID         uint   `json:"store_id" binding:"required"`
	FromDate        string `json:"from_date" binding:"required"`
	TechnicianCount int    `json:"technician_count" binding:"required,min=0"`
}

