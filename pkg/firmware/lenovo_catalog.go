package firmware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/filtering"
)

const (
	endpointUrl = "https://support.lenovo.com/services/ContentService/"

	softwareUpdateComponentXcc  = "XCC"
	softwareUpdateComponentUefi = "UEFI"
	softwareUpdateComponentLxpm = "LXPM"

	softwareUpdateTypeFix = "Fix"

	firmwareUpdateKeyUefi  = "uefi"
	firmwareUpdateKeyBmc   = "bmc"
	firmwareUpdateKeyLxpm  = "lxpm"
	firmwareUpdateKeyOther = "other"
)

type lenovoCatalog struct {
	Data []softwareUpdate `json:"Data"`
}

type softwareUpdate struct {
	FixID            string                     `json:"FixID"`
	ComponentID      string                     `json:"ComponentID"`
	Files            []lenovoSoftwareUpdateFile `json:"Files"`
	RequisitesFixIDs []string                   `json:"RequisitesFixIDs"`
}

type lenovoSoftwareUpdateFile struct {
	Type string `json:"Type"`
	URL  string `json:"URL"`
}

type UpdateType map[string]map[string]string
type UpdateRequiredType map[string]map[string][]string

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

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// Search the lenovo support site for the server firmware update information. A JSON response is returned and is saved in the local catalog path folder from the raw config file.
func generateLenovoCatalog(catalogFolder, machineType, serialNumber string) (*lenovoCatalog, error) {
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

	if fileExists(path) {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &lenovoCatalog)
		if err != nil {
			return nil, err
		}
	} else {
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

func extractAvailableFirmwareUpdates(lenovoCatalog *lenovoCatalog) (UpdateType, UpdateRequiredType) {
	firmwareUpdates := lenovoCatalog.Data

	firmwareUpdate := map[string]map[string]string{
		firmwareUpdateKeyUefi:  make(map[string]string),
		firmwareUpdateKeyBmc:   make(map[string]string),
		firmwareUpdateKeyLxpm:  make(map[string]string),
		firmwareUpdateKeyOther: make(map[string]string),
	}
	firmwareUpdateRequired := map[string]map[string][]string{
		firmwareUpdateKeyUefi:  make(map[string][]string),
		firmwareUpdateKeyBmc:   make(map[string][]string),
		firmwareUpdateKeyLxpm:  make(map[string][]string),
		firmwareUpdateKeyOther: make(map[string][]string),
	}
	softwareUpdateMap := make(map[string]softwareUpdate)

	for _, softwareUpdate := range firmwareUpdates {
		softwareUpdateMap[softwareUpdate.FixID] = softwareUpdate

		softwareUpdateVersion := extractVersion(softwareUpdate.FixID)
		firmwareFix := findFirmwareFix(softwareUpdate.Files, softwareUpdateTypeFix)

		if softwareUpdateVersion == "" || firmwareFix == nil {
			continue
		}

		firmwareUpdateKey := firmwareUpdateKeyOther
		if softwareUpdate.ComponentID == softwareUpdateComponentXcc {
			firmwareUpdateKey = firmwareUpdateKeyBmc
		} else if softwareUpdate.ComponentID == softwareUpdateComponentUefi {
			firmwareUpdateKey = firmwareUpdateKeyUefi
		} else if softwareUpdate.ComponentID == softwareUpdateComponentLxpm {
			firmwareUpdateKey = firmwareUpdateKeyLxpm
		}

		if firmwareUpdateKey != firmwareUpdateKeyOther {
			firmwareUpdate[firmwareUpdateKey][softwareUpdateVersion] = firmwareFix.URL
			if len(softwareUpdate.RequisitesFixIDs) > 0 {
				firmwareUpdateRequired[firmwareUpdateKey][softwareUpdateVersion] = resolveRequisites(softwareUpdate.RequisitesFixIDs, softwareUpdateMap)
			}
		} else {
			firmwareUpdate[firmwareUpdateKey][softwareUpdate.ComponentID+"-"+softwareUpdateVersion] = firmwareFix.URL
			if len(softwareUpdate.RequisitesFixIDs) > 0 {
				firmwareUpdateRequired[firmwareUpdateKey][softwareUpdate.ComponentID+"-"+softwareUpdateVersion] = resolveRequisites(softwareUpdate.RequisitesFixIDs, softwareUpdateMap)
			}
		}
	}

	return firmwareUpdate, firmwareUpdateRequired
}

func resolveRequisites(requisites []string, softwareUpdateMap map[string]softwareUpdate) []string {
	var result []string

	for _, requisite := range requisites {
		if softwareUpdate, ok := softwareUpdateMap[requisite]; ok {
			firmwareFix := findFirmwareFix(softwareUpdate.Files, softwareUpdateTypeFix)
			if firmwareFix != nil {
				result = append(result, firmwareFix.URL)
			}
		}
	}

	return result
}

func extractVersion(lenovoUpdateName string) string {
	version := ""
	components := strings.Split(lenovoUpdateName, "-")
	if len(components) > 1 {
		version = strings.Split(components[1], "_")[0]
	}
	return version
}

func findFirmwareFix(files []lenovoSoftwareUpdateFile, fileType string) *lenovoSoftwareUpdateFile {
	for _, file := range files {
		if file.Type == fileType {
			return &file
		}
	}
	return nil
}

func parseLenovoCatalog(configFile rawConfigFile, client metalcloud.MetalCloudClient, serverTypesFilter string, uploadToRepo bool) (firmwareCatalog, []*firmwareBinary, error) {
	_, err := processLenovoCatalog(client, configFile, serverTypesFilter)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	return firmwareCatalog{}, nil, fmt.Errorf("not implemented")

	// 	firmwareBinaryCollection := []*firmwareBinary{}

	// 	firmwareUpdates, _ := extractAvailableFirmwareUpdates(currentLenovoCatalog)

	// 	for componentType, updateVersions := range firmwareUpdates {
	// 		for version, downloadURL := range updateVersions {
	// 			componentVendorConfiguration := map[string]string{}

	// 			firmwareBinary := firmwareBinary{
	// 				ExternalId:             downloadURL,
	// 				Name:                   componentType,
	// 				Description:            componentType,
	// 				PackageId:              "",
	// 				PackageVersion:         version,
	// 				RebootRequired:         true,
	// 				UpdateSeverity:         updateSeverityUnknown,
	// 				SupportedDevices:       []map[string]string{},
	// 				SupportedSystems:       []map[string]string{},
	// 				VendorProperties:       componentVendorConfiguration,
	// 				VendorReleaseTimestamp: time.Now().Format(time.RFC3339),
	// 				CreatedTimestamp:       time.Now().Format(time.RFC3339),
	// 				DownloadURL:            downloadURL,
	// 				RepoURL:                downloadURL,
	// 			}

	// 			firmwareBinaryCollection = append(firmwareBinaryCollection, &firmwareBinary)
	// 		}
	// 	}

	// 	for _, firmwareBinary := range firmwareBinaryCollection {
	// 		prettyFirmwareBinary, err := json.MarshalIndent(firmwareBinary, "", "  ")
	// 		if err != nil {
	// 			return firmwareCatalog{}, nil, err
	// 		}

	// 		fmt.Println(string(prettyFirmwareBinary))
	// 	}
	// }

	// return catalog, firmwareBinaryCollection, nil
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
			})
		}
	}

	return serverInfoMap, nil
}

