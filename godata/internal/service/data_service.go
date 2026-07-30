package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/usecase"
)

// DataService is a thin pass-through wrapper around DataUsecase.
type DataService struct {
	uc *usecase.DataUsecase
}

// NewDataService creates a new DataService with the given usecase.
func NewDataService(uc *usecase.DataUsecase) *DataService {
	return &DataService{uc: uc}
}

// ──────────────────────────── Agent ────────────────────────────

func (s *DataService) CreateAgent(ctx context.Context, entity *model.Agent) (*model.Agent, error) {
	return s.uc.CreateAgent(ctx, entity)
}
func (s *DataService) UpdateAgent(ctx context.Context, entity *model.Agent) error {
	return s.uc.UpdateAgent(ctx, entity)
}
func (s *DataService) DeleteAgent(ctx context.Context, id string) error {
	return s.uc.DeleteAgent(ctx, id)
}
func (s *DataService) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	return s.uc.GetAgentByID(ctx, id)
}
func (s *DataService) PageAgent(ctx context.Context, page, size int, query *model.Agent) ([]*model.Agent, int64, error) {
	return s.uc.PageAgent(ctx, page, size, query)
}
func (s *DataService) ListAgent(ctx context.Context) ([]*model.Agent, error) {
	return s.uc.ListAgent(ctx)
}
func (s *DataService) PublishAgent(ctx context.Context, id string) error {
	return s.uc.PublishAgent(ctx, id)
}
func (s *DataService) OfflineAgent(ctx context.Context, id string) error {
	return s.uc.OfflineAgent(ctx, id)
}
func (s *DataService) GenerateAgentAPIKey(ctx context.Context, agentID string) (string, error) {
	return s.uc.GenerateAgentAPIKey(ctx, agentID)
}
func (s *DataService) ResetAgentAPIKey(ctx context.Context, agentID string) (string, error) {
	return s.uc.ResetAgentAPIKey(ctx, agentID)
}
func (s *DataService) DeleteAgentAPIKey(ctx context.Context, agentID string) error {
	return s.uc.DeleteAgentAPIKey(ctx, agentID)
}
func (s *DataService) ToggleAgentAPIKeyEnabled(ctx context.Context, agentID string) error {
	return s.uc.ToggleAgentAPIKeyEnabled(ctx, agentID)
}
func (s *DataService) GetAgentAPIKeyMasked(ctx context.Context, agentID string) (string, error) {
	return s.uc.GetAgentAPIKeyMasked(ctx, agentID)
}

// ──────────────────────────── AgentCategory ────────────────────────────

func (s *DataService) CreateAgentCategory(ctx context.Context, entity *model.AgentCategory) (*model.AgentCategory, error) {
	return s.uc.CreateAgentCategory(ctx, entity)
}
func (s *DataService) UpdateAgentCategory(ctx context.Context, entity *model.AgentCategory) error {
	return s.uc.UpdateAgentCategory(ctx, entity)
}
func (s *DataService) DeleteAgentCategory(ctx context.Context, id string) error {
	return s.uc.DeleteAgentCategory(ctx, id)
}
func (s *DataService) GetAgentCategoryByID(ctx context.Context, id string) (*model.AgentCategory, error) {
	return s.uc.GetAgentCategoryByID(ctx, id)
}
func (s *DataService) PageAgentCategory(ctx context.Context, page, size int, query *model.AgentCategory) ([]*model.AgentCategory, int64, error) {
	return s.uc.PageAgentCategory(ctx, page, size, query)
}
func (s *DataService) TreeAgentCategory(ctx context.Context) ([]*usecase.CategoryTreeVO, error) {
	return s.uc.TreeAgentCategory(ctx)
}

// ──────────────────────────── AgentDatasource ────────────────────────────

