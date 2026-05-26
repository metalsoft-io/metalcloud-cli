## metalcloud-cli site device-auth-provider

Manage device authentication providers (e.g. TACACS+)

### Synopsis

Manage device authentication providers used by network devices for AAA
(authentication, authorization, accounting). Currently supports TACACS+ providers.

Available Commands:
  list                  List all device auth providers
  get                   Retrieve a specific device auth provider
  create                Create a new device auth provider from JSON
  update                Update an existing device auth provider from JSON
  delete                Delete a device auth provider
  credentials           Retrieve decrypted credentials for a provider
  update-shared-secret  Rotate the shared secret for a provider
  config-example        Print a JSON template suitable for create

### Options

```
  -h, --help   help for device-auth-provider
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

* [metalcloud-cli site](metalcloud-cli_site.md)	 - Manage sites (datacenters) and their configurations
* [metalcloud-cli site device-auth-provider config-example](metalcloud-cli_site_device-auth-provider_config-example.md)	 - Print a JSON template for creating a device auth provider
* [metalcloud-cli site device-auth-provider create](metalcloud-cli_site_device-auth-provider_create.md)	 - Create a device auth provider from JSON
* [metalcloud-cli site device-auth-provider credentials](metalcloud-cli_site_device-auth-provider_credentials.md)	 - Show the decrypted credentials for a device auth provider
* [metalcloud-cli site device-auth-provider delete](metalcloud-cli_site_device-auth-provider_delete.md)	 - Delete a device auth provider
* [metalcloud-cli site device-auth-provider get](metalcloud-cli_site_device-auth-provider_get.md)	 - Retrieve a device auth provider by ID or label
* [metalcloud-cli site device-auth-provider list](metalcloud-cli_site_device-auth-provider_list.md)	 - List all device auth providers
* [metalcloud-cli site device-auth-provider update](metalcloud-cli_site_device-auth-provider_update.md)	 - Update a device auth provider from JSON
* [metalcloud-cli site device-auth-provider update-shared-secret](metalcloud-cli_site_device-auth-provider_update-shared-secret.md)	 - Rotate the shared secret for a device auth provider

