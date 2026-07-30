# Phoenix Go 重写 Phase 2 — 权限认证模块实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 迁移 phoenix-privilege — 12 GORM Entity + Repository + Domain + Usecase + 12 REST Handler + JWT/Casbin 中间件

**Architecture:** GORM Entity → Repository interface → GORM impl → Domain logic → Usecase → Service → Handler

**Tech Stack:** GORM, golang-jwt v5, Casbin v2, bcrypt, go-redis, bigcache, Gin

## Global Constraints

- Module root: `github.com/phoenix-agent-go`
- Java `ReturnVo<T>` → `response.Success/Error(c, data)` (int code)
- Java `BaseEntity` → embedded `BaseModel` (ID string snowflake, CreateTime/UpdateTime/CreateBy/UpdateBy/DelFlag)
- Java `Page<T>` → `response.SuccessPage(c, data, total, page, size)`
- Java `Mono<T>` → Go plain synchronous handlers
- Java Sa-Token → Go JWT + Casbin (Casbin model uses tenant-aware RBAC)
- All table/column names match existing Java schema exactly
- API paths match Java Controller routes verbatim

---

### Task 1: BaseModel + 12 GORM Entity

**Files:**
- Create: `internal/model/base.go`
- Create: `internal/model/privilege_entity.go`

- [ ] **Step 1: BaseModel**

```go
// internal/model/base.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type BaseModel struct {
    ID         string         `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
    CreateTime time.Time      `gorm:"column:create_time;autoCreateTime" json:"createTime"`
    UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
    CreateBy   *string        `gorm:"column:create_by;type:varchar(64)" json:"createBy,omitempty"`
    UpdateBy   *string        `gorm:"column:update_by;type:varchar(64)" json:"updateBy,omitempty"`
    DelFlag    int            `gorm:"column:del_flag;default:0" json:"delFlag"`
}
```

- [ ] **Step 2: 12 GORM Entities**

```go
// internal/model/privilege_entity.go
package model

import "time"

// PrivilegeUser — tbl_privilege_user
type PrivilegeUser struct {
    BaseModel
    EmployeeID   string     `gorm:"column:employee_id;type:varchar(64)" json:"employeeId"`
    Code         string     `gorm:"column:code;type:varchar(64)" json:"code"`
    RealName     string     `gorm:"column:real_name;type:varchar(64)" json:"realName"`
    Username     string     `gorm:"column:username;type:varchar(64);uniqueIndex" json:"username"`
    Password     string     `gorm:"column:password;type:varchar(256)" json:"-"`
    Tel          string     `gorm:"column:tel;type:varchar(32)" json:"tel"`
    Phone        string     `gorm:"column:phone;type:varchar(32)" json:"phone"`
    Mobile       string     `gorm:"column:mobile;type:varchar(32)" json:"mobile"`
    Email        string     `gorm:"column:email;type:varchar(128)" json:"email"`
    Image        []byte     `gorm:"column:image;type:bytea" json:"image,omitempty"`
    CompanyID    string     `gorm:"column:company_id;type:varchar(64)" json:"companyId"`
    DeptID       string     `gorm:"column:dept_id;type:varchar(64)" json:"deptId"`
    ItUserID     string     `gorm:"column:it_user_id;type:varchar(64)" json:"itUserId"`
    ItUserName   string     `gorm:"column:it_user_name;type:varchar(64)" json:"itUserName"`
    IsLeader     int        `gorm:"column:is_leader;default:0" json:"isLeader"`
    Sex          int        `gorm:"column:sex;default:2" json:"sex"`
    Address      string     `gorm:"column:address;type:varchar(256)" json:"address"`
    Fax          string     `gorm:"column:fax;type:varchar(32)" json:"fax"`
    FailMonth    int        `gorm:"column:fail_month;default:0" json:"failMonth"`
    FailureTime  *time.Time `gorm:"column:failure_time" json:"failureTime,omitempty"`
    ACLTimestamp int        `gorm:"column:acl_timestamp;default:0" json:"aclTimestamp"`
    PwdFtime     *time.Time `gorm:"column:pwd_ftime" json:"pwdFtime,omitempty"`
    PwdInit      int        `gorm:"column:pwd_init;default:0" json:"pwdInit"`
    UserType     int        `gorm:"column:user_type;default:0" json:"userType"`
    Status       int        `gorm:"column:status;default:1" json:"status"`
}

func (PrivilegeUser) TableName() string { return "tbl_privilege_user" }

// PrivilegeRole — tbl_privilege_role
type PrivilegeRole struct {
    BaseModel
    Name       string `gorm:"column:name;type:varchar(64)" json:"name"`
    SN         string `gorm:"column:sn;type:varchar(64)" json:"sn"`
    RoleLevel  string `gorm:"column:role_level;type:varchar(32)" json:"roleLevel"`
    Note       string `gorm:"column:note;type:varchar(256)" json:"note"`
    ValidState int    `gorm:"column:valid_state;default:1" json:"validState"`
    CompanyID  int64  `gorm:"column:company_id" json:"companyId"`
    SystemID   string `gorm:"column:system_id;type:varchar(64)" json:"systemId"`
}

func (PrivilegeRole) TableName() string { return "tbl_privilege_role" }

// PrivilegeUserRole — tbl_privilege_user_role
type PrivilegeUserRole struct {
    BaseModel
    UserID string `gorm:"column:user_id;type:varchar(64);index" json:"userId"`
    RoleID string `gorm:"column:role_id;type:varchar(64);index" json:"roleId"`
}

func (PrivilegeUserRole) TableName() string { return "tbl_privilege_user_role" }

// PrivilegeModule — tbl_privilege_module (菜单/模块)
type PrivilegeModule struct {
    BaseModel
    Name     string  `gorm:"column:name;type:varchar(64)" json:"name"`
    Code     string  `gorm:"column:code;type:varchar(64)" json:"code"`
    PID      *string `gorm:"column:pid;type:varchar(64)" json:"pid"`
    URL      string  `gorm:"column:url;type:varchar(256)" json:"url"`
    Icon     string  `gorm:"column:icon;type:varchar(128)" json:"icon"`
    Sort     int     `gorm:"column:sort;default:0" json:"sort"`
    SystemID string  `gorm:"column:system_id;type:varchar(64)" json:"systemId"`
    Type     int     `gorm:"column:type;default:0" json:"type"` // 0: menu, 1: button
}

