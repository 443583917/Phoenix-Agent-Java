package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/phoenix-agent-go/internal/dao/cache"
	"github.com/phoenix-agent-go/internal/domain/privilege"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"github.com/phoenix-agent-go/infra/id"
)

// ──────────────────────────── AppError ────────────────────────────

// AppError is a business-logic error with a numeric code and a
// human-readable message.  Code numbers match the Java side.
type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string { return e.Msg }

// Login-specific errors — match Java privilege-api exactly.
var (
	ErrInvalidCredentials = &AppError{Code: 23003, Msg: "用户名或密码错误"}
	ErrUserDisabled       = &AppError{Code: 23004, Msg: "用户已被禁用"}
	ErrPasswordWrong      = &AppError{Code: 23007, Msg: "密码错误"}
)

// General privilege errors.
var (
	ErrUserNotFound     = &AppError{Code: 401006, Msg: "用户不存在"}
	ErrUsernameExists   = &AppError{Code: 401003, Msg: "用户名已存在"}
	ErrMobileExists     = &AppError{Code: 401004, Msg: "手机号已存在"}
	ErrOldPasswordWrong = &AppError{Code: 401005, Msg: "原密码错误"}
	ErrRoleNotFound     = &AppError{Code: 402001, Msg: "角色不存在"}
	ErrDeptNotFound     = &AppError{Code: 403001, Msg: "部门不存在"}
	ErrCompanyNotFound  = &AppError{Code: 404001, Msg: "公司不存在"}
)

// ──────────────────────────── Usecase ────────────────────────────

// PrivilegeUsecase orchestrates privilege domain operations,
// coordinating repository access, domain logic, and cache.
type PrivilegeUsecase struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	moduleRepo   repository.ModuleRepository
	aclRepo      repository.ACLRepository
	deptRepo     repository.DepartmentRepository
	companyRepo  repository.CompanyRepository
	employeeRepo repository.EmployeeRepository
	dictRepo     repository.DictionaryRepository
	pvalueRepo   repository.PvalueRepository
	loginLogRepo repository.LoginLogRepository
	cache        *cache.PrivilegeCache
}

// NewPrivilegeUsecase constructs a PrivilegeUsecase with all required
// repositories and the shared privilege cache.
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

// ──────────────────────────── helpers ────────────────────────────

func genID() string {
	return strconv.FormatUint(id.MustGenerateID(), 10)
}

// ──────────────────────────── User ────────────────────────────

// Login authenticates a user.  Error codes match the Java privilege-api
// exactly: 23003 for bad credentials, 23004 for disabled user, 23007 for
// wrong password.  The returned LoginUserInfoVO.Token is left empty — the
// caller (handler) is responsible for generating a JWT.
func (u *PrivilegeUsecase) Login(ctx context.Context, dto model.LoginInfoDTO, ip string) (*model.LoginUserInfoVO, error) {
	// 1. Find user by username.
	user, err := u.userRepo.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials // 23003
	}

	// 2. Check disabled — status == 1 means DISABLED (inverted).
	if user.Status == 1 {
		return nil, ErrUserDisabled // 23004
	}

	// 3. Verify password: MD5("phoenix" + plainPassword).
	if !privilege.CheckPassword(user.Password, dto.Password) {
		return nil, ErrPasswordWrong // 23007
	}

	// 4. Save login log.
	_ = u.loginLogRepo.Create(ctx, &model.PrivilegeLoginLog{
		UserID:    user.ID,
		Username:  user.Username,
		LoginIP:   ip,
		LoginTime: time.Now(),
	})

	// 5. Collect role names.
	userRoles, _ := u.userRoleRepo.FindByUserID(ctx, user.ID)
	roleNames := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		role, _ := u.roleRepo.FindByID(ctx, ur.RoleID)
		if role != nil {
			roleNames = append(roleNames, role.SN)
		}
	}

	// 6. Collect menus (full module tree).
	menus, _ := u.moduleRepo.Tree(ctx)

	return &model.LoginUserInfoVO{
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Roles:    roleNames,
		Menus:    convertModuleTreeToUserMenus(menus),
	}, nil
}

// Logout is a no-op placeholder. Token invalidation is handled by the
// auth middleware / handler layer.
func (u *PrivilegeUsecase) Logout(ctx context.Context, userID string) error {
	// No server-side session to destroy; tokens are stateless.
	return nil
}

