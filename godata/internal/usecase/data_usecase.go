package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"github.com/phoenix-agent-go/infra/id"
)

var (
	ErrAgentNotFound                 = &AppError{Code: 601001, Msg: "智能体不存在"}
	ErrAgentCategoryNotFound         = &AppError{Code: 602001, Msg: "分类不存在"}
	ErrAgentDatasourceNotFound       = &AppError{Code: 603001, Msg: "智能体数据源关联不存在"}
	ErrAgentKnowledgeNotFound        = &AppError{Code: 604001, Msg: "知识库条目不存在"}
	ErrAgentPresetQuestionNotFound   = &AppError{Code: 605001, Msg: "预设问题不存在"}
	ErrChatSessionNotFound           = &AppError{Code: 606001, Msg: "会话不存在"}
	ErrDatasourceNotFound            = &AppError{Code: 607001, Msg: "数据源不存在"}
	ErrLogicalRelationNotFound       = &AppError{Code: 608001, Msg: "逻辑关系不存在"}
	ErrModelConfigNotFound           = &AppError{Code: 609001, Msg: "模型配置不存在"}
	ErrPromptConfigNotFound          = &AppError{Code: 610001, Msg: "Prompt配置不存在"}
	ErrSemanticModelNotFound         = &AppError{Code: 611001, Msg: "语义模型不存在"}
	ErrBusinessKnowledgeNotFound     = &AppError{Code: 612001, Msg: "业务知识不存在"}
)

type DataUsecase struct {
	agentRepo                 repository.AgentRepository
	agentCategoryRepo         repository.AgentCategoryRepository
	agentDatasourceRepo       repository.AgentDatasourceRepository
	agentKnowledgeRepo        repository.AgentKnowledgeRepository
	agentPresetQuestionRepo   repository.AgentPresetQuestionRepository
	agentDatasourceTablesRepo repository.AgentDatasourceTablesRepository
	chatSessionRepo           repository.ChatSessionRepository
	chatMessageRepo           repository.ChatMessageRepository
	datasourceRepo            repository.DatasourceRepository
	logicalRelationRepo       repository.LogicalRelationRepository
	modelConfigRepo           repository.ModelConfigRepository
	userPromptConfigRepo      repository.UserPromptConfigRepository
	semanticModelRepo         repository.SemanticModelRepository
	businessKnowledgeRepo     repository.BusinessKnowledgeRepository
	dsAccessor                repository.DatasourceAccessor
}

func NewDataUsecase(
	agentRepo repository.AgentRepository,
	agentCategoryRepo repository.AgentCategoryRepository,
	agentDatasourceRepo repository.AgentDatasourceRepository,
	agentKnowledgeRepo repository.AgentKnowledgeRepository,
	agentPresetQuestionRepo repository.AgentPresetQuestionRepository,
	agentDatasourceTablesRepo repository.AgentDatasourceTablesRepository,
	chatSessionRepo repository.ChatSessionRepository,
	chatMessageRepo repository.ChatMessageRepository,
	datasourceRepo repository.DatasourceRepository,
	logicalRelationRepo repository.LogicalRelationRepository,
	modelConfigRepo repository.ModelConfigRepository,
	userPromptConfigRepo repository.UserPromptConfigRepository,
	semanticModelRepo repository.SemanticModelRepository,
	businessKnowledgeRepo repository.BusinessKnowledgeRepository,
	dsAccessor repository.DatasourceAccessor,
) *DataUsecase {
	return &DataUsecase{
		agentRepo:                 agentRepo,
		agentCategoryRepo:         agentCategoryRepo,
		agentDatasourceRepo:       agentDatasourceRepo,
		agentKnowledgeRepo:        agentKnowledgeRepo,
		agentPresetQuestionRepo:   agentPresetQuestionRepo,
		agentDatasourceTablesRepo: agentDatasourceTablesRepo,
		chatSessionRepo:           chatSessionRepo,
		chatMessageRepo:           chatMessageRepo,
		datasourceRepo:            datasourceRepo,
		logicalRelationRepo:       logicalRelationRepo,
		modelConfigRepo:           modelConfigRepo,
		userPromptConfigRepo:      userPromptConfigRepo,
		semanticModelRepo:         semanticModelRepo,
		businessKnowledgeRepo:     businessKnowledgeRepo,
		dsAccessor:                dsAccessor,
	}
}

