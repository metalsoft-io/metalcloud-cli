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
				"METALCLOUD_INSECURE_SKIP_VERIFY",
				"METALCLOUD_ENDPOINT",
				"METALCLOUD_USER_EMAIL",
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
