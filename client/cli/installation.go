package cli

import (
	"fmt"
	"time"

	"github.com/kuzxnia/loadcli/installation"
	"github.com/spf13/cobra"
)

const (
	CommandInstall = "install"

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

func provideInstallationHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandInstall + " <config-file>",
		Aliases: []string{"i"},
		Short:   "Install workload driver",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
	}

	flags := cmd.Flags()
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

	return &cmd
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

	request := installation.Request{
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

	if err := installation.NewInstallationProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	logger.Info("✅ Installation process succeeded")

	return nil
}
