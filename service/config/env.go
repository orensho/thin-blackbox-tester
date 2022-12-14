package config

// Filled using ldflags at compile time
var (
	BuildTime   = "unknown"
	Branch      = "unknown"
	Commit      = "unknown"
	BuildNumber = "dev"
)
