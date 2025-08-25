## metalcloud-cli file-share update-meta

Update file share metadata

### Synopsis

Update the metadata of an existing file share with new information.

This command allows you to modify metadata properties of a file share such as
labels, tags, descriptions, and other custom metadata fields.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to update

Required Flags:
  --config-source               Source of the file share metadata updates
                               Can be 'pipe' for stdin input or path to a JSON file

Metadata Format:
The configuration should be a JSON object containing the metadata properties to update:
- labels: Key-value pairs for labeling
- tags: Array of tags for categorization
- description: Detailed description of the file share
- custom_fields: Custom metadata fields

Examples:
  # Update file share metadata from a JSON file
  metalcloud-cli file-share update-meta my-infrastructure 12345 --config-source metadata.json

  # Update using pipe input
  echo '{"labels":{"env":"production","team":"devops"}}' | metalcloud-cli file-share update-meta my-infrastructure 12345 --config-source pipe

  # Update with infrastructure ID
  metalcloud-cli file-share update-meta 100 12345 --config-source /path/to/metadata.json

```
metalcloud-cli file-share update-meta infrastructure_id_or_label file_share_id [flags]
```

### Options

```
      --config-source string   Source of the file share metadata updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-meta
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API key
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli file-share](metalcloud-cli_file-share.md)	 - Manage file shares for infrastructure resources

