package instance

import (
	"encoding/json"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"

	"github.com/metalsoft-io/metalcloud-cli/internal/command"
)

func TestInstanceGetCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	instance := metalcloud.Instance{
		InstanceID:      110,
		InstanceArrayID: 10,
		InstanceCredentials: metalcloud.InstanceCredentials{
			SSH: &metalcloud.SSH{
				Username:        "testu",
				InitialPassword: "testp",
				Port:            22,
			},
		},
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArraySubdomain:    "tst",
		InstanceArrayID:           10,
		InstanceArrayDeployStatus: "not_started",
		InstanceArrayDeployType:   "edit",
	}

	ia := metalcloud.InstanceArray{
		InstanceArraySubdomain:     "tst",
		InstanceArrayID:            10,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "ordered",
	}
	infra := metalcloud.Infrastructure{
		InfrastructureID:    10,
		InfrastructureLabel: "tsassd",
	}

	client.EXPECT().
		InstanceGet(gomock.Any()).
		Return(&instance, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(gomock.Any()).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(gomock.Any()).
		Return(&infra, nil).
		AnyTimes()

	cmd := command.MakeCommand(map[string]interface{}{"instance_id": 110, "show_credentials": true})

	ret, err := instanceGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("ID"))

}

func TestInstanceServerReplaceCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	instance := metalcloud.Instance{
		InstanceID:      110,
		InstanceArrayID: 10,
		InstanceCredentials: metalcloud.InstanceCredentials{
			SSH: &metalcloud.SSH{
				Username:        "testu",
				InitialPassword: "testp",
				Port:            22,
			},
		},
	}

	ia := metalcloud.InstanceArray{
		InstanceArraySubdomain: "tst",
		InstanceArrayID:        10,
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10,
		InfrastructureLabel: "tsassd",
	}

	client.EXPECT().
		InstanceGet(gomock.Any()).
		Return(&instance, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(gomock.Any()).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(gomock.Any()).
		Return(&infra, nil).
		AnyTimes()

	var server metalcloud.Server
	json.Unmarshal([]byte(_serverFixture2), &server)

	client.EXPECT().
		ServerGet(gomock.Any(), false).
		Return(&server, nil).
		AnyTimes()

	client.EXPECT().
		InstanceServerReplace(110, 100).
		Return(500, nil).
		MinTimes(2)

	cmd := command.MakeCommand(map[string]interface{}{
		"instance_id": 110,
		"server_id":   100,
		"autoconfirm": true})

	_, err := instanceServerReplaceCmd(&cmd, client)
	Expect(err).To(BeNil())

	cmd = command.MakeCommand(map[string]interface{}{
		"instance_id":   110,
		"server_id":     100,
		"return_afc_id": true,
		"autoconfirm":   true})

	ret, err := instanceServerReplaceCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Equal("500"))
}

const _serverFixture2 = "{\n        \"server_id\": 16,\n        \"agent_id\": null,\n        \"datacenter_name\": \"us02-chi-qts01-dc\",\n        \"server_uuid\": \"4c4c4544-0051-3810-8057-b7c04f533532\",\n        \"server_serial_number\": null,\n        \"server_product_name\": null,\n        \"server_vendor\": \"Dell Inc.\",\n        \"server_vendor_sku_id\": null,\n        \"server_ipmi_host\": \"172.18.44.42\",\n        \"server_ipmi_internal_username\": \"root\",\n        \"server_ipmi_internal_password\": \"testcccc\",\n        \"server_ipmi_version\": \"2\",\n        \"server_ram_gbytes\": 0,\n        \"server_processor_count\": 0,\n        \"server_processor_core_mhz\": 0,\n        \"server_processor_core_count\": 0,\n        \"server_processor_name\": null,\n        \"server_processor_cpu_mark\": null,\n        \"server_processor_threads\": null,\n        \"server_type_id\": null,\n        \"server_status\": \"registering\",\n        \"server_comments\": null,\n        \"server_details_xml\": null,\n        \"server_network_total_capacity_mbps\": 0,\n        \"server_ipmi_channel\": 1,\n        \"server_power_status\": \"off\",\n        \"server_power_status_last_update_timestamp\": \"2022-05-23T13:24:41Z\",\n        \"server_ilo_reset_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_boot_last_update_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_bdk_debug\": false,\n        \"server_dhcp_status\": \"quarantine\",\n        \"server_bios_info_json\": null,\n        \"server_vendor_info_json\": null,\n        \"server_class\": \"bigdata\",\n        \"server_created_timestamp\": \"2022-05-23T13:22:11Z\",\n        \"subnet_oob_id\": 5,\n        \"subnet_oob_index\": 42,\n        \"server_boot_type\": \"classic\",\n        \"server_disk_wipe\": false,\n        \"server_disk_count\": 0,\n        \"server_disk_size_mbytes\": 0,\n        \"server_disk_type\": \"none\",\n        \"server_requires_manual_cleaning\": false,\n        \"chassis_rack_id\": null,\n        \"server_custom_json\": null,\n        \"server_instance_custom_json\": null,\n        \"server_last_cleanup_start\": null,\n        \"server_allocation_timestamp\": null,\n        \"server_dhcp_packet_sniffing_is_enabled\": true,\n        \"snmp_community_password_dcencrypted\": null,\n        \"server_mgmt_snmp_community_password_dcencrypted\": null,\n        \"server_mgmt_snmp_port\": 161,\n        \"server_mgmt_snmp_version\": 2,\n        \"server_dhcp_relay_security_is_enabled\": true,\n        \"server_keys_json\": null,\n        \"server_info_json\": null,\n        \"server_ipmi_credentials_need_update\": false,\n        \"server_gpu_count\": 0,\n        \"server_gpu_vendor\": null,\n        \"server_gpu_model\": null,\n        \"server_bmc_mac_address\": null,\n        \"server_metrics_metadata_json\": null,\n        \"server_secure_boot_is_enabled\": false,\n        \"server_chipset_name\": null,\n        \"server_requires_reregister\": false,\n        \"server_rack_name\": null,\n        \"server_rack_position_upper_unit\": null,\n        \"server_rack_position_lower_unit\": null,\n        \"server_inventory_id\": null,\n        \"server_registered_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_interfaces\": [],\n        \"server_disks\": [],\n        \"server_tags\": [],\n        \"type\": \"Server\"\n    }"
