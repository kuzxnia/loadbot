package cli

import (
	"net/rpc"

	"github.com/kuzxnia/loadbot/orchiestrator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

	// only when k8s is used
	request := orchiestrator.InstallationRequest{
		// install options
		// local-process
		// local-docker
		// k8s

		KubeconfigPath:   srcKubeconfigPath,
		Context:          srcContext,
		Namespace:        srcNS,
		HelmTimeout:      helmTimeout,
		HelmValuesFiles:  helmValues,
		HelmValues:       helmSet,
		HelmStringValues: helmSetString,
		HelmFileValues:   helmSetFile,
	}
	var reply int

	log.Info("🚀 Starting installation process")

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Found errors trying to connect to lbot-agent:", err)
	}

	err = client.Call("InstallationProcess.Run", request, &reply)
	if err != nil {
		log.Fatal("InstallationProcess error:", err)
	}

	log.Info("✅ Installation process succeeded")

	return nil
}
