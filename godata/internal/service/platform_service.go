package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/usecase"
)

// PlatformService is a thin pass-through wrapper around PlatformUsecase.
type PlatformService struct {
	uc *usecase.PlatformUsecase
}

// NewPlatformService creates a new PlatformService with the given usecase.
func NewPlatformService(uc *usecase.PlatformUsecase) *PlatformService {
	return &PlatformService{uc: uc}
}

// ──────────────────────────── Auth / Login ────────────────────────────

func (s *PlatformService) Login(ctx context.Context, dto model.AccountLoginDTO) (*model.AccountInfo, error) {
	return s.uc.Login(ctx, dto)
}

func (s *PlatformService) ThirdPartyLogin(ctx context.Context, dto model.ThirdPartyLoginDTO) (*model.AccountInfo, error) {
	return s.uc.ThirdPartyLogin(ctx, dto)
}

func (s *PlatformService) UpdatePassword(ctx context.Context, dto model.UpdatePwdDTO) error {
	return s.uc.UpdatePassword(ctx, dto)
}

// ──────────────────────────── GroupInfo ────────────────────────────

func (s *PlatformService) CreateGroupInfo(ctx context.Context, entity *model.GroupInfo) (*model.GroupInfo, error) {
	return s.uc.CreateGroupInfo(ctx, entity)
}

func (s *PlatformService) UpdateGroupInfo(ctx context.Context, entity *model.GroupInfo) error {
	return s.uc.UpdateGroupInfo(ctx, entity)
}

func (s *PlatformService) DeleteGroupInfo(ctx context.Context, id string) error {
	return s.uc.DeleteGroupInfo(ctx, id)
}

func (s *PlatformService) GetGroupInfoByID(ctx context.Context, id string) (*model.GroupInfo, error) {
	return s.uc.GetGroupInfoByID(ctx, id)
}

func (s *PlatformService) PageGroupInfo(ctx context.Context, page, size int, query *model.GroupInfo) ([]*model.GroupInfo, int64, error) {
	return s.uc.PageGroupInfo(ctx, page, size, query)
}

func (s *PlatformService) ListGroupInfo(ctx context.Context) ([]*model.GroupInfo, error) {
	return s.uc.ListGroupInfo(ctx)
}

func (s *PlatformService) ToggleGroupInfoStatus(ctx context.Context, id string) error {
	return s.uc.ToggleGroupInfoStatus(ctx, id)
}

func (s *PlatformService) GetEnabledGroupInfo(ctx context.Context) ([]*model.GroupInfo, error) {
	return s.uc.GetEnabledGroupInfo(ctx)
}

// ──────────────────────────── GroupAgentInfo ────────────────────────────

func (s *PlatformService) CreateGroupAgentInfo(ctx context.Context, entity *model.GroupAgentInfo) (*model.GroupAgentInfo, error) {
	return s.uc.CreateGroupAgentInfo(ctx, entity)
}

func (s *PlatformService) UpdateGroupAgentInfo(ctx context.Context, entity *model.GroupAgentInfo) error {
	return s.uc.UpdateGroupAgentInfo(ctx, entity)
}

func (s *PlatformService) DeleteGroupAgentInfo(ctx context.Context, id string) error {
	return s.uc.DeleteGroupAgentInfo(ctx, id)
}

func (s *PlatformService) GetGroupAgentInfoByID(ctx context.Context, id string) (*model.GroupAgentInfo, error) {
	return s.uc.GetGroupAgentInfoByID(ctx, id)
}

func (s *PlatformService) ListGroupAgentInfo(ctx context.Context) ([]*model.GroupAgentInfo, error) {
	return s.uc.ListGroupAgentInfo(ctx)
}

func (s *PlatformService) PageGroupAgentInfo(ctx context.Context, page, size int, query *model.GroupAgentInfo) ([]*model.GroupAgentInfo, int64, error) {
	return s.uc.PageGroupAgentInfo(ctx, page, size, query)
}

func (s *PlatformService) FindGroupAgentInfoByGroupID(ctx context.Context, groupID string) ([]*model.GroupAgentInfo, error) {
	return s.uc.FindGroupAgentInfoByGroupID(ctx, groupID)
}

func (s *PlatformService) FindGroupAgentInfoByAgentID(ctx context.Context, agentID int64) ([]*model.GroupAgentInfo, error) {
	return s.uc.FindGroupAgentInfoByAgentID(ctx, agentID)
}

