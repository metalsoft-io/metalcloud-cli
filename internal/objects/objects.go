package objects

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v2"
)

const YamlSeparator = "\n---"

// ReadSingleObjectFromCommand reads a single object from file, throws error if zero or more than one objects are returned
func ReadSingleObjectFromCommand(c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.Applier, error) {
	objects, err := ReadObjectsFromCommand(c, client)
	if err != nil {
		return nil, err
	}

	if len(objects) != 1 {
		return nil, fmt.Errorf("the file should contain a single object")
	}
	return &objects[0], nil
}

func ReadObjectsFromCommand(c *command.Command, client metalcloud.MetalCloudClient) ([]metalcloud.Applier, error) {
	var err error
	content := []byte{}
	var results []metalcloud.Applier

	if filePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"]); ok {
		content, err = configuration.ReadInputFromFile(filePath)
	} else {
		return nil, fmt.Errorf("-f <input_file_name> is required")
	}

	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("Content cannot be empty")
	}

	objects := strings.Split(string(content), YamlSeparator)

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

func InitializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			InitializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			InitializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

// RenderRawObject wraps tableformatter.RenderRawObject to prepend the kind and version: entry if the format is yaml
func RenderRawObject(obj interface{}, format string, prefixToStripOrObject string) (string, error) {
	ret, err := tableformatter.RenderRawObject(obj, format, prefixToStripOrObject)
	if err != nil {
		return "", err
	}
	switch format {
	case "yaml", "YAML":
		return fmt.Sprintf("kind: %s\napiVersion: 1.0\n%s", prefixToStripOrObject, ret), nil
	default:
		return ret, nil
	}
}
