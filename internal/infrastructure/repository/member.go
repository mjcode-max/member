package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/member"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

var _ member.IMemberRepository = (*MemberRepository)(nil)
var _ member.IUsageRepository = (*UsageRepository)(nil)

// MemberRepository 会员仓储实现
type MemberRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewMemberRepository 创建会员仓储实例
func NewMemberRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *MemberRepository {
	persistence.Register(&MemberModel{})
	return &MemberRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// MemberModel 会员数据库模型（包含套餐信息）
type MemberModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Phone           string         `gorm:"size:20;index" json:"phone"`
	PackageName     string         `gorm:"size:100" json:"package_name"`
	ServiceType     string         `gorm:"size:20" json:"service_type"`
	Price           float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	UsedTimes       int            `gorm:"default:0" json:"used_times"`
	ValidityDuration int           `gorm:"default:0" json:"validity_duration"`
	ValidFrom       time.Time      `gorm:"type:date" json:"valid_from"`
	ValidTo         time.Time      `gorm:"type:date" json:"valid_to"`
	StoreID         uint           `gorm:"index" json:"store_id"`
	PurchaseAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"purchase_amount"`
	PurchaseTime    time.Time      `json:"purchase_time"`
	Status          string         `gorm:"size:20;default:'active';not null" json:"status"`
	Description     string         `gorm:"type:text" json:"description"`
	CreatedBy       uint           `gorm:"index" json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MemberModel) TableName() string {
	return "members"
}

