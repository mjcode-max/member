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

// TimeSlotRule 时段规则（时间段）
type TimeSlotRule struct {
	StartTime string `json:"start_time"` // 开始时间 (HH:MM格式)
	EndTime   string `json:"end_time"`   // 结束时间 (HH:MM格式)
}

// Template 时段模板实体
type Template struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`       // 模板名称
	Status    string         `json:"status"`     // 状态: active, inactive
	TimeSlots []TimeSlotRule `json:"time_slots"` // 时间段列表（每天使用相同的时段）
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ITemplateRepository 时段模板仓储接口
type ITemplateRepository interface {
	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, id uint) (*Template, error)
	// FindAll 查找所有模板
	FindAll(ctx context.Context) ([]*Template, error)
	// FindByStatus 根据状态查找模板列表
	FindByStatus(ctx context.Context, status string) ([]*Template, error)
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
	ErrInvalidTimeSlot     = errors.ErrInvalidParams("无效的时段规则")
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

// GetAll 获取所有模板列表
func (s *TemplateService) GetAll(ctx context.Context) ([]*Template, error) {
	s.logger.Debug("获取所有时段模板列表")

	templates, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("查找所有时段模板列表失败", logger.NewField("error", err.Error()))
		return nil, err
	}

	return templates, nil
}

// GetByStatus 根据状态获取模板列表
func (s *TemplateService) GetByStatus(ctx context.Context, status string) ([]*Template, error) {
	s.logger.Debug("获取指定状态的时段模板列表", logger.NewField("status", status))

	templates, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		s.logger.Error("查找指定状态的时段模板列表失败", logger.NewField("status", status), logger.NewField("error", err.Error()))
		return nil, err
	}

	return templates, nil
}

// Create 创建模板
func (s *TemplateService) Create(ctx context.Context, template *Template) error {
	s.logger.Info("创建时段模板",
		logger.NewField("name", template.Name),
	)

	// 验证模板名称
	if template.Name == "" {
		return ErrInvalidTemplateName
	}

	// 验证时间段规则
	if err := s.validateTimeSlots(template.TimeSlots); err != nil {
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

	// 验证时间段规则（如果提供了）
	if len(template.TimeSlots) > 0 {
		if err := s.validateTimeSlots(template.TimeSlots); err != nil {
			return err
		}
	} else {
		// 保留原有规则
		template.TimeSlots = existing.TimeSlots
	}

	// 保留原有字段（如果新字段为空）
	if template.Name == "" {
		template.Name = existing.Name
	}
	if template.Status == "" {
		template.Status = existing.Status
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

// validateTimeSlots 验证时间段规则
func (s *TemplateService) validateTimeSlots(slots []TimeSlotRule) error {
	if len(slots) == 0 {
		return errors.ErrInvalidParams("至少需要设置一个时间段")
	}

	// 检查时间段是否重叠
	for i, slot1 := range slots {
		// 验证时间格式 (HH:MM)
		if !isValidTimeFormat(slot1.StartTime) || !isValidTimeFormat(slot1.EndTime) {
			return ErrInvalidTimeSlot
		}

		// 验证开始时间早于结束时间
		if slot1.StartTime >= slot1.EndTime {
			return errors.ErrInvalidParams("开始时间必须早于结束时间")
		}

		// 检查是否与其他时间段重叠
		for j, slot2 := range slots {
			if i >= j {
				continue
			}
			// 检查时间段是否重叠
			if slot1.StartTime < slot2.EndTime && slot1.EndTime > slot2.StartTime {
				return errors.ErrInvalidParams("时间段不能重叠")
			}
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
