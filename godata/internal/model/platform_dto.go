package model

// AccountLoginDTO — POST /auth/login
// 账号登录请求
type AccountLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ThirdPartyLoginDTO — POST /auth/thirdLogin
// 第三方登录请求
type ThirdPartyLoginDTO struct {
	ThirdPartyID string `json:"thirdPartyId" binding:"required"`
	Type         string `json:"type" binding:"required"`
}

// UpdatePwdDTO — PUT /auth/updatePassword
// 修改密码请求
type UpdatePwdDTO struct {
	UserID      string `json:"userId"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
