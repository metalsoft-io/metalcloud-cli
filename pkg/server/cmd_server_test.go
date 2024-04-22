package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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

	interfaces := []metalcloud.SwitchInterfaceSearchResult{{
		ServerID: 100,
	}}

	list := []metalcloud.ServerSearchResult{
		server,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		ServersSearch("").
		Return(&list, nil).
		AnyTimes()

	client.EXPECT().
		SwitchInterfaceSearch("*").
		Return(&interfaces, nil).
		AnyTimes()

	//test json
	format := "json"
	emptyStr := ""
	cmd := command.Command{
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

	//test plaintext
	format = ""
	cmd = command.Command{
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

	bTrue := true

	cmd = command.Command{
		Arguments: map[string]interface{}{
			"filter":         &emptyStr,
			"format":         &format,
			"show_rack_info": &bTrue,
		},
	}

	ret, err = serversListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 100)))

	Expect(csv[1][8]).To(Equal("id-20040424"))
	Expect(csv[1][9]).To(Equal("Rack Name"))
	Expect(csv[1][10]).To(Equal("L-2004"))
	Expect(csv[1][11]).To(Equal("U-2404"))
}

func TestServersListWithCredsCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	var serverList []metalcloud.ServerSearchResult
	err := json.Unmarshal([]byte(_serverListFixture1), &serverList)
	Expect(err).To(BeNil())

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		ServersSearch("").
		Return(&serverList, nil).
		AnyTimes()

	var server metalcloud.Server
	err = json.Unmarshal([]byte(_serverFixture2), &server)

	client.EXPECT().
		ServerGet(16, true).
		Return(&server, nil).
		AnyTimes()

	//test json
	format := ""
	emptyStr := ""
	bTrue := true
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"filter":           &emptyStr,
			"format":           &format,
			"show_credentials": &bTrue,
			"no_color":         &bTrue,
		},
	}

	ret, err := serversListCmd(&cmd, client)
	t.Logf("%s", ret)

	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("root"))
	Expect(ret).To(ContainSubstring("testccc"))

}

func TestServerGetCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	serverType := metalcloud.ServerType{
		ServerTypeID:   100,
		ServerTypeName: "testtype",
	}

	ii := "id-20040424"
	rackName := "Rack Name"
	lowerU := "L-123"
	upperU := "U-123"
	server := metalcloud.Server{
		ServerID:                    10,
		ServerProductName:           "test",
		ServerTypeID:                100,
		ServerInventoryId:           &ii,
		ServerRackName:              &rackName,
		ServerRackPositionLowerUnit: &lowerU,
		ServerRackPositionUpperUnit: &upperU,
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

	cmd := command.Command{
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
	Expect(r["INVENTORY_ID"].(string)).To(Equal(*server.ServerInventoryId))
	Expect(r["RACK_NAME"].(string)).To(Equal(*server.ServerRackName))
	Expect(r["RACK_POSITION_LOWER_UNIT"].(string)).To(Equal(*server.ServerRackPositionLowerUnit))
	Expect(r["RACK_POSITION_UPPER_UNIT"].(string)).To(Equal(*server.ServerRackPositionUpperUnit))

	//test plaintext
	format = ""
	bTrue := true
	cmd = command.Command{
		Arguments: map[string]interface{}{
			"server_id_or_uuid": &id,
			"format":            &format,
			"show_rack_data":    bTrue,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = command.Command{
		Arguments: map[string]interface{}{
			"server_id_or_uuid": &id,
			"format":            &format,
			"show_rack_data":    bTrue,
		},
	}

	ret, err = serverGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 10)))

	Expect(csv[1][3]).To(Equal("id-20040424"))
	Expect(csv[1][4]).To(Equal("Rack Name"))
	Expect(csv[1][5]).To(Equal("L-123"))
	Expect(csv[1][6]).To(Equal("U-123"))
	Expect(csv[1][10]).To(Equal("test"))
}

