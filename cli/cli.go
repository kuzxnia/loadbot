package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	AgentUri = "agent-uri"
)

var Conn *grpc.ClientConn

func New(version string, commit string, date string) *cobra.Command {
	cmd := cobra.Command{
		Use:     "loadbot",
		Short:   "A command-line database workload driver ",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			f := cmd.Flags()
			// move to driver group
			agentUri, _ := f.GetString(AgentUri)
			Conn, err = grpc.Dial(agentUri, grpc.WithInsecure())
			// valiedate connection
			if err != nil {
				log.Fatal("Found errors trying to connect to loadbot-agent:", err)
				return
			}

			return
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) (err error) {
			defer Conn.Close()

			return
		},
	}
	pf := cmd.PersistentFlags()
	// move to driver group
	pf.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

	// setup supcommands
	// cmd.AddGroup(&OrchiestrationGroup)
	// cmd.AddCommand(provideOrchiestrationCommands()...)
	cmd.AddGroup(&AgentGroup)
	cmd.AddCommand(provideAgentCommands()...)
	cmd.AddGroup(&WorkloadGroup)
	cmd.AddCommand(provideWorkloadCommands()...)

	// by default run in docker container
	// agent save config file in /tmp/lbot/ .*
	// if you want to change file you need to reconfigure or kill process and start

	// todo: validate connection to agent when calling without args
	// default localhost
	// add arg param agent-uri, if agent is somewhere else

	// jeśli stworzone było lokalnie to bij do lokalnego,
	// jeśli na k8s to bijesz po k8s-selector, jeśli wiele to bijesz do wielu,

	return &cmd
}

const (
	CommandStartWorkload    = "start"
	CommandStopWorkload     = "stop"
	CommandWatchWorkload    = "watch"
	CommandProgressWorkload = "progress"
	CommandConfigWorkload   = "config"

	// config args
	ConfigFile = "config-file"
	Interval   = "interval"
	StdIn      = "stdin"
)

var WorkloadGroup = cobra.Group{
	ID:    "workload",
	Title: "Workload Commands:",
}

func provideWorkloadCommands() []*cobra.Command {
	startCommand := cobra.Command{
		Use:   CommandStartWorkload,
		Short: "Start workload",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// building parameters for start
			// check for params
			flags := cmd.Flags()

			progress, _ := flags.GetBool("progress")
			interval, _ := flags.GetDuration(Interval)

			if progress {
				request := proto.StartWithProgressRequest{
					RefreshInterval: interval.String(),
				}
				return StartWithProgressDriver(Conn, &request)
			} else {
				// todo: switch to local model aka cli.StartRequest
				request := proto.StartRequest{
					Watch: false,
				}

				return StartWorkload(Conn, &request)
			}

		},
		GroupID: WorkloadGroup.ID,
	}

	startCommandFlags := startCommand.Flags()
	startCommandFlags.BoolP("progress", "p", false, "Show progress of stress test")
	defaultProgressInterval, _ := time.ParseDuration("1s")
	startCommandFlags.DurationP(Interval, "i", defaultProgressInterval, "Progress refresh interval")

	stopCommand := cobra.Command{
		Use:   CommandStopWorkload,
		Short: "Stop workload",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// building parameters for stop
			// check for params

			// todo: switch to local model aka cli.StartRequest
			request := proto.StopRequest{}
			// response model could have worlkload id?

			return StopWorkload(Conn, &request)
		},
		GroupID: WorkloadGroup.ID,
	}

	watchCommand := cobra.Command{
		Use:     CommandWatchWorkload,
		Aliases: []string{"i"},
		Short:   "Watch stress test",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// building parameters for stop
			// check for params

			// todo: switch to local model aka cli.StartRequest
			request := proto.WatchRequest{}
			// response model could have worlkload id?

			return WatchWorkload(Conn, &request)
		},
		GroupID: WorkloadGroup.ID,
	}
	progressCommand := cobra.Command{
		Use:     CommandProgressWorkload,
		Aliases: []string{"i"},
		Short:   "Watch workload progress",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()
			interval, _ := flags.GetDuration(Interval)

			request := proto.ProgressRequest{
				RefreshInterval: interval.String(),
			}

			return WorkloadProgress(Conn, &request)
		},
		GroupID: WorkloadGroup.ID,
	}
	progressCommandFlags := progressCommand.Flags()
	defaultInterval, _ := time.ParseDuration("1s")
	progressCommandFlags.DurationP(Interval, "i", defaultInterval, "Progress refresh interval")

	configCommand := cobra.Command{
		Use:   CommandConfigWorkload,
		Short: "Config",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()
			configFile, _ := flags.GetString(ConfigFile)
			stdin, _ := flags.GetBool(StdIn)

			if configFile == "" && !stdin {
				return errors.New("You need to provide configuration from either " + ConfigFile + " or " + StdIn)
			}

			var parsedConfig *lbot.ConfigRequest
			if stdin {
				parsedConfig, err = lbot.ParseStdInConfig()
				if err != nil {
					return err
				}
			}

			if configFile != "" {
				parsedConfig, err = lbot.ParseConfigFile(configFile)
				if err != nil {
					return err
				}
			}

			return SetConfigWorkload(Conn, parsedConfig)
		},
		GroupID: WorkloadGroup.ID,
	}
	configCommandFlags := configCommand.Flags()
	configCommandFlags.StringP(ConfigFile, "f", "", "file with workload configuration")
	configCommandFlags.Bool(StdIn, false, "get workload configuration from stdin")

	return []*cobra.Command{&startCommand, &stopCommand, &watchCommand, &configCommand, &progressCommand}
}

