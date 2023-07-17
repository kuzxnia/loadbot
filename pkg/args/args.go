package args

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/kuzxnia/mongoload/pkg/config"
)

var (
	// to rename
	CLI struct {
		ConnectionString string        `arg:"" help:"Database connection string" default:"mongodb://localhost:27017"`
		Connections      uint64        `short:"c" help:"Number of concurrent connections" default:"10"`
		Pace             uint64        `short:"p" name:"pace" help:"Pace - RPS limit"`
		Duration         time.Duration `short:"d" name:"duration" help:"Duration (ex. 10s, 5m, 1h)"`
		Operations       uint64        `short:"o" name:"operations" help:"Operations (read/write/update) to perform"`
		BatchSize        uint64        `short:"b" help:"Batch size"`
		Timeout          time.Duration `short:"t" help:"Timeout for requests" default:"5s"`
		ConfigFile       string        `short:"f" type:"path" help:"Config file path"`
		Debug            bool          `help:"Displaying additional diagnostic information" default:"false"`
	}

	// MongoDatabase   string `help:"Database name" default:"load_test"`
	// MongoCollection string `help:"Collection name" default:"load_test_coll"`
	// PoolSize              uint64        `help:"Active connections pool size(default: 0 - no limit)" default:"0"`
	// DebugFile   string        `type:"path" help:"Redirection debug information to file"`

	// removed in simple only inserts or will be avaliabe by jobType
	// WriteRatio  uint64        `short:"w" help:"Write ratio"`
	// ReadRatio   uint64        `short:"r" help:"Read ratio"`
	// UpdateRatio uint64        `short:"u" help:"Update ratio"`

	// todo: change to data size
	// DataLenght uint64 `short:"s" help:"Lenght of single item data(chars)" default:"100"`

	FileConfigCLI struct {
		ConfigFile string `type:"path" help:"Config file path"`
	}
	defaultArgsParser = NewArgsParser()
)

func Parse() (*config.Config, error) {
	return defaultArgsParser.Parse()
}

type (
	Parser     func(interface{}) *config.Config
	FileParser func(string, interface{}) (*config.Config, error)
)

type ArgsParser struct {
	commandLineParser Parser
	configFileParser  FileParser
}

func NewArgsParser() *ArgsParser {
	return &ArgsParser{
		commandLineParser: ParseCommandLineArgs,
		configFileParser:  ParseFileConfigArgs,
	}
}

func (ap *ArgsParser) Parse() (cfg *config.Config, error error) {
	if kong.Parse(&CLI); CLI.ConfigFile != "" {
		cfg, error = ap.configFileParser(CLI.ConfigFile, &CLI)
	} else {
		cfg = ap.commandLineParser(&CLI)
	}

	if error != nil {
		return nil, error
	}

	error = cfg.Validate()

	return
}

func ParseCommandLineArgs(cli interface{}) *config.Config {
	jobs := []*config.Job{
		{
			Connections: CLI.Connections,
			Pace:        CLI.Pace,
			Duration:    CLI.Duration,
			Operations:  CLI.Operations,
			BatchSize:   CLI.BatchSize,
			Timeout:     CLI.Timeout,
		},
	}
	cfg := config.Config{
		ConnectionString: CLI.ConnectionString,
		Debug:            CLI.Debug,
	}
	for _, job := range jobs {
		job.Parent = &cfg
	}
	return &cfg
}

func ParseFileConfigArgs(filePath string, cli interface{}) (*config.Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	err = json.Unmarshal(content, &cfg)

	if err != nil {
		return nil, errors.New("Error during Unmarshal(): " + err.Error())
	}

	for _, job := range cfg.Jobs {
		job.Parent = &cfg
	}

	return &cfg, err
}
