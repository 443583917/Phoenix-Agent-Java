package repository

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
)

// ──────────────────────────── Agent ────────────────────────────

type AgentRepository interface {
	FindByID(ctx context.Context, id string) (*model.Agent, error)
	FindBySn(ctx context.Context, sn string) (*model.Agent, error)
	Page(ctx context.Context, page, size int, query *model.Agent) ([]*model.Agent, int64, error)
	List(ctx context.Context) ([]*model.Agent, error)
	Create(ctx context.Context, agent *model.Agent) error
	Update(ctx context.Context, agent *model.Agent) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AgentCategory ────────────────────────────

type AgentCategoryRepository interface {
	FindByID(ctx context.Context, id string) (*model.AgentCategory, error)
	FindByPID(ctx context.Context, pid string) ([]*model.AgentCategory, error)
	Page(ctx context.Context, page, size int, query *model.AgentCategory) ([]*model.AgentCategory, int64, error)
	List(ctx context.Context) ([]*model.AgentCategory, error)
	Create(ctx context.Context, category *model.AgentCategory) error
	Update(ctx context.Context, category *model.AgentCategory) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AgentDatasource ────────────────────────────

type AgentDatasourceRepository interface {
	FindByID(ctx context.Context, id string) (*model.AgentDatasource, error)
	FindByAgentID(ctx context.Context, agentID int64) ([]*model.AgentDatasource, error)
	Page(ctx context.Context, page, size int, query *model.AgentDatasource) ([]*model.AgentDatasource, int64, error)
	Create(ctx context.Context, ad *model.AgentDatasource) error
	Update(ctx context.Context, ad *model.AgentDatasource) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AgentKnowledge ────────────────────────────

type AgentKnowledgeRepository interface {
	FindByID(ctx context.Context, id string) (*model.AgentKnowledge, error)
	FindByAgentID(ctx context.Context, agentID int) ([]*model.AgentKnowledge, error)
	Page(ctx context.Context, page, size int, query *model.AgentKnowledge) ([]*model.AgentKnowledge, int64, error)
	Create(ctx context.Context, knowledge *model.AgentKnowledge) error
	Update(ctx context.Context, knowledge *model.AgentKnowledge) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AgentPresetQuestion ────────────────────────────

type AgentPresetQuestionRepository interface {
	FindByID(ctx context.Context, id string) (*model.AgentPresetQuestion, error)
	FindByAgentID(ctx context.Context, agentID int64) ([]*model.AgentPresetQuestion, error)
	Page(ctx context.Context, page, size int, query *model.AgentPresetQuestion) ([]*model.AgentPresetQuestion, int64, error)
	Create(ctx context.Context, pq *model.AgentPresetQuestion) error
	Update(ctx context.Context, pq *model.AgentPresetQuestion) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AgentDatasourceTables ────────────────────────────

type AgentDatasourceTablesRepository interface {
	FindByAgentDatasourceID(ctx context.Context, agentDatasourceID int) ([]*model.AgentDatasourceTables, error)
	SaveBatch(ctx context.Context, agentDatasourceID int, tables []string) error
	DeleteByAgentDatasourceID(ctx context.Context, agentDatasourceID int) error
}

// ──────────────────────────── ChatSession ────────────────────────────

type ChatSessionRepository interface {
	FindByID(ctx context.Context, id string) (*model.ChatSession, error)
	FindByAgentIDAndUserID(ctx context.Context, agentID int, userID string) ([]*model.ChatSession, error)
	FindBySessionID(ctx context.Context, sessionID string) (*model.ChatSession, error)
	Create(ctx context.Context, session *model.ChatSession) error
	Update(ctx context.Context, session *model.ChatSession) error
	Delete(ctx context.Context, id string) error
	DeleteByAgentID(ctx context.Context, agentID int) error
}

// ──────────────────────────── ChatMessage ────────────────────────────

type ChatMessageRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) ([]*model.ChatMessage, error)
	Create(ctx context.Context, msg *model.ChatMessage) error
}

// ──────────────────────────── Datasource ────────────────────────────

type DatasourceRepository interface {
	FindByID(ctx context.Context, id string) (*model.Datasource, error)
	Page(ctx context.Context, page, size int, query *model.Datasource) ([]*model.Datasource, int64, error)
	Create(ctx context.Context, ds *model.Datasource) error
	Update(ctx context.Context, ds *model.Datasource) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── LogicalRelation ────────────────────────────

type LogicalRelationRepository interface {
	FindByID(ctx context.Context, id string) (*model.LogicalRelation, error)
	FindByDatasourceID(ctx context.Context, datasourceID int) ([]*model.LogicalRelation, error)
	Create(ctx context.Context, lr *model.LogicalRelation) error
	Update(ctx context.Context, lr *model.LogicalRelation) error
	Delete(ctx context.Context, id string) error
	DeleteByDatasourceID(ctx context.Context, datasourceID int) error
}

// ──────────────────────────── ModelConfig ────────────────────────────

type ModelConfigRepository interface {
	FindByID(ctx context.Context, id string) (*model.ModelConfig, error)
	Page(ctx context.Context, page, size int) ([]*model.ModelConfig, int64, error)
	FindActive(ctx context.Context) (*model.ModelConfig, error)
	Create(ctx context.Context, mc *model.ModelConfig) error
	Update(ctx context.Context, mc *model.ModelConfig) error
	Delete(ctx context.Context, id string) error
	DeactivateAll(ctx context.Context) error
}

// ──────────────────────────── UserPromptConfig ────────────────────────────

type UserPromptConfigRepository interface {
	FindByID(ctx context.Context, id string) (*model.UserPromptConfig, error)
	Page(ctx context.Context, page, size int, promptType string) ([]*model.UserPromptConfig, int64, error)
	FindByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error)
	FindActiveByType(ctx context.Context, promptType string) (*model.UserPromptConfig, error)
	FindActiveAllByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error)
	Create(ctx context.Context, pc *model.UserPromptConfig) error
	Update(ctx context.Context, pc *model.UserPromptConfig) error
	Delete(ctx context.Context, id string) error
	BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error
}

// ──────────────────────────── SemanticModel ────────────────────────────

type SemanticModelRepository interface {
	FindByID(ctx context.Context, id string) (*model.SemanticModel, error)
	Page(ctx context.Context, page, size int, query *model.SemanticModel) ([]*model.SemanticModel, int64, error)
	Create(ctx context.Context, sm *model.SemanticModel) error
	Update(ctx context.Context, sm *model.SemanticModel) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	BatchUpdateStatus(ctx context.Context, ids []string, status int) error
	BatchCreate(ctx context.Context, sms []*model.SemanticModel) error
}

// ──────────────────────────── BusinessKnowledge ────────────────────────────

type BusinessKnowledgeRepository interface {
	FindByID(ctx context.Context, id string) (*model.BusinessKnowledge, error)
	Page(ctx context.Context, page, size int, query *model.BusinessKnowledge) ([]*model.BusinessKnowledge, int64, error)
	Create(ctx context.Context, bk *model.BusinessKnowledge) error
	Update(ctx context.Context, bk *model.BusinessKnowledge) error
	Delete(ctx context.Context, id string) error
	UpdateRecall(ctx context.Context, id string, isRecall int) error
	FindByAgentID(ctx context.Context, agentID int64) ([]*model.BusinessKnowledge, error)
	BatchResetEmbedding(ctx context.Context, agentID int64) (int64, error)
}

// ──────────────────────────── DatasourceAccessor ────────────────────────────

type DatasourceAccessor interface {
	TestConnection(ds *model.Datasource) error
	GetTables(ds *model.Datasource) (interface{}, error)
	GetColumns(ds *model.Datasource, tableName string) (interface{}, error)
}
