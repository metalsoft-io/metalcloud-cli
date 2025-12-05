package system

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

const REQUIRED_PERMISSION = "requiredPermission"

const (
	PERMISSION_ADMIN_ACCESS                                       = "admin_access"
	PERMISSION_AI_READ                                            = "ai_read"
	PERMISSION_AUTHENTICATION_PROVIDERS_READ                      = "authentication_providers_read"
	PERMISSION_AUTHENTICATION_PROVIDERS_WRITE                     = "authentication_providers_write"
	PERMISSION_SERVERS_READ                                       = "servers_read"
	PERMISSION_SERVERS_WRITE                                      = "servers_write"
	PERMISSION_SERVER_CONSOLE_ACCESS                              = "server_console_access"
	PERMISSION_SERVER_TYPES_READ                                  = "server_types_read"
	PERMISSION_SERVER_TYPES_WRITE                                 = "server_types_write"
	PERMISSION_SERVER_TYPE_UTILIZATION_REPORT_READ                = "server_type_utilization_report_read"
	PERMISSION_SERVER_INSTANCES_READ                              = "server_instances_read"
	PERMISSION_SERVER_INSTANCES_WRITE                             = "server_instances_write"
	PERMISSION_SERVER_INSTANCE_GROUPS_READ                        = "server_instance_groups_read"
	PERMISSION_SERVER_INSTANCE_GROUPS_WRITE                       = "server_instance_groups_write"
	PERMISSION_SWITCHES_READ                                      = "switches_read"
	PERMISSION_SWITCHES_WRITE                                     = "switches_write"
	PERMISSION_STORAGE_READ                                       = "storage_read"
	PERMISSION_STORAGE_WRITE                                      = "storage_write"
	PERMISSION_SUBNETS_READ                                       = "subnets_read"
	PERMISSION_SUBNETS_WRITE                                      = "subnets_write"
	PERMISSION_INFRASTRUCTURES_READ                               = "infrastructures_read"
	PERMISSION_INFRASTRUCTURES_WRITE                              = "infrastructures_write"
	PERMISSION_CUSTOM_ISO_READ                                    = "custom_iso_read"
	PERMISSION_CUSTOM_ISO_WRITE                                   = "custom_iso_write"
	PERMISSION_TEMPLATES_READ                                     = "templates_read"
	PERMISSION_TEMPLATES_WRITE                                    = "templates_write"
	PERMISSION_TEMPLATES_ADMIN                                    = "templates_admin"
	PERMISSION_EVENTS_READ                                        = "events_read"
	PERMISSION_EVENTS_WRITE                                       = "events_write"
	PERMISSION_JOB_QUEUE_READ                                     = "job_queue_read"
	PERMISSION_JOB_QUEUE_WRITE                                    = "job_queue_write"
	PERMISSION_WORKFLOWS_READ                                     = "workflows_read"
	PERMISSION_WORKFLOWS_WRITE                                    = "workflows_write"
	PERMISSION_VARIABLES_AND_SECRETS_READ                         = "variables_and_secrets_read"
	PERMISSION_VARIABLES_AND_SECRETS_WRITE                        = "variables_and_secrets_write"
	PERMISSION_FIRMWARE_UPGRADE_READ                              = "firmware_upgrade_read"
	PERMISSION_FIRMWARE_UPGRADE_WRITE                             = "firmware_upgrade_write"
	PERMISSION_FIRMWARE_BASELINES_READ                            = "firmware_baselines_read"
	PERMISSION_FIRMWARE_BASELINES_WRITE                           = "firmware_baselines_write"
	PERMISSION_ROLES_READ                                         = "roles_read"
	PERMISSION_ROLES_WRITE                                        = "roles_write"
	PERMISSION_ROLES_ASSIGN                                       = "roles_assign"
	PERMISSION_LICENSES_READ                                      = "licenses_read"
	PERMISSION_LICENSES_WRITE                                     = "licenses_write"
	PERMISSION_SUBSCRIPTIONS_READ                                 = "subscriptions_read"
	PERMISSION_SUBSCRIPTIONS_WRITE                                = "subscriptions_write"
	PERMISSION_UTILIZATION_REPORTS_READ                           = "utilization_reports_read"
	PERMISSION_SUSPEND_REASONS_READ                               = "suspend_reasons_read"
	PERMISSION_SUSPEND_REASONS_WRITE                              = "suspend_reasons_write"
	PERMISSION_CLUSTER_READ                                       = "cluster_read"
	PERMISSION_CLUSTER_WRITE                                      = "cluster_write"
	PERMISSION_CONTAINER_PLATFORM_READ                            = "container_platform_read"
	PERMISSION_CONTAINER_PLATFORM_WRITE                           = "container_platform_write"
	PERMISSION_DATASET_READ                                       = "dataset_read"
	PERMISSION_DATASET_WRITE                                      = "dataset_write"
	PERMISSION_DATASTORE_READ                                     = "datastore_read"
	PERMISSION_DATASTORE_WRITE                                    = "datastore_write"
	PERMISSION_AFC_READ                                           = "afc_read"
	PERMISSION_AFC_WRITE                                          = "afc_write"
	PERMISSION_MAINTENANCE_READ                                   = "maintenance_read"
	PERMISSION_MAINTENANCE_WRITE                                  = "maintenance_write"
	PERMISSION_ADMIN_MAINTENANCE_READ                             = "admin_maintenance_read"
	PERMISSION_ADMIN_MAINTENANCE_WRITE                            = "admin_maintenance_write"
	PERMISSION_SKIP_USER_LIMITS                                   = "skip_user_limits"
	PERMISSION_SKIP_AUTHENTICATOR                                 = "skip_authenticator"
	PERMISSION_RESOURCES_WRITE                                    = "resources_write"
	PERMISSION_NETWORK_PROFILES_READ                              = "network_profiles_read"
	PERMISSION_NETWORK_PROFILES_WRITE                             = "network_profiles_write"
	PERMISSION_GLOBAL_CONFIGURATIONS_WRITE                        = "global_configurations_write"
	PERMISSION_GLOBAL_CONFIGURATIONS_READ                         = "global_configurations_read"
	PERMISSION_LICENSE_ADMIN                                      = "license_admin"
	PERMISSION_VM_POOLS_READ                                      = "vm_pools_read"
	PERMISSION_VM_POOLS_WRITE                                     = "vm_pools_write"
	PERMISSION_VMS_READ                                           = "vms_read"
	PERMISSION_VMS_WRITE                                          = "vms_write"
	PERMISSION_VM_INSTANCES_READ                                  = "vm_instances_read"
	PERMISSION_VM_INSTANCES_WRITE                                 = "vm_instances_write"
	PERMISSION_VM_INSTANCE_GROUPS_READ                            = "vm_instance_groups_read"
	PERMISSION_VM_INSTANCE_GROUPS_WRITE                           = "vm_instance_groups_write"
	PERMISSION_VM_TYPES_READ                                      = "vm_types_read"
	PERMISSION_VM_TYPES_WRITE                                     = "vm_types_write"
	PERMISSION_VM_PROFILES_READ                                   = "vm_profiles_read"
	PERMISSION_VM_PROFILES_WRITE                                  = "vm_profiles_write"
	PERMISSION_FILE_SHARE_READ                                    = "file_shares_read"
	PERMISSION_BUCKET_READ                                        = "buckets_read"
	PERMISSION_LOGICAL_NETWORKS_READ                              = "logical_networks_read"
	PERMISSION_LOGICAL_NETWORKS_WRITE                             = "logical_networks_write"
	PERMISSION_FILE_SHARE_WRITE                                   = "file_shares_write"
	PERMISSION_DRIVES_READ                                        = "drives_read"
	PERMISSION_DRIVES_WRITE                                       = "drives_write"
	PERMISSION_BUCKETS_READ                                       = "buckets_read"
	PERMISSION_BUCKETS_WRITE                                      = "buckets_write"
	PERMISSION_RESOURCE_POOLS_READ                                = "resource_pools_read"
	PERMISSION_RESOURCE_POOLS_WRITE                               = "resource_pools_write"
	PERMISSION_RESOURCE_POOL_USER_ACCESS_READ                     = "resource_pool_user_access_read"
	PERMISSION_RESOURCE_POOL_USER_ACCESS_WRITE                    = "resource_pool_user_access_write"
	PERMISSION_USERS_READ                                         = "users_read"
	PERMISSION_USERS_WRITE                                        = "users_write"
	PERMISSION_USERS_AND_PERMISSIONS_WRITE                        = "users_and_permissions_write"
	PERMISSION_EXTENSIONS_READ                                    = "extensions_read"
	PERMISSION_EXTENSIONS_WRITE                                   = "extensions_write"
	PERMISSION_EXTENSION_INSTANCES_READ                           = "extension_instances_read"
	PERMISSION_EXTENSION_INSTANCES_WRITE                          = "extension_instances_write"
	PERMISSION_SITES_READ                                         = "sites_read"
	PERMISSION_SITES_WRITE                                        = "sites_write"
	PERMISSION_SERVER_DEFAULT_CREDENTIALS_READ                    = "server_default_credentials_read"
	PERMISSION_SERVER_DEFAULT_CREDENTIALS_WRITE                   = "server_default_credentials_write"
	PERMISSION_SERVER_CLEANUP_POLICIES_READ                       = "cleanup_policies_read"
	PERMISSION_SERVER_CLEANUP_POLICIES_WRITE                      = "cleanup_policies_write"
	PERMISSION_NETWORK_FABRICS_READ                               = "network_fabrics_read"
	PERMISSION_NETWORK_FABRICS_WRITE                              = "network_fabrics_write"
	PERMISSION_NETWORK_ENDPOINT_GROUPS_READ                       = "network_endpoint_groups_read"
	PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE                      = "network_endpoint_groups_write"
	PERMISSION_EMAIL_TEMPLATES_READ                               = "email_templates_read"
	PERMISSION_EMAIL_TEMPLATES_WRITE                              = "email_templates_write"
	PERMISSION_NETWORK_ENDPOINT_GROUP_WITH_LOGICAL_NETWORKS_READ  = "network_endpoint_group_with_logical_networks_read"
	PERMISSION_NETWORK_ENDPOINT_GROUP_WITH_LOGICAL_NETWORKS_WRITE = "network_endpoint_group_with_logical_networks_write"
	PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ        = "network_device_configuration_templates_read"
	PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE       = "network_device_configuration_templates_write"
	PERMISSION_IMPERSONATE                                        = "impersonate"
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
