package config

import (
	"encoding/json"
	"time"
)

// global config -- to rename as config


func (c *Job) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name        string
		Type        string
		Template    string
		Connections uint64 // Maximum number of concurrent connections
		Pace        uint64 // rps limit / peace - if not set max
		DataSize    uint64 // data size in bytes
		BatchSize   uint64
		Duration    string
		Operations  uint64
		Timeout     string // if not set, default
		// Params ex. for read / update
		//     * filter: { "_id": "#_id"}
	}

	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if c.Duration, err = time.ParseDuration(tmp.Duration); err != nil {
		return err
	}
	c.Timeout, err = time.ParseDuration(tmp.Timeout)

	return
}

// type parser
