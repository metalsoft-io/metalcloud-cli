package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
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

	badKey1 := fmt.Sprintf("asdasd:asd%s", RandStringBytes(67))
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
	}
	//remember the current env values, clear them during the test
	currentEnvVals := map[string]string{}
	for _, e := range envs {
		if v, ok := os.LookupEnv(e); ok {
			currentEnvVals[e] = v
			os.Unsetenv(e)
		}
	}

	if _, err := initClient("METALCLOUD_ENDPOINT"); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_USER_EMAIL", "user")

	if _, err := initClient("METALCLOUD_ENDPOINT"); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_API_KEY", fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(63)))

	if _, err := initClient("METALCLOUD_ENDPOINT"); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	os.Setenv("METALCLOUD_ENDPOINT", "endpoint")

	if _, err := initClient("METALCLOUD_ENDPOINT"); err == nil {
		t.Errorf("Should have been able to test for missing env")
	}

	client, err := initClient("METALCLOUD_ENDPOINT")
	if client == nil || err == nil {
		t.Errorf("cannot initialize metalcloud client %v", err)
	}

	//put back the env values
	for k, v := range currentEnvVals {
		os.Setenv(k, v)
	}

}

func TestInitClients(t *testing.T) {
	RegisterTestingT(t)

	envs := []string{
		"METALCLOUD_USER_EMAIL",
		"METALCLOUD_API_KEY",
		"METALCLOUD_ENDPOINT",
		"METALCLOUD_ADMIN",
	}

	currentEnvVals := map[string]string{}
	for _, e := range envs {
		if v, ok := os.LookupEnv(e); ok {
			currentEnvVals[e] = v
			os.Unsetenv(e)
		}
	}

	os.Setenv("METALCLOUD_USER_EMAIL", "user@user.com")
	os.Setenv("METALCLOUD_API_KEY", fmt.Sprintf("%d:%s", rand.Intn(100), RandStringBytes(63)))
	os.Setenv("METALCLOUD_ENDPOINT", "http://test1/1")

	clients, err := initClients()
	Expect(err).To(BeNil())
	Expect(clients).To(Not(BeNil()))
	Expect(clients[UserEndpoint]).To(Not(BeNil()))
	Expect(clients[ExtendedEndpoint]).To(BeNil())
	Expect(clients[DeveloperEndpoint]).To(BeNil())

	os.Setenv("METALCLOUD_ADMIN", "true")

	clients, err = initClients()
	Expect(clients).To(Not(BeNil()))
	Expect(clients[UserEndpoint]).To(Not(BeNil()))
	Expect(clients[ExtendedEndpoint]).To(Not(BeNil()))
	Expect(clients[DeveloperEndpoint]).To(Not(BeNil()))

	//put back the env values
	for k, v := range currentEnvVals {
		os.Setenv(k, v)
	}
}

func TestExecuteCommand(t *testing.T) {
	RegisterTestingT(t)

	execFuncExecuted := false
	initFuncExecuted := false
	execFuncExecutedOnDeveloperEndpoint := false
	commands := []Command{
		{
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
			ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
				execFuncExecuted = true
				execFuncExecutedOnDeveloperEndpoint = client.GetEndpoint() == "developer"
				return "", nil
			},
		},
	}

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	client.EXPECT().GetEndpoint().Return("user").AnyTimes()
	clients := map[string]metalcloud.MetalCloudClient{
		UserEndpoint: client,
		"":           client,
	}
	//check with wrong commands first, should return err
	err := executeCommand([]string{"", "test", "test"}, commands, clients)
	Expect(err).NotTo(BeNil())

	execFuncExecuted = false
	initFuncExecuted = false

	//should execute stuff help and not return error
	err = executeCommand([]string{"", "s", "p"}, commands, clients)
	Expect(err).To(BeNil())
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())

	execFuncExecuted = false
	initFuncExecuted = false

	//should execute stuff help and not return error
	err = executeCommand([]string{"", "tests", "testp"}, commands, clients)
	Expect(err).To(BeNil())
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())
	Expect(execFuncExecutedOnDeveloperEndpoint).To(BeFalse())

	//should refuse to execute call on unset endpoint
	commands[0].Endpoint = DeveloperEndpoint
	err = executeCommand([]string{"", "tests", "testp"}, commands, clients)
	Expect(err).NotTo(BeNil())

	//check with correct endpoint
	devClient := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	devClient.EXPECT().GetEndpoint().Return("developer").Times(1)

	//should execute the call if endoint set, on the right endpoint
	clients[DeveloperEndpoint] = devClient

	err = executeCommand([]string{"", "tests", "testp"}, commands, clients)
	Expect(err).To(BeNil())
	Expect(execFuncExecutedOnDeveloperEndpoint).To(BeTrue())

	//should show list of possible predicates if correct subject provided
	err = executeCommand([]string{"", "tests"}, commands, clients)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("testp"))
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())

	execFuncExecuted = false
	initFuncExecuted = false

	//should not show list of possible predicates if correct subject provided
	// but subject has nil predicate
	commands[0].Predicate = _nilDefaultStr
	devClient.EXPECT().GetEndpoint().Return("developer").Times(1)
	err = executeCommand([]string{"", "tests"}, commands, clients)
	Expect(err).To(BeNil())
	Expect(execFuncExecuted).To(BeTrue())
	Expect(initFuncExecuted).To(BeTrue())

	execFuncExecuted = false
	initFuncExecuted = false

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
		ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", nil
		}}

	cmd.InitFunc(&cmd)
	s := getCommandHelp(cmd, true)
	Expect(s).To(ContainSubstring(cmd.Description))
	Expect(s).To(ContainSubstring("Random param"))

}

func TestGetHelp(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	clients := map[string]metalcloud.MetalCloudClient{
		"": client,
	}
	cmds := getCommands(clients)

	s := getHelp(clients, true)
	for _, c := range cmds {
		Expect(s).To(ContainSubstring(c.Description))
	}
}

func TestRequestInputString(t *testing.T) {
	RegisterTestingT(t)
	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)

	stdin.WriteString("test")

	//check without autoconfirm
	ret, err := requestInputString("test")
	Expect(ret).To(Equal("test"))
	Expect(err).To(BeNil())
}

func TestRequestInput(t *testing.T) {
	RegisterTestingT(t)
	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)

	bytes := []byte{13, 100, 20}
	stdin.Write(bytes)

	//check without autoconfirm
	ret, err := requestInput("test")
	Expect(ret).To(Equal(bytes))
	Expect(err).To(BeNil())
}

func TestRequestConfirmation(t *testing.T) {
	RegisterTestingT(t)
	var stdin bytes.Buffer
	var stdout bytes.Buffer

	SetConsoleIOChannel(&stdin, &stdout)

	stdin.WriteString("yes\n")

	//check without autoconfirm
	ok, err := requestConfirmation("test")
	Expect(ok).To(BeTrue())
	Expect(err).To(BeNil())

}
