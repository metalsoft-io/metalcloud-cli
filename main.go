package main

// to regenerate the interface and mocking object for the underlying sdk client run 'go generate'. Make sure you have pulled or used go get on the sdk

//go:generate mockgen -source=$GOPATH/src/github.com/bigstepinc/metal-cloud-sdk-go//metal_cloud_client.go -destination=helpers/mock_client.go

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go/v2"

	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	version string
	date    string
	commit  string
	builtBy string
)

//UserEndpoint exposes regular user functions
const UserEndpoint = "user"

//ExtendedEndpoint exposes power functions
const ExtendedEndpoint = "extended"

//DeveloperEndpoint exposes admin functions
const DeveloperEndpoint = "developer"

//GetUserEmail returns the API key's owner
func GetUserEmail() string {
	return os.Getenv("METALCLOUD_USER_EMAIL")
}

//GetDatacenter returns the default datacenter
func GetDatacenter() string {
	return os.Getenv("METALCLOUD_DATACENTER")
}

func main() {

	SetConsoleIOChannel(os.Stdin, os.Stdout)

	clients, err := initClients()
	if err != nil {
		fmt.Fprintf(GetStdout(), "Could not initialize metal cloud client %s\n", err)
		os.Exit(-1)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(GetStdout(), "Error: Syntax error. Use %s help for more details.\n", os.Args[0])
		os.Exit(-1)
	}

	if os.Args[1] == "help" {
		fmt.Fprintf(GetStdout(), "%s\n", getHelp(clients, false))
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		fmt.Fprintf(GetStdout(), "Error: Syntax error. Use %s help for more details.\n", os.Args[0])
		os.Exit(-1)
	}

	commands := getCommands(clients)

	err = executeCommand(os.Args, commands, clients)

	if err != nil {
		fmt.Fprintf(GetStdout(), "%s\n", err)
		os.Exit(-2)
	}
}

func validateArguments(args []string) (string, string, int) {
	subject, predicate := _nilDefaultStr, _nilDefaultStr
	count := 1

	if !strings.HasPrefix(args[1], "-") {
		subject = args[1]
		count++
	}
	if !strings.HasPrefix(args[2], "-") {
		predicate = args[2]
		count++
	}

	return subject, predicate, count
}

func helpMessage(err error, subject string, predicate string) error {
	if predicate != _nilDefaultStr {
		return fmt.Errorf("%s Use '%s %s -h' for syntax help", err, subject, predicate)
	}

	return fmt.Errorf("%s Use '%s -h' for syntax help", err, subject)
}

func commandHelp(cmd Command) string {
	var c string
	if cmd.Predicate != _nilDefaultStr {
		c = fmt.Sprintf("%s %s", cmd.Subject, cmd.Predicate)
		return fmt.Sprintf("Command: %-40s %s (alternatively use \"%s %s\")\n", c, cmd.Description, cmd.AltSubject, cmd.AltPredicate)
	}

	c = fmt.Sprintf("%s", cmd.Subject)
	return fmt.Sprintf("Command: %-40s %s \n", c, cmd.Description)
}

func executeCommand(args []string, commands []Command, clients map[string]metalcloud.MetalCloudClient) error {
	subject, predicate, count := validateArguments(args)

	cmd := locateCommand(predicate, subject, commands)

	if cmd == nil {
		return fmt.Errorf("%s %s is not a valid command. Use %s help for more details", subject, predicate, args[0])
	}

	cmd.InitFunc(cmd)

	//disable default usage
	cmd.FlagSet.Usage = func() {}

	for _, a := range args {
		if a == "-h" || a == "-help" || a == "--help" {
			return fmt.Errorf(getCommandHelp(*cmd, true))
		}
	}

	err := cmd.FlagSet.Parse(args[count:])

	if err != nil {
		return helpMessage(err, subject, predicate)
	}

	client, ok := clients[cmd.Endpoint]
	if !ok {
		return fmt.Errorf("Client not set for endpoint %s on command %s %s", cmd.Endpoint, subject, predicate)
	}

	ret, err := cmd.ExecuteFunc(cmd, client)
	if err != nil {
		return helpMessage(err, subject, predicate)
	}

	fmt.Fprintf(GetStdout(), ret)

	return nil
}