func processLenovoCatalog(client metalcloud.MetalCloudClient, configFile rawConfigFile, serverTypesFilter string) (firmwareCatalog, error) {
	var serverInfoMap map[string][]serverInfo

	_, supportedServerTypeNames, err := retrieveSupportedServerTypes(client, serverTypesFilter)
	if err != nil {
		return firmwareCatalog{}, err
	}

	serverInfoMap, err = getSubmodelsAndSerialNumbers(client, supportedServerTypeNames)

	if err != nil {
		return firmwareCatalog{}, err
	}

	serverFilteredInfoMap := map[string][]serverInfo{}

	if len(configFile.ServersList) != 0 {
		for _, server := range configFile.ServersList {
			validServer := false
			for serverTypeName, servers := range serverInfoMap {
				if validServer{
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

				return firmwareCatalog{}, fmt.Errorf("server with machine type %s and serial number %s was not found. Valid servers parameters: %+v", server.MachineType, server.SerialNumber, validServers)
			}
		}
	} else {
		serverFilteredInfoMap = serverInfoMap
	}

	fmt.Printf("Server map: %+v\n", serverFilteredInfoMap)

	for _, servers := range serverInfoMap {
		for _, server := range servers {
			generateLenovoCatalog(configFile.LocalCatalogPath, server.MachineType, server.SerialNumber)
		}
	}

	catalogConfiguration := map[string]string{}
	vendorId := configFile.Vendor
	checkStringSize(vendorId)

	catalog := firmwareCatalog{
		Name:                   configFile.Name,
		Description:            configFile.Description,
		Vendor:                 configFile.Vendor,
		VendorID:               vendorId,
		VendorURL:              configFile.CatalogUrl,
		VendorReleaseTimestamp: time.Now().Format(time.RFC3339),
		UpdateType:             getUpdateType(configFile),
		ServerTypesSupported:   []string{},
		Configuration:          catalogConfiguration,
		CreatedTimestamp:       time.Now().Format(time.RFC3339),
	}

	fmt.Printf("Created catalog object %+v\n", catalog)

	return catalog, nil
}