func TestGetMultipleServerCreateUnmanagedInternalFromYamlFile(t *testing.T) {
	RegisterTestingT(t)
	//ctrl := gomock.NewController(t)

	f1, err := os.CreateTemp(os.TempDir(), "test-*.yaml")
	if err != nil {
		panic(err)
	}

	f1.WriteString(_multiYamlServerFixture1)

	f1.Close()

	records, err := getMultipleServerCreateUnmanagedInternalFromYamlFile(f1.Name())
	Expect(err).To(BeNil())
	Expect(len(records)).To(Equal(3))
	Expect(records[1].ServerCreateUnmanaged.ServerSerialNumber).To(Equal("FMAAB"))
}

/*
func TestServerUpdateCmd(t *testing.T) {

		RegisterTestingT(t)
		ctrl := gomock.NewController(t)

		client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

		var srv metalcloud.Server

		err := json.Unmarshal([]byte(_serverFixture1), &srv)
		if err != nil {
			t.Error(err)
		}

		f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
		if err != nil {
			t.Error(err)
		}

		//create an input json file
		f.WriteString(_serverFixture1)
		f.Close()
		defer syscall.Unlink(f.Name())

		f2, err := os.CreateTemp(os.TempDir(), "testconf-*.yaml")
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
					"server_id_or_uuid":     310,
					"read_config_from_file": f.Name(),
					"format":                "json",
				}),
				Good: true,
			},
			{
				Name: "get-from-yaml-good1",
				Cmd: command.MakeCommand(map[string]interface{}{
					"server_id_or_uuid":     310,
					"read_config_from_file": f2.Name(),
					"format":                "yaml",
				}),
				Good: true,
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				_, err := serverUpdateCmd(&c.Cmd, client)
				if c.Good && err != nil {
					t.Error(err)
				}
			})
		}
	}
*/
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

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"server_id":   1,
				"operation":   "on",
				"autoconfirm": true,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "good2",
			Cmd: command.MakeCommand(map[string]interface{}{
				"server_id": 1,
				"operation": "on",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "missing server_id",
			Cmd: command.MakeCommand(map[string]interface{}{
				"operation":   "on",
				"autoconfirm": true,
			}),
			Good: false,
			Id:   0,
		},
		{
			Name: "missing operation",
			Cmd: command.MakeCommand(map[string]interface{}{
				"server_id":   1,
				"autoconfirm": true,
			}),
			Good: false,
			Id:   0,
		},
	}

	command.TestCreateCommand(serverPowerControlCmd, cases, client, t)
}

func TestServerRegisterCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	serverID := 1

	client.EXPECT().
		ServerCreateAndRegister(gomock.Any()).
		Return(serverID, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "good",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter":    "name",
				"server_vendor": "vendor",
				"mgmt_address":  "127.0.0.1",
				"mgmt_user":     "user",
				"mgmt_pass":     "pass",
				"return_id":     true,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "missing datacenter",
			Cmd: command.MakeCommand(map[string]interface{}{
				"server_vendor": "vendor",
				"mgmt_address":  "127.0.0.1",
				"mgmt_user":     "user",
				"mgmt_pass":     "pass",
				"return_id":     true,
			}),
			Good: false,
			Id:   0,
		},
	}

	command.TestCreateCommand(serverRegisterCmd, cases, client, t)
}

func TestServerImportCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	serverTypeLabel := "serverType"
	srv := ServerCreateUnmanagedInternal{
		ServerTypeLabel: &serverTypeLabel,
		ServerCreateUnmanaged: metalcloud.ServerCreateUnmanaged{
			ServerSerialNumber: "serialnumber",
			ServerInterfaces: []metalcloud.ServerInterfaceCreate{
				{
					ServerInterfaceMACAddress:                 "aa:bb:cc:dd",
					NetworkEquipmentIdentifierString:          "sw1",
					NetworkEquipmentInterfaceIdentifierString: "eth1",
				},
				{
					ServerInterfaceMACAddress:                 "aa:bb:cc:dd",
					NetworkEquipmentIdentifierString:          "sw1",
					NetworkEquipmentInterfaceIdentifierString: "eth1",
				},
			},
		},
	}

	email := "test@test.com"

	//create an input yaml file
	f1Name := createYamlFileFromObject(srv)
	defer syscall.Unlink(f1Name)
	infraLabel := "test1"
	srv2 := ServerCreateUnmanagedInternal{
		InfrastructureLabel: &infraLabel,
		UserEmail:           &email,
		ServerCreateUnmanaged: metalcloud.ServerCreateUnmanaged{
			ServerTypeID:       103,
			ServerSerialNumber: "serialnumber",

			ServerInterfaces: []metalcloud.ServerInterfaceCreate{
				{
					ServerInterfaceMACAddress:                 "aa:bb:cc:dd",
					NetworkEquipmentIdentifierString:          "sw1",
					NetworkEquipmentInterfaceIdentifierString: "eth1",
				},
				{
					ServerInterfaceMACAddress:                 "aa:bb:cc:dd",
					NetworkEquipmentIdentifierString:          "sw1",
					NetworkEquipmentInterfaceIdentifierString: "eth1",
				},
			},
		},
	}

	//create an input yaml file
	f2Name := createYamlFileFromObject(srv2)
	defer syscall.Unlink(f2Name)

	fullServer := metalcloud.Server{
		ServerID: 100,
	}
	client.EXPECT().
		ServerUnmanagedImport(gomock.Any()).
		Return(&fullServer, nil).
		AnyTimes()

	client.EXPECT().
		GetUserEmail().
		Return(email).
		AnyTimes()

	user := metalcloud.User{
		UserEmail: email,
	}
	client.EXPECT().
		UserGetByEmail(gomock.Any()).
		Return(&user, nil).
		AnyTimes()

	serverType := metalcloud.ServerType{
		ServerTypeID: 101,
	}
	client.EXPECT().
		ServerTypeGetByLabel(gomock.Any()).
		Return(&serverType, nil).
		AnyTimes()

	//already exists
	infras := []metalcloud.InfrastructuresSearchResult{
		{
			InfrastructureID:    100,
			InfrastructureLabel: "testinfra",
		},
	}
	client.EXPECT().
		InfrastructureSearch(gomock.Any()).
		Return(&infras, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureOperationCancel(100).
		Return(nil).
		AnyTimes()

	ia := metalcloud.InstanceArray{
		InstanceArrayID: 200,
	}
	client.EXPECT().
		InstanceArrayCreate(100, gomock.Any()).
		Return(&ia, nil).
		AnyTimes()

	instances := map[string]metalcloud.Instance{
		"instance-300": {
			InstanceID: 300,
		},
	}
	client.EXPECT().
		InstanceArrayInstances(200).
		Return(&instances, nil).
		AnyTimes()

	instance := metalcloud.Instance{
		InstanceID: 300,
	}
	client.EXPECT().
		InstanceEdit(300, gomock.Any()).
		Return(&instance, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureDeployWithOptions(100, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	addToInfra := "testinfra"

	cases := []command.CommandTestCase{
		{
			Name: "with server type as label",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f1Name,
				"format":                "yaml",
				"return_id":             true,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "with server type as id",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f2Name,
				"format":                "yaml",
				"return_id":             true,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "with server type as id with add to infra",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f2Name,
				"format":                "yaml",
				"add_to_infra":          &addToInfra,
				"return_id":             true,
			}),
			Good: true,
			Id:   0,
		},
	}

	command.TestCreateCommand(serverImportCmd, cases, client, t)
}

func TestAddServerToInfrastructure(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	infra := metalcloud.Infrastructure{
		InfrastructureID: 99,
	}
	client.EXPECT().
		InfrastructureGetByLabel("99").
		Return(&infra, nil).
		AnyTimes()

	server := metalcloud.Server{
		ServerTypeID:       89,
		ServerID:           59,
		ServerSerialNumber: "ABC",
	}
	client.EXPECT().
		ServerGet(59, false).
		Return(&server, nil).
		AnyTimes()

	ia := metalcloud.InstanceArray{
		InstanceArrayID: 109,
	}
	client.EXPECT().
		InstanceArrayCreate(gomock.Any(), gomock.Any()).
		Return(&ia, nil).
		AnyTimes()

	instances := map[string]metalcloud.Instance{
		"209": {
			InstanceID: 209,
		},
	}

	client.EXPECT().
		InstanceArrayInstances(109).
		Return(&instances, nil).
		AnyTimes()

	i := metalcloud.Instance{
		InstanceID: 209,
	}
	client.EXPECT().
		InstanceEdit(209, gomock.Any()).
		Return(&i, nil).
		AnyTimes()

	c := command.MakeCommand(map[string]interface{}{
		"server_id":                  59,
		"infrastructure_id_or_label": "99",
		"return_id":                  true,
	})

	_, err := serverAddToInfraCmd(&c, client)
	Expect(err).To(BeNil())
}