func (PrivilegeModule) TableName() string { return "tbl_privilege_module" }

// PrivilegeAcl — tbl_privilege_acl (访问控制列表)
type PrivilegeAcl struct {
    BaseModel
    RoleID      string `gorm:"column:role_id;type:varchar(64);index" json:"roleId"`
    ModuleID    string `gorm:"column:module_id;type:varchar(64);index" json:"moduleId"`
    Permission  string `gorm:"column:permission;type:varchar(32)" json:"permission"` // r/w/rw
    ReleaseID   string `gorm:"column:release_id;type:varchar(64)" json:"releaseId"`
    CheckStatus int    `gorm:"column:check_status;default:0" json:"checkStatus"`
}

func (PrivilegeAcl) TableName() string { return "tbl_privilege_acl" }

// PrivilegeDepartment — tbl_privilege_department
type PrivilegeDepartment struct {
    BaseModel
    Name      string  `gorm:"column:name;type:varchar(64)" json:"name"`
    Code      string  `gorm:"column:code;type:varchar(64)" json:"code"`
    PID       *string `gorm:"column:pid;type:varchar(64)" json:"pid"`
    CompanyID string  `gorm:"column:company_id;type:varchar(64)" json:"companyId"`
    Sort      int     `gorm:"column:sort;default:0" json:"sort"`
}

func (PrivilegeDepartment) TableName() string { return "tbl_privilege_department" }

// PrivilegeCompany — tbl_privilege_company
type PrivilegeCompany struct {
    BaseModel
    CName   string `gorm:"column:cname;type:varchar(128)" json:"cname"`
    EName   string `gorm:"column:ename;type:varchar(128)" json:"ename"`
    Code    string `gorm:"column:code;type:varchar(64)" json:"code"`
    SN      string `gorm:"column:sn;type:varchar(64)" json:"sn"`
    Manager string `gorm:"column:manager;type:varchar(64)" json:"manager"`
    Note    string `gorm:"column:note;type:varchar(256)" json:"note"`
}

func (PrivilegeCompany) TableName() string { return "tbl_privilege_company" }

// PrivilegeEmployee — tbl_privilege_employee
type PrivilegeEmployee struct {
    BaseModel
    UserCode string `gorm:"column:user_code;type:varchar(64)" json:"userCode"`
    EmpCode  string `gorm:"column:emp_code;type:varchar(64)" json:"empCode"`
    DeptID   string `gorm:"column:dept_id;type:varchar(64)" json:"deptId"`
}

func (PrivilegeEmployee) TableName() string { return "tbl_privilege_employee" }

// PrivilegeDictionary — tbl_privilege_dictionary
type PrivilegeDictionary struct {
    BaseModel
    Code       string  `gorm:"column:code;type:varchar(64)" json:"code"`
    Name       string  `gorm:"column:name;type:varchar(64)" json:"name"`
    PCode      *string `gorm:"column:pcode;type:varchar(64)" json:"pcode"`
    SystemSN   string  `gorm:"column:system_sn;type:varchar(64)" json:"systemSn"`
    Sort       int     `gorm:"column:sort;default:0" json:"sort"`
    SystemID   string  `gorm:"column:system_id;type:varchar(64)" json:"systemId"`
    Sn         string  `gorm:"column:sn;type:varchar(64)" json:"sn"`
    HasChild   int     `gorm:"column:has_child;default:0" json:"hasChild"`
}

func (PrivilegeDictionary) TableName() string { return "tbl_privilege_dictionary" }

// PrivilegePvalue — tbl_privilege_pvalue (权限值)
type PrivilegePvalue struct {
    BaseModel
    Code     string `gorm:"column:code;type:varchar(64)" json:"code"`
    Name     string `gorm:"column:name;type:varchar(64)" json:"name"`
    SystemID string `gorm:"column:system_id;type:varchar(64)" json:"systemId"`
}

func (PrivilegePvalue) TableName() string { return "tbl_privilege_pvalue" }

// PrivilegeLoginLog — tbl_privilege_login_log
type PrivilegeLoginLog struct {
    BaseModel
    UserID    string `gorm:"column:user_id;type:varchar(64)" json:"userId"`
    Username  string `gorm:"column:username;type:varchar(64)" json:"username"`
    LoginIP   string `gorm:"column:login_ip;type:varchar(64)" json:"loginIp"`
    LoginTime time.Time `gorm:"column:login_time" json:"loginTime"`
}

func (PrivilegeLoginLog) TableName() string { return "tbl_privilege_login_log" }
```

- [ ] **Step 3: `go build ./internal/model/`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add BaseModel and 12 privilege GORM entities`

---

### Task 2: DTOs + VOs

**Files:** `internal/model/privilege_dto.go`, `internal/model/privilege_vo.go`

- [ ] **Step 1: DTOs (request bodies + page queries)**

```go
// internal/model/privilege_dto.go
package model

// LoginInfoDTO — POST /api/privilege/auth/login
type LoginInfoDTO struct {
    Type       string `json:"type" binding:"required"`
    Username   string `json:"username" binding:"required"`
    Password   string `json:"password" binding:"required"`
    CaptchaKey string `json:"captchaKey"`
    CaptchaCode string `json:"captchaCode"`
}

type PasswordUpdateDTO struct {
    UserID      string `json:"userId" binding:"required"`
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword" binding:"required"`
}

type PrivilegeUserDTO struct {
    EmployeeID string `json:"employeeId"`
    Code       string `json:"code"`
    RealName   string `json:"realName"`
    Username   string `json:"username" binding:"required"`
    Mobile     string `json:"mobile"`
    Email      string `json:"email"`
    CompanyID  string `json:"companyId"`
    DeptID     string `json:"deptId"`
    Status     int    `json:"status"`
}

