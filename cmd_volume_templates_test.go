package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"

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
	cmd := Command{
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
	cmd = Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = volumeTemplatesListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"
	cmd = Command{
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
		VolumeTemplateID:                   10,
		VolumeTemplateLabel:                "centos7-10",
		VolumeTemplateSizeMBytes:           40960,
		VolumeTemplateDisplayName:          "Centos7",
		VolumeTemplateDescription:          "centos7-10",
		VolumeTemplateLocalDiskSupported:   true,
		VolumeTemplateBootMethodsSupported: "pxe_iscsi",
		VolumeTemplateDeprecationStatus:    "not_deprecated",
		VolumeTemplateRepoURL:              "centos7_repo_url",
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

	cases := []CommandTestCase{
		{
			name: "good",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":     11,
				"label":        vt.VolumeTemplateLabel,
				"display_name": vt.VolumeTemplateDisplayName,
				"description":  vt.VolumeTemplateDescription,
			}),
			good: true,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"label":           vt.VolumeTemplateLabel,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			good: true,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"label":           vt.VolumeTemplateLabel,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"tags":            "tag1,tag2",
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			good: true,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing label",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing description",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":     11,
				"label":        vt.VolumeTemplateLabel,
				"display_name": vt.VolumeTemplateDisplayName,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing display name",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":    11,
				"label":       vt.VolumeTemplateLabel,
				"description": vt.VolumeTemplateDescription,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing os_type",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"label":           vt.VolumeTemplateLabel,
				"description":     vt.VolumeTemplateDescription,
				"os_version":      vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing os_version",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":        11,
				"display_name":    vt.VolumeTemplateDisplayName,
				"label":           vt.VolumeTemplateLabel,
				"description":     vt.VolumeTemplateDescription,
				"os_type":         vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_architecture": vt.VolumeTemplateOperatingSystem.OperatingSystemArchitecture,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
		{
			name: "missing os_architecture",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id":     11,
				"display_name": vt.VolumeTemplateDisplayName,
				"label":        vt.VolumeTemplateLabel,
				"description":  vt.VolumeTemplateDescription,
				"os_type":      vt.VolumeTemplateOperatingSystem.OperatingSystemType,
				"os_version":   vt.VolumeTemplateOperatingSystem.OperatingSystemVersion,
			}),
			good: false,
			id:   vt.VolumeTemplateID,
		},
	}

	testCreateCommand(volumeTemplateCreateFromDriveCmd, cases, client, t)
}
