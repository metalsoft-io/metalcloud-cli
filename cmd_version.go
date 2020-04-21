package main

import (
	"flag"
	"fmt"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var versionCmds = []Command{

	{
		Description:  "Show version",
		Subject:      "version",
		AltSubject:   "version",
		Predicate:    "show",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("list variables", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Variable's usage"),
			}
		},
		ExecuteFunc: versionShowCmd,
		Endpoint:    UserEndpoint,
	},
}

func versionShowCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	if version == "" {
		return fmt.Sprintf("manual build\n"), nil
	}
	return fmt.Sprintf("Version %s, build %s, date %s\n", version, commit, date), nil
}