// GetUserMenus returns the full module tree for a user. In a more complete
// implementation this would be filtered by the user's role ACLs.
func (u *PrivilegeUsecase) GetUserMenus(ctx context.Context, userID string) ([]model.UserMenuVO, error) {
	tree, err := u.moduleRepo.Tree(ctx)
	if err != nil {
		return nil, err
	}
	return convertModuleTreeToUserMenus(tree), nil
}

// CreateUser creates a new privilege user.  Checks for duplicate username
// and mobile before persisting.
func (u *PrivilegeUsecase) CreateUser(ctx context.Context, dto model.PrivilegeUserDTO) (*model.PrivilegeUser, error) {
	existing, _ := u.userRepo.FindByUsername(ctx, dto.Username)
	if existing != nil {
		return nil, ErrUsernameExists
	}
	if dto.Mobile != "" {
		existing, _ = u.userRepo.FindByMobile(ctx, dto.Mobile)
		if existing != nil {
			return nil, ErrMobileExists
		}
	}

	hashed, _ := privilege.HashPassword("123456") // default password
	user := &model.PrivilegeUser{
		BaseModel:  model.BaseModel{ID: genID()},
		EmployeeID: dto.EmployeeID,
		Code:       dto.Code,
		RealName:   dto.RealName,
		Username:   dto.Username,
		Password:   hashed,
		Mobile:     dto.Mobile,
		Email:      dto.Email,
		CompanyID:  dto.CompanyID,
		DeptID:     dto.DeptID,
		Status:     dto.Status,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	_ = u.cache.SetUser(ctx, user)
	return user, nil
}

// UpdateUser updates an existing privilege user identified by its ID.
func (u *PrivilegeUsecase) UpdateUser(ctx context.Context, dto model.PrivilegeUserDTO, id string) error {
	existing, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	existing.EmployeeID = dto.EmployeeID
	existing.Code = dto.Code
	existing.RealName = dto.RealName
	existing.Mobile = dto.Mobile
	existing.Email = dto.Email
	existing.CompanyID = dto.CompanyID
	existing.DeptID = dto.DeptID
	existing.Status = dto.Status

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return err
	}
	_ = u.cache.InvalidateUser(ctx, id)
	return nil
}

// DeleteUser soft-deletes a privilege user.
func (u *PrivilegeUsecase) DeleteUser(ctx context.Context, id string) error {
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return err
	}
	_ = u.cache.InvalidateUser(ctx, id)
	_ = u.cache.InvalidateUserRoles(ctx, id)
	return nil
}

// ResetPassword generates a random 8-character password, hashes it, updates
// the user, and returns the plaintext password.
func (u *PrivilegeUsecase) ResetPassword(ctx context.Context, id string) (string, error) {
	existing, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", ErrUserNotFound
	}

	newPass := privilege.GenerateRandomPassword()
	hashed, _ := privilege.HashPassword(newPass)
	existing.Password = hashed

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return "", err
	}
	_ = u.cache.InvalidateUser(ctx, id)
	return newPass, nil
}

// UpdatePassword changes a user's password after verifying the old one.
func (u *PrivilegeUsecase) UpdatePassword(ctx context.Context, dto model.PasswordUpdateDTO) error {
	existing, err := u.userRepo.FindByID(ctx, dto.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	if !privilege.CheckPassword(existing.Password, dto.OldPassword) {
		return ErrOldPasswordWrong
	}

	hashed, _ := privilege.HashPassword(dto.NewPassword)
	existing.Password = hashed

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return err
	}
	_ = u.cache.InvalidateUser(ctx, dto.UserID)
	return nil
}

// SetPassword sets a user's password directly (admin operation, no old-password check).
func (u *PrivilegeUsecase) SetPassword(ctx context.Context, id, password string) error {
	existing, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	hashed, _ := privilege.HashPassword(password)
	existing.Password = hashed

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return err
	}
	_ = u.cache.InvalidateUser(ctx, id)
	return nil
}

// PageUsers returns a paginated list of users in VO form.
func (u *PrivilegeUsecase) PageUsers(ctx context.Context, query model.PrivilegeUserPageQuery) ([]*model.PrivilegeUserVO, int64, error) {
	list, total, err := u.userRepo.Page(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]*model.PrivilegeUserVO, len(list))
	for i, user := range list {
		vos[i] = userToVO(user)
	}
	return vos, total, nil
}

