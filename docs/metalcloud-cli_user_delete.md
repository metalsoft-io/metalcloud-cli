## metalcloud-cli user delete

Delete a user and erase their personal information

### Synopsis

Delete a user account and irreversibly erase their personal information.

WARNING: this is NOT the same as 'user archive'. Archiving only marks the account as
inactive and can be undone with 'user unarchive'. Deleting archives the user AND
permanently removes their personally identifiable information; it cannot be undone
and 'user unarchive' will not bring the information back.

Required Arguments:
  user_id                 The numeric ID of the user to delete

Examples:
  metalcloud-cli user delete 12345
  metalcloud-cli user rm 12345

```
metalcloud-cli user delete user_id [flags]
```

### Options

```
  -h, --help   help for delete
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

* [metalcloud-cli user](metalcloud-cli_user.md)	 - Manage user accounts and their properties

