package datasource

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusError    = "error"
)

var SupportedTypes = []string{
	"MySQL", "PostgreSQL", "Oracle", "SQLServer", "ClickHouse", "MongoDB", "Redis",
}

func DefaultDatasourceStatus() string {
	return StatusActive
}

func IsSupportedType(dsType string) bool {
	for _, t := range SupportedTypes {
		if t == dsType {
			return true
		}
	}
	return false
}

func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}
