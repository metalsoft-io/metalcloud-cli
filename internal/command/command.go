package command

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/atomicgo/cursor"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
)

// CommandExecuteFunc a function type a command can take for executing the content
type CommandExecuteFunc = func(c *Command, client metalcloud.MetalCloudClient) (string, error)
type CommandExecuteFunc2 = func(ctx context.Context, c *Command, client *metalcloud2.APIClient) (string, error)

// CommandInitFunc a function type a command can take for initializing the command
type CommandInitFunc = func(c *Command)

// Command defines a command, arguments, description etc
type Command struct {
	Description         string
	Subject             string
	AltSubject          string
	Predicate           string
	AltPredicate        string
	FlagSet             *flag.FlagSet
	Arguments           map[string]interface{}
	InitFunc            CommandInitFunc
	ExecuteFunc         CommandExecuteFunc
	ExecuteFunc2        CommandExecuteFunc2
	Endpoint            string
	Example             string
	UserOnly            bool   //set if command is to be visible only to users regardless of endpoint
	AdminOnly           bool   //set if command is to be visible only to admins regardless of endpoint
	AdminEndpoint       string //if set will be used instead of Endpoint for admins
	PermissionsRequired []string
}

type CommandTestCase struct {
	Name string
	Cmd  Command
	Good bool
	Id   int
}

const NilDefaultStr = "__NIL__"
const NilDefaultInt = -14234

// confirms command
func ConfirmCommand(c *Command, f func() string) (bool, error) {

	if GetBoolParam(c.Arguments["autoconfirm"]) {
		return true, nil
	}

	return RequestConfirmation(f())
}

// getPtrValueIfExistsOk returns a string or an int from a map of pointers if the key exists
func getPtrValueIfExistsOk(m map[string]interface{}, key string) (interface{}, bool) {

	if v := m[key]; v != nil {
		switch v.(type) {
		case *int:
			if *v.(*int) != NilDefaultInt {
				return *v.(*int), true
			}
		case *string:
			if *v.(*string) != NilDefaultStr {
				return *v.(*string), true
			}
		}
	}
	return nil, false
}

// getIDFromStringOk returns the id and true if valid number
func getIDFromStringOk(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	return i, err == nil
}

// verifyParam returns error if param is not present
func GetParam(c *Command, label string, name string) (interface{}, error) {
	v := c.Arguments[label]
	if v == nil {
		return nil, fmt.Errorf("-%s cannot be nil", name)
	}
	switch v.(type) {
	case *int:
		if *v.(*int) <= 0 {
			return nil, fmt.Errorf("-%s cannot be <=0", name)
		}
		if *v.(*int) == NilDefaultInt {
			return nil, fmt.Errorf("-%s is required", name)
		}
	case *string:
		if *v.(*string) == "" {
			return nil, fmt.Errorf("-%s cannot be empty", name)
		}
		if *v.(*string) == NilDefaultStr {
			return nil, fmt.Errorf("-%s is required", name)
		}
	}
	return v, nil
}

func IdOrLabelString(v string) (int, string, bool) {
	if i, ok := getIDFromStringOk(v); ok {
		return i, "", true
	}
	return 0, v, false
}

