package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"syscall"
	"testing"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"

	//. "github.com/onsi/gomega"
	gomock "github.com/golang/mock/gomock"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestCheckForDuplicates(t *testing.T) {

	var commands []Command

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	client.EXPECT().GetEndpoint().Return("user").AnyTimes()
	clients := map[string]metalcloud.MetalCloudClient{
		"": client,
	}

	commands = getCommands(clients)

	for i := 0; i < len(commands); i++ {
		for j := i + 1; j < len(commands); j++ {

			a := commands[i]
			b := commands[j]

			if a.Description == b.Description {
				t.Errorf("commands have same description:\na=%+v\nb=%+v", a, b)
			}

			if sameCommand(&a, &b) {
				t.Errorf("commands have same commands:\na=%+v\nb=%+v", a, b)
			}

			sf1 := reflect.ValueOf(a.ExecuteFunc)
			sf2 := reflect.ValueOf(b.ExecuteFunc)

			if sf1.Pointer() == sf2.Pointer() {
				t.Errorf("commands have same executeFunc:\na=%+v\nb=%+v", a, b)
			}

		}
	}
}

func TestSimpleArgument(t *testing.T) {

	var executed = false

	cmd := Command{
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "c",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_instance_count": c.FlagSet.Int("instance_count", 0, "Instance count of this instance array"),
				"instance_array_instance_label": c.FlagSet.String("label", "", "Instance array's label"),
			}
		},
		ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
			executed = true
			return "retstr", nil
		},
	}

	cmd.InitFunc(&cmd)

	argv := []string{
		"-instance_count=3",
		"-label=test",
	}

	err := cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	iaCount := cmd.Arguments["instance_array_instance_count"].(*int)
	if iaCount == nil || *iaCount != 3 {
		t.Errorf("instance_array_instance_count expected to be %d\n\twas %d", 3, *iaCount)
	}

	iaLabel := cmd.Arguments["instance_array_instance_label"].(*string)

	if iaLabel == nil || *iaLabel != "test" {
		t.Errorf("instance_array_label expected to be %s\n\twas %s", "test", *iaLabel)
	}

	argv = []string{
		"instance_countasdad=3",
		"la33bel=\"test\"",
	}

	err = cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	ret, err := cmd.ExecuteFunc(&cmd, nil)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !executed || ret != "retstr" {
		t.Errorf("ExecuteFunction not called properly")
	}
}

func TestConfirmFunc(t *testing.T) {
	RegisterTestingT(t)
	v1 := 10

	cmd := Command{
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "c",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "autoconfirm text"),
			}
		},
		ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", nil
		},
	}

	cmd.InitFunc(&cmd)

	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)

	stdin.Write([]byte("yes\n"))

	//check without autoconfirm
	ok, err := confirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

	s, err := stdout.ReadString(byte('\n'))
	Expect(s).To(ContainSubstring("Reverting infrastructure"))

	//check with autoconfirm
	argv := []string{
		"-autoconfirm",
	}

	err = cmd.FlagSet.Parse(argv)
	Expect(err).To(BeNil())

	ok, err = confirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

	s, err = stdout.ReadString(byte('\n'))
	Expect(s).To(BeEmpty())

}

func TestGetIfNotDefaultOk(t *testing.T) {
	RegisterTestingT(t)

	i := 10
	s := "test"
	f := 10.2
	m := map[string]interface{}{
		"testInt":       &i,
		"testString":    &s,
		"testWrongType": &f,
	}

	v, ok := getPtrValueIfExistsOk(m, "testInt")
	Expect(v).To(Equal((10)))
	Expect(ok).To(BeTrue())

	v, ok = getPtrValueIfExistsOk(m, "testString")
	Expect(v).To(Equal(("test")))
	Expect(ok).To(BeTrue())

	v, ok = getPtrValueIfExistsOk(m, "testWrongString")
	Expect(v).To(BeNil())
	Expect(ok).To(BeFalse())

	v, ok = getPtrValueIfExistsOk(m, "testWrongType")
	Expect(v).To(BeNil())
	Expect(ok).To(BeFalse())
}

