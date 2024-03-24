package config

import (
	"time"
)

type Config struct {
	ConnectionString string    `json:"connection_string"`
	Agent            *Agent    `json:"agent,omitempty"`
	Jobs             []*Job    `json:"jobs,omitempty"`
	Schemas          []*Schema `json:"schemas,omitempty"`
	Debug            bool      `json:"debug,omitempty"`
}

func (c *Config) GetSchema(name string) *Schema {
	for _, schema := range c.Schemas {
		if schema.Name == name {
			return schema
		}
	}
	return nil
}

type Agent struct {
	Name                         string `json:"name,omitempty"`
	Port                         string `json:"port,omitempty"`
	MetricsExportUrl             string `json:"metrics_export_url,omitempty"`
	MetricsExportIntervalSeconds uint64 `json:"metrics_export_interval_seconds,omitempty"`
	MetricsExportPort            string `json:"metrics_export_port,omitempty"`
}

type Job struct {
	Name        string                 `json:"name,omitempty"`
	Database    string                 `json:"database,omitempty"`
	Collection  string                 `json:"collection,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Schema      string                 `json:"schema,omitempty"`
	Connections uint64                 `json:"connections,omitempty"` // Maximum number of concurrent connections
	Pace        uint64                 `json:"pace,omitempty"` // rps limit / peace - if not set max
	DataSize    uint64                 `json:"data_size,omitempty"` // data size in bytes
	BatchSize   uint64                 `json:"batch_size,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	Operations  uint64                 `json:"operations,omitempty"`
	Timeout     time.Duration          `json:"timeout,omitempty"` // if not set, default
	Filter      map[string]interface{} `json:"filter,omitempty"`
}

type Schema struct {
	Name       string                 `json:"name,omitempty"`
	Database   string                 `json:"database,omitempty"`
	Collection string                 `json:"collection,omitempty"`
	Schema     map[string]interface{} `json:"schema,omitempty"` // todo: introducte new type and parse
	Save       []string               `json:"save,omitempty"`
}
