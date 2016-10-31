package common

const (
	// EnvironmentVariablePrefix - the prefix for all nexus related configuration environment variables.
	EnvironmentVariablePrefix = "NEXUS_"

	// RollingRestart - the environment variable to signify a rolling restart
	RollingRestart = EnvironmentVariablePrefix + "ROLLING_RESTART_TRIGGER"
)
