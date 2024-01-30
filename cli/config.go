package cli

// 1. without args, just prints configfile

// 2. with --set= update config

// 3. generate samle config, for mongodb, postgres itp

import (
	"fmt"
	"net/rpc"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/spf13/cobra"
)

// checks if process is running in local system

// command for setting full new config
func setConfigDriverHandler(cmd *cobra.Command, args []string) (err error) {
	// flags := cmd.Flags()

	flags := cmd.Flags()

	configFile, _ := flags.GetString(ConfigFile)
	request, err := lbot.ParseConfigFile(configFile)
	if err != nil {
		return err
	}

	// to change
	var reply int

	Logger.Info("🚀 Setting new config")

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}

	err = client.Call("SetConfigProcess.Run", request, &reply)
	if err != nil {
		return fmt.Errorf("Setting config failed: %w", err)
	}

	Logger.Info("✅ Setting config succeeded")

	return
}

// todo: command for setting only one field
// ex. --set=cos.tam.tam=2
