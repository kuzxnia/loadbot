package config

import (
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
	// todo
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
		// c.validateWriteRatio,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *Config) validateWriteRatio() error {
	// if c.WriteRatio < 0.0 || c.WriteRatio > 1.0 {
	// 	return fmt.Errorf("Write ratio must be in range 0..1")
	// }
	return nil
}

// add more validators