// GetUserByID returns a single user in VO form.
func (u *PrivilegeUsecase) GetUserByID(ctx context.Context, id string) (*model.PrivilegeUserVO, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return userToVO(user), nil
}

// GetUserByCode returns a single user by their code.
func (u *PrivilegeUsecase) GetUserByCode(ctx context.Context, code string) (*model.PrivilegeUserVO, error) {
	user, err := u.userRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return userToVO(user), nil
}

// ──────────────────────────── Role ────────────────────────────

// CreateRole creates a new privilege role.
func (u *PrivilegeUsecase) CreateRole(ctx context.Context, dto model.PrivilegeRoleDTO) (*model.PrivilegeRole, error) {
	role := &model.PrivilegeRole{
		BaseModel:  model.BaseModel{ID: genID()},
		Name:       dto.Name,
		SN:         dto.SN,
		RoleLevel:  dto.RoleLevel,
		Note:       dto.Note,
		ValidState: dto.ValidState,
		CompanyID:  dto.CompanyID,
		SystemID:   dto.SystemID,
	}
	if role.ValidState == 0 {
		role.ValidState = 1
	}
	if err := u.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	_ = u.cache.SetRole(ctx, role)
	return role, nil
}

// UpdateRole updates an existing privilege role.
func (u *PrivilegeUsecase) UpdateRole(ctx context.Context, dto model.PrivilegeRoleDTO, id string) error {
	existing, err := u.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	existing.Name = dto.Name
	existing.SN = dto.SN
	existing.RoleLevel = dto.RoleLevel
	existing.Note = dto.Note
	existing.ValidState = dto.ValidState
	existing.CompanyID = dto.CompanyID
	existing.SystemID = dto.SystemID

	if err := u.roleRepo.Update(ctx, existing); err != nil {
		return err
	}
	_ = u.cache.InvalidateRole(ctx, id)
	return nil
}

// DeleteRole soft-deletes a privilege role.
func (u *PrivilegeUsecase) DeleteRole(ctx context.Context, id string) error {
	if err := u.roleRepo.Delete(ctx, id); err != nil {
		return err
	}
	_ = u.cache.InvalidateRole(ctx, id)
	return nil
}

// PageRoles returns a paginated list of roles.
func (u *PrivilegeUsecase) PageRoles(ctx context.Context, query model.PrivilegeRoleQuery) ([]*model.PrivilegeRole, int64, error) {
	return u.roleRepo.Page(ctx, query)
}

// GetRoleByID returns a single role in VO form.
func (u *PrivilegeUsecase) GetRoleByID(ctx context.Context, id string) (*model.PrivilegeRoleVO, error) {
	role, err := u.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return &model.PrivilegeRoleVO{
		BaseModel:  role.BaseModel,
		Name:       role.Name,
		SN:         role.SN,
		RoleLevel:  role.RoleLevel,
		Note:       role.Note,
		ValidState: role.ValidState,
		CompanyID:  role.CompanyID,
		SystemID:   role.SystemID,
	}, nil
}

// GetRoleAcls returns the ACL entries for a role, decorated with module names.
func (u *PrivilegeUsecase) GetRoleAcls(ctx context.Context, roleID string) ([]model.RoleAclVO, error) {
	acls, err := u.aclRepo.FindByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	result := make([]model.RoleAclVO, len(acls))
	for i, acl := range acls {
		moduleName := ""
		if mod, _ := u.moduleRepo.FindByID(ctx, acl.ModuleID); mod != nil {
			moduleName = mod.Name
		}
		result[i] = model.RoleAclVO{
			ModuleID:   acl.ModuleID,
			ModuleName: moduleName,
			Permission: acl.Permission,
			Checked:    true,
		}
	}
	return result, nil
}

// ──────────────────────────── User-Role ────────────────────────────

// SaveUserRoles assigns a single role to a user.
func (u *PrivilegeUsecase) SaveUserRoles(ctx context.Context, dto model.PrivilegeUserRoleDTO) error {
	if err := u.userRoleRepo.SaveBatch(ctx, dto.UserID, []string{dto.RoleID}); err != nil {
		return err
	}
	_ = u.cache.InvalidateUserRoles(ctx, dto.UserID)
	return nil
}

