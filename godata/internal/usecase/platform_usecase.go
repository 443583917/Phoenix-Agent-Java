package usecase

import (
	"context"
	"strconv"

	"github.com/phoenix-agent-go/internal/domain/privilege"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"github.com/phoenix-agent-go/infra/id"
)

// ──────────────────────────── Platform Errors ────────────────────────────

var (
	ErrPlatformAccountNotFound     = &AppError{Code: 501001, Msg: "账号不存在"}
	ErrPlatformInvalidCredentials  = &AppError{Code: 501002, Msg: "用户名或密码错误"}
	ErrPlatformAccountDisabled     = &AppError{Code: 501003, Msg: "账号已被禁用"}
	ErrPlatformPasswordWrong       = &AppError{Code: 501004, Msg: "密码错误"}
	ErrPlatformOldPasswordWrong    = &AppError{Code: 501005, Msg: "原密码错误"}
	ErrPlatformGroupNotFound       = &AppError{Code: 502001, Msg: "组织不存在"}
	ErrPlatformTenantNotFound      = &AppError{Code: 503001, Msg: "租户不存在"}
	ErrPlatformPlatformNotFound    = &AppError{Code: 504001, Msg: "三方平台不存在"}
	ErrPlatformGroupAgentNotFound  = &AppError{Code: 505001, Msg: "组织智能体关联不存在"}
	ErrPlatformAccountGroupNotFound = &AppError{Code: 506001, Msg: "账号组织关联不存在"}
	ErrPlatformAccountTenantNotFound = &AppError{Code: 507001, Msg: "账号租户关联不存在"}
)

// ──────────────────────────── PlatformUsecase ────────────────────────────

// PlatformUsecase orchestrates platform domain operations.
type PlatformUsecase struct {
	groupInfoRepo         repository.GroupInfoRepository
	groupAgentInfoRepo    repository.GroupAgentInfoRepository
	accountInfoRepo       repository.AccountInfoRepository
	accountGroupInfoRepo  repository.AccountGroupInfoRepository
	accountTenantInfoRepo repository.AccountTenantInfoRepository
	tenantInfoRepo        repository.TenantInfoRepository
	platformInfoRepo      repository.PlatformInfoRepository
}

// NewPlatformUsecase constructs a PlatformUsecase with all required repositories.
func NewPlatformUsecase(
	groupInfoRepo repository.GroupInfoRepository,
	groupAgentInfoRepo repository.GroupAgentInfoRepository,
	accountInfoRepo repository.AccountInfoRepository,
	accountGroupInfoRepo repository.AccountGroupInfoRepository,
	accountTenantInfoRepo repository.AccountTenantInfoRepository,
	tenantInfoRepo repository.TenantInfoRepository,
	platformInfoRepo repository.PlatformInfoRepository,
) *PlatformUsecase {
	return &PlatformUsecase{
		groupInfoRepo:         groupInfoRepo,
		groupAgentInfoRepo:    groupAgentInfoRepo,
		accountInfoRepo:       accountInfoRepo,
		accountGroupInfoRepo:  accountGroupInfoRepo,
		accountTenantInfoRepo: accountTenantInfoRepo,
		tenantInfoRepo:        tenantInfoRepo,
		platformInfoRepo:      platformInfoRepo,
	}
}

// ──────────────────────────── helpers ────────────────────────────

func platformGenID() string {
	return strconv.FormatUint(id.MustGenerateID(), 10)
}

// ──────────────────────────── Auth / Login ────────────────────────────

// Login authenticates an account by username + password (MD5).
func (u *PlatformUsecase) Login(ctx context.Context, dto model.AccountLoginDTO) (*model.AccountInfo, error) {
	account, err := u.accountInfoRepo.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrPlatformInvalidCredentials
	}
	if account.Status == "1" {
		return nil, ErrPlatformAccountDisabled
	}
	if !privilege.CheckPassword(account.Password, dto.Password) {
		return nil, ErrPlatformPasswordWrong
	}
	return account, nil
}

// ThirdPartyLogin finds an account by third-party ID.
func (u *PlatformUsecase) ThirdPartyLogin(ctx context.Context, dto model.ThirdPartyLoginDTO) (*model.AccountInfo, error) {
	account, err := u.accountInfoRepo.FindByThirdPartyID(ctx, dto.ThirdPartyID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrPlatformAccountNotFound
	}
	if account.Status == "1" {
		return nil, ErrPlatformAccountDisabled
	}
	return account, nil
}