// IdOrLabel returns an int or a string contained in the interface. The last param is true if int is returned.
func IdOrLabel(v interface{}) (int, string, bool) {
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

func GetIDOrDo(IdOrLabel string, f getIDOrDoFunc) (int, error) {
	id, label, isID := IdOrLabelString(IdOrLabel)
	if !isID {
		return f(label)
	}
	return id, nil
}

func GetIntParam(v interface{}) int {
	if v != nil && *v.(*int) != NilDefaultInt {
		return *v.(*int)
	}
	return 0
}

func GetStringParam(v interface{}) string {
	if v != nil && *v.(*string) != NilDefaultStr {
		return *v.(*string)
	}
	return ""
}

func GetBoolParam(v interface{}) bool {
	return v != nil && *v.(*bool)
}

func GetStringParamOk(v interface{}) (string, bool) {
	if v != nil && *v.(*string) != NilDefaultStr {
		return *v.(*string), true
	}
	return "", false
}

func GetIntParamOk(v interface{}) (int, bool) {
	if v != nil && *v.(*int) != NilDefaultInt {
		return *v.(*int), true
	}
	return 0, false
}

func GetBoolParamOk(v interface{}) (bool, bool) {
	if v == nil {
		return false, false
	}
	return v != nil && *v.(*bool), true
}

func UpdateIfIntParamSet(v interface{}, p *int) {
	if v, ok := GetIntParamOk(v); ok {
		*p = v
	}
}

func UpdateIfStringParamSet(v interface{}, p *string) {
	if v, ok := GetStringParamOk(v); ok {
		*p = v
	}
}

func UpdateIfBoolParamSet(v interface{}, p *bool) {
	if v, ok := GetBoolParamOk(v); ok {
		*p = v
	}
}

func GetRawObjectFromCommand(c *Command, obj interface{}) error {
	readContentfromPipe := GetBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = configuration.ReadInputFromPipe()
	} else {

		if configFilePath, ok := GetStringParamOk(c.Arguments["read_config_from_file"]); ok {
			content, err = configuration.ReadInputFromFile(configFilePath)
		} else {
			return fmt.Errorf("-raw-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return fmt.Errorf("Content cannot be empty")
	}

	format := GetStringParam(c.Arguments["format"])
	switch format {
	case "json":
		err := json.Unmarshal(content, obj)
		if err != nil {
			return fmt.Errorf("error unmarshalling json: %v. Make sure the raw config file is in the correct format", err)
		}
	case "yaml":
		err := yaml.Unmarshal(content, obj)
		if err != nil {
			return fmt.Errorf("error unmarshalling yaml: %v. Make sure the raw config file is in the correct format", err)
		}
	default:
		return fmt.Errorf("input format \"%s\" not supported", format)
	}

	return nil
}

// Watch prints the return of the f function every refreshInterval intervals. The interval is in human readable format 1m 1s etc.
func Watch(f func() (string, error), refreshInterval string) error {
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

		str += "\n" + colors.WhiteOnRed(timeStr)

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

func FuncWithWatch(c *Command, client metalcloud.MetalCloudClient, f func(*Command, metalcloud.MetalCloudClient) (string, error)) (string, error) {
	interval, ok := GetStringParamOk(c.Arguments["watch"])
	if ok {
		Watch(func() (string, error) {
			return f(c, client)
		},
			interval)
	}

	return f(c, client)
}

// GetKeyValueMapFromString returns a key value map from a kv string such as key1=value,key2=value.
// the function first does urldecode on the string
// this means that the values can be provided in normal format key1=value,key2=value but also key1%3Dvalue%2Ckey2%3Dvalue
func GetKeyValueMapFromString(kvmap string) (map[string]string, error) {
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

// GetKeyValueStringFromMap is the reverse operation from GetKeyValueMapFromString encoding the value into the key=value,key=value pairs
func GetKeyValueStringFromMap(kvmap interface{}) string {

	switch m := kvmap.(type) {
	case map[string]interface{}:
		pairs := []string{}
		for k, v := range m {
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
		}
		return strings.Join(pairs, ",")
	case []interface{}:
		return ""
	}

	return ""
}

func RequestConfirmation(s string) (bool, error) {
	yes, err := RequestInput(s)
	if err != nil {
		return false, err
	}

	return strings.Trim(string(yes), "\r\n ") == "yes", nil
}

func RequestInputSilent(s string) ([]byte, error) {

	fmt.Fprintf(configuration.GetStdout(), s)
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return []byte{}, err
	}

	content, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return []byte{}, err
	}

	defer terminal.Restore(0, oldState)
	return content, nil
}

func RequestInput(s string) ([]byte, error) {
	fmt.Fprintf(configuration.GetStdout(), s)
	reader := bufio.NewReader(configuration.GetStdin())
	content, err := reader.ReadBytes('\n')

	if err != nil && err != io.EOF {
		return content, err
	}

	return content, nil
}

func RequestInputString(s string) (string, error) {

	fmt.Fprintf(configuration.GetStdout(), s)
	reader := bufio.NewReader(configuration.GetStdin())
	content, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return content, err
	}

	return content, nil
}

func validateArguments(args []string) (string, string, int) {
	subject, predicate := NilDefaultStr, NilDefaultStr
	count := 0

	if len(args) >= 2 && !strings.HasPrefix(args[1], "-") {
		subject = args[1]
		count++
	}
	if len(args) >= 3 && !strings.HasPrefix(args[2], "-") {
		predicate = args[2]
		count++
	}

	return subject, predicate, count
}

func helpMessage(err error, subject string, predicate string) error {
	message := err.Error()

	swaggerErr, ok := err.(metalcloud2.GenericSwaggerError)
	if ok {
		if strings.Contains(swaggerErr.Error(), "404") || strings.Contains(swaggerErr.Error(), "Not Found") {
			message = "This version of CLI (" + configuration.Version + ") is only compatible with a controller versioned 6.4 and above."
		} else {
			message = swaggerErr.Error() + " [ " + string(swaggerErr.Body()) + " ]"
		}
	}

	if predicate != NilDefaultStr {
		return fmt.Errorf("%s\nUse '%s %s -h' for syntax help", message, subject, predicate)
	}

	return fmt.Errorf("%s\nUse '%s -h' for syntax help", message, subject)
}

func ExecuteCommand(args []string, commands []Command, clients map[string]metalcloud.MetalCloudClient, client2 *metalcloud2.APIClient, permissions []string) error {
	subject, predicate, count := validateArguments(args)

	if count == 1 {
		commandsForSubject := filterCommandsBySubject(subject, commands)
		if len(commandsForSubject) > 0 {
			foundNilPredicate := false

			for _, c := range commandsForSubject {
				if c.Predicate == NilDefaultStr {
					foundNilPredicate = true
				}
			}

			if !foundNilPredicate {
				return fmt.Errorf("invalid command: %s", getPossiblePredicatesForSubjectHelp(subject, commandsForSubject))
			}
		}
	}

	cmd := locateCommand(predicate, subject, commands)

	if cmd == nil {
		return fmt.Errorf("invalid command! Use 'help' for a list of commands")
	}

	cmd.InitFunc(cmd)

	if flag := cmd.FlagSet.Lookup("no-color"); flag == nil {
		cmd.Arguments["no_color"] = cmd.FlagSet.Bool("no-color", false, colors.Green("(Flag)")+" Disable coloring.")
	}

	//disable default usage
	cmd.FlagSet.Usage = func() {}

	colors.SetColoringEnabled(true)

	noColorEnabled := false
	commandHelp := false

	for _, a := range args {
		if a == "--no-color" || a == "-no-color" {
			noColorEnabled = true
		}

		if a == "-h" || a == "-help" || a == "--help" {
			commandHelp = true
		}
	}

	if noColorEnabled {
		colors.SetColoringEnabled(false)
	}

	if commandHelp {
		return fmt.Errorf(GetCommandHelp(*cmd, true))
	}

	err := cmd.FlagSet.Parse(args[count+1:])

	if err != nil {
		return helpMessage(err, subject, predicate)
	}

	endpoint := cmd.Endpoint

	if slices.Contains(permissions, ADMIN_ACCESS) && cmd.AdminEndpoint != "" {
		endpoint = cmd.AdminEndpoint
	}

	var ret string
	if cmd.ExecuteFunc2 != nil {
		ret, err = cmd.ExecuteFunc2(context.Background(), cmd, client2)
	} else {
		client, ok := clients[endpoint]
		if !ok {
			return fmt.Errorf("Client not set for endpoint %s on command %s %s", endpoint, subject, predicate)
		}

		ret, err = cmd.ExecuteFunc(cmd, client)
	}
	if err != nil {
		return helpMessage(err, subject, predicate)
	}

	fmt.Fprintf(configuration.GetStdout(), ret)

	return nil
}

// identifies command, returns nil if no matching command found
func locateCommand(predicate string, subject string, commands []Command) *Command {
	for _, c := range commands {
		if (c.Subject == subject || c.AltSubject == subject) &&
			(c.Predicate == predicate || c.AltPredicate == predicate) {
			return &c
		}
	}
	return nil
}

// identifies commands for given subjects
func filterCommandsBySubject(subject string, commands []Command) []Command {
	cmds := []Command{}
	for _, c := range commands {
		if c.Subject == subject {
			cmds = append(cmds, c)
		}
	}
	return cmds
}

func getArgumentHelp(f *flag.Flag) string {
	if len(f.Name) == 1 {
		return fmt.Sprintf("\t  -%-25s %s\n", f.Name, f.Usage)
	}

	return fmt.Sprintf("\t  --%-24s %s\n", f.Name, f.Usage)

}

func commandHelpSummary(cmd Command) string {
	var sb strings.Builder

	command := cmd.Subject

	if cmd.Predicate != NilDefaultStr {
		command = fmt.Sprintf("%s %s", cmd.Subject, cmd.Predicate)
	}

	alternate := ""
	if cmd.AltPredicate != NilDefaultStr && cmd.AltSubject != NilDefaultStr {
		alternate = fmt.Sprintf(" (alternatively use \"%s %s\")", colors.Bold(cmd.AltSubject), colors.Bold(cmd.AltPredicate))
	}
	cmdHelpSummary := fmt.Sprintf("Command: %-36s %s%s\n",
		colors.Bold(command),
		cmd.Description,
		alternate)

	sb.WriteString(cmdHelpSummary)

	return sb.String()
}

func getPossiblePredicatesForSubjectHelp(subject string, cmds []Command) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Possible commands for subject %s:\n", colors.Bold(subject)))
	for _, cmd := range cmds {
		sb.WriteString(fmt.Sprintf("\t%s %s - %s\n", colors.Bold(cmd.Subject), colors.Bold(cmd.Predicate), cmd.Description))
	}
	return sb.String()
}

