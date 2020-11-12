package main

import (
	"fmt"
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
			OSAssetID:    10,
			OSAssetUsage: "test",
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

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetCreate(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

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
		{
			name: "good2, associate a template (id)",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf2",
				"usage":                  "testf2",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
				"path":                   "test2",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good3, associate a template",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf3",
				"usage":                  "testf3",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test3",
				"variables_json":         "['1': 'test']",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good4, associate a template (name)",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf4",
				"usage":                  "testf4",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test4",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, non-existant template",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf5",
				"usage":                  "testf5",
				"read_content_from_pipe": true,
				"template_id_or_name":    "tmpl1",
				"path":                   "test5",
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, missing path",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf6",
				"usage":                  "testf6",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
			}),
			good: false,
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

func TestEditAssetCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	assetf := metalcloud.OSAsset{
		OSAssetID: 101,
	}

	assetl := map[string]metalcloud.OSAsset{
		"1": asset,
	}

	client.EXPECT().
		OSAssets().
		Return(&assetl, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetGet(asset.OSAssetID).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetGet(asset.OSAssetFileName).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetGet(assetf.OSAssetID).
		Return(nil, fmt.Errorf("test")).
		Times(1)

	client.EXPECT().
		OSAssetUpdate(gomock.Any(), gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":        100,
				"filename":                "testf1",
				"usage":                   "testf1",
				"read_content_from_pipe":  true,
				"variable_names_required": "1,2,3",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good2, associate a template (id)",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf2",
				"usage":                  "testf2",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
				"path":                   "test2",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good3, associate a template",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf3",
				"usage":                  "testf3",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test3",
				"variables_json":         "['1': 'test']",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good4, associate a template (name)",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf4",
				"usage":                  "testf4",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test4",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, non-existant template",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf5",
				"usage":                  "testf5",
				"read_content_from_pipe": true,
				"template_id_or_name":    "tmpl1",
				"path":                   "test5",
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, missing path",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf6",
				"usage":                  "testf6",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":        101,
				"filename":                "testf1",
				"usage":                   "testf1",
				"read_content_from_pipe":  true,
				"variable_names_required": "1,2,3",
			}),
			good: false,
			id:   0,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":        "file",
				"filename":                "testf1",
				"usage":                   "testf1",
				"read_content_from_pipe":  true,
				"variable_names_required": "1,2,3",
			}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(assetEditCmd, cases, client, t)
}
