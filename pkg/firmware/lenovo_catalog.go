package firmware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const ENDPOINT_URL = "https://support.lenovo.com/services/ContentService/"

type LenovoCatalog struct {
	Data []SoftwareUpdate `json:"Data"`
}

type SoftwareUpdate struct {
	FixID            string   `json:"FixID"`
	ComponentID      string   `json:"ComponentID"`
	Files            []File   `json:"Files"`
	RequisitesFixIDs []string `json:"RequisitesFixIDs"`
}

type File struct {
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

	req, err := http.NewRequest(http.MethodPost, ENDPOINT_URL+"SearchDrivers", bytes.NewBuffer(jsonParams))
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

func SearchLenovoCatalog(machineType string, serialNumber string) (*LenovoCatalog, error) {
	targetInfos := map[string]string{
		"MachineType":  machineType,
		"SerialNumber": serialNumber,
	}
	response, err := retrieveAvailableFirmwareUpdates(targetInfos)
	if err != nil {
		return nil, err
	}

	lenovoCatalog := LenovoCatalog{}
	err = json.Unmarshal([]byte(response), &lenovoCatalog)
	if err != nil {
		return nil, err
	}

	return &lenovoCatalog, nil
}

func ExtractAvailableFirmwareUpdates(lenovoCatalog *LenovoCatalog) (UpdateType, UpdateRequiredType) {
	firmwareUpdates := lenovoCatalog.Data

	firmwareUpdate := map[string]map[string]string{
		"uefi":  make(map[string]string),
		"bmc":   make(map[string]string),
		"lxpm":  make(map[string]string),
		"other": make(map[string]string),
	}
	firmwareUpdateRequired := map[string]map[string][]string{
		"uefi":  make(map[string][]string),
		"bmc":   make(map[string][]string),
		"lxpm":  make(map[string][]string),
		"other": make(map[string][]string),
	}
	softwareUpdateMap := make(map[string]SoftwareUpdate)

	for _, softwareUpdate := range firmwareUpdates {
		softwareUpdateMap[softwareUpdate.FixID] = softwareUpdate

		softwareUpdateVersion := extractVersion(softwareUpdate.FixID)
		firmwareFix := findFirmwareFix(softwareUpdate.Files, "Fix")

		if softwareUpdateVersion == "" || firmwareFix == nil {
			continue
		}

		firmwareUpdateKey := "other"
		if softwareUpdate.ComponentID == "XCC" {
			firmwareUpdateKey = "bmc"
		} else if softwareUpdate.ComponentID == "UEFI" {
			firmwareUpdateKey = "uefi"
		} else if softwareUpdate.ComponentID == "LXPM" {
			firmwareUpdateKey = "lxpm"
		}

		if firmwareUpdateKey != "other" {
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

func resolveRequisites(requisites []string, softwareUpdateMap map[string]SoftwareUpdate) []string {
	var result []string

	for _, requisite := range requisites {
		if softwareUpdate, ok := softwareUpdateMap[requisite]; ok {
			firmwareFix := findFirmwareFix(softwareUpdate.Files, "Fix")
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

func findFirmwareFix(files []File, fileType string) *File {
	for _, file := range files {
		if file.Type == fileType {
			return &file
		}
	}
	return nil
}

func test1() {

	lenovoCatalog, err := SearchLenovoCatalog("7Y51", "J10227CF")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	firmwareUpdate, _ := ExtractAvailableFirmwareUpdates(lenovoCatalog)

	prettyFirmwareUpdate, err := json.MarshalIndent(firmwareUpdate, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(prettyFirmwareUpdate))
}
