package member

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// MemberUsage 会员使用记录实体
type MemberUsage struct {
	ID             uint      `json:"id"`
	MemberID       uint      `json:"member_id"`        // 会员ID
	PackageName    string    `json:"package_name"`     // 套餐名称（冗余字段）
	ServiceItem    string    `json:"service_item"`     // 服务项目
	StoreID        uint      `json:"store_id"`         // 使用门店ID
	StoreName      string    `json:"store_name"`       // 门店名称（冗余字段）
	TechnicianID   *uint     `json:"technician_id"`    // 美甲师ID
	TechnicianName string    `json:"technician_name"`  // 美甲师姓名（冗余字段）
	UsageDate      time.Time `json:"usage_date"`       // 使用日期
	Remark         string    `json:"remark"`           // 备注
	CreatedBy      uint      `json:"created_by"`       // 操作人ID
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// IUsageRepository 使用记录仓储接口
type IUsageRepository interface {
	// FindByID 根据ID查找使用记录
	FindByID(ctx context.Context, id uint) (*MemberUsage, error)
	// FindList 获取使用记录列表（不分页，支持筛选）
	FindList(ctx context.Context, memberID *uint, storeID *uint, technicianID *uint) ([]*MemberUsage, error)
	// Create 创建使用记录
	Create(ctx context.Context, usage *MemberUsage) error
	// Delete 删除使用记录
	Delete(ctx context.Context, id uint) error
	// FindByMemberID 根据会员ID查找使用记录列表
	FindByMemberID(ctx context.Context, memberID uint) ([]*MemberUsage, error)
}

// 使用记录领域错误定义
var (
	ErrUsageNotFound = errors.ErrNotFound("使用记录不存在")
)

// UsageService 使用记录服务
type UsageService struct {
	repo       IUsageRepository
	memberRepo IMemberRepository
	logger     logger.Logger
}

// NewUsageService 创建使用记录服务
func NewUsageService(repo IUsageRepository, memberRepo IMemberRepository, log logger.Logger) *UsageService {
	return &UsageService{
		repo:       repo,
		memberRepo: memberRepo,
		logger:     log,
	}
}

// CreateUsage 记录使用，递增会员已使用次数（used_times）
func (s *UsageService) CreateUsage(ctx context.Context, usage *MemberUsage) error {
	s.logger.Info("创建使用记录", logger.NewField("member_id", usage.MemberID))

	// 验证会员是否存在且有效
	member, err := s.memberRepo.FindByID(ctx, usage.MemberID)
	if err != nil {
		s.logger.Error("查找会员失败", logger.NewField("member_id", usage.MemberID), logger.NewField("error", err.Error()))
		return err
	}
	if member == nil {
		return ErrMemberNotFound
	}

	// 检查会员是否有效（状态为active且未过期）
	if !member.IsValid() {
		if member.Status == StatusExpired {
			return ErrMemberExpired
		}
		if member.Status == StatusInactive {
			return ErrMemberInactive
		}
		return ErrMemberExpired
	}

	// 设置使用日期（如果未指定，默认为当前日期）
	if usage.UsageDate.IsZero() {
		usage.UsageDate = time.Now()
	}

	// 设置时间戳
	usage.CreatedAt = time.Now()
	usage.UpdatedAt = time.Now()

	// 在事务中创建使用记录并递增会员已使用次数
	// 注意：这里需要数据库事务支持，实际实现时需要在仓储层使用事务
	if err := s.repo.Create(ctx, usage); err != nil {
		s.logger.Error("创建使用记录失败", logger.NewField("error", err.Error()))
		return err
	}

	// 递增会员已使用次数
	if err := s.memberRepo.IncrementUsedTimes(ctx, usage.MemberID); err != nil {
		s.logger.Error("递增会员已使用次数失败", logger.NewField("member_id", usage.MemberID), logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("创建使用记录成功", logger.NewField("usage_id", usage.ID))
	return nil
}

// GetUsageList 获取使用记录列表（不分页）
func (s *UsageService) GetUsageList(ctx context.Context, memberID *uint, storeID *uint, technicianID *uint) ([]*MemberUsage, error) {
	s.logger.Debug("获取使用记录列表",
		logger.NewField("member_id", memberID),
		logger.NewField("store_id", storeID),
		logger.NewField("technician_id", technicianID),
	)

	usages, err := s.repo.FindList(ctx, memberID, storeID, technicianID)
	if err != nil {
		s.logger.Error("获取使用记录列表失败", logger.NewField("error", err.Error()))
		return nil, err
	}

	return usages, nil
}

// GetUsageByID 获取使用记录详情
func (s *UsageService) GetUsageByID(ctx context.Context, id uint) (*MemberUsage, error) {
	s.logger.Debug("获取使用记录详情", logger.NewField("usage_id", id))

	usage, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找使用记录失败", logger.NewField("usage_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if usage == nil {
		s.logger.Warn("使用记录不存在", logger.NewField("usage_id", id))
		return nil, ErrUsageNotFound
	}

	return usage, nil
}

// DeleteUsage 删除使用记录（仅管理员），回退会员已使用次数
func (s *UsageService) DeleteUsage(ctx context.Context, id uint) error {
	s.logger.Info("删除使用记录", logger.NewField("usage_id", id))

	// 获取使用记录
	usage, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if usage == nil {
		return ErrUsageNotFound
	}

	// 删除使用记录
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("删除使用记录失败", logger.NewField("error", err.Error()))
		return err
	}

	// 回退会员已使用次数（递减1，但不能小于0）
	if err := s.memberRepo.DecrementUsedTimes(ctx, usage.MemberID); err != nil {
		s.logger.Error("回退会员已使用次数失败", logger.NewField("member_id", usage.MemberID), logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("删除使用记录成功", logger.NewField("usage_id", id))
	return nil
}

// GetUsageByMemberID 根据会员ID获取使用记录列表
func (s *UsageService) GetUsageByMemberID(ctx context.Context, memberID uint) ([]*MemberUsage, error) {
	s.logger.Debug("根据会员ID获取使用记录", logger.NewField("member_id", memberID))

	usages, err := s.repo.FindByMemberID(ctx, memberID)
	if err != nil {
		s.logger.Error("根据会员ID查找使用记录失败", logger.NewField("member_id", memberID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return usages, nil
}

