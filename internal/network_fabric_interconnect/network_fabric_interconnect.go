package network_fabric_interconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var interconnectPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Order:    2,
			MaxWidth: 30,
		},
		"Name": {
			Order:    3,
			MaxWidth: 30,
		},
		"InterconnectType": {
			Title: "Type",
			Order: 4,
		},
		"BgpConfigurationTemplateId": {
			Title: "BGP Template",
			Order: 5,
		},
		"TransportId": {
			Title: "Transport",
			Order: 6,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"DeployId": {
			Title: "Deploy Job",
			Order: 8,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

var linkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"InterconnectId": {
			Title: "Interconnect",
			Order: 2,
		},
		"FabricId": {
			Title: "Fabric",
			Order: 3,
		},
		"NetworkEquipmentId": {
			Title: "Network Device",
			Order: 4,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var linkValidationPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"LinkId": {
			Title: "Link",
			Order: 1,
		},
		"NetworkDeviceId": {
			Title: "Network Device",
			Order: 2,
		},
		"CanActivate": {
			Title: "Can Activate",
			Order: 3,
		},
		"Errors": {
			Title:       "Errors",
			Order:       4,
			MaxWidth:    80,
			Transformer: formatJSONValue,
		},
	},
}

// formatJSONValue renders arbitrary nested values (e.g. validation error
// objects) as compact JSON so they fit in a table cell.
func formatJSONValue(value interface{}) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(encoded)
}

// PrintConfig exposes the interconnect table layout so other packages that list
// interconnects (e.g. the fabric command group) render them identically.
func PrintConfig() *formatter.PrintConfig {
	return &interconnectPrintConfig
}

func InterconnectList(ctx context.Context, filterStatus []string) error {
	logger.Get().Info().Msgf("Listing network fabric interconnects")

	client := api.GetApiClient(ctx)

	request := client.NetworkFabricInterconnectAPI.GetNetworkFabricInterconnects(ctx).SortBy([]string{"id:ASC"})
	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &interconnectPrintConfig)
}

func InterconnectGet(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Get network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnect, err := GetInterconnectByIdOrLabel(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &interconnectPrintConfig)
}

