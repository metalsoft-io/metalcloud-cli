package main

import (
	"flag"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
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
	},
}

func applyCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c, client)

	if err != nil {
		return "", err
	}

	for _, object := range objects {
		err = object.Validate()
		err = object.CreateOrUpdate(client)
		if err != nil {
			fmt.Printf("errror is %s\n", err)
			return "", err
		}
	}

	return "", nil
}

func deleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	objects, err := readObjectsFromCommand(c, client)

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

func readObjectsFromCommand(c *Command, client interfaces.MetalCloudClient) ([]metalcloud.Applier, error) {
	initTypeRegistry()
	var err error
	content := []byte{}
	var results []metalcloud.Applier

	if filePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {
		content, err = readInputFromFile(filePath)
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

		newType, err := getObjectByKind(kind)
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

var typeRegistry = make(map[string]reflect.Type)

func initTypeRegistry() {
	myTypes := []metalcloud.Applier{
		&metalcloud.InstanceArray{},
		&metalcloud.Datacenter{},
		&metalcloud.DriveArray{},
		&metalcloud.Infrastructure{},
		&metalcloud.Network{},
		&metalcloud.OSAsset{},
		&metalcloud.OSTemplate{},
		&metalcloud.Secret{},
		&metalcloud.Server{},
		&metalcloud.SharedDrive{},
		&metalcloud.StageDefinition{},
		&metalcloud.Workflow{},
		&metalcloud.SubnetPool{},
		&metalcloud.SwitchDevice{},
		&metalcloud.Variable{},
	}

	for _, v := range myTypes {
		t := reflect.ValueOf(v).Elem()
		u := reflect.TypeOf(v).Elem()
		typeRegistry[u.Name()] = t.Type()
	}
}

func getObjectByKind(name string) (reflect.Value, error) {
	t, ok := typeRegistry[name]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%s was not recongnized as a valid product", name)
	}

	v := reflect.New(t)
	return v, nil
}