// DeleteUserRoles removes a single role from a user.
func (u *PrivilegeUsecase) DeleteUserRoles(ctx context.Context, dto model.PrivilegeUserRoleDTO) error {
	if err := u.userRoleRepo.DeleteBatch(ctx, dto.UserID, []string{dto.RoleID}); err != nil {
		return err
	}
	_ = u.cache.InvalidateUserRoles(ctx, dto.UserID)
	return nil
}

// BatchSaveUserRoles replaces all roles for a user with the given list.
func (u *PrivilegeUsecase) BatchSaveUserRoles(ctx context.Context, dto model.UserRoleBatchDTO) error {
	if err := u.userRoleRepo.SaveBatch(ctx, dto.UserID, dto.RoleIDs); err != nil {
		return err
	}
	_ = u.cache.InvalidateUserRoles(ctx, dto.UserID)
	return nil
}

// BatchDeleteUserRoles removes a set of roles from a user.
func (u *PrivilegeUsecase) BatchDeleteUserRoles(ctx context.Context, dto model.UserRoleBatchDTO) error {
	if err := u.userRoleRepo.DeleteBatch(ctx, dto.UserID, dto.RoleIDs); err != nil {
		return err
	}
	_ = u.cache.InvalidateUserRoles(ctx, dto.UserID)
	return nil
}

// ──────────────────────────── Module / ACL ────────────────────────────

// GetModuleTree returns the full module tree sorted by sort order.
func (u *PrivilegeUsecase) GetModuleTree(ctx context.Context) ([]*model.ModuleTreeVO, error) {
	return u.moduleRepo.Tree(ctx)
}

// SaveACLs replaces all ACL entries for the role of the first ACL in the list.
func (u *PrivilegeUsecase) SaveACLs(ctx context.Context, dtos []model.PrivilegeAclDTO) error {
	acls := make([]*model.PrivilegeAcl, 0, len(dtos))
	for _, d := range dtos {
		acls = append(acls, &model.PrivilegeAcl{
			BaseModel:  model.BaseModel{ID: genID()},
			RoleID:     d.RoleID,
			ModuleID:   d.ModuleID,
			Permission: d.Permission,
			ReleaseID:  d.ReleaseID,
		})
	}
	return u.aclRepo.SaveAll(ctx, acls)
}

// SaveModuleACL saves a single ACL module entry.
func (u *PrivilegeUsecase) SaveModuleACL(ctx context.Context, dto model.PrivilegeAclDTO) error {
	acl := &model.PrivilegeAcl{
		BaseModel:  model.BaseModel{ID: genID()},
		RoleID:     dto.RoleID,
		ModuleID:   dto.ModuleID,
		Permission: dto.Permission,
		ReleaseID:  dto.ReleaseID,
	}
	return u.aclRepo.SaveModule(ctx, acl)
}

// ──────────────────────────── Department ────────────────────────────

// CreateDepartment creates a new department.
func (u *PrivilegeUsecase) CreateDepartment(ctx context.Context, dto model.PrivilegeDepartmentDTO) (*model.PrivilegeDepartment, error) {
	dept := &model.PrivilegeDepartment{
		BaseModel: model.BaseModel{ID: genID()},
		Name:      dto.Name,
		Code:      dto.Code,
		CompanyID: dto.CompanyID,
		Sort:      dto.Sort,
	}
	if dto.PID != "" {
		dept.PID = &dto.PID
	}
	if err := u.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// UpdateDepartment updates an existing department.
func (u *PrivilegeUsecase) UpdateDepartment(ctx context.Context, dto model.PrivilegeDepartmentDTO, id string) error {
	existing, err := u.deptRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrDeptNotFound
	}

	existing.Name = dto.Name
	existing.Code = dto.Code
	existing.CompanyID = dto.CompanyID
	existing.Sort = dto.Sort
	if dto.PID != "" {
		existing.PID = &dto.PID
	} else {
		existing.PID = nil
	}
	return u.deptRepo.Update(ctx, existing)
}

// DeleteDepartment soft-deletes a department.
func (u *PrivilegeUsecase) DeleteDepartment(ctx context.Context, id string) error {
	return u.deptRepo.Delete(ctx, id)
}

