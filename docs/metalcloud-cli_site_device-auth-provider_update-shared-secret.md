## metalcloud-cli site device-auth-provider update-shared-secret

Rotate the shared secret for a device auth provider

### Synopsis

Rotate the shared secret used to encrypt communication with the auth server.
The provider must be in 'maintenance' or 'disabled' status and must have no
network devices linked to it.

```
metalcloud-cli site device-auth-provider update-shared-secret provider_id_or_label [flags]
```

### Options

```
  -h, --help            help for update-shared-secret
      --secret string   The new shared secret value.
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

* [metalcloud-cli site device-auth-provider](metalcloud-cli_site_device-auth-provider.md)	 - Manage device authentication providers (e.g. TACACS+)

