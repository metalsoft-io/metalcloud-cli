## metalcloud-cli bucket get-credentials

Get access credentials for a bucket

### Synopsis

Retrieve access credentials for a specific bucket within an infrastructure.

This command displays the credentials required to access the bucket programmatically,
including access keys, secrets, and endpoint information. These credentials can be
used with S3-compatible tools and SDKs to interact with the bucket.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket

Examples:
  # Get credentials for bucket with ID 42
  metalcloud-cli bucket get-credentials 100 42

  # Get credentials using infrastructure label
  metalcloud-cli bucket get-credentials production bucket-abc123

  # Display credentials using alias
  metalcloud-cli bucket credentials staging my-bucket-id

```
metalcloud-cli bucket get-credentials infrastructure_id_or_label bucket_id [flags]
```

### Options

```
  -h, --help   help for get-credentials
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

