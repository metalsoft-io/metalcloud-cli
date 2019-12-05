package main

import (
	"encoding/json"
	"fmt"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
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

	i := 10
	newlabel := "newlabel"
	cmd := Command{
		Arguments: map[string]interface{}{
			"instance_array_id_or_label": &ia.InstanceArrayID,
			"instance_array_label":       &newlabel,
			"volume_template_id":         &i,
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
