package config

import (
	"time"

	"github.com/samber/lo"
)

type Config struct {
	ConnectionString string             `json:"connection_string"`
	Jobs             []*Job             `json:"jobs"`
	Schemas          []*Schema          `json:"schemas"`
	ReportingFormats []*ReportingFormat `json:"reporting_formats"`
	Debug            bool               `json:"debug"`
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
	Filter          map[string]interface{}
}

type Schema struct {
	Name       string                 `json:"name"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
	Save       []string               `json:"save"`
}

type ReportingFormat struct {
	Name     string
	Interval time.Duration
	Template string
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
	reportingFormat := lo.If(j.ReportingFormat != "", j.ReportingFormat).Else(j.Type)
	reportingFormats := append(j.Parent.ReportingFormats, DefaultReportFormats...)

	return lo.FindOrElse(reportingFormats, DefaultReportFormat, func(rf *ReportingFormat) bool { return rf.Name == reportingFormat })
}