// ──────────────────────────── helpers ────────────────────────────

func dataGenID() string {
	return strconv.FormatUint(id.MustGenerateID(), 10)
}

func generateAPIKey() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return "ag_" + hex.EncodeToString(b)
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// ──────────────────────────── Agent ────────────────────────────

func (u *DataUsecase) CreateAgent(ctx context.Context, entity *model.Agent) (*model.Agent, error) {
	entity.ID = dataGenID()
	if entity.Status == "" {
		entity.Status = "draft"
	}
	if err := u.agentRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateAgent(ctx context.Context, entity *model.Agent) error {
	existing, err := u.agentRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	return u.agentRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteAgent(ctx context.Context, id string) error {
	existing, err := u.agentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	return u.agentRepo.Delete(ctx, id)
}

func (u *DataUsecase) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	entity, err := u.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrAgentNotFound
	}
	return entity, nil
}

func (u *DataUsecase) PageAgent(ctx context.Context, page, size int, query *model.Agent) ([]*model.Agent, int64, error) {
	return u.agentRepo.Page(ctx, page, size, query)
}

func (u *DataUsecase) ListAgent(ctx context.Context) ([]*model.Agent, error) {
	return u.agentRepo.List(ctx)
}

// PublishAgent sets agent status to "published".
func (u *DataUsecase) PublishAgent(ctx context.Context, id string) error {
	existing, err := u.agentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	existing.Status = "published"
	return u.agentRepo.Update(ctx, existing)
}

// OfflineAgent sets agent status to "offline".
func (u *DataUsecase) OfflineAgent(ctx context.Context, id string) error {
	existing, err := u.agentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	existing.Status = "offline"
	return u.agentRepo.Update(ctx, existing)
}

// ──────────────────────────── Agent API Key ────────────────────────────

// GenerateAgentAPIKey creates a new API key for the agent and returns the raw key.
func (u *DataUsecase) GenerateAgentAPIKey(ctx context.Context, agentID string) (string, error) {
	existing, err := u.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", ErrAgentNotFound
	}
	rawKey := generateAPIKey()
	existing.ApiKey = rawKey
	existing.ApiKeyEnabled = 1
	if err := u.agentRepo.Update(ctx, existing); err != nil {
		return "", err
	}
	return rawKey, nil
}

// ResetAgentAPIKey replaces the agent's API key with a new one and returns the raw key.
func (u *DataUsecase) ResetAgentAPIKey(ctx context.Context, agentID string) (string, error) {
	existing, err := u.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", ErrAgentNotFound
	}
	rawKey := generateAPIKey()
	existing.ApiKey = rawKey
	existing.ApiKeyEnabled = 1
	if err := u.agentRepo.Update(ctx, existing); err != nil {
		return "", err
	}
	return rawKey, nil
}

// DeleteAgentAPIKey removes the API key from the agent and disables it.
func (u *DataUsecase) DeleteAgentAPIKey(ctx context.Context, agentID string) error {
	existing, err := u.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	existing.ApiKey = ""
	existing.ApiKeyEnabled = 0
	return u.agentRepo.Update(ctx, existing)
}

// ToggleAgentAPIKeyEnabled toggles the api_key_enabled flag.
func (u *DataUsecase) ToggleAgentAPIKeyEnabled(ctx context.Context, agentID string) error {
	existing, err := u.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentNotFound
	}
	if existing.ApiKeyEnabled == 1 {
		existing.ApiKeyEnabled = 0
	} else {
		existing.ApiKeyEnabled = 1
	}
	return u.agentRepo.Update(ctx, existing)
}

