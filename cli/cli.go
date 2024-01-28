package cli

import (
	"fmt"
	"strings"
	"time"

	applog "github.com/kuzxnia/loadbot/cli/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

var Logger *log.Entry

func New(rootLogger *log.Entry, version string, commit string, date string) *cobra.Command {
	Logger = rootLogger

	cmd := cobra.Command{
		Use:     "lbot",
		Short:   "A command-line database workload ",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			loglvl, _ := f.GetString(FlagLogLevel)
			logfmt, _ := f.GetString(FlagLogFormat)
			err := applog.Configure(Logger, loglvl, logfmt)
			if err != nil {
				return fmt.Errorf("failed to configure logger: %w", err)
			}

			return nil
		},
	}


  // by default run in docker container
  // agent save config file in /tmp/lbot/ .* 
  // if you want to change file you need to reconfigure or kill process and start

  // todo: validate connection to agent when calling without args
  // default localhost
  // add arg param agent-uri, if agent is somewhere else


  // jeśli stworzone było lokalnie to bij do lokalnego, 
  // jeśli na k8s to bijesz po k8s-selector, jeśli wiele to bijesz do wielu, 

	pf := cmd.PersistentFlags()
	pf.String(FlagLogLevel, applog.LevelInfo, fmt.Sprintf("log level, must be one of: %s", strings.Join(applog.Levels, ", ")))
	pf.String(FlagLogFormat, applog.FormatFancy, fmt.Sprintf("log format, must be one of: %s", strings.Join(applog.Formats, ", ")))

	cmd.AddGroup(&OrchiestrationGroup)
	cmd.AddCommand(provideOrchiestrationCommands()...)
	cmd.AddGroup(&DriverGroup)
	cmd.AddCommand(provideDriverCommands()...)

	return &cmd
}

const (
	CommandStartDriver = "start"
	CommandStopDriver  = "stop"
	CommandWatchDriver = "watch"
)

var DriverGroup = cobra.Group{
	ID:    "driver",
	Title: "Driver Commands:",
}

func provideDriverCommands() []*cobra.Command {
	startCommand := cobra.Command{
		Use:     CommandStartDriver,
		Aliases: []string{"i"},
		Short:   "Start stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    startingDriverHandler,
		GroupID: DriverGroup.ID,
	}

	stopCommand := cobra.Command{
		Use:     CommandStopDriver,
		Aliases: []string{"i"},
		Short:   "Stopping stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    stoppingDriverHandler,
		GroupID: DriverGroup.ID,
	}

	watchCommand := cobra.Command{
		Use:     CommandWatchDriver,
		Aliases: []string{"i"},
		Short:   "Watch stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    watchingDriverHandler,
		GroupID: DriverGroup.ID,
	}

	return []*cobra.Command{&startCommand, &stopCommand, &watchCommand}
}

// todo: generate complection

const (
	CommandInstall   = "install"
	CommandUnInstall = "uninstall"

	// if not set will install localy without k8s
	FlagSourceKubeconfig = "k8s-config"
	FlagSourceContext    = "k8s-context"
	FlagSourceNamespace  = "k8s-namespace"

	FlagHelmTimeout   = "helm-timeout"
	FlagHelmValues    = "helm-values"
	FlagHelmSet       = "helm-set"
	FlagHelmSetString = "helm-set-string"
	FlagHelmSetFile   = "helm-set-file"

	installationArgsNum = 1
)

var OrchiestrationGroup = cobra.Group{
	ID:    "orchiestration",
	Title: "Resource Orchiestration Commands:",
}

func provideOrchiestrationCommands() []*cobra.Command {
	installationCommand := cobra.Command{
		Use:     CommandInstall + " <config-file>",
		Aliases: []string{"i"},
		Short:   "Install workload driver with helm charts on k8s or only with docker locally",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
		GroupID: OrchiestrationGroup.ID,
	}

	flags := installationCommand.Flags()
	// flags
	flags.StringP(FlagSourceKubeconfig, "k", "", "path of the kubeconfig file of the source PVC")
	flags.StringP(FlagSourceContext, "c", "", "context in the kubeconfig file of the source PVC")
	flags.StringP(FlagSourceNamespace, "n", "", "namespace of the source PVC")

	flags.DurationP(FlagHelmTimeout, "t", 1*time.Minute, "install/uninstall timeout for helm releases")
	flags.StringSliceP(FlagHelmValues, "f", nil, "set additional Helm values by a YAML file or a URL (can specify multiple)")
	flags.StringSlice(FlagHelmSet, nil, "set additional Helm values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	flags.StringSlice(FlagHelmSetString, nil, "set additional Helm STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	flags.StringSlice(FlagHelmSetFile, nil, "set additional Helm values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")

	// if no flags provided, install as local, simply start
	// skipCleanup
	// helmTimeout
	// helmValues

	unInstallationCommand := cobra.Command{
		// todo: where to keep configuration? there will be couple workloads at the same time
		Use:     CommandUnInstall,
		Aliases: []string{"i"},
		Short:   "Uninstall workload driver",
		RunE:    unInstallationHandler,
		GroupID: OrchiestrationGroup.ID,
	}

	return []*cobra.Command{&installationCommand, &unInstallationCommand}
}
