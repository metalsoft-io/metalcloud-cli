package datacenter

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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

	var dcConf metalcloud.DatacenterWithConfig
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")
	fileName := filepath.Join(basePath, "examples", "datacenter.yaml")
	content, err := configuration.ReadInputFromFile(filepath.Join(basePath, "examples", "datacenter.yaml"))
	Expect(err).To(BeNil())
	Expect(content).NotTo(BeNil())

	err = yaml.Unmarshal(content, &dcConf)
	Expect(err).To(BeNil())

	client.EXPECT().
		DatacenterCreateFromDatacenterWithConfig(dcConf).
		Return(&dcConf, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(_userFixture1.UserEmail).
		Return(&_userFixture1, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "dc-create-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
				"datacenter_display_name": _dcFixture1.DatacenterDisplayName,
				"read_config_from_file":   fileName,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "missing all",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
		{
			Name: "missing read_config_from_file",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
			}),
			Good: false,
		},
	}

	command.TestCreateCommand(datacenterCreateCmd, cases, client, t)
}

func TestDatacenterYamlUnmarshal(t *testing.T) {
	RegisterTestingT(t)

	var dcConf metalcloud.DatacenterWithConfig

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")
	content, err := configuration.ReadInputFromFile(filepath.Join(basePath, "examples", "datacenter.yaml"))
	Expect(err).To(BeNil())
	Expect(content).NotTo(BeNil())

	err = yaml.Unmarshal(content, &dcConf)
	Expect(err).To(BeNil())

	Expect(dcConf.Config.ServerRegisterUsingGeneratedIPMICredentialsEnabled).To(BeTrue())
}

func TestDatacenterUpdate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var dcConf metalcloud.DatacenterWithConfig
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")
	fileName := filepath.Join(basePath, "examples", "datacenter.yaml")
	content, err := configuration.ReadInputFromFile(fileName)
	Expect(err).To(BeNil())
	Expect(content).NotTo(BeNil())

	err = yaml.Unmarshal(content, &dcConf)
	Expect(err).To(BeNil())

	client.EXPECT().
		DatacenterUpdateFromDatacenterWithConfig(dcConf).
		Return(&dcConf, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(_userFixture1.UserEmail).
		Return(&_userFixture1, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "dc-create-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":       dcConf.Metadata.DatacenterName,
				"read_config_from_file": fileName,
				"format":                "json",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "missing all",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
		{
			Name: "missing read_config_from_file",
			Cmd: command.MakeCommand(map[string]interface{}{
				"datacenter_name":         _dcFixture1.DatacenterName,
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
		DatacenterWithConfigGet(_dcFixture1.DatacenterName).
		Return(&_dcFixture1WithConfig, nil).
		AnyTimes()

	var dcConf metalcloud.DatacenterWithConfig

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "..")
	content, err := configuration.ReadInputFromFile(filepath.Join(basePath, "examples", "datacenter.yaml"))
	Expect(err).To(BeNil())
	Expect(content).NotTo(BeNil())

	err = yaml.Unmarshal(content, &dcConf)
	Expect(err).To(BeNil())

	client.EXPECT().
		DatacenterWithConfigGet(dcConf.Metadata.DatacenterName).
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

	_, err = datacenterGetCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name": _dcFixture1.DatacenterName,
	})

	ret, err := datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring(_dcFixture1.DatacenterName))

	//verify config url present in output
	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name":   _dcFixture1.DatacenterName,
		"return_config_url": true,
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("nfs://test-nfs/"))

	//verify json format
	cmd = command.MakeCommand(map[string]interface{}{
		"datacenter_name": _dcFixture1.DatacenterName,
		"format":          "json",
	})
	ret, err = datacenterGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	t.Logf("json: %s", string(ret))
	var m metalcloud.DatacenterWithConfig
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

}

var _dcFixture1 = metalcloud.Datacenter{
	DatacenterName:          "test",
	DatacenterDisplayName:   "datacenterDisplayName",
	UserID:                  1,
	DatacenterIsMaster:      false,
	DatacenterIsMaintenance: false,
	DatacenterHidden:        true,
	DatacenterTags:          []string{"t1", "t2"},
	DatacenterNameParent:    "test",
}

var _dcFixture1WithConfig = metalcloud.DatacenterWithConfig{
	Metadata: metalcloud.Datacenter{
		DatacenterName:          "test",
		DatacenterDisplayName:   "datacenterDisplayName",
		UserID:                  1,
		DatacenterIsMaster:      false,
		DatacenterIsMaintenance: false,
		DatacenterHidden:        true,
		DatacenterTags:          []string{"t1", "t2"},
		DatacenterNameParent:    "test",
	},
	Config: metalcloud.DatacenterConfig{
		NFSServer: "nfs://test-nfs/",
	},
}

var _dcFixture2 = metalcloud.Datacenter{
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