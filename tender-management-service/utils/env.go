package utils

const (

	// Common -----------------
	ListenAddressEnvKey           = "PORT"
	ServiceNameEnvKey             = "SERVICE_NAME"
	TenderManagementBaseUrlEnvKey = "TENDER_MANAGEMENT_BASE_URL" // to generate proper Swagger docs

	// Security ---------------
	SecurityTokenTTLKey = "TOKEN_TTL"
	AccessSecret        = "ACCESS_SECRET"
	RefreshSecret       = "REFRESH_SECRET"

	// Debug -------------------
	AppVersionEnvKey = "APP_VERSION"
	LogLevelEnvKey   = "LOG_LEVEL"

	// Work logic --------------
	RunSessionCron               = "RUN_SESSION_CRON"
	SessionUpdateFrequencyEnvKey = "SESSION_UPDATE_FREQUENCY"
)
