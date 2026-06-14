package device_auth_provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var DeviceAuthProviderPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Kind": {
			Order: 4,
		},
		"SiteId": {
			Title: "Site",
			Order: 5,
		},
		"IpAddress": {
			Title: "IP",
			Order: 6,
		},
		"Port": {
			Order: 7,
		},
		"Username": {
			MaxWidth: 30,
			Order:    8,
		},
		"HasSharedSecret": {
			Title:       "Shared Secret",
			Transformer: formatter.FormatBooleanValue,
			Order:       9,
		},
		"HasPassword": {
			Title:       "Password",
			Transformer: formatter.FormatBooleanValue,
			Order:       10,
		},
		"Status": {
			Order:       11,
			Transformer: formatter.FormatStatusValue,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       12,
		},
	},
}

var DeviceAuthProviderCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Username": {
			Order: 1,
		},
		"Password": {
			Order: 2,
		},
		"SharedSecret": {
			Title: "Shared Secret",
			Order: 3,
		},
	},
}

func GetDeviceAuthProviderByIdOrLabel(ctx context.Context, idOrLabel string) (*sdk.DeviceAuthProvider, error) {
	client := api.GetApiClient(ctx)

	if id, err := strconv.ParseInt(idOrLabel, 10, 64); err == nil {
		provider, httpRes, err := client.SiteAPI.GetDeviceAuthProviderById(ctx, id).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, err
		}
		return provider, nil
	}

	list, httpRes, err := client.SiteAPI.ListDeviceAuthProviders(ctx).Search(idOrLabel).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for _, provider := range list.Data {
		if provider.Label == idOrLabel {
			result := provider
			return &result, nil
		}
	}

	err = fmt.Errorf("device auth provider '%s' not found", idOrLabel)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

func DeviceAuthProviderList(ctx context.Context, filterSiteId, filterKind, filterStatus []string) error {
	logger.Get().Info().Msgf("Listing all device auth providers")

	client := api.GetApiClient(ctx)

	request := client.SiteAPI.ListDeviceAuthProviders(ctx)
	if len(filterSiteId) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filterSiteId))
	}
	if len(filterKind) > 0 {
		request = request.FilterKind(utils.ProcessFilterStringSlice(filterKind))
	}
	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}

	list, httpRes, err := request.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(list, &DeviceAuthProviderPrintConfig)
}

func DeviceAuthProviderGet(ctx context.Context, idOrLabel string) error {
	logger.Get().Info().Msgf("Get device auth provider '%s'", idOrLabel)

	provider, err := GetDeviceAuthProviderByIdOrLabel(ctx, idOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(provider, &DeviceAuthProviderPrintConfig)
}

func DeviceAuthProviderCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating device auth provider")

	var createRequest sdk.CreateDeviceAuthProvider
	if err := utils.UnmarshalContent(config, &createRequest); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	provider, httpRes, err := client.SiteAPI.
		CreateDeviceAuthProvider(ctx).
		CreateDeviceAuthProvider(createRequest).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(provider, &DeviceAuthProviderPrintConfig)
}

func DeviceAuthProviderUpdate(ctx context.Context, idOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Updating device auth provider '%s'", idOrLabel)

	provider, err := GetDeviceAuthProviderByIdOrLabel(ctx, idOrLabel)
	if err != nil {
		return err
	}

	var updateRequest sdk.UpdateDeviceAuthProvider
	if err := utils.UnmarshalContent(config, &updateRequest); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updated, httpRes, err := client.SiteAPI.
		UpdateDeviceAuthProvider(ctx, int64(provider.Id)).
		UpdateDeviceAuthProvider(updateRequest).
		IfMatch(strconv.Itoa(int(provider.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updated, &DeviceAuthProviderPrintConfig)
}

func DeviceAuthProviderDelete(ctx context.Context, idOrLabel string) error {
	logger.Get().Info().Msgf("Deleting device auth provider '%s'", idOrLabel)

	provider, err := GetDeviceAuthProviderByIdOrLabel(ctx, idOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SiteAPI.
		DeleteDeviceAuthProvider(ctx, int64(provider.Id)).
		IfMatch(strconv.Itoa(int(provider.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Device auth provider '%s' deleted successfully", idOrLabel)
	return nil
}

func DeviceAuthProviderGetCredentials(ctx context.Context, idOrLabel string) error {
	logger.Get().Info().Msgf("Getting credentials for device auth provider '%s'", idOrLabel)

	provider, err := GetDeviceAuthProviderByIdOrLabel(ctx, idOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.SiteAPI.
		GetDeviceAuthProviderCredentials(ctx, int64(provider.Id)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &DeviceAuthProviderCredentialsPrintConfig)
}

func DeviceAuthProviderUpdateSharedSecret(ctx context.Context, idOrLabel string, sharedSecret string) error {
	logger.Get().Info().Msgf("Updating shared secret for device auth provider '%s'", idOrLabel)

	provider, err := GetDeviceAuthProviderByIdOrLabel(ctx, idOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SiteAPI.
		UpdateDeviceAuthProviderSharedSecret(ctx, int64(provider.Id)).
		UpdateDeviceAuthProviderSharedSecret(sdk.UpdateDeviceAuthProviderSharedSecret{SharedSecret: sharedSecret}).
		IfMatch(strconv.Itoa(int(provider.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Shared secret updated for device auth provider '%s'", idOrLabel)
	return nil
}

func DeviceAuthProviderConfigExample(ctx context.Context) error {
	example := sdk.CreateDeviceAuthProvider{
		Label:        "tacacs-primary",
		Name:         "Primary TACACS+ Server",
		SiteId:       1,
		Kind:         "tacacs",
		IpAddress:    "10.0.0.10",
		Port:         49,
		SharedSecret: "shared-secret",
		Username:     "admin",
		Password:     sdk.PtrString("password"),
		Status:       sdk.PtrString("active"),
		Annotations:  &map[string]string{"environment": "production"},
	}

	return formatter.PrintResult(example, nil)
}
