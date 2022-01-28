package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

//driveArrayCmds commands affecting instance arrays
var driveSnapshotCmds = []Command{

	{
		Description:  "Creates a drive snapshot.",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("drive snapshots create", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id":  c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " The id of the drive to create a snapshot from"),
				"return_id": c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: driveSnapshotCreateCmd,
	},
	{
		Description:  "Lists drive snapshots.",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("drive snapshots list", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " The id of the drive for which to list snapshots."),
				"format":   c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveSnapshotListCmd,
	},
	{
		Description:  "Delete a snapshot.",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("drive snapshots delete", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_snapshot_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " The id of the drive snapshot"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: driveSnapshotDeleteCmd,
	},
	{
		Description:  "Rollback a snapshot.",
		Subject:      "drive-snapshot",
		AltSubject:   "snapshot",
		Predicate:    "rollback",
		AltPredicate: "rollback",
		FlagSet:      flag.NewFlagSet("drive snapshots rollback", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_snapshot_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " The id of the drive snapshot"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: driveSnapshotRollbackCmd,
	},
}

func driveSnapshotCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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

func driveSnapshotListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	driveID, ok := getIntParamOk(c.Arguments["drive_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "DRIVE_ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
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

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	subtitle := fmt.Sprintf("Snapshots of drive #%d", driveID)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Snapshots", subtitle, getStringParam(c.Arguments["format"]))
}

func driveSnapshotDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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

func driveSnapshotRollbackCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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
