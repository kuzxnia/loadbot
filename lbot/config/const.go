package config

import "time"

type JobType string

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

const (
	Write          JobType = "write"
	BulkWrite      JobType = "bulk_write"
	Read           JobType = "read"
	Update         JobType = "update"
	Sleep          JobType = "sleep"
	DropCollection JobType = "drop_collection"
)

const (
	AgentsHeartbeatInterval   = time.Second * 2
	AgentsHeartbeatExpiration = -time.Second * 4
)

const (
	DB                    = "admin"
	LockCollection        = "lbotLock"
	CommandCollection     = "lbotCmd"
	ConfigCollection      = "lbotConfig"
	AgentStatusCollection = "lbotAgent"

	// new commands
)