func TestImportServersBatch(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	infra := metalcloud.Infrastructure{
		InfrastructureID: 99,
	}
	client.EXPECT().
		InfrastructureGetByLabel("99").
		Return(&infra, nil).
		AnyTimes()

	ia := metalcloud.InstanceArray{
		InstanceArrayID: 109,
	}
	client.EXPECT().
		InstanceArrayCreate(gomock.Any(), gomock.Any()).
		Return(&ia, nil).
		MinTimes(1) //make sure this is called (meaning it went down the path of creating instance arrays on the infra)

	instances := map[string]metalcloud.Instance{
		"209": {
			InstanceID: 209,
		},
	}

	client.EXPECT().
		InstanceArrayInstances(109).
		Return(&instances, nil).
		AnyTimes()

	i := metalcloud.Instance{
		InstanceID: 209,
	}
	client.EXPECT().
		InstanceEdit(209, gomock.Any()).
		Return(&i, nil).
		AnyTimes()

	servers := map[string]metalcloud.Server{
		"10": {
			ServerID:           10,
			ServerSerialNumber: "FMAAA",
		},
		"11": {
			ServerID:           11,
			ServerSerialNumber: "FMAAB",
		},
		"12": {
			ServerID:           12,
			ServerSerialNumber: "FMAAC",
		},
	}
	client.EXPECT().
		ServerUnmanagedImportBatch(gomock.Any()).
		Return(&servers, nil).
		AnyTimes()

	s1 := servers["10"]
	client.EXPECT().
		ServerGet(10, false).
		Return(&s1, nil).
		AnyTimes()

	s2 := servers["11"]
	client.EXPECT().
		ServerGet(11, false).
		Return(&s2, nil).
		AnyTimes()

	s3 := servers["12"]
	client.EXPECT().
		ServerGet(12, false).
		Return(&s3, nil).
		AnyTimes()

	f1, err := os.CreateTemp(os.TempDir(), "test-*.yaml")
	if err != nil {
		panic(err)
	}

	f1.WriteString(_multiYamlServerFixture1)

	f1.Close()

	c := command.MakeCommand(map[string]interface{}{
		"read_config_from_file": f1.Name(),
		"add_to_infra":          "99",
		"return_id":             true,
	})

	//should return the server ids
	ret, err := serverImportBatchCmd(&c, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("10"))
	Expect(ret).To(ContainSubstring("11"))
	Expect(ret).To(ContainSubstring("12"))
}

func createYamlFileFromObject(obj interface{}) string {
	//create an input yaml file
	f1, err := os.CreateTemp(os.TempDir(), "test-*.yaml")
	if err != nil {
		panic(err)
	}

	s, err := yaml.Marshal(obj)
	Expect(err).To(BeNil())

	f1.WriteString(string(s))
	f1.Close()
	return f1.Name()
}

