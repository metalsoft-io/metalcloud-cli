package reports

import (
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"

	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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

	cmd := command.MakeCommand(map[string]interface{}{})

	ret, err := devicesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("2"))
}

const _storageListFixture = "[\r\n                {\r\n                    \"storage_pool_id\": 1,\r\n                    \"storage_pool_name\": \"UnityVSA\",\r\n                    \"storage_pool_status\": \"active\",\r\n                    \"storage_pool_in_maintenance\": false,\r\n                    \"datacenter_name\": \"us02-chi-qts01-dc\",\r\n                    \"storage_type\": \"iscsi_ssd\",\r\n                    \"user_id\": null,\r\n                    \"storage_pool_iscsi_host\": \"100.96.0.2\",\r\n                    \"storage_pool_iscsi_port\": 3260,\r\n                    \"storage_pool_capacity_total_cached_real_mbytes\": 505344,\r\n                    \"storage_pool_capacity_usable_cached_real_mbytes\": 505344,\r\n                    \"storage_pool_capacity_free_cached_real_mbytes\": 496128,\r\n                    \"storage_pool_capacity_used_cached_virtual_mbytes\": 122880\r\n                }\r\n            ]"
const _datacenterList = "{\"test\":{\"datacenter_id\":6,\"datacenter_name\":\"test\",\"datacenter_name_parent\":null,\"user_id\":null,\"datacenter_is_master\":false,\"datacenter_is_maintenance\":false,\"datacenter_type\":\"metal_cloud\",\"datacenter_display_name\":\"US02 Chi QTS01 DC\",\"datacenter_hidden\":false,\"datacenter_created_timestamp\":\"2022-02-11T11:14:08Z\",\"datacenter_updated_timestamp\":\"2022-06-09T13:32:56Z\",\"type\":\"Datacenter\",\"datacenter_tags\":[]}}"
const _serverListFixture1 = "[\n                {\n                    \"server_id\": 16,\n                    \"server_type_name\": null,\n                    \"server_type_boot_type\": null,\n                    \"server_product_name\": null,\n                    \"datacenter_name\": \"us02-chi-qts01-dc\",\n                    \"server_status\": \"registering\",\n                    \"server_class\": \"bigdata\",\n                    \"server_created_timestamp\": \"2022-05-23T13:22:11Z\",\n                    \"server_vendor\": \"Dell Inc.\",\n                    \"server_serial_number\": null,\n                    \"server_uuid\": \"4c4c4544-0051-3810-8057-b7c04f533532\",\n                    \"server_vendor_sku_id\": null,\n                    \"server_boot_type\": \"classic\",\n                    \"server_allocation_timestamp\": null,\n                    \"instance_label\": [\n                        null\n                    ],\n                    \"instance_id\": [\n                        null\n                    ],\n                    \"instance_array_id\": [\n                        null\n                    ],\n                    \"infrastructure_id\": [\n                        null\n                    ],\n                    \"server_inventory_id\": null,\n                    \"server_rack_name\": null,\n                    \"server_rack_position_lower_unit\": null,\n                    \"server_rack_position_upper_unit\": null,\n                    \"server_ipmi_host\": \"172.18.44.42\",\n                    \"server_ipmi_internal_username\": \"root\",\n                    \"server_processor_name\": null,\n                    \"server_processor_count\": 0,\n                    \"server_processor_core_count\": 0,\n                    \"server_processor_core_mhz\": 0,\n                    \"server_processor_threads\": null,\n                    \"server_processor_cpu_mark\": null,\n                    \"server_disk_type\": \"none\",\n                    \"server_disk_count\": 0,\n                    \"server_disk_size_mbytes\": 0,\n                    \"server_ram_gbytes\": 0,\n                    \"server_network_total_capacity_mbps\": 0,\n                    \"server_dhcp_status\": \"quarantine\",\n                    \"server_dhcp_packet_sniffing_is_enabled\": true,\n                    \"server_dhcp_relay_security_is_enabled\": true,\n                    \"server_disk_wipe\": false,\n                    \"server_power_status\": \"off\",\n                    \"server_power_status_last_update_timestamp\": \"2022-05-23T13:24:41Z\",\n                    \"user_id\": [\n                        [\n                            null\n                        ]\n                    ],\n                    \"user_id_owner\": [\n                        null\n                    ],\n                    \"user_email\": [\n                        [\n                            null\n                        ]\n                    ],\n                    \"infrastructure_user_id\": [\n                        [\n                            null\n                        ]\n                    ]\n                }\n            ]"
const _switchDeviceListFixture1 = "{\"4\":{\"network_equipment_id\":4,\"network_equipment_status\":\"active\",\"datacenter_name\":\"us02-chi-qts01-dc\",\"network_equipment_driver\":\"os_10\",\"network_equipment_position\":\"leaf\",\"network_equipment_provisioner_type\":\"evpnvxlanl2\",\"network_equipment_identifier_string\":\"sw1-env2\",\"network_equipment_description\":\"OS10 Enterprise. OS Version\",\"network_equipment_management_address\":\"10.0.5.6\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"admin\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"192.168.254.0\",\"network_equipment_quarantine_subnet_end\":\"192.168.254.255\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"192.168.254.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":null,\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":null,\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":3,\"network_equipment_primary_wan_ipv6_subnet_cidr\":null,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_cached_updated_timestamp\":\"0000-00-00T00:00:00Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":null,\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_tags_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":true,\"network_equipment_interfaces_blacklist_json\":null,\"network_equipment_controller_id\":null,\"network_equipment_is_border_device\":false,\"network_equipment_is_storage_switch\":true,\"network_equipment_network_types_allowed_json\":\"[\\\"wan\\\", \\\"quarantine\\\", \\\"san\\\"]\",\"network_equipment_order_index\":10,\"network_equipment_management_password\":\"Use bsidev.password_decrypt:eyJycWkiOiJici5BS3RET0E1Tm4tQjA1eUZ1TDdLMTZXeGhpcmF6UURWYzRySktWRko3Nzd1QWwzODkxR0E4ZmNqdVJUajRqRjFnZ1VVVElad2RWZkcyMm1aSXBVVmpMX2RDSWhaLXBYb1hVU0FJTDNkNGhjRGNEWjhOeDRQVE9CVk94V2VKYkRILXQ1ZEdxbkF5RzYzQ0NlOHZiT0JxTFEiLCJ2IjoieGZ2dnZPbFQxZEZrZW5FR2pSY0E1QSJ9\"},\"5\":{\"network_equipment_id\":5,\"network_equipment_status\":\"active\",\"datacenter_name\":\"us02-chi-qts01-dc\",\"network_equipment_driver\":\"os_10\",\"network_equipment_position\":\"leaf\",\"network_equipment_provisioner_type\":\"evpnvxlanl2\",\"network_equipment_identifier_string\":\"sw2-env2\",\"network_equipment_description\":\"OS10 Enterprise. OS Version\",\"network_equipment_management_address\":\"10.0.5.7\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"admin\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"192.168.254.0\",\"network_equipment_quarantine_subnet_end\":\"192.168.254.255\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"192.168.254.2\",\"network_equipment_primary_wan_ipv4_subnet_pool\":null,\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":null,\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":3,\"network_equipment_primary_wan_ipv6_subnet_cidr\":null,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_cached_updated_timestamp\":\"0000-00-00T00:00:00Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":null,\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_tags_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":true,\"network_equipment_interfaces_blacklist_json\":null,\"network_equipment_controller_id\":null,\"network_equipment_is_border_device\":false,\"network_equipment_is_storage_switch\":true,\"network_equipment_network_types_allowed_json\":\"[\\\"wan\\\", \\\"quarantine\\\", \\\"san\\\"]\",\"network_equipment_order_index\":20,\"network_equipment_management_password\":\"Use bsidev.password_decrypt:eyJycWkiOiJici4ybUpyYUNLZXlTYktKWDFvYXFPcnY1SzhMSWFQSWxLQ04xeTRTbFk4blJjdUs0WmV5M3gySU9JejcxTFprTk83MDg3VDRzNElsbnd1MjdMVFNXMm5iYkhaR3gzbGQwWWdjTnZNcTRjMlVCeWZZQW93YnQ3QU5uT09KQVdhRjk4RWZfaXhLNkFSRWFteWpFNloxbFFYSXciLCJ2IjoiT2Nxakhza2M0UEt3ZkRVaUdnekl2USJ9\"}}"