func GetCommandHelp(cmd Command, showArguments bool) string {
	var sb strings.Builder
	var c string
	if cmd.Predicate != NilDefaultStr {
		c = fmt.Sprintf("%s %s", cmd.Subject, cmd.Predicate)
	} else {
		c = fmt.Sprintf("%s", cmd.Subject)
	}

	if showArguments {
		sb.WriteString(commandHelpSummary(cmd))
		cmd.FlagSet.VisitAll(func(f *flag.Flag) {
			sb.WriteString(getArgumentHelp(f))
		})

		h := flag.Flag{
			Name:  "h",
			Usage: "Show command help and exit.",
		}

		sb.WriteString(getArgumentHelp(&h))

		if cmd.Example != "" {
			sb.WriteString("\nExample:\n")
			sb.WriteString(cmd.Example)
			sb.WriteString("\n")
		}

	} else {
		sb.WriteString(fmt.Sprintf("\t%-40s %-24s", c, cmd.Description))
	}

	return sb.String()
}

// MakeCommand utility function that creates a command from a kv map
func MakeCommand(arguments map[string]interface{}) Command {
	cmd := Command{
		Arguments: map[string]interface{}{},
	}

	for k, v := range arguments {
		switch v.(type) {
		case int:
			x := v.(int)
			cmd.Arguments[k] = &x
		case string:
			x := v.(string)
			cmd.Arguments[k] = &x
		case bool:
			x := v.(bool)
			cmd.Arguments[k] = &x
		}
	}

	return cmd
}

