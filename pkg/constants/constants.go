package constants

const (
	ConfigFilepath = "pagerduty-cli/config.json"

	APIKeyURL   = "https://support.pagerduty.com/docs/generating-api-keys#section-generating-a-general-access-rest-api-key"
	APIKeyRegex = "^[a-z|A-Z0-9+_-]{20}$"

	// Sample API key for testing
	SampleKey = "y_NbAkKc66ryYTWUXYEu"

	AlertsLimit = 10

	// PagerDuty IDs
	TeamID     = "PASPK4G"
	SilentTest = "P8QS6CC"

	// PagerDuty Incident Statuses
	StatusTriggered    = "triggered"
	StatusAcknowledged = "acknowledged"
	StatusResolved     = "resolved"
)
