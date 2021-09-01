package main

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"

	. "github.com/onsi/gomega"
)

const deleteTestCasesDir = "./cmd_apply_test_cases/delete/"

func TestApply(t *testing.T) {
	RegisterTestingT(t)
	// dcBytes, err := yaml.Marshal(_osTemplate1)

	// Expect(err).To(BeNil())

	// fmt.Printf("yaml is %s\n", string(dcBytes))

	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)
	assetList := map[string]metalcloud.OSAsset{
		"1": _osAsset1,
	}

	var err error

	client.EXPECT().
		SharedDriveGet(gomock.Any()).
		Return(&_sharedDrive1, nil).
		AnyTimes()

	client.EXPECT().
		SharedDriveEdit(gomock.Any(), gomock.Any()).
		Return(&_sharedDrive1, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(gomock.Any()).
		Return(&_instanceArray1, nil).
		AnyTimes()
	client.EXPECT().
		InstanceArrayEditByLabel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&_instanceArray1, nil).
		AnyTimes()

	client.EXPECT().
		DriveArrayGetByLabel(gomock.Any()).
		Return(&_driveArray1, nil).
		AnyTimes()
	client.EXPECT().
		DriveArrayEdit(gomock.Any(), gomock.Any()).
		Return(&_driveArray1, nil).
		AnyTimes()

	client.EXPECT().
		DatacenterGet(gomock.Any()).
		Return(&_datacenter1, nil).
		AnyTimes()
	client.EXPECT().
		DatacenterConfigUpdate(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(gomock.Any()).
		Return(&_infrastructure1, nil).
		AnyTimes()
	client.EXPECT().
		InfrastructureEdit(gomock.Any(), gomock.Any()).
		Return(&_infrastructure1, nil).
		AnyTimes()
	client.EXPECT().
		SecretGet(gomock.Any()).
		Return(&_secret1, nil).
		AnyTimes()
	client.EXPECT().
		SecretUpdate(gomock.Any(), gomock.Any()).
		Return(&_secret1, nil).
		AnyTimes()
	client.EXPECT().
		NetworkGet(gomock.Any()).
		Return(&_network1, nil).
		AnyTimes()
	client.EXPECT().
		NetworkEdit(gomock.Any(), gomock.Any()).
		Return(&_network1, nil).
		AnyTimes()
	client.EXPECT().
		OSTemplateGet(gomock.Any(), gomock.Any()).
		Return(&_osTemplate1, nil).
		AnyTimes()
	client.EXPECT().
		OSTemplateUpdate(gomock.Any(), gomock.Any()).
		Return(&_osTemplate1, nil).
		AnyTimes()
	client.EXPECT().
		OSAssetGet(gomock.Any()).
		Return(&_osAsset1, nil).
		AnyTimes()
	client.EXPECT().
		OSAssets().
		Return(&assetList, nil).
		AnyTimes()
	client.EXPECT().
		OSAssetUpdate(gomock.Any(), gomock.Any()).
		Return(&_osAsset1, nil).
		AnyTimes()
	client.EXPECT().
		WorkflowGet(gomock.Any()).
		Return(&_workflow1, nil).
		AnyTimes()
	client.EXPECT().
		WorkflowUpdate(gomock.Any(), gomock.Any()).
		Return(&_workflow1, nil).
		AnyTimes()
	client.EXPECT().
		VariableGet(gomock.Any()).
		Return(&_variable1, nil).
		AnyTimes()
	client.EXPECT().
		VariableUpdate(gomock.Any(), gomock.Any()).
		Return(&_variable1, nil).
		AnyTimes()
	client.EXPECT().
		ServerGet(gomock.Any(), false).
		Return(&_server1, nil).
		AnyTimes()
	client.EXPECT().
		ServerEditComplete(gomock.Any(), gomock.Any()).
		Return(&_server1, nil).
		AnyTimes()
	client.EXPECT().
		StageDefinitionGet(gomock.Any()).
		Return(&_stageDefinition1, nil).
		AnyTimes()
	client.EXPECT().
		StageDefinitionUpdate(gomock.Any(), gomock.Any()).
		Return(&_stageDefinition1, nil).
		AnyTimes()
	client.EXPECT().
		SwitchDeviceGet(gomock.Any(), false).
		Return(&_switchDevice1, nil).
		AnyTimes()
	client.EXPECT().
		SwitchDeviceUpdate(gomock.Any(), _switchDevice1, gomock.Any()).
		Return(&_switchDevice1, nil).
		AnyTimes()
	client.EXPECT().
		SubnetPoolGet(gomock.Any()).
		Return(&_subnetPool1, nil).
		AnyTimes()
	client.EXPECT().
		SubnetPoolGet(gomock.Any()).
		Return(nil, err).
		AnyTimes()
	client.EXPECT().
		SubnetPoolCreate(_subnetPool2).
		Return(&_subnetPool2, nil).
		AnyTimes()
	cases := []CommandTestCase{
		{
			name: "missing file name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
		{
			name: "missing file/not a file",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": "./examples",
			}),
			good: false,
		},
	}

	for _, c := range applyTestCases {
		f, err := ioutil.TempFile("./", "testapply-*.yaml")
		if err != nil {
			t.Error(err)
		}

		f.WriteString(c)
		f.Close()
		defer syscall.Unlink(f.Name())

		testCase := CommandTestCase{
			name: fmt.Sprintf("apply good %s", f.Name()),
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": f.Name(),
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

	client.EXPECT().
		SharedDriveDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		SharedDriveGetByLabel(gomock.Any()).
		Return(&_sharedDrive1, nil).
		AnyTimes()
	client.EXPECT().
		InstanceArrayDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		InstanceArrayGetByLabel("ia-test").
		Return(&_instanceArray1, nil).
		AnyTimes()
	client.EXPECT().
		DriveArrayDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		DriveArrayGetByLabel(gomock.Any()).
		Return(&_driveArray1, nil).
		AnyTimes()
	client.EXPECT().
		InfrastructureDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		InfrastructureGetByLabel(gomock.Any()).
		Return(&_infrastructure1, nil).
		AnyTimes()
	client.EXPECT().
		SecretDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		NetworkDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		NetworkGetByLabel(gomock.Any()).
		Return(&_network1, nil).
		AnyTimes()
	client.EXPECT().
		OSTemplateDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		OSAssetDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		WorkflowDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		VariableDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		ServerDelete(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		StageDefinitionDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		SwitchDeviceDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	client.EXPECT().
		SubnetPoolDelete(gomock.Any()).
		Return(nil).
		AnyTimes()
	cases := []CommandTestCase{
		{
			name: "missing file name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
		{
			name: "missing file/not a file",
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": "./examples",
			}),
			good: false,
		},
	}

	for _, c := range applyTestCases {
		f, err := ioutil.TempFile("./", "testdelete-*.yaml")
		if err != nil {
			t.Error(err)
		}

		f.WriteString(c)
		f.Close()
		defer syscall.Unlink(f.Name())

		testCase := CommandTestCase{
			name: fmt.Sprintf("delete good %s", f.Name()),
			cmd: MakeCommand(map[string]interface{}{
				"read_config_from_file": f.Name(),
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
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

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

		objects, err := readObjectsFromCommand(&cmd, client)

		Expect(err).To(BeNil())

		for index, object := range c.objects {
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
		content: _sharedDriveFixtureYaml1,
		objects: []metalcloud.Applier{
			_sharedDrive1,
		},
	},
	{
		content: fmt.Sprintf("%s\n%s\n%s", _instanceArrayFixtureYaml1, yamlSeparator, _driveArrayFixtureYaml1),
		objects: []metalcloud.Applier{
			_instanceArray1,
			_driveArray1,
		},
	},
	{
		content: _datacenterFixtureYaml1,
		objects: []metalcloud.Applier{
			_datacenter1,
		},
	},
	{
		content: _infrastructureFixtureYaml1,
		objects: []metalcloud.Applier{
			_infrastructure1,
		},
	},
	{
		content: _networkFixtureYaml1,
		objects: []metalcloud.Applier{
			_network1,
		},
	},
	{
		content: _osAssetFixtureYaml1,
		objects: []metalcloud.Applier{
			_osAsset1,
		},
	},
	{
		content: _osTemplateFixtureYaml1,
		objects: []metalcloud.Applier{
			_osTemplate1,
		},
	},
	{
		content: _secretFixtureYaml1,
		objects: []metalcloud.Applier{
			_secret1,
		},
	},
	{
		content: _workflowFixtureYaml1,
		objects: []metalcloud.Applier{
			_workflow1,
		},
	},
	{
		content: _variableFixtureYaml1,
		objects: []metalcloud.Applier{
			_variable1,
		},
	},
	{
		content: _serverFixtureYaml1,
		objects: []metalcloud.Applier{
			_server1,
		},
	},
	{
		content: _stageDefinitionFixtureYaml1,
		objects: []metalcloud.Applier{
			_stageDefinition1,
		},
	},
	{
		content: _switchDeviceFixtureYaml1,
		objects: []metalcloud.Applier{
			_switchDevice1,
		},
	},
	{
		content: _subnetPoolFixtureYaml1,
		objects: []metalcloud.Applier{
			_subnetPool1,
		},
	},
}

var applyTestCases = []string{
	_instanceArrayFixtureYaml1,
	_driveArrayFixtureYaml1,
	_datacenterFixtureYaml1,
	_sharedDriveFixtureYaml1,
	_infrastructureFixtureYaml1,
	_networkFixtureYaml1,
	_osAssetFixtureYaml1,
	_osTemplateFixtureYaml1,
	_secretFixtureYaml1,
	_workflowFixtureYaml1,
	_variableFixtureYaml1,
	_serverFixtureYaml1,
	_stageDefinitionFixtureYaml1,
	_switchDeviceFixtureYaml1,
	_subnetPoolFixtureYaml1,
	_subnetPoolFixtureYaml2,
}

const _datacenterFixtureYaml1 = "kind: Datacenter\napiVersion: 1.0\nuserid: 1\nname: dctest\nconfig:\n    BSIMachinesSubnetIPv4CIDR: 10.255.226.0/24\n    BSIVRRPListenIPv4: 172.16.10.6\n    BSIMachineListenIPv4List:\n        - 172.16.10.6\n    BSIExternallyVisibleIPv4: 89.36.24.2\n    repoURLRoot: https://repointegrationpublic.bigstepcloud.com\n    repoURLRootQuarantineNetwork: https://repointegrationpublic.bigstepcloud.com\n    SANRoutedSubnet: 100.64.0.0/21\n    NTPServers:\n        - 84.40.58.44\n        - 84.40.58.45\n    DNSServers:\n        - 84.40.63.27\n    TFTPServerWANVRRPListenIPv4: 172.16.10.6\n    dataLakeEnabled: false\n    serverRegisterUsingGeneratedIPMICredentialsEnabled: false\n    datacenterNetworkIsLayer2Only: false\n    enableTenantAccessToIPMI: false\n    proxyURL: \"\"\n    proxyUsername: \"\"\n    proxyPassword: \"\"\n    enableProxyURL: false"
const _instanceArrayFixtureYaml1 = "kind: InstanceArray\napiVersion: 1.0\ninstanceID: 100\ninfrastructureID: 2\nlabel: ia-test\noperation:\n  id: 100\n  label: ia-test\n  changeID: 200"
const _driveArrayFixtureYaml1 = "kind: DriveArray\napiVersion: 1.0\ninfrastructureID: 25524\nlabel: drive-array-45928\ncount: 2\nvolumeTemplateID: 78\nserviceStatus: active\nstorageType: iscsi_ssd\noperation:\n  \n  label: drive-array-45928\n  count: 2\n  volumeTemplateID: 78\n  storageType: iscsi_ssd\n  changeID: 215701\n  id: 45928\n  sizeMBytes: 40960\n  instanceArrayID: 35516\n  expandWithInstanceArray: true"
const _sharedDriveFixtureYaml1 = "kind: SharedDrive\napiVersion: 1.0\nid: 100\ninfrastructureID: 1\nlabel: shared-drive-test\nstorageType: iscsi_ssd\nhasGFS: false\nsizeMBytes: 2048\nsubdomain: csivolumename.test-kube-csi.7.bigstep.io\nattachedInstaceArrays:\n  - 37824\noperation:\n  infrastructureID: 1\n  label: shared-drive-test\n  attachedInstanceArrays:\n  - 37824\n  storageType: iscsi_ssd\n  hasGFS: false\n  sizeMBytes: 2048\n  subdomain: csivolumename.test-kube-csi.7.bigstep.io\n  changeID: 16508\n  id: 100\n"
const _networkFixtureYaml1 = "kind: Network\napiVersion: 1.0\nid: 101\nlabel: net-test\nsubdomain: sub-test.test\ntype: test-net-type\ninfrastructureID: 1\noperation:\n    id: 101\n    label: net-test\n    infrastructureID: 1\n    changeID: 3"
const _osAssetFixtureYaml1 = "kind: OSAsset\napiVersion: 1.0\nid: 100\nownerID: 1\nfileName: os-test\nfileMime: testMime\ncontentBase64: content\nusage: testUsage"
const _osTemplateFixtureYaml1 = "kind: OSTemplate\napiVersion: 1.0\nid: 100\nname: test-display-template\nbootType: test-boot\nos:\n    type: os-type\n    version: os-version\n    architecture: os-arch"
const _secretFixtureYaml1 = "kind: Secret\napiVersion: 1.0\nid: 1\nname: secret-test"
const _infrastructureFixtureYaml1 = "kind: Infrastructure\napiVersion: 1.0\nid: 4103\nlabel: demo\ndatacenter: us-santaclara\nsubdomain: demo.2.poc.metalcloud.io\nownerID: 2\ntouchUnixTime: \"1573829237.9229\"\nserviceStatus: active\ncreatedTimestamp: \"2019-11-12T20:44:04Z\"\nupdatedTimestamp: \"2019-11-12T20:44:04Z\"\nchangeID: 8805\ndeployID: 10420\noperation:\n    id: 4103\n    label: demo\n    datacenter: us-santaclara\n    deployStatus: finished\n    deployType: create\n    subdomain: demo.2.poc.metalcloud.io\n    ownerID: 2\n    updatedTimestamp: \"2019-11-12T20:44:04Z\"\n    changeID: 8805"
const _workflowFixtureYaml1 = "kind: Workflow\napiVersion: 1.0\nid: 100\nusage: test-usage"
const _variableFixtureYaml1 = "kind: Variable\napiVersion: 1.0\nid: 100"
const _serverFixtureYaml1 = "kind: Server\napiVersion: 1.0\nid: 100"
const _stageDefinitionFixtureYaml1 = "kind: StageDefinition\napiVersion: 1.0\nid: 100\ntype: sd-test-type\ntitle: sd-test-title"
const _switchDeviceFixtureYaml1 = "kind: SwitchDevice\napiVersion: 1.0\nid: 100"
const _subnetPoolFixtureYaml1 = "kind: SubnetPool\napiVersion: 1.0\nid: 100"
const _subnetPoolFixtureYaml2 = "kind: SubnetPool\napiVersion: 1.0\nid: 101"

var _subnetPool1 = metalcloud.SubnetPool{
	SubnetPoolID: 100,
}

var _subnetPool2 = metalcloud.SubnetPool{
	SubnetPoolID: 101,
}

var _switchDevice1 = metalcloud.SwitchDevice{
	NetworkEquipmentID: 100,
}

var _stageDefinition1 = metalcloud.StageDefinition{
	StageDefinitionID:    100,
	StageDefinitionType:  "sd-test-type",
	StageDefinitionTitle: "sd-test-title",
}

var _server1 = metalcloud.Server{
	ServerID: 100,
}

var _variable1 = metalcloud.Variable{
	VariableID: 100,
}

var _workflow1 = metalcloud.Workflow{
	WorkflowID:    100,
	WorkflowUsage: "test-usage",
}

var _secret1 = metalcloud.Secret{
	SecretID:   1,
	SecretName: "secret-test",
}

var _osTemplate1 = metalcloud.OSTemplate{
	VolumeTemplateID:          100,
	VolumeTemplateDisplayName: "test-display-template",
	VolumeTemplateBootType:    "test-boot",
	VolumeTemplateOperatingSystem: &metalcloud.OperatingSystem{
		OperatingSystemType:         "os-type",
		OperatingSystemVersion:      "os-version",
		OperatingSystemArchitecture: "os-arch",
	},
}

var _osAsset1 = metalcloud.OSAsset{
	OSAssetID:             100,
	UserIDOwner:           1,
	OSAssetFileName:       "os-test",
	OSAssetFileMime:       "testMime",
	OSAssetContentsBase64: "content",
	OSAssetUsage:          "testUsage",
}
var _network1 = metalcloud.Network{
	NetworkID:                 101,
	NetworkLabel:              "net-test",
	InfrastructureID:          1,
	NetworkSubdomain:          "sub-test.test",
	NetworkType:               "test-net-type",
	NetworkLANAutoAllocateIPs: false,
	NetworkOperation: &metalcloud.NetworkOperation{
		NetworkID:        101,
		NetworkLabel:     "net-test",
		InfrastructureID: 1,
		NetworkChangeID:  3,
	},
}
var _datacenter1 = metalcloud.Datacenter{
	DatacenterName: "dctest",
	UserID:         1,
	DatacenterConfig: &metalcloud.DatacenterConfig{
		SANRoutedSubnet:                       "100.64.0.0/21",
		BSIVRRPListenIPv4:                     "172.16.10.6",
		BSIMachineListenIPv4List:              []string{"172.16.10.6"},
		BSIMachinesSubnetIPv4CIDR:             "10.255.226.0/24",
		BSIExternallyVisibleIPv4:              "89.36.24.2",
		RepoURLRoot:                           "https://repointegrationpublic.bigstepcloud.com",
		RepoURLRootQuarantineNetwork:          "https://repointegrationpublic.bigstepcloud.com",
		DNSServers:                            []string{"84.40.63.27"},
		NTPServers:                            []string{"84.40.58.44", "84.40.58.45"},
		KMS:                                   "",
		TFTPServerWANVRRPListenIPv4:           "172.16.10.6",
		DataLakeEnabled:                       false,
		MonitoringGraphitePlainTextSocketHost: "",
		MonitoringGraphiteRenderURLHost:       "",
		Latitude:                              0,
		Longitude:                             0,
		EnableTenantAccessToIPMI:              false,
	},
}

var _infrastructure1 = metalcloud.Infrastructure{
	InfrastructureID:               4103,
	DatacenterName:                 "us-santaclara",
	UserIDowner:                    2,
	InfrastructureLabel:            "demo",
	InfrastructureCreatedTimestamp: "2019-11-12T20:44:04Z",
	InfrastructureSubdomain:        "demo.2.poc.metalcloud.io",
	InfrastructureChangeID:         8805,
	InfrastructureServiceStatus:    "active",
	InfrastructureTouchUnixtime:    "1573829237.9229",
	InfrastructureUpdatedTimestamp: "2019-11-12T20:44:04Z",
	InfrastructureDeployID:         10420,
	InfrastructureDesignIsLocked:   false,
	InfrastructureOperation: metalcloud.InfrastructureOperation{
		InfrastructureChangeID:         8805,
		InfrastructureID:               4103,
		DatacenterName:                 "us-santaclara",
		UserIDOwner:                    2,
		InfrastructureLabel:            "demo",
		InfrastructureSubdomain:        "demo.2.poc.metalcloud.io",
		InfrastructureDeployType:       "create",
		InfrastructureDeployStatus:     "finished",
		InfrastructureUpdatedTimestamp: "2019-11-12T20:44:04Z",
	},
}
var _driveArray1 = metalcloud.DriveArray{
	VolumeTemplateID:        78,
	DriveArrayStorageType:   "iscsi_ssd",
	InfrastructureID:        25524,
	DriveArrayServiceStatus: "active",
	DriveArrayCount:         2,
	DriveArrayLabel:         "drive-array-45928",
	DriveArrayOperation: &metalcloud.DriveArrayOperation{
		DriveArrayID:                      45928,
		DriveArrayChangeID:                215701,
		VolumeTemplateID:                  78,
		DriveArrayLabel:                   "drive-array-45928",
		DriveArrayStorageType:             "iscsi_ssd",
		DriveArrayCount:                   2,
		DriveSizeMBytesDefault:            40960,
		InstanceArrayID:                   35516,
		DriveArrayExpandWithInstanceArray: true,
	},
}

var _instanceArray1 = metalcloud.InstanceArray{
	InstanceArrayID:    100,
	InstanceArrayLabel: "ia-test",
	InfrastructureID:   2,
	InstanceArrayOperation: &metalcloud.InstanceArrayOperation{
		InstanceArrayID:       100,
		InstanceArrayLabel:    "ia-test",
		InstanceArrayChangeID: 200,
	},
}

var _sharedDrive1 = metalcloud.SharedDrive{
	InfrastructureID:       1,
	SharedDriveID:          100,
	SharedDriveLabel:       "shared-drive-test",
	SharedDriveSizeMbytes:  2048,
	SharedDriveSubdomain:   "csivolumename.test-kube-csi.7.bigstep.io",
	SharedDriveHasGFS:      false,
	SharedDriveStorageType: "iscsi_ssd",
	SharedDriveAttachedInstanceArrays: []int{
		37824,
	},
	SharedDriveOperation: metalcloud.SharedDriveOperation{
		InfrastructureID:       1,
		SharedDriveID:          100,
		SharedDriveLabel:       "shared-drive-test",
		SharedDriveSizeMbytes:  2048,
		SharedDriveSubdomain:   "csivolumename.test-kube-csi.7.bigstep.io",
		SharedDriveHasGFS:      false,
		SharedDriveStorageType: "iscsi_ssd",
		SharedDriveAttachedInstanceArrays: []int{
			37824,
		},
		SharedDriveChangeID: 16508,
	},
}