// PageDepartments returns a paginated list of departments.
func (u *PrivilegeUsecase) PageDepartments(ctx context.Context, page, size int, name, code, companyID string) ([]*model.PrivilegeDepartment, int64, error) {
	filter := &model.PrivilegeDepartment{}
	if name != "" {
		filter.Name = name
	}
	if code != "" {
		filter.Code = code
	}
	if companyID != "" {
		filter.CompanyID = companyID
	}
	return u.deptRepo.Page(ctx, page, size, filter)
}

// ──────────────────────────── Company ────────────────────────────

// CreateCompany creates a new company.
func (u *PrivilegeUsecase) CreateCompany(ctx context.Context, dto model.PrivilegeCompanyDTO) (*model.PrivilegeCompany, error) {
	company := &model.PrivilegeCompany{
		BaseModel: model.BaseModel{ID: genID()},
		CName:     dto.CName,
		EName:     dto.EName,
		Code:      dto.Code,
		SN:        dto.SN,
		Manager:   dto.Manager,
		Note:      dto.Note,
	}
	if err := u.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

// UpdateCompany updates an existing company.
func (u *PrivilegeUsecase) UpdateCompany(ctx context.Context, dto model.PrivilegeCompanyDTO, id string) error {
	existing, err := u.companyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCompanyNotFound
	}

	existing.CName = dto.CName
	existing.EName = dto.EName
	existing.Code = dto.Code
	existing.SN = dto.SN
	existing.Manager = dto.Manager
	existing.Note = dto.Note
	return u.companyRepo.Update(ctx, existing)
}

// DeleteCompany soft-deletes a company.
func (u *PrivilegeUsecase) DeleteCompany(ctx context.Context, id string) error {
	return u.companyRepo.Delete(ctx, id)
}

// PageCompanies returns a paginated list of companies.
func (u *PrivilegeUsecase) PageCompanies(ctx context.Context, query model.PrivilegeCompanyQuery) ([]*model.PrivilegeCompany, int64, error) {
	return u.companyRepo.Page(ctx, query)
}

// ──────────────────────────── Employee ────────────────────────────

// CreateEmployee creates a new employee binding.
func (u *PrivilegeUsecase) CreateEmployee(ctx context.Context, dto model.PrivilegeEmployeeDTO) (*model.PrivilegeEmployee, error) {
	emp := &model.PrivilegeEmployee{
		BaseModel: model.BaseModel{ID: genID()},
		UserCode:  dto.UserCode,
		EmpCode:   dto.EmpCode,
		DeptID:    dto.DeptID,
	}
	if err := u.employeeRepo.Create(ctx, emp); err != nil {
		return nil, err
	}
	return emp, nil
}

// UpdateEmployee updates an existing employee binding.
func (u *PrivilegeUsecase) UpdateEmployee(ctx context.Context, dto model.PrivilegeEmployeeDTO, id string) error {
	existing, err := u.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return &AppError{Code: 405001, Msg: "员工不存在"}
	}

	existing.UserCode = dto.UserCode
	existing.EmpCode = dto.EmpCode
	existing.DeptID = dto.DeptID
	return u.employeeRepo.Update(ctx, existing)
}

// DeleteEmployee soft-deletes an employee binding.
func (u *PrivilegeUsecase) DeleteEmployee(ctx context.Context, id string) error {
	return u.employeeRepo.Delete(ctx, id)
}

// PageEmployees returns a paginated list of employees in VO form.
func (u *PrivilegeUsecase) PageEmployees(ctx context.Context, page, size int) ([]*model.PrivilegeEmployeeVO, int64, error) {
	list, total, err := u.employeeRepo.Page(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]*model.PrivilegeEmployeeVO, len(list))
	for i, emp := range list {
		vos[i] = employeeToVO(emp)
	}
	return vos, total, nil
}

// ──────────────────────────── Dictionary ────────────────────────────

// CreateDictionary creates a new dictionary entry.
func (u *PrivilegeUsecase) CreateDictionary(ctx context.Context, dto model.PrivilegeDictionaryDTO) (*model.PrivilegeDictionary, error) {
	dict := &model.PrivilegeDictionary{
		BaseModel: model.BaseModel{ID: genID()},
		Code:      dto.Code,
		Name:      dto.Name,
		SystemSN:  dto.SystemSN,
		SystemID:  dto.SystemID,
		Sort:      dto.Sort,
	}
	if dto.PCode != "" {
		dict.PCode = &dto.PCode
	}
	if err := u.dictRepo.Create(ctx, dict); err != nil {
		return nil, err
	}
	return dict, nil
}

