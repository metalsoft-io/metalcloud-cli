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
	"github.com/metalsoft-io/metalcloud-cli/internal/fixtures"
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

	err := json.Unmarshal([]byte(fixtures.SwitchDeviceControllerFixture1), &swCtrl)
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
	f.WriteString(fixtures.SwitchDeviceControllerFixture1)
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
	f.WriteString(fixtures.SwitchDeviceControllerFixture1)
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
	err := json.Unmarshal([]byte(fixtures.SwitchDeviceControllerFixture1), &sw)
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

	err := json.Unmarshal([]byte(fixtures.SwitchDeviceFixture2), &sw)
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
