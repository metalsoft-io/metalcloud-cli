package configuration

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	Version string
	Date    string
	Commit  string
	BuiltBy string
)

const (
	defaultSSHPort = "22"
	defaultSSHUser = "root"
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
	if userGivenFirmwareRepositoryHostname := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_URL"); userGivenFirmwareRepositoryHostname == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_URL must be set when uploading firmware binaries.")
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_URL"), nil
}

func GetFirmwareRepositorySSHPath() (string, error) {
	if userGivenRemoteDirectoryPath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"); userGivenRemoteDirectoryPath == "" {
		return "", fmt.Errorf("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH must be set when uploading firmware binaries.")
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"), nil
}

func GetFirmwareRepositorySSHPort() string {
	if userGivenSSHPort := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT"); userGivenSSHPort == "" {
		// If no port is given, use the default SSH port.
		return defaultSSHPort
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT")
}

func GetFirmwareRepositorySSHUser() string {
	if userGivenSSHPort := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_USER"); userGivenSSHPort == "" {
		// If no user is given, use the default SSH user.
		return defaultSSHUser
	}

	return os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_USER")
}

func GetUserPrivateSSHKeyPath() (string, error) {
	if userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"); userPrivateSSHKeyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		defaultPrivateSSHKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa")
		if _, err := os.Stat(defaultPrivateSSHKeyPath); errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH must be set when uploading firmware binaries to the repository. Tried default private key path %s but file does not exist.", defaultPrivateSSHKeyPath)
		}

		return defaultPrivateSSHKeyPath, nil
	}

	return os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"), nil
}

func GetKnownHostsPath() (string, error) {
	var knownHostsFilePath string

	if userGivenHostsFilePath := os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH"); userGivenHostsFilePath != "" {
		knownHostsFilePath = os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		knownHostsFilePath = filepath.Join(homeDir, ".ssh", "known_hosts")

		// Create the known hosts file if it does not exist.
		if _, err := os.Stat(knownHostsFilePath); errors.Is(err, os.ErrNotExist) {
			hostsFile, err := os.Create(knownHostsFilePath)

			if err != nil {
				return "", err
			}

			hostsFile.Close()
		}
	}

	return knownHostsFilePath, nil
}

func GetAPIKey() (string, error) {
	if apiKey := os.Getenv("METALCLOUD_API_KEY"); apiKey == "" {
		return "", fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	return os.Getenv("METALCLOUD_API_KEY"), nil
}

func GetEndpoint() (string, error) {
	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		return "", fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	return os.Getenv("METALCLOUD_ENDPOINT"), nil
}
