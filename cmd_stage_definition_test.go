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

func TestStageDefinitionsListCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	ab := metalcloud.AnsibleBundle{
		AnsibleBundleArchiveFilename: "asdads",
	}

	stage1 := metalcloud.StageDefinition{
		StageDefinitionID:    10,
		StageDefinitionLabel: "test",
		StageDefinition:      ab,
		StageDefinitionType:  "AnsibleBundle",
	}

	req := metalcloud.HTTPRequest{
		URL: "http://asdad/asdasd/ass",
	}

	stage2 := metalcloud.StageDefinition{
		StageDefinitionID:    11,
		StageDefinitionLabel: "test2",
		StageDefinition:      req,
		StageDefinitionType:  "HTTPRequest",
	}

	list := map[string]metalcloud.StageDefinition{
		"test1": stage1,
		"test2": stage2,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		StageDefinitions().
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err := stageDefinitionsListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	//Expect(int(r["ID"].(float64))).To(Equal(11))
	//Expect(r["LABEL"].(string)).To(Equal(stage1.StageDefinitionLabel))
	Expect(r).NotTo(BeNil())
	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = stageDefinitionsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = stageDefinitionsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 10)))
	Expect(csv[1][1]).To(Equal("test"))

}

/*
func TestStageDefinitionsGetCmd(t *testing.T) {
	RegisterTestingT(t)

	ctrl := gomock.NewController(t)

	ab := metalcloud.AnsibleBundle{
		AnsibleBundleArchiveFilename: "asdads",
	}

	stage1 := metalcloud.StageDefinition{
		StageDefinitionID:    10,
		StageDefinitionLabel: "test",
		StageDefinition:      ab,
		StageDefinitionType:  "AnsibleBundle",
	}

	req := metalcloud.HTTPRequest{
		URL: "http://asdad/asdasd/ass",
	}

	stage2 := metalcloud.StageDefinition{
		StageDefinitionID:    11,
		StageDefinitionLabel: "test2",
		StageDefinition:      req,
		StageDefinitionType:  "HTTPRequest",
	}

	list := map[string]metalcloud.StageDefinition{
		"test1": stage1,
		"test2": stage2,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		StageDefinitions().
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err := stageDefinitionGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	//Expect(int(r["ID"].(float64))).To(Equal(11))
	//Expect(r["LABEL"].(string)).To(Equal(stage1.StageDefinitionLabel))
	Expect(r).NotTo(BeNil())
	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = stageDefinitionsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err = stageDefinitionsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 10)))
	Expect(csv[1][1]).To(Equal("test"))

}
*/
