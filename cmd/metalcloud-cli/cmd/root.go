package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "metalcloud-cli",
		Short: "MetalSoft MetalCloud CLI",
		Long: `A CLI for interacting with MetalSoft MetalCloud instance.

This CLI requires the correct version of the CLI to be used with the MetalSoft MetalCloud instance.`,
		PersistentPreRunE:  rootPersistentPreRun,
		PersistentPostRunE: rootPersistentPostRun,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// Add the global persistent flags
	rootCmd.PersistentFlags().StringP(system.ConfigEndpoint, "e", "", "MetalCloud API Endpoint")
	rootCmd.PersistentFlags().StringP(system.ConfigApiKey, "k", "", "MetalCloud API Key")
	rootCmd.PersistentFlags().StringP(logger.ConfigVerbosity, "v", "INFO", "Set the log level verbosity")
	rootCmd.PersistentFlags().StringP(logger.ConfigLogFile, "l", "", "Set the log file path")
	rootCmd.PersistentFlags().StringP(formatter.ConfigFormat, "f", "text", "The output format. Supported values are 'text','csv','md','json','yaml'.")
	rootCmd.PersistentFlags().BoolP(system.ConfigDebug, "d", false, "Set to true to enable debug logging")

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bind flags for root: %v", err)
		os.Exit(-1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(system.ConfigName)
		viper.SetConfigType(system.ConfigType)
		viper.AddConfigPath(system.ConfigPath1)
		viper.AddConfigPath(system.ConfigPath2)
		viper.AddConfigPath(system.ConfigPath3)
	}

	viper.SetEnvPrefix(system.ConfigPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) error {
	err := logger.Init()
	if err != nil {
		return err
	}

	// Create API client
	ctx := api.SetApiClient(cmd.Context(), viper.GetString(system.ConfigEndpoint), viper.GetString(system.ConfigApiKey), viper.GetBool(system.ConfigDebug))

	// Validate the version of the CLI
	err = system.ValidateVersion(ctx)
	if err != nil {
		return err
	}

	userId, userPermissions, err := system.GetUserPermissions(ctx)
	if err != nil {
		return err
	}

	ctx = api.SetUserId(ctx, userId)

	// TODO: At this point the help function is already processed and hiding the commands will not work
	hideUnavailableCommands(cmd, userPermissions)

	cmd.SetContext(ctx)

	return nil
}

func rootPersistentPostRun(cmd *cobra.Command, args []string) error {
	return nil
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

func Execute() error {
	return rootCmd.Execute()
}
