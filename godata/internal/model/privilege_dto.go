package model

// LoginInfoDTO — POST /api/privilege/auth/login
type LoginInfoDTO struct {
	Type        string `json:"type" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaKey  string `json:"captchaKey"`
	CaptchaCode string `json:"captchaCode"`
}

// PasswordUpdateDTO — 密码修改请求
type PasswordUpdateDTO struct {
	UserID      string `json:"userId" binding:"required"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// PrivilegeUserDTO — 用户创建/更新请求
type PrivilegeUserDTO struct {
	ID         string `json:"id"`
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

// PrivilegeUserPageQuery — 用户分页查询 (GET query params)
type PrivilegeUserPageQuery struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Username string `form:"username"`
	RealName string `form:"realName"`
	Status   *int   `form:"status"`
}

// PrivilegeRoleDTO — 角色创建/更新请求
type PrivilegeRoleDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name" binding:"required"`
	SN         string `json:"sn"`
	RoleLevel  string `json:"roleLevel"`
	Note       string `json:"note"`
	ValidState int    `json:"validState"`
	CompanyID  int64  `json:"companyId"`
	SystemID   string `json:"systemId"`
}

// PrivilegeRoleQuery — 角色分页查询
type PrivilegeRoleQuery struct {
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Name      string `json:"name"`
	CompanyID int64  `json:"companyId"`
}

// PrivilegeUserRoleDTO — 用户-角色关联请求
type PrivilegeUserRoleDTO struct {
	UserID string `json:"userId" binding:"required"`
	RoleID string `json:"roleId" binding:"required"`
}

// UserRoleBatchDTO — 批量用户-角色关联请求
type UserRoleBatchDTO struct {
	UserID  string   `json:"userId" binding:"required"`
	RoleIDs []string `json:"roleIds" binding:"required"`
}

// PrivilegeModuleDTO — 菜单/模块创建/更新请求
type PrivilegeModuleDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code"`
	PID      string `json:"pid"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	SystemID string `json:"systemId"`
	Type     int    `json:"type"`
}

// PrivilegeAclDTO — 访问控制列表请求
type PrivilegeAclDTO struct {
	ID         string `json:"id"`
	RoleID     string `json:"roleId"`
	ModuleID   string `json:"moduleId"`
	Permission string `json:"permission"`
	ReleaseID  string `json:"releaseId"`
}

// PrivilegeDepartmentDTO — 部门创建/更新请求
type PrivilegeDepartmentDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code"`
	PID       string `json:"pid"`
	CompanyID string `json:"companyId"`
	Sort      int    `json:"sort"`
}

// PrivilegeCompanyDTO — 公司创建/更新请求
type PrivilegeCompanyDTO struct {
	ID      string `json:"id"`
	CName   string `json:"cname" binding:"required"`
	EName   string `json:"ename"`
	Code    string `json:"code"`
	SN      string `json:"sn"`
	Manager string `json:"manager"`
	Note    string `json:"note"`
}

// PrivilegeCompanyQuery — 公司分页查询
type PrivilegeCompanyQuery struct {
	Page int    `json:"page"`
	Size int    `json:"size"`
	Name string `json:"name"`
}

// PrivilegeEmployeeDTO — 员工绑定请求
type PrivilegeEmployeeDTO struct {
	ID       string `json:"id"`
	UserCode string `json:"userCode" binding:"required"`
	EmpCode  string `json:"empCode" binding:"required"`
	DeptID   string `json:"deptId"`
}

// PrivilegeDictionaryDTO — 字典创建/更新请求
type PrivilegeDictionaryDTO struct {
	ID       string `json:"id"`
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	PCode    string `json:"pcode"`
	SystemSN string `json:"systemSn"`
	SystemID string `json:"systemId"`
	Sort     int    `json:"sort"`
}

// PrivilegePvalueDTO — 权限值创建/更新请求
type PrivilegePvalueDTO struct {
	ID       string `json:"id"`
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	SystemID string `json:"systemId"`
}

// PrivilegePvalueQuery — 权限值分页查询
type PrivilegePvalueQuery struct {
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	SystemID string `json:"systemId"`
}

// PrivilegeLoginLogDTO — 登录日志查询请求
type PrivilegeLoginLogDTO struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	LoginIP  string `json:"loginIp"`
}