type PrivilegeUserPageQuery struct {
    Page     int    `form:"page"`
    Size     int    `form:"size"`
    Username string `form:"username"`
    RealName string `form:"realName"`
    Status   *int   `form:"status"`
}

type PrivilegeRoleDTO struct {
    Name       string `json:"name" binding:"required"`
    SN         string `json:"sn"`
    RoleLevel  string `json:"roleLevel"`
    Note       string `json:"note"`
    ValidState int    `json:"validState"`
    CompanyID  int64  `json:"companyId"`
    SystemID   string `json:"systemId"`
}

type PrivilegeRoleQuery struct {
    Page      int    `json:"page"`
    Size      int    `json:"size"`
    Name      string `json:"name"`
    CompanyID int64  `json:"companyId"`
}

type PrivilegeUserRoleDTO struct {
    UserID string `json:"userId" binding:"required"`
    RoleID string `json:"roleId" binding:"required"`
}

type UserRoleBatchDTO struct {
    UserID  string   `json:"userId" binding:"required"`
    RoleIDs []string `json:"roleIds" binding:"required"`
}

type PrivilegeModuleDTO struct {
    Name     string `json:"name" binding:"required"`
    Code     string `json:"code"`
    PID      string `json:"pid"`
    URL      string `json:"url"`
    Icon     string `json:"icon"`
    Sort     int    `json:"sort"`
    SystemID string `json:"systemId"`
    Type     int    `json:"type"`
}

type PrivilegeAclDTO struct {
    RoleID     string `json:"roleId"`
    ModuleID   string `json:"moduleId"`
    Permission string `json:"permission"`
    ReleaseID  string `json:"releaseId"`
}

type PrivilegeDepartmentDTO struct {
    Name      string `json:"name" binding:"required"`
    Code      string `json:"code"`
    PID       string `json:"pid"`
    CompanyID string `json:"companyId"`
    Sort      int    `json:"sort"`
}

type PrivilegeCompanyDTO struct {
    CName   string `json:"cname" binding:"required"`
    EName   string `json:"ename"`
    Code    string `json:"code"`
    SN      string `json:"sn"`
    Manager string `json:"manager"`
    Note    string `json:"note"`
}

type PrivilegeCompanyQuery struct {
    Page int    `json:"page"`
    Size int    `json:"size"`
    Name string `json:"name"`
}

type PrivilegeEmployeeDTO struct {
    UserCode string `json:"userCode" binding:"required"`
    EmpCode  string `json:"empCode" binding:"required"`
    DeptID   string `json:"deptId"`
}

type PrivilegeDictionaryDTO struct {
    Code     string `json:"code" binding:"required"`
    Name     string `json:"name" binding:"required"`
    PCode    string `json:"pcode"`
    SystemSN string `json:"systemSn"`
    SystemID string `json:"systemId"`
    Sort     int    `json:"sort"`
}

type PrivilegePvalueDTO struct {
    Code     string `json:"code" binding:"required"`
    Name     string `json:"name" binding:"required"`
    SystemID string `json:"systemId"`
}

type PrivilegePvalueQuery struct {
    Page     int    `json:"page"`
    Size     int    `json:"size"`
    SystemID string `json:"systemId"`
}

type PrivilegeLoginLogDTO struct {
    UserID   string `json:"userId"`
    Username string `json:"username"`
    LoginIP  string `json:"loginIp"`
}
```

- [ ] **Step 2: VOs (response objects)**

```go
// internal/model/privilege_vo.go
package model

import "time"

type CaptchaVO struct {
    CaptchaKey string `json:"captchaKey"`
    Image      string `json:"image"` // base64
}

type LoginUserInfoVO struct {
    UserID   string       `json:"userId"`
    Username string       `json:"username"`
    RealName string       `json:"realName"`
    Token    string       `json:"token"`
    Roles    []string     `json:"roles"`
    Menus    []UserMenuVO `json:"menus"`
}

type UserMenuVO struct {
    ID       string        `json:"id"`
    Name     string        `json:"name"`
    Code     string        `json:"code"`
    URL      string        `json:"url"`
    Icon     string        `json:"icon"`
    Sort     int           `json:"sort"`
    Type     int           `json:"type"`
    Children []UserMenuVO  `json:"children,omitempty"`
}

type PrivilegeUserVO struct {
    ID          string  `json:"id"`
    EmployeeID  string  `json:"employeeId"`
    Code        string  `json:"code"`
    RealName    string  `json:"realName"`
    Username    string  `json:"username"`
    Tel         string  `json:"tel"`
    Mobile      string  `json:"mobile"`
    Email       string  `json:"email"`
    CompanyID   string  `json:"companyId"`
    CompanyName string  `json:"companyName"`
    DeptID      string  `json:"deptId"`
    DeptName    string  `json:"deptName"`
    Status      int     `json:"status"`
    UserType    int     `json:"userType"`
    CreateTime  time.Time `json:"createTime"`
}

type PrivilegeRoleVO struct {
    BaseModel
    Name        string `json:"name"`
    SN          string `json:"sn"`
    RoleLevel   string `json:"roleLevel"`
    Note        string `json:"note"`
    ValidState  int    `json:"validState"`
    CompanyID   int64  `json:"companyId"`
    CompanyName string `json:"companyName"`
    SystemID    string `json:"systemId"`
}

type RoleAclVO struct {
    ModuleID   string `json:"moduleId"`
    ModuleName string `json:"moduleName"`
    Permission string `json:"permission"`
    Checked    bool   `json:"checked"`
}

type ModuleTreeVO struct {
    PrivilegeModule
    Children []ModuleTreeVO `json:"children,omitempty"`
}

type OrganizationTreeVO struct {
    ID       string               `json:"id"`
    Name     string               `json:"name"`
    Code     string               `json:"code"`
    Type     string               `json:"type"` // "company" | "department"
    Children []OrganizationTreeVO `json:"children,omitempty"`
}