func (s *PlatformService) GetGroupAgentInfoByGroupIdAndAgentId(ctx context.Context, groupID string, agentID int64) (*model.GroupAgentInfo, error) {
	return s.uc.GetGroupAgentInfoByGroupIdAndAgentId(ctx, groupID, agentID)
}

// ──────────────────────────── AccountInfo ────────────────────────────

func (s *PlatformService) CreateAccountInfo(ctx context.Context, entity *model.AccountInfo) (*model.AccountInfo, error) {
	return s.uc.CreateAccountInfo(ctx, entity)
}

func (s *PlatformService) UpdateAccountInfo(ctx context.Context, entity *model.AccountInfo) error {
	return s.uc.UpdateAccountInfo(ctx, entity)
}

func (s *PlatformService) DeleteAccountInfo(ctx context.Context, id string) error {
	return s.uc.DeleteAccountInfo(ctx, id)
}

func (s *PlatformService) GetAccountInfoByID(ctx context.Context, id string) (*model.AccountInfo, error) {
	return s.uc.GetAccountInfoByID(ctx, id)
}

func (s *PlatformService) GetAccountInfoByUsername(ctx context.Context, username string) (*model.AccountInfo, error) {
	return s.uc.GetAccountInfoByUsername(ctx, username)
}

func (s *PlatformService) GetAccountInfoByThirdPartyID(ctx context.Context, thirdPartyID string) (*model.AccountInfo, error) {
	return s.uc.GetAccountInfoByThirdPartyID(ctx, thirdPartyID)
}

func (s *PlatformService) GetAccountInfoByCode(ctx context.Context, code string) (*model.AccountInfo, error) {
	return s.uc.GetAccountInfoByCode(ctx, code)
}

func (s *PlatformService) GetAccountInfoByStatus(ctx context.Context, status string) ([]*model.AccountInfo, error) {
	return s.uc.GetAccountInfoByStatus(ctx, status)
}

func (s *PlatformService) PageAccountInfo(ctx context.Context, page, size int, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	return s.uc.PageAccountInfo(ctx, page, size, query)
}

func (s *PlatformService) GetUnGroupPage(ctx context.Context, page, size int, groupID string, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	return s.uc.GetUnGroupPage(ctx, page, size, groupID, query)
}

func (s *PlatformService) GetMyAgents(ctx context.Context, accountID string) ([]*model.GroupAgentInfo, error) {
	return s.uc.GetMyAgents(ctx, accountID)
}

func (s *PlatformService) ListAccountInfo(ctx context.Context) ([]*model.AccountInfo, error) {
	return s.uc.ListAccountInfo(ctx)
}

func (s *PlatformService) BatchUpdateAccountStatus(ctx context.Context, ids []string, status string) error {
	return s.uc.BatchUpdateAccountStatus(ctx, ids, status)
}

// ──────────────────────────── AccountGroupInfo ────────────────────────────

func (s *PlatformService) CreateAccountGroupInfo(ctx context.Context, entity *model.AccountGroupInfo) (*model.AccountGroupInfo, error) {
	return s.uc.CreateAccountGroupInfo(ctx, entity)
}

func (s *PlatformService) UpdateAccountGroupInfo(ctx context.Context, entity *model.AccountGroupInfo) error {
	return s.uc.UpdateAccountGroupInfo(ctx, entity)
}

func (s *PlatformService) DeleteAccountGroupInfo(ctx context.Context, id string) error {
	return s.uc.DeleteAccountGroupInfo(ctx, id)
}

func (s *PlatformService) GetAccountGroupInfoByID(ctx context.Context, id string) (*model.AccountGroupInfo, error) {
	return s.uc.GetAccountGroupInfoByID(ctx, id)
}

func (s *PlatformService) FindAccountGroupInfoByGroupID(ctx context.Context, groupID string) ([]*model.AccountGroupInfo, error) {
	return s.uc.FindAccountGroupInfoByGroupID(ctx, groupID)
}

func (s *PlatformService) FindAccountGroupInfoByAccountID(ctx context.Context, accountID string) ([]*model.AccountGroupInfo, error) {
	return s.uc.FindAccountGroupInfoByAccountID(ctx, accountID)
}

func (s *PlatformService) PageAccountGroupInfo(ctx context.Context, page, size int, query *model.AccountGroupInfo) ([]*model.AccountGroupInfo, int64, error) {
	return s.uc.PageAccountGroupInfo(ctx, page, size, query)
}

