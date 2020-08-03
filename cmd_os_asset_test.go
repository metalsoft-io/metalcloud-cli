package main

import (
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestAssetsListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	assetList := map[string]metalcloud.OSAsset{
		"test": {
			OSAssetID:       10,
			OSAssetUsage:    "test",
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		OSAssets().
		Return(&assetList, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"ID":    10,
		"USAGE": "test",
	}

	testListCommand(assetsListCmd, nil, client, expectedFirstRow, t)

}

func TestCreateAssetCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetCreate(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf1",
				"usage":                  "testf1",
				"read_content_from_pipe": true,
			}),
			good: true,
			id:   asset.OSAssetID,
		},
	}

	testCreateCommand(assetCreateCmd, cases, client, t)
}

func TestDeleteAssetCmd(t *testing.T) {
	RegisterTestingT(t)
	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetGet(asset.OSAssetID).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetDelete(asset.OSAssetID).
		Return(nil).
		MinTimes(1)

	cmd := MakeCommand(map[string]interface{}{"asset_id_or_name": asset.OSAssetID})
	testCommandWithConfirmation(assetDeleteCmd, cmd, client, t)
}
