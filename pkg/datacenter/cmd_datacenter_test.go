package datacenter

import (
	"encoding/json"
	"os"
	"syscall"
	"testing"
	"path/filepath"
	"runtime"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

func TestDatacenterListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	dcList := map[string]metalcloud.Datacenter{
		"test_hidden": _dcFixture1,
		"test2":       _dcFixture2,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		Datacenters(true).
		Return(&dcList, nil).
		AnyTimes()

	client.EXPECT().
		UserGet(1).
		Return(&_userFixture1, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"LABEL": _dcFixture2.DatacenterName,
		"NAME":  _dcFixture2.DatacenterDisplayName,
	}

	command.TestListCommand(datacenterListCmd, nil, client, expectedFirstRow, t)
}

func TestDatacenterCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var dcConf metalcloud.DatacenterConfig
	err := json.Unmarshal([]byte(_dcConfigFixture1), &dcConf)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		DatacenterCreate(_dcFixture1, dcConf).
		Return(&_dcFixture1, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(_userFixture1.UserEmail).
		Return(&_userFixture1, nil).
		AnyTimes()

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	f.WriteString(_dcConfigFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	cases := []command.CommandTestCase{
		{
			Name: "dc-create-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
				"read_config_from_file":   f.Name(),
				"create_hidden":           true,
				"user_id":                 _userFixture1.UserEmail,
				"format":                  "json",
				"tags":                    "t1,t2",
				"datacenter_name_parent":  "test",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "missing label",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
		{
			Name: "missing title",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name": _dcFixture1.DatacenterName,
			}),
			Good: false,
		},
		{
			Name: "missing read_config_from_file",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
			}),
			Good: false,
		},
		{
			Name: "both read_config_from_file and pipe",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
				"read_config_from_file":   f.Name(),
				"read_config_from_pipe":   true,
			}),
			Good: false,
		},
	}

	command.TestCreateCommand(datacenterCreateCmd, cases, client, t)

}

func TestDatacenterYamlUnmarshal(t *testing.T) {
	RegisterTestingT(t)

	var dcConf metalcloud.DatacenterConfig

	_, b, _, _ := runtime.Caller(0)
    basePath := filepath.Join(filepath.Dir(b), "..", "..")

	content, err := configuration.ReadInputFromFile(filepath.Join(basePath, "examples", "datacenter.yaml"))
	Expect(err).To(BeNil())
	Expect(content).NotTo(BeNil())

	err = yaml.Unmarshal(content, &dcConf)
	Expect(err).To(BeNil())

	Expect(dcConf.ServerRegisterUsingGeneratedIPMICredentialsEnabled).To(BeTrue())
}

func TestDatacenterUpdate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var dcConf metalcloud.DatacenterConfig
	err := json.Unmarshal([]byte(_dcConfigFixture1), &dcConf)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		DatacenterConfigUpdate("test", dcConf).
		Return(nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(_userFixture1.UserEmail).
		Return(&_userFixture1, nil).
		AnyTimes()

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	f.WriteString(_dcConfigFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	cases := []command.CommandTestCase{
		{
			Name: "dc-create-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":       _dcFixture1.DatacenterName,
				"read_config_from_file": f.Name(),
				"format":                "json",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "missing label",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
		{
			Name: "missing read_config_from_file",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
			}),
			Good: false,
		},
		{
			Name: "both read_config_from_file and pipe",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
				"read_config_from_file":   f.Name(),
				"read_config_from_pipe":   true,
			}),
			Good: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

			_, err := datacenterUpdateCmd(&c.Cmd, client)
			if c.Good {

				if err != nil {
					t.Errorf("error thrown: %v", err)
				}

			} else {
				if err == nil {
					t.Errorf("Should have thrown error")
				}
			}
		})
	}
}

