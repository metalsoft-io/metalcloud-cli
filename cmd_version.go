package main

import (
	"flag"
	"fmt"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go/v2"
)

var versionCmds = []Command{

	{
		Description:  "Show version.",
		Subject:      "version",
		AltSubject:   "version",
		Predicate:    "show",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("list variables", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: versionShowCmd,
		Endpoint:    UserEndpoint,
	},
}

func versionShowCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	if version == "" {
		return fmt.Sprintf("manual build\n"), nil
	}
	return fmt.Sprintf("Version %s, build %s, date %s\n", version, commit, date), nil
}