// UpdatePassword changes the password after verifying the old one.
func (u *PlatformUsecase) UpdatePassword(ctx context.Context, dto model.UpdatePwdDTO) error {
	account, err := u.accountInfoRepo.FindByID(ctx, dto.UserID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrPlatformAccountNotFound
	}
	if !privilege.CheckPassword(account.Password, dto.OldPassword) {
		return ErrPlatformOldPasswordWrong
	}
	hashed, _ := privilege.HashPassword(dto.NewPassword)
	account.Password = hashed
	return u.accountInfoRepo.Update(ctx, account)
}

// ──────────────────────────── GroupInfo ────────────────────────────

func (u *PlatformUsecase) CreateGroupInfo(ctx context.Context, entity *model.GroupInfo) (*model.GroupInfo, error) {
	entity.ID = platformGenID()
	if err := u.groupInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateGroupInfo(ctx context.Context, entity *model.GroupInfo) error {
	existing, err := u.groupInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformGroupNotFound
	}
	return u.groupInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteGroupInfo(ctx context.Context, id string) error {
	existing, err := u.groupInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformGroupNotFound
	}
	return u.groupInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetGroupInfoByID(ctx context.Context, id string) (*model.GroupInfo, error) {
	entity, err := u.groupInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformGroupNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) PageGroupInfo(ctx context.Context, page, size int, query *model.GroupInfo) ([]*model.GroupInfo, int64, error) {
	return u.groupInfoRepo.Page(ctx, page, size, query)
}

func (u *PlatformUsecase) ListGroupInfo(ctx context.Context) ([]*model.GroupInfo, error) {
	return u.groupInfoRepo.List(ctx)
}

func (u *PlatformUsecase) ToggleGroupInfoStatus(ctx context.Context, id string) error {
	return u.groupInfoRepo.ToggleStatus(ctx, id)
}

func (u *PlatformUsecase) GetEnabledGroupInfo(ctx context.Context) ([]*model.GroupInfo, error) {
	return u.groupInfoRepo.GetEnabled(ctx)
}

// ──────────────────────────── GroupAgentInfo ────────────────────────────

func (u *PlatformUsecase) CreateGroupAgentInfo(ctx context.Context, entity *model.GroupAgentInfo) (*model.GroupAgentInfo, error) {
	entity.ID = platformGenID()
	if err := u.groupAgentInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateGroupAgentInfo(ctx context.Context, entity *model.GroupAgentInfo) error {
	existing, err := u.groupAgentInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformGroupAgentNotFound
	}
	return u.groupAgentInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteGroupAgentInfo(ctx context.Context, id string) error {
	existing, err := u.groupAgentInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformGroupAgentNotFound
	}
	return u.groupAgentInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetGroupAgentInfoByID(ctx context.Context, id string) (*model.GroupAgentInfo, error) {
	entity, err := u.groupAgentInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformGroupAgentNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) ListGroupAgentInfo(ctx context.Context) ([]*model.GroupAgentInfo, error) {
	return u.groupAgentInfoRepo.List(ctx)
}

func (u *PlatformUsecase) FindGroupAgentInfoByGroupID(ctx context.Context, groupID string) ([]*model.GroupAgentInfo, error) {
	return u.groupAgentInfoRepo.FindByGroupID(ctx, groupID)
}

func (u *PlatformUsecase) FindGroupAgentInfoByAgentID(ctx context.Context, agentID int64) ([]*model.GroupAgentInfo, error) {
	return u.groupAgentInfoRepo.FindByAgentID(ctx, agentID)
}

func (u *PlatformUsecase) GetGroupAgentInfoByGroupIdAndAgentId(ctx context.Context, groupID string, agentID int64) (*model.GroupAgentInfo, error) {
	entity, err := u.groupAgentInfoRepo.GetByGroupIdAndAgentId(ctx, groupID, agentID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformGroupAgentNotFound
	}
	return entity, nil
}

// ──────────────────────────── AccountInfo ────────────────────────────

func (u *PlatformUsecase) CreateAccountInfo(ctx context.Context, entity *model.AccountInfo) (*model.AccountInfo, error) {
	entity.ID = platformGenID()
	if entity.Password != "" {
		hashed, _ := privilege.HashPassword(entity.Password)
		entity.Password = hashed
	}
	if err := u.accountInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateAccountInfo(ctx context.Context, entity *model.AccountInfo) error {
	existing, err := u.accountInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountNotFound
	}
	// Preserve password if not provided in update.
	if entity.Password == "" {
		entity.Password = existing.Password
	}
	return u.accountInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteAccountInfo(ctx context.Context, id string) error {
	existing, err := u.accountInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountNotFound
	}
	return u.accountInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetAccountInfoByID(ctx context.Context, id string) (*model.AccountInfo, error) {
	entity, err := u.accountInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) GetAccountInfoByUsername(ctx context.Context, username string) (*model.AccountInfo, error) {
	entity, err := u.accountInfoRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) GetAccountInfoByThirdPartyID(ctx context.Context, thirdPartyID string) (*model.AccountInfo, error) {
	entity, err := u.accountInfoRepo.FindByThirdPartyID(ctx, thirdPartyID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) GetAccountInfoByCode(ctx context.Context, code string) (*model.AccountInfo, error) {
	entity, err := u.accountInfoRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) GetAccountInfoByStatus(ctx context.Context, status string) ([]*model.AccountInfo, error) {
	return u.accountInfoRepo.FindByStatus(ctx, status)
}

func (u *PlatformUsecase) PageAccountInfo(ctx context.Context, page, size int, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	return u.accountInfoRepo.Page(ctx, page, size, query)
}

func (u *PlatformUsecase) GetUnGroupPage(ctx context.Context, page, size int, groupID string, query *model.AccountInfo) ([]*model.AccountInfo, int64, error) {
	return u.accountInfoRepo.GetUnGroupPage(ctx, page, size, groupID, query)
}

func (u *PlatformUsecase) GetMyAgents(ctx context.Context, accountID string) ([]*model.GroupAgentInfo, error) {
	return u.accountInfoRepo.GetMyAgents(ctx, accountID)
}

func (u *PlatformUsecase) ListAccountInfo(ctx context.Context) ([]*model.AccountInfo, error) {
	return u.accountInfoRepo.List(ctx)
}

func (u *PlatformUsecase) BatchUpdateAccountStatus(ctx context.Context, ids []string, status string) error {
	return u.accountInfoRepo.BatchUpdateStatus(ctx, ids, status)
}

// ──────────────────────────── AccountGroupInfo ────────────────────────────

func (u *PlatformUsecase) CreateAccountGroupInfo(ctx context.Context, entity *model.AccountGroupInfo) (*model.AccountGroupInfo, error) {
	entity.ID = platformGenID()
	if err := u.accountGroupInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateAccountGroupInfo(ctx context.Context, entity *model.AccountGroupInfo) error {
	existing, err := u.accountGroupInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountGroupNotFound
	}
	return u.accountGroupInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteAccountGroupInfo(ctx context.Context, id string) error {
	existing, err := u.accountGroupInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountGroupNotFound
	}
	return u.accountGroupInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetAccountGroupInfoByID(ctx context.Context, id string) (*model.AccountGroupInfo, error) {
	entity, err := u.accountGroupInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountGroupNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) FindAccountGroupInfoByGroupID(ctx context.Context, groupID string) ([]*model.AccountGroupInfo, error) {
	return u.accountGroupInfoRepo.FindByGroupID(ctx, groupID)
}

func (u *PlatformUsecase) FindAccountGroupInfoByAccountID(ctx context.Context, accountID string) ([]*model.AccountGroupInfo, error) {
	return u.accountGroupInfoRepo.FindByAccountID(ctx, accountID)
}

func (u *PlatformUsecase) PageAccountGroupInfo(ctx context.Context, page, size int, query *model.AccountGroupInfo) ([]*model.AccountGroupInfo, int64, error) {
	return u.accountGroupInfoRepo.Page(ctx, page, size, query)
}

// ──────────────────────────── AccountTenantInfo ────────────────────────────

func (u *PlatformUsecase) CreateAccountTenantInfo(ctx context.Context, entity *model.AccountTenantInfo) (*model.AccountTenantInfo, error) {
	entity.ID = platformGenID()
	if err := u.accountTenantInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateAccountTenantInfo(ctx context.Context, entity *model.AccountTenantInfo) error {
	existing, err := u.accountTenantInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountTenantNotFound
	}
	return u.accountTenantInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteAccountTenantInfo(ctx context.Context, id string) error {
	existing, err := u.accountTenantInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformAccountTenantNotFound
	}
	return u.accountTenantInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetAccountTenantInfoByID(ctx context.Context, id string) (*model.AccountTenantInfo, error) {
	entity, err := u.accountTenantInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformAccountTenantNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) FindAccountTenantInfoByAccountID(ctx context.Context, accountID string) ([]*model.AccountTenantInfo, error) {
	return u.accountTenantInfoRepo.FindByAccountID(ctx, accountID)
}

func (u *PlatformUsecase) FindAccountTenantInfoByTenantID(ctx context.Context, tenantID string) ([]*model.AccountTenantInfo, error) {
	return u.accountTenantInfoRepo.FindByTenantID(ctx, tenantID)
}

func (u *PlatformUsecase) PageAccountTenantInfo(ctx context.Context, page, size int, query *model.AccountTenantInfo) ([]*model.AccountTenantInfo, int64, error) {
	return u.accountTenantInfoRepo.Page(ctx, page, size, query)
}

// ──────────────────────────── TenantInfo ────────────────────────────

func (u *PlatformUsecase) CreateTenantInfo(ctx context.Context, entity *model.TenantInfo) (*model.TenantInfo, error) {
	entity.ID = platformGenID()
	if err := u.tenantInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdateTenantInfo(ctx context.Context, entity *model.TenantInfo) error {
	existing, err := u.tenantInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformTenantNotFound
	}
	return u.tenantInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeleteTenantInfo(ctx context.Context, id string) error {
	existing, err := u.tenantInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformTenantNotFound
	}
	return u.tenantInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetTenantInfoByID(ctx context.Context, id string) (*model.TenantInfo, error) {
	entity, err := u.tenantInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformTenantNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) GetTenantInfoBySN(ctx context.Context, sn string) (*model.TenantInfo, error) {
	entity, err := u.tenantInfoRepo.FindBySN(ctx, sn)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformTenantNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) PageTenantInfo(ctx context.Context, page, size int, query *model.TenantInfo) ([]*model.TenantInfo, int64, error) {
	return u.tenantInfoRepo.Page(ctx, page, size, query)
}

// ──────────────────────────── PlatformInfo ────────────────────────────

func (u *PlatformUsecase) CreatePlatformInfo(ctx context.Context, entity *model.PlatformInfo) (*model.PlatformInfo, error) {
	entity.ID = platformGenID()
	if err := u.platformInfoRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *PlatformUsecase) UpdatePlatformInfo(ctx context.Context, entity *model.PlatformInfo) error {
	existing, err := u.platformInfoRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformPlatformNotFound
	}
	return u.platformInfoRepo.Update(ctx, entity)
}

func (u *PlatformUsecase) DeletePlatformInfo(ctx context.Context, id string) error {
	existing, err := u.platformInfoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPlatformPlatformNotFound
	}
	return u.platformInfoRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) GetPlatformInfoByID(ctx context.Context, id string) (*model.PlatformInfo, error) {
	entity, err := u.platformInfoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrPlatformPlatformNotFound
	}
	return entity, nil
}

func (u *PlatformUsecase) FindPlatformInfoByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	return u.platformInfoRepo.FindByType(ctx, typ)
}

func (u *PlatformUsecase) FindPlatformInfoEnabledByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error) {
	return u.platformInfoRepo.FindEnabledByType(ctx, typ)
}

func (u *PlatformUsecase) PagePlatformInfo(ctx context.Context, page, size int, query *model.PlatformInfo) ([]*model.PlatformInfo, int64, error) {
	return u.platformInfoRepo.Page(ctx, page, size, query)
}

func (u *PlatformUsecase) TogglePlatformInfoStatus(ctx context.Context, id string) error {
	return u.platformInfoRepo.ToggleStatus(ctx, id)
}
