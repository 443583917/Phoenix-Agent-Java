package agent

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusOffline   = "offline"
)

func DefaultAgentStatus() string {
	return StatusDraft
}

func IsValidStatusTransition(from, to string) bool {
	switch {
	case from == StatusDraft && to == StatusPublished:
		return true
	case from == StatusPublished && to == StatusOffline:
		return true
	case from == StatusOffline && to == StatusPublished:
		return true
	default:
		return false
	}
}

func CanPublish(status string) bool {
	return status == StatusDraft || status == StatusOffline
}

func CanOffline(status string) bool {
	return status == StatusPublished
}
