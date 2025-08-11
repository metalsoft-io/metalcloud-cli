package cmd

import (
	"fmt"
	"os"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string

	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Display CLI version, configuration details, and environment information",
		Long: `Display comprehensive version and configuration information for the Metalcloud CLI.

This command shows:
- CLI version information and compatible Metalsoft version range
- Current configuration settings (endpoint, security mode, logging)
- User authentication details (when authenticated)
- Relevant environment variables (proxy settings)

The command uses the current configuration from:
- Configuration file (if present)
- Environment variables
- Command-line flags from previous commands

No additional flags are supported by this command. All information is gathered
from the current CLI configuration and environment.

Examples:
  # Display basic version information
  metalcloud-cli version

  # Use short alias
  metalcloud-cli ver

  # View version info with different verbosity (set globally)
  metalcloud-cli --verbosity debug version

  # Check version with specific endpoint configuration
  metalcloud-cli --endpoint https://my.metalcloud.com version`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("CLI Version: %s\n", Version)

			minVersion, maxVersion := system.GetMinMaxVersion()
			fmt.Printf("Minimum Metalsoft Version: %s\n", minVersion)
			fmt.Printf("Maximum Metalsoft Version: %s\n", maxVersion)

			fmt.Printf("Metalsoft Endpoint: %s\n", viper.GetString(system.ConfigEndpoint))
			if viper.GetBool(system.ConfigInsecure) {
				fmt.Printf("Insecure Mode: %t\n", viper.GetBool(system.ConfigInsecure))
			}

			fmt.Printf("Log File: %s\n", viper.GetString(logger.ConfigLogFile))
			fmt.Printf("Log Verbosity: %s\n", viper.GetString(logger.ConfigVerbosity))

			fmt.Printf("Debug Mode: %t\n", viper.GetBool(system.ConfigDebug))

			// Try to get user access level if API client is available
			if cmd.Context() != nil && api.GetApiClient(cmd.Context()) != nil {
				if user, _, err := api.GetApiClient(cmd.Context()).AuthenticationAPI.GetCurrentUser(cmd.Context()).Execute(); err == nil {
					email := ""
					if user.Email != "" {
						email = fmt.Sprintf(", Email: %s", user.Email)
					}

					accessLevel := "Unknown"
					if user.AccessLevel != "" {
						accessLevel = user.AccessLevel
					}

					fmt.Printf("User ID: %d%s, Access Level: %s\n", int(user.Id), email, accessLevel)
				}
			}

			// Environment variables
			fmt.Printf("\nEnvironment Variables:\n")
			envVars := []string{
				"HTTP_PROXY",
				"HTTPS_PROXY",
				"NO_PROXY",
				"http_proxy",
				"https_proxy",
				"no_proxy",
			}

			for _, envVar := range envVars {
				value := os.Getenv(envVar)
				if value != "" {
					fmt.Printf("%s=%s\n", envVar, value)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
