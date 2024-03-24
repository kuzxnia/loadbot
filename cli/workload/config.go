package workload

// 1. without args, just prints configfile

// 2. with --set= update config

// 3. generate samle config, for mongodb, postgres itp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/emptypb"
)

// checks if process is running in local system
// here should be cli config request - not lbot one
func SetWorkloadConfig(conn grpc.ClientConnInterface, parsedConfig *lbot.ConfigRequest) (err error) {
	requestConfig := BuildConfigRequest(parsedConfig)

	fmt.Println("🚀 Setting new config")

	client := proto.NewConfigServiceClient(conn)
	_, err = client.SetConfig(context.TODO(), requestConfig)
	if err != nil {
		return fmt.Errorf("Setting config failed: %w", err)
	}
	fmt.Println("✅ Setting config succeeded")

	return
}

func GetWorkloadConfig(conn grpc.ClientConnInterface) (err error) {
	client := proto.NewConfigServiceClient(conn)
	cfg, err := client.GetConfig(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("Getting config failed: %w", err)
	}

	fmt.Println(prototext.MarshalOptions{Multiline: true}.Format(cfg))

	return
}

func GenerateConfigWorkload() (err error) {
	cfg := lbot.ConfigRequest{}
	cfg.ConnectionString = GetStringInput("Please provide connection string to database:", "mongodb://admin:pass@127.0.0.1:27017")

	metrics := GetSelectWithDescriptionInput(
		"Do you want to export metrics?",
		[]choose.Choice{
			{Text: "no", Note: "Metrics won't be exported"},
			{Text: "on port", Note: "Export metrics on provided port"},
			{Text: "push", Note: "Push metrics to prometheus url with interval"},
		},
	)
	cfg.Agent = &lbot.AgentRequest{}
	switch metrics {
	case "no": // just skip
	case "on port":
		cfg.Agent.Name = GetStringInput("Provide agent name(used for metrics labeling):", "")
		cfg.Agent.MetricsExportPort = GetStringInput("Provide metrics export port:", "6060")
	case "push":
		cfg.Agent.Name = GetStringInput("Provide agent name(used for metrics labeling):", "")
		cfg.Agent.MetricsExportUrl = GetStringInput("Provide metrics import url:", "http://victoria-metrics:8428/api/v1/import/prometheus")
		cfg.Agent.MetricsExportIntervalSeconds = uint64(GetDurationInput("Provide metrics push interval:", "5s").Seconds())
	}

	job := lbot.JobRequest{}
	cfg.Jobs = []*lbot.JobRequest{&job}
	job.Type = GetSelectInput(
		"Choose job type",
		[]string{string(config.Write), string(config.BulkWrite), string(config.Read), string(config.Update)},
		choose.WithTheme(choose.ThemeLine),
		choose.WithKeyMap(choose.HorizontalKeyMap),
	)
	if job.Type == string(config.BulkWrite) {
		job.BatchSize = uint64(GetNumberInput("Provide batch size:", "100"))
	}
	job.Name = GetStringInput("Provide workload name:", "Test workload")
	job.Database = GetStringInput("Provide database name:", "tmp-database")
	job.Collection = GetStringInput("Provide collection name:", "tmp-collection")
	job.Pace = uint64(GetNumberInput("How many requests per second?", "100"))
	job.Connections = uint64(GetNumberInput("Provide workload concurency:", "100"))
	job.DataSize = uint64(GetNumberInput("How big document should be? (whole document will be <data_size> + 33 in bytes):", "250"))

	limitType := GetSelectWithDescriptionInput(
		"Choose workload limit type:",
		[]choose.Choice{
			{Text: "operations", Note: "Executes number of operations"},
			{Text: "duration", Note: "Runs workload until time ends"},
			{Text: "infinite", Note: "Runs workload endlessly"},
		},
	)
	switch limitType {
	case "operations":
		job.Operations = uint64(GetNumberInput("How many operations to perform?", "100"))
	case "duration":
		job.Duration = *GetDurationInput("How long workload should take?", "10m30s")
	default:
		// set noting to duration or operations
	}

	filename := GetStringInput("Where to save config file:", "workload_config.json")
	if _, error := os.Stat(filename); !os.IsNotExist(error) {
		overwrite := GetBooleanInput("File exists! Do you want to overwrite?")
		if !overwrite {
			fmt.Println("Exiting...")
			os.Exit(1)
		}
	}

	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0o644)
	if err != nil {
		return err
	}

	return
}

func CheckPromptErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

func GetStringInput(label string, defaultValue string) string {
	result, err := prompt.New().Ask(label).Input(defaultValue, input.WithHelp(true))
	CheckPromptErr(err)
	return result
}

func GetBooleanInput(label string) bool {
	result := GetSelectInput(
		label,
		[]string{"t", "F"},
		choose.WithTheme(choose.ThemeLine),
		choose.WithKeyMap(choose.HorizontalKeyMap),
	)
	resultBool, err := strconv.ParseBool(result)
	CheckPromptErr(err)
	return resultBool
}

func GetNumberInput(label string, defaultValue string) int {
	result, err := prompt.New().Ask(label).Input(defaultValue, input.WithInputMode(input.InputInteger), input.WithHelp(true))
	CheckPromptErr(err)
	resultInt, err := strconv.Atoi(result)
	CheckPromptErr(err)
	return resultInt
}

func GetSelectInput(label string, options []string, ops ...choose.Option) string {
	result, err := prompt.New().Ask(label).Choose(options, ops...)
	CheckPromptErr(err)
	return result
}

func GetSelectWithDescriptionInput(label string, options []choose.Choice, ops ...choose.Option) string {
	result, err := prompt.New().Ask(label).AdvancedChoose(options, ops...)
	CheckPromptErr(err)
	return result
}

func GetDurationInput(label string, defaultValue string) *time.Duration {
	result, err := prompt.New().Ask(label).Input(defaultValue, input.WithHelp(true))
	CheckPromptErr(err)
	resultDuration, err := time.ParseDuration(result)
	CheckPromptErr(err)
	return &resultDuration
}

func BuildConfigRequest(request *lbot.ConfigRequest) *proto.ConfigRequest {
	cfg := &proto.ConfigRequest{
		ConnectionString: request.ConnectionString,
		Agent: &proto.AgentRequest{
			Name:                         request.Agent.Name,
			Port:                         request.Agent.Port,
			MetricsExportUrl:             request.Agent.MetricsExportUrl,
			MetricsExportIntervalSeconds: request.Agent.MetricsExportIntervalSeconds,
			MetricsExportPort:            request.Agent.MetricsExportPort,
		},
		Jobs:    make([]*proto.JobRequest, len(request.Jobs)),
		Schemas: make([]*proto.SchemaRequest, len(request.Schemas)),
		Debug:   request.Debug,
	}
	for i, job := range request.Jobs {
		cfg.Jobs[i] = &proto.JobRequest{
			Name:        job.Name,
			Database:    job.Database,
			Collection:  job.Collection,
			Type:        job.Type,
			Schema:      job.Schema,
			Connections: job.Connections,
			Pace:        job.Pace,
			DataSize:    job.DataSize,
			BatchSize:   job.BatchSize,
			Duration:    job.Duration.String(),
			Operations:  job.Operations,
			Timeout:     job.Timeout.String(),
			// todo: setup filters and schema inside
			// Filter:          job.Filter,
		}
	}
	for i, schema := range request.Schemas {
		cfg.Schemas[i] = &proto.SchemaRequest{
			Name:       schema.Name,
			Database:   schema.Database,
			Collection: schema.Collection,
			// Schema:     schema.Schema,
			Save: schema.Save,
		}
	}

	return cfg
}

// todo: command for setting only one field
// ex. --set=cos.tam.tam=2
