package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestNetworkListCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	iaNetworkProfiles := map[int]int{
		10: 100,
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    1234,
		InstanceArrayLabel: "test",
		InstanceArrayInterfaces: []metalcloud.InstanceArrayInterface{
			{
				InstanceArrayID:             1234,
				NetworkID:                   10,
				InstanceArrayInterfaceIndex: 0,
			},
			{
				InstanceArrayID:             1234,
				NetworkID:                   10,
				InstanceArrayInterfaceIndex: 1,
			},
			{
				InstanceArrayID:             1234,
				NetworkID:                   10,
				InstanceArrayInterfaceIndex: 2,
			},
			{
				InstanceArrayID:             1234,
				NetworkID:                   10,
				InstanceArrayInterfaceIndex: 3,
			},
		},
	}

	nw := metalcloud.Network{
		NetworkType: "wan",
		NetworkID:   10,
	}

	np := metalcloud.NetworkProfile{
		NetworkProfileLabel: "test",
		NetworkProfileID:    100,
	}

	client.EXPECT().
		NetworkProfileListByInstanceArray(ia.InstanceArrayID).
		Return(&iaNetworkProfiles, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		NetworkGet(nw.NetworkID).
		Return(&nw, nil).
		AnyTimes()

	client.EXPECT().
		NetworkProfileGet(np.NetworkProfileID).
		Return(&np, nil).
		AnyTimes()

	//check the json output
	format := "json"
	id := ia.InstanceArrayID
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &id,
			"format":                     &format,
		},
	}

	ret, err := networkListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	nwIndex := "#" + strconv.Itoa(ia.InstanceArrayInterfaces[0].InstanceArrayInterfaceIndex+1)
	nwType := nw.NetworkType + "(#" + strconv.Itoa(nw.NetworkID) + ")"

	r := m[0].(map[string]interface{})
	Expect(r["Port"].(string)).To(Equal(nwIndex))
	Expect(r["Network"].(string)).To(Equal(nwType))
	Expect(r["Profile"].(string)).To(Equal(np.NetworkProfileLabel + " (#" + strconv.Itoa(np.NetworkProfileID) + ")"))

	//check the csv output
	format = "csv"
	cmd.Arguments["format"] = &format
	ret, err = networkListCmd(&cmd, client)
	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(err).To(BeNil())
	Expect(csv[1][0]).To(Equal(nwIndex))
	Expect(csv[1][1]).To(Equal(nwType))

	//check the human readable output, just check for not empty

	format = "json"
	cmd.Arguments["format"] = &format
	ret, err = networkListCmd(&cmd, client)
	Expect(ret).NotTo(BeEmpty())
	Expect(err).To(BeNil())

	i := 105
	iaId := fmt.Sprintf("%d", i)
	cmd.Arguments["instance_array_id_or_label"] = &iaId

	client.EXPECT().
		InstanceArrayGet(i).
		Return(&ia, fmt.Errorf("testerror")).
		AnyTimes()

	_, err = networkListCmd(&cmd, client)
	Expect(err).NotTo(BeNil())
}
