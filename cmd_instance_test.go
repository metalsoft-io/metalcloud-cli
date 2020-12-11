package main

import (
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestInstanceCredentialsCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)
	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	instance := metalcloud.Instance{
		InstanceID:      110,
		InstanceArrayID: 10,
		InstanceCredentials: metalcloud.InstanceCredentials{
			SSH: &metalcloud.SSH{
				Username:        "testu",
				InitialPassword: "testp",
				Port:            22,
			},
		},
	}

	ia := metalcloud.InstanceArray{
		InstanceArraySubdomain: "tst",
		InstanceArrayID:        10,
	}

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10,
		InfrastructureLabel: "tsassd",
	}

	client.EXPECT().
		InstanceGet(gomock.Any()).
		Return(&instance, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(gomock.Any()).
		Return(&ia, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(gomock.Any()).
		Return(&infra, nil).
		AnyTimes()

	cmd := MakeCommand(map[string]interface{}{"instance_id": 110})

	ret, err := instanceCredentialsCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).To(ContainSubstring("ID"))

}
