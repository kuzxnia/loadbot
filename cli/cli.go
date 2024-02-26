package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/kuzxnia/loadbot/cli/workload"
	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func New(version string, commit string, date string) *cobra.Command {
	cmd := cobra.Command{
		Use:     "loadbot",
		Short:   "A command-line database workload driver ",
		Version: fmt.Sprintf("%s (commit: %s) (build date: %s)", version, commit, date),
	}
	cmd.AddCommand(
		provideWorkloadCommands(),
		provideAgentCommand(),
	)
	cmd.Root().CompletionOptions.HiddenDefaultCmd = true

	return &cmd
}

var (
	Conn                       *grpc.ClientConn
	DefaultProgressInterval, _ = time.ParseDuration("200ms")
)

const (
	WorkloadRootCommand = "workload"

	CommandStartWorkload    = "start"
	CommandStopWorkload     = "stop"
	CommandWatchWorkload    = "watch"
	CommandProgressWorkload = "progress"
	CommandConfigWorkload   = "config"

	// config args
	ConfigFile = "config-file"
	AgentUri   = "agent-uri"
	Interval   = "interval"
	StdIn      = "stdin"
)

func provideWorkloadCommands() *cobra.Command {
	persistentPreRunE := func(cmd *cobra.Command, args []string) (err error) {
		f := cmd.Flags()
		agentUri, _ := f.GetString(AgentUri)
		Conn, err = grpc.Dial(agentUri, grpc.WithInsecure())
		// valiedate connection
		if err != nil {
			log.Fatal("Found errors trying to connect to loadbot-agent:", err)
			return
		}
		return
	}
	persistentPostRun := func(cmd *cobra.Command, args []string) {
		Conn.Close()
	}
	workloadRootCommand := cobra.Command{
		Use:               WorkloadRootCommand,
		Short:             "Start workload",
		PersistentPreRunE: persistentPreRunE,
		PersistentPostRun: persistentPostRun,
	}

	startCommand := cobra.Command{
		Use:               CommandStartWorkload,
		Short:             "Start workload",
		PersistentPreRunE: persistentPreRunE,
		PersistentPostRun: persistentPostRun,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()

			progress, _ := flags.GetBool("progress")
			interval, _ := flags.GetDuration(Interval)

			if progress {
				request := proto.StartWithProgressRequest{
					RefreshInterval: interval.String(),
				}
				return workload.StartWorkloadWithProgress(Conn, &request)
			} else {
				// todo: switch to local model aka cli.StartRequest
				request := proto.StartRequest{
					Watch: false,
				}

				return workload.StartWorkload(Conn, &request)
			}
		},
	}

	startCommandFlags := startCommand.Flags()
	startCommandFlags.BoolP("progress", "p", false, "Show progress of stress test")
	startCommandFlags.DurationP(Interval, "i", DefaultProgressInterval, "Progress refresh interval")
	// todo: add parent command and inherit this flag
	startCommandFlags.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

	stopCommand := cobra.Command{
		Use:               CommandStopWorkload,
		Short:             "Stop workload",
		PersistentPreRunE: persistentPreRunE,
		PersistentPostRun: persistentPostRun,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// todo: switch to local model aka cli.StartRequest
			request := proto.StopRequest{}
			// response model could have worlkload id?

			return workload.StopWorkload(Conn, &request)
		},
	}
	stopCommandFlags := stopCommand.Flags()
	stopCommandFlags.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

	watchCommand := cobra.Command{
		Use:               CommandWatchWorkload,
		Short:             "Watch stress test",
		PersistentPreRunE: persistentPreRunE,
		PersistentPostRun: persistentPostRun,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// building parameters for stop
			// check for params

			// todo: switch to local model aka cli.StartRequest
			request := proto.WatchRequest{}
			// response model could have worlkload id?

			return workload.WatchWorkload(Conn, &request)
		},
	}
	watchCommandFlags := watchCommand.Flags()
	watchCommandFlags.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

	progressCommand := cobra.Command{
		Use:               CommandProgressWorkload,
		Short:             "Watch workload progress",
		PersistentPreRunE: persistentPreRunE,
		PersistentPostRun: persistentPostRun,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			flags := cmd.Flags()
			interval, _ := flags.GetDuration(Interval)

			request := proto.ProgressRequest{
				RefreshInterval: interval.String(),
			}

			return workload.WorkloadProgress(Conn, &request)
		},
	}
	progressCommandFlags := progressCommand.Flags()
	progressCommandFlags.DurationP(Interval, "i", DefaultProgressInterval, "Progress refresh interval")
	progressCommandFlags.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

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

			return workload.SetConfigWorkload(Conn, parsedConfig)
		},
	}
	configCommandFlags := configCommand.Flags()
	configCommandFlags.StringP(ConfigFile, "f", "", "file with workload configuration")
	configCommandFlags.Bool(StdIn, false, "get workload configuration from stdin")
	configCommandFlags.StringP(AgentUri, "u", "127.0.0.1:1234", "loadbot agent uri (default: 127.0.0.1:1234)")

	workloadRootCommand.AddCommand(&startCommand, &stopCommand, &watchCommand, &configCommand, &progressCommand)
	return &workloadRootCommand
}

const (
	AgentRootCommand  = "agent"
	AgentStartCommand = "start"

	// agent args
	AgentName                    = "name"
	AgentPort                    = "port"
	MetricsExportUrl             = "metrics_export_url"
	MetricsExportIntervalSeconds = "metrics_export_interval_seconds"
	MetricsExportPort            = "metrics_export_port"
)

func provideAgentCommand() *cobra.Command {
	agentRootCommand := cobra.Command{
		Use:   AgentRootCommand,
		Short: "Agent Commands",
	}

	startAgentCommand := cobra.Command{
		Use:   AgentStartCommand,
		Short: "Start agent",
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

			return StartAgent(
				cmd.Context(), agentConfig, stdin, configFile,
				// cluster agrs
				// snapshotDir, internalCommunicationPort, nodeID, initCluster,
			)
		},
	}

	flags := startAgentCommand.Flags()
	flags.StringP(AgentName, "n", "", "Agent name")
	flags.StringP(ConfigFile, "f", "", "Config file for loadbot-agent")
	flags.Bool(StdIn, false, "Provide configuration from stdin.")
	flags.StringP(AgentPort, "p", "", "Agent port")
	flags.String(MetricsExportUrl, "", "Prometheus export url used for pushing metrics")
	flags.Uint64(MetricsExportIntervalSeconds, 0, "Prometheus export push interval")
	flags.String(MetricsExportPort, "", "Expose metrics on port instead pushing to prometheus")

	agentRootCommand.AddCommand(&startAgentCommand)

	return &agentRootCommand
}

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

func provideOrchiestrationCommands() []*cobra.Command {
	installationCommand := cobra.Command{
		Use:     CommandInstall + " <config-file>",
		Aliases: []string{"i"},
		Short:   "Install workload driver with helm charts on k8s or only with docker locally",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
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

	unInstallationCommand := cobra.Command{
		// todo: where to keep configuration? there will be couple workloads at the same time
		Use:     CommandUnInstall,
		Aliases: []string{"i"},
		Short:   "Uninstall workload driver",
		RunE:    unInstallationHandler,
	}

	return []*cobra.Command{&installationCommand, &unInstallationCommand}
}

// todo: generate complection
