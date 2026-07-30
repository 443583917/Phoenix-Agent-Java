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