type PrivilegeUserRoleVO struct {
    BaseModel
    UserID   string `json:"userId"`
    RoleID   string `json:"roleId"`
    RoleName string `json:"roleName"`
    Username string `json:"username"`
}

type PrivilegeCompanyVO struct {
    PrivilegeCompany
}

type PrivilegeDepartmentVO struct {
    PrivilegeDepartment
    CompanyName string `json:"companyName"`
}

type PrivilegeEmployeeVO struct {
    PrivilegeEmployee
    Username string `json:"username"`
    RealName string `json:"realName"`
    DeptName string `json:"deptName"`
}

type PrivilegeDictionaryVO struct {
    PrivilegeDictionary
}

type PrivilegePvalueVO struct {
    PrivilegePvalue
}

type PrivilegeLoginLogVO struct {
    PrivilegeLoginLog
}

type PrivilegeModuleVO struct {
    PrivilegeModule
}

type PrivilegeAclVO struct {
    PrivilegeAcl
    ModuleName string `json:"moduleName"`
}
```

- [ ] **Step 3: `go build ./internal/model/`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add privilege DTOs and VOs`

---

### Task 3: Repository Interfaces + GORM Implementations

**Files:** `internal/repository/privilege_repo.go`, `internal/dao/db/privilege_repo.go`

- [ ] **Step 1: Repository interfaces**

```go
// internal/repository/privilege_repo.go
package repository

import (
    "context"
    "github.com/phoenix-agent-go/internal/model"
)

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeUser, error)
    FindByUsername(ctx context.Context, username string) (*model.PrivilegeUser, error)
    FindByCode(ctx context.Context, code string) (*model.PrivilegeUser, error)
    FindByMobile(ctx context.Context, mobile string) (*model.PrivilegeUser, error)
    Page(ctx context.Context, query model.PrivilegeUserPageQuery) ([]*model.PrivilegeUser, int64, error)
    Create(ctx context.Context, user *model.PrivilegeUser) error
    Update(ctx context.Context, user *model.PrivilegeUser) error
    Delete(ctx context.Context, id string) error
}

type RoleRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeRole, error)
    FindByCompanyID(ctx context.Context, companyID int64) ([]*model.PrivilegeRole, error)
    Page(ctx context.Context, query model.PrivilegeRoleQuery) ([]*model.PrivilegeRole, int64, error)
    Create(ctx context.Context, role *model.PrivilegeRole) error
    Update(ctx context.Context, role *model.PrivilegeRole) error
    Delete(ctx context.Context, id string) error
}

type UserRoleRepository interface {
    FindByUserID(ctx context.Context, userID string) ([]*model.PrivilegeUserRole, error)
    FindByRoleID(ctx context.Context, roleID string) ([]*model.PrivilegeUserRole, error)
    SaveBatch(ctx context.Context, userID string, roleIDs []string) error
    DeleteBatch(ctx context.Context, userID string, roleIDs []string) error
}

type ModuleRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeModule, error)
    FindByPID(ctx context.Context, pid string) ([]*model.PrivilegeModule, error)
    FindAll(ctx context.Context) ([]*model.PrivilegeModule, error)
    Tree(ctx context.Context) ([]*model.ModuleTreeVO, error)
    Create(ctx context.Context, module *model.PrivilegeModule) error
    Update(ctx context.Context, module *model.PrivilegeModule) error
}

type ACLRepository interface {
    FindByRoleID(ctx context.Context, roleID string) ([]*model.PrivilegeAcl, error)
    SaveAll(ctx context.Context, acls []*model.PrivilegeAcl) error
    SaveModule(ctx context.Context, acl *model.PrivilegeAcl) error
    FindByReleaseID(ctx context.Context, releaseID string) ([]*model.PrivilegeAcl, error)
}

type DepartmentRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeDepartment, error)
    FindByPID(ctx context.Context, pid string) ([]*model.PrivilegeDepartment, error)
    FindByCompanyID(ctx context.Context, companyID string) ([]*model.PrivilegeDepartment, error)
    OrgTree(ctx context.Context) ([]*model.OrganizationTreeVO, error)
    Page(ctx context.Context, page, size int, dept *model.PrivilegeDepartment) ([]*model.PrivilegeDepartment, int64, error)
    Create(ctx context.Context, dept *model.PrivilegeDepartment) error
    Update(ctx context.Context, dept *model.PrivilegeDepartment) error
    Delete(ctx context.Context, id string) error
}

type CompanyRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeCompany, error)
    FindByCode(ctx context.Context, code string) (*model.PrivilegeCompany, error)
    Page(ctx context.Context, query model.PrivilegeCompanyQuery) ([]*model.PrivilegeCompany, int64, error)
    Create(ctx context.Context, company *model.PrivilegeCompany) error
    Update(ctx context.Context, company *model.PrivilegeCompany) error
    Delete(ctx context.Context, id string) error
}

type EmployeeRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeEmployee, error)
    FindByEmpCode(ctx context.Context, empCode string) (*model.PrivilegeEmployee, error)
    Page(ctx context.Context, page, size int) ([]*model.PrivilegeEmployee, int64, error)
    Create(ctx context.Context, emp *model.PrivilegeEmployee) error
    Update(ctx context.Context, emp *model.PrivilegeEmployee) error
    Delete(ctx context.Context, id string) error
}

type DictionaryRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegeDictionary, error)
    FindBySystemSN(ctx context.Context, systemSN string) ([]*model.PrivilegeDictionary, error)
    FindByPCode(ctx context.Context, pcode string) ([]*model.PrivilegeDictionary, error)
    Page(ctx context.Context, page, size int) ([]*model.PrivilegeDictionary, int64, error)
    Create(ctx context.Context, dict *model.PrivilegeDictionary) error
    Update(ctx context.Context, dict *model.PrivilegeDictionary) error
    Delete(ctx context.Context, id string) error
}

type PvalueRepository interface {
    FindByID(ctx context.Context, id string) (*model.PrivilegePvalue, error)
    Page(ctx context.Context, query model.PrivilegePvalueQuery) ([]*model.PrivilegePvalue, int64, error)
    Create(ctx context.Context, pv *model.PrivilegePvalue) error
    Update(ctx context.Context, pv *model.PrivilegePvalue) error
    Delete(ctx context.Context, id string) error
}

type LoginLogRepository interface {
    Create(ctx context.Context, log *model.PrivilegeLoginLog) error
    Page(ctx context.Context, page, size int) ([]*model.PrivilegeLoginLog, int64, error)
}
```