const _serverFixture1 = "{\"server_id\":310,\"agent_id\":44,\"datacenter_name\":\"es-madrid\",\"server_uuid\":\"44454C4C-5900-1033-8032-B9C04F434631\",\"server_serial_number\":\"9Y32CF1\",\"server_product_name\":\"PowerEdge 1950\",\"server_vendor\":\"Dell Inc.\",\"server_vendor_sku_id\":\"0\",\"server_ipmi_host\":\"10.255.237.28\",\"server_ipmi_internal_username\":\"ddd\",\"server_ipmi_internal_password_encrypted\":\"BSI\\\\JSONRPC\\\\Server\\\\Security\\\\Authorization\\\\DeveloperAuthorization: Not leaking database encrypted values for extra security.\",\"server_ipmi_version\":\"2\",\"server_ram_gbytes\":8,\"server_processor_count\":2,\"server_processor_core_mhz\":2333,\"server_processor_core_count\":4,\"server_processor_name\":\"Intel(R) Xeon(R) CPU           E5345  @ 2.33GHz\",\"server_processor_cpu_mark\":0,\"server_processor_threads\":1,\"server_type_id\":14,\"server_status\":\"available\",\"server_comments\":\"a\",\"server_details_xml\":null,\"server_network_total_capacity_mbps\":4000,\"server_ipmi_channel\":0,\"server_power_status\":\"off\",\"server_power_status_last_update_timestamp\":\"2020-08-19T08:42:22Z\",\"server_ilo_reset_timestamp\":\"0000-00-00T00:00:00Z\",\"server_boot_last_update_timestamp\":null,\"server_bdk_debug\":false,\"server_dhcp_status\":\"deny_requests\",\"server_bios_info_json\":\"{\\\"server_bios_vendor\\\":\\\"Dell Inc.\\\",\\\"server_bios_version\\\":\\\"2.7.0\\\"}\",\"server_vendor_info_json\":\"{\\\"management\\\":\\\"iDRAC\\\",\\\"version\\\":\\\"er] rpcRoundRobinConnectedAgentsOfType() failed with error: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/task_queues.js:97:5) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8) Exception: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n FetchError: request to https:\\\\/\\\\/10.255.237.28\\\\/cgi-bin\\\\/webcgi\\\\/about failed, reason: write EPROTO 38858976:error:1425F102:SSL routines:ssl_choose_client_version:unsupported protocol:..\\\\/deps\\\\/openssl\\\\/openssl\\\\/ssl\\\\/statem\\\\/statem_lib.c:1922:\\\\n\\\\n    at ClientRequest.<anonymous> (\\\\/var\\\\/datacenter-agents-binary-compiled-temp\\\\/Power\\\\/Power.portable.js:8:469877)\\\\n    at ClientRequest.emit (events.js:209:13)\\\\n    at TLSSocket.socketErrorListener (_http_client.js:406:9)\\\\n    at TLSSocket.emit (events.js:209:13)\\\\n    at errorOrDestroy (internal\\\\/streams\\\\/destroy.js:107:12)\\\\n    at onwriteError (_stream_writable.js:449:5)\\\\n    at onwrite (_stream_writable.js:470:5)\\\\n    at internal\\\\/streams\\\\/destroy.js:49:7\\\\n    at TLSSocket.Socket._destroy (net.js:595:3)\\\\n    at TLSSocket.destroy (internal\\\\/streams\\\\/destroy.js:37:8)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/task_queues.js:97:5)\\\\n    at \\\\/var\\\\/vhosts\\\\/bsiintegration.bigstepcloud.com\\\\/BSIWebSocketServer\\\\/node_modules\\\\/jsonrpc-bidirectional\\\\/src\\\\/Client.js:331:37\\\\n    at runMicrotasks (<anonymous>)\\\\n    at processTicksAndRejections (internal\\\\/process\\\\/tas\\\"}\",\"server_class\":\"bigdata\",\"server_created_timestamp\":\"2019-07-02T07:57:19Z\",\"subnet_oob_id\":2,\"subnet_oob_index\":28,\"server_boot_type\":\"classic\",\"server_disk_wipe\":true,\"server_disk_count\":0,\"server_disk_size_mbytes\":0,\"server_disk_type\":\"none\",\"server_requires_manual_cleaning\":false,\"chassis_rack_id\":null,\"server_custom_json\":\"{\\\"previous_ipmi_username\\\":\\\"a\\\",\\\"previous_ipmi_password_encrypted\\\":\\\"rq|aes-cbc|urfNNCbe2ouIRX3reLrILyM7tBD5I1aMPycR3YkCeFo1DGEGnNI3n6u7z63sBWpW\\\"}\",\"server_instance_custom_json\":null,\"server_last_cleanup_start\":\"2020-08-12T14:26:47Z\",\"server_allocation_timestamp\":null,\"server_dhcp_packet_sniffing_is_enabled\":true,\"snmp_community_password_dcencrypted\":null,\"server_mgmt_snmp_community_password_dcencrypted\":\"BSI\\\\JSONRPC\\\\Server\\\\Security\\\\Authorization\\\\DeveloperAuthorization: Not leaking database encrypted values for extra security.\",\"server_mgmt_snmp_port\":161,\"server_mgmt_snmp_version\":2,\"server_dhcp_relay_security_is_enabled\":true,\"server_keys_json\":\"{\\\"keys\\\": {\\\"r1\\\": {\\\"created\\\": \\\"2019-07-02T07:59:17Z\\\", \\\"salt_encrypted\\\": \\\"rq|aes-cbc|9721g561woNQzA0a3yWTcHcEYxJo7vXNc1SHmEUCxYdeOqsiVbT+X+leOHHP+XsR1gfOgs8lMhdXLOw0UUBP8g==\\\", \\\"aes_key_encrypted\\\": \\\"rq|aes-cbc|/V4Y7FMu9Uo4PyktBKl+jsAKpogNh+UC2F03jxMtJI2ieacgx/Ogso0Z9d3XlL99zh1pxAPVF24gzAogNIla0L0xBgUgLicJt41ajRYvdIo=\\\"}}, \\\"active_index\\\": \\\"r1\\\", \\\"keys_partition\\\": \\\"server_id_310\\\"}\",\"server_info_json\":null,\"server_ipmi_credentials_need_update\":false,\"server_gpu_count\":0,\"server_gpu_vendor\":\"\",\"server_gpu_model\":\"\",\"server_bmc_mac_address\":null,\"server_metrics_metadata_json\":null,\"server_interfaces\":[{\"server_interface_mac_address\":\"00:1d:09:64:f0:2b\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:1d:09:64:f0:2d\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:15:17:c0:4c:e6\",\"type\":\"ServerInterface\"},{\"server_interface_mac_address\":\"00:15:17:c0:4c:e7\",\"type\":\"ServerInterface\"}],\"server_disks\":[],\"server_tags\":[],\"type\":\"Server\"}"
const _serverListFixture1 = "[\n                {\n                    \"server_id\": 16,\n                    \"server_type_name\": null,\n                    \"server_type_boot_type\": null,\n                    \"server_product_name\": null,\n                    \"datacenter_name\": \"us02-chi-qts01-dc\",\n                    \"server_status\": \"registering\",\n                    \"server_class\": \"bigdata\",\n                    \"server_created_timestamp\": \"2022-05-23T13:22:11Z\",\n                    \"server_vendor\": \"Dell Inc.\",\n                    \"server_serial_number\": null,\n                    \"server_uuid\": \"4c4c4544-0051-3810-8057-b7c04f533532\",\n                    \"server_vendor_sku_id\": null,\n                    \"server_boot_type\": \"classic\",\n                    \"server_allocation_timestamp\": null,\n                    \"instance_label\": [\n                        null\n                    ],\n                    \"instance_id\": [\n                        null\n                    ],\n                    \"instance_array_id\": [\n                        null\n                    ],\n                    \"infrastructure_id\": [\n                        null\n                    ],\n                    \"server_inventory_id\": null,\n                    \"server_rack_name\": null,\n                    \"server_rack_position_lower_unit\": null,\n                    \"server_rack_position_upper_unit\": null,\n                    \"server_ipmi_host\": \"172.18.44.42\",\n                    \"server_ipmi_internal_username\": \"root\",\n                    \"server_processor_name\": null,\n                    \"server_processor_count\": 0,\n                    \"server_processor_core_count\": 0,\n                    \"server_processor_core_mhz\": 0,\n                    \"server_processor_threads\": null,\n                    \"server_processor_cpu_mark\": null,\n                    \"server_disk_type\": \"none\",\n                    \"server_disk_count\": 0,\n                    \"server_disk_size_mbytes\": 0,\n                    \"server_ram_gbytes\": 0,\n                    \"server_network_total_capacity_mbps\": 0,\n                    \"server_dhcp_status\": \"quarantine\",\n                    \"server_dhcp_packet_sniffing_is_enabled\": true,\n                    \"server_dhcp_relay_security_is_enabled\": true,\n                    \"server_disk_wipe\": false,\n                    \"server_power_status\": \"off\",\n                    \"server_power_status_last_update_timestamp\": \"2022-05-23T13:24:41Z\",\n                    \"user_id\": [\n                        [\n                            null\n                        ]\n                    ],\n                    \"user_id_owner\": [\n                        null\n                    ],\n                    \"user_email\": [\n                        [\n                            null\n                        ]\n                    ],\n                    \"infrastructure_user_id\": [\n                        [\n                            null\n                        ]\n                    ]\n                }\n            ]"
const _serverFixture2 = "{\n        \"server_id\": 16,\n        \"agent_id\": null,\n        \"datacenter_name\": \"us02-chi-qts01-dc\",\n        \"server_uuid\": \"4c4c4544-0051-3810-8057-b7c04f533532\",\n        \"server_serial_number\": null,\n        \"server_product_name\": null,\n        \"server_vendor\": \"Dell Inc.\",\n        \"server_vendor_sku_id\": null,\n        \"server_ipmi_host\": \"172.18.44.42\",\n        \"server_ipmi_internal_username\": \"root\",\n        \"server_ipmi_internal_password\": \"testcccc\",\n        \"server_ipmi_version\": \"2\",\n        \"server_ram_gbytes\": 0,\n        \"server_processor_count\": 0,\n        \"server_processor_core_mhz\": 0,\n        \"server_processor_core_count\": 0,\n        \"server_processor_name\": null,\n        \"server_processor_cpu_mark\": null,\n        \"server_processor_threads\": null,\n        \"server_type_id\": null,\n        \"server_status\": \"registering\",\n        \"server_comments\": null,\n        \"server_details_xml\": null,\n        \"server_network_total_capacity_mbps\": 0,\n        \"server_ipmi_channel\": 1,\n        \"server_power_status\": \"off\",\n        \"server_power_status_last_update_timestamp\": \"2022-05-23T13:24:41Z\",\n        \"server_ilo_reset_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_boot_last_update_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_bdk_debug\": false,\n        \"server_dhcp_status\": \"quarantine\",\n        \"server_bios_info_json\": null,\n        \"server_vendor_info_json\": null,\n        \"server_class\": \"bigdata\",\n        \"server_created_timestamp\": \"2022-05-23T13:22:11Z\",\n        \"subnet_oob_id\": 5,\n        \"subnet_oob_index\": 42,\n        \"server_boot_type\": \"classic\",\n        \"server_disk_wipe\": false,\n        \"server_disk_count\": 0,\n        \"server_disk_size_mbytes\": 0,\n        \"server_disk_type\": \"none\",\n        \"server_requires_manual_cleaning\": false,\n        \"chassis_rack_id\": null,\n        \"server_custom_json\": null,\n        \"server_instance_custom_json\": null,\n        \"server_last_cleanup_start\": null,\n        \"server_allocation_timestamp\": null,\n        \"server_dhcp_packet_sniffing_is_enabled\": true,\n        \"snmp_community_password_dcencrypted\": null,\n        \"server_mgmt_snmp_community_password_dcencrypted\": null,\n        \"server_mgmt_snmp_port\": 161,\n        \"server_mgmt_snmp_version\": 2,\n        \"server_dhcp_relay_security_is_enabled\": true,\n        \"server_keys_json\": null,\n        \"server_info_json\": null,\n        \"server_ipmi_credentials_need_update\": false,\n        \"server_gpu_count\": 0,\n        \"server_gpu_vendor\": null,\n        \"server_gpu_model\": null,\n        \"server_bmc_mac_address\": null,\n        \"server_metrics_metadata_json\": null,\n        \"server_secure_boot_is_enabled\": false,\n        \"server_chipset_name\": null,\n        \"server_requires_reregister\": false,\n        \"server_rack_name\": null,\n        \"server_rack_position_upper_unit\": null,\n        \"server_rack_position_lower_unit\": null,\n        \"server_inventory_id\": null,\n        \"server_registered_timestamp\": \"0000-00-00T00:00:00Z\",\n        \"server_interfaces\": [],\n        \"server_disks\": [],\n        \"server_tags\": [],\n        \"type\": \"Server\"\n    }"
const _multiYamlServerFixture1 = `
label: test
serialNumber: FMAAA
---
label: test
serialNumber: FMAAB
---
label: test
serialNumber: FMAAC
`
