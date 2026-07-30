package model

// GroupInfo — tbl_platform_group_info
// 组织信息表
type GroupInfo struct {
	PlatformBaseModel
	ID          string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	Name        string `gorm:"column:name;type:varchar(128)" json:"name"`
	SN          string `gorm:"column:sn;type:varchar(64)" json:"sn"`
	Description string `gorm:"column:description;type:varchar(256)" json:"description"`
	Status      int    `gorm:"column:status;default:0" json:"status"`
}

func (GroupInfo) TableName() string { return "tbl_platform_group_info" }

// GroupAgentInfo — tbl_platform_group_agent_info
// 组织-智能体关联表
type GroupAgentInfo struct {
	PlatformBaseModel
	ID      string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	GroupID string `gorm:"column:group_id;type:varchar(64);index" json:"groupId"`
	AgentID int64  `gorm:"column:agent_id" json:"agentId"`
}

func (GroupAgentInfo) TableName() string { return "tbl_platform_group_agent_info" }

// AccountInfo — tbl_platform_account_info
// 前台账号信息表
type AccountInfo struct {
	PlatformBaseModel
	ID           string      `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	Code         string      `gorm:"column:code;type:varchar(64)" json:"code"`
	Username     string      `gorm:"column:username;type:varchar(64);uniqueIndex" json:"username"`
	Password     string      `gorm:"column:password;type:varchar(256)" json:"-"`
	RealName     string      `gorm:"column:real_name;type:varchar(64)" json:"realName"`
	NickName     string      `gorm:"column:nick_name;type:varchar(64)" json:"nickName"`
	Birthday     string      `gorm:"column:birthday;type:varchar(32)" json:"birthday"`
	Email        string      `gorm:"column:email;type:varchar(128)" json:"email"`
	Phone        string      `gorm:"column:phone;type:varchar(32)" json:"phone"`
	AvatarURL    string      `gorm:"column:avatar_url;type:varchar(256)" json:"avatarUrl"`
	Gender       string      `gorm:"column:gender;type:varchar(16)" json:"gender"`
	Status       string      `gorm:"column:status;type:varchar(16);default:0" json:"status"`
	ThirdPartyID string      `gorm:"column:third_party_id;type:varchar(64)" json:"thirdPartyId"`
	EmployeeID   string      `gorm:"column:employee_id;type:varchar(64)" json:"employeeId"`
	DeptID       string      `gorm:"column:dept_id;type:varchar(64)" json:"deptId"`
	DeptName     string      `gorm:"column:dept_name;type:varchar(128)" json:"deptName"`
	Groups       interface{} `gorm:"-" json:"groups,omitempty"`
}

func (AccountInfo) TableName() string { return "tbl_platform_account_info" }

// AccountGroupInfo — tbl_platform_account_group_info
// 账号-组织关联表
type AccountGroupInfo struct {
	PlatformBaseModel
	ID          string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	GroupID     string `gorm:"column:group_id;type:varchar(64)" json:"groupId"`
	AccountID   string `gorm:"column:account_id;type:varchar(64)" json:"accountId"`
	GroupName   string `gorm:"column:group_name;type:varchar(128)" json:"groupName"`
	AccountName string `gorm:"column:account_name;type:varchar(64)" json:"accountName"`
}

func (AccountGroupInfo) TableName() string { return "tbl_platform_account_group_info" }

// AccountTenantInfo — tbl_platform_account_tenant_info
// 账号-租户关联表
type AccountTenantInfo struct {
	PlatformBaseModel
	ID        string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	AccountID string `gorm:"column:account_id;type:varchar(64)" json:"accountId"`
	TenantID  string `gorm:"column:tenant_id;type:varchar(64)" json:"tenantId"`
}

func (AccountTenantInfo) TableName() string { return "tbl_platform_account_tenant_info" }

// TenantInfo — tbl_platform_tenant_info
// 租户信息表
type TenantInfo struct {
	PlatformBaseModel
	ID          string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	SN          string `gorm:"column:sn;type:varchar(64)" json:"sn"`
	Name        string `gorm:"column:name;type:varchar(128)" json:"name"`
	Description string `gorm:"column:description;type:varchar(256)" json:"description"`
}

func (TenantInfo) TableName() string { return "tbl_platform_tenant_info" }
