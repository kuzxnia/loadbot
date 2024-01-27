package config

import (
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// driver Settings

// global config -- to rename as config
type JobType string

const (
	Write          JobType = "write"
	BulkWrite      JobType = "bulk_write"
	Read           JobType = "read"
	Update         JobType = "update"
	Sleep          JobType = "sleep"
	Paralel        JobType = "parallel"
	CreateIndex    JobType = "create_index"
	DropCollection JobType = "drop_collection"
)

type Config struct {
	ConnectionString string             `json:"connection_string"`
	Jobs             []*Job             `json:"jobs"`
	Schemas          []*Schema          `json:"schemas"`
	ReportingFormats []*ReportingFormat `json:"reporting_formats"`
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
	Filter          map[string]interface{}
	Indexes         []*Index
}

type Index struct {
	Keys    interface{}
	Options IndexOptions
}

type IndexOptions struct {
	Background              *bool              `json:"background,omitempty"`
	ExpireAfterSeconds      *int32             `json:"expire_after_seconds,omitempty"`
	Name                    *string            `json:"name,omitempty"`
	Sparse                  *bool              `json:"sparse,omitempty"`
	StorageEngine           interface{}        `json:"storage_engine,omitempty"`
	Unique                  *bool              `json:"unique,omitempty"`
	Version                 *int32             `json:"version,omitempty"`
	DefaultLanguage         *string            `json:"default_language,omitempty"`
	LanguageOverride        *string            `json:"language_override,omitempty"`
	TextVersion             *int32             `json:"text_version,omitempty"`
	Weights                 interface{}        `json:"weights,omitempty"`
	SphereVersion           *int32             `json:"sphere_version,omitempty"`
	Bits                    *int32             `json:"bits,omitempty"`
	Max                     *float64           `json:"max,omitempty"`
	Min                     *float64           `json:"min,omitempty"`
	BucketSize              *int32             `json:"bucket_size,omitempty"`
	PartialFilterExpression interface{}        `json:"partial_filter_expression,omitempty"`
	Collation               *options.Collation `json:"collation,omitempty"`
	WildcardProjection      interface{}        `json:"wildcard_projection,omitempty"`
	Hidden                  *bool              `json:"hidden,omitempty"`
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

type Schema struct {
	Name       string                 `json:"name"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
	Save       []string               `json:"save"`
}

var DefaultReportFormats = []*ReportingFormat{
	{
		Name:     "default",
		Interval: 15 * time.Second,
		Template: `{{.Now}} {{if .JobName -}}Job: "{{.JobName}}" {{else -}}Job type: "{{.JobType}}"{{end}}
Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}}, s:{{.SuccessReqs}}/err:{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}%
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
	{
		Name:     "simple",
		Interval: 15 * time.Second,
		Template: "{{.Now}} Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}} s:{{.SuccessReqs}}/err:{{.ErrorReqs}}\n\n",
	},
	{
		Name:     "write",
		Interval: 15 * time.Second,
		Template: `{{.Now}} Reqs: {{.TotalReqs}}, RPS {{f2 .Rps}}, s:{{.SuccessReqs}}/err:{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}%
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
	{
		Name:     "bulk_write",
		Interval: 15 * time.Second,
		Template: `{{.Now}} Reqs: {{.TotalReqs}}, OPS: {{.TotalOps}}, RPS {{f2 .Rps}}, OPS {{f2 .Ops}}, s:{{.SuccessReqs}}/err{{.ErrorReqs}}/tout:{{.TimeoutErr}}/errRate:{{f1 .ErrorRate}}
AVG: {{msf3 .Avg}}ms P50: {{msf3 .P50}}ms, P90: {{msf3 .P90}}ms P99: {{msf3 .P99}}ms

`,
	},
}
var DefaultReportFormat = DefaultReportFormats[0]

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
