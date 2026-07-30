# Phoenix Go Phase 3 — Platform Management Plan

**Goal:** Migrate phoenix-platform + phoenix-common (7 models, 8 CRUD controllers, Auth login)

**Architecture:** Same DDD pattern as Phase 2: Entity → Repo → Usecase → Service → Handler

**Key difference from Phase 2:** Platform uses `BaseModel` (createTime/creator/updateTime/updator/delFlag) NOT `BaseEntity` (id/createBy/updateBy). Each entity declares its own `@Id String`.

---

### Task 1: 7 Platform Entities + BaseModel

File: `internal/model/platform_entity.go`

7 GORM entities with TableName():
- GroupInfo — tbl_platform_group_info
- GroupAgentInfo — tbl_platform_group_agent_info  
- AccountInfo — tbl_platform_account_info
- AccountGroupInfo — tbl_platform_account_group_info
- AccountTenantInfo — tbl_platform_account_tenant_info
- TenantInfo — tbl_platform_tenant_info

Plus: `internal/model/common_entity.go` — PlatformInfo — tbl_platform_platform_info

Each has its own `ID string` PK. BaseModel fields embedded.

### Task 2: Platform DTOs + Repos + Usecase + Service

Single task covering all layers since it's pure CRUD:
- DTOs: AccountLoginDTO, UpdatePwdDTO, ThirdPartyLoginDTO
- Repos: 7 GORM CRUD repos
- Usecase: Account login (MD5 same as privilege), CRUD for all
- Service: thin wrapper

### Task 3: 8 Platform Handlers + Routes

8 handler files + register in router.go:

**`/auth`** (no prefix — different from privilege auth):
- POST /login, POST /logout, POST /thirdLogin, PUT /updatePassword

**`/platform/group-info`** — CRUD + toggle-status + removeAgent
**`/platform/group-agent-info`** — CRUD + list  
**`/platform/account-info`** — CRUD + list + batch-status + getMyAgents
**`/platform/account-group-info`** — CRUD
**`/platform/account-tenant-info`** — CRUD
**`/platform/platform-info`** — CRUD + type queries + enabled
**`/platform/sync`** — sync stubs (POST all/departments/users/depts)

### Task 4: Build + Integration Verify
