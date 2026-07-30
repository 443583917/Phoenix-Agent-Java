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
	UserID    string    `gorm:"column:user_id;type:varchar(64)" json:"userId"`
	Username  string    `gorm:"column:username;type:varchar(64)" json:"username"`
	LoginIP   string    `gorm:"column:login_ip;type:varchar(64)" json:"loginIp"`
	LoginTime time.Time `gorm:"column:login_time" json:"loginTime"`
}

func (PrivilegeLoginLog) TableName() string { return "tbl_privilege_login_log" }
