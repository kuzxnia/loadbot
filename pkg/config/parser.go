package config

import (
	"encoding/json"
	"time"
)

func (c *Job) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Database    string `json:"database"`
		Collection  string `json:"collection"`
		Schema      string `json:"template"`
		Connections uint64 `json:"connections"`
		Pace        uint64 `json:"pace"`
		DataSize    uint64 `json:"data_size"`
		BatchSize   uint64 `json:"batch_size"`
		Duration    string `json:"duration"`
		Operations  uint64 `json:"operations"`
		Timeout     string `json:"timeout"` // if not set, default
		// Params ex. for read / update
		//     * filter: { "_id": "#_id"}
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

	return
}

// type parser