func (s *DataService) CreateAgentDatasource(ctx context.Context, entity *model.AgentDatasource) (*model.AgentDatasource, error) {
	return s.uc.CreateAgentDatasource(ctx, entity)
}
func (s *DataService) UpdateAgentDatasource(ctx context.Context, entity *model.AgentDatasource) error {
	return s.uc.UpdateAgentDatasource(ctx, entity)
}
func (s *DataService) DeleteAgentDatasource(ctx context.Context, id string) error {
	return s.uc.DeleteAgentDatasource(ctx, id)
}
func (s *DataService) GetAgentDatasourceByID(ctx context.Context, id string) (*model.AgentDatasource, error) {
	return s.uc.GetAgentDatasourceByID(ctx, id)
}
func (s *DataService) PageAgentDatasource(ctx context.Context, page, size int, query *model.AgentDatasource) ([]*model.AgentDatasource, int64, error) {
	return s.uc.PageAgentDatasource(ctx, page, size, query)
}
func (s *DataService) ToggleAgentDatasourceActive(ctx context.Context, id string) error {
	return s.uc.ToggleAgentDatasourceActive(ctx, id)
}
func (s *DataService) GetAgentDatasourceTables(ctx context.Context, agentDatasourceID int) ([]*model.AgentDatasourceTables, error) {
	return s.uc.GetAgentDatasourceTables(ctx, agentDatasourceID)
}
func (s *DataService) SaveAgentDatasourceTables(ctx context.Context, agentDatasourceID int, tables []string) error {
	return s.uc.SaveAgentDatasourceTables(ctx, agentDatasourceID, tables)
}

// ──────────────────────────── AgentKnowledge ────────────────────────────

func (s *DataService) CreateAgentKnowledge(ctx context.Context, entity *model.AgentKnowledge) (*model.AgentKnowledge, error) {
	return s.uc.CreateAgentKnowledge(ctx, entity)
}
func (s *DataService) UpdateAgentKnowledge(ctx context.Context, entity *model.AgentKnowledge) error {
	return s.uc.UpdateAgentKnowledge(ctx, entity)
}
func (s *DataService) DeleteAgentKnowledge(ctx context.Context, id string) error {
	return s.uc.DeleteAgentKnowledge(ctx, id)
}
func (s *DataService) GetAgentKnowledgeByID(ctx context.Context, id string) (*model.AgentKnowledge, error) {
	return s.uc.GetAgentKnowledgeByID(ctx, id)
}
func (s *DataService) PageAgentKnowledge(ctx context.Context, page, size int, query *model.AgentKnowledge) ([]*model.AgentKnowledge, int64, error) {
	return s.uc.PageAgentKnowledge(ctx, page, size, query)
}
func (s *DataService) ToggleAgentKnowledgeRecall(ctx context.Context, id string) error {
	return s.uc.ToggleAgentKnowledgeRecall(ctx, id)
}
func (s *DataService) RetryAgentKnowledgeEmbedding(ctx context.Context, id string) error {
	return s.uc.RetryAgentKnowledgeEmbedding(ctx, id)
}

// ──────────────────────────── AgentPresetQuestion ────────────────────────────

func (s *DataService) CreateAgentPresetQuestion(ctx context.Context, entity *model.AgentPresetQuestion) (*model.AgentPresetQuestion, error) {
	return s.uc.CreateAgentPresetQuestion(ctx, entity)
}
func (s *DataService) UpdateAgentPresetQuestion(ctx context.Context, entity *model.AgentPresetQuestion) error {
	return s.uc.UpdateAgentPresetQuestion(ctx, entity)
}
func (s *DataService) DeleteAgentPresetQuestion(ctx context.Context, id string) error {
	return s.uc.DeleteAgentPresetQuestion(ctx, id)
}
func (s *DataService) GetAgentPresetQuestionByID(ctx context.Context, id string) (*model.AgentPresetQuestion, error) {
	return s.uc.GetAgentPresetQuestionByID(ctx, id)
}
func (s *DataService) PageAgentPresetQuestion(ctx context.Context, page, size int, query *model.AgentPresetQuestion) ([]*model.AgentPresetQuestion, int64, error) {
	return s.uc.PageAgentPresetQuestion(ctx, page, size, query)
}

