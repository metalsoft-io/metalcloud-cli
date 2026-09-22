package cmd

import (
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/ai"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	aiFlags = struct {
		configSource string
		datacenter   string
	}{}

	aiCmd = &cobra.Command{
		Use:     "ai [command]",
		Aliases: []string{"eli"},
		Short:   "MetalSoft AI assistant",
		Long: `Ask the MetalSoft AI assistant (Eli) about the platform.

Command categories:
  Ask:   generate`,
	}

	aiGenerateCmd = &cobra.Command{
		Use:     "generate [prompt...]",
		Aliases: []string{"ask", "prompt"},
		Short:   "Ask the AI assistant a question",
		Long: `Send a prompt to the MetalSoft AI assistant and print its answer.

The prompt is given as positional text or read from a configuration document with
--config-source. The API requires the question to be scoped to a datacenter.

Required Arguments (one of):
  prompt...                The question to ask. All positional arguments are joined with spaces.
  --config-source string   Source of the request document ('pipe' or path to a JSON/YAML file)
                           with a "prompt" and optionally a "datacenter" field.

Required Flags:
  --datacenter string      The datacenter the question is about. May instead be
                           supplied by the --config-source document.

Text output prints the answer as prose; use -f json or -f yaml for the full
response object.

Examples:
  metalcloud ai generate --datacenter dc1 "which servers are available?"
  metalcloud ai generate --datacenter dc1 --config-source prompt.json
  echo '{"datacenter":"dc1","prompt":"list my infrastructures"}' | metalcloud ai generate --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_AI_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			if aiFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(aiFlags.configSource)
				if err != nil {
					return err
				}

				return ai.AIGenerateFromConfig(cmd.Context(), aiFlags.datacenter, config)
			}

			return ai.AIGenerate(cmd.Context(), aiFlags.datacenter, strings.Join(args, " "))
		},
	}
)

func init() {
	rootCmd.AddCommand(aiCmd)

	aiCmd.AddCommand(aiGenerateCmd)
	aiGenerateCmd.Flags().StringVar(&aiFlags.configSource, "config-source", "", "Source of the AI request. Can be 'pipe' or path to a JSON/YAML file.")
	aiGenerateCmd.Flags().StringVar(&aiFlags.datacenter, "datacenter", "", "The datacenter the question is about.")
}
