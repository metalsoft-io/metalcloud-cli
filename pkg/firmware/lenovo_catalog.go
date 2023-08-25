package firmware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/filtering"
	"github.com/metalsoft-io/metalcloud-cli/internal/networking"

	"golang.org/x/exp/slices"
)

const (
	endpointUrl = "https://support.lenovo.com/services/ContentService/"

	softwareUpdateComponentXcc  = "XCC"
	softwareUpdateComponentUefi = "UEFI"
	softwareUpdateComponentLxpm = "LXPM"

	softwareUpdateTypeFix        = "Fix"
	softwareUpdateTypeInstallXML = "InstallXML"

	firmwareUpdateKeyUefi  = "uefi"
	firmwareUpdateKeyBmc   = "bmc"
	firmwareUpdateKeyLxpm  = "lxpm"
	firmwareUpdateKeyOther = "other"
)

type lenovoCatalog struct {
	Data []*softwareUpdate `json:"Data"`
}

type softwareUpdate struct {
	FixID            string                     `json:"FixID"`
	ComponentID      string                     `json:"ComponentID"`
	Files            []lenovoSoftwareUpdateFile `json:"Files"`
	RequisitesFixIDs []string                   `json:"RequisitesFixIDs"`
	Version          string
	UpdateKey        string
}

type lenovoSoftwareUpdateFile struct {
	Type        string `json:"Type"`
	Description string `json:"Description"`
	URL         string `json:"URL"`
	FileHash    string `json:"FileHash"`
}

func retrieveAvailableFirmwareUpdates(targetInfos map[string]string) (string, error) {
	searchParams := map[string]interface{}{
		"Category":            "",
		"FixIds":              "",
		"IsIncludeData":       "true",
		"IsIncludeMetaData":   "true",
		"IsIncludeRequisites": "true",
		"IsLatest":            "true",
		"QueryType":           "SUP",
		"SelectSupersedes":    "3",
		"SubmitterName":       "",
		"SubmitterVersion":    "",
		"TargetInfos":         []map[string]string{targetInfos},
		"XmlUpdateType":       "",
	}

	jsonParams, err := json.Marshal(searchParams)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, endpointUrl+"SearchDrivers", bytes.NewBuffer(jsonParams))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// Search the lenovo support site for the server firmware update information. A JSON response is returned and is saved in the local catalog path folder from the raw config file.
func generateLenovoCatalog(catalogFolder, machineType, serialNumber string, overwriteCatalog bool) (*lenovoCatalog, error) {
	if machineType == "" || serialNumber == "" {
		return nil, fmt.Errorf("machine type and serial number must be specified when searching for a lenovo catalog")
	}

	targetInfos := map[string]string{
		"MachineType":  machineType,
		"SerialNumber": serialNumber,
	}

	catalogName := fmt.Sprintf("lenovo_%s_%s.json", machineType, serialNumber)
	path := filepath.Join(catalogFolder, catalogName)

	lenovoCatalog := lenovoCatalog{}

	absoluteFilePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if fileExists(path) && !overwriteCatalog {
		fmt.Printf("Using existing Lenovo catalog %s\n", absoluteFilePath)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &lenovoCatalog)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Printf("Generating Lenovo catalog at path %s\n", absoluteFilePath)
		response, err := retrieveAvailableFirmwareUpdates(targetInfos)
		if err != nil {
			return nil, err
		}

		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, err = file.WriteString(response)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(response), &lenovoCatalog)
		if err != nil {
			return nil, err
		}
	}

	return &lenovoCatalog, nil
}

func extractAvailableFirmwareUpdates(lenovoCatalog *lenovoCatalog) map[string]*softwareUpdate {
	firmwareUpdates := lenovoCatalog.Data

	softwareUpdateMap := make(map[string]*softwareUpdate)

	for _, softwareUpdate := range firmwareUpdates {
		softwareUpdateMap[softwareUpdate.FixID] = softwareUpdate
	}

	for _, softwareUpdate := range firmwareUpdates {
		softwareUpdateVersion := extractVersion(softwareUpdate.FixID)
		firmwareFix := searchByFileType(softwareUpdate.Files, softwareUpdateTypeFix)

		if softwareUpdateVersion == "" || firmwareFix == nil {
			continue
		}

		softwareUpdate.Version = softwareUpdateVersion

		firmwareUpdateKey := firmwareUpdateKeyOther
		if softwareUpdate.ComponentID == softwareUpdateComponentXcc {
			firmwareUpdateKey = firmwareUpdateKeyBmc
		} else if softwareUpdate.ComponentID == softwareUpdateComponentUefi {
			firmwareUpdateKey = firmwareUpdateKeyUefi
		} else if softwareUpdate.ComponentID == softwareUpdateComponentLxpm {
			firmwareUpdateKey = firmwareUpdateKeyLxpm
		}

		softwareUpdate.UpdateKey = firmwareUpdateKey
	}

	return softwareUpdateMap
}

