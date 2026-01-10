package store

import (
	"context"
	"regexp"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 门店状态常量
const (
	StatusOperating = "operating" // 营业中
	StatusClosed    = "closed"    // 停业
	StatusShutdown  = "shutdown"  // 关闭
)

// Store 门店实体
type Store struct {
	ID                 uint      `json:"id"`
	Name               string    `json:"name"`                 // 门店名称
	Address            string    `json:"address"`              // 门店地址
	Phone              string    `json:"phone"`                // 联系电话
	ContactPerson      string    `json:"contact_person"`       // 联系人
	Status             string    `json:"status"`               // 状态: operating, closed, shutdown
	BusinessHoursStart string    `json:"business_hours_start"` // 营业开始时间 (HH:MM格式)
	BusinessHoursEnd   string    `json:"business_hours_end"`   // 营业结束时间 (HH:MM格式)
	DepositAmount      float64   `json:"deposit_amount"`       // 押金金额
	TemplateID         *uint     `json:"template_id"`          // 时段模板ID（可选）
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// IStoreRepository 门店仓储接口
type IStoreRepository interface {
	// FindByID 根据ID查找门店
	FindByID(ctx context.Context, id uint) (*Store, error)
	// FindList 获取门店列表（支持筛选和分页）
	FindList(ctx context.Context, status, name string, page, pageSize int) ([]*Store, int64, error)
	// Create 创建门店
	Create(ctx context.Context, store *Store) error
	// Update 更新门店
	Update(ctx context.Context, store *Store) error
	// Delete 删除门店（软删除）
	Delete(ctx context.Context, id uint) error
	// FindByManagerID 根据店长ID查找门店
	FindByManagerID(ctx context.Context, managerID uint) (*Store, error)
	// FindByStatus 根据状态查找门店列表
	FindByStatus(ctx context.Context, status string) ([]*Store, error)
}

// 领域错误定义
var (
	ErrStoreNotFound     = errors.ErrNotFound("门店不存在")
	ErrInvalidStatus     = errors.ErrInvalidParams("无效的门店状态")
	ErrInvalidTimeFormat = errors.ErrInvalidParams("营业时间格式错误，应为HH:MM格式")
	ErrInvalidDeposit    = errors.ErrInvalidParams("押金金额不能为负数")
	ErrNameRequired      = errors.ErrInvalidParams("门店名称不能为空")
)

// StoreService 门店服务
type StoreService struct {
	repo   IStoreRepository
	logger logger.Logger
}

// NewStoreService 创建门店服务
func NewStoreService(repo IStoreRepository, log logger.Logger) *StoreService {
	return &StoreService{
		repo:   repo,
		logger: log,
	}
}

// GetByID 根据ID获取门店
func (s *StoreService) GetByID(ctx context.Context, id uint) (*Store, error) {
	s.logger.Debug("获取门店", logger.NewField("store_id", id))

	store, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找门店失败", logger.NewField("store_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if store == nil {
		s.logger.Warn("门店不存在", logger.NewField("store_id", id))
		return nil, ErrStoreNotFound
	}

	return store, nil
}

// GetList 获取门店列表（支持筛选和分页）
func (s *StoreService) GetList(ctx context.Context, status, name string, page, pageSize int) ([]*Store, int64, error) {
	s.logger.Debug("获取门店列表",
		logger.NewField("status", status),
		logger.NewField("name", name),
		logger.NewField("page", page),
		logger.NewField("page_size", pageSize),
	)

	stores, total, err := s.repo.FindList(ctx, status, name, page, pageSize)
	if err != nil {
		s.logger.Error("获取门店列表失败", logger.NewField("error", err.Error()))
		return nil, 0, err
	}

	return stores, total, nil
}

// Create 创建门店
func (s *StoreService) Create(ctx context.Context, store *Store) error {
	s.logger.Info("创建门店",
		logger.NewField("name", store.Name),
	)

	// 验证门店名称
	if store.Name == "" {
		return ErrNameRequired
	}

	// 验证状态
	if store.Status != "" && !isValidStatus(store.Status) {
		return ErrInvalidStatus
	}

	// 验证营业时间格式
	if store.BusinessHoursStart != "" {
		if !isValidTimeFormat(store.BusinessHoursStart) {
			return ErrInvalidTimeFormat
		}
	}
	if store.BusinessHoursEnd != "" {
		if !isValidTimeFormat(store.BusinessHoursEnd) {
			return ErrInvalidTimeFormat
		}
	}

	// 验证押金金额
	if store.DepositAmount < 0 {
		return ErrInvalidDeposit
	}

	// 设置默认值
	if store.Status == "" {
		store.Status = StatusOperating
	}
	store.CreatedAt = time.Now()
	store.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, store); err != nil {
		s.logger.Error("创建门店失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("创建门店成功", logger.NewField("store_id", store.ID))
	return nil
}

// Update 更新门店
func (s *StoreService) Update(ctx context.Context, store *Store) error {
	s.logger.Info("更新门店", logger.NewField("store_id", store.ID))

	// 检查门店是否存在
	existing, err := s.repo.FindByID(ctx, store.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrStoreNotFound
	}

	// 验证状态
	if store.Status != "" && !isValidStatus(store.Status) {
		return ErrInvalidStatus
	}

	// 验证营业时间格式
	if store.BusinessHoursStart != "" {
		if !isValidTimeFormat(store.BusinessHoursStart) {
			return ErrInvalidTimeFormat
		}
	}
	if store.BusinessHoursEnd != "" {
		if !isValidTimeFormat(store.BusinessHoursEnd) {
			return ErrInvalidTimeFormat
		}
	}

	// 验证押金金额
	if store.DepositAmount < 0 {
		return ErrInvalidDeposit
	}

	// 保留原有字段（如果新字段为空）
	if store.Name == "" {
		store.Name = existing.Name
	}
	if store.Address == "" {
		store.Address = existing.Address
	}
	if store.Phone == "" {
		store.Phone = existing.Phone
	}
	if store.ContactPerson == "" {
		store.ContactPerson = existing.ContactPerson
	}
	if store.Status == "" {
		store.Status = existing.Status
	}
	if store.BusinessHoursStart == "" {
		store.BusinessHoursStart = existing.BusinessHoursStart
	}
	if store.BusinessHoursEnd == "" {
		store.BusinessHoursEnd = existing.BusinessHoursEnd
	}
	if store.DepositAmount == 0 && existing.DepositAmount != 0 {
		store.DepositAmount = existing.DepositAmount
	}
	// TemplateID为nil时保留原有值，否则更新为新值
	if store.TemplateID == nil {
		store.TemplateID = existing.TemplateID
	}

	store.UpdatedAt = time.Now()

	// 检测配置变更（用于通知预约系统）
	configChanged := s.detectConfigChange(existing, store)

	if err := s.repo.Update(ctx, store); err != nil {
		s.logger.Error("更新门店失败", logger.NewField("error", err.Error()))
		return err
	}

	// 如果配置变更，记录日志（后续可以扩展为事件通知）
	if configChanged {
		s.logger.Info("门店配置已变更，可能影响预约时段生成",
			logger.NewField("store_id", store.ID),
			logger.NewField("business_hours_start", store.BusinessHoursStart),
			logger.NewField("business_hours_end", store.BusinessHoursEnd),
			logger.NewField("status", store.Status),
		)
	}

	s.logger.Info("更新门店成功", logger.NewField("store_id", store.ID))
	return nil
}

// Delete 删除门店（软删除）
func (s *StoreService) Delete(ctx context.Context, id uint) error {
	s.logger.Info("删除门店", logger.NewField("store_id", id))

	// 检查门店是否存在
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrStoreNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("删除门店失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("删除门店成功", logger.NewField("store_id", id))
	return nil
}

// UpdateStatus 更新门店状态
func (s *StoreService) UpdateStatus(ctx context.Context, id uint, status string) error {
	s.logger.Info("更新门店状态", logger.NewField("store_id", id), logger.NewField("status", status))

	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	store, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if store == nil {
		return ErrStoreNotFound
	}

	oldStatus := store.Status
	store.Status = status
	store.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, store); err != nil {
		s.logger.Error("更新门店状态失败", logger.NewField("error", err.Error()))
		return err
	}

	// 如果状态变更，记录日志（后续可以扩展为事件通知）
	if oldStatus != status {
		s.logger.Info("门店状态已变更，可能影响预约时段生成",
			logger.NewField("store_id", id),
			logger.NewField("old_status", oldStatus),
			logger.NewField("new_status", status),
		)
	}

	s.logger.Info("更新门店状态成功", logger.NewField("store_id", id), logger.NewField("status", status))
	return nil
}

// GetByManagerID 根据店长ID获取门店
func (s *StoreService) GetByManagerID(ctx context.Context, managerID uint) (*Store, error) {
	s.logger.Debug("根据店长ID获取门店", logger.NewField("manager_id", managerID))

	store, err := s.repo.FindByManagerID(ctx, managerID)
	if err != nil {
		s.logger.Error("查找门店失败", logger.NewField("manager_id", managerID), logger.NewField("error", err.Error()))
		return nil, err
	}

	if store == nil {
		s.logger.Warn("门店不存在", logger.NewField("manager_id", managerID))
		return nil, ErrStoreNotFound
	}

	return store, nil
}

// GetByStatus 根据状态获取门店列表
func (s *StoreService) GetByStatus(ctx context.Context, status string) ([]*Store, error) {
	s.logger.Debug("根据状态获取门店列表", logger.NewField("status", status))

	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	stores, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		s.logger.Error("查找门店列表失败", logger.NewField("error", err.Error()))
		return nil, err
	}

	return stores, nil
}

// isValidStatus 验证门店状态是否有效
func isValidStatus(status string) bool {
	return status == StatusOperating || status == StatusClosed || status == StatusShutdown
}

// isValidTimeFormat 验证时间格式是否为HH:MM
func isValidTimeFormat(timeStr string) bool {
	matched, _ := regexp.MatchString(`^([0-1][0-9]|2[0-3]):[0-5][0-9]$`, timeStr)
	return matched
}

// detectConfigChange 检测配置变更（影响预约时段生成的配置）
func (s *StoreService) detectConfigChange(oldStore, newStore *Store) bool {
	return oldStore.BusinessHoursStart != newStore.BusinessHoursStart ||
		oldStore.BusinessHoursEnd != newStore.BusinessHoursEnd ||
		oldStore.Status != newStore.Status
}
