package network

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	. "github.com/onsi/gomega"
)

func TestNetworkProfileListCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	dc := metalcloud.Datacenter{
		DatacenterID:   100,
		DatacenterName: "test",
	}
	vlanid1 := 1
	vlanid2 := 2
	np1 := metalcloud.NetworkProfile{
		NetworkProfileID:    10,
		NetworkProfileLabel: "test1",
		NetworkType:         "wan",
		NetworkProfileVLANs: []metalcloud.NetworkProfileVLAN{
			{
				VlanID: &vlanid1,
			},
			{
				VlanID: &vlanid2,
			},
		},
		NetworkProfileCreatedTimestamp: "",
		NetworkProfileUpdatedTimestamp: "",
	}

	np2 := metalcloud.NetworkProfile{
		NetworkProfileID:    11,
		NetworkProfileLabel: "test2",
		NetworkType:         "wan",
	}

	npList := map[int]metalcloud.NetworkProfile{
		np1.NetworkProfileID: np1,
		np2.NetworkProfileID: np2,
	}

	client.EXPECT().
		NetworkProfiles(dc.DatacenterName).
		Return(&npList, nil).
		AnyTimes()

	//check the json output
	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}
	cmd.Arguments["datacenter"] = &dc.DatacenterName
	ret, err := networkProfileListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	vlans := strconv.Itoa(*np1.NetworkProfileVLANs[0].VlanID) + "," + strconv.Itoa(*np1.NetworkProfileVLANs[1].VlanID)

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(np1.NetworkProfileID))
	Expect(r["LABEL"].(string)).To(ContainSubstring(np1.NetworkProfileLabel))
	Expect(r["NETWORK TYPE"].(string)).To(Equal(np1.NetworkType))
	Expect(r["VLANs"].(string)).To(Equal(vlans))
	Expect(r["CREATED"].(string)).To(Equal(np1.NetworkProfileCreatedTimestamp))
	Expect(r["UPDATED"].(string)).To(Equal(np1.NetworkProfileUpdatedTimestamp))

	//check the csv output
	format = "csv"
	cmd.Arguments["format"] = &format
	ret, err = networkProfileListCmd(&cmd, client)
	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()

	Expect(err).To(BeNil())
	Expect(csv[1][0]).To(Equal(strconv.Itoa(np1.NetworkProfileID)))
	Expect(csv[1][1]).To(Equal(np1.NetworkProfileLabel))

	//check the human readable output, just check for not empty

	format = "json"
	cmd.Arguments["format"] = &format
	ret, err = networkProfileListCmd(&cmd, client)
	Expect(ret).NotTo(BeEmpty())
	Expect(err).To(BeNil())

	dcName := "tes"
	cmd.Arguments["datacenter"] = &dcName

	client.EXPECT().
		NetworkProfiles(dcName).
		Return(&npList, fmt.Errorf("testerror")).
		AnyTimes()

	_, err = networkProfileListCmd(&cmd, client)
	Expect(err).NotTo(BeNil())
}

func TestNetworkProfileVlansListCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	vlanid1 := 1
	vlanid2 := 2

	np := metalcloud.NetworkProfile{
		NetworkProfileID: 10,
		NetworkProfileVLANs: []metalcloud.NetworkProfileVLAN{
			{
				VlanID:   &vlanid1,
				PortMode: "trunk",
				ExternalConnectionIDs: []int{
					1,
				},
			},
			{
				VlanID:   &vlanid2,
				PortMode: "trunk",
			},
		},
	}

	ec := metalcloud.ExternalConnection{
		ExternalConnectionID:    1,
		ExternalConnectionLabel: "test",
	}

	client.EXPECT().
		ExternalConnectionGet(ec.ExternalConnectionID).
		Return(&ec, nil).
		AnyTimes()

	client.EXPECT().
		NetworkProfileGet(np.NetworkProfileID).
		Return(&np, nil).
		AnyTimes()

	//check the json output
	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}
	cmd.Arguments["network_profile_id"] = &np.NetworkProfileID
	ret, err := networkProfileVlansListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	ecString := ec.ExternalConnectionLabel + " (#" + strconv.Itoa(ec.ExternalConnectionID) + ")"

	r := m[0].(map[string]interface{})
	Expect(r["VLAN"].(string)).To(Equal(strconv.Itoa(*np.NetworkProfileVLANs[0].VlanID)))
	Expect(r["Port mode"].(string)).To(Equal(np.NetworkProfileVLANs[0].PortMode))
	Expect(r["External connections"].(string)).To(Equal(ecString))
	Expect(r["Provision subnet gateways"].(bool)).To(Equal(np.NetworkProfileVLANs[0].ProvisionSubnetGateways))

	format = "csv"
	cmd.Arguments["format"] = &format
	ret, err = networkProfileVlansListCmd(&cmd, client)
	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()

	Expect(err).To(BeNil())
	Expect(csv[1][0]).To(Equal(strconv.Itoa(*np.NetworkProfileVLANs[0].VlanID)))
	Expect(csv[1][1]).To(Equal(np.NetworkProfileVLANs[0].PortMode))

	//check the human readable output, just check for not empty

	format = "json"
	cmd.Arguments["format"] = &format
	ret, err = networkProfileVlansListCmd(&cmd, client)
	Expect(ret).NotTo(BeEmpty())
	Expect(err).To(BeNil())

	iaId := 12
	cmd.Arguments["network_profile_id"] = &iaId

	client.EXPECT().
		NetworkProfileGet(iaId).
		Return(&np, fmt.Errorf("testerror")).
		AnyTimes()

	client.EXPECT().
		NetworkProfileGet(iaId).
		Return(nil, fmt.Errorf("testerror")).
		AnyTimes()

	_, err = networkProfileVlansListCmd(&cmd, client)
	Expect(err).NotTo(BeNil())
}

func TestNetworkProfileGetCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	sbnPoolId := 12
	np := metalcloud.NetworkProfile{
		NetworkProfileID: 10,
		NetworkProfileVLANs: []metalcloud.NetworkProfileVLAN{
			{
				SubnetPools: []metalcloud.NetworkProfileSubnetPool{{
					SubnetPoolID:   &sbnPoolId,
					SubnetPoolType: "ipv4",
					SubnetPoolProvidesDefaultRoute: false,
				}},
				ExternalConnectionIDs: []int{
					10,
				},
			},
			{
				SubnetPools: []metalcloud.NetworkProfileSubnetPool{{
					SubnetPoolID:   nil, //this is important as it crashed previously
					SubnetPoolType: "ipv4",
					SubnetPoolProvidesDefaultRoute: false,
				}},
				ExternalConnectionIDs: []int{
					10,
				},
			},
		},
	}

	client.EXPECT().
		NetworkProfileGet(np.NetworkProfileID).
		Return(&np, nil).
		AnyTimes()

	extC := metalcloud.ExternalConnection{
		ExternalConnectionID:          10,
		ExternalConnectionDescription: "asdasd",
	}

	client.EXPECT().
		ExternalConnectionGet(10).
		Return(&extC, nil).
		AnyTimes()

	subnPool := metalcloud.SubnetPool{
		SubnetPoolID:                  12,
		SubnetPoolPrefixHumanReadable: "192.168.0.1",
		SubnetPoolPrefixSize:          24,
	}

	client.EXPECT().
		SubnetPoolGet(12).
		Return(&subnPool, nil).
		AnyTimes()

	expectedFirstRow := map[string]interface{}{
		"ID":         "#10",
		"DATACENTER": "test",
		"LABEL":      "test",
	}

	cases := []command.CommandTestCase{
		{
			Name: "np-get-json1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_profile_id": 10,
				"format":             "json",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "np-get-yaml1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_profile_id": 10,
				"format":             "yaml",
			}),
			Good: true,
			Id:   1,
		},
		{
			Name: "no id",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
	}

	command.TestGetCommand(networkProfileGetCmd, cases, client, expectedFirstRow, t)
}

func TestNetworkProfileDeleteCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	np := metalcloud.NetworkProfile{
		NetworkProfileID: 10,
	}

	client.EXPECT().
		NetworkProfileGet(np.NetworkProfileID).
		Return(&np, nil).
		AnyTimes()

	client.EXPECT().
		NetworkProfileDelete(np.NetworkProfileID).
		Return(nil).
		AnyTimes()

	autoconf := true
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"network_profile_id": &np.NetworkProfileID,
			"autoconfirm":        &autoconf,
		},
	}

	ret, err := networkProfileDeleteCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())
}

func TestNetworkProfileAssociateToInstanceArrayCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	ia := metalcloud.InstanceArray{
		InstanceArrayID: 15,
	}

	net := metalcloud.Network{
		NetworkID: 10,
	}

	np := metalcloud.NetworkProfile{
		NetworkProfileID: 100,
	}

	result := map[int]int{
		net.NetworkID: np.NetworkProfileID,
	}

	client.EXPECT().
		InstanceArrayNetworkProfileSet(ia.InstanceArrayID, net.NetworkID, np.NetworkProfileID).
		Return(&result, nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_profile_id": 100,
				"network_id":         10,
				"instance_array_id":  15,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "associate a network profile, missing network_id",
			Cmd: command.MakeCommand(map[string]interface{}{
				"network_profile_id": 100,
				"instance_array_id":  15,
			}),
			Good: false,
			Id:   0,
		},
	}

	command.TestCreateCommand(networkProfileAssociateToInstanceArrayCmd, cases, client, t)
}
