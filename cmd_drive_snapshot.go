package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

//driveArrayCmds commands affecting instance arrays
var driveSnapshotCmds = []Command{

	Command{
		Description:  "Creates a drive snapshot.",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("drive snapshots create", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id":  c.FlagSet.Int("id", _nilDefaultInt, "(Required) The id of the drive to create a snapshot from"),
				"return_id": c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: driveSnapshotCreateCmd,
	},
	Command{
		Description:  "Lists drive snapshots",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("drive snapshots list", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) The id of the drive to create a snapshot from"),
				"format":   c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveSnapshotListCmd,
	},
	Command{
		Description:  "Delete snapshot",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("drive snapshots delete", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_snapshot_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) The id of the drive snapshot"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: driveSnapshotDeleteCmd,
	},
	Command{
		Description:  "Rollback snapshot",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "rollback",
		AltPredicate: "rollback",
		FlagSet:      flag.NewFlagSet("drive snapshots rollback", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_snapshot_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) The id of the drive snapshot"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: driveSnapshotRollbackCmd,
	},
}

func driveSnapshotCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	driveID, ok := getIntParamOk(c.Arguments["drive_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	ret, err := client.DriveSnapshotCreate(driveID)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.DriveSnapshotID), nil
	}

	return "", err
}

func driveSnapshotListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	driveID, ok := getIntParamOk(c.Arguments["drive_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 30,
		},
		SchemaField{
			FieldName: "DRIVE_ID",
			FieldType: TypeInt,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 10,
		},
	}

	snapshots, err := client.DriveSnapshots(driveID)
	if err != nil {
		return "", err
	}

	data := [][]interface{}{}
	for _, s := range *snapshots {

		data = append(data, []interface{}{
			s.DriveSnapshotID,
			s.DriveSnapshotLabel,
			s.DriveID,
			s.DriveSnapshotCreatedTimestamp,
		})

	}

	subtitle := fmt.Sprintf("Snapshots of drive #%d", driveID)
	return renderTable("Snapshots", subtitle, getStringParam(c.Arguments["format"]), data, schema)
}

func driveSnapshotDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	driveSnapshotID, ok := getIntParamOk(c.Arguments["drive_snapshot_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	snapshot, err := client.DriveSnapshotGet(driveSnapshotID)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Deleting snapshot %s (%d).  Are you sure? Type \"yes\" to continue:",
			snapshot.DriveSnapshotLabel,
			snapshot.DriveSnapshotID,
		)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage

	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.DriveSnapshotDelete(snapshot.DriveSnapshotID)
	}

	return "", err
}

func driveSnapshotRollbackCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	driveSnapshotID, ok := getIntParamOk(c.Arguments["drive_snapshot_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	snapshot, err := client.DriveSnapshotGet(driveSnapshotID)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Rolling back snapshot %s (%d) to date %s on drive %d.  Are you sure? Type \"yes\" to continue:",
			snapshot.DriveSnapshotLabel,
			snapshot.DriveSnapshotID,
			snapshot.DriveSnapshotCreatedTimestamp,
			snapshot.DriveID,
		)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage

	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.DriveSnapshotRollback(snapshot.DriveSnapshotID)
	}

	return "", err
}
