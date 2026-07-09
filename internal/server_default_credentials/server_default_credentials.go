package server_default_credentials

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)
type serverDefaultCredentialsRaw struct {
	Id                           interface{} `json:"id"`
	SiteId                       interface{} `json:"siteId"`
	ServerSerialNumber           *string     `json:"serverSerialNumber"`
	ServerMacAddress             *string     `json:"serverMacAddress"`
	DefaultUsername              *string     `json:"defaultUsername"`
	DefaultRackName              *string     `json:"defaultRackName"`
	DefaultRackPositionLowerUnit *string     `json:"defaultRackPositionLowerUnit"`
	DefaultRackPositionUpperUnit *string     `json:"defaultRackPositionUpperUnit"`
	DefaultInventoryId           *string     `json:"defaultInventoryId"`
	DefaultUuid                  *string     `json:"defaultUuid"`
}

var serverDefaultCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
		rawItems, meta, err := utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := req.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, page, limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[serverDefaultCredentialsRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse server default credentials: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &serverDefaultCredentialsPrintConfig)
	}

	if limit > 0 {
		rawItems, meta, err := utils.FetchUpToRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := req.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[serverDefaultCredentialsRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse server default credentials: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &serverDefaultCredentialsPrintConfig)
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
		_, httpRes, _ := req.Page(p).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[serverDefaultCredentialsRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse server default credentials: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &serverDefaultCredentialsPrintConfig)
}

func ServerDefaultCredentialsGet(ctx context.Context, credentialsId string) error {
	logger.Get().Info().Msgf("Get server default credentials '%s'", credentialsId)

	credentialsIdNumber, err := strconv.ParseInt(credentialsId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerDefaultCredentialsAPI.GetServerDefaultCredentialsInfo(ctx, credentialsIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &serverDefaultCredentialsPrintConfig)
}

func ServerDefaultCredentialsGetCredentials(ctx context.Context, credentialsId string) error {
	logger.Get().Info().Msgf("Get server default credentials unencrypted password for '%s'", credentialsId)

	credentialsIdNumber, err := strconv.ParseInt(credentialsId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerDefaultCredentialsAPI.GetServerDefaultCredentialsCredentials(ctx, credentialsIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func ServerDefaultCredentialsCreate(ctx context.Context, siteId float32, serialNumber string, macAddress string, username string, password string, rackName string, rackPositionLower string, rackPositionUpper string, inventoryId string, uuid string) error {
	logger.Get().Info().Msgf("Creating server default credentials for server '%s'", serialNumber)

	client := api.GetApiClient(ctx)

	createCredentials := sdk.NewCreateServerDefaultCredentials(int64(siteId), username, password)

	if serialNumber != "" {
		createCredentials.SetServerSerialNumber(serialNumber)
	}

	if macAddress != "" {
		createCredentials.SetServerMacAddress(macAddress)
	}

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

	credentialsIdNumber, err := strconv.ParseInt(credentialsId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server default credentials ID: '%s'", credentialsId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerDefaultCredentialsAPI.DeleteServerDefaultCredentials(ctx, credentialsIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server default credentials with ID %s deleted successfully", credentialsId)
	return nil
}
