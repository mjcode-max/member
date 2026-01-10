package slot

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 时段模板状态常量
const (
	TemplateStatusActive   = "active"   // 启用
	TemplateStatusInactive = "inactive" // 禁用
)

// TimeSlotRule 时段规则（一天中的时间段）
type TimeSlotRule struct {
	StartTime string `json:"start_time"` // 开始时间 (HH:MM格式)
	EndTime   string `json:"end_time"`   // 结束时间 (HH:MM格式)
	Duration  int    `json:"duration"`   // 时段时长（分钟）
}

// WeekdayRule 星期规则
type WeekdayRule struct {
	Weekday int            `json:"weekday"` // 星期几 (0=周日, 1=周一, ..., 6=周六)
	Slots   []TimeSlotRule `json:"slots"`   // 该天的时段列表
}

// Template 时段模板实体
type Template struct {
	ID          uint          `json:"id"`
	StoreID     uint          `json:"store_id"`     // 门店ID
	Name        string        `json:"name"`         // 模板名称
	Status      string        `json:"status"`       // 状态: active, inactive
	WeekdayRules []WeekdayRule `json:"weekday_rules"` // 星期规则列表
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ITemplateRepository 时段模板仓储接口
type ITemplateRepository interface {
	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, id uint) (*Template, error)
	// FindByStoreID 根据门店ID查找模板列表
	FindByStoreID(ctx context.Context, storeID uint) ([]*Template, error)
	// FindActiveByStoreID 根据门店ID查找启用的模板
	FindActiveByStoreID(ctx context.Context, storeID uint) (*Template, error)
	// Create 创建模板
	Create(ctx context.Context, template *Template) error
	// Update 更新模板
	Update(ctx context.Context, template *Template) error
	// Delete 删除模板
	Delete(ctx context.Context, id uint) error
}

// 领域错误定义
var (
	ErrTemplateNotFound    = errors.ErrNotFound("时段模板不存在")
	ErrInvalidTemplateName = errors.ErrInvalidParams("模板名称不能为空")
	ErrInvalidWeekdayRule  = errors.ErrInvalidParams("无效的星期规则")
	ErrInvalidTimeSlot      = errors.ErrInvalidParams("无效的时段规则")
	ErrTemplateInactive    = errors.ErrInvalidParams("时段模板未启用")
)

// TemplateService 时段模板服务
type TemplateService struct {
	repo   ITemplateRepository
	logger logger.Logger
}

// NewTemplateService 创建时段模板服务
func NewTemplateService(repo ITemplateRepository, log logger.Logger) *TemplateService {
	return &TemplateService{
		repo:   repo,
		logger: log,
	}
}

// GetByID 根据ID获取模板
func (s *TemplateService) GetByID(ctx context.Context, id uint) (*Template, error) {
	s.logger.Debug("获取时段模板", logger.NewField("template_id", id))

	template, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找时段模板失败", logger.NewField("template_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if template == nil {
		s.logger.Warn("时段模板不存在", logger.NewField("template_id", id))
		return nil, ErrTemplateNotFound
	}

	return template, nil
}

