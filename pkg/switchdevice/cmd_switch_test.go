package switchdevice

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

	command.TestListCommand(switchListCmd, nil, client, expectedFirstRow, t)

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

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_switchDeviceFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := os.CreateTemp(os.TempDir(), "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(sw)
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
				"read_config_from_file": filepath.Join(basePath, "examples", "switch.yaml"),
				"format":                "yaml",
			}),
			Good: true,
			Id:   1,
		},
	}

	command.TestCreateCommand(switchCreateCmd, cases, client, t)

}

func TestSwitchUpdateCmd(t *testing.T) {

	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	sw, err := getSwitchFixture1()
	if err != nil {
		t.Error(err)
	}

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_switchDeviceFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := os.CreateTemp(os.TempDir(), "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(sw)
	Expect(err).To(BeNil())

	f2.WriteString(string(s))
	f2.Close()
	defer syscall.Unlink(f2.Name())

	client.EXPECT().
		SwitchDeviceGet(310, false).
		Return(&sw, nil).
		AnyTimes()

	client.EXPECT().
		SwitchDeviceUpdate(sw.NetworkEquipmentID, gomock.Any(), false).
		Return(&sw, nil).
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
				"network_device_id_or_identifier_string": 310,
				"read_config_from_file":                  f.Name(),
				"format":                                 "json",
			}),
			Good: true,
		},
		{
			Name: "get-from-yaml-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string": 310,
				"read_config_from_file":                  f2.Name(),
				"format":                                 "yaml",
			}),
			Good: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := switchUpdateCmd(&c.Cmd, client)
			if c.Good && err != nil {
				t.Error(err)
			}
		})
	}

}

func getSwitchFixture1() (metalcloud.SwitchDevice, error) {
	var sw metalcloud.SwitchDevice
	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	return sw, err
}

func TestSwitchGet(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	sw, err := getSwitchFixture1()
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceGet(100, false).
		Return(&sw, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "sw-get-json1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string": 100,
				"format":                                 "json",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "sw-get-json1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_device_id_or_identifier_string": 100,
				"format":                                 "yaml",
			}),
			Good: true,
			Id:   1,
		},
	}

	expectedFirstRow := map[string]interface{}{
		"ID":         10,
		"IDENTIFIER": "test",
	}

	command.TestGetCommand(switchGetCmd, cases, client, expectedFirstRow, t)

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

	cmd := command.MakeCommand(map[string]interface{}{
		"network_device_id_or_identifier_string": 100,
		"autoconfirm":                            true,
	})

	_, err := switchDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())
}

func TestSwitchObjectYamlUnmarshal(t *testing.T) {
	RegisterTestingT(t)

	var obj metalcloud.SwitchDevice

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")

	cmd := command.MakeCommand(map[string]interface{}{
		"read_config_from_file": filepath.Join(basePath, "examples", "switch2.yaml"),
		"format":                "yaml",
	})

	err := command.GetRawObjectFromCommand(&cmd, &obj)
	Expect(err).To(BeNil())
	Expect(obj.NetworkEquipmentProvisionerPosition).To(Equal("other"))

	j, err := json.MarshalIndent(obj, "", "\t")
	Expect(err).To(BeNil())
	t.Log(string(j))
	Expect(string(j)).To(ContainSubstring("other"))

}

