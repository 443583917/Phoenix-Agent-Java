package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"github.com/phoenix-agent-go/infra/id"
)

// ──────────────────────────── Data Errors ────────────────────────────

var (
	ErrAgentNotFound             = &AppError{Code: 601001, Msg: "智能体不存在"}
	ErrAgentCategoryNotFound     = &AppError{Code: 602001, Msg: "分类不存在"}
	ErrAgentDatasourceNotFound   = &AppError{Code: 603001, Msg: "智能体数据源关联不存在"}
	ErrAgentKnowledgeNotFound    = &AppError{Code: 604001, Msg: "知识库条目不存在"}
	ErrAgentPresetQuestionNotFound = &AppError{Code: 605001, Msg: "预设问题不存在"}
)

// ──────────────────────────── DataUsecase ────────────────────────────

type DataUsecase struct {
	agentRepo               repository.AgentRepository
	agentCategoryRepo       repository.AgentCategoryRepository
	agentDatasourceRepo     repository.AgentDatasourceRepository
	agentKnowledgeRepo      repository.AgentKnowledgeRepository
	agentPresetQuestionRepo repository.AgentPresetQuestionRepository
	agentDatasourceTablesRepo repository.AgentDatasourceTablesRepository
}

func NewDataUsecase(
	agentRepo repository.AgentRepository,
	agentCategoryRepo repository.AgentCategoryRepository,
	agentDatasourceRepo repository.AgentDatasourceRepository,
	agentKnowledgeRepo repository.AgentKnowledgeRepository,
	agentPresetQuestionRepo repository.AgentPresetQuestionRepository,
	agentDatasourceTablesRepo repository.AgentDatasourceTablesRepository,
) *DataUsecase {
	return &DataUsecase{
		agentRepo:               agentRepo,
		agentCategoryRepo:       agentCategoryRepo,
		agentDatasourceRepo:     agentDatasourceRepo,
		agentKnowledgeRepo:      agentKnowledgeRepo,
		agentPresetQuestionRepo: agentPresetQuestionRepo,
		agentDatasourceTablesRepo: agentDatasourceTablesRepo,
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
