package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/atomicgo/cursor"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"gopkg.in/yaml.v3"
)

//CommandExecuteFunc a function type a command can take for executing the content
type CommandExecuteFunc = func(c *Command, client metalcloud.MetalCloudClient) (string, error)

//CommandInitFunc a function type a command can take for initializing the command
type CommandInitFunc = func(c *Command)

//Command defines a command, arguments, description etc
type Command struct {
	Description  string
	Subject      string
	AltSubject   string
	Predicate    string
	AltPredicate string
	FlagSet      *flag.FlagSet
	Arguments    map[string]interface{}
	InitFunc     CommandInitFunc
	ExecuteFunc  CommandExecuteFunc
	Endpoint     string
	Example      string
	UserOnly     bool //set if command is to be visible only to users regardless of endpoint
}

func sameCommand(a *Command, b *Command) bool {
	return a.Subject == b.Subject &&
		a.AltSubject == b.AltSubject &&
		a.Predicate == b.Predicate &&
		a.AltPredicate == b.AltPredicate
}

const _nilDefaultStr = "__NIL__"
const _nilDefaultInt = -14234

//confirms command
func confirmCommand(c *Command, f func() string) (bool, error) {

	if getBoolParam(c.Arguments["autoconfirm"]) {
		return true, nil
	}

	return requestConfirmation(f())
}

//getPtrValueIfExistsOk returns a string or an int from a map of pointers if the key exists
func getPtrValueIfExistsOk(m map[string]interface{}, key string) (interface{}, bool) {

	if v := m[key]; v != nil {
		switch v.(type) {
		case *int:
			if *v.(*int) != _nilDefaultInt {
				return *v.(*int), true
			}
		case *string:
			if *v.(*string) != _nilDefaultStr {
				return *v.(*string), true
			}
		}
	}
	return nil, false
}

//getIDFromStringOk returns the id and true if valid number
func getIDFromStringOk(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	return i, err == nil
}

//verifyParam returns error if param is not present
func getParam(c *Command, label string, name string) (interface{}, error) {
	v := c.Arguments[label]
	if v == nil {
		return nil, fmt.Errorf("-%s cannot be nil", name)
	}
	switch v.(type) {
	case *int:
		if *v.(*int) <= 0 {
			return nil, fmt.Errorf("-%s cannot be <=0", name)
		}
		if *v.(*int) == _nilDefaultInt {
			return nil, fmt.Errorf("-%s is required", name)
		}
	case *string:
		if *v.(*string) == "" {
			return nil, fmt.Errorf("-%s cannot be empty", name)
		}
		if *v.(*string) == _nilDefaultStr {
			return nil, fmt.Errorf("-%s is required", name)
		}
	}
	return v, nil
}

func idOrLabelString(v string) (int, string, bool) {
	if i, ok := getIDFromStringOk(v); ok {
		return i, "", true
	}
	return 0, v, false
}

//idOrLabel returns an int or a string contained in the interface. The last param is true if int is returned.
func idOrLabel(v interface{}) (int, string, bool) {
	switch v.(type) {
	case *int:
		return *v.(*int), "", true
	case *string:
		if i, ok := getIDFromStringOk(*v.(*string)); ok {
			return i, "", true
		}
		return 0, *v.(*string), false
	}
	return -1, "", false
}

type getIDOrDoFunc func(i string) (int, error)

func getIDOrDo(idOrLabel string, f getIDOrDoFunc) (int, error) {
	id, label, isID := idOrLabelString(idOrLabel)
	if !isID {
		return f(label)
	}
	return id, nil
}

func getIntParam(v interface{}) int {
	if v != nil && *v.(*int) != _nilDefaultInt {
		return *v.(*int)
	}
	return 0
}

func getStringParam(v interface{}) string {
	if v != nil && *v.(*string) != _nilDefaultStr {
		return *v.(*string)
	}
	return ""
}

func getBoolParam(v interface{}) bool {
	return v != nil && *v.(*bool)
}

func getStringParamOk(v interface{}) (string, bool) {
	if v != nil && *v.(*string) != _nilDefaultStr {
		return *v.(*string), true
	}
	return "", false
}