// GetAgentAPIKeyMasked returns a masked version of the agent's API key.
func (u *DataUsecase) GetAgentAPIKeyMasked(ctx context.Context, agentID string) (string, error) {
	existing, err := u.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", ErrAgentNotFound
	}
	if existing.ApiKey == "" {
		return "", nil
	}
	return maskAPIKey(existing.ApiKey), nil
}

// ──────────────────────────── AgentCategory ────────────────────────────

type CategoryTreeVO struct {
	ID       string           `json:"id"`
	Pid      string           `json:"pid"`
	Name     string           `json:"name"`
	Sn       string           `json:"sn"`
	Children []*CategoryTreeVO `json:"children,omitempty"`
}

func (u *DataUsecase) CreateAgentCategory(ctx context.Context, entity *model.AgentCategory) (*model.AgentCategory, error) {
	entity.ID = dataGenID()
	if err := u.agentCategoryRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateAgentCategory(ctx context.Context, entity *model.AgentCategory) error {
	existing, err := u.agentCategoryRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentCategoryNotFound
	}
	return u.agentCategoryRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteAgentCategory(ctx context.Context, id string) error {
	existing, err := u.agentCategoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentCategoryNotFound
	}
	return u.agentCategoryRepo.Delete(ctx, id)
}

func (u *DataUsecase) GetAgentCategoryByID(ctx context.Context, id string) (*model.AgentCategory, error) {
	entity, err := u.agentCategoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrAgentCategoryNotFound
	}
	return entity, nil
}

func (u *DataUsecase) PageAgentCategory(ctx context.Context, page, size int, query *model.AgentCategory) ([]*model.AgentCategory, int64, error) {
	return u.agentCategoryRepo.Page(ctx, page, size, query)
}

// TreeAgentCategory returns the full category tree.
func (u *DataUsecase) TreeAgentCategory(ctx context.Context) ([]*CategoryTreeVO, error) {
	all, err := u.agentCategoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(all, ""), nil
}

