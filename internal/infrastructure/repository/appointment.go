package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/appointment"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

var _ appointment.IAppointmentRepository = (*AppointmentRepository)(nil)
var _ appointment.IPaymentRepository = (*PaymentRepository)(nil)
var _ appointment.IRefundRepository = (*RefundRepository)(nil)

// 在包初始化时注册模型，确保迁移时能检测到
func init() {
	persistence.Register(&AppointmentModel{})
	persistence.Register(&PaymentModel{})
	persistence.Register(&RefundModel{})
}

// ==================== 预约仓储 ====================

// AppointmentRepository 预约仓储实现
type AppointmentRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewAppointmentRepository 创建预约仓储实例
func NewAppointmentRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *AppointmentRepository {
	return &AppointmentRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// AppointmentModel 预约数据库模型
type AppointmentModel struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CustomerID        uint           `gorm:"index;not null" json:"customer_id"`
	StoreID           uint           `gorm:"index;not null" json:"store_id"`
	SlotID            uint           `gorm:"index;not null" json:"slot_id"`
	TechnicianID      *uint          `gorm:"index" json:"technician_id"`
	ServiceName       string         `gorm:"size:100;not null" json:"service_name"`
	ServicePrice      float64        `gorm:"not null" json:"service_price"`
	DepositAmount     float64        `gorm:"not null" json:"deposit_amount"`
	DepositPaid       bool           `gorm:"default:false;not null" json:"deposit_paid"`
	DepositPaidAt     *time.Time     `json:"deposit_paid_at"`
	DepositRefunded   bool           `gorm:"default:false;not null" json:"deposit_refunded"`
	DepositRefundedAt *time.Time     `json:"deposit_refunded_at"`
	Status            string         `gorm:"size:20;default:'pending';not null;index" json:"status"`
	CancelledAt       *time.Time     `json:"cancelled_at"`
	CancelledBy       string         `gorm:"size:20" json:"cancelled_by"`
	CancelReason      string         `gorm:"size:500" json:"cancel_reason"`
	Remark            string         `gorm:"size:500" json:"remark"`
	AppointmentTime   time.Time      `gorm:"not null;index" json:"appointment_time"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AppointmentModel) TableName() string {
	return "appointments"
}

// ToEntity 转换为领域实体
func (m *AppointmentModel) ToEntity() *appointment.Appointment {
	if m == nil {
		return nil
	}
	return &appointment.Appointment{
		ID:                 m.ID,
		CustomerID:         m.CustomerID,
		StoreID:            m.StoreID,
		SlotID:             m.SlotID,
		TechnicianID:       m.TechnicianID,
		ServiceName:        m.ServiceName,
		ServicePrice:       m.ServicePrice,
		DepositAmount:      m.DepositAmount,
		DepositPaid:        m.DepositPaid,
		DepositPaidAt:      m.DepositPaidAt,
		DepositRefunded:    m.DepositRefunded,
		DepositRefundedAt:  m.DepositRefundedAt,
		Status:             m.Status,
		CancelledAt:        m.CancelledAt,
		CancelledBy:        m.CancelledBy,
		CancelReason:       m.CancelReason,
		Remark:             m.Remark,
		AppointmentTime:    m.AppointmentTime,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *AppointmentModel) FromEntity(a *appointment.Appointment) {
	if a == nil {
		return
	}
	m.ID = a.ID
	m.CustomerID = a.CustomerID
	m.StoreID = a.StoreID
	m.SlotID = a.SlotID
	m.TechnicianID = a.TechnicianID
	m.ServiceName = a.ServiceName
	m.ServicePrice = a.ServicePrice
	m.DepositAmount = a.DepositAmount
	m.DepositPaid = a.DepositPaid
	m.DepositPaidAt = a.DepositPaidAt
	m.DepositRefunded = a.DepositRefunded
	m.DepositRefundedAt = a.DepositRefundedAt
	m.Status = a.Status
	m.CancelledAt = a.CancelledAt
	m.CancelledBy = a.CancelledBy
	m.CancelReason = a.CancelReason
	m.Remark = a.Remark
	m.AppointmentTime = a.AppointmentTime
	m.CreatedAt = a.CreatedAt
	m.UpdatedAt = a.UpdatedAt
}

// FindByID 根据ID查找预约
func (r *AppointmentRepository) FindByID(ctx context.Context, id uint) (*appointment.Appointment, error) {
	r.logger.Debug("查找预约：根据ID", logger.NewField("appointment_id", id))

	var model AppointmentModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("预约不存在：根据ID", logger.NewField("appointment_id", id))
			return nil, nil
		}
		r.logger.Error("查找预约失败：根据ID",
			logger.NewField("appointment_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindByCustomerID 根据顾客ID查找预约列表
func (r *AppointmentRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据顾客ID", logger.NewField("customer_id", customerID))

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("appointment_time DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据顾客ID",
			logger.NewField("customer_id", customerID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindByStoreID 根据门店ID查找预约列表
func (r *AppointmentRepository) FindByStoreID(ctx context.Context, storeID uint) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据门店ID", logger.NewField("store_id", storeID))

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("store_id = ?", storeID).
		Order("appointment_time DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据门店ID",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindByTechnicianID 根据美甲师ID查找预约列表
func (r *AppointmentRepository) FindByTechnicianID(ctx context.Context, technicianID uint) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据美甲师ID", logger.NewField("technician_id", technicianID))

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("technician_id = ?", technicianID).
		Order("appointment_time DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据美甲师ID",
			logger.NewField("technician_id", technicianID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindBySlotID 根据时段ID查找预约列表
func (r *AppointmentRepository) FindBySlotID(ctx context.Context, slotID uint) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据时段ID", logger.NewField("slot_id", slotID))

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("slot_id = ?", slotID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据时段ID",
			logger.NewField("slot_id", slotID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindByStatus 根据状态查找预约列表
func (r *AppointmentRepository) FindByStatus(ctx context.Context, status string) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据状态", logger.NewField("status", status))

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("status = ?", status).
		Order("appointment_time DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据状态",
			logger.NewField("status", status),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindByDateRange 根据日期范围查找预约列表
func (r *AppointmentRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找预约列表：根据日期范围",
		logger.NewField("start_date", startDate.Format("2006-01-02")),
		logger.NewField("end_date", endDate.Format("2006-01-02")),
	)

	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("appointment_time >= ? AND appointment_time <= ?", startDate, endDate).
		Order("appointment_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找预约列表失败：根据日期范围",
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// FindUpcomingByCustomerID 查找顾客即将到来的预约（未取消、未完成）
func (r *AppointmentRepository) FindUpcomingByCustomerID(ctx context.Context, customerID uint) ([]*appointment.Appointment, error) {
	r.logger.Debug("查找即将到来的预约：根据顾客ID", logger.NewField("customer_id", customerID))

	now := time.Now()
	var models []AppointmentModel
	if err := r.db.DB().WithContext(ctx).
		Where("customer_id = ? AND status IN ? AND appointment_time >= ?",
			customerID,
			[]string{appointment.AppointmentStatusPending, appointment.AppointmentStatusPaid, appointment.AppointmentStatusConfirmed},
			now).
		Order("appointment_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找即将到来的预约失败：根据顾客ID",
			logger.NewField("customer_id", customerID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	appointments := make([]*appointment.Appointment, 0, len(models))
	for _, model := range models {
		appointments = append(appointments, model.ToEntity())
	}

	return appointments, nil
}

// Create 创建预约
func (r *AppointmentRepository) Create(ctx context.Context, a *appointment.Appointment) error {
	r.logger.Info("创建预约",
		logger.NewField("customer_id", a.CustomerID),
		logger.NewField("slot_id", a.SlotID),
	)

	model := &AppointmentModel{}
	model.FromEntity(a)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建预约失败",
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	a.ID = model.ID
	a.CreatedAt = model.CreatedAt
	a.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新预约
func (r *AppointmentRepository) Update(ctx context.Context, a *appointment.Appointment) error {
	r.logger.Info("更新预约", logger.NewField("appointment_id", a.ID))

	model := &AppointmentModel{}
	model.FromEntity(a)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&AppointmentModel{}).Where("id = ?", a.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新预约失败",
			logger.NewField("appointment_id", a.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	a.UpdatedAt = model.UpdatedAt
	return nil
}

// UpdateStatus 更新预约状态
func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	r.logger.Info("更新预约状态", logger.NewField("appointment_id", id), logger.NewField("status", status))

	if err := r.db.DB().WithContext(ctx).
		Model(&AppointmentModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		r.logger.Error("更新预约状态失败",
			logger.NewField("appointment_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	return nil
}

// Delete 删除预约
func (r *AppointmentRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除预约", logger.NewField("appointment_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&AppointmentModel{}, id).Error; err != nil {
		r.logger.Error("删除预约失败",
			logger.NewField("appointment_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	return nil
}

// ==================== 支付仓储 ====================

// PaymentRepository 支付仓储实现
type PaymentRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewPaymentRepository 创建支付仓储实例
func NewPaymentRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *PaymentRepository {
	return &PaymentRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// PaymentModel 支付数据库模型
type PaymentModel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentID uint           `gorm:"index;not null" json:"appointment_id"`
	OrderNo       string         `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	PaymentType   string         `gorm:"size:20;not null" json:"payment_type"`
	Amount        float64        `gorm:"not null" json:"amount"`
	PaymentMethod string         `gorm:"size:20;not null" json:"payment_method"`
	Status        string         `gorm:"size:20;default:'pending';not null" json:"status"`
	TransactionID string         `gorm:"size:100" json:"transaction_id"`
	PaymentParams string         `gorm:"type:text" json:"payment_params"`
	PaymentResult string         `gorm:"type:text" json:"payment_result"`
	ExpiredAt     time.Time      `json:"expired_at"`
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PaymentModel) TableName() string {
	return "payments"
}

