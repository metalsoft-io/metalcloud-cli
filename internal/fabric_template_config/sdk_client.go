package fabric_template_config

import (
	"context"
	"encoding/base64"
	"strconv"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type sdkTemplateClient struct {
	ctx context.Context
	api *sdk.APIClient
}

// NewSDKClient builds a TemplateClient over the MetalSoft SDK.
func NewSDKClient(ctx context.Context, api *sdk.APIClient) TemplateClient {
	return &sdkTemplateClient{ctx: ctx, api: api}
}

func (c *sdkTemplateClient) GetFabric(fabricId int64) (*int64, string, error) {
	fabric, httpRes, err := c.api.NetworkFabricAPI.GetNetworkFabricById(c.ctx, fabricId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, "", err
	}
	return fabric.SiteId, fabric.Name, nil
}

func deviceRecordFromSDK(d *sdk.NetworkDevice) *deviceRecord {
	id, _ := strconv.ParseInt(d.Id, 10, 64)
	rec := &deviceRecord{
		Device: fsc.Device{
			Id:                id,
			Position:          d.Position,
			ManagementAddress: d.ManagementAddress,
			IdentifierString:  d.IdentifierString,
			Driver:            string(d.Driver),
			TagsMap:           d.TagsMap,
		},
		LoopbackAddressIpv4: d.LoopbackAddressIpv4,
		CustomVariables:     d.CustomVariables,
		Revision:            strconv.FormatInt(d.Revision, 10),
	}
	if d.Asn != 0 {
		asn := d.Asn
		rec.Asn = &asn
	}
	return rec
}

func (c *sdkTemplateClient) ListFabricDevices(fabricId int64) ([]*deviceRecord, error) {
	devices, _, err := utils.FetchAllPages(c.api.NetworkFabricAPI.GetFabricNetworkDevices(c.ctx, fabricId))
	if err != nil {
		return nil, err
	}
	out := make([]*deviceRecord, 0, len(devices))
	for i := range devices {
		out = append(out, deviceRecordFromSDK(&devices[i]))
	}
	return out, nil
}

func (c *sdkTemplateClient) ListDevicesBySite(siteId int64) ([]*deviceRecord, error) {
	devices, _, err := utils.FetchAllPages(c.api.NetworkDeviceAPI.GetNetworkDevices(c.ctx).FilterSiteId([]string{strconv.FormatInt(siteId, 10)}))
	if err != nil {
		return nil, err
	}
	out := make([]*deviceRecord, 0, len(devices))
	for i := range devices {
		out = append(out, deviceRecordFromSDK(&devices[i]))
	}
	return out, nil
}

func (c *sdkTemplateClient) ListTemplates() ([]*templateRecord, error) {
	templates, _, err := utils.FetchAllPages(c.api.DeviceConfigurationTemplateAPI.GetDeviceConfigurationTemplates(c.ctx))
	if err != nil {
		return nil, err
	}
	out := make([]*templateRecord, 0, len(templates))
	for i := range templates {
		t := &templates[i]
		rec := &templateRecord{Id: t.Id, Label: t.Label, Revision: t.Revision}
		if t.TemplateContent != nil {
			rec.TemplateB64 = *t.TemplateContent
			rec.HasContent = true
		}
		if t.Annotations != nil {
			rec.Annotations = *t.Annotations
		}
		out = append(out, rec)
	}
	return out, nil
}

func (c *sdkTemplateClient) GetTemplateContent(id int64) (string, string, error) {
	t, httpRes, err := c.api.DeviceConfigurationTemplateAPI.GetDeviceConfigurationTemplate(c.ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", "", err
	}
	content := ""
	if t.TemplateContent != nil {
		content = *t.TemplateContent
	}
	return content, t.Revision, nil
}

func (c *sdkTemplateClient) CreateTemplate(t templateCreate) (int64, error) {
	body := sdk.CreateDeviceConfigurationTemplate{
		Label:           t.Label,
		Name:            sdk.PtrString(t.Label),
		Description:     sdk.PtrString(t.Description),
		DeviceDriver:    sdk.SWITCHDRIVER_CUMULUS_LINUX,
		ExecutionType:   sdk.NETWORKTEMPLATEEXECUTIONTYPE_CLI,
		TemplateContent: sdk.PtrString(t.ContentB64),
	}
	if t.Annotations != nil {
		ann := t.Annotations
		body.Annotations = &ann
	}
	created, httpRes, err := c.api.DeviceConfigurationTemplateAPI.CreateDeviceConfigurationTemplate(c.ctx).CreateDeviceConfigurationTemplate(body).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, err
	}
	return created.Id, nil
}

