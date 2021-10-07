package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

func TestServersListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	server := metalcloud.ServerSearchResult{
		ServerID:                    100,
		ServerProductName:           "test",
		ServerInventoryId:           "id-20040424",
		ServerRackName:              "Rack Name",
		ServerRackPositionLowerUnit: "L-2004",
		ServerRackPositionUpperUnit: "U-2404",
	}

	list := []metalcloud.ServerSearchResult{
		server,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		ServersSearch("").
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	emptyStr := ""
	cmd := Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err := serversListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(100))
	Expect(r["PRODUCT_NAME"].(string)).To(Equal(server.ServerProductName))
	Expect(r["INVENTORY_ID"].(string)).To(Equal(server.ServerInventoryId))
	Expect(r["RACK_NAME"].(string)).To(Equal(server.ServerRackName))
	Expect(r["RACK_POSITION_LOWER_UNIT"].(string)).To(Equal(server.ServerRackPositionLowerUnit))
	Expect(r["RACK_POSITION_UPPER_UNIT"].(string)).To(Equal(server.ServerRackPositionUpperUnit))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err = serversListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"filter": &emptyStr,
			"format": &format,
		},
	}

	ret, err = serversListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 100)))
	Expect(csv[1][5]).To(Equal("test"))
	Expect(csv[1][11]).To(Equal("id-20040424"))
	Expect(csv[1][12]).To(Equal("Rack Name"))
	Expect(csv[1][13]).To(Equal("L-2004"))
	Expect(csv[1][14]).To(Equal("U-2404"))
}

func TestServerGetCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	serverType := metalcloud.ServerType{
		ServerTypeID:   100,
		ServerTypeName: "testtype",
	}

	server := metalcloud.Server{
		ServerID:                    10,
		ServerProductName:           "test",
		ServerTypeID:                100,
		ServerInventoryId:           "id-20040424",
		ServerRackName:              "Rack Name",
		ServerRackPositionLowerUnit: "L-2004",
		ServerRackPositionUpperUnit: "U-2404",
	}

	client.EXPECT().
		ServerGet(10, false).
		Return(&server, nil).
		AnyTimes()

	client.EXPECT().
		ServerTypeGet(100).
		Return(&serverType, nil).
		AnyTimes()

	//test json
	id := 10
	format := "json"

	cmd := Command{
		Arguments: map[string]interface{}{
			"server_id_or_uuid": &id,
			"format":            &format,
		},
	}

	ret, err := serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(10))
	Expect(r["PRODUCT_NAME"].(string)).To(Equal(server.ServerProductName))
	Expect(r["INVENTORY_ID"].(string)).To(Equal(server.ServerInventoryId))
	Expect(r["RACK_NAME"].(string)).To(Equal(server.ServerRackName))
	Expect(r["RACK_POSITION_LOWER_UNIT"].(string)).To(Equal(server.ServerRackPositionLowerUnit))
	Expect(r["RACK_POSITION_UPPER_UNIT"].(string)).To(Equal(server.ServerRackPositionUpperUnit))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"server_id_or_uuid": &id,
			"format":            &format,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"server_id_or_uuid": &id,
			"format":            &format,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 10)))
	Expect(csv[1][2]).To(Equal("id-20040424"))
	Expect(csv[1][3]).To(Equal("Rack Name"))
	Expect(csv[1][4]).To(Equal("L-2004"))
	Expect(csv[1][5]).To(Equal("U-2404"))
	Expect(csv[1][9]).To(Equal("test"))
}

func TestServerEditCmd(t *testing.T) {

	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var srv metalcloud.Server

	err := json.Unmarshal([]byte(_serverFixture1), &srv)
	if err != nil {
		t.Error(err)
	}

	f, err := ioutil.TempFile("/tmp", "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_serverFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	f2, err := ioutil.TempFile("/tmp", "testconf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	//create an input yaml file
	s, err := yaml.Marshal(srv)
	Expect(err).To(BeNil())

	f2.WriteString(string(s))
	f2.Close()
	defer syscall.Unlink(f2.Name())

	client.EXPECT().
		ServerGet(310, false).
		Return(&srv, nil).
		AnyTimes()

	client.EXPECT().
		ServerEditComplete(310, gomock.Any()).
		Return(&srv, nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "missing-id",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": f.Name(),
				"format":                "json",
			}),
			good: false,
		},
		{
			name: "get-from-json-good1",
			cmd: MakeCommand(map[string]interface{}{
				"server_id_or_uuid":     310,
				"read_config_from_file": f.Name(),
				"format":                "json",
			}),
			good: true,
		},
		{
			name: "get-from-yaml-good1",
			cmd: MakeCommand(map[string]interface{}{
				"server_id_or_uuid":     310,
				"read_config_from_file": f2.Name(),
				"format":                "yaml",
			}),
			good: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := serverEditCmd(&c.cmd, client)
			if c.good && err != nil {
				t.Error(err)
			}
		})
	}

}

func TestServerPowerControlCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	server := metalcloud.Server{
		ServerID: 1,
	}

	client.EXPECT().
		ServerGet(gomock.Any(), false).
		Return(&server, nil).
		AnyTimes()

	client.EXPECT().
		ServerPowerSet(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"server_id":   1,
				"operation":   "on",
				"autoconfirm": true,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"server_id": 1,
				"operation": "on",
			}),
			good: true,
			id:   0,
		},
		{
			name: "missing server_id",
			cmd: MakeCommand(map[string]interface{}{
				"operation":   "on",
				"autoconfirm": true,
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing operation",
			cmd: MakeCommand(map[string]interface{}{
				"server_id":   1,
				"autoconfirm": true,
			}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(serverPowerControlCmd, cases, client, t)
}

const _serverFixture1 = "{\"server_id\":310,\"agent_id\":44,\"datacenter_name\":\"es-madrid\",\"server_uuid\":\"44454C4C-5900-1033-8032-B9C04F434631\",\"server_serial_number\":\"9Y32CF1\",\"server_product_name\":\"PowerEdge 1950\",\"server_vendor\":\"Dell Inc.\",\"server_vendor_sku_id\":\"0\",\"server_ipmi_host\":\"10.255.237.28\",\"server_ipmi_internal_username\":\"ddd\",\"server_ipmi_internal_password_encrypted\":\"BSI\\\\JSONRPC\\\\Server\\\\Security\\\\Authorization\\\\DeveloperAuthorization: Not leaking database encrypted values for extra security.\",\"server_ipmi_version\":\"2\",\"server_ram_gbytes\":8,\"server_processor_count\":2,\"server_processor_core_mhz\":2333,\"server_processor_core_count\":4,\"server_processor_name\":\"Intel(R) Xeon(R) CPU           E5345  @ 2.33GHz\",\"server_processor_cpu_mark\":0,\"server_processor_threads\":1,\"server_type_id\":14,\"server_status\":\"available\",\"server_comments\":\"a\",\"server_details_xml\":null,\"server_network_total_capacity_mbps\":4000,\"server_ipmi_channel\":0,\"server_power_status\":\"off\",\"server_power_status_last_update_timestamp\":\"2020-08-19T08:42:22Z\",\"server_ilo_reset_timestamp\":\"0000-00-00T00:00:00Z\",\"server_boot_last_update_timestamp\":null,\"server_bdk_debug\":false,\"server_dhcp_status\":\"deny_requests\",\"server_bios_info_json\":\"{\\\"server_bios_vendor\\\":\\\"Dell Inc.\\\",\\\"server_bios_version\\\":\\\"2.7.0\\\"}\",\"server_vendor_info_json\":\"{\\\"management\\\":\\\"iDRAC\\\",\\\"version\\\":\\\"er] rpcRoundRobinConnectedAgentsOfType() failed with error: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/task_queues.js:97:5) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/task_queues.js:97:5)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/tas\\\"}\",\"server_class\":\"bigdata\",\"server_created_timestamp\":\"2019-07-02T07:57:19Z\",\"subnet_oob_id\":2,\"subnet_oob_index\":28,\"server_boot_type\":\"classic\",\"server_disk_wipe\":true,\"server_disk_count\":0,\"server_disk_size_mbytes\":0,\"server_disk_type\":\"none\",\"server_requires_manual_cleaning\":false,\"chassis_rack_id\":null,\"server_custom_json\":\"{\\\"previous_ipmi_username\\\":\\\"a\\\",\\\"previous_ipmi_password_encrypted\\\":\\\"rq|aes-cbc|urfNNCbe2ouIRX3reLrILyM7tBD5I1aMPycR3YkCeFo1DGEGnNI3n6u7z63sBWpW\\\"}\",\"server_instance_custom_json\":null,\"server_last_cleanup_start\":\"2020-08-12T14:26:47Z\",\"server_allocation_timestamp\":null,\"server_dhcp_packet_sniffing_is_enabled\":true,\"snmp_community_password_dcencrypted\":null,\"server_mgmt_snmp_community_password_dcencrypted\":\"BSI\\\\JSONRPC\\\\Server\\\\Security\\\\Authorization\\\\DeveloperAuthorization: Not leaking database encrypted values for extra security.\",\"server_mgmt_snmp_port\":161,\"server_mgmt_snmp_version\":2,\"server_dhcp_relay_security_is_enabled\":true,\"server_keys_json\":\"{\\\"keys\\\": {\\\"r1\\\": {\\\"created\\\": \\\"2019-07-02T07:59:17Z\\\", \\\"salt_encrypted\\\": \\\"rq|aes-cbc|9721g561woNQzA0a3yWTcHcEYxJo7vXNc1SHmEUCxYdeOqsiVbT+X+leOHHP+XsR1gfOgs8lMhdXLOw0UUBP8g==\\\", \\\"aes_key_encrypted\\\": \\\"rq|aes-cbc|/V4Y7FMu9Uo4PyktBKl+jsAKpogNh+UC2F03jxMtJI2ieacgx/Ogso0Z9d3XlL99zh1pxAPVF24gzAogNIla0L0xBgUgLicJt41ajRYvdIo=\\\"}}, \\\"active_index\\\": \\\"r1\\\", \\\"keys_partition\\\": \\\"server_id_310\\\"}\",\"server_info_json\":null,\"server_ipmi_credentials_need_update\":false,\"server_gpu_count\":0,\"server_gpu_vendor\":\"\",\"server_gpu_model\":\"\",\"server_bmc_mac_address\":null,\"server_metrics_metadata_json\":null,\"server_interfaces\":[{\"server_interface_mac_address\":\"00:1d:09:64:f0:2b\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:1d:09:64:f0:2d\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:15:17:c0:4c:e6\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:15:17:c0:4c:e7\",\"type\":\"ServerInterface\"}],\"server_disks\":[],\"server_tags\":[],\"type\":\"Server\"}"
