package firmware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/networking"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

/**
 * Firmware related environment variables (all required):
 	METALCLOUD_FIRMWARE_REPOSITORY_URL		- the URL of the HTTP repository, for example: http://192.168.20.10/firmware
	METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH - the path to the SSH repository, for example: /var/www/html/firmware
	METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT	- the port of the SSH repository, for example: 22

 * SCP related environment variables:
	METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH (required) 						- the path to the private OpenSSH key used for authentication, for example: ~/.ssh/my-openssh-key
	METALCLOUD_KNOWN_HOSTS_FILE_PATH (optional, defaults to ~/.ssh/known_hosts) - the path to the known_hosts file, for example: ~/.ssh/known_hosts
*/

const (
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
	serverTypesAll            = "*"
)

type serverInfo struct {
	MachineType  string `json:"machineType" yaml:"machine_type"`
	SerialNumber string `json:"serialNumber" yaml:"serial_number"`
	VendorSkuId  string
}

type rawConfigFile struct {
	Name              string       `json:"name" yaml:"name"`
	Description       string       `json:"description" yaml:"description"`
	Vendor            string       `json:"vendor" yaml:"vendor"`
	CatalogUrl        string       `json:"catalogUrl" yaml:"catalog_url"`
	ServersList       []serverInfo `json:"serversList" yaml:"servers_list"`
	LocalCatalogPath  string       `json:"localCatalogPath" yaml:"local_catalog_path"`
	OverwriteCatalogs bool         `json:"overwriteCatalogs" yaml:"overwrite_catalogs"`
	LocalFirmwarePath string       `json:"localFirmwarePath" yaml:"local_firmware_path"`
}

type firmwareCatalog struct {
	Name                          string
	Description                   string
	Vendor                        string
	VendorID                      string
	VendorURL                     string
	VendorReleaseTimestamp        string
	UpdateType                    string
	MetalSoftServerTypesSupported []string
	ServerTypesSupported          []string
	Configuration                 map[string]string
	CreatedTimestamp              string
}

type firmwareBinary struct {
	ExternalId             string
	Name                   string
	FileName               string
	Description            string
	PackageId              string
	PackageVersion         string
	RebootRequired         bool
	UpdateSeverity         string
	Hash                   string
	HashingAlgorithm       string
	SupportedDevices       []map[string]string
	SupportedSystems       []map[string]string
	VendorProperties       map[string]string
	VendorReleaseTimestamp string
	CreatedTimestamp       string
	DownloadURL            string
	RepoURL                string
	LocalPath              string
	HasErrors              bool
	ErrorMessage           string
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

	optionalFields := []string{"CatalogUrl", "ServersList", "LocalFirmwarePath"}

	for i := 0; i < fieldNum; i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name

		isSet := field.IsValid() && (!field.IsZero() || field.Kind() == reflect.Bool)

		if !isSet && !slices.Contains[string](optionalFields, fieldName) {
			value, err := configFileParameterToValue(fieldName, configFormat)
			if err != nil {
				return err
			}

			return fmt.Errorf("the '%s' field must be set in the raw-config file", value)
		}
	}

	if downloadBinaries {
		switch configFile.Vendor {
		case catalogVendorDell:
			if configFile.LocalFirmwarePath == "" || configFile.CatalogUrl == "" {
				firmwarePathValue, catalogUrlValue := "localFirmwarePath", "catalogUrl"

				if configFormat == configFormatYAML {
					firmwarePathValue, catalogUrlValue = "local_firmware_path", "catalog_url"
				}

				return fmt.Errorf("the '%s' and '%s' fields must be set in the raw-config file when downloading Dell binaries", firmwarePathValue, catalogUrlValue)
			}
		case catalogVendorLenovo:
			if configFile.LocalFirmwarePath == "" {
				firmwarePathValue := "localFirmwarePath"

				if configFormat == configFormatYAML {
					firmwarePathValue = "local_firmware_path"
				}

				return fmt.Errorf("the '%s' field must be set in the raw-config file when downloading Lenovo binaries", firmwarePathValue)
			}
		case catalogVendorHp:
			return fmt.Errorf("TODO: HP firmware binaries are not supported yet")
		}
	}

	if configFile.LocalFirmwarePath != "" && !folderExists(configFile.LocalFirmwarePath) {
		value := "localFirmwarePath"

		if configFormat == configFormatYAML {
			value = "local_firmware_path"
		}

		return fmt.Errorf("the '%s' field must be a valid folder path", value)
	}

	if configFile.LocalCatalogPath != "" {
		value := "localCatalogPath"

		if configFormat == configFormatYAML {
			value = "local_catalog_path"
		}

		if configFile.Vendor == catalogVendorDell && !fileExists(configFile.LocalCatalogPath) {
			return fmt.Errorf("the '%s' field must be a valid file path", value)
		}

		if configFile.Vendor == catalogVendorLenovo && !folderExists(configFile.LocalCatalogPath) {
			return fmt.Errorf("the '%s' field must be a valid folder path", value)
		}

	}

	checkStringSize(configFile.Name)
	checkStringSize(configFile.Description)
	checkStringSize(configFile.CatalogUrl)

	return nil
}

