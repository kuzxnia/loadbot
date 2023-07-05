package args

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/kuzxnia/mongoload/pkg/config"
)

var CLI struct {
	MongoURI              string        `short:"u" name:"uri" help:"Database hostname url" default:"mongodb://localhost:27017"`
	MongoDatabase         string        `help:"Database name" default:"load_test"`
	MongoCollection       string        `help:"Collection name" default:"load_test_coll"`
	ConcurrentConnections uint64        `short:"c" help:"Concurrent connections amount" default:"100"`
	PoolSize              uint64        `help:"Active connections pool size(default: 0 - no limit)" default:"0"`
	RpsLimit              uint64        `name:"rps" help:"RPS limit"`
	DurationLimit         time.Duration `short:"d" name:"duration" help:"Duration limit (ex. 10s, 5m, 1h)"`
	OpsAmount             uint64        `short:"r" name:"requests" help:"Requests to perform"`
	BatchSize             uint64        `short:"b" help:"Batch size"`
	DataLenght            uint64        `short:"s" help:"Lenght of single item data(chars)" default:"100"`
	WriteRatio            float64       `short:"w" help:"Write ratio (ex. 0.2 will result with 20% writes)" default:"0.5"`
	Timeout               time.Duration `short:"t" help:"Timeout for requests" default:"5s"`
}

func Parse() (*config.Config, error) {
	kong.Parse(&CLI)

	cfg := config.Config{
		MongoURI:              CLI.MongoURI,
		MongoDatabase:         CLI.MongoDatabase,
		MongoCollection:       CLI.MongoCollection,
		ConcurrentConnections: CLI.ConcurrentConnections,
		PoolSize:              CLI.PoolSize,
		RpsLimit:              CLI.RpsLimit,
		DurationLimit:         CLI.DurationLimit,
		OpsAmount:             CLI.OpsAmount,
		BatchSize:             CLI.BatchSize,
		DataLenght:            CLI.DataLenght,
		WriteRatio:            CLI.WriteRatio,
		Timeout:               CLI.Timeout,
	}
	error := cfg.Validate()
	if error != nil {
		return nil, error
	}
	return &cfg, nil
}
