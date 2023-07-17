package args

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/kuzxnia/mongoload/pkg/config"
)

var defaultArgsParser = NewArgsParser()

type CLI struct {
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

func Parse() (*config.Config, error) {
	return defaultArgsParser.Parse()
}

type (
	Parser     func(*CLI) *config.Config
	FileParser func(*CLI) (*config.Config, error)
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
	cli := CLI{}
	kong.Parse(&cli)

	if cli.ConfigFile != "" {
		cfg, error = ap.configFileParser(&cli)
	} else {
		cfg = ap.commandLineParser(&cli)
	}

	if error != nil {
		return nil, error
	}

	error = cfg.Validate()

	return
}

func ParseCommandLineArgs(cli *CLI) *config.Config {
	jobs := []*config.Job{
		{
			Connections: cli.Connections,
			Pace:        cli.Pace,
			Duration:    cli.Duration,
			Operations:  cli.Operations,
			BatchSize:   cli.BatchSize,
			Timeout:     cli.Timeout,
		},
	}
	cfg := config.Config{
		ConnectionString: cli.ConnectionString,
		Debug:            cli.Debug,
	}
	// todo: move building relationships to different layer
	for _, job := range jobs {
		job.Parent = &cfg
	}
	return &cfg
}

func ParseFileConfigArgs(cli *CLI) (*config.Config, error) {
	content, err := os.ReadFile(cli.ConfigFile)
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
