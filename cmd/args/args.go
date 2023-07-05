package args

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/kuzxnia/mongoload/pkg/config"
)

var CLI struct {
	// mongo params
	MongoURI              string        `arg:"" help:"Database hostname url" default:"mongodb://localhost:27017"`
	MongoDatabase         string        `help:"Database name" default:"load_test"`
	MongoCollection       string        `help:"Collection name" default:"load_test_coll"`
	PoolSize              uint64        `help:"Active connections pool size(default: 0 - no limit)" default:"0"`
	ConcurrentConnections uint64        `short:"c" help:"Concurrent connections amount" default:"100"`
	RpsLimit              uint64        `name:"rps" help:"RPS limit"`
	DurationLimit         time.Duration `short:"d" name:"duration" help:"Duration limit (ex. 10s, 5m, 1h)"`
	OpsAmount             uint64        `short:"o" name:"operations" help:"Operations (read/write/update) to perform"`
	BatchSize             uint64        `short:"b" help:"Batch size"`
	DataLenght            uint64        `short:"s" help:"Lenght of single item data(chars)" default:"100"`
	WriteRatio            uint64        `short:"w" help:"Write ratio"`
	ReadRatio             uint64        `short:"r" help:"Read ratio"`
	UpdateRatio           uint64        `short:"u" help:"Update ratio"`
	Timeout               time.Duration `short:"t" help:"Timeout for requests" default:"5s"`
	Debug                 bool          `help:"Displaying additional diagnostic information" default:"false"`
	DebugFile             string        `type:"path" help:"Redirection debug information to file"`
}

func Parse() (*config.Config, error) {
	kong.Parse(&CLI)

	wr, rr, ur := 0, 0, 0

	if ratioFactor := CLI.WriteRatio + CLI.ReadRatio + CLI.UpdateRatio; ratioFactor == 0 {
		wr = 100
	} else {
		wr = int(float64(CLI.WriteRatio) / float64(ratioFactor) * 100)
		rr = int(float64(CLI.ReadRatio)/float64(ratioFactor)*100) + wr
		ur = int(float64(CLI.UpdateRatio)/float64(ratioFactor)*100) + rr
	}
	fmt.Println(wr, rr, ur)

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
		DataLenght:            CLI.DataLenght, // chagne to datasize
		WriteRatio:            uint64(wr),
		ReadRatio:             uint64(rr),
		UpdateRatio:           uint64(ur),
		Timeout:               CLI.Timeout,
		Debug:                 CLI.Debug || bool(CLI.DebugFile != ""),
		DebugFilePath:         CLI.DebugFile,
	}
	error := cfg.Validate()
	if error != nil {
		return nil, error
	}
	return &cfg, nil
}
