
package cli

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

	request := lbot.InstallationRequest{
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

  client, err := rpc.NewRpcClient("")
  clie
  err = client.Call("InstallationProcess.Run", args, &reply)

	if err := NewInstallationProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}


	logger.Info("✅ Installation process succeeded")

	return nil
}

