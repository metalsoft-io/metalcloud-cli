package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"

	"github.com/metalsoft-io/metalcloud-cli/pkg/apply"
	"github.com/metalsoft-io/metalcloud-cli/pkg/custom_isos"
	"github.com/metalsoft-io/metalcloud-cli/pkg/datacenter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/drive"
	"github.com/metalsoft-io/metalcloud-cli/pkg/extension"
	"github.com/metalsoft-io/metalcloud-cli/pkg/firewall"
	"github.com/metalsoft-io/metalcloud-cli/pkg/firmware"
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
	"github.com/metalsoft-io/metalcloud-cli/pkg/subnetoob"
	"github.com/metalsoft-io/metalcloud-cli/pkg/subnetpool"
	"github.com/metalsoft-io/metalcloud-cli/pkg/switchcontroller"
	"github.com/metalsoft-io/metalcloud-cli/pkg/switchdevice"
	"github.com/metalsoft-io/metalcloud-cli/pkg/user"
	"github.com/metalsoft-io/metalcloud-cli/pkg/variable"
	"github.com/metalsoft-io/metalcloud-cli/pkg/version"
	"github.com/metalsoft-io/metalcloud-cli/pkg/vm"
	"github.com/metalsoft-io/metalcloud-cli/pkg/volumetemplate"
	"github.com/metalsoft-io/metalcloud-cli/pkg/workflows"
)

func initClients() (map[string]metalcloud.MetalCloudClient, *metalcloud2.APIClient, error) {
	clients := map[string]metalcloud.MetalCloudClient{}
	endpointSuffixes := map[string]string{
		configuration.DeveloperEndpoint: "/api/developer/developer",
		configuration.ExtendedEndpoint:  "/api/extended",
		configuration.UserEndpoint:      "/api",
		"":                              "/api",
	}

	for clientName, suffix := range endpointSuffixes {
		client, err := initClient(suffix)
		if err != nil {
			return nil, nil, err
		}

		clients[clientName] = client
	}

	config, err := clientConfiguration()
	if err != nil {
		return nil, nil, err
	}

	client2 := metalcloud2.NewAPIClient(config)

	return clients, client2, nil
}

func getUser(userId int, client metalcloud.MetalCloudClient) (*metalcloud.User, error) {
	return client.UserGet(userId)
}

func getUserPermissions(userId int, client metalcloud.MetalCloudClient) ([]string, error) {

	user, err := getUser(userId, client)
	if err != nil {
		return []string{}, err
	}

	return (*user).UserPermissions, nil
}

func getUserIdFromAPIKey(apiKey string) (int, error) {
	components := strings.Split(apiKey, ":")
	if len(components) != 2 {
		return 0, fmt.Errorf("The API key is not in the correct format")
	}
	return strconv.Atoi(components[0])
}

func initClient(endpointSuffix string) (metalcloud.MetalCloudClient, error) {

	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	apiKey := os.Getenv("METALCLOUD_API_KEY")
	err := validateAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	userId, err := getUserIdFromAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	endpointHost := strings.TrimRight(os.Getenv("METALCLOUD_ENDPOINT"), "/")
	endpoint := fmt.Sprintf("%s%s", endpointHost, endpointSuffix)

	insecureSkipVerify := false
	if strings.ToLower(os.Getenv("METALCLOUD_INSECURE_SKIP_VERIFY")) == "true" {
		insecureSkipVerify = true
	}

	timeout := 5 * time.Minute
	if os.Getenv("METALCLOUD_TIMEOUT_SECONDS") != "" {
		timeout_seconds, err := strconv.Atoi(os.Getenv("METALCLOUD_TIMEOUT_SECONDS"))
		if err != nil {
			return nil, fmt.Errorf("cannot parse timeout, use seconds")
		}
		timeout = time.Second * time.Duration(timeout_seconds)
	}

	options := metalcloud.ClientOptions{
		ApiKey:               apiKey,
		Endpoint:             endpoint,
		LoggingEnabled:       isLoggingEnabled(),
		InsecureSkipVerify:   insecureSkipVerify,
		User:                 "",
		UserID:               userId,
		AuthenticationMethod: metalcloud.AuthMethodBearer,
		Timeout:              timeout,
	}

	return metalcloud.GetMetalcloudClientWithOptions(options)
}

