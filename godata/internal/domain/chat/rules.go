package chat

const (
	SessionStatusActive  = "active"
	SessionStatusClosed  = "closed"
	SessionStatusArchived = "archived"
)

func DefaultSessionStatus() string {
	return SessionStatusActive
}

func IsValidSessionStatus(status string) bool {
	switch status {
	case SessionStatusActive, SessionStatusClosed, SessionStatusArchived:
		return true
	}
	return false
}

func IsValidRole(role string) bool {
	switch role {
	case "user", "assistant", "system", "tool":
		return true
	}
	return false
}
