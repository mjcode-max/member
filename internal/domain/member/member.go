package member

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 服务类型常量
const (
	ServiceTypeNail    = "nail"    // 美甲
	ServiceTypeEyelash = "eyelash" // 美睫
	ServiceTypeCombo   = "combo"   // 组合
)

// 会员状态常量
const (
	StatusActive   = "active"   // 有效
	StatusExpired  = "expired"  // 过期
	StatusInactive = "inactive" // 停用
)

// Member 会员实体（包含套餐信息）
type Member struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`              // 会员姓名
	Phone           string    `json:"phone"`             // 手机号
	PackageName     string    `json:"package_name"`     // 套餐名称
	ServiceType     string    `json:"service_type"`      // 服务类型: nail, eyelash, combo
	Price           float64   `json:"price"`             // 套餐价格
	UsedTimes       int       `json:"used_times"`        // 已使用次数
	ValidityDuration int      `json:"validity_duration"` // 固定时长天数
	ValidFrom       time.Time `json:"valid_from"`       // 有效期开始时间
	ValidTo         time.Time `json:"valid_to"`        // 有效期结束时间
	StoreID         uint      `json:"store_id"`         // 购买门店ID
	PurchaseAmount  float64   `json:"purchase_amount"`   // 购买金额
	PurchaseTime    time.Time `json:"purchase_time"`     // 购买时间
	Status          string    `json:"status"`           // 状态: active, expired, inactive
	Description     string    `json:"description"`      // 套餐描述/备注
	FaceID          string    `json:"face_id"`         // 华为云人脸ID
	CreatedBy       uint      `json:"created_by"`       // 创建人ID
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// IMemberRepository 会员仓储接口
type IMemberRepository interface {
	// FindByID 根据ID查找会员
	FindByID(ctx context.Context, id uint) (*Member, error)
	// FindList 获取会员列表（支持筛选和分页）
	FindList(ctx context.Context, name, phone string, storeID *uint, status, serviceType, packageName string, page, pageSize int) ([]*Member, int64, error)
	// Create 创建会员
	Create(ctx context.Context, member *Member) error
	// Update 更新会员
	Update(ctx context.Context, member *Member) error
	// Delete 删除会员（软删除）
	Delete(ctx context.Context, id uint) error
	// FindByPhone 根据手机号查找会员列表
	FindByPhone(ctx context.Context, phone string) ([]*Member, error)
	// IncrementUsedTimes 递增会员已使用次数
	IncrementUsedTimes(ctx context.Context, id uint) error
	// DecrementUsedTimes 递减会员已使用次数（不能小于0）
	DecrementUsedTimes(ctx context.Context, id uint) error
}

// 领域错误定义
var (
	ErrMemberNotFound      = errors.ErrNotFound("会员不存在")
	ErrInvalidStatus       = errors.ErrInvalidParams("无效的会员状态")
	ErrInvalidServiceType  = errors.ErrInvalidParams("无效的服务类型")
	ErrNameRequired        = errors.ErrInvalidParams("会员姓名不能为空")
	ErrPhoneRequired       = errors.ErrInvalidParams("手机号不能为空")
	ErrPackageNameRequired = errors.ErrInvalidParams("套餐名称不能为空")
	ErrInvalidPrice        = errors.ErrInvalidParams("套餐价格不能为负数")
	ErrInvalidValidity      = errors.ErrInvalidParams("有效期结束时间必须晚于开始时间")
	ErrMemberExpired       = errors.ErrForbidden("会员已过期")
	ErrMemberInactive      = errors.ErrForbidden("会员已停用")
)

// MemberService 会员服务
type MemberService struct {
	repo   IMemberRepository
	logger logger.Logger
}

// NewMemberService 创建会员服务
func NewMemberService(repo IMemberRepository, log logger.Logger) *MemberService {
	return &MemberService{
		repo:   repo,
		logger: log,
	}
}

// GetByID 根据ID获取会员
func (s *MemberService) GetByID(ctx context.Context, id uint) (*Member, error) {
	s.logger.Debug("获取会员", logger.NewField("member_id", id))

	member, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找会员失败", logger.NewField("member_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if member == nil {
		s.logger.Warn("会员不存在", logger.NewField("member_id", id))
		return nil, ErrMemberNotFound
	}

	return member, nil
}

// GetList 获取会员列表（支持多条件筛选）
func (s *MemberService) GetList(ctx context.Context, name, phone string, storeID *uint, status, serviceType, packageName string, page, pageSize int) ([]*Member, int64, error) {
	s.logger.Debug("获取会员列表",
		logger.NewField("name", name),
		logger.NewField("phone", phone),
		logger.NewField("store_id", storeID),
		logger.NewField("status", status),
		logger.NewField("service_type", serviceType),
		logger.NewField("package_name", packageName),
		logger.NewField("page", page),
		logger.NewField("page_size", pageSize),
	)

	members, total, err := s.repo.FindList(ctx, name, phone, storeID, status, serviceType, packageName, page, pageSize)
	if err != nil {
		s.logger.Error("获取会员列表失败", logger.NewField("error", err.Error()))
		return nil, 0, err
	}

	return members, total, nil
}

// CreateMember 创建会员（购买记录，包含套餐信息）
func (s *MemberService) CreateMember(ctx context.Context, member *Member) error {
	s.logger.Info("创建会员", logger.NewField("name", member.Name))

	// 验证必填字段
	if member.Name == "" {
		return ErrNameRequired
	}
	if member.Phone == "" {
		return ErrPhoneRequired
	}
	if member.PackageName == "" {
		return ErrPackageNameRequired
	}
	if !isValidServiceType(member.ServiceType) {
		return ErrInvalidServiceType
	}
	if member.Price < 0 {
		return ErrInvalidPrice
	}

	// 有效期处理
	if err := s.processValidity(member); err != nil {
		return err
	}

	// 初始化已使用次数
	member.UsedTimes = 0

	// 设置状态
	if member.Status == "" {
		member.Status = StatusActive
	}
	if !isValidStatus(member.Status) {
		return ErrInvalidStatus
	}

	// 设置购买时间
	if member.PurchaseTime.IsZero() {
		member.PurchaseTime = time.Now()
	}

	// 设置时间戳
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, member); err != nil {
		s.logger.Error("创建会员失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("创建会员成功", logger.NewField("member_id", member.ID))
	return nil
}

// UpdateMember 更新会员信息（包括套餐信息）
func (s *MemberService) UpdateMember(ctx context.Context, member *Member) error {
	s.logger.Info("更新会员", logger.NewField("member_id", member.ID))

	// 检查会员是否存在
	existing, err := s.repo.FindByID(ctx, member.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrMemberNotFound
	}

	// 验证服务类型
	if member.ServiceType != "" && !isValidServiceType(member.ServiceType) {
		return ErrInvalidServiceType
	}

	// 验证价格
	if member.Price < 0 {
		return ErrInvalidPrice
	}

	// 有效期处理（如果提供了有效期相关字段）
	if member.ValidityDuration > 0 || !member.ValidFrom.IsZero() || !member.ValidTo.IsZero() {
		// 如果只更新了部分字段，需要从现有记录中获取其他字段
		if member.ValidityDuration == 0 {
			member.ValidityDuration = existing.ValidityDuration
		}
		if member.ValidFrom.IsZero() {
			member.ValidFrom = existing.ValidFrom
		}
		if member.ValidTo.IsZero() {
			member.ValidTo = existing.ValidTo
		}
		if err := s.processValidity(member); err != nil {
			return err
		}
	}

	// 验证状态
	if member.Status != "" && !isValidStatus(member.Status) {
		return ErrInvalidStatus
	}

	member.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, member); err != nil {
		s.logger.Error("更新会员失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("更新会员成功", logger.NewField("member_id", member.ID))
	return nil
}

// GetMemberByPhone 根据手机号获取会员列表
func (s *MemberService) GetMemberByPhone(ctx context.Context, phone string) ([]*Member, error) {
	s.logger.Debug("根据手机号获取会员", logger.NewField("phone", phone))

	if phone == "" {
		return nil, ErrPhoneRequired
	}

	members, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		s.logger.Error("根据手机号查找会员失败", logger.NewField("phone", phone), logger.NewField("error", err.Error()))
		return nil, err
	}

	return members, nil
}

// UpdateMemberStatus 更新会员状态
func (s *MemberService) UpdateMemberStatus(ctx context.Context, id uint, status string) error {
	s.logger.Info("更新会员状态", logger.NewField("member_id", id), logger.NewField("status", status))

	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	member, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotFound
	}

	member.Status = status
	member.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, member); err != nil {
		s.logger.Error("更新会员状态失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("更新会员状态成功", logger.NewField("member_id", id), logger.NewField("status", status))
	return nil
}

// CheckMemberValidity 检查会员有效期并更新状态
func (s *MemberService) CheckMemberValidity(ctx context.Context, id uint) error {
	member, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotFound
	}

	now := time.Now()
	if member.Status == StatusActive && now.After(member.ValidTo) {
		member.Status = StatusExpired
		member.UpdatedAt = time.Now()
		return s.repo.Update(ctx, member)
	}

	return nil
}

// IsValid 判断会员是否有效（状态为active且未过期，且在有效期内）
func (m *Member) IsValid() bool {
	if m.Status != StatusActive {
		return false
	}
	now := time.Now()
	return !now.Before(m.ValidFrom) && !now.After(m.ValidTo)
}

// CalculateValidityFromDuration 根据固定时长计算结束日期
func (m *Member) CalculateValidityFromDuration() {
	if m.ValidityDuration > 0 && !m.ValidFrom.IsZero() {
		m.ValidTo = m.ValidFrom.AddDate(0, 0, m.ValidityDuration)
	}
}

// CalculateDurationFromDates 根据开始和结束日期计算时长
func (m *Member) CalculateDurationFromDates() {
	if !m.ValidFrom.IsZero() && !m.ValidTo.IsZero() {
		duration := m.ValidTo.Sub(m.ValidFrom)
		m.ValidityDuration = int(duration.Hours() / 24)
	}
}

// IncrementUsedTimes 增加已使用次数
func (m *Member) IncrementUsedTimes() {
	m.UsedTimes++
}

// processValidity 处理有效期（支持固定时长和日期两种方式互相转换）
func (s *MemberService) processValidity(member *Member) error {
	// 如果输入了validity_duration（固定时长）
	if member.ValidityDuration > 0 {
		// valid_from默认为当前日期（如果未指定）
		if member.ValidFrom.IsZero() {
			member.ValidFrom = time.Now()
		}
		// 计算valid_to
		member.ValidTo = member.ValidFrom.AddDate(0, 0, member.ValidityDuration)
	} else if !member.ValidFrom.IsZero() && !member.ValidTo.IsZero() {
		// 如果输入了valid_from和valid_to（手动指定日期）
		// 验证valid_to > valid_from
		if !member.ValidTo.After(member.ValidFrom) {
			return ErrInvalidValidity
		}
		// 自动计算validity_duration
		member.CalculateDurationFromDates()
	}

	// 如果同时输入了两种方式，优先使用日期方式，并重新计算validity_duration
	if member.ValidityDuration > 0 && !member.ValidFrom.IsZero() && !member.ValidTo.IsZero() {
		// 优先使用日期，重新计算时长
		member.CalculateDurationFromDates()
	}

	return nil
}

// isValidStatus 验证会员状态是否有效
func isValidStatus(status string) bool {
	return status == StatusActive || status == StatusExpired || status == StatusInactive
}

// isValidServiceType 验证服务类型是否有效
func isValidServiceType(serviceType string) bool {
	return serviceType == ServiceTypeNail || serviceType == ServiceTypeEyelash || serviceType == ServiceTypeCombo
}

