package main

import (
	"encoding/json"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestInstanceArrayCreate(t *testing.T) {
	RegisterTestingT(t)

	responseBody = `{"result": ` + _InstanceArrayGetfixture1 + ` ,"jsonrpc": "2.0","id": 0}`

	client, err := metalcloud.GetMetalcloudClient("user", "APIKey", httpServer.URL, false)
	Expect(err).To(BeNil())

	i := 10
	s := "test"
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id":             &i,
			"instance_array_instance_count": &i,
			"volume_template_id":            &i,
			"instance_array_ram_gbytes":     &i,
			"instance_array_label":          &s,
		},
	}

	ret, err1 := instanceArrayCreateCmd(&cmd, client)

	Expect(err1).To(BeNil())
	Expect(ret).To(BeEmpty())

	reqBody := (<-requestChan).body
	Expect(reqBody).NotTo(BeNil())
	var m map[string]interface{}

	err = json.Unmarshal([]byte(reqBody), &m)
	params := m["params"].([]interface{})
	infraIDSubmitted := int(params[0].(float64))
	Expect(int(infraIDSubmitted)).To(Equal(i))

	objSubmitted := params[1].(map[string]interface{})

	Expect(int(objSubmitted["instance_array_instance_count"].(float64))).To(Equal(i))
	Expect(int(objSubmitted["instance_array_ram_gbytes"].(float64))).To(Equal(i))
	Expect(int(objSubmitted["volume_template_id"].(float64))).To(Equal(i))
	Expect(objSubmitted["instance_array_label"]).To(Equal(s))
	Expect(err).To(BeNil())

}

func TestInstanceArrayListCmdHumanReadable(t *testing.T) {
	RegisterTestingT(t)

	responseBody = `{"result": ` + _InstanceArraysFixture1 + `,"jsonrpc": "2.0","id": 0}`

	client, err := metalcloud.GetMetalcloudClient("user", "APIKey", httpServer.URL, false)
	Expect(err).To(BeNil())

	infraID := 10
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id": &infraID,
		},
	}

	ret, err1 := instanceArrayListCmd(&cmd, client)
	Expect(err1).To(BeNil())

	reqBody := (<-requestChan).body
	Expect(reqBody).NotTo(BeNil())

	Expect(ret).NotTo(BeEmpty())

}

func TestInstanceArrayEdit(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:        11,
		InstanceArrayLabel:     "testia",
		InfrastructureID:       infra.InfrastructureID,
		InstanceArrayOperation: &iao,
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

	i := 10
	newlabel := "newlabel"
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id":    &ia.InstanceArrayID,
			"instance_array_label": &newlabel,
			"volume_template_id":   &i,
		},
	}

	expectedOperationObject := iao
	expectedOperationObject.InstanceArrayLabel = "newlabel"
	expectedOperationObject.VolumeTemplateID = i

	client.EXPECT().
		InstanceArrayEdit(ia.InstanceArrayID, expectedOperationObject, nil, nil, nil, nil).
		Return(&ia, nil).
		AnyTimes()

	ret, err := instanceArrayEditCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

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
			"infrastructure_id": &infra.InfrastructureID,
			"format":            &format,
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
