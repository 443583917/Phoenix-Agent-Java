package knowledge

const (
	EmbeddingStatusPending  = "pending"
	EmbeddingStatusComplete = "complete"
	EmbeddingStatusFailed   = "failed"
)

func DefaultEmbeddingStatus() string {
	return EmbeddingStatusPending
}

func DefaultRecallEnabled() int {
	return 1
}

func IsEmbeddingTerminal(status string) bool {
	return status == EmbeddingStatusComplete || status == EmbeddingStatusFailed
}

func CanRetryEmbedding(status string) bool {
	return status == EmbeddingStatusFailed
}
