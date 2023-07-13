package config

import (
	"encoding/json"
	"time"
)

func (c *Job) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Name        string
		Type        string
		Template    string
		Connections uint64
		Pace        uint64
		DataSize    uint64
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
