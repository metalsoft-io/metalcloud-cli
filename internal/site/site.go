package site

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

var sitePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Location": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Address": {
					Order: 3,
				},
			},
		},
		"IsHidden": {
			Title: "Hidden",
			Order: 4,
		},
		"IsInMaintenance": {
			Title: "Maintenance",
			Order: 5,
		},
	},
}

func GetSiteByIdOrLabel(ctx context.Context, siteIdOrLabel string) (*sdk.Site, error) {
	client := api.GetApiClient(ctx)

	siteList, httpRes, err := client.SiteAPI.GetSites(ctx).Search(siteIdOrLabel).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if len(siteList.Data) == 0 {
		err := fmt.Errorf("site '%s' not found", siteIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	var siteInfo sdk.Site
	for _, site := range siteList.Data {
		if site.Name == siteIdOrLabel {
			siteInfo = site
			break
		}

		if strconv.Itoa(int(site.Id)) == siteIdOrLabel {
			siteInfo = site
			break
		}
	}

	if siteInfo.Id == 0 {
		err := fmt.Errorf("site '%s' not found", siteIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return &siteInfo, nil
}

func SiteList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all sites")

	client := api.GetApiClient(ctx)

	siteList, httpRes, err := client.SiteAPI.GetSites(ctx).SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(siteList, &sitePrintConfig)
}

func SiteGet(ctx context.Context, siteIdOrName string) error {
	logger.Get().Info().Msgf("Get site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	return formatter.PrintResult(siteInfo, &sitePrintConfig)
}

func SiteCreate(ctx context.Context, siteName string) error {
	logger.Get().Info().Msgf("Create site '%s'", siteName)

	createSite := sdk.SiteCreate{
		Name: siteName,
		Slug: utils.CreateSlug(siteName),
	}

	client := api.GetApiClient(ctx)

	siteInfo, httpRes, err := client.SiteAPI.CreateSite(ctx).SiteCreate(createSite).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(siteInfo, &sitePrintConfig)
}

func SiteUpdate(ctx context.Context, siteIdOrName string, label string) error {
	logger.Get().Info().Msgf("Update site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	updateSite := sdk.SiteUpdate{
		Name: &label,
	}

	client := api.GetApiClient(ctx)

	siteInfo, httpRes, err := client.SiteAPI.UpdateSite(ctx, float32(siteInfo.Id)).
		SiteUpdate(updateSite).
		IfMatch(strconv.Itoa(int(siteInfo.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(siteInfo, &sitePrintConfig)
}

func SiteDecommission(ctx context.Context, siteIdOrName string) error {
	logger.Get().Info().Msgf("Delete site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SiteAPI.
		DecommissionSite(ctx, float32(siteInfo.Id)).
		IfMatch(strconv.Itoa(int(siteInfo.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func SiteGetAgents(ctx context.Context, siteIdOrName string) error {
	logger.Get().Info().Msgf("Get agents for site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	agents, httpRes, err := client.SiteAPI.GetAgents(ctx, float32(siteInfo.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(agents, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"SiteName": {
				Title: "Site",
				Order: 3,
			},
			"AgentType": {
				Title: "Agent Type",
				Order: 4,
			},
			"AgentVersion": {
				Title: "Version",
				Order: 5,
			},
			"AgentSeenIpAddress": {
				Title: "IP",
				Order: 6,
			},
			"AgentSeenTimestamp": {
				Title:       "Last Seen",
				Transformer: formatter.FormatDateTimeValue,
				Order:       7,
			},
			"AgentConnectedInfo": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"AgentId": {
						Title: "ID",
						Order: 1,
					},
					"Hostname": {
						Title: "Hostname",
						Order: 2,
					},
				},
			},
		},
	})
}