func configFileParameterToValue(parameter, format string) (string, error) {
	switch parameter {
	case "Name":
		return "name", nil
	case "Description":
		return "description", nil
	case "Vendor":
		return "vendor", nil
	case "LocalCatalogPath":
		if format == configFormatJSON {
			return "localCatalogPath", nil
		} else if format == configFormatYAML {
			return "local_catalog_path", nil
		}
	case "OverwriteCatalogs":
		if format == configFormatJSON {
			return "overwriteCatalogs", nil
		} else if format == configFormatYAML {
			return "overwrite_catalogs", nil
		}
	default:
		return "", fmt.Errorf("invalid parameter '%s'", parameter)
	}

	return "", nil
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
	if rawConfigFile.Vendor == catalogVendorLenovo {
		return catalogUpdateTypeOnline
	}

	if rawConfigFile.CatalogUrl != "" {
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

func downloadBinariesFromCatalog(binaryCollection []*firmwareBinary) error {
	fmt.Println("Downloading binaries.")

	for _, firmwareBinary := range binaryCollection {
		if !networking.CheckValidUrl(firmwareBinary.DownloadURL) {
			return fmt.Errorf("download URL '%s' is not valid.", firmwareBinary.DownloadURL)
		}

		err := DownloadFirmwareBinary(firmwareBinary)

		if err != nil {
			return err
		}
	}

	fmt.Println("Finished downloading binaries.")
	return nil
}

func uploadBinariesToRepository(binaryCollection []*firmwareBinary, replaceIfExists, skipHostKeyChecking bool) error {
	firmwareRepositoryURL, err := configuration.GetFirmwareRepositoryURL()
	if err != nil {
		return err
	}

	remoteURL, err := url.Parse(firmwareRepositoryURL)
	if err != nil {
		return err
	}

	firmwareBinaryRepositoryHostname := remoteURL.Hostname()

	firmwareRepositorySSHPort, err := configuration.GetFirmwareRepositorySSHPort()
	if err != nil {
		return err
	}

	firmwareRepositorySSHPath, err := configuration.GetFirmwareRepositorySSHPath()
	if err != nil {
		return err
	}

	fmt.Println("Checking if binaries already exist in the repository.")
	firmwareBinaryNames := make([]string, len(binaryCollection))
	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryNames = append(firmwareBinaryNames, firmwareBinary.FileName)
	}

	missingBinaries, err := networking.GetMissingRemoteFiles(remoteURL.String(), firmwareBinaryNames)

	if err != nil {
		return err
	}

	if len(missingBinaries) == 0 && !replaceIfExists {
		fmt.Println("All binaries already exist in the repository. Skipping upload.")
		return nil
	}

	scpClient, sshClient, err := networking.CreateSSHConnection(skipHostKeyChecking)

	if err != nil {
		return fmt.Errorf("Couldn't establish a connection to the remote server: %s", err)
	}

	defer scpClient.Close()

	sshRepositoryHostname := firmwareBinaryRepositoryHostname + ":" + firmwareRepositorySSHPort
	fmt.Printf("Established connection to hostname %s.\n", sshRepositoryHostname)

	if !replaceIfExists {
		fmt.Printf("Detected %d missing binaries.\n", len(missingBinaries))
	} else {
		fmt.Println("The 'replace-if-exists' parameter is set to true. All binaries will be replaced.")
	}

	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryExists := !slices.Contains[string](missingBinaries, firmwareBinary.FileName)

		if firmwareBinaryExists && !replaceIfExists {
			continue
		}

		remotePath := firmwareRepositorySSHPath + "/" + firmwareBinary.FileName
		err := uploadBinaryToRepository(firmwareBinary, &scpClient, sshClient, firmwareBinaryExists, replaceIfExists, remotePath)

		if err != nil {
			return err
		}
	}

	fmt.Println("Finished uploading binaries.")
	return nil
}