func MakeEmptyCommand() Command {
	return MakeCommand(map[string]interface{}{})
}

// checks command with and without return_id
func TestCreateCommand(f CommandExecuteFunc, cases []CommandTestCase, client metalcloud.MetalCloudClient, t *testing.T) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			//test without return id

			_, err := f(&c.Cmd, client)
			if c.Good {

				if err != nil {
					t.Errorf("error thrown: %v", err)
				}

			} else {
				if err == nil {
					t.Errorf("Should have thrown error")
				}

			}
			if c.Id != 0 {
				//test with return id
				cmdWithReturn := c.Cmd
				bTrue := true
				cmdWithReturn.Arguments["return_id"] = &bTrue
				ret, err := f(&c.Cmd, client)
				if c.Good {
					if ret != fmt.Sprintf("%d", c.Id) {
						t.Errorf("id not returned or incorrect. Expected %s got %s", ret, fmt.Sprintf("%d", c.Id))
					}
					if err != nil {
						t.Errorf("error thrown: %v", err)
					}

				} else {
					if err == nil {
						t.Errorf("Should have thrown error")
					}
				}
			}
		})
	}
}

func TestCommandWithConfirmation(f CommandExecuteFunc, cmd Command, client metalcloud.MetalCloudClient, t *testing.T) {
	var stdin bytes.Buffer
	var stdout bytes.Buffer

	configuration.SetConsoleIOChannel(&stdin, &stdout)
	internalCmd := cmd
	//test with no autoconfirm should throw error
	bTrue := false
	internalCmd.Arguments["autoconfirm"] = &bTrue

	_, err := f(&internalCmd, client)
	if err != nil {
		t.Errorf("error thrown: %v", err)
	}

	//test with  autoconfirm should not throw error

	bTrue = true
	internalCmd.Arguments["autoconfirm"] = &bTrue

	_, err = f(&internalCmd, client)
	if err != nil {
		t.Errorf("error thrown: %v", err)
	}
}

// checks the various outputs
func TestGetCommand(f CommandExecuteFunc, cases []CommandTestCase, client metalcloud.MetalCloudClient, firstRow map[string]interface{}, t *testing.T) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := f(&c.Cmd, client)
			if c.Good {
				if err != nil {
					t.Errorf("error thrown: %v", err)
				}

			} else {
				if err == nil {
					t.Errorf("Should have thrown error")
				}
			}
		})
	}
}

