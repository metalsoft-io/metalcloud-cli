# metalcloud-cli

![Build](https://github.com/metalsoft-io/metalcloud-cli/actions/workflows/build.yml/badge.svg)

The MetalCloud CLI allows management of all MetalCloud instance objects via the command line.

[www.metalsoft.io](https://metalsoft.io)

## Installation

> [!IMPORTANT]
> Refer to the [latest release](https://github.com/metalsoft-io/metalcloud-cli/releases/latest) for the correct package you need to install.

To install on Mac OS X:

```bash
brew tap metalsoft-io/homebrew-repo
brew install metalcloud-cli
```

In case you need to refresh the local homebrew cache after a new version of the CLI is released:

```bash
brew update
```

To install on any CentOS/Redhat Linux distribution:

```bash
sudo yum install $(curl -s https://api.github.com/repos/metalsoft-io/metalcloud-cli/releases/latest | grep -i browser_download_url | grep "amd64" | grep -i linux | grep rpm | head -n 1 | cut -d'"' -f4)
```

To install on any Debian/Ubuntu distributions:

```bash
curl -skL $(curl -s https://api.github.com/repos/metalsoft-io/metalcloud-cli/releases/latest | grep -i browser_download_url  | grep "$(dpkg --print-architecture)" | grep deb | head -n 1 | cut -d'"' -f4) -o metalcloud-cli.deb && sudo dpkg -i metalcloud-cli.deb
```

To install on Windows:
Binaries are available [here](https://github.com/metalsoft-io/metalcloud-cli/releases/latest):

```bash
https://github.com/metalsoft-io/metalcloud-cli/releases/latest
```

To install using `go get` (this should also work on Windows):

```bash
go get github.com/metalsoft-io/metalcloud-cli
```

## Getting an API token

In the MetalCloud's UI go to the upper left corner and click on your initials. Then go to **Account Settings** > **API Tokens**, create a token and copy its value (it starts with `msat_`). The value is shown only once, when the token is created or regenerated.

To run administration commands (servers, accounts, users and other resources that need admin permissions), choose the **System** scope when creating the token. A token with the default **Account** scope only grants what any user of your account can do, and those commands fail with `403 Forbidden`. The scope choice is shown only to users with the `global_token_write` permission.

API keys created before version 8.0 (of the form `<user id>:<key>`) keep working as API tokens with the **System** scope until they expire, so an existing configuration keeps working unchanged. API tokens are created, regenerated and deleted only in the UI.

Configure credentials as environment variables:

```bash
export METALCLOUD_ENDPOINT="https://metal.mycompany.com"
export METALCLOUD_API_KEY="<your API token>"
```

Alternatively you can put the endpoint and API token configuration in a `metalcloud.yaml` configuration file:

```yaml
endpoint: 'https://metal.mycompany.com'
api_key: '<your API token>'
```

The environment variable and the configuration key keep the name `api_key` for compatibility with existing configurations; they accept both an API token and a legacy API key.

## Signing in

Instead of an API token you can sign in as yourself, with whatever you use to sign in to the UI (password with MFA, LDAP, SAML or OIDC):

```bash
metalcloud-cli login
```

The command prints a URL and a code, and opens the URL in your browser. Check that the page shows the same code, confirm your identity and approve. On a machine without a browser, run `metalcloud-cli login --no-browser` and open the printed URL on any other device. The endpoint (`--endpoint`, `METALCLOUD_ENDPOINT` or `metalcloud.yaml`) must be set; no API token is needed.

The session is saved in `$HOME/.metalcloud/sessions.json` (the folder is readable only by you, and so is the file), one entry per endpoint. `metalcloud.yaml` is never changed.

Each command uses, in this order:

1. the API token from `--api_key`, `METALCLOUD_API_KEY` or `metalcloud.yaml`;
2. else the service account token file, if one is set (see below);
3. else the saved session for the endpoint, if it has not expired;
4. else it stops and asks you to sign in.

A session ends by the same rules as a browser session, for example after an hour without use or after a password change. The next command then removes it and asks you to run `metalcloud-cli login` again. Only one CLI session per user is kept: signing in again, on any machine, ends the previous one.

`metalcloud-cli logout` ends the session on the server and removes it from the file, even if the server can't be reached.

A CLI session has the same access as your UI session, except for what could be misused if a sign-in were approved by mistake: it can't manage API tokens, change your password, email or two-factor authentication, or re-authenticate. Those commands fail with `403 This route does not accept a CLI session; use the UI or an API token.`, and the session stays. The user also needs the `api_key_use` permission, as for API tokens. Older servers without this feature answer: `This MetalCloud doesn't support signing in from the CLI; use an API token.`

### Signing in as a service account

Automation (a CI job, a Kubernetes workload) can sign in as a service account with a token from your OIDC provider, without a person or a browser:

```bash
metalcloud-cli login --service-account-token-file /var/run/secrets/metalcloud/token
```

Set `--service-account-token-file` or the environment variable `METALCLOUD_SERVICE_ACCOUNT_TOKEN_FILE` (or `service-account-token-file` in `metalcloud.yaml`) to the path of a file that holds the provider token. Your tooling keeps the file current; the CLI reads it, trimmed, at each sign-in and never keeps a copy.

While the token file is set, every command signs in with it, so `login` is optional for automation. The session is saved in `$HOME/.metalcloud/sessions.json` as a service account session, separate from your own session for the same endpoint, which is ignored but not removed. Each token file has its own session, so two service accounts can share one `$HOME`. A command signs in again from the file when there is no such session or it ends within 60 seconds, or when the server has ended the session earlier (the command's first calls then get a `401`; the CLI signs in once and repeats them). If the session can't be saved, for example on a read-only `$HOME`, the command still runs and prints a warning on stderr, and the next command signs in again. If a command gets a `401` after it has started, it fails with `The service account's session ended; the next command signs in again from the token file.` `--api_key` and `METALCLOUD_API_KEY` still take precedence.

A rejected token answers `The provider token wasn't accepted; check the token file and the service account's subject.`, and a MetalCloud that has service accounts turned off answers `This MetalCloud doesn't accept service accounts.` `metalcloud-cli logout` removes only your own session.

An administrator creates the service account with the subject (`sub`) that the provider puts in its tokens; a service account has no password:

```bash
metalcloud-cli user create --email robot@mycompany.com --service-account-subject client-1 --access-level user
```

`metalcloud-cli user get` shows the provider (`service-account`) in its text output, and `-f json` or `-f yaml` of `user get` show `provider` and `serviceAccountSubject`; `user list -f json` or `-f yaml` shows `provider`.

## Getting a list of supported commands

Use `metalcloud-cli --help` for a list of supported commands.

## Getting started

To create an infrastructure:

```bash
metalcloud-cli infrastructure create 11 demo

metalcloud-cli infrastructure list
```

```txt
┌──────┬───────────┬──────────────┬─────────┬───────┬──────┬──────────────────────┬──────────────────────┬───────────────┬───────────┐
│    # │ LABEL     │ CONFIG LABEL │ STATUS  │ OWNER │ SITE │ CREATED              │ UPDATED              │ DEPLOY STATUS │ DEPLOY ID │
├──────┼───────────┼──────────────┼─────────┼───────┼──────┼──────────────────────┼──────────────────────┼───────────────┼───────────┤
│ 1689 │ demo      │ demo         │ ordered │    10 │   11 │ 24 Mar 25 09:52 EET  │ 24 Mar 25 10:19 EET  │ not_started   │           │
└──────┴───────────┴──────────────┴─────────┴───────┴──────┴──────────────────────┴──────────────────────┴───────────────┴───────────┘
```

To create an server instance group with one instance and label `master` in that infrastructure, get the ID of the infrastructure from above (1689), the server type ID (12):

```bash
metalcloud-cli server-instance-group create 1689 master 12 1
```

To view the details of the created server instance group:

```bash
metalcloud-cli server-instance-group list 1689
```

```txt
┌─────┬─────────┬──────────┬─────────┬─────────────────────┬─────────────────────┐
│   # │ LABEL   │ INFRA ID │ STATUS  │ CREATED             │ UPDATED             │
├─────┼─────────┼──────────┼─────────┼─────────────────────┼─────────────────────┤
│ 824 │ master  │     1689 │ ordered │ 24 Mar 25 21:49 EET │ 24 Mar 25 21:49 EET │
└─────┴─────────┴──────────┴─────────┴─────────────────────┴─────────────────────┘
```

To create a drive array and attach it to the previous instance array:

```bash
echo '{"label": "my-drive", "logicalNetworkId": 5, "sizeMb": 2024}' | metalcloud-cli drive create 1689 --config-source pipe
```

To view the current status of the infrastructure:

```bash
metalcloud-cli infrastructure get 1689
```

```txt
┌──────┬───────────┬──────────────┬─────────┬───────┬──────┬──────────────────────┬──────────────────────┬───────────────┬───────────┐
│    # │ LABEL     │ CONFIG LABEL │ STATUS  │ OWNER │ SITE │ CREATED              │ UPDATED              │ DEPLOY STATUS │ DEPLOY ID │
├──────┼───────────┼──────────────┼─────────┼───────┼──────┼──────────────────────┼──────────────────────┼───────────────┼───────────┤
│ 1689 │ demo      │ demo         │ ordered │    10 │   11 │ 24 Mar 25 09:52 EET  │ 24 Mar 25 10:19 EET  │ not_started   │           │
└──────┴───────────┴──────────────┴─────────┴───────┴──────┴──────────────────────┴──────────────────────┴───────────────┴───────────┘
```

## Apply support

Apply creates or updates a resource from a file. The supported format is yaml.

```bash
metalcloud-cli apply -f resources.yaml
```

The type of the requested resource needs to be specified using the field *kind*.

```yaml
cat resources.yaml

kind: InstanceArray
apiVersion: 1.0
label: my-instance-array

---

kind: Secret
apiVersion: 1.0
name: my-secret

```

The objects and their fields can be found in the [SDK documentation](https://pkg.go.dev/github.com/metalsoft-io/metalcloud-sdk-go). The fields will be in the format specified in the yaml tag.

## Aliases

The CLI also provides aliases for most of it's commands:

* server-instance-group = ia
* infrastructure = infra
* list = ls
* delete = rm
...

This allows commands such as:

```bash
metalcloud-cli infra ls
```

## Using label instead of IDs

Most commands also take a label instead of an id as a parameter. For example:

```bash
metalcloud-cli site get demo-site
```

## Permissions

Some commands depend on various permissions. For instance you cannot access another user's infrastructure unless you are a delegate of it.

## Admin commands

If the user has admin permissions, additional commands will be available.

## Debugging information

To enable debugging information in the output/CLI add the `-d` flag to the command, this will print out the raw requests being made and it's usefull to identify API communication issues.

To run the CLI in VS Code in debug mode add the following configuration in `launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch metalcloud-cli",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/metalcloud-cli",
            "buildFlags": "-ldflags='-X main.version=v7.0.0 -X main.allowDevelop=true'",
            "args": ["extension", "get", "24", "-d"],
            "showLog": false,
            "dlvFlags": ["--check-go-version=false"],
            "env": {
                "METALCLOUD_API_KEY": "<your API token>",
                "METALCLOUD_ENDPOINT": "https://metal.mycompany.com"
            }
        }
    ]
}
```

## Building the CLI

### Running tests

**Unit tests** (no external dependencies, run anywhere):

```bash
go test ./...
```

**Unit tests with coverage report:**

```bash
go test $(go list ./... | grep -v '/internal/testutils\|/internal/vm$\|/pkg/api\|/pkg/logger\|/pkg/repo\|/cmd/metalcloud-cli$\|/system$') -coverprofile=coverage.out && go tool cover -func=coverage.out | tail -1
```

**Integration tests** (read-only, requires a real MetalCloud endpoint):

```bash
export METALCLOUD_ENDPOINT="https://your-instance.metalsoft.io"
export METALCLOUD_API_KEY="<your API token>"
go test -tags integration -v ./test/integration/
```

The tests default to `./metalcloud-cli` in the repo root. Build it first with `./ron_build` or `go build ./cmd/metalcloud-cli/`. Override with `METALCLOUD_CLI_BINARY`:

```bash
METALCLOUD_CLI_BINARY=/usr/local/bin/metalcloud-cli go test -tags integration ./test/integration/
```

Integration tests are read-only (list/get only). They skip automatically when `METALCLOUD_ENDPOINT` is not set.

**CRUD integration tests** (creates and immediately deletes real resources — use a dedicated test environment):

```bash
export METALCLOUD_ENDPOINT="https://your-instance.metalsoft.io"
export METALCLOUD_API_KEY="<your API token>"
export METALCLOUD_TEST_SITE_ID="<site id>"
# Optional — enables additional tests:
# export METALCLOUD_TEST_FABRIC_ID="<fabric id>"      # logical-network-profile CRUD
# export METALCLOUD_TEST_INFRA_ID="<infra id>"        # server-instance-group CRUD
# export METALCLOUD_TEST_SERVER_TYPE_ID="<type id>"   # server-instance-group CRUD
go test -tags integration_crud -v ./test/integration/
```

To build manually:

```bash
go build ./cmd/metalcloud-cli/
```

To build manually with FIPS compliant crypto libraries:

```bash
GOEXPERIMENT=boringcrypto go build ./cmd/metalcloud-cli/
```

To build for development add the following flags:

```bash
go build -ldflags="-X main.version=v7.0.0 -X main.allowDevelop=true" ./cmd/metalcloud-cli`
```

The build process is automated by travis. Just push into the repository using the appropriate tag and the binaries will be created
for Windows/Linux/Mac and also pushed to the [Homebrew Private Repo](https://github.com/metalsoft-io/homebrew-repo).

Use `git tag` to get the last tag:

```bash
git tag
v1.6.7
v1.6.8
v1.6.9
v1.7.0
v1.7.1
v1.7.2
v1.7.3
...
v1.7.4
v1.7.5
v1.7.6
v1.7.7
v1.7.8
```

Push new changes with new tag:

```bash
git add .
git commit -m "commit comment"
git tag v7.0.1
git push --tags
```

It is a good idea to update the main branch as well (with no tag):

```bash
git push
```

## Updating the SDK

To update the SDK follow the instructions in the [BUILDING.md](https://github.com/metalsoft-io/metalcloud-sdk-go/blob/main/BUILDING.md)
