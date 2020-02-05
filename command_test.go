package main

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"testing"

	//. "github.com/onsi/gomega"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestCheckForDuplicates(t *testing.T) {

	var commands []Command

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	client.EXPECT().GetEndpoint().Return("user").AnyTimes()
	clients := map[string]interfaces.MetalCloudClient{
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
		ExecuteFunc: func(c *Command, client interfaces.MetalCloudClient) (string, error) {
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
		ExecuteFunc: func(c *Command, client interfaces.MetalCloudClient) (string, error) {
			return "", nil
		},
	}

	cmd.InitFunc(&cmd)

	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)

	stdin.Write([]byte("yes\n"))

	//check without autoconfirm
	ok := confirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
	Expect(ok).To(BeTrue())
	s, err := stdout.ReadString(byte('\n'))
	Expect(s).To(ContainSubstring("Reverting infrastructure"))

	//check with autoconfirm
	argv := []string{
		"-autoconfirm",
	}

	err = cmd.FlagSet.Parse(argv)
	Expect(err).To(BeNil())

	ok = confirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
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
