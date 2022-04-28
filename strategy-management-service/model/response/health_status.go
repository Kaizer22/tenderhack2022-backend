package response

type connectionStatus string

const (
	StatusUp      connectionStatus = "UP"
	StatusDown    connectionStatus = "DOWN"
	StatusUnknown connectionStatus = "UNKNOWN"
)

type HealthStatus struct {
	Status connectionStatus `json:"status"`
	Info   string           `json:"info"`
}
