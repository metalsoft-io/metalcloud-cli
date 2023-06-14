package version

import (
	"flag"
	"fmt"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
)

var VersionCmds = []command.Command{
	{
		Description:  "Show version.",
		Subject:      "version",
		AltSubject:   "version",
		Predicate:    "show",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("list variables", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: versionShowCmd,
		Endpoint:    configuration.UserEndpoint,
	},
}

func versionShowCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	if configuration.Version == "" {
		return fmt.Sprintf("manual build\n"), nil
	}
	return fmt.Sprintf("Version %s, build %s, date %s\n", configuration.Version, configuration.Commit, configuration.Date), nil
}