// checks the various outputs
func TestListCommand(f CommandExecuteFunc, cmd *Command, client metalcloud.MetalCloudClient, firstRow map[string]interface{}, t *testing.T) {
	//RegisterTestingT(t)

	c := cmd
	if cmd == nil {
		emptyCmd := MakeEmptyCommand()
		c = &emptyCmd
	}

	//test plaintext
	ret, err := f(c, client)
	Expect(err).To(BeNil())

	//test json output
	cmdWithFormat := c
	format := "json"
	cmdWithFormat.Arguments["format"] = &format

	ret, err = f(cmdWithFormat, client)
	Expect(err).To(BeNil())
	Expect(JSONFirstRowEquals(ret, firstRow)).To(BeNil())

	//test csv output
	format = "csv"
	cmdWithFormat.Arguments["format"] = &format

	ret, err = f(cmdWithFormat, client)
	Expect(err).To(BeNil())

	//this is not reliable as the first row sometimes changes.
	//Expect(CSVFirstRowEquals(ret, firstRow)).To(BeNil())
}

// JSONFirstRowEquals checks if values of the table returned in the json match the values provided. Type is not checked (we check string equality)
func JSONFirstRowEquals(jsonString string, testVals map[string]interface{}) error {
	m, err := JSONUnmarshal(jsonString)
	if err != nil {
		return err
	}

	firstRow := m[0].(map[string]interface{})

	for k, v := range testVals {
		if fmt.Sprintf("%+v", firstRow[k]) != fmt.Sprintf("%+v", v) {
			return fmt.Errorf("values for key %s do not match:  expected '%+v' provided '%+v'", k, v, firstRow[k])
		}
	}

	return nil
}

func JSONUnmarshal(jsonString string) ([]interface{}, error) {
	var m []interface{}
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func CSVUnmarshal(csvString string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(csvString))

	return reader.ReadAll()
}

// KeysOfMapAsString returns the keys of a map as a string separated by " "
func KeysOfMapAsString(m map[string]interface{}) string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, " ")
}

// GenerateCommandTestCases generate commands with wrong arguments by cycling through all of them
func GenerateCommandTestCases(arguments map[string]interface{}) []CommandTestCase {
	cmds := []CommandTestCase{}

	//turn keys to array
	keys := []string{}
	for k := range arguments {
		keys = append(keys, k)
	}

	l := len(arguments)
	if l == 1 {
		return []CommandTestCase{
			{
				Name: KeysOfMapAsString(arguments),
				Cmd:  MakeCommand(arguments),
				Good: false,
			},
		}
	}

	for i := 0; i < l; i++ {
		//add all the arguments list
		args := map[string]interface{}{}
		for j := 0; j < l; j++ {
			if i == j {
				continue
			}
			args[keys[j]] = arguments[keys[j]]
		}

		newCmds := GenerateCommandTestCases(args)
		newCmds = append(newCmds,
			CommandTestCase{
				Name: KeysOfMapAsString(args),
				Cmd:  MakeCommand(args),
				Good: false,
			})

		for _, newCmd := range newCmds {
			duplicate := false
			for _, v := range cmds {
				if reflect.DeepEqual(v.Cmd.Arguments, newCmd.Cmd.Arguments) {
					duplicate = true
				}
			}
			if !duplicate {
				cmds = append(cmds, newCmd)
			}
		}

	}

	return cmds
}

func GetUserFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.User, error) {
	user, err := GetParam(c, "user_id", paramName)
	if err != nil {
		return nil, err
	}

	id, email, isID := IdOrLabel(user)

	if isID {
		return client.UserGet(id)
	} else {
		return client.UserGetByEmail(email)
	}
}

