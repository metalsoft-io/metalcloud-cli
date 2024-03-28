package drive

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
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

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd: command.MakeCommand(map[string]interface{}{
				"drive_id": 11,
			}),
			Good: true,
			Id:   s.DriveSnapshotID,
		},
	}

	command.TestCreateCommand(driveSnapshotCreateCmd, cases, client, t)
}

func TestSnapshotsListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	s := map[string]metalcloud.Snapshot{
		"test": {
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

	cases := []command.CommandTestCase{
		{
			Name: "good1",
			Cmd:  command.MakeCommand(map[string]interface{}{"drive_id": 10}),
			Good: true,
		},
		{
			Name: "no id",
			Cmd:  command.MakeCommand(map[string]interface{}{}),
			Good: false,
		},
	}

	command.TestGetCommand(driveSnapshotListCmd, cases, client, expectedFirstRow, t)

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

	cmd := command.MakeCommand(map[string]interface{}{"drive_snapshot_id": s.DriveSnapshotID})
	command.TestCommandWithConfirmation(driveSnapshotDeleteCmd, cmd, client, t)
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

	cmd := command.MakeCommand(map[string]interface{}{"drive_snapshot_id": s.DriveSnapshotID})
	command.TestCommandWithConfirmation(driveSnapshotRollbackCmd, cmd, client, t)
}