- [ ] **Step 2: GORM implementations** — Create `internal/dao/db/privilege_repo.go` with GORM implementations of ALL interfaces above. Each implementation wraps `gorm.DB` with `WithContext(ctx)`. Key methods:

```go
// internal/dao/db/privilege_repo.go
package db

import (
    "context"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/phoenix-agent-go/internal/repository"
    "gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) repository.UserRepository { return &userRepo{db} }

func (r *userRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeUser, error) {
    var user model.PrivilegeUser
    err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&user).Error
    if err != nil { return nil, err }
    return &user, nil
}
// ... ALL other methods

type roleRepo struct{ db *gorm.DB }
func NewRoleRepository(db *gorm.DB) repository.RoleRepository { return &roleRepo{db} }
// ... ALL other methods

type userRoleRepo struct{ db *gorm.DB }
func NewUserRoleRepository(db *gorm.DB) repository.UserRoleRepository { return &userRoleRepo{db} }
// ... ALL other methods

// ... ALL 10 repo structs with full implementations
```

**IMPORTANT:** Implement EVERY method for EVERY repository. Do not skip any. Each method calls `r.db.WithContext(ctx)` and the appropriate GORM operation. For DELETE, use soft-delete: `Update("del_flag", 1)`. For Page, use `Offset((page-1)*size).Limit(size).Find(&list)` and `Count(&total)`.

- [ ] **Step 3: `go build ./internal/...`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add privilege repository interfaces and GORM implementations`

---

### Task 4: Domain Logic + Cache Layer

**Files:** `internal/domain/privilege/domain.go`, `internal/dao/cache/privilege_cache.go`

- [ ] **Step 1: Domain logic** — `internal/domain/privilege/domain.go`

```go
package privilege

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/phoenix-agent-go/internal/model"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(hashed, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func GenerateRandomPassword() string {
    // 8 chars: letters + digits
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, 8)
    for i := range b {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        b[i] = charset[n.Int64()]
    }
    return string(b)
}

func IsValidUser(user *model.PrivilegeUser) bool {
    return user != nil && user.Status == 1
}
```
Add imports for `crypto/rand`, `math/big`, `golang.org/x/crypto/bcrypt`.

- [ ] **Step 2: Cache layer** — `internal/dao/cache/privilege_cache.go`

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/allegro/bigcache/v3"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/redis/go-redis/v9"
)

type PrivilegeCache struct {
    redis *redis.Client
    local *bigcache.BigCache
    ttl   time.Duration
}

func NewPrivilegeCache(redis *redis.Client, local *bigcache.BigCache) *PrivilegeCache {
    return &PrivilegeCache{redis: redis, local: local, ttl: 10 * time.Minute}
}

func (c *PrivilegeCache) key(prefix, id string) string { return "privilege:" + prefix + ":" + id }

func (c *PrivilegeCache) GetUser(ctx context.Context, id string) (*model.PrivilegeUser, error) {
    // L1: BigCache
    if data, err := c.local.Get(c.key("user", id)); err == nil {
        var user model.PrivilegeUser
        json.Unmarshal(data, &user)
        return &user, nil
    }
    // L2: Redis
    data, err := c.redis.Get(ctx, c.key("user", id)).Bytes()
    if err != nil { return nil, err }
    var user model.PrivilegeUser
    json.Unmarshal(data, &user)
    c.local.Set(c.key("user", id), data) // backfill L1
    return &user, nil
}

func (c *PrivilegeCache) SetUser(ctx context.Context, user *model.PrivilegeUser) error {
    data, _ := json.Marshal(user)
    c.local.Set(c.key("user", user.ID), data)
    return c.redis.Set(ctx, c.key("user", user.ID), data, c.ttl).Err()
}

func (c *PrivilegeCache) InvalidateUser(ctx context.Context, id string) error {
    c.local.Delete(c.key("user", id))
    return c.redis.Del(ctx, c.key("user", id)).Err()
}

// Similarly: SetRole, GetRole, InvalidateRole, SetUserRoles, GetUserRoles
```

- [ ] **Step 3: `go build ./...`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add privilege domain logic and cache layer`

---

### Task 5: Usecase + Service Layer

**Files:** `internal/usecase/privilege_usecase.go`, `internal/service/privilege_service.go`

- [ ] **Step 1: Usecase** — orchestrates repository + domain + cache

```go
// internal/usecase/privilege_usecase.go
package usecase

import (
    "context"
    "github.com/phoenix-agent-go/internal/dao/cache"
    "github.com/phoenix-agent-go/internal/domain/privilege"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/phoenix-agent-go/internal/repository"
    "github.com/phoenix-agent-go/infra/id"
)

type PrivilegeUsecase struct {
    userRepo       repository.UserRepository
    roleRepo       repository.RoleRepository
    userRoleRepo   repository.UserRoleRepository
    moduleRepo     repository.ModuleRepository
    aclRepo        repository.ACLRepository
    deptRepo       repository.DepartmentRepository
    companyRepo    repository.CompanyRepository
    employeeRepo   repository.EmployeeRepository
    dictRepo       repository.DictionaryRepository
    pvalueRepo     repository.PvalueRepository
    loginLogRepo   repository.LoginLogRepository
    cache          *cache.PrivilegeCache
}

