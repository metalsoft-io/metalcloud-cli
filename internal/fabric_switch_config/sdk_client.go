package fabric_switch_config

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// sdkClient is the production Client backed by the MetalSoft SDK.
type sdkClient struct {
	ctx context.Context
	api *sdk.APIClient
}

// NewSDKClient builds a Client over the given API client and context.
func NewSDKClient(ctx context.Context, api *sdk.APIClient) Client {
	return &sdkClient{ctx: ctx, api: api}
}

var revisionMismatchRe = regexp.MustCompile(`found (\d+)`)

func (c *sdkClient) GetFabric(fabricId int64) (*FabricInfo, error) {
	fabric, httpRes, err := c.api.NetworkFabricAPI.GetNetworkFabricById(c.ctx, fabricId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	id, _ := strconv.ParseInt(fabric.Id, 10, 64)
	return &FabricInfo{Id: id, Name: fabric.Name, SiteId: fabric.SiteId}, nil
}

func deviceRecordFromSDK(d *sdk.NetworkDevice) *DeviceRecord {
	id, _ := strconv.ParseInt(d.Id, 10, 64)
	rec := &DeviceRecord{
		Device: Device{
			Id:                id,
			Position:          d.Position,
			ManagementAddress: d.ManagementAddress,
			IdentifierString:  d.IdentifierString,
			Driver:            string(d.Driver),
			TagsMap:           d.TagsMap,
		},
		Asn:                                   d.Asn,
		LoopbackAddressIpv4:                   d.LoopbackAddressIpv4,
		ApplyIdentifierAsHostnameOnNextDeploy: d.ApplyIdentifierAsHostnameOnNextDeploy,
		Revision:                              d.Revision,
	}
	return rec
}

func (c *sdkClient) ListFabricDevices(fabricId int64) ([]*DeviceRecord, error) {
	request := c.api.NetworkFabricAPI.GetFabricNetworkDevices(c.ctx, fabricId)
	devices, _, err := utils.FetchAllPages(request)
	if err != nil {
		return nil, err
	}
	out := make([]*DeviceRecord, 0, len(devices))
	for i := range devices {
		out = append(out, deviceRecordFromSDK(&devices[i]))
	}
	return out, nil
}

func (c *sdkClient) ListDevicesBySite(siteId int64) ([]*DeviceRecord, error) {
	request := c.api.NetworkDeviceAPI.GetNetworkDevices(c.ctx).FilterSiteId([]string{strconv.FormatInt(siteId, 10)})
	devices, _, err := utils.FetchAllPages(request)
	if err != nil {
		return nil, err
	}
	out := make([]*DeviceRecord, 0, len(devices))
	for i := range devices {
		out = append(out, deviceRecordFromSDK(&devices[i]))
	}
	return out, nil
}

func (c *sdkClient) UpdateDevice(deviceId int64, body DeviceUpdate, revision int64) error {
	update := sdk.UpdateNetworkDevice{}
	if body.IdentifierString != nil {
		update.SetIdentifierString(*body.IdentifierString)
	}
	if body.ApplyIdentifierAsHostnameOnNextDeploy != nil {
		update.SetApplyIdentifierAsHostnameOnNextDeploy(*body.ApplyIdentifierAsHostnameOnNextDeploy)
	}
	if body.Asn != nil {
		update.SetAsn(*body.Asn)
	}
	if body.LoopbackAddress != nil {
		update.SetLoopbackAddress(*body.LoopbackAddress)
	}
	_, httpRes, err := c.api.NetworkDeviceAPI.
		UpdateNetworkDevice(c.ctx, deviceId).
		UpdateNetworkDevice(update).
		IfMatch(strconv.FormatInt(revision, 10)).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func (c *sdkClient) ListPorts(deviceId int64) ([]*PortRecord, error) {
	portsInfo, httpRes, err := c.api.NetworkDeviceAPI.GetNetworkDevicePorts(c.ctx, float32(deviceId)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	out := make([]*PortRecord, 0, len(portsInfo.Data))
	for i := range portsInfo.Data {
		p := &portsInfo.Data[i]
		rec := &PortRecord{
			InterfaceId:    p.InterfaceId,
			InterfaceName:  p.InterfaceName,
			Kind:           p.Kind,
			ConfigRevision: p.Config.Revision,
		}
		if enabled, ok := p.Config.GetEnabledOk(); ok {
			rec.Enabled = enabled
		}
		if desc, ok := p.Config.GetDescriptionOk(); ok {
			rec.Description = desc
		}
		for _, addr := range p.Ipv4.Addresses {
			rec.Ipv4Addresses = append(rec.Ipv4Addresses, IpAddress{Address: addr.Address, PrefixLength: addr.PrefixLength})
		}
		out = append(out, rec)
	}
	return out, nil
}

func (c *sdkClient) UpdatePortConfig(deviceId, portId int64, enabled *bool, description *string, configRevision int64) error {
	update := sdk.UpdateNetworkEquipmentInterfaceConfig{}
	if enabled != nil {
		update.Enabled = *sdk.NewNullableBool(enabled)
	}
	if description != nil {
		update.Description = *sdk.NewNullableString(description)
	}
	_, httpRes, err := c.api.NetworkDeviceAPI.
		UpdateNetworkDevicePortConfig(c.ctx, deviceId, portId).
		UpdateNetworkEquipmentInterfaceConfig(update).
		IfMatch(strconv.FormatInt(configRevision, 10)).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func (c *sdkClient) AddPortIpv4(deviceId, portId int64, address string, prefixLength int32, configRevision int64) error {
	payload := sdk.AddNetworkEquipmentInterfaceIp{Address: address, PrefixLength: prefixLength}
	// The lock checks a one-based counter while the single-port GET exposes a
	// zero-based config.revision: send revision+1, retry once on a 409 mismatch.
	revision := configRevision + 1
	_, httpRes, err := c.api.NetworkDeviceAPI.
		AddNetworkDevicePortIp(c.ctx, deviceId, portId, "ipv4").
		AddNetworkEquipmentInterfaceIp(payload).
		IfMatch(strconv.FormatInt(revision, 10)).
		Execute()
	if err != nil && httpRes != nil && httpRes.StatusCode == 409 {
		if expected := expectedRevision(err); expected != "" {
			_, httpRes, err = c.api.NetworkDeviceAPI.
				AddNetworkDevicePortIp(c.ctx, deviceId, portId, "ipv4").
				AddNetworkEquipmentInterfaceIp(payload).
				IfMatch(expected).
				Execute()
		}
	}
	return response_inspector.InspectResponse(httpRes, err)
}

func expectedRevision(err error) string {
	var apiErr sdk.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if m := revisionMismatchRe.FindSubmatch(apiErr.Body()); m != nil {
			return string(m[1])
		}
	}
	return ""
}

func (c *sdkClient) ListP2pLinks() ([]*P2pLinkRecord, error) {
	links, httpRes, err := c.api.PointToPointLinkAPI.GetPointToPointLinks(c.ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	out := make([]*P2pLinkRecord, 0, len(links))
	for i := range links {
		link := &links[i]
		rec := &P2pLinkRecord{Id: link.Id, Revision: link.Revision}
		if iface, ok := link.GetInterfaceAOk(); ok && iface != nil && iface.Type == sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE {
			id := iface.InterfaceId
			rec.InterfaceAId = &id
		}
		if iface, ok := link.GetInterfaceBOk(); ok && iface != nil && iface.Type == sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE {
			id := iface.InterfaceId
			rec.InterfaceBId = &id
		}
		if link.Config.Ipv4 != nil && len(link.Config.Ipv4.SubnetAllocationStrategies) > 0 {
			rec.HasIpv4Strategy = true
		}
		out = append(out, rec)
	}
	return out, nil
}

func (c *sdkClient) CreateP2pLink(payload P2pLinkCreate) (*P2pLinkRecord, error) {
	create := sdk.CreatePointToPointLink{}
	create.SetInterfaceA(sdk.CreatePointToPointInterface{
		Type:        sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE,
		InterfaceId: payload.InterfaceAId,
	})
	if payload.InterfaceBId != nil {
		create.SetInterfaceB(sdk.CreatePointToPointInterface{
			Type:        sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE,
			InterfaceId: *payload.InterfaceBId,
		})
	}
	if payload.Description != nil {
		create.SetDescription(*payload.Description)
	}
	if payload.Mtu != nil {
		create.SetMtu(*payload.Mtu)
	}
	if payload.RoutingActivation != "" {
		create.SetRoutingActivation(sdk.PointToPointLinkRoutingActivation(payload.RoutingActivation))
	}
	if payload.StagedSubnetId != nil {
		strategy, err := manualStrategyRequest(*payload.StagedSubnetId, payload.StagedBinding)
		if err != nil {
			return nil, err
		}
		create.Ipv4 = &sdk.CreatePointToPointLinkIpv4Properties{
			SubnetAllocationStrategies: []sdk.CreatePointToPointLinkConfigIpv4SubnetAllocationStrategyRequest{strategy},
		}
	}
	link, httpRes, err := c.api.PointToPointLinkAPI.CreatePointToPointLink(c.ctx).CreatePointToPointLink(create).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	return &P2pLinkRecord{Id: link.Id, Revision: link.Revision}, nil
}

func (c *sdkClient) CreateP2pIpv4Strategy(linkId, subnetId int64, binding string, linkRevision int64) error {
	strategy, err := manualStrategyRequest(subnetId, binding)
	if err != nil {
		return err
	}
	_, httpRes, err := c.api.PointToPointLinkAPI.
		CreatePointToPointLinkConfigIpv4SubnetAllocationStrategy(c.ctx, float32(linkId)).
		CreatePointToPointLinkConfigIpv4SubnetAllocationStrategyRequest(strategy).
		IfMatch(strconv.FormatInt(linkRevision, 10)).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func manualStrategyRequest(subnetId int64, binding string) (sdk.CreatePointToPointLinkConfigIpv4SubnetAllocationStrategyRequest, error) {
	bindingValue, err := sdk.NewPointToPointInterfaceBindingFromValue(binding)
	if err != nil {
		return sdk.CreatePointToPointLinkConfigIpv4SubnetAllocationStrategyRequest{}, err
	}
	manual := sdk.NewCreateManualIpv4PointToPointAllocationStrategy(
		sdk.POINTTOPOINTALLOCATIONSTRATEGYKIND_MANUAL,
		sdk.CreateResourceScope{Kind: sdk.RESOURCESCOPEKIND_GLOBAL},
		subnetId,
	)
	manual.InterfaceABinding = bindingValue
	return sdk.CreatePointToPointLinkConfigIpv4SubnetAllocationStrategyRequest{
		CreateManualIpv4PointToPointAllocationStrategy: manual,
	}, nil
}

func (c *sdkClient) ListSubnetsByFabricTag(fabricId int64) ([]*SubnetRecord, error) {
	request := c.api.SubnetAPI.GetSubnets(c.ctx)
	subnets, _, err := utils.FetchAllPages(request)
	if err != nil {
		return nil, err
	}
	want := strconv.FormatInt(fabricId, 10)
	var out []*SubnetRecord
	for i := range subnets {
		s := &subnets[i]
		if s.Tags[FabricTag] != want {
			continue
		}
		out = append(out, &SubnetRecord{
			Id:             s.Id,
			NetworkAddress: s.NetworkAddress,
			PrefixLength:   s.PrefixLength,
			Tags:           s.Tags,
		})
	}
	return out, nil
}

func (c *sdkClient) CreateSubnet(payload SubnetCreate) (*SubnetRecord, error) {
	create := sdk.CreateSubnet{
		NetworkAddress: payload.NetworkAddress,
		PrefixLength:   payload.PrefixLength,
		IsPool:         false,
	}
	if payload.Name != "" {
		create.SetName(payload.Name)
	}
	if payload.Tags != nil {
		create.SetTags(payload.Tags)
	}
	subnet, httpRes, err := c.api.SubnetAPI.CreateSubnet(c.ctx).CreateSubnet(create).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	return &SubnetRecord{Id: subnet.Id, NetworkAddress: subnet.NetworkAddress, PrefixLength: subnet.PrefixLength, Tags: subnet.Tags}, nil
}
