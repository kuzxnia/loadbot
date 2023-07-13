package config

import (
	"errors"
	"time"
)

// global config -- to rename as config
type JobType string

const (
	Write      JobType = "write"
	BulkWrite  JobType = "bulk_write"
	Read       JobType = "read"
	Update     JobType = "update"
	Sleep      JobType = "sleep"
	Paralel    JobType = "parallel"
	BuildIndex JobType = "parallel"
)

type Config struct {
	ConnectionString string    `json:"connection_string"`
	Jobs             []*Job    `json:"jobs"`
	Schemas          []*Schema `json:"schemas"`
	Debug            bool      `json:"debug"`
	DebugFile        string    `json:"debug_file"`
}

type Job struct {
	Parent      *Config
	Name        string
	Type        string
	Template    string
	Connections uint64 // Maximum number of concurrent connections
	Pace        uint64 // rps limit / peace - if not set max
	DataSize    uint64 // data size in bytes
	BatchSize   uint64
	Duration    time.Duration
	Operations  uint64
	Timeout     time.Duration // if not set, default
	// Params ex. for read / update
	//     * filter: { "_id": "#_id"}
}

func (j *Job) GetTemplateSchema() *Schema {
	for _, schema := range j.Parent.Schemas {
		if schema.Name == j.Template {
			return schema
		}
	}
	return nil
}

type Schema struct {
	Name       string
	Database   string
	Collection string
	Schema     map[string]string
	// template - nested dict
}

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateAllJobTemplatesAreProvided,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *Config) validateAllJobTemplatesAreProvided() error {
	isSchemaName := func(schema *Schema, comparator string) bool { return schema.Name == comparator }

	for _, job := range c.Jobs {
		if !Contains[*Schema, string](c.Schemas, job.Template, isSchemaName) {
			return errors.New("Job: " + job.Name + " have invalid template \"" + job.Template + "\"")
		}
	}
	return nil
}

func Contains[T comparable, X comparable](array []T, comparator X, predicate func(T, X) bool) bool {
	for _, elem := range array {
		if predicate(elem, comparator) {
			return true
		}
	}
	return false
}