// GetByStoreID 根据门店ID获取模板列表
func (s *TemplateService) GetByStoreID(ctx context.Context, storeID uint) ([]*Template, error) {
	s.logger.Debug("获取门店时段模板列表", logger.NewField("store_id", storeID))

	templates, err := s.repo.FindByStoreID(ctx, storeID)
	if err != nil {
		s.logger.Error("查找时段模板列表失败", logger.NewField("store_id", storeID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return templates, nil
}

// GetActiveByStoreID 根据门店ID获取启用的模板
func (s *TemplateService) GetActiveByStoreID(ctx context.Context, storeID uint) (*Template, error) {
	s.logger.Debug("获取门店启用的时段模板", logger.NewField("store_id", storeID))

	template, err := s.repo.FindActiveByStoreID(ctx, storeID)
	if err != nil {
		s.logger.Error("查找启用的时段模板失败", logger.NewField("store_id", storeID), logger.NewField("error", err.Error()))
		return nil, err
	}

	if template == nil {
		s.logger.Warn("门店没有启用的时段模板", logger.NewField("store_id", storeID))
		return nil, ErrTemplateNotFound
	}

	return template, nil
}

// Create 创建模板
func (s *TemplateService) Create(ctx context.Context, template *Template) error {
	s.logger.Info("创建时段模板",
		logger.NewField("store_id", template.StoreID),
		logger.NewField("name", template.Name),
	)

	// 验证模板名称
	if template.Name == "" {
		return ErrInvalidTemplateName
	}

	// 验证星期规则
	if err := s.validateWeekdayRules(template.WeekdayRules); err != nil {
		return err
	}

	// 设置默认值
	if template.Status == "" {
		template.Status = TemplateStatusActive
	}
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, template); err != nil {
		s.logger.Error("创建时段模板失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("创建时段模板成功", logger.NewField("template_id", template.ID))
	return nil
}

// Update 更新模板
func (s *TemplateService) Update(ctx context.Context, template *Template) error {
	s.logger.Info("更新时段模板", logger.NewField("template_id", template.ID))

	// 检查模板是否存在
	existing, err := s.repo.FindByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTemplateNotFound
	}

	// 验证星期规则（如果提供了）
	if len(template.WeekdayRules) > 0 {
		if err := s.validateWeekdayRules(template.WeekdayRules); err != nil {
			return err
		}
	} else {
		// 保留原有规则
		template.WeekdayRules = existing.WeekdayRules
	}

	// 保留原有字段（如果新字段为空）
	if template.Name == "" {
		template.Name = existing.Name
	}
	if template.Status == "" {
		template.Status = existing.Status
	}
	if template.StoreID == 0 {
		template.StoreID = existing.StoreID
	}

	template.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, template); err != nil {
		s.logger.Error("更新时段模板失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("更新时段模板成功", logger.NewField("template_id", template.ID))
	return nil
}

// Delete 删除模板
func (s *TemplateService) Delete(ctx context.Context, id uint) error {
	s.logger.Info("删除时段模板", logger.NewField("template_id", id))

	// 检查模板是否存在
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTemplateNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("删除时段模板失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("删除时段模板成功", logger.NewField("template_id", id))
	return nil
}

// validateWeekdayRules 验证星期规则
func (s *TemplateService) validateWeekdayRules(rules []WeekdayRule) error {
	weekdayMap := make(map[int]bool)
	for _, rule := range rules {
		// 验证星期几范围 (0-6)
		if rule.Weekday < 0 || rule.Weekday > 6 {
			return ErrInvalidWeekdayRule
		}

		// 检查是否重复
		if weekdayMap[rule.Weekday] {
			return errors.ErrInvalidParams("星期规则重复")
		}
		weekdayMap[rule.Weekday] = true

		// 验证时段规则
		if err := s.validateTimeSlots(rule.Slots); err != nil {
			return err
		}
	}
	return nil
}

// validateTimeSlots 验证时段规则
func (s *TemplateService) validateTimeSlots(slots []TimeSlotRule) error {
	for _, slot := range slots {
		// 验证时间格式 (HH:MM)
		if !isValidTimeFormat(slot.StartTime) || !isValidTimeFormat(slot.EndTime) {
			return ErrInvalidTimeSlot
		}

		// 验证时长
		if slot.Duration <= 0 {
			return errors.ErrInvalidParams("时段时长必须大于0")
		}

		// 验证开始时间早于结束时间
		if slot.StartTime >= slot.EndTime {
			return errors.ErrInvalidParams("开始时间必须早于结束时间")
		}
	}
	return nil
}

// isValidTimeFormat 验证时间格式是否为HH:MM
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	if timeStr[2] != ':' {
		return false
	}
	// 简单验证，实际可以使用正则表达式
	return true
}