//identifies command, returns nil if no matching command found
func locateCommand(predicate string, subject string, commands []Command) *Command {
	for _, c := range commands {
		if (c.Subject == subject || c.AltSubject == subject) &&
			(c.Predicate == predicate || c.AltPredicate == predicate) {
			return &c
		}
	}
	return nil
}

func getArgumentHelp(f *flag.Flag) string {

	if len(f.Name) == 1 {
		return fmt.Sprintf("\t  -%-25s %s\n", f.Name, f.Usage)
	}

	return fmt.Sprintf("\t  --%-24s %s\n", f.Name, f.Usage)

}

func getCommandHelp(cmd Command, showArguments bool) string {
	var sb strings.Builder
	var c string
	if cmd.Predicate != _nilDefaultStr {
		c = fmt.Sprintf("%s %s", cmd.Subject, cmd.Predicate)
	} else {
		c = fmt.Sprintf("%s", cmd.Subject)
	}

	if showArguments {
		sb.WriteString(commandHelp(cmd))
		cmd.FlagSet.VisitAll(func(f *flag.Flag) {
			sb.WriteString(getArgumentHelp(f))
		})

		h := flag.Flag{
			Name:  "h",
			Usage: "Show command help and exit.",
		}

		sb.WriteString(getArgumentHelp(&h))
	} else {
		sb.WriteString(fmt.Sprintf("\t%-40s %-24s", c, cmd.Description))
	}

	return sb.String()
}

func getHelp(clients map[string]metalcloud.MetalCloudClient, showArguments bool) string {
	var sb strings.Builder
	cmds := getCommands(clients)
	for _, c := range cmds {
		c.InitFunc(&c)
	}
	sb.WriteString(fmt.Sprintf("Syntax: %s <command> [args]\nAccepted commands:\n", os.Args[0]))
	for _, c := range cmds {
		sb.WriteString(fmt.Sprintln(getCommandHelp(c, false)))
	}
	return sb.String()
}

func isLoggingEnabled() bool {
	return os.Getenv("METALCLOUD_LOGGING_ENABLED") == "true"

}

func isAdmin() bool {
	return os.Getenv("METALCLOUD_ADMIN") == "true"
}

func initClients() (map[string]metalcloud.MetalCloudClient, error) {

	clients := map[string]metalcloud.MetalCloudClient{}
	endpointSuffixes := map[string]string{
		DeveloperEndpoint: "/api/developer/developer",
		ExtendedEndpoint:  "/metal-cloud/extended",
		UserEndpoint:      "/metal-cloud",
		"":                "/metal-cloud",
	}

	for clientName, suffix := range endpointSuffixes {

		if (clientName == DeveloperEndpoint || clientName == ExtendedEndpoint) && !isAdmin() {
			continue
		}

		client, err := initClient(suffix)
		if err != nil {
			return nil, err
		}
		clients[clientName] = client
	}
	return clients, nil
}

func initClient(endpointSuffix string) (metalcloud.MetalCloudClient, error) {
	if v := os.Getenv("METALCLOUD_USER_EMAIL"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_USER_EMAIL must be set")
	}

	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	if v := os.Getenv("METALCLOUD_DATACENTER"); v == "" {
		return nil, fmt.Errorf("METALCLOUD_DATACENTER must be set")
	}

	apiKey := os.Getenv("METALCLOUD_API_KEY")
	user := os.Getenv("METALCLOUD_USER_EMAIL")

	endpointHost := strings.TrimRight(os.Getenv("METALCLOUD_ENDPOINT"), "/")
	endpoint := fmt.Sprintf("%s%s", endpointHost, endpointSuffix)

	loggingEnabled := isLoggingEnabled()

	err := validateAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return metalcloud.GetMetalcloudClient(user, apiKey, endpoint, loggingEnabled)

}

