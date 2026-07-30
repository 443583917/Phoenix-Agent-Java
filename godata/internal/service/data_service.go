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
