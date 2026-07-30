package db

import (
	"context"
	"errors"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── User ────────────────────────────

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepo{db}
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeUser, error) {
	var user model.PrivilegeUser
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.PrivilegeUser, error) {
	var user model.PrivilegeUser
	err := r.db.WithContext(ctx).Where("username = ? AND del_flag = 0", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepo) FindByCode(ctx context.Context, code string) (*model.PrivilegeUser, error) {
	var user model.PrivilegeUser
	err := r.db.WithContext(ctx).Where("code = ? AND del_flag = 0", code).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepo) FindByMobile(ctx context.Context, mobile string) (*model.PrivilegeUser, error) {
	var user model.PrivilegeUser
	err := r.db.WithContext(ctx).Where("mobile = ? AND del_flag = 0", mobile).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepo) Page(ctx context.Context, query model.PrivilegeUserPageQuery) ([]*model.PrivilegeUser, int64, error) {
	var list []*model.PrivilegeUser
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeUser{}).Where("del_flag = 0")
	if query.Username != "" {
		dbQuery = dbQuery.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.RealName != "" {
		dbQuery = dbQuery.Where("real_name LIKE ?", "%"+query.RealName+"%")
	}
	if query.Status != nil {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	offset := (query.Page - 1) * query.Size
	if err := dbQuery.Offset(offset).Limit(query.Size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *userRepo) Create(ctx context.Context, user *model.PrivilegeUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) Update(ctx context.Context, user *model.PrivilegeUser) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeUser{}).Where("id = ? AND del_flag = 0", user.ID).Updates(user).Error
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeUser{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── Role ────────────────────────────

type roleRepo struct{ db *gorm.DB }

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepo{db}
}

func (r *roleRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeRole, error) {
	var role model.PrivilegeRole
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (r *roleRepo) FindByCompanyID(ctx context.Context, companyID int64) ([]*model.PrivilegeRole, error) {
	var roles []*model.PrivilegeRole
	err := r.db.WithContext(ctx).Where("company_id = ? AND del_flag = 0", companyID).Find(&roles).Error
	return roles, err
}

func (r *roleRepo) Page(ctx context.Context, query model.PrivilegeRoleQuery) ([]*model.PrivilegeRole, int64, error) {
	var list []*model.PrivilegeRole
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeRole{}).Where("del_flag = 0")
	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.CompanyID != 0 {
		dbQuery = dbQuery.Where("company_id = ?", query.CompanyID)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	offset := (query.Page - 1) * query.Size
	if err := dbQuery.Offset(offset).Limit(query.Size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *roleRepo) Create(ctx context.Context, role *model.PrivilegeRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepo) Update(ctx context.Context, role *model.PrivilegeRole) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeRole{}).Where("id = ? AND del_flag = 0", role.ID).Updates(role).Error
}

func (r *roleRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeRole{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── UserRole ────────────────────────────

type userRoleRepo struct{ db *gorm.DB }

func NewUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return &userRoleRepo{db}
}

func (r *userRoleRepo) FindByUserID(ctx context.Context, userID string) ([]*model.PrivilegeUserRole, error) {
	var list []*model.PrivilegeUserRole
	err := r.db.WithContext(ctx).Where("user_id = ? AND del_flag = 0", userID).Find(&list).Error
	return list, err
}

func (r *userRoleRepo) FindByRoleID(ctx context.Context, roleID string) ([]*model.PrivilegeUserRole, error) {
	var list []*model.PrivilegeUserRole
	err := r.db.WithContext(ctx).Where("role_id = ? AND del_flag = 0", roleID).Find(&list).Error
	return list, err
}

func (r *userRoleRepo) SaveBatch(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.PrivilegeUserRole{}).Where("user_id = ? AND del_flag = 0", userID).Update("del_flag", 1).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			ur := &model.PrivilegeUserRole{UserID: userID, RoleID: roleID}
			if err := tx.Create(ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRoleRepo) DeleteBatch(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, roleID := range roleIDs {
			if err := tx.Model(&model.PrivilegeUserRole{}).Where("user_id = ? AND role_id = ? AND del_flag = 0", userID, roleID).Update("del_flag", 1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ──────────────────────────── Module ────────────────────────────

type moduleRepo struct{ db *gorm.DB }

func NewModuleRepository(db *gorm.DB) repository.ModuleRepository {
	return &moduleRepo{db}
}

func (r *moduleRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeModule, error) {
	var m model.PrivilegeModule
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *moduleRepo) FindByPID(ctx context.Context, pid string) ([]*model.PrivilegeModule, error) {
	var list []*model.PrivilegeModule
	err := r.db.WithContext(ctx).Where("pid = ? AND del_flag = 0", pid).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *moduleRepo) FindAll(ctx context.Context) ([]*model.PrivilegeModule, error) {
	var list []*model.PrivilegeModule
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *moduleRepo) Tree(ctx context.Context) ([]*model.ModuleTreeVO, error) {
	var modules []*model.PrivilegeModule
	if err := r.db.WithContext(ctx).Where("del_flag = 0").Order("sort ASC").Find(&modules).Error; err != nil {
		return nil, err
	}
	treeVOs := buildModuleTree(modules, nil)
	result := make([]*model.ModuleTreeVO, len(treeVOs))
	for i := range treeVOs {
		result[i] = &treeVOs[i]
	}
	return result, nil
}

func buildModuleTree(modules []*model.PrivilegeModule, pid *string) []model.ModuleTreeVO {
	var tree []model.ModuleTreeVO
	for _, m := range modules {
		if (pid == nil && m.PID == nil) || (pid != nil && m.PID != nil && *m.PID == *pid) {
			node := model.ModuleTreeVO{
				PrivilegeModule: *m,
				Children:        buildModuleTree(modules, &m.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

func (r *moduleRepo) Create(ctx context.Context, module *model.PrivilegeModule) error {
	return r.db.WithContext(ctx).Create(module).Error
}

func (r *moduleRepo) Update(ctx context.Context, module *model.PrivilegeModule) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeModule{}).Where("id = ? AND del_flag = 0", module.ID).Updates(module).Error
}

// ──────────────────────────── ACL ────────────────────────────

type aclRepo struct{ db *gorm.DB }

func NewACLRepository(db *gorm.DB) repository.ACLRepository {
	return &aclRepo{db}
}

func (r *aclRepo) FindByRoleID(ctx context.Context, roleID string) ([]*model.PrivilegeAcl, error) {
	var list []*model.PrivilegeAcl
	err := r.db.WithContext(ctx).Where("role_id = ? AND del_flag = 0", roleID).Find(&list).Error
	return list, err
}

func (r *aclRepo) SaveAll(ctx context.Context, acls []*model.PrivilegeAcl) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(acls) > 0 {
			roleID := acls[0].RoleID
			if err := tx.Model(&model.PrivilegeAcl{}).Where("role_id = ? AND del_flag = 0", roleID).Update("del_flag", 1).Error; err != nil {
				return err
			}
		}
		for _, acl := range acls {
			if acl.RoleID == "" || acl.ModuleID == "" {
				continue
			}
			if err := tx.Create(acl).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *aclRepo) SaveModule(ctx context.Context, acl *model.PrivilegeAcl) error {
	return r.db.WithContext(ctx).Create(acl).Error
}

func (r *aclRepo) FindByReleaseID(ctx context.Context, releaseID string) ([]*model.PrivilegeAcl, error) {
	var list []*model.PrivilegeAcl
	err := r.db.WithContext(ctx).Where("release_id = ? AND del_flag = 0", releaseID).Find(&list).Error
	return list, err
}

// ──────────────────────────── Department ────────────────────────────

type deptRepo struct{ db *gorm.DB }

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &deptRepo{db}
}

func (r *deptRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeDepartment, error) {
	var dept model.PrivilegeDepartment
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&dept).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dept, err
}

func (r *deptRepo) FindByPID(ctx context.Context, pid string) ([]*model.PrivilegeDepartment, error) {
	var list []*model.PrivilegeDepartment
	err := r.db.WithContext(ctx).Where("pid = ? AND del_flag = 0", pid).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *deptRepo) FindByCompanyID(ctx context.Context, companyID string) ([]*model.PrivilegeDepartment, error) {
	var list []*model.PrivilegeDepartment
	err := r.db.WithContext(ctx).Where("company_id = ? AND del_flag = 0", companyID).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *deptRepo) OrgTree(ctx context.Context) ([]*model.OrganizationTreeVO, error) {
	var companies []*model.PrivilegeCompany
	if err := r.db.WithContext(ctx).Where("del_flag = 0").Find(&companies).Error; err != nil {
		return nil, err
	}

	var departments []*model.PrivilegeDepartment
	if err := r.db.WithContext(ctx).Where("del_flag = 0").Order("sort ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	return buildOrgTree(companies, departments), nil
}

func buildOrgTree(companies []*model.PrivilegeCompany, departments []*model.PrivilegeDepartment) []*model.OrganizationTreeVO {
	var tree []*model.OrganizationTreeVO
	for _, c := range companies {
		node := &model.OrganizationTreeVO{
			ID:   c.ID,
			Name: c.CName,
			Code: c.Code,
			Type: "company",
		}
		node.Children = buildDeptTree(departments, c.ID, nil)
		tree = append(tree, node)
	}
	return tree
}

func buildDeptTree(depts []*model.PrivilegeDepartment, companyID string, pid *string) []model.OrganizationTreeVO {
	var tree []model.OrganizationTreeVO
	for _, d := range depts {
		if d.CompanyID != companyID {
			continue
		}
		if (pid == nil && d.PID == nil) || (pid != nil && d.PID != nil && *d.PID == *pid) {
			node := model.OrganizationTreeVO{
				ID:   d.ID,
				Name: d.Name,
				Code: d.Code,
				Type: "department",
			}
			node.Children = buildDeptTree(depts, companyID, &d.ID)
			tree = append(tree, node)
		}
	}
	return tree
}

func (r *deptRepo) Page(ctx context.Context, page, size int, dept *model.PrivilegeDepartment) ([]*model.PrivilegeDepartment, int64, error) {
	var list []*model.PrivilegeDepartment
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeDepartment{}).Where("del_flag = 0")
	if dept != nil {
		if dept.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+dept.Name+"%")
		}
		if dept.Code != "" {
			dbQuery = dbQuery.Where("code LIKE ?", "%"+dept.Code+"%")
		}
		if dept.CompanyID != "" {
			dbQuery = dbQuery.Where("company_id = ?", dept.CompanyID)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *deptRepo) Create(ctx context.Context, dept *model.PrivilegeDepartment) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *deptRepo) Update(ctx context.Context, dept *model.PrivilegeDepartment) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeDepartment{}).Where("id = ? AND del_flag = 0", dept.ID).Updates(dept).Error
}

func (r *deptRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeDepartment{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── Company ────────────────────────────

type companyRepo struct{ db *gorm.DB }

func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepo{db}
}

func (r *companyRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeCompany, error) {
	var c model.PrivilegeCompany
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *companyRepo) FindByCode(ctx context.Context, code string) (*model.PrivilegeCompany, error) {
	var c model.PrivilegeCompany
	err := r.db.WithContext(ctx).Where("code = ? AND del_flag = 0", code).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *companyRepo) Page(ctx context.Context, query model.PrivilegeCompanyQuery) ([]*model.PrivilegeCompany, int64, error) {
	var list []*model.PrivilegeCompany
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeCompany{}).Where("del_flag = 0")
	if query.Name != "" {
		dbQuery = dbQuery.Where("cname LIKE ? OR ename LIKE ?", "%"+query.Name+"%", "%"+query.Name+"%")
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	offset := (query.Page - 1) * query.Size
	if err := dbQuery.Offset(offset).Limit(query.Size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *companyRepo) Create(ctx context.Context, company *model.PrivilegeCompany) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *companyRepo) Update(ctx context.Context, company *model.PrivilegeCompany) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeCompany{}).Where("id = ? AND del_flag = 0", company.ID).Updates(company).Error
}

func (r *companyRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeCompany{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── Employee ────────────────────────────

type employeeRepo struct{ db *gorm.DB }

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepo{db}
}

func (r *employeeRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeEmployee, error) {
	var emp model.PrivilegeEmployee
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&emp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &emp, err
}

func (r *employeeRepo) FindByEmpCode(ctx context.Context, empCode string) (*model.PrivilegeEmployee, error) {
	var emp model.PrivilegeEmployee
	err := r.db.WithContext(ctx).Where("emp_code = ? AND del_flag = 0", empCode).First(&emp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &emp, err
}

func (r *employeeRepo) Page(ctx context.Context, page, size int) ([]*model.PrivilegeEmployee, int64, error) {
	var list []*model.PrivilegeEmployee
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeEmployee{}).Where("del_flag = 0")

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *employeeRepo) Create(ctx context.Context, emp *model.PrivilegeEmployee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *employeeRepo) Update(ctx context.Context, emp *model.PrivilegeEmployee) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeEmployee{}).Where("id = ? AND del_flag = 0", emp.ID).Updates(emp).Error
}

func (r *employeeRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeEmployee{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── Dictionary ────────────────────────────

type dictionaryRepo struct{ db *gorm.DB }

func NewDictionaryRepository(db *gorm.DB) repository.DictionaryRepository {
	return &dictionaryRepo{db}
}

func (r *dictionaryRepo) FindByID(ctx context.Context, id string) (*model.PrivilegeDictionary, error) {
	var dict model.PrivilegeDictionary
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&dict).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dict, err
}

func (r *dictionaryRepo) FindBySystemSN(ctx context.Context, systemSN string) ([]*model.PrivilegeDictionary, error) {
	var list []*model.PrivilegeDictionary
	err := r.db.WithContext(ctx).Where("system_sn = ? AND del_flag = 0", systemSN).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *dictionaryRepo) FindByPCode(ctx context.Context, pcode string) ([]*model.PrivilegeDictionary, error) {
	var list []*model.PrivilegeDictionary
	err := r.db.WithContext(ctx).Where("pcode = ? AND del_flag = 0", pcode).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *dictionaryRepo) Page(ctx context.Context, page, size int) ([]*model.PrivilegeDictionary, int64, error) {
	var list []*model.PrivilegeDictionary
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeDictionary{}).Where("del_flag = 0")

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *dictionaryRepo) Create(ctx context.Context, dict *model.PrivilegeDictionary) error {
	return r.db.WithContext(ctx).Create(dict).Error
}

func (r *dictionaryRepo) Update(ctx context.Context, dict *model.PrivilegeDictionary) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeDictionary{}).Where("id = ? AND del_flag = 0", dict.ID).Updates(dict).Error
}

