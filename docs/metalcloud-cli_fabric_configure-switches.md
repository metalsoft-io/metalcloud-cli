## metalcloud-cli fabric configure-switches

Configure all switches of a fabric from a declarative YAML/JSON

### Synopsis

Configure every network device attached to a fabric from one declarative
configuration: hostnames (identifierString), ASNs, loopback IPs, physical-port
enable + interface descriptions, and point-to-point links with deterministic
/31 IPAM subnets.

Each feature section is optional - omit one to skip that step. Every step is
idempotent: current state is read first and only differences are written. Use
--dry-run to compute and preview the full plan without making any changes.

The configuration is supplied EITHER as a whole document via --config-source,
OR built up from the individual per-property flags below. The two are mutually
exclusive; provide one or the other.

Arguments:
  fabric_id    The ID or label of the fabric to configure

Whole-document input:
  --config-source   'pipe' to read from stdin, or a path to a YAML/JSON config file.
                    Run 'fabric configure-switches-example' for a template.

Per-property flags (alternative to --config-source). A feature section is
enabled when its toggle flag is set or any of its sub-flags is provided:
  ordering        --ordering
  enable ports    --enable-physical-ports
  descriptions    --description-template
  hostname        --hostname, --hostname-leaf, --hostname-spine,
                  --hostname-super-spine, --hostname-skip
  asn             --asn, --asn-leaf-start, --asn-spine-start, --asn-super-spine-start
  loopback        --loopback, --loopback-subnet
  topology        --topology-leaf-spine[-links-per-pair],
                  --topology-spine-super-spine[-links-per-pair],
                  --topology-leaf-host[-node-count|-nodes|-port-pattern|
                  -nic-names|-description-template]
  p2p             --p2p, --p2p-pool-leaf-spine, --p2p-pool-spine-super-spine,
                  --p2p-pool-leaf-host, --p2p-mtu

Always available:
  --dry-run         Compute the plan and report what would change, without writing.

Examples:
  # Whole-document input from a file or stdin
  metalcloud-cli fabric configure-switches 5 --config-source fabric-config.yaml --dry-run
  cat fabric-config.yaml | metalcloud-cli fabric configure-switches 5 --config-source pipe

  # Per-property flags: hostnames + ASNs + loopbacks with the reference defaults
  metalcloud-cli fabric configure-switches 5 --hostname --asn --loopback --dry-run

  # Per-property flags: full leaf/spine fabric with links and a custom pool
  metalcloud-cli fabric configure-switches my-fabric \
    --hostname --asn --loopback \
    --topology-leaf-spine --topology-leaf-host-node-count 8 \
    --description-template "to_{peerHostname}_{peerPort}" \
    --p2p --p2p-pool-leaf-spine 10.254.0.0/16 --p2p-mtu 9216

  # Override one ASN start and skip naming the spines
  metalcloud-cli fabric configure-switches 5 --asn --asn-leaf-start 4200001000 \
    --hostname --hostname-skip spine

```
metalcloud-cli fabric configure-switches fabric_id [flags]
```

### Options

```
      --asn                                                Enable ASN assignment using the default starts.
      --asn-leaf-start int                                 Starting ASN for leaves.
      --asn-spine-start int                                Starting ASN for spine groups.
      --asn-super-spine-start int                          Shared ASN for superspines.
      --config-source string                               Source of the switch configuration. Can be 'pipe' or path to a YAML/JSON file. Mutually exclusive with the per-property flags below.
      --description-template string                        Interface description template (placeholders {peerHostname}, {peerPort}). Requires a topology section.
      --dry-run                                            Compute and preview the plan without making any changes.
      --enable-physical-ports                              Enable every physical port's staged config. (default true)
  -h, --help                                               help for configure-switches
      --hostname                                           Enable hostname computation using the built-in reference templates.
      --hostname-leaf string                               Hostname template for leaf devices.
      --hostname-skip strings                              Positions to skip (set to null), e.g. spine.
      --hostname-spine string                              Hostname template for spine devices.
      --hostname-super-spine string                        Hostname template for super_spine devices.
      --loopback                                           Enable loopback IP allocation using the default subnet.
      --loopback-subnet string                             Pool the loopback /32s are carved from.
      --ordering string                                    Device ordering: managementAddress | identifierString | id. (default "managementAddress")
      --p2p                                                Enable point-to-point link creation with reference default pools.
      --p2p-mtu int32                                      MTU applied to created links.
      --p2p-pool-leaf-host string                          Leaf->host /31 pool.
      --p2p-pool-leaf-spine string                         Leaf<->spine /31 pool.
      --p2p-pool-spine-super-spine string                  Spine<->superspine /31 pool.
      --topology-leaf-host                                 Enable leaf->host downlinks.
      --topology-leaf-host-description-template string     Leaf->host description template.
      --topology-leaf-host-nic-names strings               Remote host NIC names (even count).
      --topology-leaf-host-node-count int                  Number of host port-pairs per leaf.
      --topology-leaf-host-nodes ints                      Exact 0-based node indices (mutually exclusive with node-count).
      --topology-leaf-host-port-pattern string             Leaf host port pattern, e.g. swp{port}s{sub}.
      --topology-leaf-spine                                Enable leaf<->spine pairing.
      --topology-leaf-spine-links-per-pair string          Leaf<->spine links per pair: 'auto' or an integer.
      --topology-spine-super-spine                         Enable spine<->superspine pairing (3-tier only).
      --topology-spine-super-spine-links-per-pair string   Spine<->superspine links per pair: 'auto' or an integer.
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

