package modelconfig

func IsActiveModel(provider, modelName, baseURL string) bool {
	return provider != "" && modelName != "" && baseURL != ""
}

func IsValidTemperature(temp float64) bool {
	return temp >= 0 && temp <= 2
}

func IsValidMaxTokens(tokens int) bool {
	return tokens >= 0 && tokens <= 1000000
}
