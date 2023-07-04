package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"

	"github.com/metalsoft-io/metalcloud-cli/pkg/apply"
	"github.com/metalsoft-io/metalcloud-cli/pkg/datacenter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/drive"
	"github.com/metalsoft-io/metalcloud-cli/pkg/firewall"
	"github.com/metalsoft-io/metalcloud-cli/pkg/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/jobs"
	"github.com/metalsoft-io/metalcloud-cli/pkg/network"
	"github.com/metalsoft-io/metalcloud-cli/pkg/osasset"
	"github.com/metalsoft-io/metalcloud-cli/pkg/ostemplate"
	"github.com/metalsoft-io/metalcloud-cli/pkg/reports"
	"github.com/metalsoft-io/metalcloud-cli/pkg/secret"
	"github.com/metalsoft-io/metalcloud-cli/pkg/server"
	"github.com/metalsoft-io/metalcloud-cli/pkg/shellcompletion"
	"github.com/metalsoft-io/metalcloud-cli/pkg/stagedefinition"
	"github.com/metalsoft-io/metalcloud-cli/pkg/storage"
	"github.com/metalsoft-io/metalcloud-cli/pkg/subnetpool"
	"github.com/metalsoft-io/metalcloud-cli/pkg/switchcontroller"
	"github.com/metalsoft-io/metalcloud-cli/pkg/switchdevice"
	"github.com/metalsoft-io/metalcloud-cli/pkg/user"
	"github.com/metalsoft-io/metalcloud-cli/pkg/variable"
	"github.com/metalsoft-io/metalcloud-cli/pkg/version"
	"github.com/metalsoft-io/metalcloud-cli/pkg/volumetemplate"
	"github.com/metalsoft-io/metalcloud-cli/pkg/workflows"
)

func initClients() (map[string]metalcloud.MetalCloudClient, error) {

	clients := map[string]metalcloud.MetalCloudClient{}
	endpointSuffixes := map[string]string{
		configuration.DeveloperEndpoint: "/api/developer/developer",
		configuration.ExtendedEndpoint:  "/api/extended",
		configuration.UserEndpoint:      "/api",
		"":                              "/api",
	}

	for clientName, suffix := range endpointSuffixes {

		if (clientName == configuration.DeveloperEndpoint || clientName == configuration.ExtendedEndpoint) && !configuration.IsAdmin() {
			continue
		}

		client, err := initClient(suffix)
		if err != nil {
			return nil, err
		}
		clients[clientName] = client
	}
	return clients, nil
}

func initClient(endpointSuffix string) (metalcloud.MetalCloudClient, error) {
	if v := os.Getenv("METALCLOUD_USER_EMAIL"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_USER_EMAIL must be set")
	}

	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	apiKey := os.Getenv("METALCLOUD_API_KEY")
	user := os.Getenv("METALCLOUD_USER_EMAIL")

	endpointHost := strings.TrimRight(os.Getenv("METALCLOUD_ENDPOINT"), "/")
	endpoint := fmt.Sprintf("%s%s", endpointHost, endpointSuffix)

	loggingEnabled := isLoggingEnabled()

	err := validateAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return metalcloud.GetMetalcloudClient(user, apiKey, endpoint, loggingEnabled, "", "", "")
}

func getHelp(clients map[string]metalcloud.MetalCloudClient, showArguments bool) string {
	var sb strings.Builder
	cmds := getCommands(clients)
	for _, c := range cmds {
		c.InitFunc(&c)
	}
	sb.WriteString(fmt.Sprintf("Syntax: %s <command> [args]\nAccepted commands:\n", os.Args[0]))
	for _, c := range cmds {
		sb.WriteString(fmt.Sprintln(command.GetCommandHelp(c, false)))
	}
	return sb.String()
}

func isLoggingEnabled() bool {
	return os.Getenv("METALCLOUD_LOGGING_ENABLED") == "true"
}

func validateAPIKey(apiKey string) error {
	const pattern = "^\\d+\\:[0-9a-zA-Z]*$"

	matched, _ := regexp.MatchString(pattern, apiKey)

	if !matched {
		return fmt.Errorf("API Key is not valid. It should start with a number followed by a semicolon followed by alphanumeric characters <id>:<chars> ")
	}

	return nil
}

func getCommands(clients map[string]metalcloud.MetalCloudClient) []command.Command {
	commands := [][]command.Command{
		apply.ApplyCmds,
		datacenter.DatacenterCmds,
		drive.DriveArrayCmds,
		drive.DriveSnapshotCmds,
		drive.SharedDriveCmds,
		firewall.FirewallRuleCmds,
		infrastructure.InfrastructureCmds,
		instance.InstanceArrayCmds,
		instance.InstanceCmds,
		jobs.JobsCmds,
		network.NetworkProfileCmds,
		network.NetworkCmds,
		osasset.OsAssetsCmds,
		ostemplate.OsTemplatesCmds,
		reports.ReportsCmds,
		secret.SecretsCmds,
		server.ServersCmds,
		shellcompletion.ShellCompletionCmds,
		stagedefinition.StageDefinitionsCmds,
		storage.StorageCmds,
		subnetpool.SubnetPoolCmds,
		switchcontroller.SwitchControllerCmds,
		switchdevice.SwitchPairCmds,
		switchdevice.SwitchCmds,
		user.UserCmds,
		variable.VariablesCmds,
		version.VersionCmds,
		volumetemplate.VolumeTemplateCmds,
		workflows.WorkflowCmds,
	}

	filteredCommands := []command.Command{}
	for _, commandSet := range commands {
		commands := fitlerCommandSet(commandSet, clients)
		filteredCommands = append(filteredCommands, commands...)
	}

	return filteredCommands
}

// fitlerCommandSet Filters commands based on endpoint availability for client
func fitlerCommandSet(commandSet []command.Command, clients map[string]metalcloud.MetalCloudClient) []command.Command {
	filteredCommands := []command.Command{}
	for _, command := range commandSet {
		if endpointAvailableForCommand(command, clients) && commandVisibleForUser(command) {
			filteredCommands = append(filteredCommands, command)
		}
	}
	return filteredCommands
}

// endpointAvailableForCommand Checks if the instantiated endpoint clients include the one needed for the command
func endpointAvailableForCommand(command command.Command, clients map[string]metalcloud.MetalCloudClient) bool {
	if configuration.IsAdmin() {
		return clients[command.AdminEndpoint] != nil
	}
	return clients[command.Endpoint] != nil
}

// commandVisibleForUser returns true if the current user (which could be admin or not) has the ability to see the respective command
func commandVisibleForUser(command command.Command) bool {
	if command.UserOnly && configuration.IsAdmin() {
		return false
	}

	if command.AdminOnly && !configuration.IsAdmin() {
		return false
	}

	return true
}

func sameCommand(a *command.Command, b *command.Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}