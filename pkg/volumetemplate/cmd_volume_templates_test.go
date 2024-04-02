package volumetemplate

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"

	. "github.com/onsi/gomega"
)

func TestVolumeTemplatesListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID:                10,
		VolumeTemplateSizeMBytes:        10,
		VolumeTemplateLabel:             "testlabel",
		VolumeTemplateDescription:       "testdesc",
		VolumeTemplateDeprecationStatus: "not deprecated",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	vtList := map[string]metalcloud.VolumeTemplate{
		"centos":  vt,
		"centos3": vt,
	}

	client.EXPECT().
		VolumeTemplates().
		Return(&vtList, nil).
		AnyTimes()

	format := "json"
	cmd := command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err := volumeTemplatesListCmd(&cmd, client)

	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal(vt.VolumeTemplateDeprecationStatus))
	Expect(r["LABEL"].(string)).To(Equal(vt.VolumeTemplateLabel))

	//test plaintext
	format = ""
	cmd = command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = volumeTemplatesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"
	cmd = command.Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = volumeTemplatesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", vt.VolumeTemplateID)))
	Expect(csv[1][1]).To(Equal(vt.VolumeTemplateLabel))

}

func TestVolumeTemplateCreateFromDriveCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID:                      10,
		VolumeTemplateLabel:                   "centos7-10",
		VolumeTemplateSizeMBytes:              40960,
		VolumeTemplateDisplayName:             "Centos7",
		VolumeTemplateDescription:             "centos7-10",
		VolumeTemplateLocalDiskSupported:      true,
		VolumeTemplateBootMethodsSupported:    "pxe_iscsi",
		VolumeTemplateDeprecationStatus:       "not_deprecated",
		VolumeTemplateOsBootstrapFunctionName: "provisioner_os_cloudinit_prepare_centos",
		VolumeTemplateRepoURL:                 "centos7_repo_url",
		VolumeTemplateVersion:                 "0.0.0",
		VolumeTemplateOperatingSystem: metalcloud.OperatingSystem{
			OperatingSystemType:         "Centos",
			OperatingSystemVersion:      "7",
			OperatingSystemArchitecture: "x86_64",
		},
	}

	client.EXPECT().
		VolumeTemplateCreateFromDrive(gomock.Any(), gomock.Any()).
		Return(&vt, nil).
		MinTimes(1)

	cases := []command.CommandTestCase{
		{
			Name: "good",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":                   11,
				"label":                      vt.VolumeTemplateLabel,
				"display_name":               vt.VolumeTemplateDisplayName,
				"description":                vt.VolumeTemplateDescription,
				"os_bootstrap_function_name": vt.VolumeTemplateOsBootstrapFunctionName,
				"version":                    vt.VolumeTemplateVersion,
			}),
			Good: true,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"label":           vt.VolumeTemplateLabel,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			Good: true,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "good2",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"label":           vt.VolumeTemplateLabel,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"tags":            "tag1,tag2",
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			Good: true,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing label",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing description",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":     11,
				"label":        vt.VolumeTemplateLabel,
				"display_name": vt.VolumeTemplateDisplayName,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing display name",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":    11,
				"label":       vt.VolumeTemplateLabel,
				"description": vt.VolumeTemplateDescription,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing os_type",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"label":           vt.VolumeTemplateLabel,
				"description":     vt.VolumeTemplateDescription,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing os_version",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"label":           vt.VolumeTemplateLabel,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
		{
			Name: "missing os_architecture",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id":     11,
				"display_name": vt.VolumeTemplateDisplayName,
				"label":        vt.VolumeTemplateLabel,
				"description":  vt.VolumeTemplateDescription,
				"os_type":      vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":   vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
			}),
			Good: false,
			Id:   vt.VolumeTemplateID,
		},
	}

	command.TestCreateCommand(volumeTemplateCreateFromDriveCmd, cases, client, t)
}

func TestVolumeTemplateMakePrivateCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	vtl := map[string]metalcloud.VolumeTemplate{
		"vt1": vt,
	}

	user := metalcloud.User{
		UserID: 1,
	}

	user1 := metalcloud.User{
		UserEmail: "test",
	}

	client.EXPECT().
		VolumeTemplateGet(gomock.Any()).
		Return(&vt, nil).
		AnyTimes()

	client.EXPECT().
		VolumeTemplates().
		Return(&vtl, nil).
		MinTimes(1)

	client.EXPECT().
		UserGet(gomock.Any()).
		Return(&user, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(gomock.Any()).
		Return(&user1, nil).
		MinTimes(1)

	client.EXPECT().
		VolumeTemplateMakePrivate(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": 10,
				"user_id":             1,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "good2",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": "test",
				"user_id":             1,
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "good3",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": 10,
				"user_id":             "test",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "template not found",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": "test1",
				"user_id":             1,
			}),
			Good: false,
			Id:   0,
		},
		{
			Name: "missing template id or name",
			Cmd: command.MakeCommand(map[string]interface{}{
				"user_id": 1,
			}),
			Good: false,
			Id:   0,
		},
		{
			Name: "missing user id or email",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": "test",
			}),
			Good: false,
			Id:   0,
		},
	}

	command.TestCreateCommand(volumeTemplateMakePrivateCmd, cases, client, t)
}

func TestVolumeTemplateMakePublicCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	vtl := map[string]metalcloud.VolumeTemplate{
		"vt1": vt,
	}

	client.EXPECT().
		VolumeTemplateGet(gomock.Any()).
		Return(&vt, nil).
		AnyTimes()

	client.EXPECT().
		VolumeTemplates().
		Return(&vtl, nil).
		MinTimes(1)

	client.EXPECT().
		VolumeTemplateMakePublic(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name":        10,
				"os_bootstrap_function_name": "provisioner_os_bootstrap_dummy",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "good2",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name":        "test",
				"os_bootstrap_function_name": "provisioner_os_bootstrap_dummy",
			}),
			Good: true,
			Id:   0,
		},
		{
			Name: "template not found",
			Cmd: command.MakeCommand(map[string]interface{}{
				"template_id_or_name": "test1",
			}),
			Good: false,
			Id:   0,
		},
		{
			Name: "missing template id or name",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
			Id:   0,
		},
	}

	command.TestCreateCommand(volumeTemplateMakePublicCmd, cases, client, t)
}