func InterconnectConfigExample(ctx context.Context) error {
	example := sdk.CreateNetworkFabricInterconnect{
		InterconnectType:           sdk.NETWORKFABRICINTERCONNECTTYPE_DCI_EVPN,
		Label:                      "dc1-dc2-interconnect",
		Name:                       sdk.PtrString("DC1 to DC2 interconnect"),
		Description:                sdk.PtrString("EVPN data center interconnect between DC1 and DC2"),
		BgpConfigurationTemplateId: 1,
		TransportId:                sdk.PtrInt64(1),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func InterconnectCreate(ctx context.Context, create sdk.CreateNetworkFabricInterconnect) error {
	logger.Get().Info().Msgf("Creating network fabric interconnect '%s'", create.Label)

	client := api.GetApiClient(ctx)

	interconnect, httpRes, err := client.NetworkFabricInterconnectAPI.
		CreateNetworkFabricInterconnect(ctx).
		CreateNetworkFabricInterconnect(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &interconnectPrintConfig)
}

func InterconnectUpdate(ctx context.Context, interconnectIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Updating network fabric interconnect '%s'", interconnectIdOrLabel)

	var update sdk.UpdateNetworkFabricInterconnect
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	interconnectId, revision, err := resolveInterconnectIdAndRevision(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	interconnect, httpRes, err := client.NetworkFabricInterconnectAPI.
		UpdateNetworkFabricInterconnect(ctx, interconnectId).
		UpdateNetworkFabricInterconnect(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &interconnectPrintConfig)
}

func InterconnectDelete(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Deleting network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkFabricInterconnectAPI.
		DeleteNetworkFabricInterconnect(ctx, interconnectId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network fabric interconnect '%s' deleted", interconnectIdOrLabel)
	return nil
}

func InterconnectDeploy(ctx context.Context, interconnectIdOrLabel string, requireConfirmation bool) error {
	logger.Get().Info().Msgf("Deploying network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricInterconnectAPI.
		DeployNetworkFabricInterconnect(ctx, interconnectId).
		NetworkFabricInterconnectDeployOptions(*sdk.NewNetworkFabricInterconnectDeployOptions(requireConfirmation)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

func InterconnectDeploymentCheck(ctx context.Context, interconnectIdOrLabel string, linkIds []string) error {
	logger.Get().Info().Msgf("Checking deployment of network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	linkIdsNumeric, err := utils.GetInt64SliceFromStrings(linkIds)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: without it the SDK posts a literal "null", which the API rejects.
	validations, httpRes, err := client.NetworkFabricInterconnectAPI.
		GetNetworkFabricInterconnectDeploymentCheck(ctx, interconnectId).
		NetworkFabricInterconnectDeploymentCheckRequest(sdk.NetworkFabricInterconnectDeploymentCheckRequest{LinkIds: linkIdsNumeric}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(validations, &linkValidationPrintConfig)
}

func InterconnectDeploymentInfo(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Get deployment info of network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	info, httpRes, err := client.NetworkFabricInterconnectAPI.
		GetNetworkFabricInterconnectDeploymentInfo(ctx, interconnectId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(info, nil)
}

func InterconnectAcceptDeploy(ctx context.Context, interconnectIdOrLabel string) error {
	return interconnectDeployDecision(ctx, interconnectIdOrLabel, true)
}

func InterconnectRejectDeploy(ctx context.Context, interconnectIdOrLabel string) error {
	return interconnectDeployDecision(ctx, interconnectIdOrLabel, false)
}

// interconnectDeployDecision accepts or rejects a pending deploy that was
// started with requireConfirmation. Both endpoints return no body.
func interconnectDeployDecision(ctx context.Context, interconnectIdOrLabel string, accept bool) error {
	action, outcome := "Rejecting", "rejected"
	if accept {
		action, outcome = "Accepting", "accepted"
	}
	logger.Get().Info().Msgf("%s deploy of network fabric interconnect '%s'", action, interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	var httpRes *http.Response
	if accept {
		httpRes, err = client.NetworkFabricInterconnectAPI.AcceptNetworkFabricInterconnectDeploy(ctx, interconnectId).Execute()
	} else {
		httpRes, err = client.NetworkFabricInterconnectAPI.RejectNetworkFabricInterconnectDeploy(ctx, interconnectId).Execute()
	}
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Deploy of network fabric interconnect '%s' %s", interconnectIdOrLabel, outcome)
	return nil
}

func InterconnectDetach(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Detaching network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricInterconnectAPI.
		DetachNetworkFabricInterconnect(ctx, interconnectId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

func InterconnectTemplateGet(ctx context.Context, interconnectType string) error {
	logger.Get().Info().Msgf("Get network fabric interconnect template for type '%s'", interconnectType)

	typeValue, err := sdk.NewNetworkFabricInterconnectTypeFromValue(interconnectType)
	if err != nil {
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	template, httpRes, err := client.NetworkFabricInterconnectAPI.
		GetNetworkFabricInterconnectTemplateByType(ctx, *typeValue).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(template, nil)
}

// fabricSummary is the slim fabric shape returned by the interconnect fabric
// endpoints. It lacks fabricConfiguration (and other fields sdk.NetworkFabric
// requires), so the strict SDK model cannot be used; the fabric table layout
// still applies and simply leaves the missing columns empty.
type fabricSummary struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	SiteId      *int64  `json:"siteId,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func InterconnectFabricsGet(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing fabrics attached to network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/network-fabric-interconnects/%d/fabrics", interconnectId), nil)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	rawItems, meta, err := utils.ParseRawPage(httpRes)
	if err != nil {
		return err
	}

	return printFabricSummaries(rawItems, meta)
}

func InterconnectAvailableFabricsGet(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing fabrics available for network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/network-fabric-interconnects/%d/fabrics-available?page=%.0f&limit=100&sortBy=id:ASC", interconnectId, page), nil)
	})
	if err != nil {
		return err
	}

	return printFabricSummaries(rawItems, meta)
}

func printFabricSummaries(rawItems []json.RawMessage, meta sdk.PaginatedResponseMeta) error {
	records, err := utils.UnmarshalRawItems[fabricSummary](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse fabrics: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), fabric.PrintConfig())
}

func InterconnectLinksGet(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing links of network fabric interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkFabricInterconnectAPI.
		GetInterconnectLinks(ctx, interconnectId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &linkPrintConfig)
}

func InterconnectLinkGet(ctx context.Context, interconnectIdOrLabel string, linkId string) error {
	logger.Get().Info().Msgf("Get link '%s' of network fabric interconnect '%s'", linkId, interconnectIdOrLabel)

	link, err := getInterconnectLink(ctx, interconnectIdOrLabel, linkId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(link, &linkPrintConfig)
}

func InterconnectLinkAdd(ctx context.Context, interconnectIdOrLabel string, fabricIdOrLabel string, networkDeviceIdOrLabel string) error {
	logger.Get().Info().Msgf("Adding link for fabric '%s' and network device '%s' to network fabric interconnect '%s'",
		fabricIdOrLabel, networkDeviceIdOrLabel, interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	fabricId, err := fabric.ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}

	networkDevice, err := network_device.GetNetworkDeviceByIdOrLabel(ctx, networkDeviceIdOrLabel)
	if err != nil {
		return err
	}
	networkDeviceId, err := utils.GetInt64FromString(networkDevice.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	link, httpRes, err := client.NetworkFabricInterconnectAPI.
		CreateInterconnectLink(ctx, interconnectId).
		CreateNetworkFabricInterconnectLink(sdk.CreateNetworkFabricInterconnectLink{
			FabricId:           fabricId,
			NetworkEquipmentId: networkDeviceId,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(link, &linkPrintConfig)
}

func InterconnectLinkRemove(ctx context.Context, interconnectIdOrLabel string, linkId string) error {
	logger.Get().Info().Msgf("Removing link '%s' from network fabric interconnect '%s'", linkId, interconnectIdOrLabel)

	link, err := getInterconnectLink(ctx, interconnectIdOrLabel, linkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkFabricInterconnectAPI.
		DeleteInterconnectLink(ctx, link.InterconnectId, link.Id).
		IfMatch(link.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Link '%s' removed from network fabric interconnect '%s'", linkId, interconnectIdOrLabel)
	return nil
}

func InterconnectLinksActivate(ctx context.Context, interconnectIdOrLabel string, linkIds []string, requireConfirmation bool) error {
	logger.Get().Info().Msgf("Activating links %v of network fabric interconnect '%s'", linkIds, interconnectIdOrLabel)

	interconnectId, linkIdsNumeric, err := resolveInterconnectAndLinkIds(ctx, interconnectIdOrLabel, linkIds)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricInterconnectAPI.
		ActivateNetworkFabricInterconnectLinks(ctx, interconnectId).
		NetworkFabricInterconnectLinksDeployOptions(*sdk.NewNetworkFabricInterconnectLinksDeployOptions(linkIdsNumeric, requireConfirmation)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

func InterconnectLinksDeactivate(ctx context.Context, interconnectIdOrLabel string, linkIds []string, requireConfirmation bool) error {
	logger.Get().Info().Msgf("Deactivating links %v of network fabric interconnect '%s'", linkIds, interconnectIdOrLabel)

	interconnectId, linkIdsNumeric, err := resolveInterconnectAndLinkIds(ctx, interconnectIdOrLabel, linkIds)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricInterconnectAPI.
		DeactivateNetworkFabricInterconnectLinks(ctx, interconnectId).
		NetworkFabricInterconnectDeactivateLinks(*sdk.NewNetworkFabricInterconnectDeactivateLinks(linkIdsNumeric, requireConfirmation)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

// GetInterconnectByIdOrLabel resolves an interconnect by numeric ID first and
// falls back to an exact label match, mirroring fabric.GetFabricByIdOrLabel.
func GetInterconnectByIdOrLabel(ctx context.Context, interconnectIdOrLabel string) (*sdk.NetworkFabricInterconnect, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := utils.GetInt64FromString(interconnectIdOrLabel); err == nil {
		interconnect, httpRes, err := client.NetworkFabricInterconnectAPI.GetNetworkFabricInterconnectById(ctx, idNumeric).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return interconnect, nil
		}
	}

	list, httpRes, err := client.NetworkFabricInterconnectAPI.
		GetNetworkFabricInterconnects(ctx).
		FilterLabel([]string{interconnectIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for i := range list.Data {
		if list.Data[i].Label == interconnectIdOrLabel {
			return &list.Data[i], nil
		}
	}

	err = fmt.Errorf("network fabric interconnect '%s' not found", interconnectIdOrLabel)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

// ResolveInterconnectNumericId resolves an interconnect ID or label to its numeric ID.
func ResolveInterconnectNumericId(ctx context.Context, interconnectIdOrLabel string) (int64, error) {
	interconnectId, _, err := resolveInterconnectIdAndRevision(ctx, interconnectIdOrLabel)
	return interconnectId, err
}

func resolveInterconnectIdAndRevision(ctx context.Context, interconnectIdOrLabel string) (int64, string, error) {
	interconnect, err := GetInterconnectByIdOrLabel(ctx, interconnectIdOrLabel)
	if err != nil {
		return 0, "", err
	}

	interconnectId, err := utils.GetInt64FromString(interconnect.Id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid network fabric interconnect ID %q: %w", interconnect.Id, err)
	}

	return interconnectId, interconnect.Revision, nil
}

func resolveInterconnectAndLinkIds(ctx context.Context, interconnectIdOrLabel string, linkIds []string) (int64, []int64, error) {
	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return 0, nil, err
	}

	linkIdsNumeric, err := utils.GetInt64SliceFromStrings(linkIds)
	if err != nil {
		return 0, nil, err
	}

	return interconnectId, linkIdsNumeric, nil
}

func getInterconnectLink(ctx context.Context, interconnectIdOrLabel string, linkId string) (*sdk.NetworkFabricInterconnectLink, error) {
	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return nil, err
	}

	linkIdNumeric, err := utils.GetInt64FromString(linkId)
	if err != nil {
		err = fmt.Errorf("invalid link ID: '%s'", linkId)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	client := api.GetApiClient(ctx)

	link, httpRes, err := client.NetworkFabricInterconnectAPI.
		GetInterconnectLink(ctx, interconnectId, linkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return link, nil
}