func buildCategoryTree(all []*model.AgentCategory, pid string) []*CategoryTreeVO {
	var tree []*CategoryTreeVO
	for _, cat := range all {
		if cat.Pid == pid {
			node := &CategoryTreeVO{
				ID:   cat.ID,
				Pid:  cat.Pid,
				Name: cat.Name,
				Sn:   cat.Sn,
			}
			children := buildCategoryTree(all, cat.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// ──────────────────────────── AgentDatasource ────────────────────────────

func (u *DataUsecase) CreateAgentDatasource(ctx context.Context, entity *model.AgentDatasource) (*model.AgentDatasource, error) {
	entity.ID = dataGenID()
	if entity.IsActive == 0 {
		entity.IsActive = 1
	}
	if err := u.agentDatasourceRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateAgentDatasource(ctx context.Context, entity *model.AgentDatasource) error {
	existing, err := u.agentDatasourceRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentDatasourceNotFound
	}
	return u.agentDatasourceRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteAgentDatasource(ctx context.Context, id string) error {
	existing, err := u.agentDatasourceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentDatasourceNotFound
	}
	return u.agentDatasourceRepo.Delete(ctx, id)
}

func (u *DataUsecase) GetAgentDatasourceByID(ctx context.Context, id string) (*model.AgentDatasource, error) {
	entity, err := u.agentDatasourceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrAgentDatasourceNotFound
	}
	return entity, nil
}

func (u *DataUsecase) PageAgentDatasource(ctx context.Context, page, size int, query *model.AgentDatasource) ([]*model.AgentDatasource, int64, error) {
	return u.agentDatasourceRepo.Page(ctx, page, size, query)
}

// ToggleAgentDatasourceActive toggles the is_active flag.
func (u *DataUsecase) ToggleAgentDatasourceActive(ctx context.Context, id string) error {
	existing, err := u.agentDatasourceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentDatasourceNotFound
	}
	if existing.IsActive == 1 {
		existing.IsActive = 0
	} else {
		existing.IsActive = 1
	}
	return u.agentDatasourceRepo.Update(ctx, existing)
}

// GetAgentDatasourceTables returns the table list for a given datasource-agent link.
func (u *DataUsecase) GetAgentDatasourceTables(ctx context.Context, agentDatasourceID int) ([]*model.AgentDatasourceTables, error) {
	return u.agentDatasourceTablesRepo.FindByAgentDatasourceID(ctx, agentDatasourceID)
}

// SaveAgentDatasourceTables replaces the table selection for a datasource-agent link.
func (u *DataUsecase) SaveAgentDatasourceTables(ctx context.Context, agentDatasourceID int, tables []string) error {
	return u.agentDatasourceTablesRepo.SaveBatch(ctx, agentDatasourceID, tables)
}

// ──────────────────────────── AgentKnowledge ────────────────────────────

func (u *DataUsecase) CreateAgentKnowledge(ctx context.Context, entity *model.AgentKnowledge) (*model.AgentKnowledge, error) {
	entity.ID = dataGenID()
	if entity.EmbeddingStatus == "" {
		entity.EmbeddingStatus = "pending"
	}
	if entity.IsRecall == 0 {
		entity.IsRecall = 1
	}
	if err := u.agentKnowledgeRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateAgentKnowledge(ctx context.Context, entity *model.AgentKnowledge) error {
	existing, err := u.agentKnowledgeRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentKnowledgeNotFound
	}
	return u.agentKnowledgeRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteAgentKnowledge(ctx context.Context, id string) error {
	existing, err := u.agentKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentKnowledgeNotFound
	}
	return u.agentKnowledgeRepo.Delete(ctx, id)
}

func (u *DataUsecase) GetAgentKnowledgeByID(ctx context.Context, id string) (*model.AgentKnowledge, error) {
	entity, err := u.agentKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrAgentKnowledgeNotFound
	}
	return entity, nil
}

func (u *DataUsecase) PageAgentKnowledge(ctx context.Context, page, size int, query *model.AgentKnowledge) ([]*model.AgentKnowledge, int64, error) {
	return u.agentKnowledgeRepo.Page(ctx, page, size, query)
}

// ToggleAgentKnowledgeRecall toggles the is_recall flag.
func (u *DataUsecase) ToggleAgentKnowledgeRecall(ctx context.Context, id string) error {
	existing, err := u.agentKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentKnowledgeNotFound
	}
	if existing.IsRecall == 1 {
		existing.IsRecall = 0
	} else {
		existing.IsRecall = 1
	}
	return u.agentKnowledgeRepo.Update(ctx, existing)
}

// RetryAgentKnowledgeEmbedding resets embedding status to "pending" for retry.
func (u *DataUsecase) RetryAgentKnowledgeEmbedding(ctx context.Context, id string) error {
	existing, err := u.agentKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentKnowledgeNotFound
	}
	existing.EmbeddingStatus = "pending"
	existing.ErrorMsg = ""
	return u.agentKnowledgeRepo.Update(ctx, existing)
}

// ──────────────────────────── AgentPresetQuestion ────────────────────────────

func (u *DataUsecase) CreateAgentPresetQuestion(ctx context.Context, entity *model.AgentPresetQuestion) (*model.AgentPresetQuestion, error) {
	entity.ID = dataGenID()
	if entity.IsActive == nil {
		val := true
		entity.IsActive = &val
	}
	if err := u.agentPresetQuestionRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateAgentPresetQuestion(ctx context.Context, entity *model.AgentPresetQuestion) error {
	existing, err := u.agentPresetQuestionRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentPresetQuestionNotFound
	}
	return u.agentPresetQuestionRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteAgentPresetQuestion(ctx context.Context, id string) error {
	existing, err := u.agentPresetQuestionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAgentPresetQuestionNotFound
	}
	return u.agentPresetQuestionRepo.Delete(ctx, id)
}

func (u *DataUsecase) GetAgentPresetQuestionByID(ctx context.Context, id string) (*model.AgentPresetQuestion, error) {
	entity, err := u.agentPresetQuestionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrAgentPresetQuestionNotFound
	}
	return entity, nil
}