// ──────────────────────────── ChatSession ────────────────────────────

func (s *DataService) ListChatSessions(ctx context.Context, agentID int, userID string) ([]*model.ChatSession, error) {
	return s.uc.ListChatSessions(ctx, agentID, userID)
}
func (s *DataService) CreateChatSession(ctx context.Context, entity *model.ChatSession) (*model.ChatSession, error) {
	return s.uc.CreateChatSession(ctx, entity)
}
func (s *DataService) DeleteAllSessions(ctx context.Context, agentID int) error {
	return s.uc.DeleteAllSessions(ctx, agentID)
}
func (s *DataService) DeleteSession(ctx context.Context, id string) error {
	return s.uc.DeleteSession(ctx, id)
}
func (s *DataService) PinSession(ctx context.Context, id string, isPinned bool) error {
	return s.uc.PinSession(ctx, id, isPinned)
}
func (s *DataService) RenameSession(ctx context.Context, id string, title string) error {
	return s.uc.RenameSession(ctx, id, title)
}

// ──────────────────────────── ChatMessage ────────────────────────────

func (s *DataService) GetSessionMessages(ctx context.Context, sessionID string) ([]*model.ChatMessage, error) {
	return s.uc.GetSessionMessages(ctx, sessionID)
}
func (s *DataService) AddChatMessage(ctx context.Context, entity *model.ChatMessage) (*model.ChatMessage, error) {
	return s.uc.AddChatMessage(ctx, entity)
}

// ──────────────────────────── Datasource ────────────────────────────

func (s *DataService) PageDatasource(ctx context.Context, page, size int, query *model.Datasource) ([]*model.Datasource, int64, error) {
	return s.uc.PageDatasource(ctx, page, size, query)
}
func (s *DataService) GetDatasourceByID(ctx context.Context, id string) (*model.Datasource, error) {
	return s.uc.GetDatasourceByID(ctx, id)
}
func (s *DataService) CreateDatasource(ctx context.Context, entity *model.Datasource) (*model.Datasource, error) {
	return s.uc.CreateDatasource(ctx, entity)
}
func (s *DataService) UpdateDatasource(ctx context.Context, entity *model.Datasource) error {
	return s.uc.UpdateDatasource(ctx, entity)
}
func (s *DataService) DeleteDatasource(ctx context.Context, id string) error {
	return s.uc.DeleteDatasource(ctx, id)
}
func (s *DataService) TestDatasourceConnection(ctx context.Context, id string) error {
	return s.uc.TestDatasourceConnection(ctx, id)
}
func (s *DataService) GetDatasourceTables(ctx context.Context, id string) (interface{}, error) {
	return s.uc.GetDatasourceTables(ctx, id)
}
func (s *DataService) GetDatasourceColumns(ctx context.Context, id string, tableName string) (interface{}, error) {
	return s.uc.GetDatasourceColumns(ctx, id, tableName)
}

// ──────────────────────────── LogicalRelation ────────────────────────────

func (s *DataService) ListLogicalRelations(ctx context.Context, datasourceID int) ([]*model.LogicalRelation, error) {
	return s.uc.ListLogicalRelations(ctx, datasourceID)
}
func (s *DataService) CreateLogicalRelation(ctx context.Context, entity *model.LogicalRelation) (*model.LogicalRelation, error) {
	return s.uc.CreateLogicalRelation(ctx, entity)
}
func (s *DataService) UpdateLogicalRelations(ctx context.Context, datasourceID int, relations []*model.LogicalRelation) error {
	return s.uc.UpdateLogicalRelations(ctx, datasourceID, relations)
}
func (s *DataService) DeleteLogicalRelations(ctx context.Context, datasourceID int) error {
	return s.uc.DeleteLogicalRelations(ctx, datasourceID)
}
func (s *DataService) UpdateLogicalRelation(ctx context.Context, entity *model.LogicalRelation) error {
	return s.uc.UpdateLogicalRelation(ctx, entity)
}
func (s *DataService) DeleteLogicalRelation(ctx context.Context, id string) error {
	return s.uc.DeleteLogicalRelation(ctx, id)
}

