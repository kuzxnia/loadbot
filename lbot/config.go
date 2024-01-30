package lbot

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/tailscale/hujson"
	"golang.org/x/net/context"
)

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
