## metalcloud-cli bucket update-meta

Update bucket metadata and custom properties

### Synopsis

Update the metadata and custom properties of an existing bucket.

This command allows you to modify bucket metadata such as labels, descriptions, 
custom tags, and other non-configuration properties. The metadata updates are 
provided through either a JSON file or piped input.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket to update

Required Flags:
  --config-source string       Source of the bucket metadata updates
                               Accepts either 'pipe' for piped JSON input or a path to a JSON file

Examples:
  # Update bucket metadata from a JSON file
  metalcloud-cli bucket update-meta 100 42 --config-source metadata.json

  # Update metadata using piped input
  echo '{"label": "production-storage", "description": "Main storage bucket"}' | metalcloud-cli bucket update-meta production bucket-123 --config-source pipe

  # Update bucket metadata with file
  metalcloud-cli bucket meta-update staging my-bucket --config-source /configs/meta.json

```
metalcloud-cli bucket update-meta infrastructure_id_or_label bucket_id [flags]
```

### Options

```
      --config-source string   Source of the bucket metadata updates. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli bucket](metalcloud-cli_bucket.md)	 - Manage S3-compatible object storage buckets

