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
