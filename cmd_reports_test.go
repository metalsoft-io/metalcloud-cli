package main

import (
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestDevicesListCmdWithWatch(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var serverList []metalcloud.ServerSearchResult
	json.Unmarshal([]byte(_serverListFixture1), &serverList)

	client.EXPECT().
		ServersSearch("datacenter_name:test").
		Return(&serverList, nil).
		AnyTimes()

	var switchList map[string]metalcloud.SwitchDevice
	json.Unmarshal([]byte(_switchDeviceListFixture1), &switchList)

	client.EXPECT().
		SwitchDevices("test", "").
		Return(&switchList, nil).
		AnyTimes()

	var storageList []metalcloud.StoragePoolSearchResult
	json.Unmarshal([]byte(_storageListFixture), &storageList)

	client.EXPECT().
		StoragePoolSearch("datacenter_name:test").
		Return(&storageList, nil).
		AnyTimes()

	var dcList map[string]metalcloud.Datacenter
	json.Unmarshal([]byte(_datacenterList), &dcList)

	client.EXPECT().
		Datacenters(true).
		Return(&dcList, nil).
		AnyTimes()

	cmd := MakeCommand(map[string]interface{}{})

	ret, err := devicesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("2"))
}

const _storageListFixture = "[\r\n                {\r\n                    \"storage_pool_id\": 1,\r\n                    \"storage_pool_name\": \"UnityVSA\",\r\n                    \"storage_pool_status\": \"active\",\r\n                    \"storage_pool_in_maintenance\": false,\r\n                    \"datacenter_name\": \"us02-chi-qts01-dc\",\r\n                    \"storage_type\": \"iscsi_ssd\",\r\n                    \"user_id\": null,\r\n                    \"storage_pool_iscsi_host\": \"100.96.0.2\",\r\n                    \"storage_pool_iscsi_port\": 3260,\r\n                    \"storage_pool_capacity_total_cached_real_mbytes\": 505344,\r\n                    \"storage_pool_capacity_usable_cached_real_mbytes\": 505344,\r\n                    \"storage_pool_capacity_free_cached_real_mbytes\": 496128,\r\n                    \"storage_pool_capacity_used_cached_virtual_mbytes\": 122880\r\n                }\r\n            ]"
const _datacenterList = "{\"test\":{\"datacenter_id\":6,\"datacenter_name\":\"test\",\"datacenter_name_parent\":null,\"user_id\":null,\"datacenter_is_master\":false,\"datacenter_is_maintenance\":false,\"datacenter_type\":\"metal_cloud\",\"datacenter_display_name\":\"US02 Chi QTS01 DC\",\"datacenter_hidden\":false,\"datacenter_created_timestamp\":\"2022-02-11T11:14:08Z\",\"datacenter_updated_timestamp\":\"2022-06-09T13:32:56Z\",\"type\":\"Datacenter\",\"datacenter_tags\":[]}}"
