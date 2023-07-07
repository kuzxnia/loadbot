package config

import (
	"fmt"
	"time"
)

type Config struct {
	MongoURI              string
	MongoDatabase         string
	MongoCollection       string
	ConcurrentConnections uint64
	PoolSize              uint64
	RpsLimit              uint64
	DurationLimit         time.Duration
	OpsAmount             uint64
	BatchSize             uint64
	DataLenght            uint64
	WriteRatio            uint64
	ReadRatio             uint64
	UpdateRatio           uint64
	Timeout               time.Duration
	Debug                 bool
	DebugFilePath         string
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
	if c.WriteRatio < 0.0 || c.WriteRatio > 1.0 {
		return fmt.Errorf("Write ratio must be in range 0..1")
	}
	return nil
}

// add more validators
