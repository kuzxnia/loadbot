package cli

// 1. without args, just prints configfile

// 2. with --set= update config

// 3. generate samle config, for mongodb, postgres itp

import (
	"errors"
	"fmt"
	"net/rpc"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/spf13/cobra"
)

// checks if process is running in local system

// command for setting full new config
func setConfigDriverHandler(cmd *cobra.Command, args []string) (err error) {
	var requestConfig *lbot.ConfigRequest

	flags := cmd.Flags()
	configFile, _ := flags.GetString(ConfigFile)
	stdin, _ := flags.GetBool(StdIn)

	if configFile == "" && !stdin {
		return errors.New("You need to provide configuration from either " + ConfigFile + " or " + StdIn)
	}

	if stdin {
		requestConfig, err = lbot.ParseStdInConfig()
		fmt.Printf("%+v", requestConfig)
	}

	if configFile != "" {
		requestConfig, err = lbot.ParseConfigFile(configFile)
		if err != nil {
			return err
		}
	}

	// to change
	var reply int

	Logger.Info("🚀 Setting new config")

	client, err := rpc.DialHTTP("tcp", "0.0.0.0:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}

	err = client.Call("SetConfigProcess.Run", requestConfig, &reply)
	if err != nil {
		return fmt.Errorf("Setting config failed: %w", err)
	}

	Logger.Info("✅ Setting config succeeded")

	return
}

// todo: command for setting only one field
// ex. --set=cos.tam.tam=2
