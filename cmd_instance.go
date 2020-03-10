package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

//instanceCmds commands affecting instances
var instanceCmds = []Command{

	Command{
		Description:  "Control power an instance",
		Subject:      "instance",
		AltSubject:   "instance",
		Predicate:    "power_control",
		AltPredicate: "pwr",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) Instances's id . Note that the 'label' this be ambiguous in certain situations."),
				"operation":   c.FlagSet.String("operation", _nilDefaultStr, "(Required) Power control operation, one of: on, off, reset, soft"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: instancePowerControlCmd,
	},
}

func instancePowerControlCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	instanceID, ok := getIntParamOk(c.Arguments["instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}
	operation, ok := getStringParamOk(c.Arguments["operation"])
	if !ok {
		return "", fmt.Errorf("-operation is required (one of: on, off, reset, soft)")
	}

	instance, err := client.InstanceGet(instanceID)
	if err != nil {
		return "", err
	}

	ia, err := client.InstanceArrayGet(instance.InstanceArrayID)
	if err != nil {
		return "", err
	}

	infra, err := client.InfrastructureGet(ia.InfrastructureID)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		op := ""
		switch operation {
		case "on":
			op = "Turning on"
		case "off":
			op = "Turning off (hard)"
		case "reset":
			op = "Rebooting"
		case "sort":
			op = "Shutting down"
		}

		confirmationMessage := fmt.Sprintf("%s instance %s (%d) of instance array %s (#%d) infrastructure %s (#%d).  Are you sure? Type \"yes\" to continue:",
			op,
			instance.InstanceLabel,
			instance.InstanceID,
			ia.InstanceArrayLabel,
			ia.InstanceArrayID,
			infra.InfrastructureLabel,
			infra.InfrastructureID,
		)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage

	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.InstanceServerPowerSet(instanceID, operation)
	}

	return "", err
}
