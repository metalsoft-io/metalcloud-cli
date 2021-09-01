package main

import (
	"encoding/json"
	"io/ioutil"
	"syscall"
	"testing"

	"gopkg.in/yaml.v2"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestSwitchList(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := map[string]metalcloud.SwitchDevice{
		"sw1": {
			NetworkEquipmentID:               10,
			NetworkEquipmentIdentifierString: "test",
		},
	}

	client.EXPECT().
		SwitchDevices("", "").
		Return(&list, nil).
		AnyTimes()

	expectedFirstRow := map[string]interface{}{
		"ID":         10,
		"IDENTIFIER": "test",
	}

	testListCommand(switchListCmd, nil, client, expectedFirstRow, t)

}

func TestSwitchCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceCreate(gomock.Any(), false).
		Return(&sw, nil).
		AnyTimes()

	f, err := ioutil.TempFile("/tmp", "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_switchDeviceFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := ioutil.TempFile("/tmp", "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(sw)
	Expect(err).To(BeNil())

	f2.WriteString(string(s))
	f2.Close()
	defer syscall.Unlink(f2.Name())

	cases := []CommandTestCase{
		/*	{
				name: "sw-create-good1",
				cmd: MakeCommand(map[string]interface{}{
					"read_config_from_file": f.Name(),
					"format":                "json",
				}),
				good: true,
				id:   1,
			},
			{
				name: "sw-create-good-yaml",
				cmd: MakeCommand(map[string]interface{}{
					"read_config_from_file": f2.Name(),
					"format":                "yaml",
				}),
				good: true,
				id:   1,
			},*/
		{
			name: "sw-create-good-yaml",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": "examples/switch.yaml",
				"format":                "yaml",
			}),
			good: true,
			id:   1,
		},
	}

	testCreateCommand(switchCreateCmd, cases, client, t)

}

func TestSwitchGet(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var sw metalcloud.SwitchDevice

	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceGet(100, false).
		Return(&sw, nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "sw-get-json1",
			cmd: MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string": 100,
				"format":                                 "json",
			}),
			good: true,
			id:   1,
		},
		{
			name: "sw-get-json1",
			cmd: MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string": 100,
				"format":                                 "yaml",
			}),
			good: true,
			id:   1,
		},
	}

	expectedFirstRow := map[string]interface{}{
		"ID":         10,
		"IDENTIFIER": "test",
	}

	testGetCommand(switchGetCmd, cases, client, expectedFirstRow, t)

}

func TestSwitchDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	sw := metalcloud.SwitchDevice{
		NetworkEquipmentID: 100,
	}
	client.EXPECT().
		SwitchDeviceGet(100, false).
		Return(&sw, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceDelete(100).
		Return(nil).
		AnyTimes()

	cmd := MakeCommand(map[string]interface{}{
		"network_device_id_or_identifier_string": 100,
		"autoconfirm":                            true,
	})

	_, err := switchDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())
}

func TestSwitchObjectYamlUnmarshal(t *testing.T) {
	RegisterTestingT(t)

	var obj metalcloud.SwitchDevice

	cmd := MakeCommand(map[string]interface{}{
		"read_config_from_file": "examples/switch2.yaml",
		"format":                "yaml",
	})

	err := getRawObjectFromCommand(&cmd, &obj)
	Expect(err).To(BeNil())
	Expect(obj.NetworkEquipmentProvisionerPosition).To(Equal("other"))

	j, err := json.MarshalIndent(obj, "", "\t")
	Expect(err).To(BeNil())
	t.Log(string(j))
	Expect(string(j)).To(ContainSubstring("other"))

}

const _switchDeviceFixture1 = "{\"network_equipment_id\":1,\"datacenter_name\":\"uk-reading\",\"network_equipment_driver\":\"hp5900\",\"network_equipment_position\":\"tor\",\"network_equipment_provisioner_type\":\"vpls\",\"network_equipment_identifier_string\":\"UK_RDG_EVR01_00_0001_00A9_01\",\"network_equipment_description\":\"HP Comware Software, Version 7.1.045, Release 2311P06\",\"network_equipment_management_address\":\"10.0.0.0\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"sad\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"11.16.0.1\",\"network_equipment_quarantine_subnet_end\":\"11.16.0.00\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"11.16.0.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":\"11.24.0.2\",\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":\"100.64.0.0\",\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":1,\"network_equipment_primary_wan_ipv6_subnet_cidr\":\"2A02:0CB8:0000:0000:0000:0000:0000:0000/53\",\"network_equipment_cached_updated_timestamp\":\"2020-08-04T20:11:49Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":\"\",\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":1,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_tags\":[],\"network_equipment_management_password\":\"ddddd\"}"
