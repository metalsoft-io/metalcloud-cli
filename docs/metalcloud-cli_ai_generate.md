## metalcloud-cli ai generate

Ask the AI assistant a question

### Synopsis

Send a prompt to the MetalSoft AI assistant and print its answer.

The prompt is given as positional text or read from a configuration document with
--config-source. The API requires the question to be scoped to a datacenter.

Required Arguments (one of):
  prompt...                The question to ask. All positional arguments are joined with spaces.
  --config-source string   Source of the request document ('pipe' or path to a JSON/YAML file)
                           with a "prompt" and optionally a "datacenter" field.

Required Flags:
  --datacenter string      The datacenter the question is about. May instead be
                           supplied by the --config-source document.

Text output prints the answer as prose; use -f json or -f yaml for the full
response object.

Examples:
  metalcloud ai generate --datacenter dc1 "which servers are available?"
  metalcloud ai generate --datacenter dc1 --config-source prompt.json
  echo '{"datacenter":"dc1","prompt":"list my infrastructures"}' | metalcloud ai generate --config-source pipe

```
metalcloud-cli ai generate [prompt...] [flags]
```

### Options

```
      --config-source string   Source of the AI request. Can be 'pipe' or path to a JSON/YAML file.
      --datacenter string      The datacenter the question is about.
  -h, --help                   help for generate
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

* [metalcloud-cli ai](metalcloud-cli_ai.md)	 - MetalSoft AI assistant

