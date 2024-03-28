package stagedefinition

import (
	"os"
	"syscall"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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

	expectedFirstRow := map[string]interface{}{
		"ID":    10,
		"LABEL": "test",
	}

	command.TestListCommand(stageDefinitionsListCmd, nil, client, expectedFirstRow, t)

}

func TestStageDefinitionCreateCmdAnsible(t *testing.T) {
	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	stage1 := metalcloud.StageDefinition{
		StageDefinitionID:    10,
		StageDefinitionLabel: "test",
		StageDefinition:      metalcloud.AnsibleBundle{},
		StageDefinitionType:  "AnsibleBundle",
	}

	client.EXPECT().
		StageDefinitionCreate(gomock.Any()).
		Return(&stage1, nil).
		MinTimes(1)

	f, err := os.CreateTemp(os.TempDir(), "testansible.zip")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cases := []command.CommandTestCase{
		{
			Name: "ansibleBundle1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"label":                       "test",
				"title":                       "test",
				"type":                        stage1.StageDefinitionType,
				"ansible_bundle_filename":     f.Name,
				"http_request_body_from_pipe": true,
			}),
			Good: true,
			Id:   stage1.StageDefinitionID,
		},
		{
			Name: "missing label",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
	}

	command.TestCreateCommand(stageDefinitionCreateCmd, cases, client, t)
}

func TestStageDefinitionCreateHTTPRequestCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))
	req := metalcloud.HTTPRequest{
		URL: "http://asdad/asdasd/ass",
	}
	stage := metalcloud.StageDefinition{
		StageDefinitionID:    11,
		StageDefinitionLabel: "test2",
		StageDefinition:      req,
		StageDefinitionType:  "HTTPRequest",
	}

	client.EXPECT().
		StageDefinitionCreate(gomock.Any()).
		Return(&stage, nil).
		MinTimes(1)

	cases := []command.CommandTestCase{
		{
			Name: "httpRequest1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"label":               "test",
				"title":               "test",
				"type":                stage.StageDefinitionType,
				"http_request_url":    req.URL,
				"http_request_method": "get",
			}),
			Good: true,
			Id:   stage.StageDefinitionID,
		},
	}
	badTestCases := command.GenerateCommandTestCases(map[string]interface{}{"label": "test", "type": "HTTPRequest", "title": "test"})
	cases = append(cases, badTestCases...)

	command.TestCreateCommand(stageDefinitionCreateCmd, cases, client, t)
}
