## metalcloud-cli login

Sign in through the browser

### Synopsis

Sign in to MetalCloud as yourself, through the browser, and save the session.

The command prints a URL and a code. Open the URL in a browser on this machine or on any
other device, check that the page shows the same code, and approve the sign-in. The command
then saves the session in $HOME/.metalcloud/sessions.json (one entry per endpoint, readable
only by you). Commands use it when no API token is set with --api_key or METALCLOUD_API_KEY.

The session ends by the same rules as a browser session, for example after a period without
use, and then the next command asks you to sign in again. Some commands, such as the ones
that manage API tokens, are not available to a CLI session; use the UI or an API token for them.

To sign in as a service account instead, give the path of a file that holds the service account's
provider token, with --service-account-token-file or METALCLOUD_SERVICE_ACCOUNT_TOKEN_FILE. The
command then signs in with the token in the file and saves the session, marked as a service
account session. While the token file is set, other commands use only that session and sign in
again from the file when they have none or it ends within 60 seconds, so automation needs no login
step. A person's saved session for the same endpoint is ignored, not removed. --api_key still wins.

Examples:
  # Sign in, opening the browser on this machine
  metalcloud-cli login

  # Sign in from a machine without a browser; open the printed URL on another device
  metalcloud-cli login --no-browser

  # Sign in as a service account from a token file kept current by your tooling
  metalcloud-cli login --service-account-token-file /var/run/secrets/metalcloud/token

```
metalcloud-cli login [flags]
```

### Options

```
  -h, --help                                help for login
      --no-browser                          Do not open the browser; only print the URL and the code.
      --service-account-token-file string   Sign in as a service account with the provider token in this file (or set METALCLOUD_SERVICE_ACCOUNT_TOKEN_FILE).
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API token (msat_...), or a legacy API key (<user id>:<key>), which works only until it expires
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

