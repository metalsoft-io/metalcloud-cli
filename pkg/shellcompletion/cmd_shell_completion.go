package shellcompletion

import (
	"flag"
	"fmt"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
)

var ShellCompletionCmds = []command.Command{
	{
		Description:  "Outputs bash/zsh autocompletion script",
		Subject:      "localshell",
		AltSubject:   "localshell",
		Predicate:    "autocomplete",
		AltPredicate: "gen",
		FlagSet:      flag.NewFlagSet("shell get", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				// "format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: shellCompletionCmd,
		Endpoint:    configuration.UserEndpoint,
	},
}

func shellCompletionCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	// if version == "" {
	// 	return fmt.Sprintf("manual build\n"), nil
	// }
	script := `#/usr/bin/env bash
### This output should be redirected to /etc/bash_completion.d/metalcoud-cli
### or ~/.zshrc . Once done, reload the shell
_metalcloud-cli_completions()
{
  local secondArg firstArg help thirdArg
  firstArg="$(metalcloud-cli help|awk '{print $1}'|sort -u|egrep -v '^Accepted|^Syntax:|^$'|xargs)"
  secondArg="$(metalcloud-cli help|grep "^\s\+${COMP_WORDS[1]}\ "|grep -oP '^\s+[\w\-]+\ [\w\-]+'|awk '{print $2}'|xargs)"
  thirdArg="$(metalcloud-cli ${COMP_WORDS[1]} ${COMP_WORDS[2]} --help|grep -oP '\s+\-[\-\w]+'|awk '{print $1}'|xargs)"
  helpArg="-help --help -h"

  if [ "${#COMP_WORDS[@]}" -eq "4" ]; then
    COMPREPLY=($(compgen -W "$helpArg $thirdArg" -- "${COMP_WORDS[3]}"))
  fi

  if [ "${#COMP_WORDS[@]}" -eq "3" ]; then
    if [ "${COMP_WORDS[1]}" == 'apply' ] || [ "${COMP_WORDS[1]}" == 'delete' ];then
      thirdArg="$(metalcloud-cli ${COMP_WORDS[1]} --help|grep -oP '\s+\-[\-\w]+'|awk '{print $1}'|xargs)"
      COMPREPLY=($(compgen -W "$thirdArg" -- "${COMP_WORDS[2]}"))
    else
      COMPREPLY=($(compgen -W "$secondArg" -- "${COMP_WORDS[2]}"))
    fi
  fi
  if [ "${#COMP_WORDS[@]}" -eq "2" ]; then
    COMPREPLY=($(compgen -W "$firstArg" -- "${COMP_WORDS[1]}"))
  fi
}

complete -F _metalcloud-cli_completions metalcloud-cli
`

	return fmt.Sprintf(script), nil
}
