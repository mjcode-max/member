package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/store"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

var _ store.IStoreRepository = (*StoreRepository)(nil)

// 在包初始化时注册模型，确保迁移时能检测到
func init() {
	persistence.Register(&StoreModel{})
}

// StoreRepository 门店仓储实现
type StoreRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewStoreRepository 创建门店仓储实例
func NewStoreRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *StoreRepository {
	return &StoreRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// StoreModel 门店数据库模型
type StoreModel struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:100;not null" json:"name"`
	Address            string         `gorm:"size:255" json:"address"`
	Phone              string         `gorm:"size:20" json:"phone"`
	ContactPerson      string         `gorm:"size:50" json:"contact_person"`
	Status             string         `gorm:"size:20;default:'operating';not null" json:"status"`
	BusinessHoursStart string         `gorm:"size:10" json:"business_hours_start"`
	BusinessHoursEnd   string         `gorm:"size:10" json:"business_hours_end"`
	DepositAmount      float64        `gorm:"type:decimal(10,2);default:0" json:"deposit_amount"`
	TemplateID         *uint          `gorm:"index;default:NULL" json:"template_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreModel) TableName() string {
	return "stores"
}

// ToEntity 转换为领域实体
func (m *StoreModel) ToEntity() *store.Store {
	if m == nil {
		return nil
	}
	return &store.Store{
		ID:                 m.ID,
		Name:               m.Name,
		Address:            m.Address,
		Phone:              m.Phone,
		ContactPerson:      m.ContactPerson,
		Status:             m.Status,
		BusinessHoursStart: m.BusinessHoursStart,
		BusinessHoursEnd:   m.BusinessHoursEnd,
		DepositAmount:      m.DepositAmount,
		TemplateID:         m.TemplateID,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *StoreModel) FromEntity(s *store.Store) {
	if s == nil {
		return
	}
	m.ID = s.ID
	m.Name = s.Name
	m.Address = s.Address
	m.Phone = s.Phone
	m.ContactPerson = s.ContactPerson
	m.Status = s.Status
	m.BusinessHoursStart = s.BusinessHoursStart
	m.BusinessHoursEnd = s.BusinessHoursEnd
	m.DepositAmount = s.DepositAmount
	m.TemplateID = s.TemplateID
	m.CreatedAt = s.CreatedAt
	m.UpdatedAt = s.UpdatedAt
}

// FindByID 根据ID查找门店
func (r *StoreRepository) FindByID(ctx context.Context, id uint) (*store.Store, error) {
	r.logger.Debug("查找门店：根据ID",
		logger.NewField("store_id", id),
	)

	var model StoreModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("门店不存在：根据ID",
				logger.NewField("store_id", id),
			)
			return nil, nil
		}
		r.logger.Error("查找门店失败：根据ID",
			logger.NewField("store_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	r.logger.Debug("查找门店成功：根据ID",
		logger.NewField("store_id", id),
	)

	return model.ToEntity(), nil
}

// FindList 获取门店列表（支持筛选和分页）
func (r *StoreRepository) FindList(ctx context.Context, status, name string, page, pageSize int) ([]*store.Store, int64, error) {
	r.logger.Debug("查找门店列表",
		logger.NewField("status", status),
		logger.NewField("name", name),
		logger.NewField("page", page),
		logger.NewField("page_size", pageSize),
	)

	var models []StoreModel
	var total int64

	query := r.db.DB().WithContext(ctx).Model(&StoreModel{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("获取门店总数失败",
			logger.NewField("error", err.Error()),
		)
		return nil, 0, errors.ErrDatabase(err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error; err != nil {
		r.logger.Error("查找门店列表失败",
			logger.NewField("error", err.Error()),
		)
		return nil, 0, errors.ErrDatabase(err)
	}

	stores := make([]*store.Store, 0, len(models))
	for _, model := range models {
		stores = append(stores, model.ToEntity())
	}

	r.logger.Debug("查找门店列表成功",
		logger.NewField("count", len(stores)),
		logger.NewField("total", total),
	)

	return stores, total, nil
}

// Create 创建门店
func (r *StoreRepository) Create(ctx context.Context, s *store.Store) error {
	r.logger.Info("创建门店",
		logger.NewField("name", s.Name),
	)

	model := &StoreModel{}
	model.FromEntity(s)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建门店失败",
			logger.NewField("name", s.Name),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	s.ID = model.ID
	s.CreatedAt = model.CreatedAt
	s.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建门店成功",
		logger.NewField("store_id", s.ID),
		logger.NewField("name", s.Name),
	)

	return nil
}

// Update 更新门店
func (r *StoreRepository) Update(ctx context.Context, s *store.Store) error {
	r.logger.Info("更新门店",
		logger.NewField("store_id", s.ID),
	)

	model := &StoreModel{}
	model.FromEntity(s)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&StoreModel{}).Where("id = ?", s.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新门店失败",
			logger.NewField("store_id", s.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	s.UpdatedAt = model.UpdatedAt

	r.logger.Info("更新门店成功",
		logger.NewField("store_id", s.ID),
	)

	return nil
}

// Delete 删除门店（软删除）
func (r *StoreRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除门店",
		logger.NewField("store_id", id),
	)

	if err := r.db.DB().WithContext(ctx).Delete(&StoreModel{}, id).Error; err != nil {
		r.logger.Error("删除门店失败",
			logger.NewField("store_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	r.logger.Info("删除门店成功",
		logger.NewField("store_id", id),
	)

	return nil
}

// FindByManagerID 根据店长ID查找门店
func (r *StoreRepository) FindByManagerID(ctx context.Context, managerID uint) (*store.Store, error) {
	r.logger.Debug("查找门店：根据店长ID",
		logger.NewField("manager_id", managerID),
	)

	// 通过users表查找店长关联的门店
	var model StoreModel
	if err := r.db.DB().WithContext(ctx).
		Table("stores").
		Joins("INNER JOIN users ON users.store_id = stores.id").
		Where("users.id = ? AND users.role = ?", managerID, "store_manager").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("门店不存在：根据店长ID",
				logger.NewField("manager_id", managerID),
			)
			return nil, nil
		}
		r.logger.Error("查找门店失败：根据店长ID",
			logger.NewField("manager_id", managerID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	r.logger.Debug("查找门店成功：根据店长ID",
		logger.NewField("manager_id", managerID),
		logger.NewField("store_id", model.ID),
	)

	return model.ToEntity(), nil
}

// FindByStatus 根据状态查找门店列表
func (r *StoreRepository) FindByStatus(ctx context.Context, status string) ([]*store.Store, error) {
	r.logger.Debug("查找门店：根据状态",
		logger.NewField("status", status),
	)

	var models []StoreModel
	if err := r.db.DB().WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		r.logger.Error("查找门店失败：根据状态",
			logger.NewField("status", status),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	stores := make([]*store.Store, 0, len(models))
	for _, model := range models {
		stores = append(stores, model.ToEntity())
	}

	r.logger.Debug("查找门店成功：根据状态",
		logger.NewField("status", status),
		logger.NewField("count", len(stores)),
	)

	return stores, nil
}
