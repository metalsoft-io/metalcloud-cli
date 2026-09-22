## metalcloud-cli server-registration-profile search

Find the registration profile that applies to a site

### Synopsis

Find the server registration profile that applies to a site.

The profile assigned to the site is returned when there is one; otherwise the
system-wide default profile is returned.

Required Flags:
  --site-id    The ID of the site to search the registration profile for

Examples:
  # Find the registration profile used by site 1
  metalcloud-cli server-registration-profile search --site-id 1


```
metalcloud-cli server-registration-profile search [flags]
```

### Options

```
  -h, --help          help for search
      --site-id int   The ID of the site to search the registration profile for.
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