// ToEntity 转换为领域实体
func (m *PaymentModel) ToEntity() *appointment.Payment {
	if m == nil {
		return nil
	}
	return &appointment.Payment{
		ID:            m.ID,
		AppointmentID: m.AppointmentID,
		OrderNo:       m.OrderNo,
		PaymentType:   m.PaymentType,
		Amount:        m.Amount,
		PaymentMethod: m.PaymentMethod,
		Status:        m.Status,
		TransactionID: m.TransactionID,
		PaymentParams: m.PaymentParams,
		PaymentResult: m.PaymentResult,
		ExpiredAt:     m.ExpiredAt,
		PaidAt:        m.PaidAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *PaymentModel) FromEntity(p *appointment.Payment) {
	if p == nil {
		return
	}
	m.ID = p.ID
	m.AppointmentID = p.AppointmentID
	m.OrderNo = p.OrderNo
	m.PaymentType = p.PaymentType
	m.Amount = p.Amount
	m.PaymentMethod = p.PaymentMethod
	m.Status = p.Status
	m.TransactionID = p.TransactionID
	m.PaymentParams = p.PaymentParams
	m.PaymentResult = p.PaymentResult
	m.ExpiredAt = p.ExpiredAt
	m.PaidAt = p.PaidAt
	m.CreatedAt = p.CreatedAt
	m.UpdatedAt = p.UpdatedAt
}

// FindByID 根据ID查找支付记录
func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*appointment.Payment, error) {
	r.logger.Debug("查找支付记录：根据ID", logger.NewField("payment_id", id))

	var model PaymentModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("支付记录不存在：根据ID", logger.NewField("payment_id", id))
			return nil, nil
		}
		r.logger.Error("查找支付记录失败：根据ID",
			logger.NewField("payment_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindByAppointmentID 根据预约ID查找支付记录
func (r *PaymentRepository) FindByAppointmentID(ctx context.Context, appointmentID uint) (*appointment.Payment, error) {
	r.logger.Debug("查找支付记录：根据预约ID", logger.NewField("appointment_id", appointmentID))

	var model PaymentModel
	if err := r.db.DB().WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("支付记录不存在：根据预约ID", logger.NewField("appointment_id", appointmentID))
			return nil, nil
		}
		r.logger.Error("查找支付记录失败：根据预约ID",
			logger.NewField("appointment_id", appointmentID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindByOrderNo 根据订单号查找支付记录
func (r *PaymentRepository) FindByOrderNo(ctx context.Context, orderNo string) (*appointment.Payment, error) {
	r.logger.Debug("查找支付记录：根据订单号", logger.NewField("order_no", orderNo))

	var model PaymentModel
	if err := r.db.DB().WithContext(ctx).Where("order_no = ?", orderNo).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("支付记录不存在：根据订单号", logger.NewField("order_no", orderNo))
			return nil, nil
		}
		r.logger.Error("查找支付记录失败：根据订单号",
			logger.NewField("order_no", orderNo),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// Create 创建支付记录
func (r *PaymentRepository) Create(ctx context.Context, p *appointment.Payment) error {
	r.logger.Info("创建支付记录",
		logger.NewField("appointment_id", p.AppointmentID),
		logger.NewField("order_no", p.OrderNo),
	)

	model := &PaymentModel{}
	model.FromEntity(p)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建支付记录失败",
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新支付记录
func (r *PaymentRepository) Update(ctx context.Context, p *appointment.Payment) error {
	r.logger.Info("更新支付记录", logger.NewField("payment_id", p.ID))

	model := &PaymentModel{}
	model.FromEntity(p)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&PaymentModel{}).Where("id = ?", p.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新支付记录失败",
			logger.NewField("payment_id", p.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	p.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除支付记录
func (r *PaymentRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除支付记录", logger.NewField("payment_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&PaymentModel{}, id).Error; err != nil {
		r.logger.Error("删除支付记录失败",
			logger.NewField("payment_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	return nil
}

// ==================== 退款仓储 ====================

// RefundRepository 退款仓储实现
type RefundRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewRefundRepository 创建退款仓储实例
func NewRefundRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *RefundRepository {
	return &RefundRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// RefundModel 退款数据库模型
type RefundModel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppointmentID uint           `gorm:"index;not null" json:"appointment_id"`
	PaymentID     uint           `gorm:"index;not null" json:"payment_id"`
	OrderNo       string         `gorm:"size:50;not null" json:"order_no"`
	RefundNo      string         `gorm:"uniqueIndex;size:50;not null" json:"refund_no"`
	RefundType    string         `gorm:"size:20;not null" json:"refund_type"`
	Amount        float64        `gorm:"not null" json:"amount"`
	Status        string         `gorm:"size:20;default:'pending';not null" json:"status"`
	TransactionID string         `gorm:"size:100" json:"transaction_id"`
	RefundResult  string         `gorm:"type:text" json:"refund_result"`
	RefundedAt    *time.Time     `json:"refunded_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RefundModel) TableName() string {
	return "refunds"
}

// ToEntity 转换为领域实体
func (m *RefundModel) ToEntity() *appointment.Refund {
	if m == nil {
		return nil
	}
	return &appointment.Refund{
		ID:            m.ID,
		AppointmentID: m.AppointmentID,
		PaymentID:     m.PaymentID,
		OrderNo:       m.OrderNo,
		RefundNo:      m.RefundNo,
		RefundType:    m.RefundType,
		Amount:        m.Amount,
		Status:        m.Status,
		TransactionID: m.TransactionID,
		RefundResult:  m.RefundResult,
		RefundedAt:    m.RefundedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *RefundModel) FromEntity(r *appointment.Refund) {
	if r == nil {
		return
	}
	m.ID = r.ID
	m.AppointmentID = r.AppointmentID
	m.PaymentID = r.PaymentID
	m.OrderNo = r.OrderNo
	m.RefundNo = r.RefundNo
	m.RefundType = r.RefundType
	m.Amount = r.Amount
	m.Status = r.Status
	m.TransactionID = r.TransactionID
	m.RefundResult = r.RefundResult
	m.RefundedAt = r.RefundedAt
	m.CreatedAt = r.CreatedAt
	m.UpdatedAt = r.UpdatedAt
}

// FindByID 根据ID查找退款记录
func (r *RefundRepository) FindByID(ctx context.Context, id uint) (*appointment.Refund, error) {
	r.logger.Debug("查找退款记录：根据ID", logger.NewField("refund_id", id))

	var model RefundModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("退款记录不存在：根据ID", logger.NewField("refund_id", id))
			return nil, nil
		}
		r.logger.Error("查找退款记录失败：根据ID",
			logger.NewField("refund_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindByAppointmentID 根据预约ID查找退款记录
func (r *RefundRepository) FindByAppointmentID(ctx context.Context, appointmentID uint) ([]*appointment.Refund, error) {
	r.logger.Debug("查找退款记录列表：根据预约ID", logger.NewField("appointment_id", appointmentID))

	var models []RefundModel
	if err := r.db.DB().WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找退款记录列表失败：根据预约ID",
			logger.NewField("appointment_id", appointmentID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	refunds := make([]*appointment.Refund, 0, len(models))
	for _, model := range models {
		refunds = append(refunds, model.ToEntity())
	}

	return refunds, nil
}

// FindByPaymentID 根据支付ID查找退款记录
func (r *RefundRepository) FindByPaymentID(ctx context.Context, paymentID uint) ([]*appointment.Refund, error) {
	r.logger.Debug("查找退款记录列表：根据支付ID", logger.NewField("payment_id", paymentID))

	var models []RefundModel
	if err := r.db.DB().WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找退款记录列表失败：根据支付ID",
			logger.NewField("payment_id", paymentID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	refunds := make([]*appointment.Refund, 0, len(models))
	for _, model := range models {
		refunds = append(refunds, model.ToEntity())
	}

	return refunds, nil
}

// FindByRefundNo 根据退款单号查找退款记录
func (r *RefundRepository) FindByRefundNo(ctx context.Context, refundNo string) (*appointment.Refund, error) {
	r.logger.Debug("查找退款记录：根据退款单号", logger.NewField("refund_no", refundNo))

	var model RefundModel
	if err := r.db.DB().WithContext(ctx).Where("refund_no = ?", refundNo).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("退款记录不存在：根据退款单号", logger.NewField("refund_no", refundNo))
			return nil, nil
		}
		r.logger.Error("查找退款记录失败：根据退款单号",
			logger.NewField("refund_no", refundNo),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// Create 创建退款记录
func (r *RefundRepository) Create(ctx context.Context, ref *appointment.Refund) error {
	r.logger.Info("创建退款记录",
		logger.NewField("appointment_id", ref.AppointmentID),
		logger.NewField("refund_no", ref.RefundNo),
	)

	model := &RefundModel{}
	model.FromEntity(ref)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建退款记录失败",
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	ref.ID = model.ID
	ref.CreatedAt = model.CreatedAt
	ref.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新退款记录
func (r *RefundRepository) Update(ctx context.Context, ref *appointment.Refund) error {
	r.logger.Info("更新退款记录", logger.NewField("refund_id", ref.ID))

	model := &RefundModel{}
	model.FromEntity(ref)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&RefundModel{}).Where("id = ?", ref.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新退款记录失败",
			logger.NewField("refund_id", ref.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	ref.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除退款记录
func (r *RefundRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除退款记录", logger.NewField("refund_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&RefundModel{}, id).Error; err != nil {
		r.logger.Error("删除退款记录失败",
			logger.NewField("refund_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	return nil
}

