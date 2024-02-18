package cli

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/samber/lo"
)

// tutaj nie powinno wchodzić proto
func StartAgent(context context.Context, config *lbot.AgentRequest, stdin bool, configFile string) (err error) {
	var requestConfig *lbot.ConfigRequest

	if stdin {
		requestConfig, err = lbot.ParseStdInConfig()
		if err != nil {
			return err
		}
	}

	if configFile != "" {
		requestConfig, err = lbot.ParseConfigFile(configFile)
		if err != nil {
			return err
		}
	}

	if lo.IsNil(requestConfig) {
		requestConfig = &lbot.ConfigRequest{}
	}
	if lo.IsNil(requestConfig.Agent) {
		requestConfig.Agent = &lbot.AgentRequest{}
	}
	if lo.IsEmpty(requestConfig.Agent.MetricsExportIntervalSeconds) {
		requestConfig.Agent.MetricsExportIntervalSeconds = 10
	}
	if lo.IsEmpty(requestConfig.Agent.Port) {
		requestConfig.Agent.Port = "1234"
	}

	if lo.IsNotEmpty(config.Name) {
		requestConfig.Agent.Name = config.Name
	}
	if lo.IsNotEmpty(config.Port) {
		requestConfig.Agent.Port = config.Port
	}
	if lo.IsNotEmpty(config.MetricsExportUrl) {
		requestConfig.Agent.MetricsExportUrl = config.MetricsExportUrl
	}
	if lo.IsNotEmpty(config.MetricsExportIntervalSeconds) {
		requestConfig.Agent.MetricsExportIntervalSeconds = config.MetricsExportIntervalSeconds
	}
	if lo.IsNotEmpty(config.MetricsExportPort) {
		requestConfig.Agent.MetricsExportPort = config.MetricsExportPort
	}

	agent := lbot.NewAgent(context)
	if requestConfig != nil {
		agent.ApplyConfig(requestConfig)
	}
	agent.Listen()
	return nil
}
