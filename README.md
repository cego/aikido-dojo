# aikido-dojo

A command-line client for the [Aikido Security](https://www.aikido.dev) public REST API, built
for AI agents first and people second. Each of the API's 206 operations is a command, generated
from Aikido's published OpenAPI spec.

aikido-dojo is unofficial and not affiliated with Aikido Security.

## Install

```sh
brew install cego/tap/aikido-dojo
```

Or, with Go 1.26 or later:

```sh
go install github.com/cego/aikido-dojo@latest
```

The macOS binaries are not signed or notarized yet. Gatekeeper may refuse the first run of a
Homebrew install with "cannot be opened". To allow it, run
`xattr -d com.apple.quarantine "$(readlink -f "$(command -v aikido-dojo)")"`, or install with
`go install`, which builds the binary on your machine.

## Authenticate

Create an API client in Aikido's workspace settings, with only the scopes you need: read scopes
cover everything but writes. Then store it as a profile, and check what it can do:

```sh
aikido-dojo auth login --client-id AIK_CLIENT_xxxx
aikido-dojo auth status
```

`auth login` asks for the client secret without echoing it and keeps it in the OS keychain. The
client ID and the region go in `~/.config/aikido-dojo/config.json`. A workspace outside the EU
adds `--region us`, `au` or `me`.

In CI, set `AIKIDO_DOJO_CLIENT_ID` and `AIKIDO_DOJO_CLIENT_SECRET` instead, and nothing is
stored.

## Examples

```sh
aikido-dojo repo list --limit 5
aikido-dojo issue-group list --filter-code-repo-id 1 --jq '.[] | .title'
aikido-dojo repo current
aikido-dojo issue-group-note create 12 --note 'Accepted: only reachable in tests'
aikido-dojo team delete 7 --dry-run
aikido-dojo api GET /workspace
```

Commands are `<resource> <verb>`, with path parameters as arguments and query parameters as
flags. Output is JSON: indented on a terminal, compact in a pipe. `--jq` filters it and
`--ndjson` prints a list one item per line. `repo current` finds the Aikido repo of the git
checkout you are in, across every profile. `--dry-run` prints the request a write would send
and sends nothing. `api` calls any path under `/api/public/v1`.

## For agents

[AGENTS.md](AGENTS.md) is a primer for an agent to load. An agent finds commands with `search`,
learns one with `schema`, and runs in read-only mode:

```sh
aikido-dojo search 'ignore a finding'
aikido-dojo schema issue-group ignore
aikido-dojo --read-only team delete 7
```

Give an agent an API client with read scopes only, and set `AIKIDO_DOJO_READ_ONLY=1` so every
call that isn't a GET is refused. Destructive commands ask first on a terminal and, off one,
refuse unless `--yes` is passed.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | unexpected failure |
| 2 | usage or validation error |
| 3 | authentication failed |
| 4 | missing scope or forbidden |
| 5 | not found |
| 6 | rate-limited after retries |
| 7 | refused by read-only mode or the destructive-action guard |

An error is one JSON object on stderr: `{"error":{"code":…,"message":…,"hint":…}}`.

## License

Apache License 2.0; see [LICENSE](LICENSE). To report a vulnerability, see
[SECURITY.md](SECURITY.md).
