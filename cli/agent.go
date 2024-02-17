package cli

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot"
)

// tutaj nie powinno wchodzić proto
func StartAgent(context context.Context, stdin bool, cmdPort string, configFile string) (err error) {
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

	if cmdPort != "" {
		requestConfig.AgentPort = cmdPort
	}
	agent := lbot.NewAgent(context)
	if requestConfig != nil {
		agent.ApplyConfig(requestConfig)
	}
	agent.Listen()
	return nil
}
