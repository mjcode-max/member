package appointment

import (
	"context"
	"time"

	"member-pre/internal/domain/slot"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 预约状态常量
const (
	AppointmentStatusPending   = "pending"   // 待支付（已创建，未支付押金）
	AppointmentStatusPaid      = "paid"      // 已支付（已支付押金，等待确认）
	AppointmentStatusConfirmed = "confirmed" // 已确认（美甲师确认到店）
	AppointmentStatusCompleted = "completed" // 已完成
	AppointmentStatusCancelled = "cancelled" // 已取消
)

// 取消来源常量
const (
	CancelSourceCustomer    = "customer"    // 顾客取消
	CancelSourceTechnician  = "technician"  // 美甲师取消
	CancelSourceSystem      = "system"      // 系统取消
)

// 押金金额常量
const (
	DepositAmount = 10.00 // 预约押金金额（元）
)

// Appointment 预约实体
type Appointment struct {
	ID               uint      `json:"id"`
	CustomerID       uint      `json:"customer_id"`       // 顾客ID
	StoreID          uint      `json:"store_id"`          // 门店ID
	SlotID           uint      `json:"slot_id"`           // 时段ID
	TechnicianID     *uint     `json:"technician_id"`     // 美甲师ID（可选）
	ServiceName      string    `json:"service_name"`      // 服务名称
	ServicePrice     float64   `json:"service_price"`     // 服务价格（元）
	DepositAmount    float64   `json:"deposit_amount"`    // 押金金额（元）
	DepositPaid      bool      `json:"deposit_paid"`      // 是否已支付押金
	DepositPaidAt    *time.Time `json:"deposit_paid_at"`   // 押金支付时间
	DepositRefunded  bool      `json:"deposit_refunded"`  // 是否已退还押金
	DepositRefundedAt *time.Time `json:"deposit_refunded_at"` // 押金退还时间
	Status           string    `json:"status"`            // 状态: pending, paid, confirmed, completed, cancelled
	CancelledAt      *time.Time `json:"cancelled_at"`      // 取消时间
	CancelledBy      string    `json:"cancelled_by"`      // 取消人: customer, technician, system
	CancelReason     string    `json:"cancel_reason"`     // 取消原因
	Remark           string    `json:"remark"`            // 备注
	AppointmentTime  time.Time `json:"appointment_time"`  // 预约时间（时段开始时间）
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// IAppointmentRepository 预约仓储接口
type IAppointmentRepository interface {
	// FindByID 根据ID查找预约
	FindByID(ctx context.Context, id uint) (*Appointment, error)
	// FindByCustomerID 根据顾客ID查找预约列表
	FindByCustomerID(ctx context.Context, customerID uint) ([]*Appointment, error)
	// FindByStoreID 根据门店ID查找预约列表
	FindByStoreID(ctx context.Context, storeID uint) ([]*Appointment, error)
	// FindByTechnicianID 根据美甲师ID查找预约列表
	FindByTechnicianID(ctx context.Context, technicianID uint) ([]*Appointment, error)
	// FindBySlotID 根据时段ID查找预约列表
	FindBySlotID(ctx context.Context, slotID uint) ([]*Appointment, error)
	// FindByStatus 根据状态查找预约列表
	FindByStatus(ctx context.Context, status string) ([]*Appointment, error)
	// FindByDateRange 根据日期范围查找预约列表
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Appointment, error)
	// FindUpcomingByCustomerID 查找顾客即将到来的预约（未取消、未完成）
	FindUpcomingByCustomerID(ctx context.Context, customerID uint) ([]*Appointment, error)
	// Create 创建预约
	Create(ctx context.Context, appointment *Appointment) error
	// Update 更新预约
	Update(ctx context.Context, appointment *Appointment) error
	// UpdateStatus 更新预约状态
	UpdateStatus(ctx context.Context, id uint, status string) error
	// Delete 删除预约
	Delete(ctx context.Context, id uint) error
}

// 领域错误定义
var (
	ErrAppointmentNotFound      = errors.ErrNotFound("预约不存在")
	ErrAppointmentNotPending    = errors.ErrInvalidParams("预约不是待支付状态")
	ErrAppointmentNotPaid       = errors.ErrInvalidParams("预约未支付押金")
	ErrAppointmentAlreadyPaid   = errors.ErrInvalidParams("预约已支付押金")
	ErrDepositAmountMismatch    = errors.ErrInvalidParams("押金金额不匹配")
	ErrCancelTooLate            = errors.ErrInvalidParams("距离预约时间不足3小时，无法取消")
	ErrCancelNotAllowed         = errors.ErrInvalidParams("当前状态不允许取消")
	ErrDepositAlreadyRefunded   = errors.ErrInvalidParams("押金已退还")
	ErrTechnicianNotAssigned    = errors.ErrInvalidParams("未指定美甲师")
)

// AppointmentService 预约服务
type AppointmentService struct {
	appointmentRepo IAppointmentRepository
	paymentService  IPaymentService
	slotService     *slot.SlotService
	logger          logger.Logger
}

// NewAppointmentService 创建预约服务
func NewAppointmentService(
	appointmentRepo IAppointmentRepository,
	paymentService IPaymentService,
	slotService *slot.SlotService,
	log logger.Logger,
) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		paymentService:  paymentService,
		slotService:     slotService,
		logger:          log,
	}
}