// UpdateDictionary updates an existing dictionary entry.
func (u *PrivilegeUsecase) UpdateDictionary(ctx context.Context, dto model.PrivilegeDictionaryDTO, id string) error {
	existing, err := u.dictRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return &AppError{Code: 406001, Msg: "字典不存在"}
	}

	existing.Code = dto.Code
	existing.Name = dto.Name
	existing.SystemSN = dto.SystemSN
	existing.SystemID = dto.SystemID
	existing.Sort = dto.Sort
	if dto.PCode != "" {
		existing.PCode = &dto.PCode
	} else {
		existing.PCode = nil
	}
	return u.dictRepo.Update(ctx, existing)
}

// DeleteDictionary soft-deletes a dictionary entry.
func (u *PrivilegeUsecase) DeleteDictionary(ctx context.Context, id string) error {
	return u.dictRepo.Delete(ctx, id)
}

// PageDictionaries returns a paginated list of dictionaries.
func (u *PrivilegeUsecase) PageDictionaries(ctx context.Context, page, size int) ([]*model.PrivilegeDictionary, int64, error) {
	return u.dictRepo.Page(ctx, page, size)
}

// ──────────────────────────── Pvalue ────────────────────────────

// CreatePvalue creates a new permission value.
func (u *PrivilegeUsecase) CreatePvalue(ctx context.Context, dto model.PrivilegePvalueDTO) (*model.PrivilegePvalue, error) {
	pv := &model.PrivilegePvalue{
		BaseModel: model.BaseModel{ID: genID()},
		Code:      dto.Code,
		Name:      dto.Name,
		SystemID:  dto.SystemID,
	}
	if err := u.pvalueRepo.Create(ctx, pv); err != nil {
		return nil, err
	}
	return pv, nil
}

// UpdatePvalue updates an existing permission value.
func (u *PrivilegeUsecase) UpdatePvalue(ctx context.Context, dto model.PrivilegePvalueDTO, id string) error {
	existing, err := u.pvalueRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return &AppError{Code: 407001, Msg: "权限值不存在"}
	}

	existing.Code = dto.Code
	existing.Name = dto.Name
	existing.SystemID = dto.SystemID
	return u.pvalueRepo.Update(ctx, existing)
}

// DeletePvalue soft-deletes a permission value.
func (u *PrivilegeUsecase) DeletePvalue(ctx context.Context, id string) error {
	return u.pvalueRepo.Delete(ctx, id)
}

// PagePvalues returns a paginated list of permission values.
func (u *PrivilegeUsecase) PagePvalues(ctx context.Context, query model.PrivilegePvalueQuery) ([]*model.PrivilegePvalue, int64, error) {
	return u.pvalueRepo.Page(ctx, query)
}

// ──────────────────────────── LoginLog ────────────────────────────

// PageLoginLogs returns a paginated list of login logs.
func (u *PrivilegeUsecase) PageLoginLogs(ctx context.Context, page, size int) ([]*model.PrivilegeLoginLog, int64, error) {
	return u.loginLogRepo.Page(ctx, page, size)
}

// ──────────────────────────── Tree / Sync ────────────────────────────

// GetOrgTree returns the full organization tree (companies with nested departments).
func (u *PrivilegeUsecase) GetOrgTree(ctx context.Context) ([]*model.OrganizationTreeVO, error) {
	return u.deptRepo.OrgTree(ctx)
}

// GetDeptTree returns the department tree for a specific company.
func (u *PrivilegeUsecase) GetDeptTree(ctx context.Context, companyID string) ([]*model.OrganizationTreeVO, error) {
	company, err := u.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	departments, err := u.deptRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	root := &model.OrganizationTreeVO{
		ID:       company.ID,
		Name:     company.CName,
		Code:     company.Code,
		Type:     "company",
		Children: buildDeptTreeNodes(departments, nil),
	}
	return []*model.OrganizationTreeVO{root}, nil
}

