package platform

func IsValidAccountStatus(status string) bool {
	return status == "0" || status == "1"
}

func IsAccountDisabled(status string) bool {
	return status == "1"
}

func IsValidGroupStatus(status string) bool {
	return status == "0" || status == "1"
}
