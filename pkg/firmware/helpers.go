package firmware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/networking"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
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
	FileName               string
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
	if !downloadBinaries {
		return fmt.Errorf("Unsupported for the moment")
	}

	firmwareBinaryRepositoryHostname := configuration.GetFirmwareRepositoryHostname()
	firmwareRepositoryPath := configuration.GetFirmwareRepositoryPath()

	//TODO: change this to https
	remoteURL := "http://" + firmwareBinaryRepositoryHostname + firmwareRepositoryPath

	fmt.Println("Checking if binaries already exist in the repository.")
	firmwareBinaryNames := make([]string, len(binaryCollection))
	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryNames = append(firmwareBinaryNames, firmwareBinary.FileName)
	}

	missingBinaries, err := networking.GetMissingRemoteFiles(remoteURL, firmwareBinaryNames)

	if err != nil {
		return err
	}

	if len(missingBinaries) == 0 {
		fmt.Println("All binaries already exist in the repository. Skipping upload.")
		return nil
	}

	scpClient, sshClient, err := networking.CreateSSHConnection(strictHostKeyChecking)

	if err != nil {
		return fmt.Errorf("Couldn't establish a connection to the remote server: %s", err)
	}

	defer scpClient.Close()

	sshRepositoryHostname := configuration.GetFirmwareRepositoryHostname() + ":" + configuration.GetFirmwareRepositorySSHPort()
	fmt.Printf("Established connection to hostname %s.\n", sshRepositoryHostname)

	fmt.Printf("Detected %d missing binaries.\n", len(missingBinaries))
	if replaceIfExists {
		fmt.Println("The 'replace-if-exists' parameter is set to true. All binaries will be replaced.")
	}

	for _, firmwareBinary := range binaryCollection {
		firmwareBinaryExists := !slices.Contains[string](missingBinaries, firmwareBinary.FileName)
		remotePath := configuration.GetFirmwareRepositorySSHPath() + "/" + firmwareBinary.FileName

		if firmwareBinaryExists && !replaceIfExists {
			continue
		}

		if firmwareBinaryExists {
			fmt.Printf("Replacing firmware binary %s at path %s.\n", firmwareBinary.FileName, remotePath)
		} else {
			fmt.Printf("Uploading new firmware binary %s at path %s.\n", firmwareBinary.FileName, remotePath)
		}

		err := uploadBinaryToRepository(firmwareBinary, &scpClient, sshClient)

		if err != nil {
			return err
		}
	}

	fmt.Println("Finished uploading binaries.")
	return nil
}

func uploadBinaryToRepository(binary firmwareBinary, scpClient *scp.Client, sshClient *ssh.Client) error {
	firmwareBinaryPath := binary.LocalPath

	if firmwareBinaryPath == "" {
		return fmt.Errorf("No local path specified for firmware binary %s.", binary.FileName)
	}

	firmwareBinaryFilename := binary.FileName
	remotePath := configuration.GetFirmwareRepositorySSHPath() + "/" + firmwareBinaryFilename

	firmwareBinaryFile, err := os.Open(firmwareBinaryPath)
	if err != nil {
		return fmt.Errorf("File not found at path %s.", firmwareBinaryPath)
	}
	defer firmwareBinaryFile.Close()

	err = scpClient.CopyFile(context.Background(), firmwareBinaryFile, remotePath, "0777")

	if err != nil {
		// Regenerate the session if it was previously closed
		scpSession, err := sshClient.NewSession()
		if err != nil {
			return err
		}
	
		scpClient.Session = scpSession
		err = scpClient.CopyFile(context.Background(), firmwareBinaryFile, remotePath, "0777")

		if err != nil {
			return fmt.Errorf("Error while copying file: %s", err)
		}
	}

	return nil
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
