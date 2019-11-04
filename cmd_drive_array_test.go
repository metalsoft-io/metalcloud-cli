package main

import (
	"encoding/json"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestDriveArrayCreate(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		DriveArrayCreate(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()

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

	ret, _ := driveArrayCreateCmd(&cmd, client)

	Expect(ret).To(Equal(""))

}

func TestDriveArrayDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
		InfrastructureID:   infra.InfrastructureID,
	}

	da := metalcloud.DriveArray{
		DriveArrayID:     10,
		DriveArrayLabel:  "test",
		InstanceArrayID:  ia.InstanceArrayID,
		InfrastructureID: infra.InfrastructureID,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayGet(da.DriveArrayID).
		Return(&da, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(da.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayDelete(da.DriveArrayID).
		Return(nil).
		AnyTimes()

	autoconf := true
	id := da.DriveArrayID
	cmd := Command{
		Arguments: map[string]interface{}{
			"drive_array_id": &id,
			"autoconfirm":    &autoconf,
		},
	}

	ret, err := driveArrayDeleteCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

}

func TestDriveArrayEdit(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
		InfrastructureID:   infra.InfrastructureID,
	}

	dao := metalcloud.DriveArrayOperation{
		DriveArrayID:     10,
		DriveArrayLabel:  "test",
		InstanceArrayID:  ia.InstanceArrayID,
		InfrastructureID: infra.InfrastructureID,
		DriveArrayCount:  101,
	}

	da := metalcloud.DriveArray{
		DriveArrayID:        10,
		DriveArrayLabel:     "test",
		InstanceArrayID:     ia.InstanceArrayID,
		InfrastructureID:    infra.InfrastructureID,
		DriveArrayCount:     101,
		DriveArrayOperation: &dao,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayGet(da.DriveArrayID).
		Return(&da, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(da.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	i := 10
	newlabel := "newlabel"
	cmd := Command{
		Arguments: map[string]interface{}{
			"drive_array_id":     &da.DriveArrayID,
			"drive_array_label":  &newlabel,
			"volume_template_id": &i,
		},
	}

	expectedOperationObject := dao
	expectedOperationObject.DriveArrayLabel = "newlabel"
	expectedOperationObject.VolumeTemplateID = i

	client.EXPECT().
		DriveArrayEdit(da.DriveArrayID, expectedOperationObject).
		Return(&da, nil).
		AnyTimes()

	ret, err := driveArrayEditCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

}

func TestDriveArrayListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
		InfrastructureID:   infra.InfrastructureID,
	}

	dao := metalcloud.DriveArrayOperation{
		DriveArrayID:           10,
		DriveArrayLabel:        "test-edited",
		InstanceArrayID:        ia.InstanceArrayID,
		InfrastructureID:       infra.InfrastructureID,
		DriveArrayCount:        101,
		DriveArrayDeployType:   "edit",
		DriveArrayDeployStatus: "not_started",
	}

	da := metalcloud.DriveArray{
		DriveArrayID:            10,
		DriveArrayLabel:         "test",
		InstanceArrayID:         ia.InstanceArrayID,
		InfrastructureID:        infra.InfrastructureID,
		DriveArrayCount:         101,
		DriveArrayOperation:     &dao,
		DriveArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayGet(da.DriveArrayID).
		Return(&da, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(da.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id": &ia.InfrastructureID,
			"format":            &format,
		},
	}

	daList := map[string]metalcloud.DriveArray{
		da.DriveArrayLabel + ".vanilla": da,
	}

	client.EXPECT().
		DriveArrays(infra.InfrastructureID).
		Return(&daList, nil).
		AnyTimes()

	ret, err := driveArrayListCmd(&cmd, client)

	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal("edited"))
	Expect(r["LABEL"].(string)).To(Equal(dao.DriveArrayLabel))

}

func TestDriveArrayDeleteCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
		InfrastructureID:   infra.InfrastructureID,
	}

	dao := metalcloud.DriveArrayOperation{
		DriveArrayID:           10,
		DriveArrayLabel:        "test-edited",
		InstanceArrayID:        ia.InstanceArrayID,
		InfrastructureID:       infra.InfrastructureID,
		DriveArrayCount:        101,
		DriveArrayDeployType:   "edit",
		DriveArrayDeployStatus: "not_started",
	}

	da := metalcloud.DriveArray{
		DriveArrayID:            10,
		DriveArrayLabel:         "test",
		InstanceArrayID:         ia.InstanceArrayID,
		InfrastructureID:        infra.InfrastructureID,
		DriveArrayCount:         101,
		DriveArrayOperation:     &dao,
		DriveArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayGet(da.DriveArrayID).
		Return(&da, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(da.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	cmd := Command{
		Arguments: map[string]interface{}{
			"drive_array_id": &da.DriveArrayID,
		},
	}

	client.EXPECT().
		DriveArrayDelete(da.DriveArrayID).
		Return(nil).
		AnyTimes()

	ret, err := driveArrayDeleteCmd(&cmd, client)

	Expect(err.Error()).To(Equal("Operation not confirmed. Aborting"))

	bTrue := true
	cmd.Arguments["autoconfirm"] = &bTrue

	ret, err = driveArrayDeleteCmd(&cmd, client)

	Expect(err).To(BeNil())
	Expect(ret).To(BeEmpty())

}
