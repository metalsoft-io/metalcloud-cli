package firmware

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/metalsoft-io/metalcloud-cli/internal/networking"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	kh "golang.org/x/crypto/ssh/knownhosts"
)

/**
 * Firmware related environment variables:
 	METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME
	METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT
	METALCLOUD_FIRMWARE_REPOSITORY_ISO_PATH

 * SCP related environment variables:
	METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH
	METALCLOUD_KNOWN_HOSTS_FILE_PATH
*/

const (
	repositoryURL = "https://repo.metalsoft.com/firmware/"

	configFormatJSON          = "json"
	configFormatYAML          = "yaml"
	catalogVendorDell         = "Dell"
	catalogVendorLenovo       = "Lenovo"
	catalogVendorHp           = "Hp"
	catalogUpdateTypeOffline  = "offline"
	catalogUpdateTypeOnline   = "online"
	stringMinimumSize         = 1
	stringMaximumSize         = 255
	updateSeverityUnknown     = "unknown"
	updateSeverityRecommended = "recommended"
	updateSeverityCritical    = "critical"
	updateSeverityOptional    = "optional"

	defaultFirmwareRepositorySSHPath  = "/var/www/html/firmware"
	defaultFirmwareRepositorySSHPort  = "22"
	defaultFirmwareRepositoryHostname = "192.168.20.10"
	defaultRepositoryFirmwarePath     = "/firmware"
)

type serverInfo struct {
	MachineType  string `json:"machineType" yaml:"machine_type"`
	SerialNumber string `json:"serialNumber" yaml:"serial_number"`
}

type rawConfigFile struct {
	Name            string       `json:"name" yaml:"name"`
	Description     string       `json:"description" yaml:"description"`
	Vendor          string       `json:"vendor" yaml:"vendor"`
	CatalogUrl      string       `json:"catalogUrl" yaml:"catalog_url"`
	DownloadCatalog bool         `json:"downloadCatalog" yaml:"download_catalog"`
	CatalogPath     string       `json:"catalogPath" yaml:"catalog_path"`
	ServersList     []serverInfo `json:"serversList" yaml:"servers_list"`
	DownloadPath    string       `json:"downloadPath" yaml:"download_path"`
}

type firmwareCatalog struct {
	Name                   string
	Description            string
	Vendor                 string
	VendorID               string
	VendorURL              string
	VendorReleaseTimestamp string
	UpdateType             string
	ServerTypesSupported   []string
	Configuration          map[string]string
	CreatedTimestamp       string
}

type firmwareBinary struct {
	ExternalId             string
	Name                   string
	Description            string
	PackageId              string
	PackageVersion         string
	RebootRequired         bool
	UpdateSeverity         string
	SupportedDevices       []map[string]string
	SupportedSystems       []map[string]string
	VendorProperties       map[string]string
	VendorReleaseTimestamp string
	CreatedTimestamp       string
	DownloadURL            string
	RepoURL                string
	LocalPath              string
}

