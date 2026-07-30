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

// ──────────────────────────── ChatSession ────────────────────────────

type chatSessionRepo struct{ db *gorm.DB }

func NewChatSessionRepository(db *gorm.DB) repository.ChatSessionRepository {
	return &chatSessionRepo{db}
}

func (r *chatSessionRepo) FindByID(ctx context.Context, id string) (*model.ChatSession, error) {
	var entity model.ChatSession
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *chatSessionRepo) FindByAgentIDAndUserID(ctx context.Context, agentID int, userID string) ([]*model.ChatSession, error) {
	var list []*model.ChatSession
	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND user_id = ? AND del_flag = 0", agentID, userID).
		Order("is_pinned DESC, update_time DESC").
		Find(&list).Error
	return list, err
}

func (r *chatSessionRepo) FindBySessionID(ctx context.Context, sessionID string) (*model.ChatSession, error) {
	return r.FindByID(ctx, sessionID)
}

func (r *chatSessionRepo) Create(ctx context.Context, entity *model.ChatSession) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *chatSessionRepo) Update(ctx context.Context, entity *model.ChatSession) error {
	return r.db.WithContext(ctx).Model(&model.ChatSession{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *chatSessionRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.ChatSession{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *chatSessionRepo) DeleteByAgentID(ctx context.Context, agentID int) error {
	return r.db.WithContext(ctx).Model(&model.ChatSession{}).Where("agent_id = ? AND del_flag = 0", agentID).Update("del_flag", 1).Error
}

// ──────────────────────────── ChatMessage ────────────────────────────

type chatMessageRepo struct{ db *gorm.DB }

func NewChatMessageRepository(db *gorm.DB) repository.ChatMessageRepository {
	return &chatMessageRepo{db}
}

func (r *chatMessageRepo) FindBySessionID(ctx context.Context, sessionID string) ([]*model.ChatMessage, error) {
	var list []*model.ChatMessage
	err := r.db.WithContext(ctx).
		Where("session_id = ? AND del_flag = 0", sessionID).
		Order("create_time ASC").
		Find(&list).Error
	return list, err
}

func (r *chatMessageRepo) Create(ctx context.Context, entity *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// ──────────────────────────── Datasource ────────────────────────────

type datasourceRepo struct{ db *gorm.DB }

func NewDatasourceRepository(db *gorm.DB) repository.DatasourceRepository {
	return &datasourceRepo{db}
}

func (r *datasourceRepo) FindByID(ctx context.Context, id string) (*model.Datasource, error) {
	var entity model.Datasource
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *datasourceRepo) Page(ctx context.Context, page, size int, query *model.Datasource) ([]*model.Datasource, int64, error) {
	var list []*model.Datasource
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.Datasource{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
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

func (r *datasourceRepo) Create(ctx context.Context, entity *model.Datasource) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *datasourceRepo) Update(ctx context.Context, entity *model.Datasource) error {
	return r.db.WithContext(ctx).Model(&model.Datasource{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *datasourceRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Datasource{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── LogicalRelation ────────────────────────────

type logicalRelationRepo struct{ db *gorm.DB }

func NewLogicalRelationRepository(db *gorm.DB) repository.LogicalRelationRepository {
	return &logicalRelationRepo{db}
}

func (r *logicalRelationRepo) FindByDatasourceID(ctx context.Context, datasourceID int) ([]*model.LogicalRelation, error) {
	var list []*model.LogicalRelation
	err := r.db.WithContext(ctx).Where("datasource_id = ? AND del_flag = 0", datasourceID).Find(&list).Error
	return list, err
}

func (r *logicalRelationRepo) Create(ctx context.Context, entity *model.LogicalRelation) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *logicalRelationRepo) Update(ctx context.Context, entity *model.LogicalRelation) error {
	return r.db.WithContext(ctx).Model(&model.LogicalRelation{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *logicalRelationRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.LogicalRelation{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *logicalRelationRepo) FindByID(ctx context.Context, id string) (*model.LogicalRelation, error) {
	var entity model.LogicalRelation
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *logicalRelationRepo) DeleteByDatasourceID(ctx context.Context, datasourceID int) error {
	return r.db.WithContext(ctx).Model(&model.LogicalRelation{}).Where("datasource_id = ? AND del_flag = 0", datasourceID).Update("del_flag", 1).Error
}

// ──────────────────────────── ModelConfig ────────────────────────────

type modelConfigRepo struct{ db *gorm.DB }

func NewModelConfigRepository(db *gorm.DB) repository.ModelConfigRepository {
	return &modelConfigRepo{db}
}

func (r *modelConfigRepo) FindByID(ctx context.Context, id string) (*model.ModelConfig, error) {
	var entity model.ModelConfig
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *modelConfigRepo) Page(ctx context.Context, page, size int) ([]*model.ModelConfig, int64, error) {
	var list []*model.ModelConfig
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.ModelConfig{}).Where("del_flag = 0")
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

func (r *modelConfigRepo) FindActive(ctx context.Context) (*model.ModelConfig, error) {
	var entity model.ModelConfig
	err := r.db.WithContext(ctx).Where("is_active = true AND del_flag = 0").First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *modelConfigRepo) Create(ctx context.Context, entity *model.ModelConfig) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *modelConfigRepo) Update(ctx context.Context, entity *model.ModelConfig) error {
	return r.db.WithContext(ctx).Model(&model.ModelConfig{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *modelConfigRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.ModelConfig{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *modelConfigRepo) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&model.ModelConfig{}).Where("del_flag = 0").Update("is_active", false).Error
}

// ──────────────────────────── UserPromptConfig ────────────────────────────

type userPromptConfigRepo struct{ db *gorm.DB }

func NewUserPromptConfigRepository(db *gorm.DB) repository.UserPromptConfigRepository {
	return &userPromptConfigRepo{db}
}

func (r *userPromptConfigRepo) FindByID(ctx context.Context, id string) (*model.UserPromptConfig, error) {
	var entity model.UserPromptConfig
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *userPromptConfigRepo) Page(ctx context.Context, page, size int, promptType string) ([]*model.UserPromptConfig, int64, error) {
	var list []*model.UserPromptConfig
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.UserPromptConfig{}).Where("del_flag = 0")
	if promptType != "" {
		dbQuery = dbQuery.Where("prompt_type = ?", promptType)
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
	if err := dbQuery.Offset(offset).Limit(size).Order("display_order ASC, priority DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *userPromptConfigRepo) FindByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error) {
	var list []*model.UserPromptConfig
	err := r.db.WithContext(ctx).Where("prompt_type = ? AND del_flag = 0", promptType).Order("display_order ASC").Find(&list).Error
	return list, err
}

func (r *userPromptConfigRepo) FindActiveByType(ctx context.Context, promptType string) (*model.UserPromptConfig, error) {
	var entity model.UserPromptConfig
	err := r.db.WithContext(ctx).Where("prompt_type = ? AND enabled = true AND del_flag = 0", promptType).Order("priority DESC").First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *userPromptConfigRepo) FindActiveAllByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error) {
	var list []*model.UserPromptConfig
	err := r.db.WithContext(ctx).Where("prompt_type = ? AND enabled = true AND del_flag = 0", promptType).Order("priority DESC").Find(&list).Error
	return list, err
}

func (r *userPromptConfigRepo) Create(ctx context.Context, entity *model.UserPromptConfig) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *userPromptConfigRepo) Update(ctx context.Context, entity *model.UserPromptConfig) error {
	return r.db.WithContext(ctx).Model(&model.UserPromptConfig{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *userPromptConfigRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.UserPromptConfig{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *userPromptConfigRepo) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.UserPromptConfig{}).Where("id IN ? AND del_flag = 0", ids).Updates(updates).Error
}

// ──────────────────────────── SemanticModel ────────────────────────────

type semanticModelRepo struct{ db *gorm.DB }

func NewSemanticModelRepository(db *gorm.DB) repository.SemanticModelRepository {
	return &semanticModelRepo{db}
}

func (r *semanticModelRepo) FindByID(ctx context.Context, id string) (*model.SemanticModel, error) {
	var entity model.SemanticModel
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *semanticModelRepo) Page(ctx context.Context, page, size int, query *model.SemanticModel) ([]*model.SemanticModel, int64, error) {
	var list []*model.SemanticModel
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.SemanticModel{}).Where("del_flag = 0")
	if query != nil {
		if query.AgentId != 0 {
			dbQuery = dbQuery.Where("agent_id = ?", query.AgentId)
		}
		if query.Table != "" {
			dbQuery = dbQuery.Where("table_name LIKE ?", "%"+query.Table+"%")
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

func (r *semanticModelRepo) Create(ctx context.Context, entity *model.SemanticModel) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *semanticModelRepo) Update(ctx context.Context, entity *model.SemanticModel) error {
	return r.db.WithContext(ctx).Model(&model.SemanticModel{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *semanticModelRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.SemanticModel{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *semanticModelRepo) BatchDelete(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Model(&model.SemanticModel{}).Where("id IN ? AND del_flag = 0", ids).Update("del_flag", 1).Error
}

func (r *semanticModelRepo) BatchUpdateStatus(ctx context.Context, ids []string, status int) error {
	return r.db.WithContext(ctx).Model(&model.SemanticModel{}).Where("id IN ? AND del_flag = 0", ids).Update("status", status).Error
}

func (r *semanticModelRepo) BatchCreate(ctx context.Context, sms []*model.SemanticModel) error {
	return r.db.WithContext(ctx).CreateInBatches(sms, 100).Error
}

// ──────────────────────────── BusinessKnowledge ────────────────────────────

type businessKnowledgeRepo struct{ db *gorm.DB }

func NewBusinessKnowledgeRepository(db *gorm.DB) repository.BusinessKnowledgeRepository {
	return &businessKnowledgeRepo{db}
}

func (r *businessKnowledgeRepo) FindByID(ctx context.Context, id string) (*model.BusinessKnowledge, error) {
	var entity model.BusinessKnowledge
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *businessKnowledgeRepo) Page(ctx context.Context, page, size int, query *model.BusinessKnowledge) ([]*model.BusinessKnowledge, int64, error) {
	var list []*model.BusinessKnowledge
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.BusinessKnowledge{}).Where("del_flag = 0")
	if query != nil {
		if query.AgentId != 0 {
			dbQuery = dbQuery.Where("agent_id = ?", query.AgentId)
		}
		if query.BusinessTerm != "" {
			dbQuery = dbQuery.Where("business_term LIKE ?", "%"+query.BusinessTerm+"%")
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

func (r *businessKnowledgeRepo) Create(ctx context.Context, entity *model.BusinessKnowledge) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *businessKnowledgeRepo) Update(ctx context.Context, entity *model.BusinessKnowledge) error {
	return r.db.WithContext(ctx).Model(&model.BusinessKnowledge{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *businessKnowledgeRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.BusinessKnowledge{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

func (r *businessKnowledgeRepo) UpdateRecall(ctx context.Context, id string, isRecall int) error {
	return r.db.WithContext(ctx).Model(&model.BusinessKnowledge{}).Where("id = ? AND del_flag = 0", id).Update("is_recall", isRecall).Error
}

func (r *businessKnowledgeRepo) FindByAgentID(ctx context.Context, agentID int64) ([]*model.BusinessKnowledge, error) {
	var list []*model.BusinessKnowledge
	err := r.db.WithContext(ctx).Where("agent_id = ? AND del_flag = 0", agentID).Find(&list).Error
	return list, err
}

func (r *businessKnowledgeRepo) BatchResetEmbedding(ctx context.Context, agentID int64) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.BusinessKnowledge{}).
		Where("agent_id = ? AND del_flag = 0", agentID).
		Updates(map[string]interface{}{"embedding_status": "pending", "error_msg": ""})
	return result.RowsAffected, result.Error
}
