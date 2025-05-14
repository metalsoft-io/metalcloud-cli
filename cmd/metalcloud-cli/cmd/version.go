package cmd

import (
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string

	versionCmd = &cobra.Command{
		Use:          "version",
		Aliases:      []string{"ver"},
		Short:        "Get CLI version details.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("CLI Version: %s\n", Version)

			minVersion, maxVersion := system.GetMinMaxVersion()
			fmt.Printf("Minimum Metalsoft Version: %s\n", minVersion)
			fmt.Printf("Maximum Metalsoft Version: %s\n", maxVersion)

			fmt.Printf("Metalsoft Endpoint: %s\n", viper.GetString(system.ConfigEndpoint))

			fmt.Printf("Log File: %s\n", viper.GetString(logger.ConfigLogFile))
			fmt.Printf("Log Verbosity: %s\n", viper.GetString(logger.ConfigVerbosity))

			fmt.Printf("Debug Mode: %t\n", viper.GetBool(system.ConfigDebug))

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