func (r *dictionaryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegeDictionary{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── Pvalue ────────────────────────────

type pvalueRepo struct{ db *gorm.DB }

func NewPvalueRepository(db *gorm.DB) repository.PvalueRepository {
	return &pvalueRepo{db}
}

func (r *pvalueRepo) FindByID(ctx context.Context, id string) (*model.PrivilegePvalue, error) {
	var pv model.PrivilegePvalue
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&pv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pv, err
}

func (r *pvalueRepo) Page(ctx context.Context, query model.PrivilegePvalueQuery) ([]*model.PrivilegePvalue, int64, error) {
	var list []*model.PrivilegePvalue
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegePvalue{}).Where("del_flag = 0")
	if query.SystemID != "" {
		dbQuery = dbQuery.Where("system_id = ?", query.SystemID)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	offset := (query.Page - 1) * query.Size
	if err := dbQuery.Offset(offset).Limit(query.Size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *pvalueRepo) Create(ctx context.Context, pv *model.PrivilegePvalue) error {
	return r.db.WithContext(ctx).Create(pv).Error
}

func (r *pvalueRepo) Update(ctx context.Context, pv *model.PrivilegePvalue) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegePvalue{}).Where("id = ? AND del_flag = 0", pv.ID).Updates(pv).Error
}

func (r *pvalueRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.PrivilegePvalue{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── LoginLog ────────────────────────────

type loginLogRepo struct{ db *gorm.DB }

func NewLoginLogRepository(db *gorm.DB) repository.LoginLogRepository {
	return &loginLogRepo{db}
}

func (r *loginLogRepo) Create(ctx context.Context, log *model.PrivilegeLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *loginLogRepo) Page(ctx context.Context, page, size int) ([]*model.PrivilegeLoginLog, int64, error) {
	var list []*model.PrivilegeLoginLog
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.PrivilegeLoginLog{})

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("login_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
