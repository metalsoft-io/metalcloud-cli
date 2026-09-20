package server

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// ServerImportUnmanaged registers a server that MetalSoft does not manage the
// lifecycle of, using an externally supplied hardware description.
func ServerImportUnmanaged(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Importing unmanaged server")

	var importConfig sdk.ServerUnmanagedImport
	if err := utils.UnmarshalContent(config, &importConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The typed Server model fails to decode valid responses (ServerMetricsMetadata
	// is typed as a map but returned as an array), so the body is parsed raw.
	_, httpRes, sdkErr := client.UnmanagedServersAPI.
		ImportUnmanagedServer(ctx).
		ServerUnmanagedImport(importConfig).
		Execute()

	serverInfo, err := parseServerRaw(httpRes, sdkErr)
	if err != nil {
		return err
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}

func ServerImportUnmanagedConfigExample(ctx context.Context) error {
	importConfig := sdk.ServerUnmanagedImport{
		SiteId:                        1,
		ServerTypeId:                  2,
		ManagementAddress:             sdk.PtrString("10.0.0.1"),
		ServerSupportsOobProvisioning: sdk.PtrFloat32(0),
		ServerSerialNumber:            sdk.PtrString("ABC1234"),
		ServerUUID:                    sdk.PtrString("00000000-0000-0000-0000-000000000000"),
		Hostname:                      sdk.PtrString("server-01"),
		ServerIpmiUsername:            sdk.PtrString("admin"),
		ServerIpmiPassword:            sdk.PtrString("password"),
		ServerInterfaces: []sdk.ServerUnmanagedImportInternalInterface{
			{
				ServerInterfaceMacAddress:                 "AA:BB:CC:DD:EE:00",
				IdentifierString:                          "eth0",
				NetworkEquipmentInterfaceIdentifierString: "Ethernet0",
			},
		},
	}

	return formatter.PrintResult(importConfig, nil)
}