func TestIdOrLabel(t *testing.T) {
	RegisterTestingT(t)

	i := "test"
	id, label, isID := idOrLabel(&i)
	Expect(id).To(Equal(0))
	Expect(label).To(Equal("test"))
	Expect(isID).To(BeFalse())

	i = "100"
	id, label, isID = idOrLabel(&i)
	Expect(id).To(Equal(100))
	Expect(label).To(Equal(""))
	Expect(isID).To(BeTrue())

	ii := 100
	id, label, isID = idOrLabel(&ii)
	Expect(id).To(Equal(100))
	Expect(label).To(Equal(""))
	Expect(isID).To(BeTrue())
}

//checks the various outputs
func testListCommand(f CommandExecuteFunc, cmd *Command, client metalcloud.MetalCloudClient, firstRow map[string]interface{}, t *testing.T) {
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

//JSONFirstRowEquals checks if values of the table returned in the json match the values provided. Type is not checked (we check string equality)
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

//checks the various outputs
func testGetCommand(f CommandExecuteFunc, cases []CommandTestCase, client metalcloud.MetalCloudClient, firstRow map[string]interface{}, t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := f(&c.cmd, client)
			if c.good {

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

type CommandTestCase struct {
	name string
	cmd  Command
	good bool
	id   int
}

//KeysOfMapAsString returns the keys of a map as a string separated by " "
func KeysOfMapAsString(m map[string]interface{}) string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, " ")
}

//MakeCommand utility function that creates a command from a kv map
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

//GenerateCommandTestCases generate commands with wrong arguments by cycling through all of them
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
				name: KeysOfMapAsString(arguments),
				cmd:  MakeCommand(arguments),
				good: false,
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
				name: KeysOfMapAsString(args),
				cmd:  MakeCommand(args),
				good: false,
			})

		for _, newCmd := range newCmds {
			duplicate := false
			for _, v := range cmds {
				if reflect.DeepEqual(v.cmd.Arguments, newCmd.cmd.Arguments) {
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

//checks command with and without return_id
func testCreateCommand(f CommandExecuteFunc, cases []CommandTestCase, client metalcloud.MetalCloudClient, t *testing.T) {

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//test without return id

			_, err := f(&c.cmd, client)
			if c.good {

				if err != nil {
					t.Errorf("error thrown: %v", err)
				}

			} else {
				if err == nil {
					t.Errorf("Should have thrown error")
				}

			}
			if c.id != 0 {
				//test with return id
				cmdWithReturn := c.cmd
				bTrue := true
				cmdWithReturn.Arguments["return_id"] = &bTrue
				ret, err := f(&c.cmd, client)
				if c.good {
					if ret != fmt.Sprintf("%d", c.id) {
						t.Error("id not returned")
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

func testCommandWithConfirmation(f CommandExecuteFunc, cmd Command, client metalcloud.MetalCloudClient, t *testing.T) {
	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)
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

func TestMakeWrongCommand(t *testing.T) {
	RegisterTestingT(t)
	cmds := GenerateCommandTestCases(map[string]interface{}{"1": 1, "2": 2, "3": 3, "4": 4})

	Expect(len(cmds)).To(Equal(14))

}

func TestGetRawObjectFromCommand(t *testing.T) {
	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice
	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	if err != nil {
		t.Error(err)
	}

	f, err := ioutil.TempFile("/tmp", "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_switchDeviceFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	cmd := MakeCommand(map[string]interface{}{
		"read_config_from_file": f.Name(),
		"format":                "json",
	})

	var sw2 metalcloud.SwitchDevice

	err = getRawObjectFromCommand(&cmd, &sw2)
	Expect(err).To(BeNil())
	Expect(sw.NetworkEquipmentPrimarySANSubnetPool).To(Equal("100.64.0.0"))
}
