package config

import (
	"time"
)

// global config -- to rename as config
type JobType string

const (
	Write          JobType = "write"
	BulkWrite      JobType = "bulk_write"
	Read           JobType = "read"
	Update         JobType = "update"
	Sleep          JobType = "sleep"
	Paralel        JobType = "parallel"
	BuildIndex     JobType = "parallel"
	DropCollection JobType = "drop_collection"
)

type Config struct {
	ConnectionString string             `json:"connection_string"`
	Jobs             []*Job             `json:"jobs"`
	Schemas          []*Schema          `json:"schemas"`
	ReportingFormats []*ReportingFormat `json:"reports"`
	Debug            bool               `json:"debug"`
	DebugFile        string             `json:"debug_file"`
}

type Job struct {
	Parent          *Config
	Name            string
	Database        string
	Collection      string
	Type            string
	Schema          string
	ReportingFormat string
	Connections     uint64 // Maximum number of concurrent connections
	Pace            uint64 // rps limit / peace - if not set max
	DataSize        uint64 // data size in bytes
	BatchSize       uint64
	Duration        time.Duration
	Operations      uint64
	Timeout         time.Duration // if not set, default
	// Params ex. for read / update
	//     * filter: { "_id": "#_id"}
}

func (j *Job) GetSchema() *Schema {
	for _, schema := range j.Parent.Schemas {
		if schema.Name == j.Schema {
			return schema
		}
	}
	return nil
}

func (j *Job) GetReport() *ReportingFormat {
	for _, report := range j.Parent.ReportingFormats {
		if report.Name == j.ReportingFormat {
			return report
		}
	}
	return nil
}

type Schema struct {
	Name       string `json:"name"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
	// todo: introducte new type and parse
	Schema map[string]interface{} `json:"schema"`
}

type ReportingFormat struct {
	Name     string
	Interval time.Duration
	Template string
}

func NewConfigFromArgs() *Config {
	return nil
}

func NewConfigFromJson([]byte) *Config {
	return nil
}