func resolveRequisites(requisites []string, softwareUpdateMap map[string]*softwareUpdate) []string {
	var result []string

	for _, requisite := range requisites {
		if softwareUpdate, ok := softwareUpdateMap[requisite]; ok {
			firmwareFix := searchByFileType(softwareUpdate.Files, softwareUpdateTypeFix)
			if firmwareFix != nil {
				componentPathArr := strings.Split(firmwareFix.URL, "/")
				componentName := componentPathArr[len(componentPathArr)-1]
				result = append(result, componentName)
			}
		}
	}

	return result
}

func searchByFileType(files []lenovoSoftwareUpdateFile, fileType string) *lenovoSoftwareUpdateFile {
	for _, file := range files {
		if file.Type == fileType {
			return &file
		}
	}
	return nil
}

func parseLenovoCatalog(configFile rawConfigFile, client metalcloud.MetalCloudClient, serverTypesFilter string, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) (firmwareCatalog, []*firmwareBinary, error) {
	catalog, serverInfoToCatalogMap, err := processLenovoCatalog(client, configFile, serverTypesFilter, downloadBinaries)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	firmwareBinaryCollection, err := processLenovoBinaries(configFile, serverInfoToCatalogMap, &catalog, uploadToRepo, downloadBinaries, repoConfig)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	return catalog, firmwareBinaryCollection, nil
}

func getSubmodelsAndSerialNumbers(client metalcloud.MetalCloudClient, supportedServerTypeNames []string) (map[string][]serverInfo, error) {
	if len(supportedServerTypeNames) == 0 {
		return nil, fmt.Errorf("no supported server type IDs were found")
	}

	filter := "server_type_name:"
	for _, serverTypeName := range supportedServerTypeNames {
		filter += fmt.Sprintf("%s,", serverTypeName)
	}

	// Remove the last trailing comma
	filter = filter[:len(filter)-1]

	list, err := client.ServersSearch(filtering.ConvertToSearchFieldFormat(filter))
	if err != nil {
		return nil, err
	}

	serverInfoMap := map[string][]serverInfo{}
	for _, server := range *list {
		if server.ServerSubmodel != "" && server.ServerSerialNumber != "" {
			serverInfoMap[server.ServerTypeName] = append(serverInfoMap[server.ServerTypeName], serverInfo{
				MachineType:  server.ServerSubmodel,
				SerialNumber: server.ServerSerialNumber,
				VendorSkuId:  server.ServerVendorSKUID,
			})
		}
	}

	return serverInfoMap, nil
}

func checkValidServerList(configFile rawConfigFile, serverFilteredInfoMap map[string][]serverInfo, serverInfoMap map[string][]serverInfo) error {
	for _, server := range configFile.ServersList {
		validServer := false
		for serverTypeName, servers := range serverInfoMap {
			if validServer {
				break
			}

			for _, serverInfo := range servers {
				if serverInfo.MachineType == server.MachineType && serverInfo.SerialNumber == server.SerialNumber {
					serverFilteredInfoMap[serverTypeName] = append(serverFilteredInfoMap[serverTypeName], serverInfo)
					validServer = true
					break
				}
			}
		}

		if !validServer {
			validServers := []serverInfo{}
			for _, servers := range serverInfoMap {
				validServers = append(validServers, servers...)
			}

			return fmt.Errorf("server with machine type %s and serial number %s was not found. Existing servers: %+v", server.MachineType, server.SerialNumber, validServers)
		}
	}

	return nil
}

func processLenovoCatalog(client metalcloud.MetalCloudClient, configFile rawConfigFile, serverTypesFilter string, downloadBinaries bool) (firmwareCatalog, map[serverInfo][]*lenovoCatalog, error) {
	var serverInfoMap map[string][]serverInfo
	serverInfoToCatalogMap := map[serverInfo][]*lenovoCatalog{}

	_, supportedServerTypeNames, err := retrieveSupportedServerTypes(client, serverTypesFilter)
	if err != nil {
		return firmwareCatalog{}, serverInfoToCatalogMap, err
	}

	serverInfoMap, err = getSubmodelsAndSerialNumbers(client, supportedServerTypeNames)

	if err != nil {
		return firmwareCatalog{}, serverInfoToCatalogMap, err
	}

	serverFilteredInfoMap := map[string][]serverInfo{}

	if len(configFile.ServersList) != 0 {
		err := checkValidServerList(configFile, serverFilteredInfoMap, serverInfoMap)
		if err != nil {
			return firmwareCatalog{}, serverInfoToCatalogMap, err
		}
	} else {
		serverFilteredInfoMap = serverInfoMap
	}

	if len(serverFilteredInfoMap) == 0 {
		return firmwareCatalog{}, serverInfoToCatalogMap, fmt.Errorf("no servers were found")
	}

	metalsoftSupportedServerTypes := []string{}
	serverTypesSupported := []string{}

	for metalsoftServerType, servers := range serverFilteredInfoMap {
		if !slices.Contains[string](metalsoftSupportedServerTypes, metalsoftServerType) {
			metalsoftSupportedServerTypes = append(metalsoftSupportedServerTypes, metalsoftServerType)
		}

		for _, server := range servers {
			if !slices.Contains[string](serverTypesSupported, server.VendorSkuId) {
				serverTypesSupported = append(serverTypesSupported, server.VendorSkuId)
			}

			generatedCatalog, err := generateLenovoCatalog(configFile.LocalCatalogPath, server.MachineType, server.SerialNumber, configFile.OverwriteCatalogs)
			if err != nil {
				return firmwareCatalog{}, serverInfoToCatalogMap, err
			}

			serverInfoToCatalogMap[server] = append(serverInfoToCatalogMap[server], generatedCatalog)
		}
	}

	catalogConfiguration := map[string]any{}
	vendorId := configFile.Vendor
	err = checkStringSize(vendorId, 1, 255)
	if err != nil {
		return firmwareCatalog{}, serverInfoToCatalogMap, err
	}

	catalog := firmwareCatalog{
		Name:                          configFile.Name,
		Description:                   configFile.Description,
		Vendor:                        configFile.Vendor,
		VendorID:                      vendorId,
		VendorURL:                     configFile.CatalogUrl,
		VendorReleaseTimestamp:        time.Now().Format(time.RFC3339),
		UpdateType:                    getUpdateType(configFile),
		MetalSoftServerTypesSupported: metalsoftSupportedServerTypes,
		ServerTypesSupported:          serverTypesSupported,
		Configuration:                 catalogConfiguration,
		CreatedTimestamp:              time.Now().Format(time.RFC3339),
	}

	return catalog, serverInfoToCatalogMap, nil
}

