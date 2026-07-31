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
	UserID     string `gorm:"column:user_id;type:varchar(64);index" json:"userId"`
	UserNo     string `gorm:"column:user_no;type:varchar(32)" json:"userNo"`
	RoleID     string `gorm:"column:role_id;type:varchar(64);index" json:"roleId"`
	ValidMonth int    `gorm:"column:valid_month" json:"validMonth"`
}

func (PrivilegeUserRole) TableName() string { return "tbl_privilege_user_role" }

// PrivilegeModule — tbl_privilege_module (菜单/模块)
type PrivilegeModule struct {
	BaseModel
	Name       string  `gorm:"column:name;type:varchar(64)" json:"name"`
	Code       string  `gorm:"column:sn;type:varchar(64)" json:"code"`          // DB: sn
	PID        *string `gorm:"column:pid;type:varchar(64)" json:"pid"`
	URL        string  `gorm:"column:url;type:varchar(256)" json:"url"`
	Icon       string  `gorm:"column:image;type:varchar(128)" json:"icon"`      // DB: image
	State      string  `gorm:"column:state;type:varchar(100)" json:"state"`     // DB: state
	Component  string  `gorm:"column:component;type:varchar(120)" json:"component"`
	SystemID   string  `gorm:"column:system_id;type:varchar(64)" json:"systemId"`
	Status     int     `gorm:"column:status" json:"status"`                     // DB: status(int4)
	OrderNo    int     `gorm:"column:order_no;default:0" json:"orderNo"`        // DB: order_no
	IsShow     int     `gorm:"column:is_show;default:1" json:"isShow"`          // DB: is_show
	CategoryID int     `gorm:"column:category_id" json:"categoryId"`            // DB: category_id
	Type       string  `gorm:"column:type;type:varchar(255)" json:"type"`       // DB: type(varchar)
}

func (PrivilegeModule) TableName() string { return "tbl_privilege_module" }

// PrivilegeAcl — tbl_privilege_acl (访问控制列表)
type PrivilegeAcl struct {
	BaseModel
	ReleaseID  string `gorm:"column:release_id;type:varchar(64)" json:"releaseId"`
	ReleaseSN  string `gorm:"column:release_sn;type:varchar(10)" json:"releaseSn"`
	SystemSN   string `gorm:"column:system_sn;type:varchar(40)" json:"systemSn"`
	ModuleID   string `gorm:"column:module_id;type:varchar(64)" json:"moduleId"`
	ModuleSN   string `gorm:"column:module_sn;type:varchar(40)" json:"moduleSn"`
	AclState   string `gorm:"column:acl_state;type:varchar(100)" json:"aclState"`
}

func (PrivilegeAcl) TableName() string { return "tbl_privilege_acl" }

// PrivilegeDepartment — tbl_privilege_department
type PrivilegeDepartment struct {
	BaseModel
	CompanyID      string  `gorm:"column:company_id;type:varchar(64)" json:"companyId"`
	Name           string  `gorm:"column:name;type:varchar(100)" json:"name"`
	Code           string  `gorm:"column:code;type:varchar(20)" json:"code"`
	Note           string  `gorm:"column:note;type:varchar(80)" json:"note"`
	PID            *string `gorm:"column:pid;type:varchar(64)" json:"pid"`
	OrderNo        int     `gorm:"column:order_no" json:"orderNo"`             // DB: order_no
	Leader         int     `gorm:"column:leader;default:0" json:"leader"`
	DepartmentType int     `gorm:"column:department_type" json:"departmentType"`
	Status         int     `gorm:"column:status;default:0" json:"status"`       // DB: status(int2)
	Nature         int     `gorm:"column:nature;default:0" json:"nature"`       // DB: nature(int2)
	ThirdID        string  `gorm:"column:third_id;type:varchar(255)" json:"thirdId"`
}

func (PrivilegeDepartment) TableName() string { return "tbl_privilege_department" }

// PrivilegeCompany — tbl_privilege_company
type PrivilegeCompany struct {
	BaseModel
	PID           string `gorm:"column:pid;type:varchar(64)" json:"pid"`
	CName         string `gorm:"column:cname;type:varchar(128)" json:"cname"`
	EName         string `gorm:"column:ename;type:varchar(128)" json:"ename"`
	IDMCompanyID  string `gorm:"column:idm_company_id;type:varchar(64)" json:"idmCompanyId"`
	ShortName     string `gorm:"column:short_name;type:varchar(120)" json:"shortName"`
	Code          string `gorm:"column:code;type:varchar(64)" json:"code"`
	ThirdID       string `gorm:"column:third_id;type:varchar(255)" json:"thirdId"`
	Descr         string `gorm:"column:descr;type:varchar(200)" json:"descr"`  // DB: descr
	Status        int    `gorm:"column:status;default:1" json:"status"`         // DB: status(int4)
	Sort          int    `gorm:"column:sort;default:0" json:"sort"`             // DB: sort(int2)
}

func (PrivilegeCompany) TableName() string { return "tbl_privilege_company" }

