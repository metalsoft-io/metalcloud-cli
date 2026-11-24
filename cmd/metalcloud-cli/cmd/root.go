package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/login"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const (
	noCredentialText = "Sign in with `metalcloud-cli login`, or set an API token with --api_key or METALCLOUD_API_KEY."
	sessionEndedText = "Your session has ended; run `metalcloud-cli login`"

	serviceAccountSessionEndedText = "The service account's session ended; the next command signs in again from the token file."

	// tokenFileKey is the setting that holds the service account's token file: the login flag, the
	// config file and METALCLOUD_SERVICE_ACCOUNT_TOKEN_FILE (viper maps the dashes to underscores).
	tokenFileKey = "service-account-token-file"
)

// credentialKind tells what a command authenticates with, so a 401 is handled for the right one.
type credentialKind int

const (
	credentialAPIKey credentialKind = iota
	credentialSession
	credentialServiceAccount
)

var (
	cfgFile     string
	showVersion bool

	rootCmd = &cobra.Command{
		Use:   "metalcloud-cli",
		Short: "MetalCloud CLI",
		Long: `A CLI for interacting with MetalSoft instance.

This CLI requires the correct version of the CLI to be used with the MetalCloud instance of compatible version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				return runVersionCmd(cmd, args)
			}
			return cmd.Help()
		},
		PersistentPreRunE:  rootPersistentPreRun,
		PersistentPostRunE: rootPersistentPostRun,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file path")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "Display version information")

	// Add the global persistent flags
	rootCmd.PersistentFlags().StringP(system.ConfigEndpoint, "e", "", "MetalCloud API endpoint")
	rootCmd.PersistentFlags().StringP(system.ConfigApiKey, "k", "", "MetalCloud API token (msat_...), or a legacy API key (<user id>:<key>)")
	rootCmd.PersistentFlags().StringP(logger.ConfigVerbosity, "v", "INFO", "Log level verbosity")
	rootCmd.PersistentFlags().StringP(logger.ConfigLogFile, "l", "", "Log file path")
	rootCmd.PersistentFlags().StringP(formatter.ConfigFormat, "f", "text", "Output format. Supported values are 'text','csv','md','json','yaml'.")
	rootCmd.PersistentFlags().BoolP(system.ConfigDebug, "d", false, "Set to enable debug logging")
	rootCmd.PersistentFlags().BoolP(system.ConfigInsecure, "i", false, "Set to allow insecure transport")

	// Add hidden flag to enable development mode
	rootCmd.PersistentFlags().BoolVarP(&system.AllowDevelop, "allow_develop", "x", false, "Allow development mode")
	rootCmd.PersistentFlags().MarkHidden("allow_develop")

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

// restoreRunE puts back the RunE that the last service account run wrapped, if that wrapper is still in
// place: it is called when the command's RunE ends and at the start of the next run, so a run that
// never reaches RunE (a failing PreRunE, say) can't leave the wrapper installed.
var restoreRunE func()

func rootPersistentPreRun(cmd *cobra.Command, args []string) error {
	if restoreRunE != nil {
		restoreRunE()
	}

	err := logger.Init()
	if err != nil {
		return err
	}

	// Commands that only work with local resources don't talk to the controller at all, so they run
	// without an endpoint, an API key or any network access
	if cmd.Annotations[system.LOCAL_COMMAND] == "true" {
		return nil
	}

	endpoint := viper.GetString(system.ConfigEndpoint)

	// Commands that don't require endpoint or API key
	skipValidation := cmd.Name() == "version" || showVersion || cmd.Name() == "completion" ||
		cmd.Name() == "__complete" || cmd.Name() == "__completeNoDesc" ||
		(cmd.Parent() != nil && cmd.Parent().Name() == "completion")

	if !skipValidation && endpoint == "" {
		return fmt.Errorf("API endpoint is required. Use --%s flag or set %s_%s environment variable",
			system.ConfigEndpoint,
			strings.ToUpper(system.ConfigPrefix),
			strings.ToUpper(system.ConfigEndpoint))
	}

	// Validate the endpoint is proper URL
	u, err := url.Parse(endpoint)
	if !skipValidation && (err != nil || u.Scheme == "" || u.Host == "") {
		return fmt.Errorf("invalid API endpoint URL: %s", endpoint)
	}

	// login and logout manage the saved sign-in themselves
	if cmd.Annotations[system.NO_API_KEY_COMMAND] == "true" {
		return nil
	}

	// A command uses the API token if there is one, else the saved sign-in for the endpoint
	credential := viper.GetString(system.ConfigApiKey)
	kind := credentialAPIKey
	if !skipValidation {
		credential, kind, err = resolveCredential(cmd.Context(), credential, viper.GetString(tokenFileKey), loginClient(), time.Now())
		if err != nil {
			return err
		}
	}

	// The startup calls: the CLI version check and the current user. The server can end a session
	// before the time saved with it, so a 401 on a service account's saved session is answered by one
	// new sign-in from the token file and one repeat of both calls. A 401 on anything else, or on the
	// repeat, ends the command.
	ctx, userId, accessLevel, userPermissions, err := startupCalls(cmd.Context(), credential, skipValidation)
	if err != nil && kind == credentialServiceAccount && isUnauthorized(err) {
		tokenFile := viper.GetString(tokenFileKey)
		if delErr := login.DeleteServiceAccountSession(endpoint, tokenFile); delErr != nil {
			logger.Get().Warn().Err(delErr).Msg("Could not remove the ended session")
		}
		credential, err = login.ServiceAccountCredential(cmd.Context(), loginClient(), tokenFile, endpoint, time.Now())
		if err != nil {
			return err
		}
		ctx, userId, accessLevel, userPermissions, err = startupCalls(cmd.Context(), credential, skipValidation)
	}
	if err != nil {
		// For version and completion commands, don't fail if we can't get permissions
		if skipValidation {
			cmd.SetContext(ctx)
			return nil
		}
		return sessionEndedOr(err, kind, endpoint, viper.GetString(tokenFileKey))
	}

	ctx = api.SetUserId(ctx, userId)
	ctx = api.SetUserAccessLevel(ctx, accessLevel)

	// The API no longer exposes a per-user permission list, so command visibility
	// can no longer be gated by permission key. When no permissions are returned,
	// leave every command visible; the server still enforces authorization.
	if len(userPermissions) > 0 {
		hideUnavailableCommands(cmd, userPermissions)
	}

	cmd.SetContext(ctx)

	// A 401 after the command has started is not retried: with a service account's session it fails
	// the command with the session-ended text and forgets the session, so the next command signs in
	// again from the token file. The command's own RunE is wrapped for this one run only (see restoreRunE).
	if kind == credentialServiceAccount && cmd.RunE != nil {
		run := cmd.RunE
		restoreRunE = func() {
			cmd.RunE = run
			restoreRunE = nil
		}
		cmd.RunE = func(c *cobra.Command, args []string) error {
			defer restoreRunE()
			err := run(c, args)
			if err == nil {
				return nil
			}
			return sessionEndedOr(err, kind, endpoint, viper.GetString(tokenFileKey))
		}
	}

	return nil
}

// resolveCredential picks what authenticates a command: the API token (flag, environment variable or
// config), else the service account's token file when one is set (a person's saved session is then
// ignored), else the saved session for the endpoint if it has not expired.
func resolveCredential(ctx context.Context, apiKey, tokenFile string, c login.Client, now time.Time) (token string, kind credentialKind, err error) {
	if apiKey != "" {
		return apiKey, credentialAPIKey, nil
	}

	if tokenFile != "" {
		token, err := login.ServiceAccountCredential(ctx, c, tokenFile, c.Endpoint, now)
		if err != nil {
			return "", credentialServiceAccount, err
		}
		return token, credentialServiceAccount, nil
	}

	session, err := login.LoadSession(c.Endpoint)
	if err != nil {
		return "", credentialSession, err
	}
	if session != nil && session.Token != "" {
		expires, parseErr := time.Parse(time.RFC3339, session.ExpiresTimestamp)
		if parseErr == nil && now.Before(expires) {
			return session.Token, credentialSession, nil
		}
	}

	return "", credentialSession, errors.New(noCredentialText)
}

// unauthorizedStatus matches the API's 401 status line ("401 Unauthorized - <body>") at the start of
// an error message, or right after a ": " when a command wrapped it ("failed to X: 401 ..."), so an
// id or a name that merely contains "401 " is not taken for it. The SDK's error type (and the
// inspected response error) only expose the status through their text, so this covers %w and %v
// wrapping alike.
var unauthorizedStatus = regexp.MustCompile(`(^|: )401 `)

// isUnauthorized tells whether err is, or wraps, an HTTP 401 answer from the API.
func isUnauthorized(err error) bool {
	return unauthorizedStatus.MatchString(err.Error())
}

// sessionEndedOr turns a 401 on a saved session into the session-ended message and forgets the session:
// the person's session, or the service account's one for the token file, whichever the command used. Any other error, a
// 403 included, is returned as it is and the session stays.
func sessionEndedOr(err error, kind credentialKind, endpoint, tokenFile string) error {
	if !isUnauthorized(err) {
		return err
	}

	switch kind {
	case credentialSession:
		if delErr := login.DeleteSession(endpoint); delErr != nil {
			logger.Get().Warn().Err(delErr).Msg("Could not remove the ended session")
		}
		return errors.New(sessionEndedText)
	case credentialServiceAccount:
		if delErr := login.DeleteServiceAccountSession(endpoint, tokenFile); delErr != nil {
			logger.Get().Warn().Err(delErr).Msg("Could not remove the ended session")
		}
		return errors.New(serviceAccountSessionEndedText)
	}
	return err
}

// startupCalls creates the API client for the credential and makes the calls every command starts
// with. It returns the context that carries the client, also when a call fails.
func startupCalls(parent context.Context, credential string, skipValidation bool) (ctx context.Context, userId, accessLevel string, userPermissions []string, err error) {
	ctx = api.SetApiClient(parent,
		viper.GetString(system.ConfigEndpoint),
		credential,
		viper.GetBool(system.ConfigDebug),
		viper.GetBool(system.ConfigInsecure),
	)

	// Skip version validation for version and completion commands since they may fail with develop versions
	if !skipValidation {
		// Validate the version of the CLI
		if err = system.ValidateVersion(ctx, Version); err != nil {
			return
		}
	}

	userId, accessLevel, userPermissions, err = system.GetUserPermissions(ctx)
	return
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

func GenerateDocs() error {
	rootCmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		return err
	}
	return nil
}
