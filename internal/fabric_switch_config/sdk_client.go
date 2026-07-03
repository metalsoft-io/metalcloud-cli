package fabric_switch_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
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
		DriftDetectionSyncStatus:              d.DriftDetectionSyncStatus,
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

func (c *sdkClient) UpdateDevice(deviceId int64, body DeviceUpdate, driftStatus string, revision int64) error {
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
	// The server requires driftDetectionSyncStatus on every UpdateNetworkDevice
	// PATCH, but the generated SDK struct omits the field. Carry the current
	// value forward via AdditionalProperties so we don't mutate drift state.
	if driftStatus != "" {
		update.AdditionalProperties = map[string]interface{}{"driftDetectionSyncStatus": driftStatus}
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

// rawP2pInterface / rawP2pLink are a lenient view of the point-to-point-links
// listing. The SDK's typed decode fails here: the ipv4 subnet-allocation
// strategy is a oneOf whose "unnumbered" variant is just {kind, scope}, so every
// auto/manual strategy also matches it ("data matches more than one schema in
// oneOf"). We only need the link id/revision, the switch-interface ids, and
// whether an ipv4 strategy exists, so parse the raw body directly.
type rawP2pInterface struct {
	Type        string `json:"type"`
	InterfaceId int64  `json:"interfaceId"`
}

type rawP2pLink struct {
	Id         int64            `json:"id"`
	Revision   int64            `json:"revision"`
	InterfaceA *rawP2pInterface `json:"interfaceA"`
	InterfaceB *rawP2pInterface `json:"interfaceB"`
	Config     struct {
		Ipv4 *struct {
			SubnetAllocationStrategies []json.RawMessage `json:"subnetAllocationStrategies"`
		} `json:"ipv4"`
	} `json:"config"`
}

func (c *sdkClient) ListP2pLinks() ([]*P2pLinkRecord, error) {
	httpRes, err := api.DoJSONRequest(c.ctx, http.MethodGet, "/api/v2/point-to-point-links", nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("reading point-to-point links: %w", err)
	}
	return parseP2pLinksBody(body)
}

func parseP2pLinksBody(body []byte) ([]*P2pLinkRecord, error) {
	var links []rawP2pLink
	if err := json.Unmarshal(body, &links); err != nil {
		return nil, fmt.Errorf("parsing point-to-point links: %w", err)
	}

	const neqType = string(sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE)
	out := make([]*P2pLinkRecord, 0, len(links))
	for i := range links {
		link := &links[i]
		rec := &P2pLinkRecord{Id: link.Id, Revision: link.Revision}
		if link.InterfaceA != nil && link.InterfaceA.Type == neqType {
			id := link.InterfaceA.InterfaceId
			rec.InterfaceAId = &id
		}
		if link.InterfaceB != nil && link.InterfaceB.Type == neqType {
			id := link.InterfaceB.InterfaceId
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
	// Built as a raw body rather than via the SDK: a manual strategy's scope is
	// `{"kind":"global"}` with no resourceId, but the SDK's CreateResourceScope
	// always serializes resourceId:0, which the API rejects
	// ("scope.resourceId must not be less than 1").
	create := map[string]any{
		"interfaceA": pointToPointInterface(payload.InterfaceAId),
	}
	if payload.InterfaceBId != nil {
		create["interfaceB"] = pointToPointInterface(*payload.InterfaceBId)
	}
	if payload.Description != nil {
		create["description"] = *payload.Description
	}
	if payload.Mtu != nil {
		create["mtu"] = *payload.Mtu
	}
	if payload.RoutingActivation != "" {
		create["routingActivation"] = payload.RoutingActivation
	}
	if payload.StagedSubnetId != nil {
		create["ipv4"] = map[string]any{
			"subnetAllocationStrategies": []map[string]any{
				manualStrategyBody(*payload.StagedSubnetId, payload.StagedBinding),
			},
		}
	}

	body, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}

	httpRes, err := api.DoJSONRequest(c.ctx, http.MethodPost, "/api/v2/point-to-point-links", body)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	var link struct {
		Id       int64 `json:"id"`
		Revision int64 `json:"revision"`
	}
	if err := json.NewDecoder(httpRes.Body).Decode(&link); err != nil {
		return nil, fmt.Errorf("decoding created point-to-point link: %w", err)
	}
	return &P2pLinkRecord{Id: link.Id, Revision: link.Revision}, nil
}

func (c *sdkClient) CreateP2pIpv4Strategy(linkId, subnetId int64, binding string, linkRevision int64) error {
	// Raw request for the same scope-serialization reason as CreateP2pLink; the
	// config endpoint additionally requires If-Match with the link's revision.
	body, err := json.Marshal(manualStrategyBody(subnetId, binding))
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v2/point-to-point-links/%d/config/ipv4/subnet-allocation-strategies", linkId)
	headers := map[string]string{"If-Match": strconv.FormatInt(linkRevision, 10)}

	httpRes, err := api.DoJSONRequestWithHeaders(c.ctx, http.MethodPost, path, body, headers)
	return response_inspector.InspectResponse(httpRes, err)
}

func pointToPointInterface(interfaceId int64) map[string]any {
	return map[string]any{
		"type":        string(sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE),
		"interfaceId": interfaceId,
	}
}

// manualStrategyBody is the manual /31 allocation-strategy body. The scope is
// `{"kind":"global"}` with no resourceId (matching the reference implementation);
// see CreateP2pLink for why this is built by hand rather than via the SDK.
func manualStrategyBody(subnetId int64, binding string) map[string]any {
	return map[string]any{
		"kind":              string(sdk.POINTTOPOINTALLOCATIONSTRATEGYKIND_MANUAL),
		"subnetId":          subnetId,
		"scope":             map[string]any{"kind": string(sdk.RESOURCESCOPEKIND_GLOBAL)},
		"interfaceABinding": binding,
	}
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
