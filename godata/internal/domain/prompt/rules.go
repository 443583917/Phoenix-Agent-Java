package prompt

var SupportedTypes = []string{"system", "user", "assistant", "tool", "reasoning"}

func IsValidPromptType(promptType string) bool {
	for _, t := range SupportedTypes {
		if t == promptType {
			return true
		}
	}
	return false
}

func DefaultPriority() int {
	return 0
}

func DefaultDisplayOrder() int {
	return 0
}
