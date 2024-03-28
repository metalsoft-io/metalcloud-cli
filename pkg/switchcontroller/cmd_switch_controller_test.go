package switchcontroller

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"gopkg.in/yaml.v2"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"

	"github.com/metalsoft-io/metalcloud-cli/internal/command"
)

func TestSwitchControllerList(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := map[int]metalcloud.SwitchDeviceController{
		10: {
			NetworkEquipmentControllerID:               10,
			NetworkEquipmentControllerIdentifierString: "test",
		},
	}

	client.EXPECT().
		SwitchDeviceControllers("").
		Return(&list, nil).
		AnyTimes()

	expectedFirstRow := map[string]interface{}{
		"ID":         10,
		"IDENTIFIER": "test",
	}

	command.TestListCommand(switchControllersListCmd, nil, client, expectedFirstRow, t)

}

func TestSwitchControllerCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var swCtrl metalcloud.SwitchDeviceController

	err := json.Unmarshal([]byte(_SwitchDeviceControllerFixture), &swCtrl)
	if err != nil {
		t.Error(err)
	}

	result := metalcloud.SwitchDeviceController{
		NetworkEquipmentControllerID:               10,
		NetworkEquipmentControllerIdentifierString: "test",
	}

	client.EXPECT().
		SwitchDeviceControllerCreate(gomock.Any()).
		Return(&result, nil).
		AnyTimes()

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_SwitchDeviceControllerFixture)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := os.CreateTemp(os.TempDir(), "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(swCtrl)
	Expect(err).To(BeNil())

	f2.WriteString(string(s))
	f2.Close()
	defer syscall.Unlink(f2.Name())

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")

	cases := []command.CommandTestCase{
		{
			Name: "sw-create-good-yaml",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": filepath.Join(basePath, "examples", "switch_controller.yaml"),
				"format":                "yaml",
			}),
			Good: true,
		},
	}

	command.TestCreateCommand(switchControllerCreateCmd, cases, client, t)
}

func TestSwitchControllerEditCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	swCtrl, err := getSwitchControllerFixture()
	if err != nil {
		t.Error(err)
	}

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_SwitchDeviceControllerFixture)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := os.CreateTemp(os.TempDir(), "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(swCtrl)
	Expect(err).To(BeNil())

	f2.WriteString(string(s))
	f2.Close()
	defer syscall.Unlink(f2.Name())

	client.EXPECT().
		SwitchDeviceControllerGet(310, false).
		Return(&swCtrl, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceControllerUpdate(swCtrl.NetworkEquipmentControllerID, gomock.Any()).
		Return(&swCtrl, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "missing-id",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f.Name(),
				"format":                "json",
			}),
			Good: false,
		},
		{
			Name: "get-from-json-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_controller_id_or_identifier_string": 310,
				"read_config_from_file":                      f.Name(),
				"format":                                     "json",
			}),
			Good: true,
		},
		{
			Name: "get-from-yaml-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_controller_id_or_identifier_string": 310,
				"read_config_from_file":                      f2.Name(),
				"format":                                     "yaml",
			}),
			Good: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := switchControllerEditCmd(&c.Cmd, client)
			if c.Good && err != nil {
				t.Error(err)
			}
		})
	}

}

func getSwitchControllerFixture() (metalcloud.SwitchDeviceController, error) {
	var sw metalcloud.SwitchDeviceController
	err := json.Unmarshal([]byte(_SwitchDeviceControllerFixture), &sw)
	return sw, err
}

func TestSwitchControllerGet(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	sw, err := getSwitchControllerFixture()
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceControllerGet(2, false).
		Return(&sw, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "sw-ctrl-get-json",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_controller_id_or_identifier_string": 2,
				"format": "json",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "sw-ctrl-get-yaml",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_controller_id_or_identifier_string": 2,
				"format": "yaml",
			}),
			Good: true,
			Id:   2,
		},
	}

	expectedFirstRow := map[string]interface{}{
		"ID":         2,
		"IDENTIFIER": "test",
	}

	command.TestGetCommand(switchControllerGetCmd, cases, client, expectedFirstRow, t)
}

func TestSwitchControllerDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	swCtrl := metalcloud.SwitchDeviceController{
		NetworkEquipmentControllerID: 100,
	}

	client.EXPECT().
		SwitchDeviceControllerGet(100, false).
		Return(&swCtrl, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceControllerDelete(100).
		Return(nil).
		AnyTimes()

	cmd := command.MakeCommand(map[string]interface{}{
		"network_controller_id_or_identifier_string": 100,
		"autoconfirm": true,
	})

	_, err := switchControllerDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())
}

func TestSwitchControllerSync(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	swCtrl := metalcloud.SwitchDeviceController{
		NetworkEquipmentControllerID: 100,
	}

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_SwitchDeviceFixture), &sw)
	if err != nil {
		t.Error(err)
	}

	result := map[int]metalcloud.SwitchDevice{
		sw.NetworkEquipmentID: sw,
	}

	client.EXPECT().
		SwitchDeviceControllerGet(100, false).
		Return(&swCtrl, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceControllerSync(100).
		Return(&result, nil).
		AnyTimes()

	cmd := command.MakeCommand(map[string]interface{}{
		"network_controller_id_or_identifier_string": 100,
	})

	_, err = switchControllerSyncCmd(&cmd, client)
	Expect(err).To(BeNil())
}

func TestSwitchObjectYamlUnmarshal(t *testing.T) {
	RegisterTestingT(t)

	var obj metalcloud.SwitchDevice

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")

	cmd := command.MakeCommand(map[string]interface{}{
		"read_config_from_file": filepath.Join(basePath, "examples", "switch_controller.yaml"),
		"format":                "yaml",
	})

	err := command.GetRawObjectFromCommand(&cmd, &obj)
	Expect(err).To(BeNil())
	Expect(obj.NetworkEquipmentProvisionerPosition).To(Equal("leaf"))

	j, err := json.MarshalIndent(obj, "", "\t")
	Expect(err).To(BeNil())
	t.Log(string(j))
	Expect(string(j)).To(ContainSubstring("leaf"))

}

const _SwitchDeviceFixture = "{\"network_equipment_id\":1,\"datacenter_name\":\"ro-bucharest\",\"network_equipment_driver\":\"cisco_aci51\",\"network_equipment_position\":\"leaf\",\"network_equipment_provisioner_type\":\"sdn\",\"network_equipment_identifier_string\":\"UK_RDG_EVR01_00_0001_00A9_01\",\"network_equipment_description\":\"HP Comware Software, Version 7.1.045, Release 2311P06\",\"network_equipment_management_address\":\"10.0.0.0\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"sad\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"11.16.0.1\",\"network_equipment_quarantine_subnet_end\":\"11.16.0.00\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"11.16.0.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":\"11.24.0.2\",\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":\"100.64.0.0\",\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":1,\"network_equipment_primary_wan_ipv6_subnet_cidr\":\"2A02:0CB8:0000:0000:0000:0000:0000:0000/53\",\"network_equipment_cached_updated_timestamp\":\"2020-08-04T20:11:49Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":\"\",\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":1,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_tags\":[],\"network_equipment_management_password\":\"ddddd\"}"
const _SwitchDeviceControllerFixture = "{\"datacenter_name\":\"br-sao-001\",\"network_equipment_controller_driver\":\"cisco_aci51\",\"network_equipment_controller_id\":2,\"network_equipment_controller_identifier_string\":\"Cisco ACI 5.1\",\"network_equipment_controller_management_address\": \"10.220.13.74\",\"network_equipment_controller_management_port\":22,\"network_equipment_controller_management_username\":\"metalsoft\",\"network_equipment_controller_options\":{\"vrf_shared_name\":\"VRF_SHARED_ADI\"},\"network_equipment_controller_provisioner_type\":\"sdn\"}"