func uploadBinaryToRepository(binary *firmwareBinary, scpClient *scp.Client, sshClient *ssh.Client, firmwareBinaryExists, replaceIfExists bool, remotePath string) error {
	// Regenerate the session in the case it was previously closed, otherwise only the first file will be uploaded.
	// TODO: need a check to see if the session is still open
	scpSession, err := sshClient.NewSession()
	if err != nil {
		return err
	}

	scpClient.Session = scpSession
	firmwareBinaryPath := binary.LocalPath

	var firmwareBinaryFile *os.File
	if firmwareBinaryPath == "" {
		if binary.DownloadURL != "" && !binary.HasErrors && (!firmwareBinaryExists || replaceIfExists) {
			// We don't save the binaries on the local filesystem, so we need to download them from the catalog as temporary files and then upload them to the repository.
			firmwareBinaryFile, err = ioutil.TempFile(os.TempDir(), binary.FileName)
			if err != nil {
				return err
			}
			defer os.Remove(firmwareBinaryFile.Name())
			defer firmwareBinaryFile.Close()

			if !networking.CheckValidUrl(binary.DownloadURL) {
				return fmt.Errorf("download URL '%s' is not valid.", binary.DownloadURL)
			}

			binary.LocalPath = firmwareBinaryFile.Name()
			err := DownloadFirmwareBinary(binary)
			binary.LocalPath = ""

			if err != nil {
				return err
			}
		}
	} else {
		firmwareBinaryFile, err = os.Open(firmwareBinaryPath)
		if err != nil {
			return fmt.Errorf("file not found at path %s. Make sure the local firmware path is set correctly in the raw-config file.", firmwareBinaryPath)
		}
		defer firmwareBinaryFile.Close()
	}

	if binary.HasErrors {
		fmt.Printf("Skipping uploading binary %s because it has errors: %s\n", binary.FileName, binary.ErrorMessage)
		binary.RepoURL = ""
		return nil
	}

	if firmwareBinaryExists {
		fmt.Printf("Replacing firmware binary %s at path %s.\n", binary.FileName, remotePath)
	} else {
		fmt.Printf("Uploading new firmware binary %s at path %s.\n", binary.FileName, remotePath)
	}

	err = scpClient.CopyFile(context.Background(), firmwareBinaryFile, remotePath, "0777")

	if err != nil {
		return fmt.Errorf("Error while copying file: %s", err)
	}

	return nil
}

// Returns a map, the key being the actual server type and the value being the Metalsoft internal one.
// Output example: map[PowerEdge R430:[M.32.64.2 M.40.32.1.v2] PowerEdge R730:[M.32.32.2]] [M.32.64.2 M.40.32.1.v2 M.32.32.2]
func retrieveSupportedServerTypes(client metalcloud.MetalCloudClient, input string) (map[string][]string, []string, error) {
	supportedServerTypes := map[string][]string{}
	metalsoftServerTypes := []string{}

	serversTypeObject, err := client.ServerTypes(false)

	if err != nil {
		return nil, nil, fmt.Errorf("Error getting server types: %v", err)
	}

	for _, serverTypeObject := range *serversTypeObject {
		var serverTypes []string
		err := json.Unmarshal([]byte(serverTypeObject.ServerTypeAllowedVendorSkuIdsJSON), &serverTypes)
		if err != nil {
			return nil, nil, fmt.Errorf("Error unmarshalling server types: %v", err)
		}

		supportedServerTypes[serverTypes[0]] = append(supportedServerTypes[serverTypes[0]], serverTypeObject.ServerTypeName)
		metalsoftServerTypes = append(metalsoftServerTypes, serverTypeObject.ServerTypeName)
	}

	if input == "" || input == serverTypesAll {
		return supportedServerTypes, metalsoftServerTypes, nil
	}

	filteredServerTypes := map[string][]string{}
	serverTypesList := strings.Split(input, ",")
	filteredMetalsoftServerTypes := []string{}

	for _, serverType := range serverTypesList {
		serverType = strings.TrimSpace(serverType)
		if !slices.Contains[string](metalsoftServerTypes, serverType) {
			return nil, nil, fmt.Errorf("cannot filter server type '%s' because it is not supported by Metalsoft. Supported types are %+v", serverType, metalsoftServerTypes)
		}

		for actualServerType, metalsoftServerType := range supportedServerTypes {
			if slices.Contains[string](metalsoftServerType, serverType) {
				filteredServerTypes[actualServerType] = append(filteredServerTypes[actualServerType], serverType)
				break
			}
		}

		if !slices.Contains[string](filteredMetalsoftServerTypes, serverType) {
			filteredMetalsoftServerTypes = append(filteredMetalsoftServerTypes, serverType)
		}
	}

	return filteredServerTypes, filteredMetalsoftServerTypes, nil
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
func sendBinaries(binaryCollection []*firmwareBinary) error {
	printOnce := false
	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryJson, err := json.MarshalIndent(firmwareBinary, "", " ")

		if err != nil {
			return fmt.Errorf("Error while marshalling binary to JSON: %s", err)
		}

		if !printOnce {
			fmt.Printf("Created firmware binary: %v\n", string(firmwareBinaryJson))
			printOnce = true
		}
	}

	return nil
}

func DownloadFirmwareBinary(binary *firmwareBinary) error {
	err := networking.DownloadFile(binary.DownloadURL, binary.LocalPath, binary.Hash, binary.HashingAlgorithm)

	if err != nil {
		if err.Error() == fmt.Sprintf("%d", http.StatusNotFound) {
			binary.LocalPath = ""
			binary.HasErrors = true
			binary.ErrorMessage = fmt.Sprintf("Binary not found at URL %s", binary.DownloadURL)

			fmt.Printf("Skipping binary %s: %s\n", binary.FileName, binary.ErrorMessage)
			return nil
		}

		return err
	}

	return nil
}

func folderExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
