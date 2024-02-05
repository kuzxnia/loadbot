package lbot

import (
	"fmt"

	"github.com/kuzxnia/loadbot/lbot/config"
	applog "github.com/kuzxnia/loadbot/lbot/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ConfigFile = "config-file"
	StdIn      = "stdin"
)

func BuildArgs(logger *log.Entry, version string, commit string, date string) *cobra.Command {
	cmd := cobra.Command{
		Use:     "lbot-agent",
		Short:   "Database workload driver ",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			loglvl, _ := f.GetString(config.FlagLogLevel)
			logfmt, _ := f.GetString(config.FlagLogFormat)
			err := applog.Configure(logger, loglvl, logfmt)
			if err != nil {
				return fmt.Errorf("failed to configure logger: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()

			configFile, _ := flags.GetString(ConfigFile)
			stdin, _ := flags.GetBool(StdIn)

			var requestConfig *ConfigRequest

			if stdin {
				requestConfig, err = ParseStdInConfig()
				if err != nil {
					return err
				}
			}

			if configFile != "" {
				requestConfig, err = ParseConfigFile(configFile)
				if err != nil {
					return err
				}
			}

			agent := NewAgent(cmd.Context(), logger)
			if requestConfig != nil {
				agent.ApplyConfig(requestConfig)
			}

			// agent.Listen()
			agent.ListenGRPC()
			return nil
		},
	}

	pf := cmd.PersistentFlags()
	pf.StringP(ConfigFile, "f", "", "Config file for agent")
	pf.Bool(StdIn, false, "get workload configuration from stdin")

	return &cmd
}