func (u *DataUsecase) PageAgentPresetQuestion(ctx context.Context, page, size int, query *model.AgentPresetQuestion) ([]*model.AgentPresetQuestion, int64, error) {
	return u.agentPresetQuestionRepo.Page(ctx, page, size, query)
}

// ──────────────────────────── ChatSession ────────────────────────────

func (u *DataUsecase) ListChatSessions(ctx context.Context, agentID int, userID string) ([]*model.ChatSession, error) {
	return u.chatSessionRepo.FindByAgentIDAndUserID(ctx, agentID, userID)
}

func (u *DataUsecase) CreateChatSession(ctx context.Context, entity *model.ChatSession) (*model.ChatSession, error) {
	entity.ID = dataGenID()
	if entity.Status == "" {
		entity.Status = "active"
	}
	if err := u.chatSessionRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) DeleteAllSessions(ctx context.Context, agentID int) error {
	return u.chatSessionRepo.DeleteByAgentID(ctx, agentID)
}

func (u *DataUsecase) DeleteSession(ctx context.Context, id string) error {
	existing, err := u.chatSessionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChatSessionNotFound
	}
	return u.chatSessionRepo.Delete(ctx, id)
}

func (u *DataUsecase) PinSession(ctx context.Context, id string, isPinned bool) error {
	existing, err := u.chatSessionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChatSessionNotFound
	}
	existing.IsPinned = isPinned
	return u.chatSessionRepo.Update(ctx, existing)
}

func (u *DataUsecase) RenameSession(ctx context.Context, id string, title string) error {
	existing, err := u.chatSessionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChatSessionNotFound
	}
	existing.Title = title
	return u.chatSessionRepo.Update(ctx, existing)
}

// ──────────────────────────── ChatMessage ────────────────────────────

func (u *DataUsecase) GetSessionMessages(ctx context.Context, sessionID string) ([]*model.ChatMessage, error) {
	return u.chatMessageRepo.FindBySessionID(ctx, sessionID)
}

func (u *DataUsecase) GetChatSessionByID(ctx context.Context, sessionID string) (*model.ChatSession, error) {
	entity, err := u.chatSessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrChatSessionNotFound
	}
	return entity, nil
}

func (u *DataUsecase) AddChatMessage(ctx context.Context, entity *model.ChatMessage) (*model.ChatMessage, error) {
	entity.ID = dataGenID()
	if err := u.chatMessageRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// ──────────────────────────── Datasource ────────────────────────────

func (u *DataUsecase) PageDatasource(ctx context.Context, page, size int, query *model.Datasource) ([]*model.Datasource, int64, error) {
	return u.datasourceRepo.Page(ctx, page, size, query)
}

func (u *DataUsecase) GetDatasourceByID(ctx context.Context, id string) (*model.Datasource, error) {
	entity, err := u.datasourceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrDatasourceNotFound
	}
	return entity, nil
}

func (u *DataUsecase) CreateDatasource(ctx context.Context, entity *model.Datasource) (*model.Datasource, error) {
	entity.ID = dataGenID()
	if entity.Status == "" {
		entity.Status = "active"
	}
	if err := u.datasourceRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateDatasource(ctx context.Context, entity *model.Datasource) error {
	existing, err := u.datasourceRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrDatasourceNotFound
	}
	return u.datasourceRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteDatasource(ctx context.Context, id string) error {
	existing, err := u.datasourceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrDatasourceNotFound
	}
	return u.datasourceRepo.Delete(ctx, id)
}

func (u *DataUsecase) TestDatasourceConnection(ctx context.Context, id string) error {
	ds, err := u.datasourceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ds == nil {
		return ErrDatasourceNotFound
	}
	if u.dsAccessor == nil {
		return fmt.Errorf("datasource accessor not configured")
	}
	return u.dsAccessor.TestConnection(ds)
}

func (u *DataUsecase) GetDatasourceTables(ctx context.Context, id string) (interface{}, error) {
	ds, err := u.datasourceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, ErrDatasourceNotFound
	}
	if u.dsAccessor == nil {
		return nil, fmt.Errorf("datasource accessor not configured")
	}
	return u.dsAccessor.GetTables(ds)
}

