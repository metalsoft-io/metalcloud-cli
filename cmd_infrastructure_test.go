package main

import (
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

	bFalse := false
	bTrue := true
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id":             &infra.InfrastructureID,
			"allow_data_loss":               &bTrue,
			"hard_shutdown_after_timeout":   &bTrue,
			"attempt_soft_shutdown":         &bFalse,
			"soft_shutdown_timeout_seconds": 256,
		},
	}

	expectedShutdownOptions := metalcloud.ShutdownOptions{
		HardShutdownAfterTimeout:   true,
		AttemptSoftShutdown:        false,
		SoftShutdownTimeoutSeconds: 256,
	}

	client.EXPECT().
		InfrastructureDeploy(infra.InfrastructureID, expectedShutdownOptions, true, false).
		Return(nil).
		Times(1)

	ret, err := infrastructureRevertCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

}
