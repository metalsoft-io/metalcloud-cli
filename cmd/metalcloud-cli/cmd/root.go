package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "metalcloud-cli",
	Short: "MetalSoft MetalCloud CLI",
	Long: `A CLI for interacting with MetalSoft MetalCloud instance.

This CLI requires the correct version of the CLI to be used with the MetalSoft MetalCloud instance.`,
	PersistentPreRunE:  rootPersistentPreRun,
	PersistentPostRunE: rootPersistentPostRun,
}

func init() {
	// Add the root command flags
	rootCmd.PersistentFlags().StringP(system.ConfigEndpoint, "e", "", "MetalCloud API Endpoint")
	rootCmd.PersistentFlags().StringP(system.ConfigApiKey, "k", "", "MetalCloud API Key")

	rootCmd.PersistentFlags().StringP(logger.ConfigVerbosity, "v", "INFO", "Set the log level verbosity")
	rootCmd.PersistentFlags().StringP(logger.ConfigLogFile, "l", "", "Set the log file path")

	rootCmd.PersistentFlags().StringP(formatter.ConfigFormat, "f", "text", "The output format. Supported values are 'text','csv','md','json','yaml'.")

	rootCmd.PersistentFlags().BoolP(system.ConfigDebug, "d", false, "Set to true to enable debug logging")

	// Bind the flags to the viper configuration
	viper.SetConfigName(system.ConfigName)
	viper.SetConfigType(system.ConfigType)
	viper.AddConfigPath(system.ConfigPath1)
	viper.AddConfigPath(system.ConfigPath2)
	viper.AddConfigPath(system.ConfigPath3)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "failed to parse the config file: %v", err)
			os.Exit(-1)
		}
	}

	viper.SetEnvPrefix(system.ConfigPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bind flags for root: %v", err)
		os.Exit(-1)
	}
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) error {
	err := logger.Init()
	if err != nil {
		return err
	}

	// Initialize API client using the arguments from the command line or environment variables
	cfg := sdk.NewConfiguration()
	cfg.UserAgent = "metalcloud-cli"
	cfg.Servers = []sdk.ServerConfiguration{
		{
			URL:         viper.GetString(system.ConfigEndpoint),
			Description: "MetalSoft",
		},
	}

	// Set debug mode
	cfg.Debug = viper.GetBool(system.ConfigDebug)

	// Create API client
	apiClient := sdk.NewAPIClient(cfg)

	ctx := cmd.Context()
	ctx = context.WithValue(ctx, system.ApiClientContextKey, apiClient)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, viper.GetString(system.ConfigApiKey))

	// Validate the version of the CLI
	err = system.ValidateVersion(ctx)
	if err != nil {
		return err
	}

	userPermissions, err := system.GetUserPermissions(ctx)
	if err != nil {
		return err
	}

	// TODO: At this point the help function is already processed and hiding the commands will not work
	hideUnavailableCommands(cmd, userPermissions)

	cmd.SetContext(ctx)

	return nil
}

func rootPersistentPostRun(cmd *cobra.Command, args []string) error {
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func hideUnavailableCommands(cmd *cobra.Command, userPermissionKeys []string) {
	for _, c := range cmd.Commands() {
		if requiredPermission, ok := c.Annotations[system.REQUIRED_PERMISSION]; ok {
			if !slices.Contains(userPermissionKeys, requiredPermission) {
				c.Hidden = true
				c.SilenceUsage = true
			}
		}

		if c.HasSubCommands() {
			hideUnavailableCommands(c, userPermissionKeys)
		}
	}
}
