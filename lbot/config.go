package lbot

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/tailscale/hujson"
	"golang.org/x/net/context"
)

func NewConfig(request *ConfigRequest) *config.Config {
	cfg := &config.Config{
		ConnectionString: request.ConnectionString,
		Jobs:             make([]*config.Job, len(request.Jobs)),
		Schemas:          make([]*config.Schema, len(request.Schemas)),
		ReportingFormats: make([]*config.ReportingFormat, len(request.ReportingFormats)),
		Debug:            request.Debug,
	}
	for i, job := range request.Jobs {
		cfg.Jobs[i] = &config.Job{
			Name:            job.Name,
			Database:        job.Database,
			Collection:      job.Collection,
			Type:            job.Type,
			Schema:          job.Schema,
			ReportingFormat: job.ReportingFormat,
			Connections:     job.Connections,
			Pace:            job.Pace,
			DataSize:        job.DataSize,
			BatchSize:       job.BatchSize,
			Duration:        job.Duration,
			Operations:      job.Operations,
			Timeout:         job.Timeout,
			Filter:          job.Filter,
		}
	}
	for i, schema := range request.Schemas {
		cfg.Schemas[i] = &config.Schema{
			Name:       schema.Name,
			Database:   schema.Database,
			Collection: schema.Collection,
			Schema:     schema.Schema,
			Save:       schema.Save,
		}
	}
	for i, rf := range request.ReportingFormats {
		cfg.ReportingFormats[i] = &config.ReportingFormat{
			Name:     rf.Name,
			Interval: rf.Interval,
			Template: rf.Template,
		}
	}

	return cfg
}

// todo: should be pointers
type ConfigRequest struct {
	ConnectionString string                    `json:"connection_string"`
	Jobs             []*JobRequest             `json:"jobs"`
	Schemas          []*SchemaRequest          `json:"schemas"`
	ReportingFormats []*ReportingFormatRequest `json:"reporting_formats"`
	Debug            bool                      `json:"debug"`
}

