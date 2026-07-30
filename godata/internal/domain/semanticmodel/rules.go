package semanticmodel

const (
	StatusDisabled = 0
	StatusEnabled  = 1
)

func DefaultSemanticModelStatus() int {
	return StatusEnabled
}

func IsEnabled(status int) bool {
	return status == StatusEnabled
}