func (u *DataUsecase) GetDatasourceColumns(ctx context.Context, id string, tableName string) (interface{}, error) {
	ds, err := u.datasourceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, ErrDatasourceNotFound
	}
	if u.dsAccessor == nil {
		return nil, fmt.Errorf("datasource accessor not configured")
	}
	return u.dsAccessor.GetColumns(ds, tableName)
}

// ──────────────────────────── LogicalRelation ────────────────────────────

func (u *DataUsecase) ListLogicalRelations(ctx context.Context, datasourceID int) ([]*model.LogicalRelation, error) {
	return u.logicalRelationRepo.FindByDatasourceID(ctx, datasourceID)
}

func (u *DataUsecase) CreateLogicalRelation(ctx context.Context, entity *model.LogicalRelation) (*model.LogicalRelation, error) {
	entity.ID = dataGenID()
	if err := u.logicalRelationRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateLogicalRelations(ctx context.Context, datasourceID int, relations []*model.LogicalRelation) error {
	if err := u.logicalRelationRepo.DeleteByDatasourceID(ctx, datasourceID); err != nil {
		return err
	}
	for _, rel := range relations {
		rel.ID = dataGenID()
		rel.DatasourceId = datasourceID
		if err := u.logicalRelationRepo.Create(ctx, rel); err != nil {
			return err
		}
	}
	return nil
}

func (u *DataUsecase) DeleteLogicalRelations(ctx context.Context, datasourceID int) error {
	return u.logicalRelationRepo.DeleteByDatasourceID(ctx, datasourceID)
}

func (u *DataUsecase) UpdateLogicalRelation(ctx context.Context, entity *model.LogicalRelation) error {
	existing, err := u.logicalRelationRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrLogicalRelationNotFound
	}
	return u.logicalRelationRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteLogicalRelation(ctx context.Context, id string) error {
	existing, err := u.logicalRelationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrLogicalRelationNotFound
	}
	return u.logicalRelationRepo.Delete(ctx, id)
}

// ──────────────────────────── ModelConfig ────────────────────────────

func (u *DataUsecase) PageModelConfig(ctx context.Context, page, size int) ([]*model.ModelConfig, int64, error) {
	return u.modelConfigRepo.Page(ctx, page, size)
}

func (u *DataUsecase) CreateModelConfig(ctx context.Context, entity *model.ModelConfig) (*model.ModelConfig, error) {
	entity.ID = dataGenID()
	if err := u.modelConfigRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateModelConfig(ctx context.Context, entity *model.ModelConfig) error {
	existing, err := u.modelConfigRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrModelConfigNotFound
	}
	return u.modelConfigRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteModelConfig(ctx context.Context, id string) error {
	existing, err := u.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrModelConfigNotFound
	}
	return u.modelConfigRepo.Delete(ctx, id)
}

func (u *DataUsecase) ActivateModelConfig(ctx context.Context, id string) error {
	existing, err := u.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrModelConfigNotFound
	}
	if err := u.modelConfigRepo.DeactivateAll(ctx); err != nil {
		return err
	}
	existing.IsActive = true
	return u.modelConfigRepo.Update(ctx, existing)
}

func (u *DataUsecase) CheckModelConfigReady(ctx context.Context) (bool, error) {
	active, err := u.modelConfigRepo.FindActive(ctx)
	if err != nil {
		return false, err
	}
	return active != nil, nil
}

func (u *DataUsecase) TestModelConfig(ctx context.Context, entity *model.ModelConfig) (bool, error) {
	if entity.BaseUrl == "" || entity.ModelName == "" {
		return false, fmt.Errorf("baseUrl and modelName are required")
	}
	return true, nil
}

// ──────────────────────────── UserPromptConfig ────────────────────────────

func (u *DataUsecase) PagePromptConfig(ctx context.Context, page, size int, promptType string) ([]*model.UserPromptConfig, int64, error) {
	return u.userPromptConfigRepo.Page(ctx, page, size, promptType)
}