// ──────────────────────────── AccountTenantInfo ────────────────────────────

func (s *PlatformService) CreateAccountTenantInfo(ctx context.Context, entity *model.AccountTenantInfo) (*model.AccountTenantInfo, error) {
	return s.uc.CreateAccountTenantInfo(ctx, entity)
}

func (s *PlatformService) UpdateAccountTenantInfo(ctx context.Context, entity *model.AccountTenantInfo) error {
	return s.uc.UpdateAccountTenantInfo(ctx, entity)
}

func (s *PlatformService) DeleteAccountTenantInfo(ctx context.Context, id string) error {
	return s.uc.DeleteAccountTenantInfo(ctx, id)
}

func (s *PlatformService) GetAccountTenantInfoByID(ctx context.Context, id string) (*model.AccountTenantInfo, error) {
	return s.uc.GetAccountTenantInfoByID(ctx, id)
}

func (s *PlatformService) FindAccountTenantInfoByAccountID(ctx context.Context, accountID string) ([]*model.AccountTenantInfo, error) {
	return s.uc.FindAccountTenantInfoByAccountID(ctx, accountID)
}

func (s *PlatformService) FindAccountTenantInfoByTenantID(ctx context.Context, tenantID string) ([]*model.AccountTenantInfo, error) {
	return s.uc.FindAccountTenantInfoByTenantID(ctx, tenantID)
}

func (s *PlatformService) PageAccountTenantInfo(ctx context.Context, page, size int, query *model.AccountTenantInfo) ([]*model.AccountTenantInfo, int64, error) {
	return s.uc.PageAccountTenantInfo(ctx, page, size, query)
}

// ──────────────────────────── TenantInfo ────────────────────────────

func (s *PlatformService) CreateTenantInfo(ctx context.Context, entity *model.TenantInfo) (*model.TenantInfo, error) {
	return s.uc.CreateTenantInfo(ctx, entity)
}

func (s *PlatformService) UpdateTenantInfo(ctx context.Context, entity *model.TenantInfo) error {
	return s.uc.UpdateTenantInfo(ctx, entity)
}

func (s *PlatformService) DeleteTenantInfo(ctx context.Context, id string) error {
	return s.uc.DeleteTenantInfo(ctx, id)
}

func (s *PlatformService) GetTenantInfoByID(ctx context.Context, id string) (*model.TenantInfo, error) {
	return s.uc.GetTenantInfoByID(ctx, id)
}

func (s *PlatformService) GetTenantInfoBySN(ctx context.Context, sn string) (*model.TenantInfo, error) {
	return s.uc.GetTenantInfoBySN(ctx, sn)
}

func (s *PlatformService) PageTenantInfo(ctx context.Context, page, size int, query *model.TenantInfo) ([]*model.TenantInfo, int64, error) {
	return s.uc.PageTenantInfo(ctx, page, size, query)
}

// ──────────────────────────── PlatformInfo ────────────────────────────

func (s *PlatformService) CreatePlatformInfo(ctx context.Context, entity *model.PlatformInfo) (*model.PlatformInfo, error) {
	return s.uc.CreatePlatformInfo(ctx, entity)
}

func (s *PlatformService) UpdatePlatformInfo(ctx context.Context, entity *model.PlatformInfo) error {
	return s.uc.UpdatePlatformInfo(ctx, entity)
}

func (s *PlatformService) DeletePlatformInfo(ctx context.Context, id string) error {
	return s.uc.DeletePlatformInfo(ctx, id)
}

func (s *PlatformService) GetPlatformInfoByID(ctx context.Context, id string) (*model.PlatformInfo, error) {
	return s.uc.GetPlatformInfoByID(ctx, id)
}

func (s *PlatformService) FindPlatformInfoByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	return s.uc.FindPlatformInfoByType(ctx, typ)
}

func (s *PlatformService) FindPlatformInfoEnabledByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	return s.uc.FindPlatformInfoEnabledByType(ctx, typ)
}

func (s *PlatformService) PagePlatformInfo(ctx context.Context, page, size int, query *model.PlatformInfo) ([]*model.PlatformInfo, int64, error) {
	return s.uc.PagePlatformInfo(ctx, page, size, query)
}

func (s *PlatformService) TogglePlatformInfoStatus(ctx context.Context, id string) error {
	return s.uc.TogglePlatformInfoStatus(ctx, id)
}
