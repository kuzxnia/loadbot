package main

import (
	"fmt"

	"github.com/kuzxnia/loadbot/lbot"
	applog "github.com/kuzxnia/loadbot/lbot/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ConfigFile = "config-file"

func BuildArgs(logger *log.Entry, version string, commit string, date string) *cobra.Command {
	cmd := cobra.Command{
		Use:     "lbot-agent",
		Short:   "Database workload driver ",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			loglvl, _ := f.GetString(lbot.FlagLogLevel)
			logfmt, _ := f.GetString(lbot.FlagLogFormat)
			err := applog.Configure(logger, loglvl, logfmt)
			if err != nil {
				return fmt.Errorf("failed to configure logger: %w", err)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()

			configFilePath, _ := flags.GetString(ConfigFile)

			agent := lbot.NewAgent(cmd.Context(), logger)
			if configFilePath != "" {
        agent.ApplyConfig(configFilePath)
			}

			agent.Listen()
		},
	}

	pf := cmd.PersistentFlags()
	pf.StringP(ConfigFile, "f", "", "Config file for agent")

	return &cmd
}