func (c *sdkTemplateClient) UpdateTemplate(id int64, contentB64 *string, annotations *map[string]string, revision string) error {
	body := sdk.UpdateDeviceConfigurationTemplate{
		TemplateContent: contentB64,
		Annotations:     annotations,
	}
	_, httpRes, err := c.api.DeviceConfigurationTemplateAPI.
		UpdateDeviceConfigurationTemplate(c.ctx, id).
		UpdateDeviceConfigurationTemplate(body).
		IfMatch(revision).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func (c *sdkTemplateClient) ListProfiles() ([]*profileRecord, error) {
	profiles, _, err := utils.FetchAllPages(c.api.DeviceConfigurationTemplateAPI.GetDeviceConfigurationTemplateProfiles(c.ctx))
	if err != nil {
		return nil, err
	}
	out := make([]*profileRecord, 0, len(profiles))
	for i := range profiles {
		p := &profiles[i]
		rec := &profileRecord{
			Id:         p.Id,
			TemplateId: p.DeviceConfigurationTemplateId,
			Variables:  p.Variables,
			Priority:   p.Priority,
			IsEnabled:  p.IsEnabled,
			Revision:   p.Revision,
		}
		if devId, ok := p.GetNetworkDeviceIdOk(); ok && devId != nil {
			rec.DeviceId = *devId
		}
		if p.ApplyMode != nil {
			rec.ApplyMode = string(*p.ApplyMode)
		}
		out = append(out, rec)
	}
	return out, nil
}

func (c *sdkTemplateClient) CreateProfile(p profileCreate) error {
	body := sdk.CreateDeviceConfigurationTemplateProfile{
		DeviceConfigurationTemplateId: p.TemplateId,
		NetworkDeviceId:               *sdk.NewNullableInt64(sdk.PtrInt64(p.DeviceId)),
		NetworkFabricId:               *sdk.NewNullableInt64(sdk.PtrInt64(p.FabricId)),
		LifecycleStage:                lifecycleStagePtr(p.LifecycleStage),
		Variables:                     p.Variables,
		IsEnabled:                     sdk.PtrBool(p.IsEnabled),
		Priority:                      sdk.PtrFloat32(p.Priority),
		ApplyMode:                     applyModePtr(p.ApplyMode),
	}
	_, httpRes, err := c.api.DeviceConfigurationTemplateAPI.CreateDeviceConfigurationTemplateProfile(c.ctx).CreateDeviceConfigurationTemplateProfile(body).Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func (c *sdkTemplateClient) UpdateProfile(id int64, p profileUpdate, revision string) error {
	body := sdk.UpdateDeviceConfigurationTemplateProfile{
		Variables: p.Variables,
		IsEnabled: sdk.PtrBool(p.IsEnabled),
		Priority:  sdk.PtrFloat32(p.Priority),
		ApplyMode: applyModePtr(p.ApplyMode),
	}
	_, httpRes, err := c.api.DeviceConfigurationTemplateAPI.
		UpdateDeviceConfigurationTemplateProfile(c.ctx, id).
		UpdateDeviceConfigurationTemplateProfile(body).
		IfMatch(revision).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func (c *sdkTemplateClient) RenderTemplate(contentB64 string, variables map[string]interface{}) (string, error) {
	body := sdk.RenderDeviceConfigurationTemplate{TemplateContent: contentB64, Variables: variables}
	rendered, httpRes, err := c.api.DeviceConfigurationTemplateAPI.RenderDeviceConfigurationTemplate(c.ctx).RenderDeviceConfigurationTemplate(body).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", err
	}
	return rendered.Rendered, nil
}

func (c *sdkTemplateClient) GetDeviceCustomVariables(deviceId int64) (map[string]interface{}, string, string, error) {
	dev, httpRes, err := c.api.NetworkDeviceAPI.GetNetworkDevice(c.ctx, deviceId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, "", "", err
	}
	return dev.CustomVariables, dev.DriftDetectionSyncStatus, strconv.FormatInt(dev.Revision, 10), nil
}

func (c *sdkTemplateClient) UpdateDeviceCustomVariables(deviceId int64, customVariables map[string]interface{}, driftStatus string, revision string) error {
	body := sdk.UpdateNetworkDevice{CustomVariables: customVariables}
	// The server requires driftDetectionSyncStatus on every UpdateNetworkDevice
	// PATCH, but the generated SDK struct omits the field. Carry the current
	// value forward via AdditionalProperties so we don't mutate drift state.
	if driftStatus != "" {
		body.AdditionalProperties = map[string]interface{}{"driftDetectionSyncStatus": driftStatus}
	}
	_, httpRes, err := c.api.NetworkDeviceAPI.
		UpdateNetworkDevice(c.ctx, deviceId).
		UpdateNetworkDevice(body).
		IfMatch(revision).
		Execute()
	return response_inspector.InspectResponse(httpRes, err)
}

func lifecycleStagePtr(s string) *sdk.DeviceConfigurationProfileLifecycleStage {
	v := sdk.DeviceConfigurationProfileLifecycleStage(s)
	return &v
}

func applyModePtr(s string) *sdk.DeviceConfigurationProfileApplyMode {
	v := sdk.DeviceConfigurationProfileApplyMode(s)
	return &v
}

// base64Encode is the template-content encoding the engine expects.
func base64Encode(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }

func base64Decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}
