package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestInstanceArrayCreateCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayLabel:          "testia",
		InstanceArrayInstanceCount:  10,
		InstanceArrayProcessorCount: 10,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label":     &infra.InfrastructureID,
			"instance_array_label":           &ia.InstanceArrayLabel,
			"instance_array_instance_count":  &ia.InstanceArrayInstanceCount,
			"instance_array_processor_count": &ia.InstanceArrayProcessorCount,
		},
	}

	retIA := ia
	retIA.InstanceArrayID = 1222

	client.EXPECT().
		InstanceArrayCreate(infra.InfrastructureID, ia).
		Return(&retIA, nil).
		AnyTimes()

	//check with no return_id
	ret, err := instanceArrayCreateCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Equal(""))

	bTrue := true
	cmd.Arguments["return_id"] = &bTrue

	//check with return_id
	ret, err = instanceArrayCreateCmd(&cmd, client)

	Expect(ret).To(Equal(fmt.Sprintf("%d", retIA.InstanceArrayID)))
	Expect(err).To(BeNil())
}

func TestInstanceArrayEdit(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:         11,
		InstanceArrayLabel:      "testia",
		InstanceArrayBootMethod: "pxe",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:         11,
		InstanceArrayLabel:      "testia",
		InstanceArrayBootMethod: "pxe",
		InfrastructureID:        infra.InfrastructureID,
		InstanceArrayOperation:  &iao,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	i := "10"
	newlabel := "newlabel"
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label":  &ia.InstanceArrayID,
			"instance_array_label":        &newlabel,
			"volume_template_id_or_label": &i,
		},
	}

	expectedOperationObject := iao
	expectedOperationObject.InstanceArrayLabel = "newlabel"
	expectedOperationObject.VolumeTemplateID = 10

	client.EXPECT().
		InstanceArrayEdit(ia.InstanceArrayID, expectedOperationObject, nil, nil, nil, nil).
		Return(&ia, nil).
		AnyTimes()

	ret, err := instanceArrayEditCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

	//verify that untouched params stay untouched
	Expect(expectedOperationObject.InstanceArrayBootMethod).To(Equal(ia.InstanceArrayBootMethod))
}

func TestInstanceArrayListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label": &infra.InfrastructureID,
			"format":                     &format,
		},
	}

	iaList := map[string]metalcloud.InstanceArray{
		ia.InstanceArrayLabel + ".vanilla": ia,
	}
	client.EXPECT().
		InstanceArrays(infra.InfrastructureID).
		Return(&iaList, nil).
		Times(1)

	ret, err := instanceArrayListCmd(&cmd, client)

	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal("edited"))
	Expect(r["LABEL"].(string)).To(Equal(iao.InstanceArrayLabel))

}

func TestInstanceArrayDeleteCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayID,
		},
	}

	client.EXPECT().
		InstanceArrayDelete(ia.InstanceArrayID).
		Return(nil).
		Times(1)

	//test autoconfirm
	_, err := instanceArrayDeleteCmd(&cmd, client)
	Expect(err.Error()).To(Equal("Operation not confirmed. Aborting"))

	bTrue := true
	cmd.Arguments["autoconfirm"] = &bTrue

	ret, err := instanceArrayDeleteCmd(&cmd, client)

	Expect(err).To(BeNil())
	Expect(ret).To(BeEmpty())

}

func TestGetInstanceArrayFromCommand(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	iao2 := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia2",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia2 := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia2",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao2,
		InstanceArrayServiceStatus: "active",
	}

	iaList := map[string]metalcloud.InstanceArray{
		ia.InstanceArrayLabel + ".vanilla":  ia,
		ia2.InstanceArrayLabel + ".vanilla": ia2,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InstanceArrays(infra.InfrastructureID).
		Return(&iaList, nil).
		AnyTimes()

	//if requested with id
	client.EXPECT().
		InstanceArrayGet((ia.InstanceArrayID)).
		Return(&ia, nil).
		Times(1)

	//check with int
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayID,
			"infrastructure_id_or_label": &ia.InfrastructureID,
		},
	}

	ret, err := getInstanceArrayFromCommand("id", &cmd, client)

	Expect(err).To(BeNil())
	Expect(ret.InstanceArrayID).To(Equal(ia.InstanceArrayID))

	//check with label

	client.EXPECT().
		InstanceArrayGetByLabel(ia.InstanceArrayLabel).
		Return(&ia, nil).
		Times(1)

	cmd = Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayLabel,
			"infrastructure_id_or_label": &ia.InfrastructureID,
		},
	}

	ret, err = getInstanceArrayFromCommand("id", &cmd, client)

	Expect(err).To(BeNil())
	Expect(ret.InstanceArrayID).To(Equal(ia.InstanceArrayID))

}

func TestInstanceArrayGetCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	expectedFirstRow := map[string]interface{}{
		"ID":             "#10",
		"INFRASTRUCTURE": "test",
		"LABEL":          "test",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    10,
		InstanceArrayLabel: "test",
		InfrastructureID:   1,
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID: 1,
	}

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "ia-get-json1",
			cmd: MakeCommand(map[string]interface{}{
				"instance_array_id_or_label": 10,
				"format":                     "json",
			}),
			good: true,
			id:   1,
		},
		{
			name: "ia-get-yaml1",
			cmd: MakeCommand(map[string]interface{}{
				"instance_array_id_or_label": 10,
				"format":                     "yaml",
			}),
			good: true,
			id:   1,
		},
		{
			name: "no id",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
	}

	testGetCommand(instanceArrayGetCmd, cases, client, expectedFirstRow, t)
}

func TestInstanceArrayInstancesListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           0,
		InstanceArrayOperation:     nil,
		InstanceArrayServiceStatus: "active",
	}

	ips := []metalcloud.IP{
		{
			IPType:          "ipv4",
			IPHumanReadable: "192.168.1.1",
		},
	}

	n := metalcloud.Network{
		NetworkID:   105,
		NetworkType: "wan",
	}

	itfs := []metalcloud.InstanceInterface{
		{
			InstanceInterfaceLabel: "ef0",
			InstanceInterfaceIPs:   ips,
			NetworkID:              105,
		},
	}

	creds := metalcloud.InstanceCredentials{
		SSH: &metalcloud.SSH{
			Port:            22,
			Username:        "root",
			InitialPassword: "asda';l321",
		},
	}

	io := metalcloud.InstanceOperation{
		InstanceID:                 100,
		InstanceLabel:              "instance-100",
		InstanceSubdomain:          "instance-100.asdasd.asdasd.asd",
		InstanceSubdomainPermanent: "instance-100.bigstep.io",
		ServerTypeID:               106,
	}

	i := metalcloud.Instance{
		InstanceID:                 100,
		InstanceLabel:              "instance-100",
		InstanceSubdomain:          "instance-100.asdasd.asdasd.asd",
		InstanceSubdomainPermanent: "instance-100.bigstep.io",
		InstanceInterfaces:         itfs,
		ServerTypeID:               106,
		InstanceOperation:          io,
		InstanceCredentials:        creds,
	}

	ilist := map[string]metalcloud.Instance{
		"instance-100": i,
	}

	st := metalcloud.ServerType{
		ServerTypeID:          106,
		ServerTypeDisplayName: "M.40.256.12D",
	}

	iaNetworkProfiles := map[int]int{
		1: 10,
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10,
		InfrastructureLabel: "test",
	}

	client.EXPECT().
		InstanceArrayInstances((ia.InstanceArrayID)).
		Return(&ilist, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet((ia.InstanceArrayID)).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		NetworkGet((itfs[0].NetworkID)).
		Return(&n, nil).
		AnyTimes()

	client.EXPECT().
		ServerTypeGet((ilist["instance-100"].ServerTypeID)).
		Return(&st, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(gomock.Any()).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		NetworkProfileListByInstanceArray(ia.InstanceArrayID).
		Return(&iaNetworkProfiles, nil).
		AnyTimes()

	//test with text output
	format := "text"
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayID,
			"format":                     &format,
		},
	}

	ret, err := instanceArrayInstancesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring(ips[0].IPHumanReadable))
	Expect(ret).To(ContainSubstring(ilist["instance-100"].InstanceSubdomainPermanent))
	Expect(ret).To(ContainSubstring(st.ServerTypeDisplayName))

	//test with credentials
	bTrue := true
	cmd = Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayID,
			"format":                     &format,
			"show_credentials":           &bTrue,
		},
	}

	ret, err = instanceArrayInstancesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring(ips[0].IPHumanReadable))
	Expect(ret).To(ContainSubstring(ilist["instance-100"].InstanceSubdomainPermanent))
	Expect(ret).To(ContainSubstring(st.ServerTypeDisplayName))
	Expect(ret).To(ContainSubstring(creds.SSH.Username))
	Expect(ret).To(ContainSubstring(creds.SSH.InitialPassword))

	//test with json output
	format = "json"
	cmd.Arguments["format"] = &format

	ret, err = instanceArrayInstancesListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["ID"]).To(Equal(float64(ilist["instance-100"].InstanceID)))
	Expect(r["SUBDOMAIN"]).To(Equal(ilist["instance-100"].InstanceSubdomainPermanent))
	Expect(r["WAN_IP"]).To(Equal(ips[0].IPHumanReadable))
	Expect(r["DETAILS"]).To(ContainSubstring(st.ServerTypeDisplayName))

	//test with csv

	format = "csv"
	cmd.Arguments["format"] = &format

	ret, err = instanceArrayInstancesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()

	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", i.InstanceID)))
	Expect(csv[1][1]).To(Equal(i.InstanceSubdomainPermanent))
	Expect(csv[1][2]).To(Equal(ips[0].IPHumanReadable))

}