func TestDatacenterGet(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		DatacenterGet(_dcFixture1.DatacenterName).
		Return(&_dcFixture1, nil).
		AnyTimes()

	var dcConf metalcloud.DatacenterConfig
	err := json.Unmarshal([]byte(_dcConfigFixture1), &dcConf)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		DatacenterConfigGet(_dcFixture1.DatacenterName).
		Return(&dcConf, nil).
		AnyTimes()

	client.EXPECT().
		UserGet(1).
		Return(&_userFixture1, nil).
		AnyTimes()

	client.EXPECT().
		DatacenterAgentsConfigJSONDownloadURL(_dcFixture1.DatacenterName, true).
		Return("https:/asasd/asdasd", nil).
		AnyTimes()

	//should throw error for missing label
	cmd := command.MakeCommand(map[string]interface{}{})

	ret, err := datacenterGetCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name": _dcFixture1.DatacenterName,
	})

	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring(_dcFixture1.DatacenterName))

	//verify config url present in output
	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name":        _dcFixture1.DatacenterName,
		"show_secret_config_url": true,
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("https:/asasd/asdasd"))

	//verify return url
	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name":   _dcFixture1.DatacenterName,
		"return_config_url": true,
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Equal("https:/asasd/asdasd"))

	//verify config
	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name":        _dcFixture1.DatacenterName,
		"show_datacenter_config": true,
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("Datacenter name"))

	//verify json format

	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name":        _dcFixture1.DatacenterName,
		"show_datacenter_config": true,
		"show_secret_config_url": true,
		"format":                 "json",
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["LABEL"].(string)).To(Equal(_dcFixture1.DatacenterName))
	Expect(r["TITLE"].(string)).To(Equal(_dcFixture2.DatacenterDisplayName))

}

var _dcFixture1 metalcloud.Datacenter = metalcloud.Datacenter{
	DatacenterName:          "test",
	DatacenterDisplayName:   "datacenterDisplayName",
	UserID:                  1,
	DatacenterIsMaster:      false,
	DatacenterIsMaintenance: false,
	DatacenterHidden:        true,
	DatacenterTags:          []string{"t1", "t2"},
	DatacenterNameParent:    "test",
}

var _dcFixture2 metalcloud.Datacenter = metalcloud.Datacenter{
	DatacenterName:          "test",
	DatacenterDisplayName:   "datacenterDisplayName",
	UserID:                  1,
	DatacenterIsMaster:      false,
	DatacenterIsMaintenance: false,
	DatacenterHidden:        false,
	DatacenterTags:          []string{"t1", "t2"},
	DatacenterNameParent:    "test",
}

var _userFixture1 = metalcloud.User{
	UserID:    1,
	UserEmail: "test@test.com",
}

const _dcConfigFixture1 = "{\"SANRoutedSubnet\":\"100.96.0.0/24\",\"BSIVRRPListenIPv4\":\"172.31.240.126\",\"BSIMachineListenIPv4List\":[\"172.31.240.124\",\"172.31.240.125\"],\"BSIMachinesSubnetIPv4CIDR\":\"172.31.240.96/27\",\"BSIExternallyVisibleIPv4\":\"10.255.231.54\",\"repoURLRoot\":\"https://repointegration.bigstepcloud.com\",\"repoURLRootQuarantineNetwork\":\"http://10.255.239.35\",\"DNSServers\":[\"10.255.231.44\",\"10.255.231.45\"],\"NTPServers\":[\"10.255.231.28\",\"10.255.231.29\"],\"KMS\":\"10.255.235.41:1688\",\"TFTPServerWANVRRPListenIPv4\":\"172.31.240.126\",\"dataLakeEnabled\":true,\"monitoringGraphitePlainTextSocketHost\":\"172.31.240.148:2003\",\"monitoringGraphiteRenderURLHost\":\"172.31.240.157:80\",\"latitude\":0,\"longitude\":0,\"address\":\"\",\"VLANProvisioner\":{\"LANVLANRange\":\"200-299\",\"WANVLANRange\":\"100-199\",\"quarantineVLANID\":5},\"VPLSProvisioner\":{\"ACLSAN\":3999,\"ACLWAN\":3399,\"SANACLRange\":\"3700-3998\",\"ToRLANVLANRange\":\"400-699\",\"ToRSANVLANRange\":\"700-999\",\"ToRWANVLANRange\":\"100-399\",\"quarantineVLANID\":5,\"NorthWANVLANRange\":\"1001-2000\"},\"childDatacentersConfigDefault\":[]}"
