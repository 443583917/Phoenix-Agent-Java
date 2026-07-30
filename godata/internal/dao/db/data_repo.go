package db

import (
	"context"
	"errors"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── Agent ────────────────────────────

type agentRepo struct{ db *gorm.DB }

func NewAgentRepository(db *gorm.DB) repository.AgentRepository {
	return &agentRepo{db}
}

func (r *agentRepo) FindByID(ctx context.Context, id string) (*model.Agent, error) {
	var entity model.Agent
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentRepo) FindBySn(ctx context.Context, sn string) (*model.Agent, error) {
	var entity model.Agent
	err := r.db.WithContext(ctx).Where("sn = ? AND del_flag = 0", sn).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentRepo) Page(ctx context.Context, page, size int, query *model.Agent) ([]*model.Agent, int64, error) {
	var list []*model.Agent
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.Agent{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Type != "" {
			dbQuery = dbQuery.Where("type = ?", query.Type)
		}
		if query.Status != "" {
			dbQuery = dbQuery.Where("status = ?", query.Status)
		}
		if query.CategoryId != "" {
			dbQuery = dbQuery.Where("category_id = ?", query.CategoryId)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("order_num ASC, create_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *agentRepo) List(ctx context.Context) ([]*model.Agent, error) {
	var list []*model.Agent
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("order_num ASC, create_time DESC").Find(&list).Error
	return list, err
}

func (r *agentRepo) Create(ctx context.Context, entity *model.Agent) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *agentRepo) Update(ctx context.Context, entity *model.Agent) error {
	return r.db.WithContext(ctx).Model(&model.Agent{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *agentRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Agent{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AgentCategory ────────────────────────────

type agentCategoryRepo struct{ db *gorm.DB }

func NewAgentCategoryRepository(db *gorm.DB) repository.AgentCategoryRepository {
	return &agentCategoryRepo{db}
}

func (r *agentCategoryRepo) FindByID(ctx context.Context, id string) (*model.AgentCategory, error) {
	var entity model.AgentCategory
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentCategoryRepo) FindByPID(ctx context.Context, pid string) ([]*model.AgentCategory, error) {
	var list []*model.AgentCategory
	err := r.db.WithContext(ctx).Where("pid = ? AND del_flag = 0", pid).Order("sn ASC").Find(&list).Error
	return list, err
}

func (r *agentCategoryRepo) Page(ctx context.Context, page, size int, query *model.AgentCategory) ([]*model.AgentCategory, int64, error) {
	var list []*model.AgentCategory
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AgentCategory{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("sn ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *agentCategoryRepo) List(ctx context.Context) ([]*model.AgentCategory, error) {
	var list []*model.AgentCategory
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("sn ASC").Find(&list).Error
	return list, err
}

func (r *agentCategoryRepo) Create(ctx context.Context, entity *model.AgentCategory) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *agentCategoryRepo) Update(ctx context.Context, entity *model.AgentCategory) error {
	return r.db.WithContext(ctx).Model(&model.AgentCategory{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *agentCategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AgentCategory{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AgentDatasource ────────────────────────────

type agentDatasourceRepo struct{ db *gorm.DB }

func NewAgentDatasourceRepository(db *gorm.DB) repository.AgentDatasourceRepository {
	return &agentDatasourceRepo{db}
}

func (r *agentDatasourceRepo) FindByID(ctx context.Context, id string) (*model.AgentDatasource, error) {
	var entity model.AgentDatasource
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentDatasourceRepo) FindByAgentID(ctx context.Context, agentID int64) ([]*model.AgentDatasource, error) {
	var list []*model.AgentDatasource
	err := r.db.WithContext(ctx).Where("agent_id = ? AND del_flag = 0", agentID).Find(&list).Error
	return list, err
}

func (r *agentDatasourceRepo) Page(ctx context.Context, page, size int, query *model.AgentDatasource) ([]*model.AgentDatasource, int64, error) {
	var list []*model.AgentDatasource
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AgentDatasource{}).Where("del_flag = 0")
	if query != nil {
		if query.AgentId != 0 {
			dbQuery = dbQuery.Where("agent_id = ?", query.AgentId)
		}
		if query.DatasourceId != 0 {
			dbQuery = dbQuery.Where("datasource_id = ?", query.DatasourceId)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *agentDatasourceRepo) Create(ctx context.Context, entity *model.AgentDatasource) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *agentDatasourceRepo) Update(ctx context.Context, entity *model.AgentDatasource) error {
	return r.db.WithContext(ctx).Model(&model.AgentDatasource{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *agentDatasourceRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AgentDatasource{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AgentKnowledge ────────────────────────────

type agentKnowledgeRepo struct{ db *gorm.DB }

func NewAgentKnowledgeRepository(db *gorm.DB) repository.AgentKnowledgeRepository {
	return &agentKnowledgeRepo{db}
}

func (r *agentKnowledgeRepo) FindByID(ctx context.Context, id string) (*model.AgentKnowledge, error) {
	var entity model.AgentKnowledge
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentKnowledgeRepo) FindByAgentID(ctx context.Context, agentID int) ([]*model.AgentKnowledge, error) {
	var list []*model.AgentKnowledge
	err := r.db.WithContext(ctx).Where("agent_id = ? AND del_flag = 0", agentID).Find(&list).Error
	return list, err
}

func (r *agentKnowledgeRepo) Page(ctx context.Context, page, size int, query *model.AgentKnowledge) ([]*model.AgentKnowledge, int64, error) {
	var list []*model.AgentKnowledge
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AgentKnowledge{}).Where("del_flag = 0")
	if query != nil {
		if query.AgentId != 0 {
			dbQuery = dbQuery.Where("agent_id = ?", query.AgentId)
		}
		if query.Title != "" {
			dbQuery = dbQuery.Where("title LIKE ?", "%"+query.Title+"%")
		}
		if query.Type != "" {
			dbQuery = dbQuery.Where("type = ?", query.Type)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("create_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *agentKnowledgeRepo) Create(ctx context.Context, entity *model.AgentKnowledge) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *agentKnowledgeRepo) Update(ctx context.Context, entity *model.AgentKnowledge) error {
	return r.db.WithContext(ctx).Model(&model.AgentKnowledge{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *agentKnowledgeRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AgentKnowledge{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AgentPresetQuestion ────────────────────────────

type agentPresetQuestionRepo struct{ db *gorm.DB }

func NewAgentPresetQuestionRepository(db *gorm.DB) repository.AgentPresetQuestionRepository {
	return &agentPresetQuestionRepo{db}
}

func (r *agentPresetQuestionRepo) FindByID(ctx context.Context, id string) (*model.AgentPresetQuestion, error) {
	var entity model.AgentPresetQuestion
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *agentPresetQuestionRepo) FindByAgentID(ctx context.Context, agentID int64) ([]*model.AgentPresetQuestion, error) {
	var list []*model.AgentPresetQuestion
	err := r.db.WithContext(ctx).Where("agent_id = ? AND del_flag = 0", agentID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *agentPresetQuestionRepo) Page(ctx context.Context, page, size int, query *model.AgentPresetQuestion) ([]*model.AgentPresetQuestion, int64, error) {
	var list []*model.AgentPresetQuestion
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.AgentPresetQuestion{}).Where("del_flag = 0")
	if query != nil {
		if query.AgentId != 0 {
			dbQuery = dbQuery.Where("agent_id = ?", query.AgentId)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("sort_order ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *agentPresetQuestionRepo) Create(ctx context.Context, entity *model.AgentPresetQuestion) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *agentPresetQuestionRepo) Update(ctx context.Context, entity *model.AgentPresetQuestion) error {
	return r.db.WithContext(ctx).Model(&model.AgentPresetQuestion{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *agentPresetQuestionRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.AgentPresetQuestion{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── AgentDatasourceTables ────────────────────────────

type agentDatasourceTablesRepo struct{ db *gorm.DB }

func NewAgentDatasourceTablesRepository(db *gorm.DB) repository.AgentDatasourceTablesRepository {
	return &agentDatasourceTablesRepo{db}
}

func (r *agentDatasourceTablesRepo) FindByAgentDatasourceID(ctx context.Context, agentDatasourceID int) ([]*model.AgentDatasourceTables, error) {
	var list []*model.AgentDatasourceTables
	err := r.db.WithContext(ctx).Where("agent_datasource_id = ? AND del_flag = 0", agentDatasourceID).Find(&list).Error
	return list, err
}

func (r *agentDatasourceTablesRepo) SaveBatch(ctx context.Context, agentDatasourceID int, tables []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft-delete existing.
		if err := tx.Model(&model.AgentDatasourceTables{}).
			Where("agent_datasource_id = ? AND del_flag = 0", agentDatasourceID).
			Update("del_flag", 1).Error; err != nil {
			return err
		}
		// Insert new rows.
		for _, table := range tables {
			adt := &model.AgentDatasourceTables{
				AgentDatasourceId: agentDatasourceID,
				Table:             table,
			}
			if err := tx.Create(adt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *agentDatasourceTablesRepo) DeleteByAgentDatasourceID(ctx context.Context, agentDatasourceID int) error {
	return r.db.WithContext(ctx).Model(&model.AgentDatasourceTables{}).
		Where("agent_datasource_id = ? AND del_flag = 0", agentDatasourceID).
		Update("del_flag", 1).Error
}
