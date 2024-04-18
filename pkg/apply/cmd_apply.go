package apply

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"gopkg.in/yaml.v3"
)

const yamlSeparator = "\n---"

var ApplyCmds = []command.Command{
	{
		Description:  "Apply changes from file.",
		Subject:      "apply",
		AltSubject:   "apply",
		Predicate:    command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("apply", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, "The file "),
			}
		},
		ExecuteFunc: applyCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Delete changes from file.",
		Subject:      "delete",
		AltSubject:   "delete",
		Predicate:    command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("delete", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, "The file "),
			}
		},
		ExecuteFunc: deleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func applyCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c)

	if err != nil {
		return "", err
	}

	for _, object := range objects {
		err = object.CreateOrUpdate(client)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func deleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c)
	if err != nil {
		return "", err
	}

	for _, object := range objects {
		err := object.Delete(client)

		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func readObjectsFromCommand(c *command.Command) ([]metalcloud.Applier, error) {
	var err error
	content := []byte{}
	var results []metalcloud.Applier

	if filePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"]); ok {
		content, err = configuration.ReadInputFromFile(filePath)
	} else {
		return nil, fmt.Errorf("file name is required")
	}

	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("Content cannot be empty")
	}

	objects := strings.Split(string(content), yamlSeparator)

	getKind := func(object []byte) (string, error) {
		re := regexp.MustCompile(`kind\s*:\s*(.+)`)
		matches := re.FindAllSubmatch(object, -1)

		if len(matches) > 0 {
			return string(matches[0][1]), nil
		}

		return "", fmt.Errorf("property kind is missing")
	}

	for _, object := range objects {
		if len(strings.Trim(object, " \n\r")) == 0 {
			continue
		}
		bytes := []byte(object)
		kind, err := getKind(bytes)

		if err != nil {
			return nil, err
		}
		kind = strings.Trim(kind, " \n\r")

		newType, err := metalcloud.GetObjectByKind(kind)
		if err != nil {
			return nil, err
		}

		if err = yaml.Unmarshal(bytes, newType.Interface()); err != nil {
			return nil, err
		}
		newObject := newType.Elem().Interface()

		results = append(results, newObject.(metalcloud.Applier))
	}

	return results, nil
}
