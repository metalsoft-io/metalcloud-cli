package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"gopkg.in/yaml.v3"
)

const yamlSeparator = "\n---"

//infrastructureCmds commands affecting infrastructures
var applyCmds = []Command{

	{
		Description:  "Apply changes from file.",
		Subject:      "apply",
		AltSubject:   "apply",
		Predicate:    _nilDefaultStr,
		AltPredicate: _nilDefaultStr,
		FlagSet:      flag.NewFlagSet("apply", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", _nilDefaultStr, "The file "),
			}
		},
		ExecuteFunc: applyCmd,
		Endpoint:    DeveloperEndpoint,
	},

	{
		Description:  "Delete changes from file.",
		Subject:      "delete",
		AltSubject:   "delete",
		Predicate:    _nilDefaultStr,
		AltPredicate: _nilDefaultStr,
		FlagSet:      flag.NewFlagSet("apply", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", _nilDefaultStr, "The file "),
			}
		},
		ExecuteFunc: deleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func applyCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c, client)

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

func deleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c, client)

	if err != nil {
		return "", err
	}

	for _, object := range objects {
		if err != nil {
			return "", err
		}
		err := object.Delete(client)

		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func readObjectsFromCommand(c *Command, client metalcloud.MetalCloudClient) ([]metalcloud.Applier, error) {
	var err error
	content := []byte{}
	var results []metalcloud.Applier

	if filePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {
		content, err = readInputFromFile(filePath)
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
