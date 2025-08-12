## metalcloud-cli firmware-catalog create

Create a new firmware catalog from vendor sources

### Synopsis

Create a new firmware catalog from vendor sources.

This command creates a firmware catalog by downloading and processing firmware information
from vendor repositories. The catalog can be configured for online or offline updates.

Configuration Methods:
1. Command-line flags - Specify all options using individual flags
2. Configuration file - Use --config-source to load settings from JSON/YAML file

Required Flags (when not using --config-source):
  --name            Name of the firmware catalog
  --vendor          Vendor type: 'dell', 'hp', or 'lenovo'
  --update-type     Update method: 'online' or 'offline'

Source Configuration (mutually exclusive):
  --vendor-url                 URL of the online vendor catalog
  --vendor-local-catalog-path  Path to a local catalog file

Optional Flags:
  --description                   Description of the firmware catalog
  --vendor-token                  Authentication token for vendor API access
  --server-types                  Comma-separated list of Metalsoft server types to filter
  --vendor-systems                Comma-separated list of vendor system models to filter
  --vendor-local-binaries-path    Local directory for downloaded firmware binaries
  --download-binaries             Download firmware binaries locally
  --upload-binaries               Upload binaries to offline repository

Offline Repository Configuration (required when --upload-binaries is used):
  --repo-base-url        Base URL of the offline repository
  --repo-ssh-host        SSH hostname:port for repository upload
  --repo-ssh-user        SSH username for repository access
  --repo-ssh-path        Target directory path on SSH server

SSH Configuration (mutually exclusive):
  --user-private-key-path    Path to SSH private key (default: ~/.ssh/id_rsa)
  --known-hosts-path         Path to SSH known hosts file (default: ~/.ssh/known_hosts)
  --ignore-host-key-check    Skip SSH host key verification

```
metalcloud-cli firmware-catalog create [flags]
```

### Examples

```

Dell example (online):
metalcloud-cli firmware-catalog create \
  --name "Dell R640 Catalog" \
  --description "Dell PowerEdge R640 firmware catalog" \
  --vendor dell \
  --vendor-url https://downloads.dell.com/FOLDER06417267M/1/ESXi_Catalog.xml.gz \
  --vendor-systems "R640" \
  --server-types "M.24.64.2,M.32.64.2" \
  --update-type online

Dell example (offline with upload):
metalcloud-cli firmware-catalog create \
  --name "Dell Offline Catalog" \
  --vendor dell \
  --vendor-url https://downloads.dell.com/catalog.xml.gz \
  --download-binaries \
  --vendor-local-binaries-path ./downloads \
  --upload-binaries \
  --repo-base-url http://repo.mycloud.com/dell \
  --repo-ssh-host repo.mycloud.com:22 \
  --repo-ssh-user admin \
  --repo-ssh-path /var/www/html/dell \
  --update-type offline

HP example:
metalcloud-cli firmware-catalog create \
  --name "HP Gen11 Catalog" \
  --vendor hp \
  --vendor-url https://downloads.linux.hpe.com/SDR/repo/fwpp-gen11/current/fwrepodata/fwrepo.json \
  --server-types "M.8.8.2.v5" \
  --update-type online

Lenovo example using config file:
metalcloud-cli firmware-catalog create --config-source ./lenovo-config.json

Configuration file examples:

lenovo-config.json:
{
  "name": "Lenovo Catalog",
  "description": "Lenovo server firmware catalog",
  "vendor": "lenovo",
  "update_type": "offline",
  "vendor_local_catalog_path": "./lenovo_catalogs",
  "vendor_local_binaries_path": "./lenovo_downloads",
  "server_types_filter": ["M.8.8.2.v5"],
  "vendor_systems_filter": ["7Y51"],
  "download_binaries": true
}

hp-config.yaml:
name: HP Gen11 Catalog
description: HP ProLiant Gen11 firmware
vendor: hp
update_type: online
vendor_url: https://downloads.linux.hpe.com/SDR/repo/fwpp-gen11/current/fwrepodata/fwrepo.json
vendor_local_catalog_path: ./fwrepo.json
vendor_local_binaries_path: ./hp_downloads
server_types_filter:
  - M.8.8.2.v5
```

### Options

```
      --config-source string                Source of the new firmware catalog configuration. Can be 'pipe' or path to a JSON file.
      --description string                  Description of the firmware catalog
      --download-binaries                   Download binaries from the vendor catalog
  -h, --help                                help for create
      --ignore-host-key-check               Ignore host key check for SSH connections
      --known-hosts-path string             Path to the known hosts file for SSH connections (default "~/.ssh/known_hosts")
      --name string                         Name of the firmware catalog
      --repo-base-url string                Base URL of the offline repository
      --repo-ssh-host string                SSH host with port of the offline repository
      --repo-ssh-path string                The path to the target folder in the SSH repository
      --repo-ssh-user string                SSH user for the offline repository
      --server-types strings                List of supported Metalsoft server types (comma-separated)
      --update-type string                  Update type (e.g., 'online', 'offline')
      --upload-binaries                     Upload binaries to the offline repository
      --user-private-key-path string        Path to the user's private SSH key (default "~/.ssh/id_rsa")
      --vendor string                       Vendor type (e.g., 'dell', 'hp')
      --vendor-local-binaries-path string   Path to the local binaries directory
      --vendor-local-catalog-path string    Path to the local catalog file
      --vendor-systems strings              List of supported vendor systems (comma-separated)
      --vendor-token string                 Token for accessing the online vendor catalog
      --vendor-url string                   URL of the online vendor catalog
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

* [metalcloud-cli firmware-catalog](metalcloud-cli_firmware-catalog.md)	 - Manage firmware catalogs for server hardware updates

###### Auto generated by spf13/cobra on 11-Aug-2025