const _switchDeviceFixture1 = "{\"network_equipment_id\":1,\"datacenter_name\":\"uk-reading\",\"network_equipment_driver\":\"hp5900\",\"network_equipment_position\":\"tor\",\"network_equipment_provisioner_type\":\"vpls\",\"network_equipment_identifier_string\":\"UK_RDG_EVR01_00_0001_00A9_01\",\"network_equipment_description\":\"HP Comware Software, Version 7.1.045, Release 2311P06\",\"network_equipment_management_address\":\"10.0.0.0\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"sad\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"11.16.0.1\",\"network_equipment_quarantine_subnet_end\":\"11.16.0.00\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"11.16.0.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":\"11.24.0.2\",\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":\"100.64.0.0\",\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":1,\"network_equipment_primary_wan_ipv6_subnet_cidr\":\"2A02:0CB8:0000:0000:0000:0000:0000:0000/53\",\"network_equipment_cached_updated_timestamp\":\"2020-08-04T20:11:49Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":\"\",\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":1,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_tags\":[],\"network_equipment_management_password\":\"ddddd\"}"
const _switchDeviceListFixture1 = "{\"4\":{\"network_equipment_id\":4,\"network_equipment_status\":\"active\",\"datacenter_name\":\"us02-chi-qts01-dc\",\"network_equipment_driver\":\"os_10\",\"network_equipment_position\":\"leaf\",\"network_equipment_provisioner_type\":\"evpnvxlanl2\",\"network_equipment_identifier_string\":\"sw1-env2\",\"network_equipment_description\":\"OS10 Enterprise. OS Version\",\"network_equipment_management_address\":\"10.0.5.6\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"admin\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"192.168.254.0\",\"network_equipment_quarantine_subnet_end\":\"192.168.254.255\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"192.168.254.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":null,\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":null,\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":3,\"network_equipment_primary_wan_ipv6_subnet_cidr\":null,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_cached_updated_timestamp\":\"0000-00-00T00:00:00Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":null,\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_tags_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":true,\"network_equipment_interfaces_blacklist_json\":null,\"network_equipment_controller_id\":null,\"network_equipment_is_border_device\":false,\"network_equipment_is_storage_switch\":true,\"network_equipment_network_types_allowed_json\":\"[\\\"wan\\\", \\\"quarantine\\\", \\\"san\\\"]\",\"network_equipment_order_index\":10,\"network_equipment_management_password\":\"Use bsidev.password_decrypt:eyJycWkiOiJici5BS3RET0E1Tm4tQjA1eUZ1TDdLMTZXeGhpcmF6UURWYzRySktWRko3Nzd1QWwzODkxR0E4ZmNqdVJUajRqRjFnZ1VVVElad2RWZkcyMm1aSXBVVmpMX2RDSWhaLXBYb1hVU0FJTDNkNGhjRGNEWjhOeDRQVE9CVk94V2VKYkRILXQ1ZEdxbkF5RzYzQ0NlOHZiT0JxTFEiLCJ2IjoieGZ2dnZPbFQxZEZrZW5FR2pSY0E1QSJ9\"},\"5\":{\"network_equipment_id\":5,\"network_equipment_status\":\"active\",\"datacenter_name\":\"us02-chi-qts01-dc\",\"network_equipment_driver\":\"os_10\",\"network_equipment_position\":\"leaf\",\"network_equipment_provisioner_type\":\"evpnvxlanl2\",\"network_equipment_identifier_string\":\"sw2-env2\",\"network_equipment_description\":\"OS10 Enterprise. OS Version\",\"network_equipment_management_address\":\"10.0.5.7\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"admin\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"192.168.254.0\",\"network_equipment_quarantine_subnet_end\":\"192.168.254.255\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"192.168.254.2\",\"network_equipment_primary_wan_ipv4_subnet_pool\":null,\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":null,\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":3,\"network_equipment_primary_wan_ipv6_subnet_cidr\":null,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_cached_updated_timestamp\":\"0000-00-00T00:00:00Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":null,\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_tags_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":true,\"network_equipment_interfaces_blacklist_json\":null,\"network_equipment_controller_id\":null,\"network_equipment_is_border_device\":false,\"network_equipment_is_storage_switch\":true,\"network_equipment_network_types_allowed_json\":\"[\\\"wan\\\", \\\"quarantine\\\", \\\"san\\\"]\",\"network_equipment_order_index\":20,\"network_equipment_management_password\":\"Use bsidev.password_decrypt:eyJycWkiOiJici4ybUpyYUNLZXlTYktKWDFvYXFPcnY1SzhMSWFQSWxLQ04xeTRTbFk4blJjdUs0WmV5M3gySU9JejcxTFprTk83MDg3VDRzNElsbnd1MjdMVFNXMm5iYkhaR3gzbGQwWWdjTnZNcTRjMlVCeWZZQW93YnQ3QU5uT09KQVdhRjk4RWZfaXhLNkFSRWFteWpFNloxbFFYSXciLCJ2IjoiT2Nxakhza2M0UEt3ZkRVaUdnekl2USJ9\"}}"
