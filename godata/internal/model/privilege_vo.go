package model

import "time"

// CaptchaVO — 验证码响应
type CaptchaVO struct {
	CaptchaKey string `json:"captchaKey"`
	Image      string `json:"image"` // base64
}

// LoginUserInfoVO — 登录成功后的用户信息
type LoginUserInfoVO struct {
	UserID   string       `json:"userId"`
	Username string       `json:"username"`
	RealName string       `json:"realName"`
	Token    string       `json:"token"`
	Roles    []string     `json:"roles"`
	Menus    []UserMenuVO `json:"menus"`
}

// UserMenuVO — 用户菜单树节点
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

// PrivilegeUserVO — 用户列表/详情响应
type PrivilegeUserVO struct {
	ID          string    `json:"id"`
	EmployeeID  string    `json:"employeeId"`
	Code        string    `json:"code"`
	RealName    string    `json:"realName"`
	Username    string    `json:"username"`
	Tel         string    `json:"tel"`
	Mobile      string    `json:"mobile"`
	Email       string    `json:"email"`
	CompanyID   string    `json:"companyId"`
	CompanyName string    `json:"companyName"`
	DeptID      string    `json:"deptId"`
	DeptName    string    `json:"deptName"`
	Status      int       `json:"status"`
	UserType    int       `json:"userType"`
	CreateTime  time.Time `json:"createTime"`
}

// PrivilegeRoleVO — 角色列表/详情响应
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

// RoleAclVO — 角色权限项
type RoleAclVO struct {
	ModuleID   string `json:"moduleId"`
	ModuleName string `json:"moduleName"`
	Permission string `json:"permission"`
	Checked    bool   `json:"checked"`
}

// ModuleTreeVO — 模块树节点
type ModuleTreeVO struct {
	PrivilegeModule
	Children []ModuleTreeVO `json:"children,omitempty"`
}

// OrganizationTreeVO — 组织树节点（公司/部门）
type OrganizationTreeVO struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Code     string               `json:"code"`
	Type     string               `json:"type"` // "company" | "department"
	Children []OrganizationTreeVO `json:"children,omitempty"`
}

// PrivilegeUserRoleVO — 用户角色关联响应
type PrivilegeUserRoleVO struct {
	BaseModel
	UserID   string `json:"userId"`
	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
	Username string `json:"username"`
}

// PrivilegeCompanyVO — 公司响应
type PrivilegeCompanyVO struct {
	PrivilegeCompany
}

// PrivilegeDepartmentVO — 部门响应
type PrivilegeDepartmentVO struct {
	PrivilegeDepartment
	CompanyName string `json:"companyName"`
}

// PrivilegeEmployeeVO — 员工绑定响应
type PrivilegeEmployeeVO struct {
	PrivilegeEmployee
	Username string `json:"username"`
	RealName string `json:"realName"`
	DeptName string `json:"deptName"`
}

// PrivilegeDictionaryVO — 字典响应
type PrivilegeDictionaryVO struct {
	PrivilegeDictionary
}

// PrivilegePvalueVO — 权限值响应
type PrivilegePvalueVO struct {
	PrivilegePvalue
}

// PrivilegeLoginLogVO — 登录日志响应
type PrivilegeLoginLogVO struct {
	PrivilegeLoginLog
}

// PrivilegeModuleVO — 模块响应
type PrivilegeModuleVO struct {
	PrivilegeModule
}

// PrivilegeAclVO — 访问控制列表响应
type PrivilegeAclVO struct {
	PrivilegeAcl
	ModuleName string `json:"moduleName"`
}
