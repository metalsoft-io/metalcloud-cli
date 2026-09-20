// Package configuration exposes the global platform configuration document
// served by /api/v2/config. The document is a free-form, deeply nested map of
// service sections (auth, platform, notification, tunnel, gateway-api,
// image-builder), so it is handled as raw JSON rather than through the typed
// SDK models; the comment above ConfigurationGet explains why.
package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

// configurationBasePath is the platform configuration endpoint.
const configurationBasePath = "/api/v2/config"

// KnownServices lists the configuration sections (the {filter} path parameter
// of PUT/PATCH /api/v2/config/{filter} and the optional ?filter= query
// parameter of GET /api/v2/config). They are the top-level keys of the
// configuration document as modelled by sdk.Configurations.
var KnownServices = []string{
	"auth",
	"platform",
	"notification",
	"tunnel",
	"gateway-api",
	"image-builder",
}

// Why this package bypasses the typed SDK calls:
//
//   - ConfigurationAPI.ReplaceConfiguration / UpdateConfiguration return
//     *sdk.ReplaceConfigurationRequest, a oneOf union of the six service
//     configuration models. None of those models declares a discriminator and
//     all of them accept arbitrary additional properties, so every payload
//     matches all six branches and the generated UnmarshalJSON always fails
//     with "data matches more than one schema in oneOf(ReplaceConfigurationRequest)".
//   - ConfigurationAPI.GetConfiguration returns *sdk.Configurations, whose
//     nested sdk.TunnelConfiguration requires shared_secret, bdk and syslog;
//     a platform with a partially configured tunnel section makes the whole
//     GET fail with "no value given for required property ...".
//
// Raw JSON keeps both read and write lossless and schema-drift proof.

// ConfigurationGet prints the global platform configuration, optionally
// restricted to a single service section.
func ConfigurationGet(ctx context.Context, service string) error {
	if service == "" {
		logger.Get().Info().Msgf("Get platform configuration")
	} else {
		logger.Get().Info().Msgf("Get platform configuration for service '%s'", service)
	}

	path := configurationBasePath
	if service != "" {
		path = fmt.Sprintf("%s?filter=%s", configurationBasePath, url.QueryEscape(service))
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}

	return printConfigurationDocument(body)
}

// ConfigurationReplace replaces the whole configuration of a service section
// (PUT /api/v2/config/{filter}).
func ConfigurationReplace(ctx context.Context, service string, config []byte) error {
	logger.Get().Info().Msgf("Replacing platform configuration for service '%s'", service)

	return writeConfiguration(ctx, http.MethodPut, service, config)
}

// ConfigurationUpdate merges a partial configuration into a service section
// (PATCH /api/v2/config/{filter}).
func ConfigurationUpdate(ctx context.Context, service string, config []byte) error {
	logger.Get().Info().Msgf("Updating platform configuration for service '%s'", service)

	return writeConfiguration(ctx, http.MethodPatch, service, config)
}

func writeConfiguration(ctx context.Context, method string, service string, config []byte) error {
	service, err := ValidateService(service)
	if err != nil {
		return err
	}

	payload, err := normalizeConfigurationPayload(config)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", configurationBasePath, url.PathEscape(service))

	body, err := api.RawJSONRequest(ctx, method, path, payload, nil)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		logger.Get().Info().Msgf("Platform configuration for service '%s' updated", service)
		return nil
	}

	return printConfigurationDocument(body)
}

// ValidateService normalizes and checks the service (filter) name so that a
// typo cannot silently write to the wrong configuration section.
func ValidateService(service string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(service))
	if normalized == "" {
		err := fmt.Errorf("configuration service must be one of: %s", strings.Join(KnownServices, ", "))
		logger.Get().Error().Err(err).Msg("")
		return "", err
	}

	for _, known := range KnownServices {
		if normalized == known {
			return normalized, nil
		}
	}

	err := fmt.Errorf("unknown configuration service '%s', expected one of: %s", service, strings.Join(KnownServices, ", "))
	logger.Get().Error().Err(err).Msg("")
	return "", err
}

// normalizeConfigurationPayload accepts JSON or YAML input and returns the
// JSON body to send. Routing through utils.UnmarshalContent keeps YAML input
// behaving exactly like JSON input.
func normalizeConfigurationPayload(config []byte) ([]byte, error) {
	var document map[string]any
	if err := utils.UnmarshalContent(config, &document); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("failed to encode configuration: %w", err)
	}

	return payload, nil
}

// printConfigurationDocument renders the configuration document. json and yaml
// output keep the raw object; table formats would render the nested sections
// as unreadable Go maps, so they get YAML instead.
func printConfigurationDocument(body []byte) error {
	document, err := utils.DecodeRawObject(body)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(document, nil)
	}

	return formatter.PrintYamlResult(document)
}
