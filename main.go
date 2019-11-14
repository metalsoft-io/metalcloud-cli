package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//GetUserEmail returns the API key's owner
func GetUserEmail() string {
	return os.Getenv("METALCLOUD_USER_EMAIL")
}

//GetDatacenter returns the default datacenter
func GetDatacenter() string {
	return os.Getenv("METALCLOUD_DATACENTER")
}

func main() {

	commands := getCommands()

	if len(os.Args) < 2 {
		fmt.Printf("Error: Syntax error. Use %s help for more details.\n", os.Args[0])
		os.Exit(-1)
	}

	if os.Args[1] == "help" {
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		fmt.Printf("Error: Syntax error. Use %s help for more details.\n", os.Args[0])
		os.Exit(-1)
	}

	client, err := initClient()
	if err != nil {
		fmt.Printf("Could not initialize metal cloud client %s\n", err)
		os.Exit(-1)
	}

	err = executeCommand(os.Args, commands, client)

	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}
}

func executeCommand(args []string, commands []Command, client MetalCloudClient) error {
	predicate := args[1]
	subject := args[2]

	commandExecuted := false
	for _, c := range commands {
		c.InitFunc(&c)

		if (c.Subject == subject || c.AltSubject == subject) &&
			(c.Predicate == predicate || c.AltPredicate == predicate) {

			if len(args) == 4 && args[3] == "-h" {
				commandExecuted = true
				return fmt.Errorf(getCommandHelp(c))
			}

			c.FlagSet.Parse(args[3:])

			ret, err := c.ExecuteFunc(&c, client)
			if err != nil {
				return fmt.Errorf("%s Use '%s %s -h' for syntax help", err, predicate, subject)
			}

			fmt.Print(ret)

			commandExecuted = true
			break
		}
	}

	if !commandExecuted {
		return fmt.Errorf("%s %s is not a valid command. Use %s help for more details", predicate, subject, args[0])
	}

	return nil
}

func getCommandHelp(cmd Command) string {
	var sb strings.Builder

	c := fmt.Sprintf("%s %s", cmd.Predicate, cmd.Subject)
	sb.WriteString(fmt.Sprintf("Command: %-25s %s (alternatively use \"%s %s\")\n", c, cmd.Description, cmd.AltPredicate, cmd.AltSubject))
	cmd.FlagSet.VisitAll(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("\t  -%-25s %s\n", f.Name, f.Usage))
	})

	return sb.String()
}

func printHelp() {
	cmds := getCommands()
	for _, c := range cmds {
		c.InitFunc(&c)
	}
	fmt.Printf("Syntax: %s <command> [args]\nAccepted commands:\n", os.Args[0])
	for _, c := range cmds {
		fmt.Println(getCommandHelp(c))
	}
}

func initClient() (MetalCloudClient, error) {
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
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")
	loggingEnabled := os.Getenv("METALCLOUD_LOGGING_ENABLED") == "true"

	err := validateAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return metalcloud.GetMetalcloudClient(user, apiKey, endpoint, loggingEnabled)
}

func getCommands() []Command {
	var commands []Command

	commands = append(commands, infrastructureCmds...)
	commands = append(commands, instanceArrayCmds...)
	commands = append(commands, driveArrayCmds...)
	commands = append(commands, volumeTemplateyCmds...)
	commands = append(commands, firewallRuleCmds...)

	return commands
}

func validateAPIKey(apiKey string) error {
	const pattern = "^\\d+\\:[0-9a-zA-Z]{63}$"

	matched, _ := regexp.MatchString(pattern, apiKey)

	if !matched {
		return fmt.Errorf("API Key is not valid. It should start with a number followed by a semicolon and 63 alphanumeric characters <id>:<63 chars> ")
	}

	return nil
}

func requestConfirmation(s string) bool {

	fmt.Printf(s)
	reader := bufio.NewReader(os.Stdin)
	yes, _ := reader.ReadString('\n')

	return yes == "yes\n"
}
