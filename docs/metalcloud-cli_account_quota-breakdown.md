## metalcloud-cli account quota-breakdown

Show how the quota limits of an account are derived

### Synopsis

Show the quota limits breakdown of an account in the MetalCloud platform.

The breakdown lists, for every quota limit, the effective value together with the
value contributed by the quota profile of the account and by the quota profile of its
parent account. The effective value is the most restrictive of the two; individual
users may still be further restricted by their role or group quota profiles.

Text, CSV and Markdown output flatten the breakdown into one row per limit. JSON and
YAML output keep the original nested object returned by the API.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account

Optional Flags:
  --include-usage    Also report the aggregate resource counts of the account and of
                     the parent account (adds two columns to the table)

Examples:
  # Show the quota breakdown of an account
  metalcloud-cli account quota-breakdown 1234

  # Include the current usage
  metalcloud-cli account quota-breakdown 1234 --include-usage

  # Keep the raw object
  metalcloud-cli account quota 1234 -f json

```
metalcloud-cli account quota-breakdown account_id [flags]
```

### Options

```
  -h, --help            help for quota-breakdown
      --include-usage   Also report the aggregate resource usage of the account and of its parent account.
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

* [metalcloud-cli account](metalcloud-cli_account.md)	 - Manage user accounts and account-related operations

