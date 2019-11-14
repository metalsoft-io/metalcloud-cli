# metalcloud-cli

[![Build Status](https://travis-ci.org/bigstepinc/metalcloud-cli.svg?branch=master)](https://travis-ci.org/bigstepinc/metalcloud-cli)
[![Coverage Status](https://coveralls.io/repos/github/bigstepinc/metalcloud-cli/badge.svg?branch=master)](https://coveralls.io/github/bigstepinc/metalcloud-cli?branch=master)

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

To view the id of the previously created drive array:

```
metalcloud-cli list instance_array -infra 12345
+-------+---------------------+---------------------+-----------+
| ID    | LABEL               | STATUS              | INST_CNT  |
+-------+---------------------+---------------------+-----------+
| 54321 | master              | ordered             | 1         |
+-------+---------------------+---------------------+-----------+
Total: 1 Instance Arrays
```

To create a drive array and attach it to the previous instance array:

```
metalcloud-cli create drive_array -infra 12345 -label master-da -ia 54321
```

To view the current status of the infrastructure

```
metalcloud-cli get infrastructure -id 12345
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
Command: create infrastructure     Creates an infrastructure. (alternatively use "new infra")
	  -dc                        (Required) Infrastructure datacenter
	  -label                     (Required) Infrastructure's label
	  -return_id                 (Flag) If set will print the ID of the created infrastructure. Useful for automating tasks.

Command: list infrastructure       Lists all infrastructures. (alternatively use "ls infra")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.

Command: delete infrastructure     Delete an infrastructure. (alternatively use "rm infra")
	  -autoconfirm               (Flag) If set it does not ask for confirmation anymore
	  -infra                     (Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations.

Command: deploy infrastructure     Deploy an infrastructure. (alternatively use "apply infra")
	  -allow_data_loss           (Flag) If set, deploy will throw error if data loss is expected.
	  -autoconfirm               (Flag) If set operation procedes without asking for confirmation
	  -infra                     (Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations.
	  -no_attempt_soft_shutdown  (Flag) If set,do not atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy
	  -no_hard_shutdown_after_timeout (Flag) If set do not force a hard power off after timeout expired and the server is not powered off.
	  -skip_ansible              (Flag) If set, some automatic provisioning steps will be skipped. This parameter should generally be ignored.
	  -soft_shutdown_timeout_seconds (Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set.

Command: get infrastructure        Get an infrastructure. (alternatively use "show infra")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.
	  -infra                     (Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations.

Command: revert infrastructure     Revert all changes of an infrastructure. (alternatively use "undo infra")
	  -autoconfirm               (Flag) If set it does not ask for confirmation anymore
	  -infra                     (Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations.

Command: create instance_array     Creates an instance array. (alternatively use "new ia")
	  -boot                      InstanceArray's boot type:'pxe_iscsi','local_drives'
	  -disk_size                 InstanceArray's local disk sizes
	  -disks                     InstanceArray's number of local drives
	  -infra                     (Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations.
	  -instance_count            (Required) Instance count of this instance array
	  -label                     InstanceArray's label
	  -proc                      InstanceArray's minimum processor count
	  -proc_core_count           InstanceArray's minimum processor core count
	  -proc_freq                 InstanceArray's minimum processor frequency (Mhz)
	  -ram                       InstanceArray's minimum RAM (GB)
	  -return_id                 (Flag) If set will print the ID of the created Instance Array. Useful for automating tasks.
	  -template                  InstanceArray's volume template when booting from for local drives
	  -un_managed_fw             (Flag) If set InstanceArray's firewall management on or off

Command: list instance_array       Lists all instance arrays of an infrastructure. (alternatively use "ls ia")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.
	  -infra                     (Required) Infrastrucure ID

Command: delete instance_array     Delete instance array. (alternatively use "rm ia")
	  -autoconfirm               If true it does not ask for confirmation anymore
	  -id                        (Required) InstanceArray ID

Command: edit instance_array       Edits an instance array. (alternatively use "alter ia")
	  -boot                      InstanceArray's boot type:'pxe_iscsi','local_drives'
	  -disk_size                 InstanceArray's local disk sizes
	  -disks                     InstanceArray's number of local drives
	  -do_not_keep_detaching_drives (Flag) If set and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted.
	  -id                        (Required) InstanceArray's id
	  -instance_count            Instance count of this instance array
	  -label                     (Required) InstanceArray's label
	  -proc                      InstanceArray's minimum processor count
	  -proc_core_count           InstanceArray's minimum processor core count
	  -proc_freq                 InstanceArray's minimum processor frequency (Mhz)
	  -ram                       InstanceArray's minimum RAM (GB)
	  -swap_existing_hardware    (Flag) If set all the hardware of the Instance objects is swapped to match the new InstanceArray specifications
	  -template                  InstanceArray's volume template when booting from for local drives
	  -unmanaged_fw              (Flag) If set InstanceArray's firewall management is off

Command: create drive_array        Creates a drive array. (alternatively use "new da")
	  -count                     DriveArrays's drive count. Use this only for unconnected DriveArrays.
	  -expand_with_ia            Auto-expand when the connected instance array expands
	  -ia                        (Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays
	  -infra                     (Required) Infrastrucure ID
	  -label                     (Required) The label of the drive array
	  -return_id                 (Optional) Will print the ID of the created Drive Array. Useful for automating tasks.
	  -size                      (Optional, default = 40960) Drive arrays's size in MBytes
	  -template                  DriveArrays's volume template to clone when creating Drives
	  -type                      Possible values: iscsi_ssd, iscsi_hdd

Command: edit drive_array          Edit a drive array. (alternatively use "alter da")
	  -count                     DriveArrays's drive count. Use this only for unconnected DriveArrays.
	  -expand_with_ia            Auto-expand when the connected instance array expands
	  -ia                        (Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays
	  -id                        (Required) Drive Array's ID
	  -label                     (Required) The label of the drive array
	  -size                      (Optional, default = 40960) Drive arrays's size in MBytes
	  -template                  DriveArrays's volume template to clone when creating Drives
	  -type                      Possible values: iscsi_ssd, iscsi_hdd

Command: list drive_array          Lists all drive arrays of an infrastructure. (alternatively use "ls da")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.
	  -infra                     (Required) Infrastrucure ID

Command: delete drive_array        Delete a drive array. (alternatively use "rm da")
	  -autoconfirm               If true it does not ask for confirmation anymore
	  -id                        (Required) DriveArray's ID

Command: list volume_templates     Lists available volume templates (alternatively use "ls templates")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.
	  -local_only                Show only templates that support local install
	  -pxe_only                  Show only templates that support pxe booting

Command: list firewall_rule        Lists instance array firewall rules (alternatively use "ls fw")
	  -format                    The output format. Supported values are 'json','csv'. The default format is human readable.
	  -ia                        (Required) The instance array id

Command: add firewall_rule         Add instance array firewall rule (alternatively use "new fw")
	  -description               The firewall rule's description.
	  -destination               The destination address to filter on. It can also be a range with the start and end values separated by a dash.
	  -ia                        (Required) The instance array id
	  -ip_address_type           The IP address type of the firewall rule. Possible values: ipv4, ipv6.
	  -port                      The port to filter on. It can also be a range with the start and end values separated by a dash.
	  -protocol                  The protocol of the firewall rule. Possible values: all, icmp, tcp, udp.
	  -source                    The source address to filter on. It can also be a range with the start and end values separated by a dash.

Command: delete firewall_rule      Remove instance array firewall rule (alternatively use "rm fw")
	  -destination               The destination address to filter on. It can also be a range with the start and end values separated by a dash.
	  -ia                        (Required) The instance array id
	  -ip_address_type           The IP address type of the firewall rule. Possible values: ipv4, ipv6.
	  -port                      The port to filter on. It can also be a range with the start and end values separated by a dash.
	  -protocol                  The protocol of the firewall rule. Possible values: all, icmp, tcp, udp.
	  -source                    The source address to filter on. It can also be a range with the start and end values separated by a dash.
```