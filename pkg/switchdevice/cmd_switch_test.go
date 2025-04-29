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
	"github.com/metalsoft-io/metalcloud-cli/internal/fixtures"
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

	err := json.Unmarshal([]byte(fixtures.SwitchDeviceFixture1), &sw)
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
	f.WriteString(fixtures.SwitchDeviceFixture1)
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

func TestSwitchEditCmd(t *testing.T) {

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
	f.WriteString(fixtures.SwitchDeviceFixture1)
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
			_, err := switchEditCmd(&c.Cmd, client)
			if c.Good && err != nil {
				t.Error(err)
			}
		})
	}

}

func getSwitchFixture1() (metalcloud.SwitchDevice, error) {
	var sw metalcloud.SwitchDevice
	err := json.Unmarshal([]byte(fixtures.SwitchDeviceFixture1), &sw)
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