// ──────────────────────────── ModelConfig ────────────────────────────

func (s *DataService) PageModelConfig(ctx context.Context, page, size int) ([]*model.ModelConfig, int64, error) {
	return s.uc.PageModelConfig(ctx, page, size)
}
func (s *DataService) CreateModelConfig(ctx context.Context, entity *model.ModelConfig) (*model.ModelConfig, error) {
	return s.uc.CreateModelConfig(ctx, entity)
}
func (s *DataService) UpdateModelConfig(ctx context.Context, entity *model.ModelConfig) error {
	return s.uc.UpdateModelConfig(ctx, entity)
}
func (s *DataService) DeleteModelConfig(ctx context.Context, id string) error {
	return s.uc.DeleteModelConfig(ctx, id)
}
func (s *DataService) ActivateModelConfig(ctx context.Context, id string) error {
	return s.uc.ActivateModelConfig(ctx, id)
}
func (s *DataService) CheckModelConfigReady(ctx context.Context) (bool, error) {
	return s.uc.CheckModelConfigReady(ctx)
}
func (s *DataService) TestModelConfig(ctx context.Context, entity *model.ModelConfig) (bool, error) {
	return s.uc.TestModelConfig(ctx, entity)
}

// ──────────────────────────── UserPromptConfig ────────────────────────────

func (s *DataService) PagePromptConfig(ctx context.Context, page, size int, promptType string) ([]*model.UserPromptConfig, int64, error) {
	return s.uc.PagePromptConfig(ctx, page, size, promptType)
}
func (s *DataService) GetPromptConfigByID(ctx context.Context, id string) (*model.UserPromptConfig, error) {
	return s.uc.GetPromptConfigByID(ctx, id)
}
func (s *DataService) SavePromptConfig(ctx context.Context, entity *model.UserPromptConfig) (*model.UserPromptConfig, error) {
	return s.uc.SavePromptConfig(ctx, entity)
}
func (s *DataService) ListPromptConfigByType(ctx context.Context, promptType string, page, size int) ([]*model.UserPromptConfig, int64, error) {
	return s.uc.ListPromptConfigByType(ctx, promptType, page, size)
}
func (s *DataService) GetActivePromptConfigByType(ctx context.Context, promptType string) (*model.UserPromptConfig, error) {
	return s.uc.GetActivePromptConfigByType(ctx, promptType)
}
func (s *DataService) GetActiveAllPromptConfigByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error) {
	return s.uc.GetActiveAllPromptConfigByType(ctx, promptType)
}
func (s *DataService) DeletePromptConfig(ctx context.Context, id string) error {
	return s.uc.DeletePromptConfig(ctx, id)
}
func (s *DataService) EnablePromptConfig(ctx context.Context, id string) error {
	return s.uc.EnablePromptConfig(ctx, id)
}
func (s *DataService) DisablePromptConfig(ctx context.Context, id string) error {
	return s.uc.DisablePromptConfig(ctx, id)
}
func (s *DataService) BatchEnablePromptConfig(ctx context.Context, ids []string) error {
	return s.uc.BatchEnablePromptConfig(ctx, ids)
}
func (s *DataService) BatchDisablePromptConfig(ctx context.Context, ids []string) error {
	return s.uc.BatchDisablePromptConfig(ctx, ids)
}
func (s *DataService) SetPromptConfigPriority(ctx context.Context, id string, priority int) error {
	return s.uc.SetPromptConfigPriority(ctx, id, priority)
}
func (s *DataService) SetPromptConfigDisplayOrder(ctx context.Context, id string, displayOrder int) error {
	return s.uc.SetPromptConfigDisplayOrder(ctx, id, displayOrder)
}

