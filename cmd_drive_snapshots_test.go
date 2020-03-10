package main

import (
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestCreateDriveSnapshotCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	s := metalcloud.Snapshot{
		DriveSnapshotID: 100,
	}

	client.EXPECT().
		DriveSnapshotCreate(gomock.Any()).
		Return(&s, nil).
		MinTimes(1)

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"drive_id": 11,
			}),
			good: true,
			id:   s.DriveSnapshotID,
		},
	}

	testCreateCommand(driveSnapshotCreateCmd, cases, client, t)
}

func TestSnapshotsListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	s := map[string]metalcloud.Snapshot{
		"test": metalcloud.Snapshot{
			DriveSnapshotID: 100,
		},
	}

	client.EXPECT().
		DriveSnapshots(gomock.Any()).
		Return(&s, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"ID": 100,
	}

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd:  MakeCommand(map[string]interface{}{"drive_id": 10}),
			good: true,
		},
		{
			name: "no id",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
		},
	}

	testGetCommand(driveSnapshotListCmd, cases, client, expectedFirstRow, t)

}

func TestDeleteSnapshotCmd(t *testing.T) {
	RegisterTestingT(t)
	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	s := metalcloud.Snapshot{
		DriveSnapshotID: 100,
	}

	client.EXPECT().
		DriveSnapshotGet(s.DriveSnapshotID).
		Return(&s, nil).
		MinTimes(1)

	client.EXPECT().
		DriveSnapshotDelete(s.DriveSnapshotID).
		Return(nil).
		MinTimes(1)

	cmd := MakeCommand(map[string]interface{}{"drive_snapshot_id": s.DriveSnapshotID})
	testCommandWithConfirmation(driveSnapshotDeleteCmd, cmd, client, t)
}

func TestDriveSnapshotRollbackCmd(t *testing.T) {
	RegisterTestingT(t)
	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	s := metalcloud.Snapshot{
		DriveSnapshotID: 100,
	}

	client.EXPECT().
		DriveSnapshotGet(s.DriveSnapshotID).
		Return(&s, nil).
		MinTimes(1)

	client.EXPECT().
		DriveSnapshotRollback(s.DriveSnapshotID).
		Return(nil).
		MinTimes(1)

	cmd := MakeCommand(map[string]interface{}{"drive_snapshot_id": s.DriveSnapshotID})
	testCommandWithConfirmation(driveSnapshotRollbackCmd, cmd, client, t)
}
