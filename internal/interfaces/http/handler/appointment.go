package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/appointment"
	"member-pre/internal/domain/slot"
	"member-pre/pkg/logger"
	"member-pre/pkg/utils"
)

// AppointmentHandler 预约处理器
type AppointmentHandler struct {
	service     *appointment.AppointmentService
	slotService *slot.SlotService
	logger      logger.Logger
}

// NewAppointmentHandler 创建预约处理器
func NewAppointmentHandler(service *appointment.AppointmentService, slotService *slot.SlotService, log logger.Logger) *AppointmentHandler {
	return &AppointmentHandler{
		service:     service,
		slotService: slotService,
		logger:      log,
	}
}

// GetAppointment 获取预约详情
// @Summary 获取预约详情
// @Description 根据预约ID获取预约详细信息
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "预约ID"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/{id} [get]
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的预约ID",
			logger.NewField("request_id", requestID),
			logger.NewField("id", idStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的预约ID")
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("获取预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", id),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	if appointment == nil {
		utils.ErrorWithCode(c, http.StatusNotFound, "预约不存在")
		return
	}

	h.logger.Info("获取预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", id),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// GetMyAppointments 获取我的预约列表
// @Summary 获取我的预约列表
// @Description 获取当前顾客的预约列表
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} appointment.Appointment
// @Router /appointments/my [get]
func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 从上下文中获取当前用户ID
	customerID, err := getCurrentUserID(c)
	if err != nil {
		h.logger.Warn("获取当前用户ID失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未登录")
		return
	}

	ctx := c.Request.Context()
	appointments, err := h.service.GetByCustomerID(ctx, customerID)
	if err != nil {
		h.logger.Error("获取我的预约列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("customer_id", customerID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取我的预约列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("customer_id", customerID),
		logger.NewField("count", len(appointments)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointments)
}

// GetMyUpcomingAppointments 获取我即将到来的预约
// @Summary 获取我即将到来的预约
// @Description 获取当前顾客即将到来的预约（未取消、未完成）
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} appointment.Appointment
// @Router /appointments/my/upcoming [get]
func (h *AppointmentHandler) GetMyUpcomingAppointments(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	// 从上下文中获取当前用户ID
	customerID, err := getCurrentUserID(c)
	if err != nil {
		h.logger.Warn("获取当前用户ID失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未登录")
		return
	}

	ctx := c.Request.Context()
	appointments, err := h.service.GetUpcomingByCustomerID(ctx, customerID)
	if err != nil {
		h.logger.Error("获取我即将到来的预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("customer_id", customerID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取我即将到来的预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("customer_id", customerID),
		logger.NewField("count", len(appointments)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointments)
}

// CreateAppointment 创建预约
// @Summary 创建预约
// @Description 创建新的预约（创建时锁定时段，需要支付押金）
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAppointmentRequest true "创建预约请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建预约请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 从上下文中获取当前用户ID
	customerID, err := getCurrentUserID(c)
	if err != nil {
		h.logger.Warn("获取当前用户ID失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusUnauthorized, "未登录")
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.CreateAppointment(ctx, &appointment.CreateAppointmentRequest{
		CustomerID:   customerID,
		StoreID:      req.StoreID,
		SlotID:       req.SlotID,
		TechnicianID: req.TechnicianID,
		ServiceName:  req.ServiceName,
		ServicePrice: req.ServicePrice,
		Remark:       req.Remark,
	})
	if err != nil {
		h.logger.Error("创建预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("创建预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", appointment.ID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// PayDeposit 支付押金
// @Summary 支付押金
// @Description 支付预约押金（10元）
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PayDepositRequest true "支付押金请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/pay-deposit [post]
func (h *AppointmentHandler) PayDeposit(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req PayDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("支付押金请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.PayDeposit(ctx, req.AppointmentID, req.PaymentMethod)
	if err != nil {
		h.logger.Error("支付押金失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", req.AppointmentID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("支付押金成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", req.AppointmentID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// ConfirmArrival 美甲师确认到店
// @Summary 美甲师确认到店
// @Description 美甲师确认顾客到店，退还押金
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConfirmArrivalRequest true "确认到店请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/confirm-arrival [post]
func (h *AppointmentHandler) ConfirmArrival(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req ConfirmArrivalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("确认到店请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.ConfirmArrival(ctx, req.AppointmentID)
	if err != nil {
		h.logger.Error("确认到店失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", req.AppointmentID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("确认到店成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", req.AppointmentID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// Complete 完成预约
// @Summary 完成预约
// @Description 美甲师完成预约服务
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CompleteAppointmentRequest true "完成预约请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/complete [post]
func (h *AppointmentHandler) Complete(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CompleteAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("完成预约请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.Complete(ctx, req.AppointmentID)
	if err != nil {
		h.logger.Error("完成预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", req.AppointmentID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("完成预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", req.AppointmentID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// CancelByCustomer 顾客取消预约
// @Summary 顾客取消预约
// @Description 顾客取消预约（需提前3小时）
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CancelByCustomerRequest true "取消预约请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/cancel/customer [post]
func (h *AppointmentHandler) CancelByCustomer(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CancelByCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("顾客取消预约请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.CancelByCustomer(ctx, req.AppointmentID, req.Reason)
	if err != nil {
		h.logger.Error("顾客取消预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", req.AppointmentID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("顾客取消预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", req.AppointmentID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// CancelByTechnician 美甲师取消预约
// @Summary 美甲师取消预约
// @Description 美甲师随时可以取消预约
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CancelByTechnicianRequest true "取消预约请求"
// @Success 200 {object} appointment.Appointment
// @Router /appointments/cancel/technician [post]
func (h *AppointmentHandler) CancelByTechnician(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	var req CancelByTechnicianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("美甲师取消预约请求参数错误",
			logger.NewField("request_id", requestID),
			logger.NewField("error", err.Error()),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.service.CancelByTechnician(ctx, req.AppointmentID, req.Reason)
	if err != nil {
		h.logger.Error("美甲师取消预约失败",
			logger.NewField("request_id", requestID),
			logger.NewField("appointment_id", req.AppointmentID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("美甲师取消预约成功",
		logger.NewField("request_id", requestID),
		logger.NewField("appointment_id", req.AppointmentID),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointment)
}

// GetStoreAppointments 获取门店预约列表
// @Summary 获取门店预约列表
// @Description 根据门店ID获取预约列表（后台、店长、美甲师使用）
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param store_id query int true "门店ID"
// @Param start_date query string false "开始日期 (YYYY-MM-DD格式)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD格式)"
// @Success 200 {array} appointment.Appointment
// @Router /appointments/store [get]
func (h *AppointmentHandler) GetStoreAppointments(c *gin.Context) {
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
	var appointments []*appointment.Appointment

	// 检查是否有日期范围参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的开始日期格式，应为YYYY-MM-DD")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.ErrorWithCode(c, http.StatusBadRequest, "无效的结束日期格式，应为YYYY-MM-DD")
			return
		}

		appointments, err = h.service.GetByDateRange(ctx, startDate, endDate)
		if err != nil {
			h.logger.Error("获取门店预约列表失败",
				logger.NewField("request_id", requestID),
				logger.NewField("store_id", storeID),
				logger.NewField("error", err.Error()),
				logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
			)
			utils.Error(c, err)
			return
		}
	} else {
		appointments, err = h.service.GetByStoreID(ctx, uint(storeID))
		if err != nil {
			h.logger.Error("获取门店预约列表失败",
				logger.NewField("request_id", requestID),
				logger.NewField("store_id", storeID),
				logger.NewField("error", err.Error()),
				logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
			)
			utils.Error(c, err)
			return
		}
	}

	h.logger.Info("获取门店预约列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("store_id", storeID),
		logger.NewField("count", len(appointments)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointments)
}

// GetTechnicianAppointments 获取美甲师预约列表
// @Summary 获取美甲师预约列表
// @Description 根据美甲师ID获取预约列表
// @Tags 预约管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param technician_id query int true "美甲师ID"
// @Success 200 {array} appointment.Appointment
// @Router /appointments/technician [get]
func (h *AppointmentHandler) GetTechnicianAppointments(c *gin.Context) {
	startTime := time.Now()
	requestID := getRequestID(c)

	technicianIDStr := c.Query("technician_id")
	technicianID, err := strconv.ParseUint(technicianIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("无效的美甲师ID",
			logger.NewField("request_id", requestID),
			logger.NewField("technician_id", technicianIDStr),
		)
		utils.ErrorWithCode(c, http.StatusBadRequest, "无效的美甲师ID")
		return
	}

	ctx := c.Request.Context()
	appointments, err := h.service.GetByTechnicianID(ctx, uint(technicianID))
	if err != nil {
		h.logger.Error("获取美甲师预约列表失败",
			logger.NewField("request_id", requestID),
			logger.NewField("technician_id", technicianID),
			logger.NewField("error", err.Error()),
			logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
		)
		utils.Error(c, err)
		return
	}

	h.logger.Info("获取美甲师预约列表成功",
		logger.NewField("request_id", requestID),
		logger.NewField("technician_id", technicianID),
		logger.NewField("count", len(appointments)),
		logger.NewField("duration_ms", time.Since(startTime).Milliseconds()),
	)

	utils.Success(c, appointments)
}

// ==================== 请求结构体 ====================

// CreateAppointmentRequest 创建预约请求
type CreateAppointmentRequest struct {
	StoreID      uint    `json:"store_id" binding:"required"`
	SlotID       uint    `json:"slot_id" binding:"required"`
	TechnicianID *uint   `json:"technician_id"`
	ServiceName  string  `json:"service_name" binding:"required"`
	ServicePrice float64 `json:"service_price" binding:"required,gt=0"`
	Remark       string  `json:"remark"`
}

// PayDepositRequest 支付押金请求
type PayDepositRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"` // wechat, alipay, etc.
}

// ConfirmArrivalRequest 确认到店请求
type ConfirmArrivalRequest struct {
	AppointmentID uint `json:"appointment_id" binding:"required"`
}

// CompleteAppointmentRequest 完成预约请求
type CompleteAppointmentRequest struct {
	AppointmentID uint `json:"appointment_id" binding:"required"`
}

// CancelByCustomerRequest 顾客取消预约请求
type CancelByCustomerRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	Reason        string `json:"reason"`
}

// CancelByTechnicianRequest 美甲师取消预约请求
type CancelByTechnicianRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	Reason        string `json:"reason"`
}

// getCurrentUserID 从上下文中获取当前用户ID
func getCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, nil
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, nil
	}

	return id, nil
}

