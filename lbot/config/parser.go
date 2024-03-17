package config

import (
	"encoding/json"

	"github.com/tailscale/hujson"
)

// func ParseConfigFile(configFile string) (*Config, error) {
// 	content, err := os.ReadFile(configFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	content, err = standardizeJSON(content)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var cfg Config
// 	err = json.Unmarshal(content, &cfg)

// 	if err != nil {
// 		return nil, errors.New("Error during Unmarshal(): " + err.Error())
// 	}
// 	return &cfg, err
// }

func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func (c *Job) UnmarshalJSON(data []byte) (err error) {
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
		Duration    Duration               `json:"duration"`
		Operations  uint64                 `json:"operations"`
		Timeout     Duration               `json:"timeout"` // if not set, default
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
