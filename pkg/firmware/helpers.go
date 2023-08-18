package firmware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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
	METALCLOUD_FIRMWARE_REPOSITORY_SSH_USER	- the user of the SSH repository, for example: root

 * SCP related environment variables:
	METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH (required) 						- the path to the private OpenSSH key used for authentication, for example: ~/.ssh/my-openssh-key
	METALCLOUD_KNOWN_HOSTS_FILE_PATH (optional, defaults to ~/.ssh/known_hosts) - the path to the known_hosts file, for example: ~/.ssh/known_hosts
*/

const (
	configFormatJSON          = "json"
	configFormatYAML          = "yaml"
	catalogVendorDell         = "dell"
	catalogVendorLenovo       = "lenovo"
	catalogVendorHp           = "hp"
	catalogUpdateTypeOffline  = "offline"
	catalogUpdateTypeOnline   = "online"
	updateSeverityUnknown     = "unknown"
	updateSeverityRecommended = "recommended"
	updateSeverityCritical    = "critical"
	updateSeverityOptional    = "optional"
	serverTypesAll            = "*"

	batchSize         = 10
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

type repoConfiguration struct {
	HttpUrl    string
	SshPath    string
	SshPort    string
	SshUser    string
	SshKeyPath string
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
	Configuration                 map[string]any
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

type catalogMsObject struct {
	ServerFirmwareCatalogId                     int    `json:"server_firmware_catalog_id"`
	ServerFirmwareCatalogName                   string `json:"server_firmware_catalog_name"`
	ServerFirmwareCatalogDescription            string `json:"server_firmware_catalog_description"`
	ServerFirmwareCatalogVendor                 string `json:"server_firmware_catalog_vendor"`
	ServerFirmwareCatalogVendorId               string `json:"server_firmware_catalog_vendor_id"`
	ServerFirmwareCatalogVendorUrl              string `json:"server_firmware_catalog_vendor_url"`
	ServerFirmwareCatalogVendorReleaseTimestamp string `json:"server_firmware_catalog_vendor_release_timestamp"`
	ServerFirmwareCatalogUpdateType             string `json:"server_firmware_catalog_update_type"`
	ServerFirmwareCatalogCreatedTimestamp       string `json:"server_firmware_catalog_created_timestamp"`
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
		case catalogVendorDell, catalogVendorHp:
			if configFile.LocalFirmwarePath == "" || configFile.CatalogUrl == "" {
				firmwarePathValue, catalogUrlValue := "localFirmwarePath", "catalogUrl"

				if configFormat == configFormatYAML {
					firmwarePathValue, catalogUrlValue = "local_firmware_path", "catalog_url"
				}

				return fmt.Errorf("the '%s' and '%s' fields must be set in the raw-config file when downloading %s binaries", firmwarePathValue, catalogUrlValue, configFile.Vendor)
			}
		case catalogVendorLenovo:
			if configFile.LocalFirmwarePath == "" {
				firmwarePathValue := "localFirmwarePath"

				if configFormat == configFormatYAML {
					firmwarePathValue = "local_firmware_path"
				}

				return fmt.Errorf("the '%s' field must be set in the raw-config file when downloading %s binaries", firmwarePathValue, configFile.Vendor)
			}
		default:
			supportedVendors := []string{catalogVendorDell, catalogVendorHp, catalogVendorLenovo}
			return fmt.Errorf("the 'vendor' field must be one of %v", supportedVendors)
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

	err := checkStringSize(configFile.Name, 1, 255)
	if err != nil {
		return err
	}

	err = checkStringSize(configFile.Description, 1, 4096)
	if err != nil {
		return err
	}

	err = checkStringSize(configFile.CatalogUrl, 1, 255)
	if err != nil {
		return err
	}

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

func getRepoConfiguration(c *command.Command) repoConfiguration {
	repoConfig := repoConfiguration{
		HttpUrl:    command.GetStringParam(c.Arguments["repo_http_url"]),
		SshPath:    command.GetStringParam(c.Arguments["repo_ssh_path"]),
		SshPort:    command.GetStringParam(c.Arguments["repo_ssh_port"]),
		SshUser:    command.GetStringParam(c.Arguments["repo_ssh_user"]),
		SshKeyPath: command.GetStringParam(c.Arguments["user_private_ssh_key_path"]),
	}

	return repoConfig
}

func checkStringSize(str string, minimumSize, maximumSize int) error {
	if str == "" {
		return nil
	}

	if len(str) < minimumSize {
		return fmt.Errorf("the '%s' field must be at least %d characters", str, minimumSize)
	}

	if len(str) > maximumSize {
		return fmt.Errorf("the '%s' field must be less than %d characters", str, maximumSize)
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

func downloadBinariesFromCatalog(binaryCollection []*firmwareBinary, user, password string) error {
	fmt.Println("Downloading binaries.")

	for _, firmwareBinary := range binaryCollection {
		if !networking.CheckValidUrl(firmwareBinary.DownloadURL) {
			return fmt.Errorf("download URL '%s' is not valid.", firmwareBinary.DownloadURL)
		}

		err := DownloadFirmwareBinary(firmwareBinary, user, password)

		if err != nil {
			return err
		}
	}

	fmt.Println("Finished downloading binaries.")
	return nil
}

func uploadBinariesToRepository(binaryCollection []*firmwareBinary, replaceIfExists, skipHostKeyChecking bool, downloadUser, downloadPassword string, repoConfig repoConfiguration) error {
	firmwareRepositoryURL := repoConfig.HttpUrl
	if firmwareRepositoryURL == "" {
		var err error
		firmwareRepositoryURL, err = configuration.GetFirmwareRepositoryURL()
		if err != nil {
			return err
		}
	}

	remoteURL, err := url.Parse(firmwareRepositoryURL)
	if err != nil {
		return err
	}

	firmwareBinaryRepositoryHostname := remoteURL.Hostname()

	firmwareRepositorySSHPort := repoConfig.SshPort
	if firmwareRepositorySSHPort == "" {
		firmwareRepositorySSHPort = configuration.GetFirmwareRepositorySSHPort()
	}

	firmwareRepositorySSHPath := repoConfig.SshPath
	if firmwareRepositorySSHPath == "" {
		firmwareRepositorySSHPath, err = configuration.GetFirmwareRepositorySSHPath()
		if err != nil {
			return err
		}
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

	sshUser := repoConfig.SshUser
	if sshUser == "" {
		sshUser = configuration.GetFirmwareRepositorySSHUser()
	}

	userPrivateSSHKeyPath, err := configuration.GetUserPrivateSSHKeyPath()
	if err != nil {
		return err
	}

	scpClient, sshClient, err := networking.CreateSSHConnection(skipHostKeyChecking, firmwareRepositoryURL, firmwareRepositorySSHPort, sshUser, userPrivateSSHKeyPath)

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
		err := uploadBinaryToRepository(firmwareBinary, &scpClient, sshClient, firmwareBinaryExists, replaceIfExists, remotePath, downloadUser, downloadPassword)

		if err != nil {
			return err
		}
	}

	fmt.Println("Finished uploading binaries.")
	return nil
}

func uploadBinaryToRepository(binary *firmwareBinary, scpClient *scp.Client, sshClient *ssh.Client, firmwareBinaryExists, replaceIfExists bool, remotePath string, downloadUser, downloadPassword string) error {
	// Regenerate the session in the case it was previously closed, otherwise only the first file will be uploaded.
	scpSession, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer scpSession.Close()

	scpClient.Session = scpSession
	firmwareBinaryPath := binary.LocalPath

	var firmwareBinaryFile *os.File
	if firmwareBinaryPath == "" {
		if binary.DownloadURL != "" && !binary.HasErrors && (!firmwareBinaryExists || replaceIfExists) {
			// We don't save the binaries on the local filesystem, so we need to download them from the catalog as temporary files and then upload them to the repository.
			firmwareBinaryFile, err = os.CreateTemp(os.TempDir(), binary.FileName)
			if err != nil {
				return err
			}
			defer os.Remove(firmwareBinaryFile.Name())
			defer firmwareBinaryFile.Close()

			if !networking.CheckValidUrl(binary.DownloadURL) {
				return fmt.Errorf("download URL '%s' is not valid.", binary.DownloadURL)
			}

			binary.LocalPath = firmwareBinaryFile.Name()
			err := DownloadFirmwareBinary(binary, downloadUser, downloadPassword)
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

func DownloadFirmwareBinary(binary *firmwareBinary, user, password string) error {
	err := networking.DownloadFile(binary.DownloadURL, binary.LocalPath, binary.Hash, binary.HashingAlgorithm, user, password)

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

func sendCatalog(catalog firmwareCatalog) (catalogMsObject, error) {
	catalogMsObject := catalogMsObject{}

	endpoint, err := configuration.GetEndpoint()
	if err != nil {
		return catalogMsObject, err
	}

	catalogObject, err := createCatalogObject(catalog)
	if err != nil {
		return catalogMsObject, err
	}

	apiKey, err := configuration.GetAPIKey()
	if err != nil {
		return catalogMsObject, err
	}

	output, err := networking.SendMsRequest(networking.RequestTypePost, endpoint+networking.CatalogUrlPath, apiKey, catalogObject)
	if err != nil {
		return catalogMsObject, err
	}

	err = json.Unmarshal([]byte(output), &catalogMsObject)

	if err != nil {
		return catalogMsObject, fmt.Errorf("error parsing json response %s: %s", output, err.Error())
	}

	fmt.Printf("Created catalog with ID %d.\n", catalogMsObject.ServerFirmwareCatalogId)
	return catalogMsObject, nil
}

func createCatalogObject(catalog firmwareCatalog) ([]byte, error) {
	metalsoftServerTypesSupported, err := json.Marshal(catalog.MetalSoftServerTypesSupported)
	if err != nil {
		return nil, err
	}

	serverTypesSupported, err := json.Marshal(catalog.ServerTypesSupported)
	if err != nil {
		return nil, err
	}

	configuration, err := json.Marshal(catalog.Configuration)
	if err != nil {
		return nil, err
	}

	catalog.Name = trimToSize(catalog.Name, 255)
	catalog.Description = trimToSize(catalog.Description, 4096)
	catalog.VendorID = trimToSize(catalog.VendorID, 255)
	catalog.VendorURL = trimToSize(catalog.VendorURL, 255)

	catalogMap := map[string]any{
		"server_firmware_catalog_name":                                  catalog.Name,
		"server_firmware_catalog_description":                           catalog.Description,
		"server_firmware_catalog_vendor":                                catalog.Vendor,
		"server_firmware_catalog_update_type":                           catalog.UpdateType,
		"server_firmware_catalog_vendor_id":                             catalog.VendorID,
		"server_firmware_catalog_vendor_url":                            catalog.VendorURL,
		"server_firmware_catalog_vendor_release_timestamp":              catalog.VendorReleaseTimestamp,
		"server_firmware_catalog_metalsoft_server_types_supported_json": string(metalsoftServerTypesSupported),
		"server_firmware_catalog_vendor_server_types_supported_json":    string(serverTypesSupported),
		"server_firmware_catalog_vendor_configuration_json":             string(configuration),
	}

	for key, value := range catalogMap {
		if value == "" {
			delete(catalogMap, key)
		}
	}

	catalogObject, err := json.Marshal(catalogMap)

	if err != nil {
		return nil, err
	}

	return catalogObject, nil
}

func sendBinaries(binaryCollection []*firmwareBinary, catalogId int) error {
	endpoint, err := configuration.GetEndpoint()
	if err != nil {
		return err
	}

	binariesMap := []map[string]any{}
	counter := 0
	binariesCounter := 0
	for _, firmwareBinary := range binaryCollection {
		// Skip binaries that have errors. Most likely they were not found in the catalog.
		if firmwareBinary.HasErrors {
			continue
		}

		supportedDevices, err := json.Marshal(firmwareBinary.SupportedDevices)
		if err != nil {
			return err
		}

		supportedSystems, err := json.Marshal(firmwareBinary.SupportedSystems)
		if err != nil {
			return err
		}

		vendorProperties, err := json.Marshal(firmwareBinary.VendorProperties)
		if err != nil {
			return err
		}

		firmwareBinary.ExternalId = trimToSize(firmwareBinary.ExternalId, 255)
		firmwareBinary.VendorProperties["importantInfo"] = trimToSize(firmwareBinary.VendorProperties["importantInfo"], 255)
		firmwareBinary.DownloadURL = trimToSize(firmwareBinary.DownloadURL, 255)
		firmwareBinary.RepoURL = trimToSize(firmwareBinary.RepoURL, 255)
		firmwareBinary.Name = trimToSize(firmwareBinary.Name, 255)
		firmwareBinary.Description = trimToSize(firmwareBinary.Description, 4096)
		firmwareBinary.PackageVersion = trimToSize(firmwareBinary.PackageVersion, 255)

		binaryMap := map[string]any{
			"server_firmware_binary_catalog_id":                    catalogId,
			"server_firmware_binary_external_id":                   firmwareBinary.ExternalId,
			"server_firmware_binary_vendor_info_url":               firmwareBinary.VendorProperties["importantInfo"],
			"server_firmware_binary_vendor_download_url":           firmwareBinary.DownloadURL,
			"server_firmware_binary_cache_download_url":            firmwareBinary.RepoURL,
			"server_firmware_binary_name":                          firmwareBinary.Name,
			"server_firmware_binary_description":                   firmwareBinary.Description,
			"server_firmware_binary_package_id":                    firmwareBinary.PackageId,
			"server_firmware_binary_package_version":               firmwareBinary.PackageVersion,
			"server_firmware_binary_reboot_required":               firmwareBinary.RebootRequired,
			"server_firmware_binary_update_severity":               firmwareBinary.UpdateSeverity,
			"server_firmware_binary_vendor_supported_devices_json": string(supportedDevices),
			"server_firmware_binary_vendor_supported_systems_json": string(supportedSystems),
			"server_firmware_binary_vendor_release_timestamp":      firmwareBinary.VendorReleaseTimestamp,
			"server_firmware_binary_vendor_json":                   string(vendorProperties),
		}

		for key, value := range binaryMap {
			if value == "" {
				delete(binaryMap, key)
			}
		}

		binariesMap = append(binariesMap, binaryMap)
		counter++
		binariesCounter++
		if counter == batchSize {
			binariesObject, err := json.Marshal(binariesMap)
			if err != nil {
				return err
			}

			err = sendBinariesBatch(endpoint, binariesObject)
			if err != nil {
				return err
			}

			binariesMap = []map[string]any{}
			counter = 0
		}
	}

	if counter != 0 {
		binariesObject, err := json.Marshal(binariesMap)
		if err != nil {
			return err
		}

		err = sendBinariesBatch(endpoint, binariesObject)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Created %d binaries for catalog with ID %d.\n", binariesCounter, catalogId)
	return nil
}

func sendBinariesBatch(endpoint string, binariesObject []byte) error {
	apiKey, err := configuration.GetAPIKey()
	if err != nil {
		return err
	}

	_, err = networking.SendMsRequest(networking.RequestTypePost, endpoint+networking.BinaryUrlPath, apiKey, binariesObject)
	if err != nil {
		return err
	}

	return nil
}

func trimToSize(str string, size int) string {
	if len(str) > size {
		return str[:size]
	}

	return str
}