func getIntParamOk(v interface{}) (int, bool) {
	if v != nil && *v.(*int) != _nilDefaultInt {
		return *v.(*int), true
	}
	return 0, false
}

func getBoolParamOk(v interface{}) (bool, bool) {
	if v == nil {
		return false, false
	}
	return v != nil && *v.(*bool), true
}

func updateIfIntParamSet(v interface{}, p *int) {
	if v, ok := getIntParamOk(v); ok {
		*p = v
	}
}

func updateIfStringParamSet(v interface{}, p *string) {
	if v, ok := getStringParamOk(v); ok {
		*p = v
	}
}

func updateIfBoolParamSet(v interface{}, p *bool) {
	if v, ok := getBoolParamOk(v); ok {
		*p = v
	}
}

func getRawObjectFromCommand(c *Command, obj interface{}) error {
	readContentfromPipe := getBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = readInputFromPipe()
	} else {

		if configFilePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = readInputFromFile(configFilePath)
		} else {
			return fmt.Errorf("-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return fmt.Errorf("Content cannot be empty")
	}

	format := getStringParam(c.Arguments["format"])

	switch format {
	case "json":
		err := json.Unmarshal(content, obj)
		if err != nil {
			return err
		}
	case "yaml":

		err := yaml.Unmarshal(content, obj)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("input format \"%s\" not supported", format)
	}

	return nil
}

//watch prints the return of the f function every refreshInterval intervals. The interval is in human readable format 1m 1s etc.
func watch(f func() (string, error), refreshInterval string) error {

	interval, err := time.ParseDuration(refreshInterval)
	if err != nil {
		return err
	}

	visualBeepInterval, err := time.ParseDuration("500ms")

	prevLen := 0
	for {
		str, err := f()
		if err != nil {
			return err
		}

		if prevLen != 0 {
			cursor.ClearLinesUp(prevLen)
		}

		cursor.StartOfLine()

		timeStr := fmt.Sprintf("Refreshed at %s", time.Now().Format("01-02-2006 15:04:05"))

		str += "\n" + whiteOnRed(timeStr)

		fmt.Printf(str)

		prevLen = linesStringCount(str) - 1

		time.Sleep(visualBeepInterval)

		cursor.StartOfLine()

		fmt.Printf(timeStr)

		time.Sleep(interval - visualBeepInterval)

	}
}

func linesStringCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

func funcWithWatch(c *Command, client metalcloud.MetalCloudClient, f func(*Command, metalcloud.MetalCloudClient) (string, error)) (string, error) {
	interval, ok := getStringParamOk(c.Arguments["watch"])
	if ok {

		watch(func() (string, error) {
			return f(c, client)
		},
			interval)
	}

	return f(c, client)
}

//getKeyValueMapFromString returns a key value map from a kv string such as key1=value,key2=value.
//the function first does urldecode on the string
//this means that the values can be provided in normal format key1=value,key2=value but also key1%3Dvalue%2Ckey2%3Dvalue
func getKeyValueMapFromString(kvmap string) (map[string]string, error) {

	m := map[string]string{}

	str, err := url.QueryUnescape(kvmap)
	if err != nil {
		return map[string]string{}, err
	}

	pairs := strings.Split(str, ",")

	for _, pair := range pairs {

		pair := strings.Trim(pair, " ")
		elements := strings.Split(pair, "=")
		//if it ends in = we conclude it is an empty string
		if len(elements) == 1 && pair[len(pair)-1] != '=' {
			m[elements[0]] = ""
		}

		if (len(elements) == 2 && pair[0] == '=') || len(elements) > 2 {
			return map[string]string{}, fmt.Errorf("pair has invalid format expecting k=v, given %s", pair)
		}

		m[elements[0]] = elements[1]
	}

	return m, nil
}

//getKeyValueStringFromMap is the reverse operation from getKeyValueMapFromString encoding the value into the key=value,key=value pairs
func getKeyValueStringFromMap(kvmap interface{}) string {

	pairs := []string{}
	m := kvmap.(map[string]interface{})
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(pairs, ",")
}
