package utils

const (

	// Common -----------------
	ListenAddressEnvKey             = "PORT"
	ServiceNameEnvKey               = "SERVICE_NAME"
	StrategyManagementBaseUrlEnvKey = "STRATEGY_MANAGEMENT_BASE_URL" // to generate proper Swagger docs

	// Security ---------------
	SecurityTokenTTLKey          = "TOKEN_TTL"
	AccessSecret                 = "ACCESS_SECRET"
	RefreshSecret                = "REFRESH_SECRET"
	SessionUpdateFrequencyEnvKey = "SESSION_UPDATE_FREQUENCY"

	// Debug -------------------
	AppVersionEnvKey = "APP_VERSION"
	LogLevelEnvKey   = "LOG_LEVEL"
)