type JobRequest struct {
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

type SchemaRequest struct {
	Name       string                 `json:"name"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
	Save       []string               `json:"save"`
}

type ReportingFormatRequest struct {
	Name     string
	Interval time.Duration
	Template string
}

type SetConfigProcess struct {
	ctx  context.Context
	lbot *Lbot
}

func NewSetConfigProcess(ctx context.Context, lbot *Lbot) *SetConfigProcess {
	return &SetConfigProcess{ctx: ctx, lbot: lbot}
}

func (c *SetConfigProcess) Run(request *ConfigRequest, reply *int) error {
	// if watch arg - run watch

	// 	driver.Torment(config)
	c.lbot.SetConfig(nil)
	c.lbot.Run(c.ctx)

	// before configing process it will varify health of cluster, if pods
	return nil
}

func ParseConfigFile(configFile string) (*ConfigRequest, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	content, err = standardizeJSON(content)
	if err != nil {
		return nil, err
	}

	var cfg ConfigRequest
	err = json.Unmarshal(content, &cfg)

	if err != nil {
		return nil, errors.New("Error during Unmarshal(): " + err.Error())
	}

	return &cfg, err
}

func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func (c *JobRequest) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name            string                 `json:"name"`
		Type            string                 `json:"type"`
		Database        string                 `json:"database"`
		Collection      string                 `json:"collection"`
		Schema          string                 `json:"template"`
		ReportingFormat string                 `json:"format"`
		Connections     uint64                 `json:"connections"`
		Pace            uint64                 `json:"pace"`
		DataSize        uint64                 `json:"data_size"`
		BatchSize       uint64                 `json:"batch_size"`
		Duration        string                 `json:"duration"`
		Operations      uint64                 `json:"operations"`
		Timeout         string                 `json:"timeout"` // if not set, default
		Filter          map[string]interface{} `json:"filter"`
	}
	// default values
	tmp.Connections = 1

	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	c.Name = tmp.Name
	c.Database = tmp.Database
	c.Collection = tmp.Collection
	c.Type = tmp.Type
	c.Schema = tmp.Schema
	c.ReportingFormat = tmp.ReportingFormat
	c.Connections = tmp.Connections
	c.Pace = tmp.Pace
	c.DataSize = tmp.DataSize
	c.BatchSize = tmp.BatchSize

	if tmp.Duration != "" {
		if c.Duration, err = time.ParseDuration(tmp.Duration); err != nil {
			return err
		}
	}

	c.Operations = tmp.Operations

	if tmp.Timeout != "" {
		if c.Timeout, err = time.ParseDuration(tmp.Timeout); err != nil {
			return err
		}
	}
	c.Filter = tmp.Filter

	return
}

func (c *ReportingFormatRequest) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name     string `json:"name"`
		Interval string `json:"interval"`
		Template string `json:"template"`
	}

	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	c.Name = tmp.Name
	c.Template = tmp.Template

	if tmp.Interval != "" {
		if c.Interval, err = time.ParseDuration(tmp.Interval); err != nil {
			return err
		}
	}

	return
}

func (c *ConfigRequest) Validate() error {
	validators := []func() error{
		c.validateJobs,
		// c.validateSchemas,
		c.validateReportingFormats,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *ConfigRequest) validateJobs() error {
	for _, job := range c.Jobs {
		if error := job.Validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *ConfigRequest) validateReportingFormats() error {
	for _, reportingFormat := range c.ReportingFormats {
		if error := reportingFormat.Validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *JobRequest) Validate() error {
	validators := []func() error{
		job.validateSchema,
		job.validateReportFormat,
		job.validateDatabase,
		job.validateCollection,
		job.validateType,
		job.validateDuration,
		job.validatePace,
		job.validateConnections,
		job.validateBatchSize,
		job.validateOperations,
		job.validateDataSize,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (job *JobRequest) validateSchema() error {
	if string(config.Sleep) == job.Type || job.Schema == "" {
		return nil
	}

	// todo: to fix
	// if !Contains(job.Parent.Schemas, func(s *Schema) bool { return s.Name == job.Schema }) {
	// 	return errors.New("JobValidationError: job \"" + job.Name + "\" have invalid template \"" + job.Schema + "\"")
	// }
	return nil
}

func (job *JobRequest) validateReportFormat() error {
	if job.ReportingFormat == "" {
		return nil
	}

	// todo: to fix
	// reportingFormats := append(job.Parent.ReportingFormats, DefaultReportFormats...)
	// if !Contains(reportingFormats, func(s *ReportingFormat) bool { return s.Name == job.ReportingFormat }) {
	// 	return errors.New("JobValidationError: job \"" + job.Name + "\" have invalid report_format \"" + job.ReportingFormat + "\"")
	// }
	return nil
}

func (job *JobRequest) validateType() (err error) {
	switch job.Type {
	case string(config.Write):
	case string(config.BulkWrite):
	case string(config.Read):
	case string(config.Update):
	case string(config.DropCollection):
	case string(config.Sleep):
	default:
		err = errors.New("Job type: " + job.Type + " ")
	}
	return
}

func (job *JobRequest) validateDatabase() (err error) {
	if job.Schema != "" || job.Type == string(config.Sleep) {
		return
	}
	if job.Database == "" {
		err = errors.New("JobValidationError: field 'database' is required if 'template' or 'type' is not set")
	}
	return
}

func (job *JobRequest) validateCollection() (err error) {
	if job.Schema != "" || job.Type == string(config.Sleep) {
		return
	}
	if job.Collection == "" {
		err = errors.New("JobValidationError: field 'collection' is required if 'template' or 'type' is not set")
	}
	return
}

func (job *JobRequest) validateConnections() (err error) {
	if job.Connections == 0 {
		err = errors.New("JobValidationError: field 'connections' must be greater than 0")
	}
	if job.Type == string(config.Sleep) {
		if job.Connections != 1 {
			err = errors.New("JobValidationError: field 'connections' max number concurrent connections for job type 'sleep' is 1")
		}
	}
	return
}

func (job *JobRequest) validateDuration() (err error) {
	if job.Type == string(config.Sleep) {
		if job.Duration <= 0 {
			err = errors.New("JobValidationError: field 'duration' must be greater than 0 for job with 'sleep' type ")
		}
	}
	return
}

func (job *JobRequest) validatePace() (err error) {
	if job.Type == string(config.Sleep) {
		if job.Pace != 0 {
			err = errors.New("JobValidationError: field 'pace' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *JobRequest) validateBatchSize() (err error) {
	if job.Type == string(config.Sleep) {
		if job.BatchSize != 0 {
			err = errors.New("JobValidationError: field 'batch_size' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *JobRequest) validateDataSize() (err error) {
	if job.Type == string(config.Sleep) {
		if job.DataSize != 0 {
			err = errors.New("JobValidationError: field 'data_size' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (job *JobRequest) validateOperations() (err error) {
	if job.Type == string(config.Sleep) {
		if job.Operations != 0 {
			err = errors.New("JobValidationError: field 'operations' must be equal 0 or must be not set for job with 'sleep' type ")
		}
	}
	return
}

func (rp *ReportingFormatRequest) Validate() error {
	validators := []func() error{
		rp.validateReportingFormat,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (rpt *ReportingFormatRequest) validateReportingFormat() (err error) {
	return nil
}

// todo: add schema validation
// schema keys
// save key should be in schema

// todo: validation job type
// todo: validation duration and opertions cannot be set together

// func Contains[T comparable, X comparable](array []T, comparator X, predicate func(T, X) bool) bool {
// 	for _, elem := range array {
// 		if predicate(elem, comparator) {
// 			return true
// 		}
// 	}
// 	return false
// }

func Contains[T comparable](array []T, predicate func(T) bool) bool {
	for _, elem := range array {
		if predicate(elem) {
			return true
		}
	}
	return false
}
