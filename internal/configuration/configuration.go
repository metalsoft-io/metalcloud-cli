package configuration

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	Version string
	Date    string
	Commit  string
	BuiltBy string
)

func IsAdmin() bool {
	return os.Getenv("METALCLOUD_ADMIN") == "true"
}

func ReadInputFromPipe() ([]byte, error) {

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

func ReadInputFromFile(path string) ([]byte, error) {

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

// ConsoleIOChannel represents an IO channel, typically stdin and stdout but could be anything
type ConsoleIOChannel struct {
	Stdin  io.Reader
	Stdout io.Writer
}

var consoleIOChannelInstance ConsoleIOChannel

var once sync.Once

// GetConsoleIOChannel returns the console channel singleton
func GetConsoleIOChannel() *ConsoleIOChannel {
	once.Do(func() {

		consoleIOChannelInstance = ConsoleIOChannel{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
		}
	})

	return &consoleIOChannelInstance
}

// GetStdout returns the configured output channel
func GetStdout() io.Writer {
	return GetConsoleIOChannel().Stdout
}

// GetStdin returns the configured input channel
func GetStdin() io.Reader {
	return GetConsoleIOChannel().Stdin
}

// SetConsoleIOChannel configures the stdin and stdout to be used by all io with
func SetConsoleIOChannel(in io.Reader, out io.Writer) {
	channel := GetConsoleIOChannel()
	channel.Stdin = in
	channel.Stdout = out
}

// GetUserEmail returns the API key's owner
func GetUserEmail() string {
	return os.Getenv("METALCLOUD_USER_EMAIL")
}

func GetFirmwareRepositoryURL() (string, error) {
	hostname, err := GetFirmwareRepositoryHostname()
	if err != nil {
		return "", err
	}

	repositoryPath, err := GetFirmwareRepositoryPath()
	if err != nil {
		return "", err
	}

	return "https://" + hostname + repositoryPath, nil
}

func GetFirmwareRepositoryHostname() (string, error) {
	if userGivenFirmwareRepositoryHostname := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME"); userGivenFirmwareRepositoryHostname == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME must be set when uploading a firmware binary.")
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME"), nil
}

func GetFirmwareRepositoryPath() (string, error) {
	var firmwareRepositoryPath string

	if userGivenFirmwarePath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_PATH"); userGivenFirmwarePath == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_PATH must be set when uploading a firmware binary.")
	}

	firmwareRepositoryPath = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_PATH")

	if !strings.HasPrefix(firmwareRepositoryPath, "/") {
		firmwareRepositoryPath = "/" + firmwareRepositoryPath
	}
	
	return firmwareRepositoryPath, nil
}

func GetFirmwareRepositorySSHPath() (string, error) {
	if userGivenRemoteDirectoryPath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"); userGivenRemoteDirectoryPath == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH must be set when uploading a firmware binary.")
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"), nil
}

func GetFirmwareRepositorySSHPort() (string, error) {
	if userGivenSSHPort := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT"); userGivenSSHPort == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT must be set when uploading a firmware binary.")
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT"), nil
}

func GetUserPrivateSSHKeyPath() (string, error) {
	if userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"); userPrivateSSHKeyPath == "" {
		return "", fmt.Errorf("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH must be set when creating a firmware binary. The key is needed when uploading to the firmware binary repository.")
	}

	return os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"), nil
}

func GetKnownHostsPath() string {
	var knownHostsFilePath string

	if userGivenHostsFilePath := os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH"); userGivenHostsFilePath != "" {
		knownHostsFilePath = os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH")
	}

	return knownHostsFilePath
}