// CreateAppointment 创建预约（创建时锁定时段，需要支付押金）
func (s *AppointmentService) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*Appointment, error) {
	s.logger.Info("创建预约",
		logger.NewField("customer_id", req.CustomerID),
		logger.NewField("slot_id", req.SlotID),
		logger.NewField("store_id", req.StoreID),
	)

	// 1. 检查时段可用性
	slotInfo, err := s.slotService.GetByID(ctx, req.SlotID)
	if err != nil {
		return nil, err
	}
	if slotInfo == nil {
		return nil, slot.ErrSlotNotFound
	}

	// 2. 检查时段容量
	available := slotInfo.Capacity - slotInfo.LockedCount - slotInfo.BookedCount
	if available < 1 {
		return nil, slot.ErrInsufficientCapacity
	}

	// 3. 锁定时段
	if err := s.slotService.LockSlot(ctx, req.SlotID, 1); err != nil {
		return nil, err
	}

	// 4. 创建预约记录（初始状态为待支付）
	appointment := &Appointment{
		CustomerID:      req.CustomerID,
		StoreID:         req.StoreID,
		SlotID:          req.SlotID,
		TechnicianID:    req.TechnicianID,
		ServiceName:     req.ServiceName,
		ServicePrice:    req.ServicePrice,
		DepositAmount:   DepositAmount,
		DepositPaid:     false,
		Status:          AppointmentStatusPending,
		AppointmentTime: slotInfo.StartTime,
		Remark:          req.Remark,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.appointmentRepo.Create(ctx, appointment); err != nil {
		// 创建失败时释放时段锁定
		_ = s.slotService.UnlockSlot(ctx, req.SlotID, 1)
		s.logger.Error("创建预约失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	s.logger.Info("创建预约成功", logger.NewField("appointment_id", appointment.ID))
	return appointment, nil
}

// PayDeposit 支付押金
func (s *AppointmentService) PayDeposit(ctx context.Context, appointmentID uint, paymentMethod string) (*Appointment, error) {
	s.logger.Info("支付押金",
		logger.NewField("appointment_id", appointmentID),
		logger.NewField("payment_method", paymentMethod),
	)

	appointment, err := s.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// 检查状态
	if appointment.Status != AppointmentStatusPending {
		return nil, ErrAppointmentNotPending
	}

	// 检查是否已支付
	if appointment.DepositPaid {
		return nil, ErrAppointmentAlreadyPaid
	}

	// 调用支付服务（先调用，成功后再更新状态）
	err = s.paymentService.CreatePayment(ctx, &PaymentRequest{
		AppointmentID: appointmentID,
		Amount:        appointment.DepositAmount,
		PaymentMethod: paymentMethod,
		PaymentType:   PaymentTypeDeposit,
	})
	if err != nil {
		s.logger.Error("创建支付订单失败", logger.NewField("error", err.Error()))
		return nil, err
	}

	// 更新预约状态为已支付
	now := time.Now()
	appointment.DepositPaid = true
	appointment.DepositPaidAt = &now
	appointment.Status = AppointmentStatusPaid
	appointment.UpdatedAt = now

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		s.logger.Error("更新预约状态失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	s.logger.Info("支付押金成功", logger.NewField("appointment_id", appointmentID))
	return appointment, nil
}

// ConfirmArrival 美甲师确认到店（退还押金）
func (s *AppointmentService) ConfirmArrival(ctx context.Context, appointmentID uint) (*Appointment, error) {
	s.logger.Info("美甲师确认到店", logger.NewField("appointment_id", appointmentID))

	appointment, err := s.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// 检查状态（必须是已支付状态）
	if appointment.Status != AppointmentStatusPaid {
		return nil, errors.ErrInvalidParams("预约未支付押金，无法确认到店")
	}

	// 检查押金是否已退还
	if appointment.DepositRefunded {
		return nil, ErrDepositAlreadyRefunded
	}

	// 调用退款服务
	err = s.paymentService.CreateRefund(ctx, &RefundRequest{
		AppointmentID: appointmentID,
		Amount:        appointment.DepositAmount,
		RefundType:    RefundTypeArrival,
	})
	if err != nil {
		s.logger.Error("创建退款订单失败", logger.NewField("error", err.Error()))
		return nil, err
	}

	// 更新预约状态为已确认
	now := time.Now()
	appointment.Status = AppointmentStatusConfirmed
	appointment.DepositRefunded = true
	appointment.DepositRefundedAt = &now
	appointment.UpdatedAt = now

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		s.logger.Error("更新预约状态失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	s.logger.Info("确认到店成功，押金已退还", logger.NewField("appointment_id", appointmentID))
	return appointment, nil
}

// Complete 完成预约
func (s *AppointmentService) Complete(ctx context.Context, appointmentID uint) (*Appointment, error) {
	s.logger.Info("完成预约", logger.NewField("appointment_id", appointmentID))

	appointment, err := s.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// 检查状态
	if appointment.Status != AppointmentStatusConfirmed {
		return nil, errors.ErrInvalidParams("预约未确认到店，无法完成")
	}

	// 更新状态
	now := time.Now()
	appointment.Status = AppointmentStatusCompleted
	appointment.UpdatedAt = now

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		s.logger.Error("更新预约状态失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	// 释放时段（从已预约转为可用）
	if err := s.slotService.ReleaseSlot(ctx, appointment.SlotID, 1); err != nil {
		s.logger.Warn("释放时段失败", logger.NewField("slot_id", appointment.SlotID), logger.NewField("error", err.Error()))
		// 不返回错误，因为预约已完成
	}

	s.logger.Info("预约完成成功", logger.NewField("appointment_id", appointmentID))
	return appointment, nil
}

// CancelByCustomer 顾客取消预约（需提前3小时）
func (s *AppointmentService) CancelByCustomer(ctx context.Context, appointmentID uint, reason string) (*Appointment, error) {
	s.logger.Info("顾客取消预约", logger.NewField("appointment_id", appointmentID))

	appointment, err := s.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// 检查是否允许取消
	if appointment.Status == AppointmentStatusCancelled ||
		appointment.Status == AppointmentStatusCompleted {
		return nil, ErrCancelNotAllowed
	}

	// 检查时间（需提前3小时）
	now := time.Now()
	appointmentTime := appointment.AppointmentTime
	if appointmentTime.Sub(now) < 3*time.Hour {
		return nil, ErrCancelTooLate
	}

	// 更新预约状态
	now = time.Now()
	appointment.Status = AppointmentStatusCancelled
	appointment.CancelledAt = &now
	appointment.CancelledBy = CancelSourceCustomer
	appointment.CancelReason = reason
	appointment.UpdatedAt = now

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		s.logger.Error("更新预约状态失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	// 释放时段锁定
	if appointment.Status == AppointmentStatusPending {
		// 如果是待支付状态，直接释放锁定
		if err := s.slotService.UnlockSlot(ctx, appointment.SlotID, 1); err != nil {
			s.logger.Warn("释放时段锁定失败", logger.NewField("slot_id", appointment.SlotID), logger.NewField("error", err.Error()))
		}
	} else if appointment.Status == AppointmentStatusPaid {
		// 如果是已支付状态，释放已预约数量并退款
		if err := s.slotService.ReleaseSlot(ctx, appointment.SlotID, 1); err != nil {
			s.logger.Warn("释放时段失败", logger.NewField("slot_id", appointment.SlotID), logger.NewField("error", err.Error()))
		}

		// 如果已支付押金，需要退还
		if appointment.DepositPaid && !appointment.DepositRefunded {
			// 调用退款服务
			_ = s.paymentService.CreateRefund(ctx, &RefundRequest{
				AppointmentID: appointmentID,
				Amount:        appointment.DepositAmount,
				RefundType:    RefundTypeCancel,
			})
			// 更新退款状态
			appointment.DepositRefunded = true
			appointment.DepositRefundedAt = &now
			_ = s.appointmentRepo.Update(ctx, appointment)
		}
	}

	s.logger.Info("顾客取消预约成功", logger.NewField("appointment_id", appointmentID))
	return appointment, nil
}

// CancelByTechnician 美甲师取消预约（随时可以取消）
func (s *AppointmentService) CancelByTechnician(ctx context.Context, appointmentID uint, reason string) (*Appointment, error) {
	s.logger.Info("美甲师取消预约", logger.NewField("appointment_id", appointmentID))

	appointment, err := s.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// 检查是否允许取消
	if appointment.Status == AppointmentStatusCancelled ||
		appointment.Status == AppointmentStatusCompleted {
		return nil, ErrCancelNotAllowed
	}

	// 更新预约状态
	now := time.Now()
	appointment.Status = AppointmentStatusCancelled
	appointment.CancelledAt = &now
	appointment.CancelledBy = CancelSourceTechnician
	appointment.CancelReason = reason
	appointment.UpdatedAt = now

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		s.logger.Error("更新预约状态失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	// 释放时段
	if appointment.Status == AppointmentStatusPending {
		// 如果是待支付状态，直接释放锁定
		if err := s.slotService.UnlockSlot(ctx, appointment.SlotID, 1); err != nil {
			s.logger.Warn("释放时段锁定失败", logger.NewField("slot_id", appointment.SlotID), logger.NewField("error", err.Error()))
		}
	} else if appointment.Status == AppointmentStatusPaid {
		// 如果是已支付状态，释放已预约数量
		if err := s.slotService.ReleaseSlot(ctx, appointment.SlotID, 1); err != nil {
			s.logger.Warn("释放时段失败", logger.NewField("slot_id", appointment.SlotID), logger.NewField("error", err.Error()))
		}

		// 如果已支付押金，需要退还
		if appointment.DepositPaid && !appointment.DepositRefunded {
			// 调用退款服务
			_ = s.paymentService.CreateRefund(ctx, &RefundRequest{
				AppointmentID: appointmentID,
				Amount:        appointment.DepositAmount,
				RefundType:    RefundTypeCancel,
			})
			// 更新退款状态
			appointment.DepositRefunded = true
			appointment.DepositRefundedAt = &now
			_ = s.appointmentRepo.Update(ctx, appointment)
		}
	}

	s.logger.Info("美甲师取消预约成功", logger.NewField("appointment_id", appointmentID))
	return appointment, nil
}

// GetByID 获取预约详情
func (s *AppointmentService) GetByID(ctx context.Context, id uint) (*Appointment, error) {
	s.logger.Debug("获取预约详情", logger.NewField("appointment_id", id))

	appointment, err := s.appointmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	return appointment, nil
}

// GetByCustomerID 获取顾客的预约列表
func (s *AppointmentService) GetByCustomerID(ctx context.Context, customerID uint) ([]*Appointment, error) {
	s.logger.Debug("获取顾客预约列表", logger.NewField("customer_id", customerID))

	return s.appointmentRepo.FindByCustomerID(ctx, customerID)
}

// GetUpcomingByCustomerID 获取顾客即将到来的预约
func (s *AppointmentService) GetUpcomingByCustomerID(ctx context.Context, customerID uint) ([]*Appointment, error) {
	s.logger.Debug("获取顾客即将到来的预约", logger.NewField("customer_id", customerID))

	return s.appointmentRepo.FindUpcomingByCustomerID(ctx, customerID)
}

// GetByStoreID 获取门店的预约列表
func (s *AppointmentService) GetByStoreID(ctx context.Context, storeID uint) ([]*Appointment, error) {
	s.logger.Debug("获取门店预约列表", logger.NewField("store_id", storeID))

	return s.appointmentRepo.FindByStoreID(ctx, storeID)
}

// GetByTechnicianID 获取美甲师的预约列表
func (s *AppointmentService) GetByTechnicianID(ctx context.Context, technicianID uint) ([]*Appointment, error) {
	s.logger.Debug("获取美甲师预约列表", logger.NewField("technician_id", technicianID))

	return s.appointmentRepo.FindByTechnicianID(ctx, technicianID)
}

// GetByDateRange 根据日期范围获取预约列表
func (s *AppointmentService) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Appointment, error) {
	s.logger.Debug("获取日期范围内的预约",
		logger.NewField("start_date", startDate.Format("2006-01-02")),
		logger.NewField("end_date", endDate.Format("2006-01-02")),
	)

	return s.appointmentRepo.FindByDateRange(ctx, startDate, endDate)
}

// CreateAppointmentRequest 创建预约请求
type CreateAppointmentRequest struct {
	CustomerID    uint     `json:"customer_id" binding:"required"`
	StoreID       uint     `json:"store_id" binding:"required"`
	SlotID        uint     `json:"slot_id" binding:"required"`
	TechnicianID  *uint    `json:"technician_id"`
	ServiceName   string   `json:"service_name" binding:"required"`
	ServicePrice  float64  `json:"service_price" binding:"required,gt=0"`
	Remark        string   `json:"remark"`
}

