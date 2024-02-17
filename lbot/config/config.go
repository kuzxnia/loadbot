package config

import (
	"time"
)

type Config struct {
	ConnectionString  string    `json:"connection_string"`
	AgentName         string    `json:"agent_name"`
	AgentPort         string    `json:"agent_port"`
	MetricsExportUrl  string    `json:"metrics_export_url"`
	MetricsExportPort string    `json:"metrics_export_port"`
	Jobs              []*Job    `json:"jobs"`
	Schemas           []*Schema `json:"schemas"`
	Debug             bool      `json:"debug"`
}

type Job struct {
	Parent      *Config
	Name        string
	Database    string
	Collection  string
	Type        string
	Schema      string
	Connections uint64 // Maximum number of concurrent connections
	Pace        uint64 // rps limit / peace - if not set max
	DataSize    uint64 // data size in bytes
	BatchSize   uint64
	Duration    time.Duration
	Operations  uint64
	Timeout     time.Duration // if not set, default
	Filter      map[string]interface{}
}

type Schema struct {
	Name       string                 `json:"name"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
	Save       []string               `json:"save"`
}

func (j *Job) GetSchema() *Schema {
	for _, schema := range j.Parent.Schemas {
		if schema.Name == j.Schema {
			return schema
		}
	}
	return nil
}
