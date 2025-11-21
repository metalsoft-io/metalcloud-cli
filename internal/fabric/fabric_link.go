package fabric

import (
	"context"
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var fabricLinkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"NetworkFabricId": {
			Title: "Fabric ID",
			Order: 2,
		},
		"NetworkDeviceARef": {
			Title: "Device A",
			Order: 3,
		},
		"NetworkDeviceAInterfaceRef": {
			Title: "Interface A",
			Order: 4,
		},
		"NetworkDeviceBRef": {
			Title: "Device B",
			Order: 5,
		},
		"NetworkDeviceBInterfaceRef": {
			Title: "Interface B",
			Order: 6,
		},
		"LinkType": {
			Title: "Type",
			Order: 7,
		},
		"MlagPair": {
			Title: "Mlag Pair",
			Order: 8,
		},
		"BgpNumbering": {
			Title: "Bgp Numbering",
			Order: 9,
		},
		"BgpLinkConfiguration": {
			Title: "Bgp Config",
			Order: 10,
		},
		"CustomVariables": {
			Title: "Custom Variables",
			Order: 11,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       12,
		},
	},
}

type NetworkFabricLinkFlat struct {
	Id                          float32
	NetworkFabricId             float32
	NetworkDeviceAId            float32
	NetworkDeviceAName          string
	NetworkDeviceARef           string
	NetworkDeviceAInterfaceId   float32
	NetworkDeviceAInterfaceName string
	NetworkDeviceAInterfaceRef  string
	NetworkDeviceBId            float32
	NetworkDeviceBName          string
	NetworkDeviceBRef           string
	NetworkDeviceBInterfaceId   float32
	NetworkDeviceBInterfaceName string
	NetworkDeviceBInterfaceRef  string
	LinkType                    string
	MlagPair                    bool
	BgpNumbering                string
	BgpLinkConfiguration        string
	CustomVariables             map[string]interface{}
	Status                      string
}

func FabricLinksGet(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Get fabric '%s' links", fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	linksList, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricLinks(ctx, int32(fabricIdNumeric)).Limit(1000).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(linksList.Data, &fabricLinkPrintConfig)
	}

	fabricInterfacesMap, err := getFabricInterfacesMap(ctx, client, fabricIdNumeric)
	if err != nil {
		return err
	}

	// Create flat links list
	flatLinks := []NetworkFabricLinkFlat{}
	for _, l := range linksList.Data {
		flatLinks = append(flatLinks, NetworkFabricLinkFlat{
			Id:                          l.Id,
			NetworkFabricId:             l.NetworkFabricId,
			NetworkDeviceAId:            fabricInterfacesMap[l.NetworkDeviceAInterfaceId].NetworkDeviceId,
			NetworkDeviceAName:          fabricInterfacesMap[l.NetworkDeviceAInterfaceId].NetworkDeviceName,
			NetworkDeviceARef:           fabricInterfacesMap[l.NetworkDeviceAInterfaceId].NetworkDeviceRef,
			NetworkDeviceAInterfaceId:   l.NetworkDeviceAInterfaceId,
			NetworkDeviceAInterfaceName: fabricInterfacesMap[l.NetworkDeviceAInterfaceId].InterfaceName,
			NetworkDeviceAInterfaceRef:  fabricInterfacesMap[l.NetworkDeviceAInterfaceId].InterfaceRef,
			NetworkDeviceBId:            fabricInterfacesMap[l.NetworkDeviceBInterfaceId].NetworkDeviceId,
			NetworkDeviceBName:          fabricInterfacesMap[l.NetworkDeviceBInterfaceId].NetworkDeviceName,
			NetworkDeviceBRef:           fabricInterfacesMap[l.NetworkDeviceBInterfaceId].NetworkDeviceRef,
			NetworkDeviceBInterfaceId:   l.NetworkDeviceBInterfaceId,
			NetworkDeviceBInterfaceName: fabricInterfacesMap[l.NetworkDeviceBInterfaceId].InterfaceName,
			NetworkDeviceBInterfaceRef:  fabricInterfacesMap[l.NetworkDeviceBInterfaceId].InterfaceRef,
			LinkType:                    l.LinkType,
			MlagPair:                    l.MlagPair != 0,
			BgpNumbering:                l.BgpNumbering,
			BgpLinkConfiguration:        l.BgpLinkConfiguration,
			CustomVariables:             l.CustomVariables,
			Status:                      l.Status,
		})
	}

	return formatter.PrintResult(flatLinks, &fabricLinkPrintConfig)
}

