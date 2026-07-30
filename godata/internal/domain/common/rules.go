package common

func IsValidPlatformType(platformType string) bool {
	switch platformType {
	case "dingtalk", "feishu", "wecom":
		return true
	}
	return false
}

func IsValidSyncScope(scope string) bool {
	switch scope {
	case "all", "departments", "users":
		return true
	}
	return false
}
