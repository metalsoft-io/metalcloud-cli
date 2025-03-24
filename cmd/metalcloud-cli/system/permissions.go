package system

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

const REQUIRED_PERMISSION = "requiredPermission"

const (
	ADMIN_ACCESS                = "admin_access"
	SITE_READ                   = "site_read"
	SITE_WRITE                  = "site_write"
	FIRMWARE_UPGRADE_READ       = "firmware_upgrade_read"
	FIRMWARE_UPGRADE_WRITE      = "firmware_upgrade_write"
	FIRMWARE_BASELINES_READ     = "firmware_baselines_read"
	FIRMWARE_BASELINES_WRITE    = "firmware_baselines_write"
	JOB_QUEUE_READ              = "job_queue_read"
	JOB_QUEUE_WRITE             = "job_queue_write"
	TEMPLATES_READ              = "templates_read"
	TEMPLATES_WRITE             = "templates_write"
	SERVERS_READ                = "servers_read"
	SERVERS_WRITE               = "servers_write"
	STORAGE_READ                = "storage_read"
	STORAGE_WRITE               = "storage_write"
	SUBNETS_READ                = "subnets_read"
	SUBNETS_WRITE               = "subnets_write"
	SWITCHES_READ               = "switches_read"
	SWITCHES_WRITE              = "switches_write"
	USERS_AND_PERMISSIONS_READ  = "users_and_permissions_read"
	USERS_AND_PERMISSIONS_WRITE = "users_and_permissions_write"
	NETWORK_PROFILES_READ       = "network_profiles_read"
	NETWORK_PROFILES_WRITE      = "network_profiles_write"
	WORKFLOWS_READ              = "workflows_read"
	WORKFLOWS_WRITE             = "workflows_write"
	VM_POOLS_READ               = "vm_pools_read"
	VM_POOLS_WRITE              = "vm_pools_write"
	VM_PROFILES_READ            = "vm_profiles_read"
	VM_PROFILES_WRITE           = "vm_profiles_write"
	VM_TYPES_READ               = "vm_types_read"
	VM_TYPES_WRITE              = "vm_types_write"
	VMS_READ                    = "vms_read"
	VMS_WRITE                   = "vms_write"
	EXTENSIONS_READ             = "extensions_read"
	EXTENSIONS_WRITE            = "extensions_write"
)

func GetUserPermissions(ctx context.Context) (string, []string, error) {
	client := api.GetApiClient(ctx)

	user, httpRes, err := client.AuthenticationAPI.GetCurrentUser(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", nil, err
	}

	// TODO: The API returns the permissions of the user in a map with no specific type
	if user.Permissions == nil || user.Permissions.AdditionalProperties == nil {
		return fmt.Sprintf("%d", int(user.Id)), nil, nil
	}

	userPermissions := make([]string, 0, len(user.Permissions.AdditionalProperties))
	for k := range user.Permissions.AdditionalProperties {
		if user.Permissions.AdditionalProperties[k] == nil {
			continue
		}

		v, ok := user.Permissions.AdditionalProperties[k].(bool)
		if !ok || !v {
			continue
		}

		userPermissions = append(userPermissions, k)
	}

	return fmt.Sprintf("%d", int(user.Id)), userPermissions, nil
}