type interfaceInfo struct {
	InterfaceId       float32
	InterfaceName     string
	NetworkDeviceId   float32
	NetworkDeviceName string
	InterfaceRef      string
	NetworkDeviceRef  string
}

func getFabricInterfacesMap(ctx context.Context, client *sdk.APIClient, fabricId float32) (map[float32]interfaceInfo, error) {
	networkFabricNodes := map[float32]interfaceInfo{}

	networkDevices, httpRes, err := client.NetworkFabricAPI.GetFabricNetworkDevices(ctx, int32(fabricId)).Limit(1000).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for _, nd := range networkDevices.Data {
		nd_id, err := utils.GetFloat32FromString(nd.Id)
		if err != nil {
			continue
		}

		ports, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevicePorts(ctx, nd_id).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			continue
		}

		for _, p := range ports.Data {
			networkFabricNodes[p.InterfaceId] = interfaceInfo{
				p.InterfaceId,
				p.InterfaceName,
				nd_id,
				nd.IdentifierString,
				fmt.Sprintf("%s (%.0f)", p.InterfaceName, p.InterfaceId),
				fmt.Sprintf("%s (%.0f)", nd.IdentifierString, nd_id),
			}
		}
	}

	return networkFabricNodes, nil
}

func FabricLinkAdd(ctx context.Context, fabricId string, createLink sdk.CreateNetworkFabricLink) error {
	logger.Get().Info().Msgf("Adding link to fabric '%s'", fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	linkInfo, httpRes, err := client.NetworkFabricAPI.CreateNetworkFabricLink(ctx, int32(fabricIdNumeric)).
		CreateNetworkFabricLink(createLink).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(linkInfo, &fabricLinkPrintConfig)
}

func FabricLinkAddEx(ctx context.Context, fabricId string,
	networkDeviceA string,
	interfaceA string,
	networkDeviceB string,
	interfaceB string,
	linkType string,
	mlagPair bool,
	bgpNumbering string,
	bgpLinkConfiguration string,
	customVariables []string,
) error {
	fabricIdNumeric, err := utils.GetFloat32FromString(fabricId)
	if err != nil {
		return err
	}

	createLink := sdk.CreateNetworkFabricLink{
		LinkType:             linkType,
		BgpNumbering:         bgpNumbering,
		BgpLinkConfiguration: bgpLinkConfiguration,
		CustomVariables:      map[string]interface{}{},
	}

	// Lookup referenced devices and interfaces
	networkDeviceA = strings.ToLower(networkDeviceA)
	interfaceA = strings.ToLower(interfaceA)
	networkDeviceB = strings.ToLower(networkDeviceB)
	interfaceB = strings.ToLower(interfaceB)

	client := api.GetApiClient(ctx)

	networkDevices, httpRes, err := client.NetworkFabricAPI.GetFabricNetworkDevices(ctx, int32(fabricIdNumeric)).Limit(1000).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	for _, nd := range networkDevices.Data {
		deviceIdentifier := strings.ToLower(nd.IdentifierString)

		if deviceIdentifier != networkDeviceA && deviceIdentifier != networkDeviceB {
			continue
		}

		nd_id, err := utils.GetFloat32FromString(nd.Id)
		if err != nil {
			continue
		}

		ports, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevicePorts(ctx, nd_id).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			continue
		}

		for _, p := range ports.Data {
			if deviceIdentifier == networkDeviceA && strings.ToLower(p.InterfaceName) == interfaceA {
				createLink.NetworkDeviceAInterfaceId = p.InterfaceId
			}

			if deviceIdentifier == networkDeviceB && strings.ToLower(p.InterfaceName) == interfaceB {
				createLink.NetworkDeviceBInterfaceId = p.InterfaceId
			}
		}
	}

	if createLink.NetworkDeviceAInterfaceId == 0 || createLink.NetworkDeviceBInterfaceId == 0 {
		return fmt.Errorf("could not find match for the specified device and interface in the fabric")
	}

	// Add remaining properties
	if mlagPair {
		createLink.MlagPair = 1
	}

	for _, cv := range customVariables {
		parts := strings.Split(cv, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid custom variable format '%s' - expected key=value", cv)
		}

		createLink.CustomVariables[parts[0]] = parts[1]
	}

	return FabricLinkAdd(ctx, fabricId, createLink)
}

func FabricLinkRemove(ctx context.Context, fabricId string, linkId string) error {
	logger.Get().Info().Msgf("Removing link '%s' from fabric '%s'", linkId, fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	linkIdNumeric, err := utils.GetFloat32FromString(linkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkFabricAPI.DeleteNetworkFabricLink(ctx, int32(fabricIdNumeric), int32(linkIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}
