package subnetpool

import (
	"encoding/json"
	"os"
	"syscall"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

func TestSubnetPoolList(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	list := []metalcloud.SubnetPool{
		{
			SubnetPoolID:                  10,
			SubnetPoolPrefixHumanReadable: "asdads",
		},
	}

	stats := metalcloud.SubnetPoolUtilization{}
	client.EXPECT().
		SubnetPoolSearch("").
		Return(&list, nil).
		AnyTimes()

	client.EXPECT().
		SubnetPoolPrefixSizesStats(10).
		Return(&stats, nil).
		AnyTimes()

	expectedFirstRow := map[string]interface{}{
		"ID":     10,
		"PREFIX": "asdads/0",
	}

	command.TestListCommand(subnetPoolListCmd, nil, client, expectedFirstRow, t)

}

func TestSubnetCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	var sw metalcloud.SubnetPool

	err := json.Unmarshal([]byte(_subnetPoolFixture1), &sw)
	if err != nil {
		t.Error(err)
	}

	client.EXPECT().
		SubnetPoolCreate(gomock.Any()).
		Return(&sw, nil).
		AnyTimes()

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_subnetPoolFixture1)
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

	cases := []command.CommandTestCase{
		{
			Name: "sn-create-good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f.Name(),
				"format":                "json",
			}),
			Good: true,
			Id:   1309,
		},
		{
			Name: "sn-create-good-yaml",
			Cmd: command.MakeCommand(map[string]interface{}{
				"read_config_from_file": f2.Name(),
				"format":                "yaml",
			}),
			Good: true,
			Id:   1309,
		},
	}

	command.TestCreateCommand(subnetPoolCreateCmd, cases, client, t)

}

func TestSubnetGet(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	obj := metalcloud.SubnetPool{
		SubnetPoolID:                  100,
		SubnetPoolPrefixHumanReadable: "asdas",
	}

	obj2 := metalcloud.SubnetPoolUtilization{
		IPAddressesUsableCountFree: "10",
	}

	client.EXPECT().
		SubnetPoolGet(100).
		Return(&obj, nil).
		AnyTimes()

	client.EXPECT().
		SubnetPoolPrefixSizesStats(100).
		Return(&obj2, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "sw-get-json1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"subnet_pool_id": 100,
				"format":         "json",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "sw-get-json1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"subnet_pool_id": 100,
				"format":         "yaml",
			}),
			Good: true,
			Id:   1,
		},
	}

	expectedFirstRow := map[string]interface{}{
		"ID":         3675,
		"IDENTIFIER": "test",
	}

	command.TestGetCommand(subnetPoolGetCmd, cases, client, expectedFirstRow, t)

}

const _subnetPoolFixture1 = "{\"subnet_pool_id\":1309,\"user_id\":3675,\"subnet_pool_is_only_for_manual_allocation\":false,\"datacenter_name\":\"es-madrid\",\"subnet_pool_prefix_hex\":\"a53c0ee3\",\"subnet_pool_prefix_human_readable\":\"165.60.14.227\",\"subnet_pool_prefix_size\":25,\"subnet_pool_type\":\"ipv4\",\"subnet_pool_routable\":true,\"subnet_pool_destination\":\"WAN\",\"subnet_pool_netmask_human_readable\":\"255.255.255.128\",\"subnet_pool_netmask_hex\":\"ffffff80\",\"network_equipment_id\":null,\"subnet_pool_utilization_cached_json\":\"{\\\"prefix_count_free\\\": {\\\"27\\\": 4}, \\\"prefix_count_allocated\\\": [], \\\"ip_addresses_usable_count_free\\\": \\\"116\\\", \\\"ip_addresses_usable_count_allocated\\\": 0, \\\"ip_addresses_usable_free_percent_optimistic\\\": \\\"100\\\"}\",\"subnet_pool_cached_updated_timestamp\":\"2020-08-07T12:53:55Z\"}"