func (u *DataUsecase) GetPromptConfigByID(ctx context.Context, id string) (*model.UserPromptConfig, error) {
	entity, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPromptConfigNotFound
	}
	return entity, nil
}

func (u *DataUsecase) SavePromptConfig(ctx context.Context, entity *model.UserPromptConfig) (*model.UserPromptConfig, error) {
	if entity.ID == "" {
		entity.ID = dataGenID()
		if err := u.userPromptConfigRepo.Create(ctx, entity); err != nil {
			return nil, err
		}
		return entity, nil
	}
	if err := u.userPromptConfigRepo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) ListPromptConfigByType(ctx context.Context, promptType string, page, size int) ([]*model.UserPromptConfig, int64, error) {
	return u.userPromptConfigRepo.Page(ctx, page, size, promptType)
}

func (u *DataUsecase) GetActivePromptConfigByType(ctx context.Context, promptType string) (*model.UserPromptConfig, error) {
	return u.userPromptConfigRepo.FindActiveByType(ctx, promptType)
}

func (u *DataUsecase) GetActiveAllPromptConfigByType(ctx context.Context, promptType string) ([]*model.UserPromptConfig, error) {
	return u.userPromptConfigRepo.FindActiveAllByType(ctx, promptType)
}

func (u *DataUsecase) DeletePromptConfig(ctx context.Context, id string) error {
	existing, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPromptConfigNotFound
	}
	return u.userPromptConfigRepo.Delete(ctx, id)
}

func (u *DataUsecase) EnablePromptConfig(ctx context.Context, id string) error {
	existing, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPromptConfigNotFound
	}
	existing.Enabled = true
	return u.userPromptConfigRepo.Update(ctx, existing)
}

func (u *DataUsecase) DisablePromptConfig(ctx context.Context, id string) error {
	existing, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPromptConfigNotFound
	}
	existing.Enabled = false
	return u.userPromptConfigRepo.Update(ctx, existing)
}

func (u *DataUsecase) BatchEnablePromptConfig(ctx context.Context, ids []string) error {
	return u.userPromptConfigRepo.BatchUpdate(ctx, ids, map[string]interface{}{"enabled": true})
}

func (u *DataUsecase) BatchDisablePromptConfig(ctx context.Context, ids []string) error {
	return u.userPromptConfigRepo.BatchUpdate(ctx, ids, map[string]interface{}{"enabled": false})
}

func (u *DataUsecase) SetPromptConfigPriority(ctx context.Context, id string, priority int) error {
	existing, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPromptConfigNotFound
	}
	existing.Priority = priority
	return u.userPromptConfigRepo.Update(ctx, existing)
}

func (u *DataUsecase) SetPromptConfigDisplayOrder(ctx context.Context, id string, displayOrder int) error {
	existing, err := u.userPromptConfigRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPromptConfigNotFound
	}
	existing.DisplayOrder = displayOrder
	return u.userPromptConfigRepo.Update(ctx, existing)
}

// ──────────────────────────── SemanticModel ────────────────────────────

func (u *DataUsecase) PageSemanticModel(ctx context.Context, page, size int, query *model.SemanticModel) ([]*model.SemanticModel, int64, error) {
	return u.semanticModelRepo.Page(ctx, page, size, query)
}

func (u *DataUsecase) GetSemanticModelByID(ctx context.Context, id string) (*model.SemanticModel, error) {
	entity, err := u.semanticModelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrSemanticModelNotFound
	}
	return entity, nil
}

func (u *DataUsecase) CreateSemanticModel(ctx context.Context, entity *model.SemanticModel) (*model.SemanticModel, error) {
	entity.ID = dataGenID()
	if entity.Status == 0 {
		entity.Status = 1
	}
	if err := u.semanticModelRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateSemanticModel(ctx context.Context, entity *model.SemanticModel) error {
	existing, err := u.semanticModelRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSemanticModelNotFound
	}
	return u.semanticModelRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteSemanticModel(ctx context.Context, id string) error {
	existing, err := u.semanticModelRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSemanticModelNotFound
	}
	return u.semanticModelRepo.Delete(ctx, id)
}

