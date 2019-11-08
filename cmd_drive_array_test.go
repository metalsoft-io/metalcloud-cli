package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
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

	i := 10
	s := "test"
	sEmpty := ""

	//correct config
	da := metalcloud.DriveArray{
		DriveArrayLabel:  "test",
		InstanceArrayID:  i,
		VolumeTemplateID: i,
	}

	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id":  &i,
			"drive_array_label":  &da.DriveArrayLabel,
			"volume_template_id": &i,
			"instance_array_id":  &i,
		},
	}

	//return error, see if it's thrown
	client.EXPECT().
		DriveArrayCreate(i, da).
		Return(&da, fmt.Errorf("testerr")).
		Times(1)

	_, err := driveArrayCreateCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

	//return success, check empty return
	client.EXPECT().
		DriveArrayCreate(i, da).
		Return(&da, nil).
		Times(1)

	ret, err := driveArrayCreateCmd(&cmd, client)
	Expect(ret).To(BeEmpty())
	Expect(err).To(BeNil())

	//return success, check id return
	bTrue := true
	cmd = Command{
		Arguments: map[string]interface{}{
			"infrastructure_id":  &i,
			"drive_array_label":  &da.DriveArrayLabel,
			"volume_template_id": &i,
			"instance_array_id":  &i,
			"return_id":          &bTrue,
		},
	}
	retDA := da
	retDA.DriveArrayID = 1001

	client.EXPECT().
		DriveArrayCreate(i, da).
		Return(&retDA, nil).
		Times(1)

	ret, err = driveArrayCreateCmd(&cmd, client)
	Expect(ret).To(Equal(fmt.Sprintf("%d", retDA.DriveArrayID)))
	Expect(err).To(BeNil())

	//test no infra id

	errorArguments := []map[string]interface{}{
		map[string]interface{}{
			//"infrastructure_id": &i,
			"volume_template_id": &s,
			"instance_array_id":  &i,
		},
		map[string]interface{}{
			"infrastructure_id": &i,
			//"volume_template_id": &i,
			"instance_array_id": &i,
		},
		map[string]interface{}{
			"infrastructure_id":  &i,
			"volume_template_id": &i,
			//"instance_array_id":  &i,
		},
		map[string]interface{}{
			"infrastructure_id":  &i,
			"volume_template_id": &i,
			"instance_array_id":  &i,
			"drive_array_label":  &sEmpty,
		},
	}

	client.EXPECT().
		DriveArrayCreate(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()

	//test all error scenarios
	for _, a := range errorArguments {

		_, err := driveArrayCreateCmd(&Command{Arguments: a}, client)

		Expect(err).NotTo(BeNil())
	}

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
		Times(1)

	ret, err := driveArrayEditCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

	//check missing values
	errorArguments := []map[string]interface{}{
		map[string]interface{}{
			//"drive_array_id": &i,
		},
	}

	//test all error scenarios
	for _, a := range errorArguments {

		_, err := driveArrayEditCmd(&Command{Arguments: a}, client)

		Expect(err).NotTo(BeNil())
	}

	//check catches error at get
	i = 100
	cmd.Arguments["drive_array_id"] = &i

	client.EXPECT().
		DriveArrayGet(i).
		Return(&da, fmt.Errorf("testerr")).
		Times(1)

	client.EXPECT().
		DriveArrayEdit(i, gomock.Any()).
		Return(nil, fmt.Errorf("testerr")).
		Times(1)

	_, err = driveArrayEditCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

	//check catches error at edit
	i = 101
	cmd.Arguments["drive_array_id"] = &i

	client.EXPECT().
		DriveArrayGet(i).
		Return(&da, nil).
		Times(1)

	client.EXPECT().
		DriveArrayEdit(i, gomock.Any()).
		Return(nil, fmt.Errorf("testerr")).
		Times(1)

	_, err = driveArrayEditCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

}

func TestDriveArrayListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	//check missing params
	errorArguments := []map[string]interface{}{
		map[string]interface{}{
			//"infrastructure_id": &i,
		},
	}

	//test all missing params scenarios
	for _, a := range errorArguments {

		_, err := driveArrayListCmd(&Command{Arguments: a}, client)
		Expect(err).NotTo(BeNil())
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:    11,
		InstanceArrayLabel: "testia",
		InfrastructureID:   infra.InfrastructureID,
	}

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID:                10,
		VolumeTemplateSizeMBytes:        10,
		VolumeTemplateLabel:             "testlabel",
		VolumeTemplateDescription:       "testdesc",
		VolumeTemplateDeprecationStatus: "not deprecated",
	}

	dao := metalcloud.DriveArrayOperation{
		DriveArrayID:           10,
		DriveArrayLabel:        "test-edited",
		InstanceArrayID:        ia.InstanceArrayID,
		InfrastructureID:       infra.InfrastructureID,
		VolumeTemplateID:       vt.VolumeTemplateID,
		DriveArrayCount:        101,
		DriveArrayDeployType:   "edit",
		DriveArrayDeployStatus: "not_started",
	}

	da := metalcloud.DriveArray{
		DriveArrayID:            10,
		DriveArrayLabel:         "test",
		InstanceArrayID:         ia.InstanceArrayID,
		InfrastructureID:        infra.InfrastructureID,
		VolumeTemplateID:        vt.VolumeTemplateID,
		DriveArrayCount:         101,
		DriveArrayOperation:     &dao,
		DriveArrayServiceStatus: "active",
	}

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
		VolumeTemplateGet(vt.VolumeTemplateID).
		Return(&vt, nil).
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

	//check the csv output
	format = "csv"
	cmd.Arguments["format"] = &format
	ret, err = driveArrayListCmd(&cmd, client)
	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(err).To(BeNil())
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", da.DriveArrayID)))
	Expect(csv[1][1]).To(Equal(da.DriveArrayOperation.DriveArrayLabel))

	//check the human readable output, just check for not empty

	format = "text"
	cmd.Arguments["format"] = &format
	ret, err = driveArrayListCmd(&cmd, client)
	Expect(ret).NotTo(BeEmpty())
	Expect(err).To(BeNil())

	//check that it catches drive array list error

	i := 101
	client.EXPECT().
		DriveArrays(i).
		Return(nil, fmt.Errorf("testerror")).
		AnyTimes()

	cmd.Arguments["infrastructure_id"] = &i
	_, err = driveArrayListCmd(&cmd, client)
	Expect(err).NotTo(BeNil())

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
