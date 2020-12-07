package main

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const applyTestCasesDir = "./cmd_apply_test_cases/apply/"
const deleteTestCasesDir = "./cmd_apply_test_cases/delete/"

func TestApply(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	cases := []CommandTestCase{
		{
			name: "missing file name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
		{
			name: "missing file/not a file",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": applyTestCasesDir,
			}),
			good: false,
		},
	}

	files, err := ioutil.ReadDir(applyTestCasesDir)
	if err != nil {
		t.Error(err)
	}

	for _, f := range files {
		testCase := CommandTestCase{
			name: fmt.Sprintf("apply good %s", f.Name()),
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": fmt.Sprintf("%s%s", applyTestCasesDir, f.Name()),
			}),
			good: true,
			id:   0,
		}
		cases = append(cases, testCase)
	}

	testCreateCommand(applyCmd, cases, client, t)
}

func TestDelete(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	cases := []CommandTestCase{
		{
			name: "missing file name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
		{
			name: "missing file/not a file",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": deleteTestCasesDir,
			}),
			good: false,
		},
	}

	files, err := ioutil.ReadDir(deleteTestCasesDir)
	if err != nil {
		t.Error(err)
	}

	for _, f := range files {
		testCase := CommandTestCase{
			name: fmt.Sprintf("delete good %s", f.Name()),
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": fmt.Sprintf("%s%s", deleteTestCasesDir, f.Name()),
			}),
			good: true,
			id:   0,
		}
		cases = append(cases, testCase)
	}

	testCreateCommand(deleteCmd, cases, client, t)
}

func TestReadObjectsFromCommand(t *testing.T) {
	RegisterTestingT(t)

	for _, c := range readFromFileTestCases {
		f, err := ioutil.TempFile("./", "testread-*.yaml")
		if err != nil {
			t.Error(err)
		}

		f.WriteString(c.content)
		f.Close()
		defer syscall.Unlink(f.Name())

		cmd := MakeCommand(map[string]interface{}{
			"read_config_from_file": f.Name(),
		})
		objects, err := readObjectsFromCommand(&cmd)
		Expect(err).To(BeNil())

		expected := c.objects

		for index, object := range expected {
			Expect(object).To(Equal(objects[index]))
		}
	}
}

type ApplyTestCase struct {
	content string
	objects []metalcloud.Applier
}

var readFromFileTestCases = []ApplyTestCase{
	{
		content: "kind: SharedDrive\napiVersion: 1.0\ninfrastructureID: 2\nlabel: test-shared\nstorageType: iscsi_ssd\n",
		objects: []metalcloud.Applier{
			metalcloud.SharedDrive{
				InfrastructureID:       2,
				SharedDriveLabel:       "test-shared",
				SharedDriveStorageType: "iscsi_ssd",
			},
		},
	},
	{
		content: "kind: InstanceArray\napiVersion: 1.0\ninfrastructureID: 2\nlabel: test-ia\ninstanceCount: 1\n---\nkind: DriveArray\napiVersion: 1.0\ninfrastructureID: 2\nlabel: test-da\ncount: 2",
		objects: []metalcloud.Applier{
			metalcloud.InstanceArray{
				InfrastructureID:           2,
				InstanceArrayLabel:         "test-ia",
				InstanceArrayInstanceCount: 1,
			},
			metalcloud.DriveArray{
				InfrastructureID: 2,
				DriveArrayLabel:  "test-da",
				DriveArrayCount:  2,
			},
		},
	},
}