func NewPrivilegeUsecase(
    userRepo repository.UserRepository,
    roleRepo repository.RoleRepository,
    userRoleRepo repository.UserRoleRepository,
    moduleRepo repository.ModuleRepository,
    aclRepo repository.ACLRepository,
    deptRepo repository.DepartmentRepository,
    companyRepo repository.CompanyRepository,
    employeeRepo repository.EmployeeRepository,
    dictRepo repository.DictionaryRepository,
    pvalueRepo repository.PvalueRepository,
    loginLogRepo repository.LoginLogRepository,
    cache *cache.PrivilegeCache,
) *PrivilegeUsecase {
    return &PrivilegeUsecase{
        userRepo, roleRepo, userRoleRepo, moduleRepo, aclRepo,
        deptRepo, companyRepo, employeeRepo, dictRepo, pvalueRepo, loginLogRepo, cache,
    }
}

// --- User ---
func (u *PrivilegeUsecase) Login(ctx context.Context, dto model.LoginInfoDTO, ip string) (*model.LoginUserInfoVO, error) {
    user, err := u.userRepo.FindByUsername(ctx, dto.Username)
    if err != nil { return nil, ErrInvalidCredentials }
    if !privilege.IsValidUser(user) { return nil, ErrUserDisabled }
    if !privilege.CheckPassword(user.Password, dto.Password) { return nil, ErrInvalidCredentials }
    
    // Get roles
    userRoles, _ := u.userRoleRepo.FindByUserID(ctx, user.ID)
    roleNames := []string{}
    for _, ur := range userRoles {
        role, _ := u.roleRepo.FindByID(ctx, ur.RoleID)
        if role != nil { roleNames = append(roleNames, role.SN) }
    }
    
    // Get menus
    menus, _ := u.moduleRepo.Tree(ctx)
    
    // Log login
    u.loginLogRepo.Create(ctx, &model.PrivilegeLoginLog{
        UserID: user.ID, Username: user.Username, LoginIP: ip,
    })
    
    return &model.LoginUserInfoVO{
        UserID: user.ID, Username: user.Username, RealName: user.RealName,
        Roles: roleNames, Menus: menus,
    }, nil
}

// Add ALL other usecase methods: CreateUser, UpdateUser, DeleteUser, ResetPassword,
// UpdatePassword, PageUsers, CreateRole, PageRoles, GetRoleAcls,
// SaveUserRoles, DeleteUserRoles, GetModuleTree, SaveACLs,
// CreateDept, UpdateDept, GetOrgTree, CreateCompany, PageCompanies,
// CreateEmployee, SyncEmployees, PageDictionaries,
// CreatePvalue, PagePvalues, PageLoginLogs, etc.

var (
    ErrInvalidCredentials = &AppError{Code: 401001, Msg: "用户名或密码错误"}
    ErrUserDisabled       = &AppError{Code: 401002, Msg: "用户已被禁用"}
    ErrUsernameExists     = &AppError{Code: 401003, Msg: "用户名已存在"}
    ErrMobileExists       = &AppError{Code: 401004, Msg: "手机号已存在"}
    ErrOldPasswordWrong   = &AppError{Code: 401005, Msg: "原密码错误"}
    ErrRoleNotFound       = &AppError{Code: 402001, Msg: "角色不存在"}
)

type AppError struct {
    Code int
    Msg  string
}
func (e *AppError) Error() string { return e.Msg }
```

- [ ] **Step 2: Service** — thin wrapper for format conversion, `internal/service/privilege_service.go`

```go
// internal/service/privilege_service.go
package service

import (
    "context"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/phoenix-agent-go/internal/usecase"
)

type PrivilegeService struct {
    uc *usecase.PrivilegeUsecase
}

func NewPrivilegeService(uc *usecase.PrivilegeUsecase) *PrivilegeService {
    return &PrivilegeService{uc}
}

// Delegates to usecase, converts AppError to errcode.ErrCode in handler layer
func (s *PrivilegeService) Login(ctx context.Context, dto model.LoginInfoDTO, ip string) (*model.LoginUserInfoVO, error) {
    return s.uc.Login(ctx, dto, ip)
}

// ... pass-through methods for all usecase operations
```

- [ ] **Step 3: `go build ./...`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add privilege usecase and service layer`

---

### Task 6: JWT + Casbin Middleware (replace Phase 1 stubs)

**Files:**
- Modify: `api/middleware/auth.go`
- Modify: `api/middleware/rbac.go`
- Create: `internal/config/casbin_model.conf`

- [ ] **Step 1: JWT Auth middleware** — replace stub

```go
// api/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
    "github.com/phoenix-agent-go/infra/jwt"
    "github.com/phoenix-agent-go/infra/response"
)

func Auth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Error(c, errcode.Unauthorized)
            c.Abort()
            return
        }
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.Error(c, errcode.Unauthorized)
            c.Abort()
            return
        }
        claims, err := jwtManager.ParseToken(parts[1])
        if err != nil {
            response.ErrorWithStatus(c, http.StatusUnauthorized, errcode.Unauthorized)
            c.Abort()
            return
        }
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

Note: Phase 1 `Auth()` becomes `Auth(jwtManager)`. Update `api/router.go` to pass jwtManager.

- [ ] **Step 2: Casbin RBAC** — replace stub

```go
// api/middleware/rbac.go
package middleware

import (
    "fmt"
    "github.com/casbin/casbin/v2"
    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
    "github.com/phoenix-agent-go/infra/response"
)

func RBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")
        obj := c.Request.URL.Path
        act := c.Request.Method
        ok, err := enforcer.Enforce(fmt.Sprint(userID), obj, act)
        if err != nil || !ok {
            response.Error(c, errcode.Forbidden)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

- [ ] **Step 3: Casbin model config** — `internal/config/casbin_model.conf`

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
```

- [ ] **Step 4: `go build ./...`** — PASS
- [ ] **Step 5: Commit** `feat(phase2): implement JWT auth and Casbin RBAC middleware`

---

### Task 7: Login Handler

**Files:** `api/handler/privilege/auth.go`

- [ ] **Step 1: Captcha + Login + Logout + Menus + UserInfo**

```go
// api/handler/privilege/auth.go
package privilege

import (
    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
    "github.com/phoenix-agent-go/infra/jwt"
    "github.com/phoenix-agent-go/infra/response"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/phoenix-agent-go/internal/service"
)

type AuthHandler struct {
    svc        *service.PrivilegeService
    jwtManager *jwt.JWTManager
}

func NewAuthHandler(svc *service.PrivilegeService, jwtManager *jwt.JWTManager) *AuthHandler {
    return &AuthHandler{svc, jwtManager}
}

// GET /api/privilege/auth/captcha
func (h *AuthHandler) Captcha(c *gin.Context) {
    // Generate simple math captcha
    captchaKey := "captcha_" + id.MustGenerateIDStr()
    c.SetCookie("captcha_key", captchaKey, 300, "/", "", false, true)
    response.Success(c, gin.H{"captchaKey": captchaKey, "image": "placeholder"})
}

// POST /api/privilege/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
    var dto model.LoginInfoDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        response.Error(c, errcode.InvalidParams); return
    }
    ip := c.ClientIP()
    userInfo, err := h.svc.Login(c.Request.Context(), dto, ip)
    if err != nil {
        if appErr, ok := err.(*usecase.AppError); ok {
            response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
            return
        }
        response.Error(c, errcode.InternalError); return
    }
    token, _ := h.jwtManager.GenerateToken(userInfo.UserID, userInfo.Username)
    userInfo.Token = token
    response.Success(c, userInfo)
}

// POST /api/privilege/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
    response.Success(c, "退出成功")
}

// GET /api/privilege/auth/menus
func (h *AuthHandler) Menus(c *gin.Context) {
    userID, _ := c.Get("user_id")
    menus, err := h.svc.GetUserMenus(c.Request.Context(), userID.(string))
    if err != nil {
        response.Error(c, errcode.InternalError); return
    }
    response.Success(c, menus)
}

// GET /api/privilege/auth/getLoginUserInfo
func (h *AuthHandler) GetLoginUserInfo(c *gin.Context) {
    userID, _ := c.Get("user_id")
    user, err := h.svc.GetUserByID(c.Request.Context(), userID.(string))
    if err != nil {
        response.Error(c, errcode.Unauthorized); return
    }
    response.Success(c, user)
}
```

- [ ] **Step 2: Register auth routes in router**

In `api/router.go`, update the auth group:
```go
authHandler := privilege.NewAuthHandler(privilegeSvc, jwtManager)
auth := r.Group("/api/privilege/auth")
auth.Use(middleware.RateLimit())
auth.GET("/captcha", authHandler.Captcha)
auth.POST("/login", authHandler.Login)
auth.POST("/logout", authHandler.Logout)
auth.GET("/menus", middleware.Auth(jwtManager), authHandler.Menus)
auth.GET("/getLoginUserInfo", middleware.Auth(jwtManager), authHandler.GetLoginUserInfo)
```

- [ ] **Step 3: `go build ./...`** — PASS
- [ ] **Step 4: Commit** `feat(phase2): add login handler with captcha, login, logout, menus`

---

### Task 8: User Handler

**Files:** `api/handler/privilege/user.go`

- [ ] **Step 1: Full CRUD + Password + Reset Password**

```go
// api/handler/privilege/user.go
package privilege

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
    "github.com/phoenix-agent-go/infra/response"
    "github.com/phoenix-agent-go/internal/model"
    "github.com/phoenix-agent-go/internal/service"
)

type UserHandler struct { svc *service.PrivilegeService }

func NewUserHandler(svc *service.PrivilegeService) *UserHandler { return &UserHandler{svc} }

// GET /api/privilege/user/page
func (h *UserHandler) Page(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    var query model.PrivilegeUserPageQuery
    c.ShouldBindQuery(&query)
    query.Page = page; query.Size = size
    list, total, err := h.svc.PageUsers(c.Request.Context(), query)
    if err != nil { response.Error(c, errcode.InternalError); return }
    response.SuccessPage(c, list, total, page, size)
}

// GET /api/privilege/user/:id
func (h *UserHandler) GetByID(c *gin.Context) {
    user, err := h.svc.GetUserByID(c.Request.Context(), c.Param("id"))
    if err != nil { response.Error(c, errcode.NotFound); return }
    response.Success(c, user)
}

// GET /api/privilege/user/code/:code
func (h *UserHandler) GetByCode(c *gin.Context) {
    user, err := h.svc.GetUserByCode(c.Request.Context(), c.Param("code"))
    if err != nil { response.Error(c, errcode.NotFound); return }
    response.Success(c, user)
}

// POST /api/privilege/user
func (h *UserHandler) Create(c *gin.Context) {
    var dto model.PrivilegeUserDTO
    if err := c.ShouldBindJSON(&dto); err != nil { response.Error(c, errcode.InvalidParams); return }
    if err := h.svc.CreateUser(c.Request.Context(), dto); err != nil {
        if appErr, ok := err.(*usecase.AppError); ok {
            response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
            return
        }
        response.Error(c, errcode.InternalError); return
    }
    response.Success(c, true)
}

// PUT /api/privilege/user
func (h *UserHandler) Update(c *gin.Context) {
    var dto model.PrivilegeUserDTO
    if err := c.ShouldBindJSON(&dto); err != nil { response.Error(c, errcode.InvalidParams); return }
    if err := h.svc.UpdateUser(c.Request.Context(), dto); err != nil {
        response.Error(c, errcode.InternalError); return
    }
    response.Success(c, true)
}

// DELETE /api/privilege/user/:id
func (h *UserHandler) Delete(c *gin.Context) {
    if err := h.svc.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
        response.Error(c, errcode.InternalError); return
    }
    response.Success(c, true)
}

// PUT /api/privilege/user/password
func (h *UserHandler) UpdatePassword(c *gin.Context) {
    var dto model.PasswordUpdateDTO
    if err := c.ShouldBindJSON(&dto); err != nil { response.Error(c, errcode.InvalidParams); return }
    if err := h.svc.UpdatePassword(c.Request.Context(), dto); err != nil {
        if appErr, ok := err.(*usecase.AppError); ok {
            response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
            return
        }
        response.Error(c, errcode.InternalError); return
    }
    response.Success(c, "密码修改成功")
}