func fitlerCommandSet(commandSet []Command, clients map[string]metalcloud.MetalCloudClient) []Command {
	filteredCommands := []Command{}
	for _, command := range commandSet {
		if _, ok := clients[command.Endpoint]; ok {
			filteredCommands = append(filteredCommands, command)
		}
	}
	return filteredCommands
}

func getCommands(clients map[string]metalcloud.MetalCloudClient) []Command {

	commands := [][]Command{
		datacenterCmds,
		infrastructureCmds,
		instanceArrayCmds,
		instanceCmds,
		driveArrayCmds,
		driveSnapshotCmds,
		volumeTemplateCmds,
		firewallRuleCmds,
		secretsCmds,
		variablesCmds,
		osAssetsCmds,
		osTemplatesCmds,
		serversCmds,
		switchCmds,
		subnetPoolCmds,
		stageDefinitionsCmds,
		workflowCmds,
		versionCmds,
		applyCmds,
	}

	filteredCommands := []Command{}
	for _, commandSet := range commands {
		commands := fitlerCommandSet(commandSet, clients)
		filteredCommands = append(filteredCommands, commands...)
	}

	return filteredCommands
}

func validateAPIKey(apiKey string) error {
	const pattern = "^\\d+\\:[0-9a-zA-Z]*$"

	matched, _ := regexp.MatchString(pattern, apiKey)

	if !matched {
		return fmt.Errorf("API Key is not valid. It should start with a number followed by a semicolon followed by alphanumeric characters <id>:<chars> ")
	}

	return nil
}

func readInputFromPipe() ([]byte, error) {

	reader := bufio.NewReader(GetStdin())
	var content []byte

	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		content = append(content, input)
	}

	return content, nil
}

func requestInputSilent(s string) ([]byte, error) {

	fmt.Fprintf(GetStdout(), s)
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

func requestInput(s string) ([]byte, error) {

	fmt.Fprintf(GetStdout(), s)
	reader := bufio.NewReader(GetStdin())
	content, err := reader.ReadBytes('\n')

	if err != nil && err != io.EOF {
		return content, err
	}

	return content, nil
}

func requestInputString(s string) (string, error) {

	fmt.Fprintf(GetStdout(), s)
	reader := bufio.NewReader(GetStdin())
	content, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return content, err
	}

	return content, nil
}

func requestConfirmation(s string) (bool, error) {
	yes, err := requestInput(s)
	if err != nil {
		return false, err
	}

	return strings.Trim(string(yes), "\r\n ") == "yes", nil
}

func readInputFromFile(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//ConsoleIOChannel represents an IO channel, typically stdin and stdout but could be anything
type ConsoleIOChannel struct {
	Stdin  io.Reader
	Stdout io.Writer
}

var consoleIOChannelInstance ConsoleIOChannel

var once sync.Once

//GetConsoleIOChannel returns the console channel singleton
func GetConsoleIOChannel() *ConsoleIOChannel {
	once.Do(func() {

		consoleIOChannelInstance = ConsoleIOChannel{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
		}
	})

	return &consoleIOChannelInstance
}

//GetStdout returns the configured output channel
func GetStdout() io.Writer {
	return GetConsoleIOChannel().Stdout
}

//GetStdin returns the configured input channel
func GetStdin() io.Reader {
	return GetConsoleIOChannel().Stdin
}

//SetConsoleIOChannel configures the stdin and stdout to be used by all io with
func SetConsoleIOChannel(in io.Reader, out io.Writer) {
	channel := GetConsoleIOChannel()
	channel.Stdin = in
	channel.Stdout = out
}
