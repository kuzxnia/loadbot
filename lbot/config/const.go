package config

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
