package cli


// 1. without args, just prints configfile 

// 2. with --set= update config

// 3. generate samle config, for mongodb, postgres itp

// global config -- to rename as config
// type JobType string

// const (
// 	Write          JobType = "write"
// 	BulkWrite      JobType = "bulk_write"
// 	Read           JobType = "read"
// 	Update         JobType = "update"
// 	Sleep          JobType = "sleep"
// 	Paralel        JobType = "parallel"
// 	CreateIndex    JobType = "create_index"
// 	DropCollection JobType = "drop_collection"
// )

// type Config struct {
// 	ConnectionString string             `json:"connection_string"`
// 	Jobs             []*Job             `json:"jobs"`
// 	Schemas          []*Schema          `json:"schemas"`
// 	ReportingFormats []*ReportingFormat `json:"reporting_formats"`
// 	Debug            bool               `json:"debug"`
// 	DebugFile        string             `json:"debug_file"`
// }

// type Job struct {
// 	Parent          *Config
// 	Name            string
// 	Database        string
// 	Collection      string
// 	Type            string
// 	Schema          string
// 	ReportingFormat string
// 	Connections     uint64 // Maximum number of concurrent connections
// 	Pace            uint64 // rps limit / peace - if not set max
// 	DataSize        uint64 // data size in bytes
// 	BatchSize       uint64
// 	Duration        time.Duration
// 	Operations      uint64
// 	Timeout         time.Duration // if not set, default
// 	Filter          map[string]interface{}
// }

// type Schema struct {
// 	Name       string                 `json:"name"`
// 	Database   string                 `json:"database"`
// 	Collection string                 `json:"collection"`
// 	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
// 	Save       []string               `json:"save"`
// }
