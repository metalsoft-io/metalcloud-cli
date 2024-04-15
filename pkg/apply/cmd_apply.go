package apply

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"gopkg.in/yaml.v3"
)

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
		FlagSet:      flag.NewFlagSet("apply", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, "The file "),
			}
		},
		ExecuteFunc: deleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Generate a stub of an object",
		Subject:      "generate",
		AltSubject:   "generate",
		Predicate:    command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("apply", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"object": c.FlagSet.String("object", command.NilDefaultStr, "Object to use, if none it will list compatible objects"),
			}
		},
		ExecuteFunc: generateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func applyCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	objects, err := objects.ReadObjectsFromCommand(c, client)

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
	objects, err := objects.ReadObjectsFromCommand(c, client)

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

func generateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	types := metalcloud.GetTypesThatSupportApplierInterface()
	typesArr := []string{}
	for k := range types {
		typesArr = append(typesArr, k)
	}

	if object, ok := command.GetStringParamOk(c.Arguments["object"]); ok {
		t, ok2 := types[object]
		if !ok2 {
			return "", fmt.Errorf("%s was not supported by the apply method. Only the following types are supported: %s", object, strings.Join(typesArr, ","))
		}
		log.Printf("Type is %v", t)
		//obj := metalcloud.SubnetOOB{}
		v := reflect.New(t)
		objects.InitializeStruct(t, v.Elem())

		c := v.Interface().(*metalcloud.SubnetOOB)
		//obj := reflect.Zero(t)
		log.Printf("object is %+v", c)
		b, err := yaml.Marshal(c)
		if err != nil {
			return "", err
		}

		log.Printf(string(b))
		//return string(b), err
		return string(b), nil

	} else {

		return fmt.Sprintf("Types supported by the apply/delete/generate commands:%s", strings.Join(typesArr, ",")), nil
	}
}
