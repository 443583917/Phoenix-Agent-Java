package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/usecase"
)

// PrivilegeService is a thin pass-through wrapper around PrivilegeUsecase.
// Error-code conversion (AppError to HTTP errcode) is handled by the
// handler layer.
type PrivilegeService struct {
	uc *usecase.PrivilegeUsecase
}

// NewPrivilegeService creates a new PrivilegeService with the given usecase.
func NewPrivilegeService(uc *usecase.PrivilegeUsecase) *PrivilegeService {
	return &PrivilegeService{uc: uc}
}

// ──────────────────────────── User ────────────────────────────

// Login authenticates a user and returns login info including menus and roles.
func (s *PrivilegeService) Login(ctx context.Context, dto model.LoginInfoDTO, ip string) (*model.LoginUserInfoVO, error) {
	return s.uc.Login(ctx, dto, ip)
}

// Logout is a no-op placeholder (stateless tokens).
func (s *PrivilegeService) Logout(ctx context.Context, userID string) error {
	return s.uc.Logout(ctx, userID)
}

// GetUserMenus returns the module tree for the given user.
func (s *PrivilegeService) GetUserMenus(ctx context.Context, userID string) ([]model.UserMenuVO, error) {
	return s.uc.GetUserMenus(ctx, userID)
}

// CreateUser creates a new privilege user.
func (s *PrivilegeService) CreateUser(ctx context.Context, dto model.PrivilegeUserDTO) (*model.PrivilegeUser, error) {
	return s.uc.CreateUser(ctx, dto)
}

// UpdateUser updates an existing privilege user.
func (s *PrivilegeService) UpdateUser(ctx context.Context, dto model.PrivilegeUserDTO, id string) error {
	return s.uc.UpdateUser(ctx, dto, id)
}

// DeleteUser soft-deletes a privilege user.
func (s *PrivilegeService) DeleteUser(ctx context.Context, id string) error {
	return s.uc.DeleteUser(ctx, id)
}

// ResetPassword generates a new random password for a user and returns it.
func (s *PrivilegeService) ResetPassword(ctx context.Context, id string) (string, error) {
	return s.uc.ResetPassword(ctx, id)
}

// UpdatePassword changes a user's password after verifying the old one.
func (s *PrivilegeService) UpdatePassword(ctx context.Context, dto model.PasswordUpdateDTO) error {
	return s.uc.UpdatePassword(ctx, dto)
}

// SetPassword sets a user's password directly (admin operation).
func (s *PrivilegeService) SetPassword(ctx context.Context, id, password string) error {
	return s.uc.SetPassword(ctx, id, password)
}

// PageUsers returns a paginated list of users.
func (s *PrivilegeService) PageUsers(ctx context.Context, query model.PrivilegeUserPageQuery) ([]*model.PrivilegeUserVO, int64, error) {
	return s.uc.PageUsers(ctx, query)
}

// GetUserByID returns a single user by primary key.
func (s *PrivilegeService) GetUserByID(ctx context.Context, id string) (*model.PrivilegeUserVO, error) {
	return s.uc.GetUserByID(ctx, id)
}

// GetUserByCode returns a single user by business code.
func (s *PrivilegeService) GetUserByCode(ctx context.Context, code string) (*model.PrivilegeUserVO, error) {
	return s.uc.GetUserByCode(ctx, code)
}

// ──────────────────────────── Role ────────────────────────────

// CreateRole creates a new privilege role.
func (s *PrivilegeService) CreateRole(ctx context.Context, dto model.PrivilegeRoleDTO) (*model.PrivilegeRole, error) {
	return s.uc.CreateRole(ctx, dto)
}

// UpdateRole updates an existing privilege role.
func (s *PrivilegeService) UpdateRole(ctx context.Context, dto model.PrivilegeRoleDTO, id string) error {
	return s.uc.UpdateRole(ctx, dto, id)
}

// DeleteRole soft-deletes a privilege role.
func (s *PrivilegeService) DeleteRole(ctx context.Context, id string) error {
	return s.uc.DeleteRole(ctx, id)
}

// PageRoles returns a paginated list of roles.
func (s *PrivilegeService) PageRoles(ctx context.Context, query model.PrivilegeRoleQuery) ([]*model.PrivilegeRole, int64, error) {
	return s.uc.PageRoles(ctx, query)
}

// GetRoleByID returns a single role in VO form.
func (s *PrivilegeService) GetRoleByID(ctx context.Context, id string) (*model.PrivilegeRoleVO, error) {
	return s.uc.GetRoleByID(ctx, id)
}

// GetRoleAcls returns the ACL entries for a role.
func (s *PrivilegeService) GetRoleAcls(ctx context.Context, roleID string) ([]model.RoleAclVO, error) {
	return s.uc.GetRoleAcls(ctx, roleID)
}

