package reports

import (
	"flag"
	"fmt"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var ReportsCmds = []command.Command{
	{
		Description:  "Statistics and other reports.",
		Subject:      "report",
		AltSubject:   "report",
		Predicate:    "devices",
		AltPredicate: "devices",
		FlagSet:      flag.NewFlagSet("list active devices in all datacenters", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: devicesListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func getActiveServers(datacenter string, client metalcloud.MetalCloudClient) (*[]metalcloud.ServerSearchResult, error) {
	servers, err := client.ServersSearch("datacenter_name:" + datacenter)
	if err != nil {
		return nil, err
	}

	filteredServers := []metalcloud.ServerSearchResult{}
	for _, s := range *servers {
		switch s.ServerStatus {
		case "available", "used", "cleaning", "cleaning_required", "available_reserved":
			filteredServers = append(filteredServers, s)
		}
	}

	return &filteredServers, nil
}

func getAllActiveSwitches(datacenter string, client metalcloud.MetalCloudClient) (*[]metalcloud.SwitchDevice, error) {

	switches := []metalcloud.SwitchDevice{}

	switchList, err := client.SwitchDevices(datacenter, "")
	if err != nil {
		return nil, err
	}

	for _, s := range *switchList {
		switches = append(switches, s)
	}

	return &switches, nil
}

func getAllActiveStoragePools(datacenter string, client metalcloud.MetalCloudClient) (*[]metalcloud.StoragePoolSearchResult, error) {
	return client.StoragePoolSearch("datacenter_name:" + datacenter)
}

type devicesList struct {
	servers  *[]metalcloud.ServerSearchResult
	switches *[]metalcloud.SwitchDevice
	storages *[]metalcloud.StoragePoolSearchResult
}

func getAllActiveDevices(datacenter string, client metalcloud.MetalCloudClient) (*devicesList, error) {

	servers, err := getActiveServers(datacenter, client)
	if err != nil {
		return nil, err
	}

	switches, err := getAllActiveSwitches(datacenter, client)
	if err != nil {
		return nil, err
	}

	storages, err := getAllActiveStoragePools(datacenter, client)
	if err != nil {
		return nil, err
	}

	return &devicesList{
		servers:  servers,
		switches: switches,
		storages: storages,
	}, nil
}

func devicesListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	DCs, err := client.Datacenters(true)
	if err != nil {
		return "", err
	}

	stats := map[string]*devicesList{}

	for _, dc := range *DCs {
		if dc.DatacenterIsMaster {
			continue
		}
		deviceList, err := getAllActiveDevices(dc.DatacenterName, client)
		if err != nil {
			return "", err
		}

		stats[dc.DatacenterName] = deviceList
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "DC_IDX",
			FieldType: tableformatter.TypeString,
			FieldSize: 3,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SERVERS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SWITCHES",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STORAGES",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	data := [][]interface{}{}

	totalServers := 0
	totalSwitches := 0
	totalStorages := 0

	dc_idx := 0
	for datacenterName, dcStats := range stats {

		serverCount := len(*dcStats.servers)
		switchesCount := len(*dcStats.switches)
		storagePoolsCount := len(*dcStats.storages)

		row := []interface{}{
			fmt.Sprintf("%d", dc_idx),
			datacenterName,
			fmt.Sprintf("%d", serverCount),
			fmt.Sprintf("%d", switchesCount),
			fmt.Sprintf("%d", storagePoolsCount),
		}

		data = append(data, row)
		dc_idx++

		totalServers += serverCount
		totalSwitches += switchesCount
		totalStorages += storagePoolsCount
	}

	totalsRow := []interface{}{
		"",
		colors.Bold("TOTAL"),
		colors.Bold(fmt.Sprintf("%d", totalServers)),
		colors.Bold(fmt.Sprintf("%d", totalSwitches)),
		colors.Bold(fmt.Sprintf("%d", totalStorages)),
	}

	data = append(data, totalsRow)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	totalDevices := totalServers + totalSwitches + totalStorages

	title := fmt.Sprint("Count of active or in-use equipment per datacenter")

	return table.RenderTable(fmt.Sprintf("Records (%d active devices across all datacenters)", totalDevices), title, command.GetStringParam(c.Arguments["format"]))

}
