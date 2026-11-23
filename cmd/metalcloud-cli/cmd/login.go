package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/login"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	loginFlags = struct {
		noBrowser bool
	}{}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Sign in through the browser",
		Long: `Sign in to MetalCloud as yourself, through the browser, and save the session.

The command prints a URL and a code. Open the URL in a browser on this machine or on any
other device, check that the page shows the same code, and approve the sign-in. The command
then saves the session in $HOME/.metalcloud/sessions.json (one entry per endpoint, readable
only by you). Commands use it when no API token is set with --api_key or METALCLOUD_API_KEY.

The session ends by the same rules as a browser session, for example after a period without
use, and then the next command asks you to sign in again. Some commands, such as the ones
that manage API tokens, are not available to a CLI session; use the UI or an API token for them.

Examples:
  # Sign in, opening the browser on this machine
  metalcloud-cli login

  # Sign in from a machine without a browser; open the printed URL on another device
  metalcloud-cli login --no-browser`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.NO_API_KEY_COMMAND: "true"},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return login.Login(cmd.Context(), loginClient(), clientName(), clientVersion(), loginFlags.noBrowser, os.Stdout)
		},
	}

	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "End the saved session and forget it",
		Long: `End the session saved by 'metalcloud-cli login' on the server and remove it from
$HOME/.metalcloud/sessions.json. The saved session is removed even if the server can't be reached.

Examples:
  metalcloud-cli logout`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.NO_API_KEY_COMMAND: "true"},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(cmd.Context(), loginClient())
		},
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolVar(&loginFlags.noBrowser, "no-browser", false, "Do not open the browser; only print the URL and the code.")

	rootCmd.AddCommand(logoutCmd)
}

func loginClient() login.Client {
	return login.Client{
		Endpoint: viper.GetString(system.ConfigEndpoint),
		Insecure: viper.GetBool(system.ConfigInsecure),
	}
}

// clientName is what the approval page shows as the requester: the machine's name.
func clientName() string {
	name, err := os.Hostname()
	if err != nil || name == "" {
		name = "metalcloud-cli"
	}
	return truncate(name, 64)
}

func clientVersion() string {
	if Version == "" {
		return "unknown"
	}
	return truncate(Version, 32)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

func runLogout(ctx context.Context, c login.Client) error {
	session, err := login.LoadSession(c.Endpoint)
	if err != nil {
		return err
	}
	if session == nil {
		fmt.Printf("No saved session for %s.\n", c.Endpoint)
		return nil
	}

	logoutErr := login.Logout(ctx, c, session.Token)

	if err := login.DeleteSession(c.Endpoint); err != nil {
		return fmt.Errorf("could not remove the saved session: %w", err)
	}

	if logoutErr != nil {
		fmt.Printf("The server could not end the session (%v); the saved session was removed anyway.\n", logoutErr)
		return nil
	}
	fmt.Println("Signed out.")
	return nil
}
