// Package email_template manages the notification email templates of the
// platform. Templates are addressed by name (not by numeric ID) on every
// endpoint of /api/v2/email-templates.
package email_template

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// emailTemplatePrintConfig deliberately omits the Text and Html bodies: they
// are multi-line documents that destroy table layout. Use json or yaml output
// (or 'email-template get') to read them.
var emailTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Title:    "Name",
			Order:    2,
			MaxWidth: 40,
		},
		"Subject": {
			Title:    "Subject",
			Order:    3,
			MaxWidth: 50,
		},
		"Description": {
			Title:    "Description",
			Order:    4,
			MaxWidth: 50,
		},
		"Revision": {
			Title: "Revision",
			Order: 5,
		},
	},
}

func EmailTemplateList(ctx context.Context, filterName []string) error {
	logger.Get().Info().Msgf("Listing email templates")

	client := api.GetApiClient(ctx)

	request := client.EmailTemplateAPI.GetEmailTemplates(ctx).SortBy([]string{"id:ASC"})
	if len(filterName) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filterName))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &emailTemplatePrintConfig)
}

func EmailTemplateGet(ctx context.Context, emailTemplateName string) error {
	logger.Get().Info().Msgf("Get email template '%s'", emailTemplateName)

	emailTemplate, err := GetEmailTemplateByName(ctx, emailTemplateName)
	if err != nil {
		return err
	}

	return formatter.PrintResult(emailTemplate, &emailTemplatePrintConfig)
}

func EmailTemplateConfigExample(ctx context.Context) error {
	example := sdk.EmailTemplateCreate{
		Name:        "infrastructure-deployed",
		Subject:     "Your infrastructure {{infrastructureLabel}} has been deployed",
		Description: sdk.PtrString("Sent when an infrastructure deployment completes"),
		Text:        "Hello {{userDisplayName}},\n\nInfrastructure {{infrastructureLabel}} has been deployed.\n",
		Html:        "<p>Hello {{userDisplayName}},</p><p>Infrastructure {{infrastructureLabel}} has been deployed.</p>",
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func EmailTemplateCreate(ctx context.Context, config []byte) error {
	var create sdk.EmailTemplateCreate
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Creating email template '%s'", create.Name)

	client := api.GetApiClient(ctx)

	emailTemplate, httpRes, err := client.EmailTemplateAPI.
		CreateEmailTemplate(ctx).
		EmailTemplateCreate(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(emailTemplate, &emailTemplatePrintConfig)
}

func EmailTemplateUpdate(ctx context.Context, emailTemplateName string, config []byte) error {
	logger.Get().Info().Msgf("Updating email template '%s'", emailTemplateName)

	var update sdk.EmailTemplateUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	// The PATCH endpoint takes an If-Match entity tag, so the current revision
	// has to be read first.
	current, err := GetEmailTemplateByName(ctx, emailTemplateName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	emailTemplate, httpRes, err := client.EmailTemplateAPI.
		EmailTemplatesControllerPatchEmailTemplate(ctx, current.Name).
		EmailTemplateUpdate(update).
		IfMatch(strconv.FormatInt(current.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(emailTemplate, &emailTemplatePrintConfig)
}

func EmailTemplateDelete(ctx context.Context, emailTemplateName string) error {
	logger.Get().Info().Msgf("Deleting email template '%s'", emailTemplateName)

	name, err := validateEmailTemplateName(emailTemplateName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EmailTemplateAPI.DeleteEmailTemplate(ctx, name).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Email template '%s' deleted", emailTemplateName)
	return nil
}

// GetEmailTemplateByName resolves a template by its name. Email templates have
// no by-ID endpoint: the name is the addressable key on every route.
func GetEmailTemplateByName(ctx context.Context, emailTemplateName string) (*sdk.EmailTemplate, error) {
	name, err := validateEmailTemplateName(emailTemplateName)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	emailTemplate, httpRes, err := client.EmailTemplateAPI.GetEmailTemplateByName(ctx, name).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	// A 2xx without a body would otherwise be dereferenced below.
	if emailTemplate == nil {
		err := fmt.Errorf("email template '%s' not found", emailTemplateName)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return emailTemplate, nil
}

// validateEmailTemplateName rejects an empty or whitespace-only name so that
// the CLI does not issue a request against /api/v2/email-templates/.
func validateEmailTemplateName(emailTemplateName string) (string, error) {
	name := strings.TrimSpace(emailTemplateName)
	if name == "" {
		err := fmt.Errorf("email template name cannot be empty")
		logger.Get().Error().Err(err).Msg("")
		return "", err
	}
	return name, nil
}
