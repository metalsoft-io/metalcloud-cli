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

func TestInfrastructureRevertCmd(t *testing.T) {
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

	iList := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel: infra,
	}

	client.EXPECT().
		Infrastructures().
		Return(&iList, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	autoconfirm := true
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label": &infra.InfrastructureID,
			"autoconfirm":                &autoconfirm,
		},
	}

	client.EXPECT().
		InfrastructureOperationCancel(infra.InfrastructureID).
		Return(nil).
		Times(1)

	ret, err := infrastructureRevertCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

}

func TestInfrastructureDeployCmd(t *testing.T) {
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
	//bFalse := true
	bTrue := true
	timeout := 256
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label":    &infra.InfrastructureID,
			"allow_data_loss":               &bTrue,
			"no_attempt_soft_shutdown":      &bTrue,
			"soft_shutdown_timeout_seconds": &timeout,
		},
	}

	expectedShutdownOptions := metalcloud.ShutdownOptions{
		HardShutdownAfterTimeout:   true,
		AttemptSoftShutdown:        false,
		SoftShutdownTimeoutSeconds: timeout,
	}

	client.EXPECT().
		InfrastructureDeploy(infra.InfrastructureID, expectedShutdownOptions, true, false).
		Return(nil).
		Times(1)

	//test first without confirmation
	ret, err := infrastructureDeleteCmd(&cmd, client)
	Expect(ret).To(Equal(""))
	Expect(err).NotTo(BeNil()) //should throw error indicating confirmation not given
	Expect(err.Error()).To(Equal("Operation not confirmed. Aborting"))

	cmd.Arguments["autoconfirm"] = &bTrue

	ret, err = infrastructureDeployCmd(&cmd, client)
	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil()) //should be nil

}

func TestInfrastructureGetCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	infra2 := metalcloud.Infrastructure{
		InfrastructureID:    10003,
		InfrastructureLabel: "testinfra2",
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
		InfrastructureGet(infra2.InfrastructureID).
		Return(&infra2, nil).
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
		InstanceArrays(gomock.Any()).
		Return(&iaList, nil).
		AnyTimes()

	daList := map[string]metalcloud.DriveArray{
		da.DriveArrayLabel + ".vanilla": da,
	}
	client.EXPECT().
		DriveArrays(gomock.Any()).
		Return(&daList, nil).
		AnyTimes()

	ret, err := infrastructureGetCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal("edited"))
	Expect(r["LABEL"].(string)).To(Equal(iao.InstanceArrayLabel))

	//test with label instead of id

	infraList := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel:  infra,
		infra2.InfrastructureLabel: infra2,
	}

	client.EXPECT().
		Infrastructures().
		Return(&infraList, nil).
		AnyTimes()

	cmd = Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label": &infra.InfrastructureLabel,
			"format":                     &format,
		},
	}

	ret, err = infrastructureGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())
}

func TestInfrastructureListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	iao := metalcloud.InfrastructureOperation{
		InfrastructureLabel: "testinfra-edited",
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID:            10002,
		InfrastructureLabel:         "testinfra",
		InfrastructureServiceStatus: "active",
		InfrastructureOperation:     iao,
	}

	iao2 := metalcloud.InfrastructureOperation{
		InfrastructureLabel: "testinfra-edited",
	}

	infra2 := metalcloud.Infrastructure{
		InfrastructureID:            10003,
		InfrastructureLabel:         "testinfra2",
		InfrastructureServiceStatus: "ordered",
		InfrastructureOperation:     iao2,
	}

	infraList := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel:  infra,
		infra2.InfrastructureLabel: infra2,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		Infrastructures().
		Return(&infraList, nil).
		AnyTimes()

	//test plaintext return
	format := ""
	cmd := Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err := infrastructureListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))
	Expect(ret).To(ContainSubstring(infra.InfrastructureOperation.InfrastructureLabel))
	Expect(ret).To(ContainSubstring(infra2.InfrastructureOperation.InfrastructureLabel))

	//test json return
	format = "json"
	cmd.Arguments["format"] = &format

	ret, err = infrastructureListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})

	Expect(r["STATUS"].(string)).To(SatisfyAny(
		Equal(infra.InfrastructureServiceStatus),
		Equal(infra2.InfrastructureServiceStatus),
	))
	Expect(r["LABEL"].(string)).To(SatisfyAny(
		Equal(infra.InfrastructureOperation.InfrastructureLabel),
		Equal(infra2.InfrastructureOperation.InfrastructureLabel),
	))

	//test csv return
	format = "csv"
	cmd.Arguments["format"] = &format

	ret, err = infrastructureListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(Not(Equal("")))

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()

	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", infra.InfrastructureID)))
	Expect(csv[1][2]).To(Equal(infra.UserEmailOwner))
	Expect(csv[2][1]).To(Equal(infra2.InfrastructureOperation.InfrastructureLabel))
}

func TestGetInfrastructureFromCommand(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	infra2 := metalcloud.Infrastructure{
		InfrastructureID:    10003,
		InfrastructureLabel: "testinfra2",
	}

	infra3 := metalcloud.Infrastructure{
		InfrastructureID:    10004,
		InfrastructureLabel: "testinfra",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra2.InfrastructureID).
		Return(&infra2, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra3.InfrastructureID).
		Return(&infra2, nil).
		AnyTimes()

	infraListAmbigous := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel:        infra,
		infra2.InfrastructureLabel:       infra2,
		infra3.InfrastructureLabel + "1": infra3,
	}

	client.EXPECT().
		Infrastructures().
		Return(&infraListAmbigous, nil).
		AnyTimes()

	//check with id
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label": &infra.InfrastructureID,
		},
	}

	i, err := getInfrastructureFromCommand(&cmd, client)
	Expect(err).To(BeNil())
	Expect(i.InfrastructureID).To(Equal(infra.InfrastructureID))

	//check with ambiguous label
	cmd = Command{
		Arguments: map[string]interface{}{
			"infrastructure_id_or_label": &infra.InfrastructureLabel,
		},
	}

	i, err = getInfrastructureFromCommand(&cmd, client)
	Expect(err).NotTo(BeNil())

	//check with wrong label
	blablah := "asdasdasdasd"
	cmd.Arguments["infrastructure_id_or_label"] = &blablah

	i, err = getInfrastructureFromCommand(&cmd, client)
	Expect(err).NotTo(BeNil())

	//check with correct label
	cmd.Arguments["infrastructure_id_or_label"] = &infra2.InfrastructureLabel

	i, err = getInfrastructureFromCommand(&cmd, client)
	Expect(err).To(BeNil())
	Expect(i.InfrastructureID).To(Equal(infra2.InfrastructureID))

}
