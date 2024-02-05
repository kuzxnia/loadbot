package cli

// 1. without args, just prints configfile

// 2. with --set= update config

// 3. generate samle config, for mongodb, postgres itp

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

// checks if process is running in local system
func SetConfigDriver(conn *grpc.ClientConn, parsedConfig *lbot.ConfigRequest) (err error) {
	requestConfig := BuildConfigRequest(parsedConfig)

	Logger.Info("🚀 Setting new config")

	client := proto.NewSetConfigProcessClient(conn)
	reply, err := client.Run(context.TODO(), requestConfig)
	if err != nil {
		return fmt.Errorf("Setting config failed: %w", err)
	}

	Logger.Infof("Received: %v", reply)
	Logger.Info("✅ Setting config succeeded")

	return
}

func BuildConfigRequest(request *lbot.ConfigRequest) *proto.ConfigRequest {
	cfg := &proto.ConfigRequest{
		ConnectionString: request.ConnectionString,
		Jobs:             make([]*proto.JobRequest, len(request.Jobs)),
		Schemas:          make([]*proto.SchemaRequest, len(request.Schemas)),
		ReportingFormats: make([]*proto.ReportingFormatRequest, len(request.ReportingFormats)),
		Debug:            request.Debug,
	}
	for i, job := range request.Jobs {
		cfg.Jobs[i] = &proto.JobRequest{
			Name:            job.Name,
			Database:        job.Database,
			Collection:      job.Collection,
			Type:            job.Type,
			Schema:          job.Schema,
			ReportingFormat: job.ReportingFormat,
			Connections:     job.Connections,
			Pace:            job.Pace,
			DataSize:        job.DataSize,
			BatchSize:       job.BatchSize,
			Duration:        job.Duration.String(),
			Operations:      job.Operations,
			Timeout:         job.Timeout.String(),
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
	for i, rf := range request.ReportingFormats {
		cfg.ReportingFormats[i] = &proto.ReportingFormatRequest{
			Name:     rf.Name,
			Interval: rf.Interval.String(),
			Template: rf.Template,
		}
	}

	return cfg
}

// todo: command for setting only one field
// ex. --set=cos.tam.tam=2
