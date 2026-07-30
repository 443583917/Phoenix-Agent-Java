package repository

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
)

// ──────────────────────────── GroupInfo ────────────────────────────

type GroupInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.GroupInfo, error)
	Page(ctx context.Context, page, size int, query *model.GroupInfo) ([]*model.GroupInfo, int64, error)
	List(ctx context.Context) ([]*model.GroupInfo, error)
	Create(ctx context.Context, group *model.GroupInfo) error
	Update(ctx context.Context, group *model.GroupInfo) error
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) error
	GetEnabled(ctx context.Context) ([]*model.GroupInfo, error)
}

// ──────────────────────────── GroupAgentInfo ────────────────────────────

type GroupAgentInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.GroupAgentInfo, error)
	FindByGroupID(ctx context.Context, groupID string) ([]*model.GroupAgentInfo, error)
	FindByAgentID(ctx context.Context, agentID int64) ([]*model.GroupAgentInfo, error)
	GetByGroupIdAndAgentId(ctx context.Context, groupID string, agentID int64) (*model.GroupAgentInfo, error)
	List(ctx context.Context) ([]*model.GroupAgentInfo, error)
	Create(ctx context.Context, ga *model.GroupAgentInfo) error
	Update(ctx context.Context, ga *model.GroupAgentInfo) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AccountInfo ────────────────────────────

type AccountInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.AccountInfo, error)
	FindByUsername(ctx context.Context, username string) (*model.AccountInfo, error)
	FindByThirdPartyID(ctx context.Context, thirdPartyID string) (*model.AccountInfo, error)
	FindByCode(ctx context.Context, code string) (*model.AccountInfo, error)
	FindByStatus(ctx context.Context, status string) ([]*model.AccountInfo, error)
	Page(ctx context.Context, page, size int, query *model.AccountInfo) ([]*model.AccountInfo, int64, error)
	GetUnGroupPage(ctx context.Context, page, size int, groupID string, query *model.AccountInfo) ([]*model.AccountInfo, int64, error)
	GetMyAgents(ctx context.Context, accountID string) ([]*model.GroupAgentInfo, error)
	List(ctx context.Context) ([]*model.AccountInfo, error)
	Create(ctx context.Context, account *model.AccountInfo) error
	Update(ctx context.Context, account *model.AccountInfo) error
	Delete(ctx context.Context, id string) error
	BatchUpdateStatus(ctx context.Context, ids []string, status string) error
}

// ──────────────────────────── AccountGroupInfo ────────────────────────────

type AccountGroupInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.AccountGroupInfo, error)
	FindByGroupID(ctx context.Context, groupID string) ([]*model.AccountGroupInfo, error)
	FindByAccountID(ctx context.Context, accountID string) ([]*model.AccountGroupInfo, error)
	Page(ctx context.Context, page, size int, query *model.AccountGroupInfo) ([]*model.AccountGroupInfo, int64, error)
	Create(ctx context.Context, ag *model.AccountGroupInfo) error
	Update(ctx context.Context, ag *model.AccountGroupInfo) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── AccountTenantInfo ────────────────────────────

type AccountTenantInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.AccountTenantInfo, error)
	FindByAccountID(ctx context.Context, accountID string) ([]*model.AccountTenantInfo, error)
	FindByTenantID(ctx context.Context, tenantID string) ([]*model.AccountTenantInfo, error)
	Page(ctx context.Context, page, size int, query *model.AccountTenantInfo) ([]*model.AccountTenantInfo, int64, error)
	Create(ctx context.Context, at *model.AccountTenantInfo) error
	Update(ctx context.Context, at *model.AccountTenantInfo) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── TenantInfo ────────────────────────────

type TenantInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.TenantInfo, error)
	FindBySN(ctx context.Context, sn string) (*model.TenantInfo, error)
	Page(ctx context.Context, page, size int, query *model.TenantInfo) ([]*model.TenantInfo, int64, error)
	Create(ctx context.Context, tenant *model.TenantInfo) error
	Update(ctx context.Context, tenant *model.TenantInfo) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── PlatformInfo ────────────────────────────

type PlatformInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.PlatformInfo, error)
	FindByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error)
	FindEnabledByType(ctx context.Context, typ string) ([]*model.PlatformInfo, error)
	Page(ctx context.Context, page, size int, query *model.PlatformInfo) ([]*model.PlatformInfo, int64, error)
	Create(ctx context.Context, platform *model.PlatformInfo) error
	Update(ctx context.Context, platform *model.PlatformInfo) error
	Delete(ctx context.Context, id string) error
	ToggleStatus(ctx context.Context, id string) error
}