func processLenovoBinaries(configFile rawConfigFile, serverInfoToCatalogMap map[serverInfo][]*lenovoCatalog, catalog *firmwareCatalog, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) ([]*firmwareBinary, error) {
	firmwareBinaryCollection := []*firmwareBinary{}

	repositoryURL := repoConfig.HttpUrl
	if repositoryURL == "" {
		var err error
		repositoryURL, err = configuration.GetFirmwareRepositoryURL()
		if uploadToRepo && err != nil {
			return nil, fmt.Errorf("Error getting firmware repository URL: %v", err)
		}
	}

	for info, lenovoCatalogs := range serverInfoToCatalogMap {
		for _, lenovoCatalog := range lenovoCatalogs {
			softwareUpdateMap := extractAvailableFirmwareUpdates(lenovoCatalog)

			for _, softwareUpdate := range softwareUpdateMap {
				if softwareUpdate.UpdateKey == firmwareUpdateKeyOther {
					continue
				}

				firmwareFix := searchByFileType(softwareUpdate.Files, softwareUpdateTypeFix)
				if firmwareFix == nil {
					return nil, fmt.Errorf("no firmware fix was found for software update %s", softwareUpdate.FixID)
				}

				installXML := searchByFileType(softwareUpdate.Files, softwareUpdateTypeInstallXML)
				description := ""
				if installXML != nil {
					description = installXML.Description
				}

				componentVendorConfiguration := map[string]any{
					"requires": resolveRequisites(softwareUpdate.RequisitesFixIDs, softwareUpdateMap),
				}

				componentPathArr := strings.Split(firmwareFix.URL, "/")
				componentName := componentPathArr[len(componentPathArr)-1]
				componentRepoUrl, err := url.JoinPath(repositoryURL, componentName)
				if err != nil {
					return nil, err
				}

				localPath := ""
				if configFile.LocalFirmwarePath != "" && downloadBinaries {
					var err error
					localPath, err = filepath.Abs(filepath.Join(configFile.LocalFirmwarePath, componentName))

					if err != nil {
						return nil, fmt.Errorf("error getting download binary absolute path: %v", err)
					}
				}

				supportedDevices := []map[string]string{}

				supportedDevices = append(supportedDevices, map[string]string{
					"type": softwareUpdate.UpdateKey,
				})

				supportedSystems := []map[string]string{}

				supportedSystems = append(supportedSystems, map[string]string{
					"machineType":  info.MachineType,
					"serialNumber": info.SerialNumber,
				})

				firmwareBinary := firmwareBinary{
					ExternalId:             softwareUpdate.FixID,
					Name:                   softwareUpdate.FixID,
					FileName:               componentName,
					Description:            description,
					PackageId:              "",
					PackageVersion:         softwareUpdate.Version,
					RebootRequired:         true,
					UpdateSeverity:         updateSeverityUnknown,
					Hash:                   firmwareFix.FileHash,
					HashingAlgorithm:       networking.HashingAlgorithmSHA1,
					SupportedDevices:       supportedDevices,
					SupportedSystems:       supportedSystems,
					VendorProperties:       componentVendorConfiguration,
					VendorReleaseTimestamp: time.Now().Format(time.RFC3339),
					CreatedTimestamp:       time.Now().Format(time.RFC3339),
					DownloadURL:            firmwareFix.URL,
					RepoURL:                componentRepoUrl,
					LocalPath:              localPath,
				}

				firmwareBinaryCollection = append(firmwareBinaryCollection, &firmwareBinary)
			}
		}
	}

	return firmwareBinaryCollection, nil
}

func extractVersion(lenovoUpdateName string) string {
	version := ""
	components := strings.Split(lenovoUpdateName, "-")
	if len(components) > 1 {
		version = strings.Split(components[1], "_")[0]
	}
	return version
}
