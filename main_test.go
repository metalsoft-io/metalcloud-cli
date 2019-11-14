package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestValidateAPIKey(t *testing.T) {
	RegisterTestingT(t)

	Expect(len(RandStringBytes(64))).To(Equal(64))
	goodKey := fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(63))

	badKey1 := fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(64))
	badKey2 := fmt.Sprintf(":%s", RandStringBytes(63))

	Expect(validateAPIKey(goodKey)).To(BeNil())
	Expect(validateAPIKey(badKey1)).NotTo(BeNil())
	Expect(validateAPIKey(badKey2)).NotTo(BeNil())

}

func TestInitClient(t *testing.T) {

	envs := []string{
		"METALCLOUD_USER_EMAIL",
		"METALCLOUD_API_KEY",
		"METALCLOUD_ENDPOINT",
		"METALCLOUD_DATACENTER",
	}
	//remember the current env values, clear them during the test
	currentEnvVals := map[string]string{}
	for _, e := range envs {
		if v, ok := os.LookupEnv(e); ok {
			currentEnvVals[e] = v
			os.Unsetenv(e)
		}
	}

	if _, err := initClient(); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_USER_EMAIL", "user")

	if _, err := initClient(); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_API_KEY", fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(63)))

	if _, err := initClient(); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_ENDPOINT", "endpoint")

	if _, err := initClient(); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_DATACENTER", "dc")

	if _, err := initClient(); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	client, err := initClient()
	if client == nil || err == nil {
		t.Errorf("cannot initialize metalcloud client %v", err)
	}

	//put back the env values
	for k, v := range currentEnvVals {
		os.Setenv(k, v)
	}

}

func TestExecuteCommand(t *testing.T) {
	RegisterTestingT(t)

	execFuncExecuted := false
	initFuncExecuted := false
	commands := []Command{
		Command{
			Subject:      "tests",
			AltSubject:   "s",
			Predicate:    "testp",
			AltPredicate: "p",
			FlagSet:      flag.NewFlagSet(RandStringBytes(10), flag.ExitOnError),
			InitFunc: func(c *Command) {
				c.Arguments = map[string]interface{}{
					"cmd": c.FlagSet.Int(RandStringBytes(10), 0, "Random param"),
				}
				initFuncExecuted = true
			},
			ExecuteFunc: func(c *Command, client MetalCloudClient) (string, error) {
				execFuncExecuted = true
				return "", nil
			},
		},
	}

	//check with wrong commands first, should return err

	err := executeCommand([]string{"", "test", "test"}, commands, nil)
	Expect(err).NotTo(BeNil())

	execFuncExecuted = false
	initFuncExecuted = false

	//should execute stuff help and not return error
	err = executeCommand([]string{"", "p", "s"}, commands, nil)
	Expect(err).To(BeNil())
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())

	execFuncExecuted = false
	initFuncExecuted = false

	//should execute stuff help and not return error
	err = executeCommand([]string{"", "testp", "tests"}, commands, nil)
	Expect(err).To(BeNil())
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())

}

func TestGetDatacenter(t *testing.T) {
	RegisterTestingT(t)
	dc := RandStringBytes(10)
	os.Setenv("METALCLOUD_DATACENTER", dc)
	Expect(GetDatacenter()).To(Equal(dc))
}

func TestGetCommandHelp(t *testing.T) {
	RegisterTestingT(t)
	cmd := Command{
		Description:  "Lists available volume templates",
		Subject:      "tests",
		AltSubject:   "s",
		Predicate:    "testp",
		AltPredicate: "p",
		FlagSet:      flag.NewFlagSet(RandStringBytes(10), flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"cmd": c.FlagSet.Int(RandStringBytes(10), 0, "Random param"),
			}
		},
		ExecuteFunc: func(c *Command, client MetalCloudClient) (string, error) {
			return "", nil
		}}

	cmd.InitFunc(&cmd)
	s := getCommandHelp(cmd)
	Expect(s).To(ContainSubstring(cmd.Description))
	Expect(s).To(ContainSubstring("Random param"))

}

func TestGetHelp(t *testing.T) {
	RegisterTestingT(t)
	cmds := getCommands()

	s := getHelp()
	for _, c := range cmds {
		Expect(s).To(ContainSubstring(c.Description))

		c.FlagSet.VisitAll(func(f *flag.Flag) {
			Expect(s).To(ContainSubstring(f.Name))
			Expect(s).To(ContainSubstring(f.Usage))
		})

	}

}
