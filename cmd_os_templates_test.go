package main

import (
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestOSTemplatesListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	list := map[string]metalcloud.OSTemplate{
		"test": {
			VolumeTemplateID:        10,
			VolumeTemplateLabel:     "test",
			OSAssetBootloaderOSBoot: 100,
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		OSTemplates().
		Return(&list, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	client.EXPECT().
		OSAssetGet(list["test"].OSAssetBootloaderOSBoot).
		Return(&asset, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"ID":    10,
		"LABEL": "test",
	}

	testListCommand(templatesListCmd, client, expectedFirstRow, t)

}
