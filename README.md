# metalcloud-cli

[![Build Status](https://travis-ci.org/bigstepinc/metalcloud-cli.svg?branch=master)](https://travis-ci.org/bigstepinc/metalcloud-cli)
[![Coverage Status](https://coveralls.io/repos/github/bigstepinc/metalcloud-cli/badge.svg?branch=master)](https://coveralls.io/github/bigstepinc/metalcloud-cli?branch=master)

This tool allows the manipulation of all Bigstep Metal Cloud elements via the command line.

![metalcloud-cli](https://bigstep.com/assets/images/blog/2019/metalcloud-cli-animated.gif)

## Installation

To install on Mac OS X:
```
brew tap bigstepinc/homebrew-repo
brew install metalcloud-cli
```

To install on any CentOS/Redhat Linux distribution:
```
$ sudo rpm -i https://github.com/bigstepinc/metalcloud-cli/releases/download/v1.0.3/metalcloud-cli_1.0.3_linux_amd64.rpm
```

To install on any Debian/Ubuntu distributions:
```
curl -sLO https://github.com/bigstepinc/metalcloud-cli/releases/download/v1.0.3/metalcloud-cli_1.0.3_linux_amd64.deb && sudo dpkg -i metalcloud-cli_1.0.3_linux_amd64.deb
```

To install using `go get` (this should also work on Windows):
```bash
go get github.com/bigstepinc/metalcloud-cli
```

Configure credentials as environment variables:
```bash
export METALCLOUD_API_KEY="<your key>"
export METALCLOUD_ENDPOINT="https://api.bigstep.com"
export METALCLOUD_USER_EMAIL="<your email>"
export METALCLOUD_DATACENTER="uk-reading"
```

## Getting a list of supported commands

Use `metalcloud-cli help` for a list of supported commands.


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

## Using label instead of IDs

Most commands also take a label instead of an id as a parameter. For example:
```
metalcloud-cli show infra -id complex-demo
```

## Admin commands
To enable admin commands use the following environment variable:
```bash
export METALCLOUD_ADMIN="true"
```