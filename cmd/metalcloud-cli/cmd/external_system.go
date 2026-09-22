package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/external_system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	externalSystemFlags = struct {
		configSource string
		label        string
		name         string
		annotations  string
		filterLabel  []string
	}{}

	externalSystemCmd = &cobra.Command{
		Use:     "external-system [command]",
		Aliases: []string{"ext-system", "external-systems"},
		Short:   "External system management",
		Long: `Manage the external systems registered with the platform.

An external system is a third-party system the platform integrates with; it is
identified by a unique label and can carry free-form annotations.

Command categories:
  Read:    list, get
  Write:   create, update, delete
  Helper:  config-example`,
	}

	externalSystemListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List external systems",
		Long: `List all external systems registered with the platform.

Optional Flags:
  --filter-label strings   Filter by label. Repeatable or comma-separated.

Examples:
  metalcloud external-system list
  metalcloud ext-system ls --filter-label my-external-system`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_system.ExternalSystemList(cmd.Context(), externalSystemFlags.filterLabel)
		},
	}

	externalSystemGetCmd = &cobra.Command{
		Use:     "get external_system_id",
		Aliases: []string{"show"},
		Short:   "Get external system details",
		Long: `Get the details of an external system, including its annotations.

Required Arguments:
  external_system_id   The numeric ID of the external system

Examples:
  metalcloud external-system get 12
  metalcloud ext-system show 12 -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_system.ExternalSystemGet(cmd.Context(), args[0])
		},
	}

	externalSystemConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print an example external system configuration",
		Long: `Print an example configuration that can be edited and passed to
'external-system create --config-source'.

Examples:
  metalcloud external-system config-example > external-system.json
  metalcloud ext-system config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_system.ExternalSystemConfigExample(cmd.Context())
		},
	}

	externalSystemCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create an external system",
		Long: `Create a new external system.

The external system can be described either with a configuration document
(--config-source) or with individual flags (--label, --name, --annotations).

Required Flags (one of):
  --config-source string   Source of the new external system configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Unique label of the new external system (lowercase letters, digits and dashes)

Optional Flags (when not using --config-source):
  --name string            Display name. Defaults to the label.
  --annotations string     Annotations as a JSON object, e.g. '{"owner":"platform-team"}'

Examples:
  metalcloud external-system create --config-source external-system.json
  cat external-system.yaml | metalcloud ext-system create --config-source pipe
  metalcloud external-system create --label my-system --name "My System" --annotations '{"owner":"platform-team"}'`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateExternalSystem

			if externalSystemFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(externalSystemFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.CreateExternalSystem{
					Label: externalSystemFlags.label,
					Name:  externalSystemFlags.label,
				}
				if externalSystemFlags.name != "" {
					create.Name = externalSystemFlags.name
				}
				if externalSystemFlags.annotations != "" {
					annotations := map[string]interface{}{}
					if err := json.Unmarshal([]byte(externalSystemFlags.annotations), &annotations); err != nil {
						return fmt.Errorf("invalid --annotations, expected a JSON object: %w", err)
					}
					create.Annotations = annotations
				}
			}

			return external_system.ExternalSystemCreate(cmd.Context(), create)
		},
	}

	externalSystemUpdateCmd = &cobra.Command{
		Use:     "update external_system_id",
		Aliases: []string{"edit"},
		Short:   "Update an external system",
		Long: `Update the label, name or annotations of an external system. Fields absent from
the configuration document are left unchanged.

The current revision is read first and sent as the If-Match entity tag, so the
update fails if somebody else changed the external system in the meantime.

Required Arguments:
  external_system_id       The numeric ID of the external system

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud external-system update 12 --config-source update.json
  echo '{"name":"New name"}' | metalcloud ext-system update 12 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(externalSystemFlags.configSource)
			if err != nil {
				return err
			}

			return external_system.ExternalSystemUpdate(cmd.Context(), args[0], config)
		},
	}

	externalSystemDeleteCmd = &cobra.Command{
		Use:     "delete external_system_id",
		Aliases: []string{"rm"},
		Short:   "Delete an external system",
		Long: `Delete an external system.

Required Arguments:
  external_system_id   The numeric ID of the external system

Examples:
  metalcloud external-system delete 12
  metalcloud ext-system rm 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_SYSTEMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_system.ExternalSystemDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(externalSystemCmd)

	externalSystemCmd.AddCommand(externalSystemListCmd)
	externalSystemListCmd.Flags().StringSliceVar(&externalSystemFlags.filterLabel, "filter-label", nil, "Filter by external system label.")

	externalSystemCmd.AddCommand(externalSystemGetCmd)

	externalSystemCmd.AddCommand(externalSystemConfigExampleCmd)

	externalSystemCmd.AddCommand(externalSystemCreateCmd)
	externalSystemCreateCmd.Flags().StringVar(&externalSystemFlags.configSource, "config-source", "", "Source of the new external system configuration. Can be 'pipe' or path to a JSON/YAML file.")
	externalSystemCreateCmd.Flags().StringVar(&externalSystemFlags.label, "label", "", "Unique label of the new external system.")
	externalSystemCreateCmd.Flags().StringVar(&externalSystemFlags.name, "name", "", "Display name of the new external system.")
	externalSystemCreateCmd.Flags().StringVar(&externalSystemFlags.annotations, "annotations", "", "Annotations of the new external system, as a JSON object.")
	externalSystemCreateCmd.MarkFlagsOneRequired("config-source", "label")
	externalSystemCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	externalSystemCreateCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	externalSystemCreateCmd.MarkFlagsMutuallyExclusive("config-source", "annotations")

	externalSystemCmd.AddCommand(externalSystemUpdateCmd)
	externalSystemUpdateCmd.Flags().StringVar(&externalSystemFlags.configSource, "config-source", "", "Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.")
	externalSystemUpdateCmd.MarkFlagRequired("config-source")

	externalSystemCmd.AddCommand(externalSystemDeleteCmd)
}