func parseConfigFile(configFormat string, rawConfigFileContents []byte, configFile *rawConfigFile, downloadBinaries bool) error {
	switch configFormat {
	case configFormatJSON:
		err := json.Unmarshal(rawConfigFileContents, &configFile)

		if err != nil {
			return fmt.Errorf("error parsing json file %s: %s", rawConfigFileContents, err.Error())
		}

	case configFormatYAML:
		err := yaml.Unmarshal(rawConfigFileContents, &configFile)

		if err != nil {
			return fmt.Errorf("error parsing yaml file %s: %s", rawConfigFileContents, err.Error())
		}

	default:
		validFormats := []string{configFormatJSON, configFormatYAML}
		return fmt.Errorf("the 'config-format' parameter must be one of %v", validFormats)
	}

	structValue := reflect.ValueOf(configFile).Elem()
	fieldNum := structValue.NumField()

	optionalFields := []string{"CatalogUrl", "ServersList", "DownloadPath"}

	for i := 0; i < fieldNum; i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if !isSet && !slices.Contains[string](optionalFields, fieldName) {
			return fmt.Errorf("the '%s' field must be set in the raw-config file", fieldName)
		}
	}

	if configFile.DownloadCatalog && configFile.CatalogUrl == "" {
		if configFormat == configFormatJSON {
			return fmt.Errorf("the 'catalogUrl' field must be set in the raw-config file when 'downloadCatalog' is set to true")
		}

		if configFormat == configFormatYAML {
			return fmt.Errorf("the 'catalog_url' field must be set in the raw-config file when 'download_catalog' is set to true")
		}
	}

	if configFile.Vendor == catalogVendorLenovo && configFile.ServersList == nil {
		if configFormat == configFormatJSON {
			return fmt.Errorf("the 'serversList' field must be set in the raw-config file when 'vendor' is set to '%s'", catalogVendorLenovo)
		}

		if configFormat == configFormatYAML {
			return fmt.Errorf("the 'servers_list' field must be set in the raw-config file when 'vendor' is set to '%s'", catalogVendorLenovo)
		}
	}

	if downloadBinaries && (configFile.DownloadPath == "" || configFile.CatalogUrl == "") {
		if configFormat == configFormatJSON {
			return fmt.Errorf("the 'downloadPath' and 'catalogUrl' fields must be set in the raw-config file when downloading binaries")
		}

		if configFormat == configFormatYAML {
			return fmt.Errorf("the 'download_path' and 'catalog_url' fields must be set in the raw-config file when downloading binaries")
		}
	}

	checkStringSize(configFile.Name)
	checkStringSize(configFile.Description)
	checkStringSize(configFile.CatalogUrl)

	return nil
}

func checkStringSize(str string) error {
	if len(str) < stringMinimumSize {
		return fmt.Errorf("the '%s' field must be at least %d characters", str, stringMinimumSize)
	}

	if len(str) > stringMaximumSize {
		return fmt.Errorf("the '%s' field must be less than %d characters", str, stringMaximumSize)
	}

	return nil
}

func getUpdateType(rawConfigFile rawConfigFile) string {
	if rawConfigFile.DownloadCatalog {
		return catalogUpdateTypeOnline
	}

	return catalogUpdateTypeOffline
}

func getSeverity(input string) (string, error) {
	switch input {
	case "0":
		return updateSeverityUnknown, nil
	case "1":
		return updateSeverityRecommended, nil
	case "2":
		return updateSeverityCritical, nil
	case "3":
		return updateSeverityOptional, nil
	default:
		return "", fmt.Errorf("invalid severity value: %s", input)
	}
}

func getFirmwareRepositoryURL() string {
	return "https://" + getFirmwareRepositoryHostname() + getFirmwareRepositoryPath()
}

func getFirmwareRepositoryHostname() string {
	firmwareRepositoryHostname := defaultFirmwareRepositoryHostname

	if userGivenFirmwareRepositoryHostname := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME"); userGivenFirmwareRepositoryHostname != "" {
		firmwareRepositoryHostname = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_HOSTNAME")
	}

	return firmwareRepositoryHostname
}

func getFirmwareRepositoryPath() string {
	firmwareRepositoryPath := defaultRepositoryFirmwarePath

	if userGivenFirmwarePath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_ISO_PATH"); userGivenFirmwarePath != "" {
		firmwareRepositoryPath = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_ISO_PATH")

		if !strings.HasPrefix(firmwareRepositoryPath, "/") {
			firmwareRepositoryPath = "/" + firmwareRepositoryPath
		}
	}

	return firmwareRepositoryPath
}

func getFirmwareRepositorySSHPath() string {
	remoteDirectoryPath := defaultFirmwareRepositorySSHPath

	if userGivenRemoteDirectoryPath := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH"); userGivenRemoteDirectoryPath != "" {
		remoteDirectoryPath = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH")
	}

	return remoteDirectoryPath
}

func getFirmwareRepositorySSHPort() string {
	remoteSSHPort := defaultFirmwareRepositorySSHPort

	if userGivenSSHPort := os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT"); userGivenSSHPort != "" {
		remoteSSHPort = os.Getenv("METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT")
	}

	return remoteSSHPort
}

