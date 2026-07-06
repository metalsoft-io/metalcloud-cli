## metalcloud-cli license

Manage the system license

### Synopsis

Manage the license installed on the MetalSoft system.

A license controls which product categories are enabled (servers, switches, VMs,
storage) and the resource allowances granted to the installation.

Available Commands:
  status        Show the validity status of the current license
  get           Show the installed license document (Base64)
  allowance     Show the resource allowance granted by the license
  products      Show which product categories are licensed
  request       Show the license request document to send to MetalSoft
  add           Install a license on the system

Typical workflow:
  # 1. Generate the request document and send it to MetalSoft
  metalcloud-cli license request > license-request.txt

  # 2. Install the signed license returned by MetalSoft
  metalcloud-cli license add --source ./license.txt

  # 3. Verify it took effect
  metalcloud-cli license status

### Options

```
  -h, --help   help for license
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli license add](metalcloud-cli_license_add.md)	 - Install a license on the system
* [metalcloud-cli license allowance](metalcloud-cli_license_allowance.md)	 - Show the resource allowance granted by the license
* [metalcloud-cli license get](metalcloud-cli_license_get.md)	 - Show the installed license document
* [metalcloud-cli license products](metalcloud-cli_license_products.md)	 - Show which product categories are licensed
* [metalcloud-cli license request](metalcloud-cli_license_request.md)	 - Show the license request document
* [metalcloud-cli license status](metalcloud-cli_license_status.md)	 - Show the validity status of the current license