// GetInfrastructureFromCommand returns an Infrastructure object using the infrastructure_id_or_label argument
func GetInfrastructureFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.Infrastructure, error) {
	m, err := GetParam(c, "infrastructure_id_or_label", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := IdOrLabel(m)

	if isID {
		return client.InfrastructureGet(id)
	}

	return client.InfrastructureGetByLabel(label)
}

func GetInstanceArrayFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {
	m, err := GetParam(c, "instance_array_id_or_label", paramName)
	if err != nil {
		return nil, err
	}
	id, label, isID := IdOrLabel(m)
	if isID {
		return client.InstanceArrayGet(id)
	}
	return client.InstanceArrayGetByLabel(label)
}

func GetOSTemplateFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient, decryptPasswd bool) (*metalcloud.OSTemplate, error) {
	v, err := GetParam(c, "template_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := IdOrLabel(v)

	if isID {
		return client.OSTemplateGet(id, decryptPasswd)
	}

	list, err := client.OSTemplates()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.VolumeTemplateLabel == label {
			return &s, nil
		}
	}

	if isID {
		return nil, fmt.Errorf("template %d not found", id)
	}

	return nil, fmt.Errorf("template %s not found", label)
}

func GetWorkflowFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.Workflow, error) {
	v, err := GetParam(c, "workflow_id_or_label", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := IdOrLabel(v)

	if isID {
		return client.WorkflowGet(id)
	}

	list, err := client.Workflows()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.WorkflowLabel == label {
			return &s, nil
		}
	}

	if isID {
		return nil, fmt.Errorf("workflow %d not found", id)
	}

	return nil, fmt.Errorf("workflow %s not found", label)
}

// asset_id_or_name
func GetOSAssetFromCommand(paramName string, internalParamName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.OSAsset, error) {
	v, err := GetParam(c, internalParamName, paramName)
	if err != nil {
		return nil, err
	}

	id, name, isID := IdOrLabel(v)

	if isID {
		return client.OSAssetGet(id)
	}

	list, err := client.OSAssets()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.OSAssetFileName == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate asset with file name '%s'", name)
}

func GetNetworkOperatingSystemFromCommand(c *Command) (*metalcloud.NetworkOperatingSystem, error) {
	var networkOperatingSystem = metalcloud.NetworkOperatingSystem{}

	nosSwitchDriver := GetStringParam(c.Arguments["network_os_switch_driver"])
	if nosSwitchDriver != "" {
		networkOperatingSystem.OperatingSystemSwitchDriver = nosSwitchDriver
	} else {
		return nil, fmt.Errorf("network-os-switch-driver is required")
	}

	nosSwitchRole := GetStringParam(c.Arguments["network_os_switch_role"])
	if nosSwitchRole != "" {
		networkOperatingSystem.OperatingSystemSwitchRole = nosSwitchRole
	}

	nosVersion := GetStringParam(c.Arguments["network_os_version"])
	if nosVersion != "" {
		networkOperatingSystem.OperatingSystemVersion = nosVersion
	} else {
		return nil, fmt.Errorf("network-os-version is required")
	}

	nosArchitecture := GetStringParam(c.Arguments["network_os_architecture"])
	if nosArchitecture != "" {
		networkOperatingSystem.OperatingSystemArchitecture = nosArchitecture
	} else {
		return nil, fmt.Errorf("network-os-architecture is required")
	}

	nosVendor := GetStringParam(c.Arguments["network_os_vendor"])
	if nosVendor != "" {
		networkOperatingSystem.OperatingSystemVendor = nosVendor
	} else {
		return nil, fmt.Errorf("network-os-vendor is required")
	}

	nosMachine := GetStringParam(c.Arguments["network_os_machine"])
	if nosMachine != "" {
		networkOperatingSystem.OperatingSystemMachine = nosMachine
	} else {
		return nil, fmt.Errorf("network-os-machine is required")
	}

	nosDatacenterName := GetStringParam(c.Arguments["network_os_datacenter_name"])
	if nosDatacenterName != "" {
		networkOperatingSystem.OperatingSystemDatacenterName = nosDatacenterName
	}

	return &networkOperatingSystem, nil
}

func GetOperatingSystemFromCommand(c *Command) (*metalcloud.OperatingSystem, error) {
	var operatingSystem = metalcloud.OperatingSystem{}
	present := false

	if osType, ok := GetStringParamOk(c.Arguments["os_type"]); ok {
		present = true
		operatingSystem.OperatingSystemType = osType
	}

	if osVersion, ok := GetStringParamOk(c.Arguments["os_version"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the operating system flags are missing")
		}
		operatingSystem.OperatingSystemVersion = osVersion
	} else if present {
		return nil, fmt.Errorf("os-version is required")
	}

	if osArchitecture, ok := GetStringParamOk(c.Arguments["os_architecture"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the operating system flags are missing")
		}
		operatingSystem.OperatingSystemArchitecture = osArchitecture
	} else if present {
		return nil, fmt.Errorf("os-architecture is required")
	}

	return &operatingSystem, nil
}