// GetRolesByCompanyID returns all roles for the given company.
func (s *PrivilegeService) GetRolesByCompanyID(ctx context.Context, companyID int64) ([]*model.PrivilegeRole, error) {
	return s.uc.GetRolesByCompanyID(ctx, companyID)
}

// ──────────────────────────── User-Role ────────────────────────────

// GetUserRolesByUserID returns all user-role associations for a user.
func (s *PrivilegeService) GetUserRolesByUserID(ctx context.Context, userID string) ([]*model.PrivilegeUserRoleVO, error) {
	return s.uc.GetUserRolesByUserID(ctx, userID)
}

// GetUserRolesByRoleID returns all user-role associations for a role.
func (s *PrivilegeService) GetUserRolesByRoleID(ctx context.Context, roleID string) ([]*model.PrivilegeUserRoleVO, error) {
	return s.uc.GetUserRolesByRoleID(ctx, roleID)
}

// SaveUserRoles assigns a single role to a user.
func (s *PrivilegeService) SaveUserRoles(ctx context.Context, dto model.PrivilegeUserRoleDTO) error {
	return s.uc.SaveUserRoles(ctx, dto)
}

// DeleteUserRoles removes a single role from a user.
func (s *PrivilegeService) DeleteUserRoles(ctx context.Context, dto model.PrivilegeUserRoleDTO) error {
	return s.uc.DeleteUserRoles(ctx, dto)
}

// BatchSaveUserRoles replaces all roles for a user.
func (s *PrivilegeService) BatchSaveUserRoles(ctx context.Context, dto model.UserRoleBatchDTO) error {
	return s.uc.BatchSaveUserRoles(ctx, dto)
}

// BatchDeleteUserRoles removes a set of roles from a user.
func (s *PrivilegeService) BatchDeleteUserRoles(ctx context.Context, dto model.UserRoleBatchDTO) error {
	return s.uc.BatchDeleteUserRoles(ctx, dto)
}

// ──────────────────────────── Module / ACL ────────────────────────────

// GetModuleTree returns the full module tree.
func (s *PrivilegeService) GetModuleTree(ctx context.Context) ([]*model.ModuleTreeVO, error) {
	return s.uc.GetModuleTree(ctx)
}

// SaveACLs replaces all ACL entries for a role.
func (s *PrivilegeService) SaveACLs(ctx context.Context, dtos []model.PrivilegeAclDTO) error {
	return s.uc.SaveACLs(ctx, dtos)
}

// SaveModuleACL saves a single ACL module entry.
func (s *PrivilegeService) SaveModuleACL(ctx context.Context, dto model.PrivilegeAclDTO) error {
	return s.uc.SaveModuleACL(ctx, dto)
}

// ──────────────────────────── Department ────────────────────────────

// CreateDepartment creates a new department.
func (s *PrivilegeService) CreateDepartment(ctx context.Context, dto model.PrivilegeDepartmentDTO) (*model.PrivilegeDepartment, error) {
	return s.uc.CreateDepartment(ctx, dto)
}

// UpdateDepartment updates an existing department.
func (s *PrivilegeService) UpdateDepartment(ctx context.Context, dto model.PrivilegeDepartmentDTO, id string) error {
	return s.uc.UpdateDepartment(ctx, dto, id)
}

// DeleteDepartment soft-deletes a department.
func (s *PrivilegeService) DeleteDepartment(ctx context.Context, id string) error {
	return s.uc.DeleteDepartment(ctx, id)
}

// PageDepartments returns a paginated list of departments.
func (s *PrivilegeService) PageDepartments(ctx context.Context, page, size int, name, code, companyID string) ([]*model.PrivilegeDepartment, int64, error) {
	return s.uc.PageDepartments(ctx, page, size, name, code, companyID)
}

// ──────────────────────────── Company ────────────────────────────

// CreateCompany creates a new company.
func (s *PrivilegeService) CreateCompany(ctx context.Context, dto model.PrivilegeCompanyDTO) (*model.PrivilegeCompany, error) {
	return s.uc.CreateCompany(ctx, dto)
}

// UpdateCompany updates an existing company.
func (s *PrivilegeService) UpdateCompany(ctx context.Context, dto model.PrivilegeCompanyDTO, id string) error {
	return s.uc.UpdateCompany(ctx, dto, id)
}

// DeleteCompany soft-deletes a company.
func (s *PrivilegeService) DeleteCompany(ctx context.Context, id string) error {
	return s.uc.DeleteCompany(ctx, id)
}

