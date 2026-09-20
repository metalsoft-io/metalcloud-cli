## metalcloud-cli endpoint-instance-group network acl add

Add a security rule to a network connection

### Synopsis

Add a security rule to one network connection of an endpoint instance group.

The rule can be described either with a configuration file (--config-source) or with
individual flags (--rule-type, --direction, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Required Flags (one of):
  --config-source string       Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.
  --rule-type string           Rule type: ipv4, ipv6 or mac

Optional Flags (when not using --config-source):
  --direction string           Rule direction: in or out (default "in")
  --sequence int               Evaluation order of the rule
  --forwarding-action string   Forwarding action: allow, deny, transit or discard (default "allow")
  --enforcement-point string   Enforcement point of the rule (default "svi")
  --network-protocol string    Network protocol of the rule
  --source-address string      Source address of the rule
  --destination-address string Destination address of the rule
  --source-port string         Source port of the rule
  --destination-port string    Destination port of the rule

Examples:
  metalcloud-cli endpoint-instance-group network acl add 12 5 --rule-type ipv4 --sequence 10 --source-address 10.0.0.0/24
  metalcloud-cli eig net acl add 12 5 --config-source rule.json

```
metalcloud-cli endpoint-instance-group network acl add endpoint_instance_group_id connection_id [flags]
```

### Options

```
      --config-source string         Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.
      --destination-address string   Destination address of the rule.
      --destination-port string      Destination port of the rule.
      --direction string             Rule direction: in or out. (default "in")
      --enforcement-point string     Enforcement point of the rule. (default "svi")
      --forwarding-action string     Forwarding action: allow, deny, transit or discard. (default "allow")
  -h, --help                         help for add
      --network-protocol string      Network protocol of the rule.
      --rule-type string             Rule type: ipv4, ipv6 or mac.
      --sequence int32               Evaluation order of the rule.
      --source-address string        Source address of the rule.
      --source-port string           Source port of the rule.
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

* [metalcloud-cli endpoint-instance-group network acl](metalcloud-cli_endpoint-instance-group_network_acl.md)	 - Manage the security rules of a network connection