func (u *DataUsecase) BatchDeleteSemanticModel(ctx context.Context, ids []string) error {
	return u.semanticModelRepo.BatchDelete(ctx, ids)
}

func (u *DataUsecase) EnableSemanticModels(ctx context.Context, ids []string) error {
	return u.semanticModelRepo.BatchUpdateStatus(ctx, ids, 1)
}

func (u *DataUsecase) DisableSemanticModels(ctx context.Context, ids []string) error {
	return u.semanticModelRepo.BatchUpdateStatus(ctx, ids, 0)
}

func (u *DataUsecase) BatchCreateSemanticModels(ctx context.Context, entities []*model.SemanticModel) (int, error) {
	for i, e := range entities {
		if e.ID == "" {
			e.ID = dataGenID()
		}
		if e.Status == 0 {
			e.Status = 1
		}
		entities[i] = e
	}
	if err := u.semanticModelRepo.BatchCreate(ctx, entities); err != nil {
		return 0, err
	}
	return len(entities), nil
}

// ──────────────────────────── BusinessKnowledge ────────────────────────────

func (u *DataUsecase) PageBusinessKnowledge(ctx context.Context, page, size int, query *model.BusinessKnowledge) ([]*model.BusinessKnowledge, int64, error) {
	return u.businessKnowledgeRepo.Page(ctx, page, size, query)
}

func (u *DataUsecase) GetBusinessKnowledgeByID(ctx context.Context, id string) (*model.BusinessKnowledge, error) {
	entity, err := u.businessKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrBusinessKnowledgeNotFound
	}
	return entity, nil
}

func (u *DataUsecase) CreateBusinessKnowledge(ctx context.Context, entity *model.BusinessKnowledge) (*model.BusinessKnowledge, error) {
	entity.ID = dataGenID()
	if entity.EmbeddingStatus == "" {
		entity.EmbeddingStatus = "pending"
	}
	if entity.IsRecall == 0 {
		entity.IsRecall = 1
	}
	if err := u.businessKnowledgeRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *DataUsecase) UpdateBusinessKnowledge(ctx context.Context, entity *model.BusinessKnowledge) error {
	existing, err := u.businessKnowledgeRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrBusinessKnowledgeNotFound
	}
	return u.businessKnowledgeRepo.Update(ctx, entity)
}

func (u *DataUsecase) DeleteBusinessKnowledge(ctx context.Context, id string) error {
	existing, err := u.businessKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrBusinessKnowledgeNotFound
	}
	return u.businessKnowledgeRepo.Delete(ctx, id)
}

func (u *DataUsecase) RetryBusinessKnowledgeEmbedding(ctx context.Context, id string) error {
	existing, err := u.businessKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrBusinessKnowledgeNotFound
	}
	existing.EmbeddingStatus = "pending"
	existing.ErrorMsg = ""
	return u.businessKnowledgeRepo.Update(ctx, existing)
}

func (u *DataUsecase) ToggleBusinessKnowledgeRecall(ctx context.Context, id string, isRecall bool) error {
	existing, err := u.businessKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrBusinessKnowledgeNotFound
	}
	val := 0
	if isRecall {
		val = 1
	}
	return u.businessKnowledgeRepo.UpdateRecall(ctx, id, val)
}

func (u *DataUsecase) ToggleBusinessKnowledgeRecallOn(ctx context.Context, id string) error {
	existing, err := u.businessKnowledgeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrBusinessKnowledgeNotFound
	}
	val := 0
	if existing.IsRecall == 1 {
		val = 0
	} else {
		val = 1
	}
	return u.businessKnowledgeRepo.UpdateRecall(ctx, id, val)
}

func (u *DataUsecase) RefreshBusinessKnowledgeVectorStore(ctx context.Context, agentID int64) (int64, error) {
	return u.businessKnowledgeRepo.BatchResetEmbedding(ctx, agentID)
}