// PageCompanies returns a paginated list of companies.
func (s *PrivilegeService) PageCompanies(ctx context.Context, query model.PrivilegeCompanyQuery) ([]*model.PrivilegeCompany, int64, error) {
	return s.uc.PageCompanies(ctx, query)
}

// ──────────────────────────── Employee ────────────────────────────

// CreateEmployee creates a new employee binding.
func (s *PrivilegeService) CreateEmployee(ctx context.Context, dto model.PrivilegeEmployeeDTO) (*model.PrivilegeEmployee, error) {
	return s.uc.CreateEmployee(ctx, dto)
}

// UpdateEmployee updates an existing employee binding.
func (s *PrivilegeService) UpdateEmployee(ctx context.Context, dto model.PrivilegeEmployeeDTO, id string) error {
	return s.uc.UpdateEmployee(ctx, dto, id)
}

// DeleteEmployee soft-deletes an employee binding.
func (s *PrivilegeService) DeleteEmployee(ctx context.Context, id string) error {
	return s.uc.DeleteEmployee(ctx, id)
}

// PageEmployees returns a paginated list of employees.
func (s *PrivilegeService) PageEmployees(ctx context.Context, page, size int) ([]*model.PrivilegeEmployeeVO, int64, error) {
	return s.uc.PageEmployees(ctx, page, size)
}

// ──────────────────────────── Dictionary ────────────────────────────

// CreateDictionary creates a new dictionary entry.
func (s *PrivilegeService) CreateDictionary(ctx context.Context, dto model.PrivilegeDictionaryDTO) (*model.PrivilegeDictionary, error) {
	return s.uc.CreateDictionary(ctx, dto)
}

// UpdateDictionary updates an existing dictionary entry.
func (s *PrivilegeService) UpdateDictionary(ctx context.Context, dto model.PrivilegeDictionaryDTO, id string) error {
	return s.uc.UpdateDictionary(ctx, dto, id)
}

// DeleteDictionary soft-deletes a dictionary entry.
func (s *PrivilegeService) DeleteDictionary(ctx context.Context, id string) error {
	return s.uc.DeleteDictionary(ctx, id)
}

// PageDictionaries returns a paginated list of dictionaries.
func (s *PrivilegeService) PageDictionaries(ctx context.Context, page, size int) ([]*model.PrivilegeDictionary, int64, error) {
	return s.uc.PageDictionaries(ctx, page, size)
}

// ──────────────────────────── Pvalue ────────────────────────────

// CreatePvalue creates a new permission value.
func (s *PrivilegeService) CreatePvalue(ctx context.Context, dto model.PrivilegePvalueDTO) (*model.PrivilegePvalue, error) {
	return s.uc.CreatePvalue(ctx, dto)
}

// UpdatePvalue updates an existing permission value.
func (s *PrivilegeService) UpdatePvalue(ctx context.Context, dto model.PrivilegePvalueDTO, id string) error {
	return s.uc.UpdatePvalue(ctx, dto, id)
}

// DeletePvalue soft-deletes a permission value.
func (s *PrivilegeService) DeletePvalue(ctx context.Context, id string) error {
	return s.uc.DeletePvalue(ctx, id)
}

// PagePvalues returns a paginated list of permission values.
func (s *PrivilegeService) PagePvalues(ctx context.Context, query model.PrivilegePvalueQuery) ([]*model.PrivilegePvalue, int64, error) {
	return s.uc.PagePvalues(ctx, query)
}

// ──────────────────────────── LoginLog ────────────────────────────

// PageLoginLogs returns a paginated list of login logs.
func (s *PrivilegeService) PageLoginLogs(ctx context.Context, page, size int) ([]*model.PrivilegeLoginLog, int64, error) {
	return s.uc.PageLoginLogs(ctx, page, size)
}

// ──────────────────────────── Tree / Sync ────────────────────────────

// GetOrgTree returns the full organization tree.
func (s *PrivilegeService) GetOrgTree(ctx context.Context) ([]*model.OrganizationTreeVO, error) {
	return s.uc.GetOrgTree(ctx)
}

// GetDeptTree returns the department tree for a specific company.
func (s *PrivilegeService) GetDeptTree(ctx context.Context, companyID string) ([]*model.OrganizationTreeVO, error) {
	return s.uc.GetDeptTree(ctx, companyID)
}

// SyncDepartments batch-syncs department data.
func (s *PrivilegeService) SyncDepartments(ctx context.Context, dtos []model.PrivilegeDepartmentDTO) error {
	return s.uc.SyncDepartments(ctx, dtos)
}

// SyncEmployees batch-syncs employee bindings.
func (s *PrivilegeService) SyncEmployees(ctx context.Context, dtos []model.PrivilegeEmployeeDTO) error {
	return s.uc.SyncEmployees(ctx, dtos)
}