// SyncDepartments batch-syncs department data — each item is created if new
// or updated if a matching code already exists in the company.
func (u *PrivilegeUsecase) SyncDepartments(ctx context.Context, dtos []model.PrivilegeDepartmentDTO) error {
	for _, dto := range dtos {
		// Attempt to find by code within the same company.
		all, _, _ := u.deptRepo.Page(ctx, 1, 1000, &model.PrivilegeDepartment{Code: dto.Code, CompanyID: dto.CompanyID})
		if len(all) > 0 {
			existing := all[0]
			existing.Name = dto.Name
			existing.Sort = dto.Sort
			if dto.PID != "" {
				existing.PID = &dto.PID
			}
			if err := u.deptRepo.Update(ctx, existing); err != nil {
				return err
			}
		} else {
			dept := &model.PrivilegeDepartment{
				BaseModel: model.BaseModel{ID: genID()},
				Name:      dto.Name,
				Code:      dto.Code,
				CompanyID: dto.CompanyID,
				Sort:      dto.Sort,
			}
			if dto.PID != "" {
				dept.PID = &dto.PID
			}
			if err := u.deptRepo.Create(ctx, dept); err != nil {
				return err
			}
		}
	}
	return nil
}

// SyncEmployees batch-syncs employee bindings — each item is created if new
// or updated if a matching empCode already exists.
func (u *PrivilegeUsecase) SyncEmployees(ctx context.Context, dtos []model.PrivilegeEmployeeDTO) error {
	for _, dto := range dtos {
		existing, _ := u.employeeRepo.FindByEmpCode(ctx, dto.EmpCode)
		if existing != nil {
			existing.UserCode = dto.UserCode
			existing.DeptID = dto.DeptID
			if err := u.employeeRepo.Update(ctx, existing); err != nil {
				return err
			}
		} else {
			emp := &model.PrivilegeEmployee{
				BaseModel: model.BaseModel{ID: genID()},
				UserCode:  dto.UserCode,
				EmpCode:   dto.EmpCode,
				DeptID:    dto.DeptID,
			}
			if err := u.employeeRepo.Create(ctx, emp); err != nil {
				return err
			}
		}
	}
	return nil
}

// ──────────────────────────── converters ────────────────────────────

func userToVO(user *model.PrivilegeUser) *model.PrivilegeUserVO {
	return &model.PrivilegeUserVO{
		ID:         user.ID,
		EmployeeID: user.EmployeeID,
		Code:       user.Code,
		RealName:   user.RealName,
		Username:   user.Username,
		Tel:        user.Tel,
		Mobile:     user.Mobile,
		Email:      user.Email,
		CompanyID:  user.CompanyID,
		DeptID:     user.DeptID,
		Status:     user.Status,
		UserType:   user.UserType,
		CreateTime: user.CreateTime,
	}
}

func employeeToVO(emp *model.PrivilegeEmployee) *model.PrivilegeEmployeeVO {
	return &model.PrivilegeEmployeeVO{
		PrivilegeEmployee: *emp,
	}
}

func convertModuleTreeToUserMenus(tree []*model.ModuleTreeVO) []model.UserMenuVO {
	if tree == nil {
		return nil
	}
	result := make([]model.UserMenuVO, len(tree))
	for i, node := range tree {
		result[i] = model.UserMenuVO{
			ID:       node.ID,
			Name:     node.Name,
			Code:     node.Code,
			URL:      node.URL,
			Icon:     node.Icon,
			Sort:     node.Sort,
			Type:     node.Type,
			Children: convertModuleTreeChildren(node.Children),
		}
	}
	return result
}

func convertModuleTreeChildren(children []model.ModuleTreeVO) []model.UserMenuVO {
	result := make([]model.UserMenuVO, len(children))
	for i, child := range children {
		result[i] = model.UserMenuVO{
			ID:       child.ID,
			Name:     child.Name,
			Code:     child.Code,
			URL:      child.URL,
			Icon:     child.Icon,
			Sort:     child.Sort,
			Type:     child.Type,
			Children: convertModuleTreeChildren(child.Children),
		}
	}
	return result
}

// buildDeptTreeNodes builds a department tree from a flat list, filtering by pid.
func buildDeptTreeNodes(depts []*model.PrivilegeDepartment, pid *string) []model.OrganizationTreeVO {
	var tree []model.OrganizationTreeVO
	for _, d := range depts {
		if (pid == nil && d.PID == nil) || (pid != nil && d.PID != nil && *d.PID == *pid) {
			node := model.OrganizationTreeVO{
				ID:   d.ID,
				Name: d.Name,
				Code: d.Code,
				Type: "department",
			}
			node.Children = buildDeptTreeNodes(depts, &d.ID)
			tree = append(tree, node)
		}
	}
	return tree
}