// ──────────────────────────── SemanticModel ────────────────────────────

func (s *DataService) PageSemanticModel(ctx context.Context, page, size int, query *model.SemanticModel) ([]*model.SemanticModel, int64, error) {
	return s.uc.PageSemanticModel(ctx, page, size, query)
}
func (s *DataService) GetSemanticModelByID(ctx context.Context, id string) (*model.SemanticModel, error) {
	return s.uc.GetSemanticModelByID(ctx, id)
}
func (s *DataService) CreateSemanticModel(ctx context.Context, entity *model.SemanticModel) (*model.SemanticModel, error) {
	return s.uc.CreateSemanticModel(ctx, entity)
}
func (s *DataService) UpdateSemanticModel(ctx context.Context, entity *model.SemanticModel) error {
	return s.uc.UpdateSemanticModel(ctx, entity)
}
func (s *DataService) DeleteSemanticModel(ctx context.Context, id string) error {
	return s.uc.DeleteSemanticModel(ctx, id)
}
func (s *DataService) BatchDeleteSemanticModel(ctx context.Context, ids []string) error {
	return s.uc.BatchDeleteSemanticModel(ctx, ids)
}
func (s *DataService) EnableSemanticModels(ctx context.Context, ids []string) error {
	return s.uc.EnableSemanticModels(ctx, ids)
}
func (s *DataService) DisableSemanticModels(ctx context.Context, ids []string) error {
	return s.uc.DisableSemanticModels(ctx, ids)
}
func (s *DataService) BatchCreateSemanticModels(ctx context.Context, entities []*model.SemanticModel) (int, error) {
	return s.uc.BatchCreateSemanticModels(ctx, entities)
}

// ──────────────────────────── BusinessKnowledge ────────────────────────────

func (s *DataService) PageBusinessKnowledge(ctx context.Context, page, size int, query *model.BusinessKnowledge) ([]*model.BusinessKnowledge, int64, error) {
	return s.uc.PageBusinessKnowledge(ctx, page, size, query)
}
func (s *DataService) GetBusinessKnowledgeByID(ctx context.Context, id string) (*model.BusinessKnowledge, error) {
	return s.uc.GetBusinessKnowledgeByID(ctx, id)
}
func (s *DataService) CreateBusinessKnowledge(ctx context.Context, entity *model.BusinessKnowledge) (*model.BusinessKnowledge, error) {
	return s.uc.CreateBusinessKnowledge(ctx, entity)
}
func (s *DataService) UpdateBusinessKnowledge(ctx context.Context, entity *model.BusinessKnowledge) error {
	return s.uc.UpdateBusinessKnowledge(ctx, entity)
}
func (s *DataService) DeleteBusinessKnowledge(ctx context.Context, id string) error {
	return s.uc.DeleteBusinessKnowledge(ctx, id)
}
func (s *DataService) RetryBusinessKnowledgeEmbedding(ctx context.Context, id string) error {
	return s.uc.RetryBusinessKnowledgeEmbedding(ctx, id)
}
func (s *DataService) ToggleBusinessKnowledgeRecall(ctx context.Context, id string, isRecall bool) error {
	return s.uc.ToggleBusinessKnowledgeRecall(ctx, id, isRecall)
}
func (s *DataService) ToggleBusinessKnowledgeRecallOn(ctx context.Context, id string) error {
	return s.uc.ToggleBusinessKnowledgeRecallOn(ctx, id)
}
func (s *DataService) RefreshBusinessKnowledgeVectorStore(ctx context.Context, agentID int64) (int64, error) {
	return s.uc.RefreshBusinessKnowledgeVectorStore(ctx, agentID)
}
