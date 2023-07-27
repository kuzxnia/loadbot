package config

import (
	"encoding/json"
	"time"
)

func (c *Job) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name            string   `json:"name"`
		Type            string   `json:"type"`
		Database        string   `json:"database"`
		Collection      string   `json:"collection"`
		Schema          string   `json:"template"`
		ReportingFormat string   `json:"format"`
		Connections     uint64   `json:"connections"`
		Pace            uint64   `json:"pace"`
		DataSize        uint64   `json:"data_size"`
		BatchSize       uint64   `json:"batch_size"`
		Duration        string   `json:"duration"`
		Operations      uint64   `json:"operations"`
		Timeout         string   `json:"timeout"` // if not set, default
		Indexes         []*Index `json:"indexes"` // if not set, default
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
	c.Indexes = tmp.Indexes

  if c.Type == string(CreateIndex) {
    c.Operations = 1
  }

	return
}

func (c *ReportingFormat) UnmarshalJSON(data []byte) (err error) {
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
