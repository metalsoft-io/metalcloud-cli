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

const (
	defaultFirmwareRepositorySSHPath  = "/var/www/html/firmware"
	defaultFirmwareRepositorySSHPort  = "22"
	defaultFirmwareRepositoryHostname = "192.168.20.10"
	defaultRepositoryFirmwarePath     = "/firmware"
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

func GetFirmwareRepositoryURL() string {
	return "https://" + GetFirmwareRepositoryHostname() + GetFirmwareRepositoryPath()
}

func GetFirmwareRepositoryHostname() string {
	firmwareRepositoryHostname := defaultFirmwareRepositoryHostname

	if userGivenFirmwareRepositoryHostname := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME"); userGivenFirmwareRepositoryHostname != "" {
		firmwareRepositoryHostname = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME")
	}

	return firmwareRepositoryHostname
}

func GetFirmwareRepositoryPath() string {
	firmwareRepositoryPath := defaultRepositoryFirmwarePath

	if userGivenFirmwarePath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_ISO_PATH"); userGivenFirmwarePath != "" {
		firmwareRepositoryPath = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_ISO_PATH")

		if !strings.HasPrefix(firmwareRepositoryPath, "/") {
			firmwareRepositoryPath = "/" + firmwareRepositoryPath
		}
	}

	return firmwareRepositoryPath
}

func GetFirmwareRepositorySSHPath() string {
	remoteDirectoryPath := defaultFirmwareRepositorySSHPath

	if userGivenRemoteDirectoryPath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"); userGivenRemoteDirectoryPath != "" {
		remoteDirectoryPath = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH")
	}

	return remoteDirectoryPath
}

func GetFirmwareRepositorySSHPort() string {
	remoteSSHPort := defaultFirmwareRepositorySSHPort

	if userGivenSSHPort := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT"); userGivenSSHPort != "" {
		remoteSSHPort = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT")
	}

	return remoteSSHPort
}

func GetUserPrivateSSHKeyPath() (string, error) {
	if userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"); userPrivateSSHKeyPath == "" {
		return "", fmt.Errorf("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH must be set when creating a firmware binary. The key is needed when uploading to the firmware binary repository.")
	}

	userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH")

	return userPrivateSSHKeyPath, nil
}

func GetKnownHostsPath() string {
	var knownHostsFilePath string

	if userGivenHostsFilePath := os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH"); userGivenHostsFilePath != "" {
		knownHostsFilePath = os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH")
	}

	return knownHostsFilePath
}