func getUserPrivateSSHKeyPath() (string, error) {
	if userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH"); userPrivateSSHKeyPath == "" {
		return "", fmt.Errorf("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH must be set when creating a firmware binary. The key is needed when uploading to the firmware binary repository.")
	}

	userPrivateSSHKeyPath := os.Getenv("METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH")

	return userPrivateSSHKeyPath, nil
}

func getKnownHostsPath() string {
	var knownHostsFilePath string

	if userGivenHostsFilePath := os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH"); userGivenHostsFilePath != "" {
		knownHostsFilePath = os.Getenv("METALCLOUD_KNOWN_HOSTS_FILE_PATH")
	}

	return knownHostsFilePath
}

func downloadBinariesFromCatalog(binaryCollection []firmwareBinary) error {
	fmt.Println("Downloading binaries.")

	for _, firmwareBinary := range binaryCollection {
		if !networking.CheckValidUrl(firmwareBinary.DownloadURL) {
			return fmt.Errorf("download URL '%s' is not valid.", firmwareBinary.DownloadURL)
		}

		err := networking.DownloadFile(firmwareBinary.DownloadURL, firmwareBinary.LocalPath)

		if err != nil {
			return fmt.Errorf("error while downloading file: %s", err)
		}

		fmt.Printf("Downloaded binary '%s' from URL '%s' to path '%s'.\n", filepath.Base(firmwareBinary.DownloadURL), firmwareBinary.DownloadURL, firmwareBinary.LocalPath)
	}

	fmt.Println("Finished downloading binaries.")
	return nil
}

func uploadBinariesToRepository(binaryCollection []firmwareBinary, replaceIfExists, strictHostKeyChecking, downloadBinaries bool) error {
	//(c.Arguments["replace_if_exists"])
	//"strict_host_key_checking"

	for _, firmwareBinary := range binaryCollection {
		handleFirmwareBinaryUpload(firmwareBinary, replaceIfExists, strictHostKeyChecking, downloadBinaries)
	}

	return nil
}

func getFirmwareBinaryName(binary firmwareBinary) string {
	firmwareBinaryNameArr := strings.Split(binary.LocalPath, "/")
	firmwareBinaryName := firmwareBinaryNameArr[len(firmwareBinaryNameArr)-1]

	return firmwareBinaryName
}

