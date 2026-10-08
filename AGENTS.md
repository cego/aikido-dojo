# aikido-dojo for agents

aikido-dojo calls the Aikido Security public API. Data goes to stdout as JSON. An error is one
JSON object on stderr, with a `code`, a `message` and a `hint` that says what to do next.

## Commands

A command is `aikido-dojo <resource> <verb> [arguments] [--flags]`:

- path parameters are positional arguments;
- query parameters are flags named after the spec's, with `-` for `_`;
- body fields with a flag of their own take it; the rest go in `--body '<json>'` or
  `--body-file <path|->`;
- credentials go only in `--body-file`.

Find a command and learn it without calling the API:

```sh
aikido-dojo search 'list open issues for a repository'
aikido-dojo schema issue-group list
```

`schema` prints JSON Schema for the command's `args`, `flags`, `body` and `response`. A list
command pages by itself and prints one array; `--limit N` stops early.

Find the Aikido repo of the current git checkout, and the profile to use with it:

```sh
aikido-dojo repo current --jq '.[0].repo.id'
```

When no command fits, `aikido-dojo api GET <path>` calls any path under `/api/public/v1`.

## Output

- `--jq '<expr>'` filters any command's JSON. Strings print without quotes.
- `--ndjson` prints a list one item per line, as it arrives.

## Safety

- `--read-only`, or `AIKIDO_DOJO_READ_ONLY=1`, refuses every call that isn't a GET.
- `--dry-run` on a write prints the request it would send, with secrets redacted, and sends
  nothing.
- A destructive command, such as a delete, deactivate or rotate, asks first on a terminal. Off
  one, it refuses unless `--yes` is passed.
- Exit 1 with code `output_failed` after a write means the write succeeded. Don't repeat it.

## Exit codes

0 success, 1 unexpected, 2 usage or validation, 3 authentication failed, 4 missing scope or
forbidden, 5 not found, 6 rate-limited after retries, 7 refused by read-only mode or the
destructive-action guard.
