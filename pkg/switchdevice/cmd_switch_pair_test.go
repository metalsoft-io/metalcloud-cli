package switchdevice

import (
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"

	"github.com/metalsoft-io/metalcloud-cli/internal/command"
)

func TestSwitchPairList(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := map[int]metalcloud.SwitchDeviceLink{
		1: {
			NetworkEquipmentLinkID:   10,
			NetworkEquipmentID1:      10,
			NetworkEquipmentID2:      10,
			NetworkEquipmentLinkType: "mlag",
		},
		2: {
			NetworkEquipmentLinkID:   10,
			NetworkEquipmentID1:      10,
			NetworkEquipmentID2:      10,
			NetworkEquipmentLinkType: "mlag",
		},
	}

	sw := metalcloud.SwitchDevice{
		NetworkEquipmentID:               10,
		NetworkEquipmentIdentifierString: "sw1",
	}

	client.EXPECT().
		SwitchDeviceLinks().
		Return(&list, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceGet(10, false).
		Return(&sw, nil).
		AnyTimes()

	expectedFirstRow := map[string]interface{}{
		"ID":      10,
		"Switch1": "sw1 (#10)",
	}

	command.TestListCommand(switchPairListCmd, nil, client, expectedFirstRow, t)

}

func TestSwitchPairCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var swl metalcloud.SwitchDeviceLink

	err := json.Unmarshal([]byte(_switchDeviceLinkFixture1), &swl)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceLinkCreate(gomock.Any(), gomock.Any(), "mlag").
		Return(&swl, nil).
		AnyTimes()

	sw1, err := getSwitchFixture1()

	client.EXPECT().
		SwitchDeviceGet(7, false).
		Return(&sw1, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceGet(8, false).
		Return(&sw1, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "sw-link-create-good-yaml",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string1": 7,
				"network_device_id_or_identifier_string2": 8,
				"type":      "mlag",
				"return_id": true,
			}),
			Good: true,
			Id:   7,
		},
	}

	command.TestCreateCommand(switchPairCreateCmd, cases, client, t)

}

func TestSwitchPairDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var swl metalcloud.SwitchDeviceLink

	err := json.Unmarshal([]byte(_switchDeviceLinkFixture1), &swl)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceLinkGet(1, 1, "mlag").
		Return(&swl, nil).
		AnyTimes()

	sw1, err := getSwitchFixture1()

	client.EXPECT().
		SwitchDeviceGet(7, false).
		Return(&sw1, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceGet(8, false).
		Return(&sw1, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceLinkDelete(1, 1, "mlag").
		Return(nil).
		AnyTimes()

	cmd := command.MakeCommand(map[string]interface{}{
		"network_device_id_or_identifier_string1": 7,
		"network_device_id_or_identifier_string2": 8,
		"type":        "mlag",
		"autoconfirm": true,
	})

	_, err = switchPairDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())
}

const _switchDeviceLinkFixture1 = "{\"network_equipment_link_id\": 7,\"network_equipment_id_1\": 7,\"network_equipment_id_2\": 8,\"network_equipment_link_type\": \"mlag\",\"network_equipment_link_properties\": []}"
