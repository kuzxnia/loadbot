package lbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"github.com/tailscale/hujson"
)

func NewConfig(request *ConfigRequest) *config.Config {
	cfg := &config.Config{
		ConnectionString: request.ConnectionString,
		Agent: &config.Agent{
			Name:                         request.Agent.Name,
			Port:                         request.Agent.Port,
			MetricsExportUrl:             request.Agent.MetricsExportUrl,
			MetricsExportIntervalSeconds: request.Agent.MetricsExportIntervalSeconds,
			MetricsExportPort:            request.Agent.MetricsExportPort,
		},
		Jobs:    make([]*config.Job, len(request.Jobs)),
		Schemas: make([]*config.Schema, len(request.Schemas)),
		Debug:   request.Debug,
	}
	for i, job := range request.Jobs {
		cfg.Jobs[i] = &config.Job{
			Name: job.Name,
			// Parent:      cfg,
			Database:    job.Database,
			Collection:  job.Collection,
			Type:        job.Type,
			Schema:      job.Schema,
			Connections: job.Connections,
			Pace:        job.Pace,
			DataSize:    job.DataSize,
			BatchSize:   job.BatchSize,
			Duration:    job.Duration,
			Operations:  job.Operations,
			Timeout:     job.Timeout,
			Filter:      job.Filter,
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

	return cfg
}

func NewConfigFromProtoConfigRequest(request *proto.ConfigRequest) *config.Config {
	cfg := &config.Config{
		ConnectionString: request.ConnectionString,
		Agent: &config.Agent{
			Name:                         request.Agent.Name,
			Port:                         request.Agent.Port,
			MetricsExportUrl:             request.Agent.MetricsExportUrl,
			MetricsExportIntervalSeconds: request.Agent.MetricsExportIntervalSeconds,
			MetricsExportPort:            request.Agent.MetricsExportPort,
		},
		Jobs:    make([]*config.Job, len(request.Jobs)),
		Schemas: make([]*config.Schema, len(request.Schemas)),
		Debug:   request.Debug,
	}
	for i, job := range request.Jobs {
		duration, _ := time.ParseDuration(job.Duration)
		timeout, _ := time.ParseDuration(job.Timeout)
		cfg.Jobs[i] = &config.Job{
			Name: job.Name,
			// Parent:      cfg,
			Database:    job.Database,
			Collection:  job.Collection,
			Type:        job.Type,
			Schema:      job.Schema,
			Connections: job.Connections,
			Pace:        job.Pace,
			DataSize:    job.DataSize,
			BatchSize:   job.BatchSize,
			Duration:    duration,
			Operations:  job.Operations,
			Timeout:     timeout,
			// Filter:          job.Filter,
		}
	}
	for i, schema := range request.Schemas {
		cfg.Schemas[i] = &config.Schema{
			Name:       schema.Name,
			Database:   schema.Database,
			Collection: schema.Collection,
			// Schema:     schema.Schema,
			Save: schema.Save,
		}
	}
	return cfg
}

func NewConfigResponseFromConfig(cfg *config.Config) *proto.ConfigResponse {
	response := &proto.ConfigResponse{
		ConnectionString: cfg.ConnectionString,
		Agent: &proto.AgentRequest{
			Name:                         cfg.Agent.Name,
			Port:                         cfg.Agent.Port,
			MetricsExportUrl:             cfg.Agent.MetricsExportUrl,
			MetricsExportIntervalSeconds: cfg.Agent.MetricsExportIntervalSeconds,
			MetricsExportPort:            cfg.Agent.MetricsExportPort,
		},
		Jobs:    make([]*proto.JobRequest, len(cfg.Jobs)),
		Schemas: make([]*proto.SchemaRequest, len(cfg.Schemas)),
		Debug:   cfg.Debug,
	}
	for i, job := range cfg.Jobs {
		response.Jobs[i] = &proto.JobRequest{
			Name: job.Name,
			// Parent:      cfg,
			Database:    job.Database,
			Collection:  job.Collection,
			Type:        job.Type,
			Schema:      job.Schema,
			Connections: job.Connections,
			Pace:        job.Pace,
			DataSize:    job.DataSize,
			BatchSize:   job.BatchSize,
			Duration:    job.Duration.String(),
			Operations:  job.Operations,
			Timeout:     job.Timeout.String(),
			// Filter:          job.Filter,
		}
	}
	for i, schema := range cfg.Schemas {
		response.Schemas[i] = &proto.SchemaRequest{
			Name:       schema.Name,
			Database:   schema.Database,
			Collection: schema.Collection,
			// Schema:     schema.Schema,
			Save: schema.Save,
		}
	}
	return response 
}

// todo: should be pointers
type ConfigRequest struct {
	ConnectionString string           `json:"connection_string"`
	Agent            *AgentRequest    `json:"agent"`
	Jobs             []*JobRequest    `json:"jobs"`
	Schemas          []*SchemaRequest `json:"schemas"`
	Debug            bool             `json:"debug"`
}

// todo: change or even remove,
// yes, remove - there will be multiple agents, and agent will be configred from commandline
// todo: move agentn-name nad add new config flage - custom metrics label or similar
// purpose is to export metrics with cluster name
type AgentRequest struct {
	Name                         string `json:"name"`
	Port                         string `json:"port"`
	MetricsExportUrl             string `json:"metrics_export_url"`
	MetricsExportIntervalSeconds uint64 `json:"metrics_export_interval_seconds"`
	MetricsExportPort            string `json:"metrics_export_port"`
}

type JobRequest struct {
	Name        string                 `json:"name"`
	Database    string                 `json:"database"`
	Collection  string                 `json:"collection"`
	Type        string                 `json:"type"`
	Schema      string                 `json:"schema"`
	Connections uint64                 `json:"connections"`
	Pace        uint64                 `json:"pace"`
	DataSize    uint64                 `json:"data_size"`
	BatchSize   uint64                 `json:"batch_size"`
	Duration    time.Duration          `json:"duration"`
	Operations  uint64                 `json:"operations"`
	Timeout     time.Duration          `json:"timeout"`
	Filter      map[string]interface{} `json:"filter"`
}

type SchemaRequest struct {
	Name       string                 `json:"name"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Schema     map[string]interface{} `json:"schema"` // todo: introducte new type and parse
	Save       []string               `json:"save"`
}

type ConfigService struct {
	proto.UnimplementedConfigServiceServer
	ctx  context.Context
	lbot *Lbot
}

func NewConfigService(ctx context.Context, lbot *Lbot) *ConfigService {
	return &ConfigService{ctx: ctx, lbot: lbot}
}

func (c *ConfigService) SetConfig(ctx context.Context, request *proto.ConfigRequest) (*proto.ConfigResponse, error) {
	cfg := NewConfigFromProtoConfigRequest(request)
	c.lbot.SetConfig(cfg)

	// before configing process it will varify health of cluster, if pods
	return &proto.ConfigResponse{}, nil
}

func (c *ConfigService) GetConfig(ctx context.Context, empty *emptypb.Empty) (*proto.ConfigResponse, error) {
	// todo: should get from db
	response := NewConfigResponseFromConfig(c.lbot.Config)
  fmt.Println("in2")

	return response, nil
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

func InStdInNotEmpty() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, fmt.Errorf("you have an error in stdin:%s", err)
	}

	return (stat.Mode() & os.ModeNamedPipe) == 0, nil
}

func ParseStdInConfig() (*ConfigRequest, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	// move, repetition as above
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
		Name        string                 `json:"name"`
		Type        string                 `json:"type"`
		Database    string                 `json:"database"`
		Collection  string                 `json:"collection"`
		Schema      string                 `json:"template"`
		Connections uint64                 `json:"connections"`
		Pace        uint64                 `json:"pace"`
		DataSize    uint64                 `json:"data_size"`
		BatchSize   uint64                 `json:"batch_size"`
		Duration    config.Duration        `json:"duration"`
		Operations  uint64                 `json:"operations"`
		Timeout     config.Duration        `json:"timeout"` // if not set, default
		Filter      map[string]interface{} `json:"filter"`
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
	c.Connections = tmp.Connections
	c.Pace = tmp.Pace
	c.DataSize = tmp.DataSize
	c.BatchSize = tmp.BatchSize
	c.Duration = tmp.Duration.Duration
	c.Operations = tmp.Operations
	c.Timeout = tmp.Timeout.Duration
	c.Filter = tmp.Filter

	return
}

func (c *ConfigRequest) Values() (string, error) {
	result, err := json.Marshal(c)
	if err != nil {
		return "", nil
	}
	return string(result), nil
}

func (c *ConfigRequest) Validate() error {
	validators := []func() error{
		c.validateJobs,
		// c.validateSchemas,
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

func (job *JobRequest) Validate() error {
	validators := []func() error{
		job.validateSchema,
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