// PrivilegeEmployee — tbl_privilege_employee
type PrivilegeEmployee struct {
	BaseModel
	EmpCode        string     `gorm:"column:emp_code;type:varchar(255)" json:"empCode"`
	EmpName        string     `gorm:"column:emp_name;type:varchar(255)" json:"empName"`
	PositionCode   string     `gorm:"column:position_code;type:varchar(255)" json:"positionCode"`
	JobGradeCode   string     `gorm:"column:job_grade_code;type:varchar(255)" json:"jobGradeCode"`
	LeaderUserID   string     `gorm:"column:leader_user_id;type:varchar(255)" json:"leaderUserId"`
	LeaderUserName string     `gorm:"column:leader_user_name;type:varchar(255)" json:"leaderUserName"`
	CompanyID      string     `gorm:"column:company_id;type:varchar(255)" json:"companyId"`
	CompanyName    string     `gorm:"column:company_name;type:varchar(255)" json:"companyName"`
	DeptID         string     `gorm:"column:dept_id;type:varchar(255)" json:"deptId"`
	DeptName       string     `gorm:"column:dept_name;type:varchar(255)" json:"deptName"`
	Sex            int        `gorm:"column:sex" json:"sex"`
	Status         int        `gorm:"column:status;default:1" json:"status"`
	EnableFlag     int        `gorm:"column:enable_flag;default:3" json:"enableFlag"`
	ServiceDate    *time.Time `gorm:"column:service_date" json:"serviceDate,omitempty"`
	LeaveDate      *time.Time `gorm:"column:leave_date" json:"leaveDate,omitempty"`
	ThirdUnionID   string     `gorm:"column:third_union_id;type:varchar(255)" json:"thirdUnionId"`
	ThirdOpenID    string     `gorm:"column:third_open_id;type:varchar(255)" json:"thirdOpenId"`
	ThirdUserID    string     `gorm:"column:third_user_id;type:varchar(255)" json:"thirdUserId"`
	AvatarURL      string     `gorm:"column:avatar_url;type:varchar(255)" json:"avatarUrl"`
	Mobile         string     `gorm:"column:mobile;type:varchar(255)" json:"mobile"`
	Email          string     `gorm:"column:email;type:varchar(255)" json:"email"`
	IsDeptLeader   string     `gorm:"column:is_dept_leader;type:varchar(255)" json:"isDeptLeader"`
	Paths          string     `gorm:"column:paths;type:varchar(255)" json:"paths"`
}

func (PrivilegeEmployee) TableName() string { return "tbl_privilege_employee" }

// PrivilegeDictionary — tbl_privilege_dictionary (uses creator_by not create_by)
type PrivilegeDictionary struct {
	ID         int64     `gorm:"primaryKey;column:id;type:bigint" json:"id"`
	Code       string    `gorm:"column:code;type:varchar(64)" json:"code"`
	Name       string    `gorm:"column:name;type:varchar(128)" json:"name"`
	PCode      *string   `gorm:"column:pcode;type:varchar(64)" json:"pcode"`
	SystemSN   string    `gorm:"column:system_sn;type:varchar(32)" json:"systemSn"`
	SN         string    `gorm:"column:sn;type:varchar(32)" json:"sn"`
	OrderNo    int       `gorm:"column:order_no" json:"orderNo"`               // DB: order_no
	DelFlag    int       `gorm:"column:del_flag;default:0" json:"delFlag"`
	CreatorBy  string    `gorm:"column:creator_by;type:varchar(32)" json:"creatorBy"` // DB: creator_by (not create_by)
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	UpdateBy   *string   `gorm:"column:update_by;type:varchar(32)" json:"updateBy,omitempty"`
}

func (PrivilegeDictionary) TableName() string { return "tbl_privilege_dictionary" }

// PrivilegePvalue — tbl_privilege_pvalue (权限值)
type PrivilegePvalue struct {
	BaseModel
	Position int    `gorm:"column:position" json:"position"`                     // DB: position
	Name     string `gorm:"column:name;type:varchar(32)" json:"name"`
	OrderNo  int    `gorm:"column:order_no" json:"orderNo"`                     // DB: order_no
	Remark   string `gorm:"column:remark;type:varchar(200)" json:"remark"`
	SystemID string `gorm:"column:system_id;type:varchar(64)" json:"systemId"`  // DB: system_id(int8)
}

func (PrivilegePvalue) TableName() string { return "tbl_privilege_pvalue" }

// PrivilegeLoginLog — tbl_privilege_login_log
type PrivilegeLoginLog struct {
	ID               int64     `gorm:"primaryKey;column:id;type:bigint" json:"id"`
	OperationID      int64     `gorm:"column:operation_id" json:"operationId"`
	IP               string    `gorm:"column:ip;type:varchar(45)" json:"ip"`
	OperationUsername string   `gorm:"column:operation_username;type:varchar(32)" json:"operationUsername"`
	OperationPerson  string    `gorm:"column:operation_person;type:varchar(32)" json:"operationPerson"`
	OperationContent string    `gorm:"column:operation_content;type:varchar(128)" json:"operationContent"`
	OperationTime    time.Time `gorm:"column:operation_time" json:"operationTime"`
	DelFlag          int       `gorm:"column:del_flag" json:"delFlag"`
	CreateBy         *string   `gorm:"column:create_by;type:varchar(32)" json:"createBy,omitempty"`
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateBy         *string   `gorm:"column:update_by;type:varchar(32)" json:"updateBy,omitempty"`
	UpdateTime       time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (PrivilegeLoginLog) TableName() string { return "tbl_privilege_login_log" }

// PrivilegeGroup — tbl_privilege_group (uses creator_by not create_by)
type PrivilegeGroup struct {
	ID           int64     `gorm:"primaryKey;column:id;type:bigint;autoIncrement" json:"id"`
	SuperID      int64     `gorm:"column:super_id" json:"superId"`
	Name         string    `gorm:"column:name;type:varchar(32)" json:"name"`
	Type         int       `gorm:"column:type" json:"type"`                       // DB: type(int2)
	State        int       `gorm:"column:state" json:"state"`                     // DB: state(int2)
	Introduction string    `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	CreatorBy    string    `gorm:"column:creator_by;type:varchar(32)" json:"creatorBy"` // DB: creator_by
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	UpdateBy     *string   `gorm:"column:update_by;type:varchar(32)" json:"updateBy,omitempty"`
	DelFlag      int       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (PrivilegeGroup) TableName() string { return "tbl_privilege_group" }
