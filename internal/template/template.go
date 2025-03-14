package template

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func GetTemplateByIdOrLabel(ctx context.Context, templateIdOrLabel string) (*sdk.OSTemplate, error) {
	client := api.GetApiClient(ctx)

	templateId, err := utils.GetFloat32FromString(templateIdOrLabel)
	if err != nil {
		return nil, err
	}

	templateInfo, httpRes, err := client.OSTemplateAPI.GetOSTemplate(ctx, templateId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return templateInfo, nil
}
