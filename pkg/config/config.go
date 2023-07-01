package config

import (
	"fmt"
	"time"
)

// to refactor
type Config struct {
	MongoURI              string
	MongoDatabase         string
	MongoCollection       string
	ConcurrentConnections uint64
	PoolSize              uint64
	RpsLimit              int
	DurationLimit         time.Duration
	OpsAmount             int
	BatchSize             uint64
	DataLenght            uint64
	WriteRatio            float64
	Timeout               time.Duration
}

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateWriteRatio,
	}

	for _, validate := range validators {
		if error := validate(); error != nil {
			return error
		}
	}
	return nil
}

func (c *Config) validateWriteRatio() error {
	if c.WriteRatio <= 0.0 || c.WriteRatio > 1.0 {
		return fmt.Errorf("Write ratio must be in range 0..1")
	}
	return nil
}

// add more validators
