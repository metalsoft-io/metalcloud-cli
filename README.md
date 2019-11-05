# metalcloud-cli

[![Build Status](https://travis-ci.org/bigstepinc/metalcloud-cli.svg?branch=master)](https://travis-ci.org/bigstepinc/metalcloud-cli)

This tool allows the manipulation of all Bigstep Metal Cloud elements via the command line.

To install use:
```bash
go install github.com/bigstepinc/metalcloud-cli
```

Configure credentials as environment variables:
```bash
export METALCLOUD_API_KEY="<your key>"
export METALCLOUD_ENDPOINT="https://api.bigstep.com/metal-cloud"
export METALCLOUD_USER_EMAIL="<your email>"
export METALCLOUD_DATACENTER="uk-reading"
```

## Getting started

To create an infrastructure, in the default datacenter, configured via the `METALCLOUD_DATACENTER` environment variable):
```
metalcloud-cli create infrastructure -label test -return_id
```

```
metalcloud-cli list infrastructure 
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
| ID    | LABEL                                   | OWNER                         | REL.      | STATUS    | CREATED             | UPDATED             |
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
| 12345 | complex-demo                            | d.d@sdd.com                   | OWNER     | active    | 2019-03-28T15:23:08Z| 2019-03-28T15:23:08Z|
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
```

To create an instance array in that infrastructure, get the ID of the infrastructure from above (12345):

```
metalcloud-cli create instance_array -infra 12345 -label master -proc 1 -proc_core_count 8 -ram 16
```

To view the id of the previously created drive array use
```
metalcloud-cli list instance_array -infra 12345
+-------+---------------------+---------------------+-----------+
| ID    | LABEL               | STATUS              | INST_CNT  |
+-------+---------------------+---------------------+-----------+
| 54321 | master              | ordered             | 1         |
+-------+---------------------+---------------------+-----------+
Total: 1 Instance Arrays
```

To create a drive array and attach it to the previous instance array
```
metalcloud-cli create drive_array -infra 12345 -label master-da -ia 54321
```

To view the current status of the infrastructure
```
metalcloud-cli get infrastructure -infra 12345
alex@Alexandrus-MacBook-Pro metalcloud-cli $ ./metalcloud-cli get infrastructure -id 26345
Infrastructures I have access to (as test@test.com)
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
| ID    | OBJECT_TYPE    | LABEL                         | DETAILS                                                               | STATUS    |
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
| 36791 | InstanceArray  | master                        | 1 instances (16 RAM, 8 cores, 1 disks)                                | ordered   |
| 47398 | DriveArray     | master-da                     | 1 drives - 40.0 GB iscsi_ssd (volume_template:0) attached to: 36791   | ordered   |
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
Total: 2 elements
```

## Condensed format

The CLI also provides a "condensed format" for most of it's commands:
* instance_array = ia
* drive_array = da
* infrastructure = infra
* list = ls
* delete = rm
...

This allows commands such as:
```
metalcloud-cli ls infra
```



## Supported commands

This is the output of the `metalcloud-cli help` command:

```
Syntax: ./metalcloud-cli <command> [args]
Accepted commands:
create infrastructure - Creates an infrastructure. (alternatively use "new infra")
  -dc string
        (Required) Infrastructure datacenter (default "uk-reading")
  -label string
        (Required) Infrastructure's label
  -return_id
        (Optional) Will print the ID of the created infrastructure. Useful for automating tasks.

list infrastructure - Lists all infrastructures. (alternatively use "ls infra")
  -format string
        The output format. Supproted values are 'json','csv'. The default format is human readable.

delete infrastructure - Delete an infrastructure. (alternatively use "rm infra")
  -autoconfirm
        If true it does not ask for confirmation anymore
  -id int
        (Required) Infrastructure's id

deploy infrastructure - Deploy an infrastructure. (alternatively use "apply infra")
  -allow_data_loss
        (Optional, default false) If false, deploy will throw error if data loss is expected.
  -attempt_soft_shutdown
        (Optional, default true) If needed, atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy (default true)
  -autoconfirm
        If true it does not ask for confirmation anymore
  -hard_shutdown_after_timeout
        (Optional, default true) Force a hard power off after timeout expired and the server is not powered off. (default true)
  -id int
        (Required) Infrastructure's id
  -skip_ansible
        (Optional, default false) If true, some automatic provisioning steps will be skipped. This parameter should generally be ignored.
  -soft_shutdown_timeout_seconds int
        (Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set. (default 180)

get infrastructure - Get an infrastructure. (alternatively use "show infra")
  -format string
        The output format. Supproted values are 'json','csv'. The default format is human readable.
  -id int
        (Required) Infrastructure's id

revert infrastructure - Revert all changes of an infrastructure. (alternatively use "undo infra")
  -autoconfirm
        If true it does not ask for confirmation anymore
  -id int
        (Required) Infrastructure's id

create instance_array - Creates an instance array. (alternatively use "new ia")
  -boot string
        InstanceArray's boot type:'pxe_iscsi','local_drives' (default "__NIL__")
  -disk_size int
        InstanceArray's local disk sizes (default -14234)
  -disks int
        InstanceArray's number of local drives (default -14234)
  -infra int
        (Required) Infrastrucure ID (default -14234)
  -instance_count int
        (Required) Instance count of this instance array (default -14234)
  -label string
        InstanceArray's label (default "__NIL__")
  -managed_fw
        InstanceArray's firewall management on or off (default true)
  -proc int
        InstanceArray's minimum processor count (default -14234)
  -proc_core_count int
        InstanceArray's minimum processor core count (default -14234)
  -proc_freq int
        InstanceArray's minimum processor frequency (Mhz) (default -14234)
  -ram int
        InstanceArray's minimum RAM (GB) (default -14234)
  -return_id
        (Optional) Will print the ID of the created Instance Array. Useful for automating tasks.
  -template int
        InstanceArray's volume template when booting from for local drives (default -14234)

list instance_array - Lists all instance arrays of an infrastructure. (alternatively use "ls ia")
  -format string
        The output format. Supproted values are 'json','csv'. The default format is human readable.
  -infra int
        (Required) Infrastrucure ID

delete instance_array - Delete instance array. (alternatively use "rm ia")
  -autoconfirm
        If true it does not ask for confirmation anymore
  -id int
        (Required) InstanceArray ID

edit instance_array - Edits an instance array. (alternatively use "alter ia")
  -boot string
        InstanceArray's boot type:'pxe_iscsi','local_drives' (default "__NIL__")
  -disk_size int
        InstanceArray's local disk sizes (default -14234)
  -disks int
        InstanceArray's number of local drives (default -14234)
  -id int
        (Required) InstanceArray's id (default -14234)
  -instance_count int
        Instance count of this instance array (default -14234)
  -keep_detaching_drives
        If false and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted. (default true)
  -label string
        (Required) InstanceArray's label (default "__NIL__")
  -managed_fw
        InstanceArray's firewall management on or off (default true)
  -proc int
        InstanceArray's minimum processor count (default -14234)
  -proc_core_count int
        InstanceArray's minimum processor core count (default -14234)
  -proc_freq int
        InstanceArray's minimum processor frequency (Mhz) (default -14234)
  -ram int
        InstanceArray's minimum RAM (GB) (default -14234)
  -swap_existing_hardware
        If true, all the hardware of the Instance objects is swapped to match the new InstanceArray specifications
  -template int
        InstanceArray's volume template when booting from for local drives (default -14234)

create drive_array - Creates a drive array. (alternatively use "new da")
  -count int
        DriveArrays's drive count. Use this only for unconnected DriveArrays. (default -14234)
  -expand_with_ia
        Auto-expand when the connected instance array expands (default true)
  -ia int
        (Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays (default -14234)
  -infra int
        (Required) Infrastrucure ID (default -14234)
  -label string
        (Required) The label of the drive array (default "__NIL__")
  -return_id
        (Optional) Will print the ID of the created Drive Array. Useful for automating tasks.
  -size int
        (Optional, default = 40960) Drive arrays's size in MBytes (default -14234)
  -template int
        DriveArrays's volume template to clone when creating Drives (default -14234)
  -type string
        Possible values: iscsi_ssd, iscsi_hdd (default "__NIL__")

edit drive_array - Edit a drive array. (alternatively use "alter da")
  -count int
        DriveArrays's drive count. Use this only for unconnected DriveArrays. (default -14234)
  -expand_with_ia
        Auto-expand when the connected instance array expands (default true)
  -ia int
        (Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays (default -14234)
  -id int
        (Required) Drive Array's ID (default -14234)
  -label string
        (Required) The label of the drive array (default "__NIL__")
  -size int
        (Optional, default = 40960) Drive arrays's size in MBytes (default -14234)
  -template int
        DriveArrays's volume template to clone when creating Drives (default -14234)
  -type string
        Possible values: iscsi_ssd, iscsi_hdd (default "__NIL__")

list drive_array - Lists all drive arrays of an infrastructure. (alternatively use "ls da")
  -format string
        The output format. Supproted values are 'json','csv'. The default format is human readable.
  -infra int
        (Required) Infrastrucure ID

delete drive_array - Delete a drive array. (alternatively use "rm da")
  -autoconfirm
        If true it does not ask for confirmation anymore
  -id int
        (Required) DriveArray's ID
```