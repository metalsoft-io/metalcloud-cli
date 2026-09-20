package network_device

import (
	"context"
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var networkDeviceSecretNamesPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Names": {
			Title:       "Secret Names",
			MaxWidth:    80,
			Transformer: formatNetworkDeviceStringList,
			Order:       1,
		},
	},
}

var networkDeviceSecretCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"SecretValue": {
			Title: "Secret Value",
			Order: 1,
		},
	},
}

// formatNetworkDeviceStringList renders a list of strings as a comma separated
// value so it fits into a single table cell.
func formatNetworkDeviceStringList(value interface{}) string {
	switch list := value.(type) {
	case []string:
		return strings.Join(list, ", ")
	case []interface{}:
		parts := make([]string, 0, len(list))
		for _, item := range list {
			parts = append(parts, fmt.Sprint(item))
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprint(value)
	}
}

// NetworkDeviceSecretList lists the names of the secrets stored for a network
// device. Values are never returned by this endpoint.
func NetworkDeviceSecretList(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Listing secrets of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secretNames, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceSecrets(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secretNames, &networkDeviceSecretNamesPrintConfig)
}

// NetworkDeviceSecretSet stores (or replaces) one named secret of a network
// device. The value is stored encrypted.
func NetworkDeviceSecretSet(ctx context.Context, networkDeviceRef string, name string, secret string) error {
	logger.Get().Info().Msgf("Setting secret '%s' of network device '%s'", name, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secretNames, httpRes, err := client.NetworkDeviceAPI.
		SetNetworkDeviceSecret(ctx, networkDeviceIdNumeric).
		SetNetworkDeviceSecret(sdk.SetNetworkDeviceSecret{Name: name, Secret: secret}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secretNames, &networkDeviceSecretNamesPrintConfig)
}

// NetworkDeviceSecretGetCredentials reveals the unencrypted value of one named
// secret of a network device.
func NetworkDeviceSecretGetCredentials(ctx context.Context, networkDeviceRef string, name string) error {
	logger.Get().Info().Msgf("Getting credentials of secret '%s' of network device '%s'", name, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceSecretCredentials(ctx, networkDeviceIdNumeric, name).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &networkDeviceSecretCredentialsPrintConfig)
}

// NetworkDeviceSecretRemove deletes one named secret of a network device.
func NetworkDeviceSecretRemove(ctx context.Context, networkDeviceRef string, name string) error {
	logger.Get().Info().Msgf("Removing secret '%s' of network device '%s'", name, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DeleteNetworkDeviceSecret(ctx, networkDeviceIdNumeric, name).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Secret '%s' of network device '%s' removed", name, networkDeviceRef)
	return nil
}

// NetworkDeviceSecretRemoveAll deletes every secret stored for a network device.
func NetworkDeviceSecretRemoveAll(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Removing all secrets of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DeleteNetworkDeviceSecrets(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("All secrets of network device '%s' removed", networkDeviceRef)
	return nil
}