const (
	CommandStartAgent = "start-agent"

	// agent args
	AgentName                    = "name"
	AgentPort                    = "port"
	MetricsExportUrl             = "metrics_export_url"
	MetricsExportIntervalSeconds = "metrics_export_interval_seconds"
	MetricsExportPort            = "metrics_export_port"
)

var AgentGroup = cobra.Group{
	ID:    "Agent",
	Title: "Agent Commands:",
}

func provideAgentCommands() []*cobra.Command {
	agentCommand := cobra.Command{
		Use:   CommandStartAgent,
		Short: "Start loadbot-agent",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()

			name, _ := flags.GetString(AgentName)
			port, _ := flags.GetString(AgentPort)
			metricsExportUrl, _ := flags.GetString(MetricsExportUrl)
			metricsExportIntervalSeconds, _ := flags.GetUint64(MetricsExportIntervalSeconds)
			metricsExportPort, _ := flags.GetString(MetricsExportPort)

			agentConfig := &lbot.AgentRequest{
				Name:                         name,
				Port:                         port,
				MetricsExportUrl:             metricsExportUrl,
				MetricsExportIntervalSeconds: metricsExportIntervalSeconds,
				MetricsExportPort:            metricsExportPort,
			}

			configFile, _ := flags.GetString(ConfigFile)
			stdin, _ := flags.GetBool(StdIn)

			return StartAgent(cmd.Context(), agentConfig, stdin, configFile)
		},
		GroupID: AgentGroup.ID,
	}

	flags := agentCommand.Flags()
	flags.StringP(AgentName, "n", "", "Agent name")
	flags.StringP(ConfigFile, "f", "", "Config file for loadbot-agent")
	flags.Bool(StdIn, false, "Provide configuration from stdin.")
	flags.StringP(AgentPort, "p", "", "Agent port")
	flags.String(MetricsExportUrl, "", "Prometheus export url used for pushing metrics")
	flags.Uint64(MetricsExportIntervalSeconds, 0, "Prometheus export push interval")
	flags.String(MetricsExportPort, "", "Expose metrics on port instead pushing to prometheus")

	return []*cobra.Command{&agentCommand}
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