func clientConfiguration() (*metalcloud2.Configuration, error) {
	basePath := os.Getenv("METALCLOUD_ENDPOINT")
	if basePath == "" {
		return nil, fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	apiKey := os.Getenv("METALCLOUD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	config := metalcloud2.NewConfiguration()

	config.BasePath = basePath
	config.UserAgent = "MetalCloud CLI"
	config.AddDefaultHeader("Content-Type", "application/json")
	config.AddDefaultHeader("Accept", "application/json")
	config.AddDefaultHeader("Authorization", "Bearer "+apiKey)

	return config, nil
}

func getHelp(clients map[string]metalcloud.MetalCloudClient, permissions []string) string {
	var sb strings.Builder
	cmds := getCommands(clients, permissions)
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

func getCommands(clients map[string]metalcloud.MetalCloudClient, permissions []string) []command.Command {
	commands := [][]command.Command{
		apply.ApplyCmds,
		custom_isos.CustomISOCmds,
		datacenter.DatacenterCmds,
		drive.DriveArrayCmds,
		drive.DriveSnapshotCmds,
		drive.SharedDriveCmds,
		extension.ExtensionCmds,
		extension.ExtensionInstanceCmds,
		firewall.FirewallRuleCmds,
		firmware.FirmwareCatalogCmds,
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
		subnetoob.SubnetOOBCmds,
		switchcontroller.SwitchControllerCmds,
		switchdevice.SwitchCmds,
		switchdevice.SwitchDefaultsCmds,
		switchdevice.SwitchPairCmds,
		user.UserCmds,
		variable.VariablesCmds,
		version.VersionCmds,
		vm.VmInstanceGroupsCmds,
		vm.VmInstancesCmds,
		vm.VmPoolsCmds,
		vm.VmTypesCmds,
		volumetemplate.VolumeTemplateCmds,
		workflows.WorkflowCmds,
	}

	filteredCommands := []command.Command{}
	for _, commandSet := range commands {
		commands := fitlerCommandSet(commandSet, clients, permissions)
		filteredCommands = append(filteredCommands, commands...)
	}

	return filteredCommands
}

// fitlerCommandSet Filters commands based on endpoint availability for client
func fitlerCommandSet(commandSet []command.Command, clients map[string]metalcloud.MetalCloudClient, permissions []string) []command.Command {
	filteredCommands := []command.Command{}

	for _, command := range commandSet {
		if endpointAvailableForCommand(command, clients, permissions) &&
			commandVisibleForUser(command, permissions) ||
			command.ExecuteFunc2 != nil {
			filteredCommands = append(filteredCommands, command)
		}
	}

	return filteredCommands
}

// endpointAvailableForCommand Checks if the instantiated endpoint clients include the one needed for the command
func endpointAvailableForCommand(c command.Command, clients map[string]metalcloud.MetalCloudClient, permissions []string) bool {
	if slices.Contains(permissions, command.ADMIN_ACCESS) {

		return clients[c.AdminEndpoint] != nil
	}
	return clients[c.Endpoint] != nil
}

// commandVisibleForUser returns true if the current user (which could be admin or not) has the ability to see the respective command
func commandVisibleForUser(c command.Command, permissions []string) bool {

	if c.UserOnly && slices.Contains(permissions, command.ADMIN_ACCESS) {
		return false
	}

	if c.AdminOnly && !slices.Contains(permissions, command.ADMIN_ACCESS) {
		return false
	}

	for _, permission := range c.PermissionsRequired {
		if !slices.Contains(permissions, permission) {
			return false
		}
	}

	return true
}

func sameCommand(a *command.Command, b *command.Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}
