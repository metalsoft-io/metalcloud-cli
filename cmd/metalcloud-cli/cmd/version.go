package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

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
		RunE:         runVersionCmd,
	}
)

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("CLI Version: %s\n", Version)

	minVersion, maxVersion := system.GetMinMaxVersion()
	fmt.Printf("Minimum Metalsoft Version: %s\n", minVersion)
	fmt.Printf("Maximum Metalsoft Version: %s\n", maxVersion)

	fmt.Printf("Metalsoft Endpoint: %s\n", viper.GetString(system.ConfigEndpoint))
	if cmd.Context() != nil {
		if apiClient, err := api.GetApiClientE(cmd.Context()); err == nil && apiClient != nil {
			if ver, _, err := apiClient.SystemAPI.GetVersion(cmd.Context()).Execute(); err == nil {
				fmt.Printf("Metalsoft Endpoint Version: %s\n", ver.Version)
			}
		}
	}
	if viper.GetBool(system.ConfigInsecure) {
		fmt.Printf("Insecure Mode: %t\n", viper.GetBool(system.ConfigInsecure))
	}

	fmt.Printf("Log File: %s\n", viper.GetString(logger.ConfigLogFile))
	fmt.Printf("Log Verbosity: %s\n", viper.GetString(logger.ConfigVerbosity))

	fmt.Printf("Debug Mode: %t\n", viper.GetBool(system.ConfigDebug))

	if cmd.Context() != nil {
		fmt.Println(currentUserLine(cmd.Context()))
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
}

// currentUserLine describes the user the CLI runs as, read with the credential any other command would
// use: the API token, the service account's token file or the saved session. When the user can't be
// read, the line says why, so a refused or missing credential doesn't pass for a working one.
func currentUserLine(ctx context.Context) string {
	endpoint := viper.GetString(system.ConfigEndpoint)
	if endpoint == "" {
		return "User: cannot be read: no API endpoint is set"
	}

	tokenFile := viper.GetString(tokenFileKey)
	credential, kind, err := resolveCredential(ctx, viper.GetString(system.ConfigApiKey), tokenFile, loginClient(), time.Now())
	if err != nil {
		return "User: cannot be read: " + err.Error()
	}

	ctx = api.SetApiClient(ctx, endpoint, credential, viper.GetBool(system.ConfigDebug), viper.GetBool(system.ConfigInsecure))
	user, err := system.GetCurrentUser(ctx)
	if err != nil {
		return "User: cannot be read: " + sessionEndedOr(err, kind, endpoint, tokenFile).Error()
	}

	email := ""
	if user.Email != "" {
		email = fmt.Sprintf(", Email: %s", user.Email)
	}
	accessLevel := "Unknown"
	if user.AccessLevel != "" {
		accessLevel = user.AccessLevel
	}
	return fmt.Sprintf("User ID: %d%s, Access Level: %s", int(user.Id), email, accessLevel)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