func handleFirmwareBinaryUpload(binary firmwareBinary, replaceIfExists, strictHostKeyChecking, downloadBinaries bool) (string, error) {
	var firmwareBinaryPath string

	if downloadBinaries {
		firmwareBinaryPath = binary.LocalPath
	} else {
		return "", fmt.Errorf("Unsupported for the moment")
	}

	firmwareBinaryFilename := getFirmwareBinaryName(binary)
	remotePath := getFirmwareRepositorySSHPath() + firmwareBinaryFilename

	firmwareBinaryRepositoryHostname := getFirmwareRepositoryHostname()
	firmwareRepositoryPath := getFirmwareRepositoryPath()

	remoteURL := "https://" + firmwareBinaryRepositoryHostname + firmwareRepositoryPath
	sshRepositoryHostname := firmwareBinaryRepositoryHostname + ":" + getFirmwareRepositorySSHPort()

	firmwareBinaryExists, err := networking.CheckRemoteFileExists(remoteURL, firmwareBinaryFilename)

	if err != nil {
		return "", err
	}

	if firmwareBinaryExists && !replaceIfExists {
		fmt.Printf("Firmware binary %s already exists at path %s. Skipping upload. Use the 'replace-if-exists' parameter to replace the existing firmware binary.\n", firmwareBinaryFilename, remotePath)
		return "", nil
	}

	if firmwareBinaryExists {
		fmt.Printf("Replacing firmware binary %s at path %s.\n", firmwareBinaryFilename, remotePath)
	} else {
		fmt.Printf("Uploading new firmware binary %s at path %s.\n", firmwareBinaryFilename, remotePath)
	}

	userPrivateSSHKeyPath, err := getUserPrivateSSHKeyPath()

	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	knownHostsFilePath := getKnownHostsPath()

	if knownHostsFilePath == "" {
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

	hostKeyCallback, err := kh.New(knownHostsFilePath)

	if err != nil {
		return "", fmt.Errorf("Received following error when parsing the known_hosts file: %s.", err)
	}

	// Use SSH key authentication from the auth package.
	clientConfig, err := auth.PrivateKey(
		"root",
		userPrivateSSHKeyPath,
		ssh.HostKeyCallback(func(hostname string, remoteAddress net.Addr, publicKey ssh.PublicKey) error {
			var keyError *kh.KeyError
			hostsError := hostKeyCallback(hostname, remoteAddress, publicKey)

			// Reference: https://www.godoc.org/golang.org/x/crypto/ssh/knownhosts#KeyError
			//if keyErr.Want is not empty and
			if errors.As(hostsError, &keyError) {
				if len(keyError.Want) > 0 {
					// If host is known then there is key mismatch and the connection is rejected.
					fmt.Printf(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The key sent by the remote host is
%s.
Please contact your system administrator.
Add correct host key in %s to get rid of this message.
Host key for %s has changed and you have requested strict checking.
Host key verification failed.
`,
						networking.SerializeSSHKey(publicKey), knownHostsFilePath, hostname,
					)
					return keyError
				} else {
					// If keyErr.Want slice is empty then host is unknown.
					fmt.Printf(`
The authenticity of host '%s' can't be established.
SSH key is %s.
This key is not known by any other names.
It will be added to known_hosts file %s.
Are you sure you want to continue connecting (yes/no)?
`,
						hostname, networking.SerializeSSHKey(publicKey), knownHostsFilePath,
					)

					if strictHostKeyChecking {
						reader := bufio.NewReader(os.Stdin)
						input, err := reader.ReadString('\n')

						if err != nil {
							return err
						}

						// Remove \r and \n from input
						input = string(bytes.TrimSuffix([]byte(input), []byte("\r\n")))

						if input != "yes" {
							if input == "no" {
								fmt.Println("Aborting connection.")
							} else {
								fmt.Println("Invalid response given. Aborting connection.")
							}

							return keyError
						}
					} else {
						fmt.Printf("Skipped manual check because 'strict-host-key-checking' is set to false.")
					}

					return networking.AddHostKey(knownHostsFilePath, remoteAddress, publicKey)
				}
			}

			fmt.Printf("Public key exists for remote %s. Establishing connection.\n", hostname)
			return nil
		}),
	)

	if err != nil {
		return "", fmt.Errorf("Could not create SSH client config. Received error: %s", err)
	}

	// Create a new SCP client.
	scpClient := scp.NewClient(sshRepositoryHostname, &clientConfig)

	// Connect to the remote server.
	err = scpClient.Connect()
	if err != nil {
		return "", fmt.Errorf("Couldn't establish a connection to the remote server: %s", err)
	}

	defer scpClient.Close()

	fmt.Printf("Established connection to hostname %s.\n", sshRepositoryHostname)

	firmwareBinaryFile, err := os.Open(firmwareBinaryPath)
	if err != nil {
		return "", fmt.Errorf("File not found at path %s.", firmwareBinaryPath)
	}
	defer firmwareBinaryFile.Close()

	fmt.Printf("Starting file upload to repository at path %s.\n", remotePath)
	err = scpClient.CopyFile(context.Background(), firmwareBinaryFile, remotePath, "0777")

	if err != nil {
		return "", fmt.Errorf("Error while copying file: %s", err)
	}

	fmt.Printf("Finished file upload to repository at path %s.\n", remotePath)

	return "", nil
}

// TODO: this function should send the catalog to the gateway microservice
func sendCatalog(catalog firmwareCatalog) error {
	catalogJSON, err := json.MarshalIndent(catalog, "", " ")

	if err != nil {
		return fmt.Errorf("Error while marshalling catalog to JSON: %s", err)
	}

	fmt.Printf("Created catalog: %+v\n", string(catalogJSON))

	return nil
}

// TODO: this function should send the binaries to the gateway microservice
func sendBinaries(binaryCollection []firmwareBinary) error {
	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryJson, err := json.MarshalIndent(firmwareBinary, "", " ")

		if err != nil {
			return fmt.Errorf("Error while marshalling binary to JSON: %s", err)
		}

		fmt.Printf("Created firmware binary: %v\n", string(firmwareBinaryJson))
	}

	return nil
}
