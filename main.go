package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

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
	predicate := os.Args[1]
	subject := os.Args[2]

	commandExecuted := false
	for _, c := range commands {
		c.InitFunc(&c)

		if (c.Subject == subject || c.AltSubject == subject) &&
			(c.Predicate == predicate || c.AltPredicate == predicate) {

			if len(os.Args) == 4 && os.Args[3] == "-h" {
				commandExecuted = true
				printCommandHelp(c)
				break
			}

			c.FlagSet.Parse(os.Args[3:])

			ret, err := c.ExecuteFunc(&c, client)
			if err != nil {

				fmt.Printf("Error: %s. Use '%s %s -h' for syntax help.\n", err, predicate, subject)
				os.Exit(-2)
			}

			fmt.Print(ret)

			commandExecuted = true
			break
		}
	}

	if !commandExecuted {
		fmt.Printf("Error: %s %s is not a valid command. Use %s help for more details.\n", predicate, subject, os.Args[0])
		os.Exit(-2)
	}
}

func printCommandHelp(cmd Command) {
	fmt.Printf("%s %s - %s (alternatively use \"%s %s\")\n", cmd.Predicate, cmd.Subject, cmd.Description, cmd.AltPredicate, cmd.AltSubject)
	cmd.FlagSet.PrintDefaults()
}

func printHelp() {
	cmds := getCommands()
	for _, c := range cmds {
		c.InitFunc(&c)
	}
	fmt.Printf("Syntax: %s <command> [args]\nAccepted commands:", os.Args[0])
	for _, c := range cmds {
		fmt.Println()
		printCommandHelp(c)
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