// ToEntity 转换为领域实体
func (m *MemberModel) ToEntity() *member.Member {
	if m == nil {
		return nil
	}
	return &member.Member{
		ID:              m.ID,
		Name:            m.Name,
		Phone:           m.Phone,
		PackageName:     m.PackageName,
		ServiceType:     m.ServiceType,
		Price:           m.Price,
		UsedTimes:       m.UsedTimes,
		ValidityDuration: m.ValidityDuration,
		ValidFrom:       m.ValidFrom,
		ValidTo:         m.ValidTo,
		StoreID:         m.StoreID,
		PurchaseAmount:  m.PurchaseAmount,
		PurchaseTime:    m.PurchaseTime,
		Status:          m.Status,
		Description:     m.Description,
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *MemberModel) FromEntity(mem *member.Member) {
	if mem == nil {
		return
	}
	m.ID = mem.ID
	m.Name = mem.Name
	m.Phone = mem.Phone
	m.PackageName = mem.PackageName
	m.ServiceType = mem.ServiceType
	m.Price = mem.Price
	m.UsedTimes = mem.UsedTimes
	m.ValidityDuration = mem.ValidityDuration
	m.ValidFrom = mem.ValidFrom
	m.ValidTo = mem.ValidTo
	m.StoreID = mem.StoreID
	m.PurchaseAmount = mem.PurchaseAmount
	m.PurchaseTime = mem.PurchaseTime
	m.Status = mem.Status
	m.Description = mem.Description
	m.CreatedBy = mem.CreatedBy
	m.CreatedAt = mem.CreatedAt
	m.UpdatedAt = mem.UpdatedAt
}

// FindByID 根据ID查找会员
func (r *MemberRepository) FindByID(ctx context.Context, id uint) (*member.Member, error) {
	r.logger.Debug("查找会员：根据ID", logger.NewField("member_id", id))

	var model MemberModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("会员不存在：根据ID", logger.NewField("member_id", id))
			return nil, nil
		}
		r.logger.Error("查找会员失败：根据ID", logger.NewField("member_id", id), logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindList 获取会员列表（支持筛选和分页）
func (r *MemberRepository) FindList(ctx context.Context, name, phone string, storeID *uint, status, serviceType, packageName string, page, pageSize int) ([]*member.Member, int64, error) {
	r.logger.Debug("查找会员列表",
		logger.NewField("name", name),
		logger.NewField("phone", phone),
		logger.NewField("store_id", storeID),
		logger.NewField("status", status),
		logger.NewField("service_type", serviceType),
		logger.NewField("package_name", packageName),
		logger.NewField("page", page),
		logger.NewField("page_size", pageSize),
	)

	var models []MemberModel
	var total int64

	query := r.db.DB().WithContext(ctx).Model(&MemberModel{})

	// 筛选条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}
	if packageName != "" {
		query = query.Where("package_name LIKE ?", "%"+packageName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("获取会员总数失败", logger.NewField("error", err.Error()))
		return nil, 0, errors.ErrDatabase(err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error; err != nil {
		r.logger.Error("查找会员列表失败", logger.NewField("error", err.Error()))
		return nil, 0, errors.ErrDatabase(err)
	}

	members := make([]*member.Member, 0, len(models))
	for _, model := range models {
		members = append(members, model.ToEntity())
	}

	r.logger.Debug("查找会员列表成功",
		logger.NewField("count", len(members)),
		logger.NewField("total", total),
	)

	return members, total, nil
}

// Create 创建会员
func (r *MemberRepository) Create(ctx context.Context, m *member.Member) error {
	r.logger.Info("创建会员", logger.NewField("name", m.Name))

	model := &MemberModel{}
	model.FromEntity(m)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建会员失败", logger.NewField("name", m.Name), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	m.ID = model.ID
	m.CreatedAt = model.CreatedAt
	m.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建会员成功", logger.NewField("member_id", m.ID), logger.NewField("name", m.Name))
	return nil
}

// Update 更新会员
func (r *MemberRepository) Update(ctx context.Context, m *member.Member) error {
	r.logger.Info("更新会员", logger.NewField("member_id", m.ID))

	model := &MemberModel{}
	model.FromEntity(m)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&MemberModel{}).Where("id = ?", m.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新会员失败", logger.NewField("member_id", m.ID), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	m.UpdatedAt = model.UpdatedAt

	r.logger.Info("更新会员成功", logger.NewField("member_id", m.ID))
	return nil
}

// Delete 删除会员（软删除）
func (r *MemberRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除会员", logger.NewField("member_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&MemberModel{}, id).Error; err != nil {
		r.logger.Error("删除会员失败", logger.NewField("member_id", id), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	r.logger.Info("删除会员成功", logger.NewField("member_id", id))
	return nil
}

// FindByPhone 根据手机号查找会员列表
func (r *MemberRepository) FindByPhone(ctx context.Context, phone string) ([]*member.Member, error) {
	r.logger.Debug("查找会员：根据手机号", logger.NewField("phone", phone))

	var models []MemberModel
	if err := r.db.DB().WithContext(ctx).Where("phone = ?", phone).Find(&models).Error; err != nil {
		r.logger.Error("查找会员失败：根据手机号", logger.NewField("phone", phone), logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	members := make([]*member.Member, 0, len(models))
	for _, model := range models {
		members = append(members, model.ToEntity())
	}

	r.logger.Debug("查找会员成功：根据手机号", logger.NewField("phone", phone), logger.NewField("count", len(members)))
	return members, nil
}

// IncrementUsedTimes 递增会员已使用次数
func (r *MemberRepository) IncrementUsedTimes(ctx context.Context, id uint) error {
	r.logger.Debug("递增会员已使用次数", logger.NewField("member_id", id))

	if err := r.db.DB().WithContext(ctx).Model(&MemberModel{}).
		Where("id = ?", id).
		Update("used_times", gorm.Expr("used_times + 1")).Error; err != nil {
		r.logger.Error("递增会员已使用次数失败", logger.NewField("member_id", id), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	return nil
}

// DecrementUsedTimes 递减会员已使用次数（不能小于0）
func (r *MemberRepository) DecrementUsedTimes(ctx context.Context, id uint) error {
	r.logger.Debug("递减会员已使用次数", logger.NewField("member_id", id))

	if err := r.db.DB().WithContext(ctx).Model(&MemberModel{}).
		Where("id = ? AND used_times > 0", id).
		Update("used_times", gorm.Expr("used_times - 1")).Error; err != nil {
		r.logger.Error("递减会员已使用次数失败", logger.NewField("member_id", id), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	return nil
}

// UsageRepository 使用记录仓储实现
type UsageRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewUsageRepository 创建使用记录仓储实例
func NewUsageRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *UsageRepository {
	persistence.Register(&UsageModel{})
	return &UsageRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// UsageModel 使用记录数据库模型
type UsageModel struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	MemberID       uint      `gorm:"index" json:"member_id"`
	PackageName    string    `gorm:"size:100" json:"package_name"`
	ServiceItem    string    `gorm:"size:100" json:"service_item"`
	StoreID        uint      `gorm:"index" json:"store_id"`
	StoreName      string    `gorm:"size:100" json:"store_name"`
	TechnicianID   *uint     `gorm:"index" json:"technician_id"`
	TechnicianName string    `gorm:"size:50" json:"technician_name"`
	UsageDate      time.Time `gorm:"type:date" json:"usage_date"`
	Remark         string    `gorm:"type:text" json:"remark"`
	CreatedBy      uint      `gorm:"index" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UsageModel) TableName() string {
	return "member_usages"
}

// ToEntity 转换为领域实体
func (m *UsageModel) ToEntity() *member.MemberUsage {
	if m == nil {
		return nil
	}
	return &member.MemberUsage{
		ID:             m.ID,
		MemberID:       m.MemberID,
		PackageName:    m.PackageName,
		ServiceItem:    m.ServiceItem,
		StoreID:        m.StoreID,
		StoreName:      m.StoreName,
		TechnicianID:   m.TechnicianID,
		TechnicianName: m.TechnicianName,
		UsageDate:      m.UsageDate,
		Remark:         m.Remark,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *UsageModel) FromEntity(u *member.MemberUsage) {
	if u == nil {
		return
	}
	m.ID = u.ID
	m.MemberID = u.MemberID
	m.PackageName = u.PackageName
	m.ServiceItem = u.ServiceItem
	m.StoreID = u.StoreID
	m.StoreName = u.StoreName
	m.TechnicianID = u.TechnicianID
	m.TechnicianName = u.TechnicianName
	m.UsageDate = u.UsageDate
	m.Remark = u.Remark
	m.CreatedBy = u.CreatedBy
	m.CreatedAt = u.CreatedAt
	m.UpdatedAt = u.UpdatedAt
}

// FindByID 根据ID查找使用记录
func (r *UsageRepository) FindByID(ctx context.Context, id uint) (*member.MemberUsage, error) {
	r.logger.Debug("查找使用记录：根据ID", logger.NewField("usage_id", id))

	var model UsageModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("使用记录不存在：根据ID", logger.NewField("usage_id", id))
			return nil, nil
		}
		r.logger.Error("查找使用记录失败：根据ID", logger.NewField("usage_id", id), logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindList 获取使用记录列表（不分页，支持筛选）
func (r *UsageRepository) FindList(ctx context.Context, memberID *uint, storeID *uint, technicianID *uint) ([]*member.MemberUsage, error) {
	r.logger.Debug("查找使用记录列表",
		logger.NewField("member_id", memberID),
		logger.NewField("store_id", storeID),
		logger.NewField("technician_id", technicianID),
	)

	var models []UsageModel
	query := r.db.DB().WithContext(ctx).Model(&UsageModel{})

	// 筛选条件
	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if technicianID != nil {
		query = query.Where("technician_id = ?", *technicianID)
	}

	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		r.logger.Error("查找使用记录列表失败", logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	usages := make([]*member.MemberUsage, 0, len(models))
	for _, model := range models {
		usages = append(usages, model.ToEntity())
	}

	r.logger.Debug("查找使用记录列表成功", logger.NewField("count", len(usages)))
	return usages, nil
}

// Create 创建使用记录
func (r *UsageRepository) Create(ctx context.Context, u *member.MemberUsage) error {
	r.logger.Info("创建使用记录", logger.NewField("member_id", u.MemberID))

	model := &UsageModel{}
	model.FromEntity(u)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建使用记录失败", logger.NewField("member_id", u.MemberID), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建使用记录成功", logger.NewField("usage_id", u.ID), logger.NewField("member_id", u.MemberID))
	return nil
}

// Delete 删除使用记录
func (r *UsageRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除使用记录", logger.NewField("usage_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&UsageModel{}, id).Error; err != nil {
		r.logger.Error("删除使用记录失败", logger.NewField("usage_id", id), logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	r.logger.Info("删除使用记录成功", logger.NewField("usage_id", id))
	return nil
}

// FindByMemberID 根据会员ID查找使用记录列表
func (r *UsageRepository) FindByMemberID(ctx context.Context, memberID uint) ([]*member.MemberUsage, error) {
	r.logger.Debug("查找使用记录：根据会员ID", logger.NewField("member_id", memberID))

	var models []UsageModel
	if err := r.db.DB().WithContext(ctx).Where("member_id = ?", memberID).Order("created_at DESC").Find(&models).Error; err != nil {
		r.logger.Error("查找使用记录失败：根据会员ID", logger.NewField("member_id", memberID), logger.NewField("error", err.Error()))
		return nil, errors.ErrDatabase(err)
	}

	usages := make([]*member.MemberUsage, 0, len(models))
	for _, model := range models {
		usages = append(usages, model.ToEntity())
	}

	r.logger.Debug("查找使用记录成功：根据会员ID", logger.NewField("member_id", memberID), logger.NewField("count", len(usages)))
	return usages, nil
}

