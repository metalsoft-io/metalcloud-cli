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

func TestVolumeTemplateCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	vt := metalcloud.VolumeTemplate{
		VolumeTemplateID: 10,
	}

	client.EXPECT().
		VolumeTemplateCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&vt, nil).
		MinTimes(1)

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id": 11,
				"label":    "ass",
			}),
			good: true,
			id:   vt.VolumeTemplateID,
		},
	}

	testCreateCommand(volumeTemplateCreateCmd, cases, client, t)
}
