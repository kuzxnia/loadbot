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
		Short:   "A command-line database workload ",
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

	cmd.AddGroup(&OrchiestrationGroup)
	cmd.AddCommand(provideOrchiestrationCommands()...)
	cmd.AddGroup(&DriverGroup)
	cmd.AddCommand(provideDriverCommands()...)

	// _ = cmd.RegisterFlagCompletionFunc(FlagLogLevel, buildStaticSliceCompletionFunc(applog.Levels))
	// _ = cmd.RegisterFlagCompletionFunc(FlagLogFormat, buildStaticSliceCompletionFunc(applog.Formats))

	return &cmd
}

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
	CommandStartDriver = "start"
	CommandStopDriver  = "stop"
	CommandWatchDriver = "watch"
)

func startingDriverHandler(cmd *cobra.Command, args []string) error {
	// flags := cmd.Flags()

	request := StartRequest{}

	logger.Info("🚀 Starting stress test")

	if err := NewStartProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	logger.Info("✅ Starting stress test succeeded")

	return nil
}

func stoppingDriverHandler(cmd *cobra.Command, args []string) error {
	request := StoppingRequest{}

	logger.Info("🚀 Stopping stress test")

	if err := NewStoppingProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("Stopping stress test failed: %w", err)
	}

	logger.Info("✅ Stopping stress test succeeded")

	return nil
}

func watchingDriverHandler(cmd *cobra.Command, args []string) error {
	request := WatchingRequest{}

	logger.Info("🚀 Watching stress test")

	if err := NewWatchingProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("Watching stress test failed: %w", err)
	}

	logger.Info("✅ Watching stress test succeeded")

	return nil
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

func installationHandler(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	srcKubeconfigPath, _ := flags.GetString(FlagSourceKubeconfig)
	srcContext, _ := flags.GetString(FlagSourceContext)
	srcNS, _ := flags.GetString(FlagSourceNamespace)

	helmTimeout, _ := flags.GetDuration(FlagHelmTimeout)
	helmValues, _ := flags.GetStringSlice(FlagHelmValues)
	helmSet, _ := flags.GetStringSlice(FlagHelmSet)
	helmSetString, _ := flags.GetStringSlice(FlagHelmSetString)
	helmSetFile, _ := flags.GetStringSlice(FlagHelmSetFile)

	request := InstallationRequest{
		KubeconfigPath:   srcKubeconfigPath,
		Context:          srcContext,
		Namespace:        srcNS,
		HelmTimeout:      helmTimeout,
		HelmValuesFiles:  helmValues,
		HelmValues:       helmSet,
		HelmStringValues: helmSetString,
		HelmFileValues:   helmSetFile,
	}

	logger.Info("🚀 Starting installation process")

	if err := NewInstallationProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	logger.Info("✅ Installation process succeeded")

	return nil
}

func unInstallationHandler(cmd *cobra.Command, args []string) error {
	request := UnInstallationRequest{}

	logger.Info("🚀 Starting installation process")

	if err := NewUnInstallationProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	logger.Info("✅ Installation process succeeded")

	return nil
}