// PUT /api/privilege/user/reset-password/:id
func (h *UserHandler) ResetPassword(c *gin.Context) {
    newPwd, err := h.svc.ResetPassword(c.Request.Context(), c.Param("id"))
    if err != nil { response.ErrorWithMsg(c, errcode.NotFound, "用户不存在"); return }
    response.Success(c, newPwd)
}
```

- [ ] **Step 2: Register user routes**

```go
userHandler := privilege.NewUserHandler(privilegeSvc)
userGroup := r.Group("/api/privilege/user")
userGroup.Use(middleware.Auth(jwtManager))
{
    userGroup.GET("/page", userHandler.Page)
    userGroup.GET("/:id", userHandler.GetByID)
    userGroup.GET("/code/:code", userHandler.GetByCode)
    userGroup.POST("", userHandler.Create)
    userGroup.PUT("", userHandler.Update)
    userGroup.DELETE("/:id", userHandler.Delete)
    userGroup.PUT("/password", userHandler.UpdatePassword)
    userGroup.PUT("/reset-password/:id", userHandler.ResetPassword)
}
```

- [ ] **Step 3: Commit** `feat(phase2): add user handler with CRUD + password management`

---

### Task 9: Role + UserRole Handlers

**Files:** `api/handler/privilege/role.go`, `api/handler/privilege/user_role.go`

- [ ] **Step 1: Role Handler** — endpoints: `POST /page`, `GET /:id`, `GET /company/:companyId`, `POST`, `PUT`, `DELETE /:id`, `GET /:roleId/acls`

Full handler with all 7 endpoints, similar CRUD pattern as User handler. Map responses to `model.PrivilegeRoleVO`.

- [ ] **Step 2: UserRole Handler** — endpoints: `GET /page`, `GET /user/:userId`, `GET /role/:roleId`, `POST`, `DELETE /:id`, `POST /batch-save`, `DELETE /batch-remove`

- [ ] **Step 3: Register routes** — same pattern as user routes
- [ ] **Step 4: Commit** `feat(phase2): add role and user-role handlers`

---

### Task 10: Module + ACL + Department + Company Handlers

**Files:** `api/handler/privilege/module.go`, `api/handler/privilege/acl.go`, `api/handler/privilege/department.go`, `api/handler/privilege/company.go`

- [ ] **Step 1: Module** — `GET /page`, `GET /tree`, `GET /tree/acl`, `GET /:id`, `GET /pid/:pid`, `POST`, `PUT`
- [ ] **Step 2: ACL** — `GET /release/:releaseId`, `POST /saveAll/:releaseId/:checkStatus`, `POST /saveModule`
- [ ] **Step 3: Department** — `GET /orgTree`, `GET /page`, `GET /tree`, `GET /:id`, `GET /pid/:pid`, `GET /company/:companyId`, `GET /code/:code`, `POST`, `PUT`, `DELETE /:id`, `POST /sync`, `POST /sync-children/:deptId`
- [ ] **Step 4: Company** — `POST /page`, `GET /:id`, `GET /code/:code`, `POST`, `PUT`, `DELETE /:id`
- [ ] **Step 5: Register all routes with Auth middleware**
- [ ] **Step 6: Commit** `feat(phase2): add module, ACL, department, and company handlers`

---

### Task 11: Employee + Dictionary + Pvalue + LoginLog Handlers

**Files:** `api/handler/privilege/employee.go`, `api/handler/privilege/dictionary.go`, `api/handler/privilege/pvalue.go`, `api/handler/privilege/login_log.go`

- [ ] **Step 1: Employee** — `GET /page`, `GET /:id`, `GET /emp-code/:empCode`, `POST`, `PUT`, `DELETE /:id`, `POST /sync`, `POST /sync-by-dept/:deptId`
- [ ] **Step 2: Dictionary** — `GET /page`, `GET /:id`, `GET /system/:systemSn`, `GET /pcode/:pcode`, `POST`, `PUT`, `DELETE /:id`
- [ ] **Step 3: Pvalue** — `POST /page`, `GET /:id`, `GET /system`, `POST`, `PUT`, `DELETE /:id`
- [ ] **Step 4: LoginLog** — `GET /page`, `GET /:id`, `POST`, `DELETE /:id`
- [ ] **Step 5: Register all routes**
- [ ] **Step 6: Commit** `feat(phase2): add employee, dictionary, pvalue, and login-log handlers`

---

### Task 12: Wiring + Integration Test

**Files:**
- Modify: `cmd/api/main.go` (wire up DB, repos, usecase, service, handlers)
- Modify: `api/router.go` (register ALL privilege routes)

- [ ] **Step 1: Wire up main.go** — Add GORM DB init, create all repos, usecase, service, handlers, pass to router

- [ ] **Step 2: Complete router.go** — Register all 12 handler groups

- [ ] **Step 3: Integration test**

```bash
go build ./...
go test ./... -short
go run ./cmd/api &
sleep 2

# Test login
curl -X POST http://localhost:8066/api/privilege/auth/login \
  -H "Content-Type: application/json" \
  -d '{"type":"admin","username":"admin","password":"admin123"}'

# Test authenticated endpoint
curl http://localhost:8066/api/privilege/user/page \
  -H "Authorization: Bearer <token>"
```

- [ ] **Step 4: Commit** `feat(phase2): wire up privilege module and integration verification`

---

## 里程碑

| M | Task | 验收 |
|:---|:---|:---|
| M1 | 1-2 | 12 Entity + DTO/VO 编译通过 |
| M2 | 3-4 | Repository + Domain + Cache 编译通过 |
| M3 | 5 | Usecase + Service 编译通过 |
| M4 | 6-7 | JWT/Casbin 中间件就绪，Login API 可调用 |
| M5 | 8-11 | 全部 12 个 REST 资源 CRUD 就绪 |
| M6 | 12 | 全量集成 — 登录→获取 token→调用受保护 API→返回正确数据 |

## 回退

每个 Task 独立 commit，`git revert <commit>` 撤销。Phase 2 不修改数据库 Schema。
