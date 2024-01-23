package cli

import (
	"fmt"
	"strings"

	applog "github.com/kuzxnia/loadcli/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	appName = "workload"
)

var logger *log.Entry

func New(rootLogger *log.Entry, version string, commit string, date string) *cobra.Command {
	logger = rootLogger

	return buildRootCmd(version, commit, date)
}

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

func buildRootCmd(version string, commit string, date string) *cobra.Command {
	cmd := cobra.Command{
		Use:     appName,
		Short:   "A command-line database workload driver.",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			loglvl, _ := f.GetString(FlagLogLevel)
			logfmt, _ := f.GetString(FlagLogFormat)
			err := applog.Configure(logger, loglvl, logfmt)
			if err != nil {
				return fmt.Errorf("failed to configure logger: %w", err)
			}

			return nil
		},
	}

	pf := cmd.PersistentFlags()
	pf.String(FlagLogLevel, applog.LevelInfo, fmt.Sprintf("log level, must be one of: %s", strings.Join(applog.Levels, ", ")))
	pf.String(FlagLogFormat, applog.FormatFancy, fmt.Sprintf("log format, must be one of: %s", strings.Join(applog.Formats, ", ")))

	cmd.AddCommand(provideInstallationHandler())
	cmd.AddCommand(provideStartDriverHandler())
	cmd.AddCommand(provideStopDriverHandler())
	cmd.AddCommand(provideWatchDriverHandler())

	// _ = cmd.RegisterFlagCompletionFunc(FlagLogLevel, buildStaticSliceCompletionFunc(applog.Levels))
	// _ = cmd.RegisterFlagCompletionFunc(FlagLogFormat, buildStaticSliceCompletionFunc(applog.Formats))

	return &cmd
}

// todo: generate complection
