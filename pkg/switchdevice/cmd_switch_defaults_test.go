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
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"

	"github.com/metalsoft-io/metalcloud-cli/internal/command"
)

func TestSwitchDefaultsList(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := []metalcloud.SwitchDeviceDefaults{
		{
			NetworkEquipmentDefaultsID:           100,
			NetworkEquipmentSerialNumber:         "ABCDEF",
			NetworkEquipmentManagementMacAddress: "00:00:00:00:00:00",
		},
	}

	client.EXPECT().
		SwitchDeviceDefaults("test").
		Return(&list, nil).
		AnyTimes()

	datacenter := "test"
	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"datacenter_name": &datacenter,
			"format": &format,
		},
	}

	ret, err := switchDefaultsListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(100))
	Expect(r["Serial Number"]).To(Equal("ABCDEF"))
	Expect(r["Management MAC"]).To(Equal("00:00:00:00:00:00"))
}

func TestSwitchDefaultsCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var swCtrl metalcloud.SwitchDeviceDefaults

	err := json.Unmarshal([]byte(_SwitchDefaultsFixture), &swCtrl)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SwitchDeviceDefaultsCreate(gomock.Any()).
		Return(nil).
		AnyTimes()

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_SwitchDefaultsFixture)
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
				"read_config_from_file": filepath.Join(basePath, "examples", "switch_defaults.yaml"),
				"format":                "yaml",
			}),
			Good: true,
		},
	}

	command.TestCreateCommand(switchDefaultsCreateCmd, cases, client, t)
}

func TestSwitchControllerDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	ids := []int{100, 101}
	client.EXPECT().
		SwitchDeviceDefaultsDelete(ids).
		Return(nil).
		AnyTimes()

	cmd := command.MakeCommand(map[string]interface{}{
		"switch_defaults_ids": "100, 101",
	})

	_, err := switchDefaultsDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())
}

const _SwitchDefaultsFixture = "{\"datacenter_name\":\"test\",\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"network_equipment_serial_number\":\"ABCDEF\"}"
