package server_default_credentials

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverDefaultCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"SiteId": {
			Title: "Site ID",
			Order: 2,
		},
		"ServerSerialNumber": {
			Title:    "Serial Number",
			MaxWidth: 30,
			Order:    3,
		},
		"ServerMacAddress": {
			Title: "MAC Address",
			Order: 4,
		},
		"DefaultUsername": {
			Title: "Username",
			Order: 5,
		},
		"DefaultRackName": {
			Title: "Rack Name",
			Order: 6,
		},
		"DefaultRackPositionLowerUnit": {
			Title: "Rack Lower Unit",
			Order: 7,
		},
		"DefaultRackPositionUpperUnit": {
			Title: "Rack Upper Unit",
			Order: 8,
		},
		"DefaultInventoryId": {
			Title: "Inventory ID",
			Order: 9,
		},
		"DefaultUuid": {
			Title: "UUID",
			Order: 10,
		},
	},
}

func ServerDefaultCredentialsList(ctx context.Context, page int, limit int) error {
	logger.Get().Info().Msgf("Listing server default credentials")

	client := api.GetApiClient(ctx)

	req := client.ServerDefaultCredentialsAPI.GetServersDefaultCredentials(ctx)

	if page > 0 {
		req = req.Page(float32(page))
	}

	if limit > 0 {
		req = req.Limit(float32(limit))
	}

	credentialsList, httpRes, err := req.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentialsList.Data, &serverDefaultCredentialsPrintConfig)
}

func ServerDefaultCredentialsGet(ctx context.Context, credentialsId string) error {
	logger.Get().Info().Msgf("Get server default credentials '%s'", credentialsId)

	credentialsIdNumber, err := strconv.ParseFloat(credentialsId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerDefaultCredentialsAPI.GetServerDefaultCredentialsInfo(ctx, float32(credentialsIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &serverDefaultCredentialsPrintConfig)
}

func ServerDefaultCredentialsGetCredentials(ctx context.Context, credentialsId string) error {
	logger.Get().Info().Msgf("Get server default credentials unencrypted password for '%s'", credentialsId)

	credentialsIdNumber, err := strconv.ParseFloat(credentialsId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerDefaultCredentialsAPI.GetServerDefaultCredentialsCredentials(ctx, float32(credentialsIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func ServerDefaultCredentialsCreate(ctx context.Context, siteId float32, serialNumber string, macAddress string, username string, password string, rackName string, rackPositionLower string, rackPositionUpper string, inventoryId string, uuid string) error {
	logger.Get().Info().Msgf("Creating server default credentials for server '%s'", serialNumber)

	client := api.GetApiClient(ctx)

	createCredentials := sdk.NewCreateServerDefaultCredentials(siteId, serialNumber, macAddress, username, password)

	if rackName != "" {
		createCredentials.SetDefaultRackName(rackName)
	}

	if rackPositionLower != "" {
		createCredentials.SetDefaultRackPositionLowerUnit(rackPositionLower)
	}

	if rackPositionUpper != "" {
		createCredentials.SetDefaultRackPositionUpperUnit(rackPositionUpper)
	}

	if inventoryId != "" {
		createCredentials.SetDefaultInventoryId(inventoryId)
	}

	if uuid != "" {
		createCredentials.SetDefaultUuid(uuid)
	}

	credentials, httpRes, err := client.ServerDefaultCredentialsAPI.CreateServerDefaultCredentials(ctx).
		CreateServerDefaultCredentials(*createCredentials).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &serverDefaultCredentialsPrintConfig)
}

func ServerDefaultCredentialsDelete(ctx context.Context, credentialsId string) error {
	logger.Get().Info().Msgf("Deleting server default credentials '%s'", credentialsId)

	credentialsIdNumber, err := strconv.ParseFloat(credentialsId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerDefaultCredentialsAPI.DeleteServerDefaultCredentials(ctx, float32(credentialsIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server default credentials with ID %s deleted successfully", credentialsId)
	return nil
}
