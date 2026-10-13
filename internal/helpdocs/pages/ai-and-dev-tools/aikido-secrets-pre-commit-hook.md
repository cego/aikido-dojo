# Secrets Pre-Commit Hook

Aikido prevents secrets at three points before they reach production: in the IDE, at pre-commit, and at PR gating. Our pre-commit githook scans your staged code for secrets, passwords and API keys. It stops sensitive data from ever reaching your repository, which reduces the risk of leaks and accidental exposure.

## Installation

{% tabs %}
{% tab title="From IDE" %}
When the [Aikido IDE plugin](broken://pages/TtCgRvpjJR8sPsXdszal) is installed you can use the Aikido Expansion Packs to install the pre commit hook with one click.

[Learn more in the Expansion Packs docs.](/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs.md)
{% endtab %}

{% tab title="Mac/Linux" %}
To install the Aikido Secrets pre-commit hook for all git repositories, run:

```shellscript
curl -fsSL https://raw.githubusercontent.com/AikidoSec/pre-commit/6cc79e039ee78b206520f143d618a44665c904b3/installation-samples/install-global/install-aikido-hook.sh | bash
```

This will download the Aikido pre-commit scanner used for secrets detection and install the pre-commit hook script in the global hooks directory.
{% endtab %}

{% tab title="Windows" %}
To install the Aikido Secrets pre-commit hook for all git repositories, run the following in PowerShell:

```powershell
iex (iwr "https://raw.githubusercontent.com/AikidoSec/pre-commit/6cc79e039ee78b206520f143d618a44665c904b3/installation-samples/install-global/install-aikido-hook.ps1" -UseBasicParsing)
```

This will download the Aikido pre-commit scanner used for secrets detection and install the pre-commit hook script in the global hooks directory.
{% endtab %}

{% tab title="Pre-Commit Framework" %}
If you're already using the [pre-commit](https://pre-commit.com/) framework, add this to your `.pre-commit-config.yaml`:

```
repos:
  - repo: https://github.com/AikidoSec/pre-commit
    rev: main  # or pin to a specific commit
    hooks:
      - id: aikido-local-scanner
```

Then install the hooks:

```
pre-commit install
```

**Note:** The `aikido-local-scanner` binary must be installed separately. Run the global installation script first:

**macOS/Linux:**

```
curl -fsSL https://raw.githubusercontent.com/AikidoSec/pre-commit/2235fe0536f9135aa561ce108702fac708b38977/installation-samples/install-global/install-aikido-hook.sh | bash -s -- --download-only
```

**Windows (PowerShell):**

```
irm https://raw.githubusercontent.com/AikidoSec/pre-commit/2235fe0536f9135aa561ce108702fac708b38977/installation-samples/install-global/install-aikido-hook.ps1 | % { iex \"& { $_ } -DownloadOnly\" }
```

This installs the scanner to `~/.local/bin/aikido-local-scanner`.
{% endtab %}
{% endtabs %}

The [source of the script and more information about its workings](https://github.com/AikidoSec/pre-commit) are available on Github.

## Testing the pre-commit hook

To test the pre-commit hook after you've set it up, create a `sample.js` file in a repository:

```javascript
const password = "eRwjQKVUSRX7uYV017B0cRHVKv45Gv8G"
```

Add this file to your staged changes. If you try commit this file, the pre-commit hook will run and block the commit with the following message:

```log
Detected secrets in staged files!
Secret #1:
  File: sample.js
  Line: 1
  Secret: password = "****************************Gv8G"
  Description: Detected a Generic API Key, potentially exposing access to various services and sensitive operations.
```

## Skipping a specific secret

To skip a [specific secret from being flagged](/code-scanning/scanning-practices/ignoring-secrets-via-code-comments.md), add a comment on the line of the detected secret:

```javascript
const password = "eRwjQKVUSRX7uYV017B0cRHVKv45Gv8G" // gitleaks:allow
```

## Disable the Aikido Secrets pre-commit scan

Temporarily bypass pre-commit hooks for a single commit

```sh
git commit --no-verify
```

Temporarily bypass the Aikido Secrets pre-commit hook for a single commit

```bash
AIKIDO_SKIP_PRE_COMMIT=1 git commit
```

## Uninstall

Use the uninstall script or follow the step below to manually uninstall the hook.

{% tabs %}
{% tab title="Mac/Lunix" %}
If you've installed the Aikido pre-commit hook using the install script and want to uninstall, run:

```bash
curl -fsSL https://raw.githubusercontent.com/AikidoSec/pre-commit/e6f541e65378dd30f3f320628000f837cfba0ec4/installation-samples/install-global/uninstall-aikido-hook.sh | bash
```

{% endtab %}

{% tab title="Windows" %}
If you've installed the Aikido pre-commit hook using the install script and want to uninstall, run:

```powershell
iex (iwr "https://raw.githubusercontent.com/AikidoSec/pre-commit/e6f541e65378dd30f3f320628000f837cfba0ec4/installation-samples/install-global/uninstall-aikido-hook.ps1" -UseBasicParsing)
```

{% endtab %}
{% endtabs %}

### Manual uninstall of global pre-commit hooks

This fully removes all global Git hooks and the Aikido binary.

1. Remove the global hooks directory:
   * Unix/Linux/macOS: `rm -rf ~/.git-hooks`
   * Windows: `Remove-Item -Recurse -Force $env:USERPROFILE\.git-hooks`
2. Reset Git hooks path: `git config --global --unset core.hooksPath`
3. Optionally remove the binary:
   * Unix/Linux/macOS: `rm ~/.local/bin/aikido-local-scanner`
   * Windows: `Remove-Item $env:USERPROFILE\.local\bin\aikido-local-scanner.exe`

### Manual uninstall of only Aikido Git Hook

If you already had your own global Git hooks and want to keep them, do not delete the hooks directory.

Instead:

1. Open the pre-commit file in your global hooks directory (for example \~/.git-hooks/pre-commit).
2. Remove only the lines that invoke aikido-local-scanner or are clearly marked as added by Aikido.
3. Save the file.

Git will keep using your existing hooks, without running Aikido Secrets.

## Related articles

* [Ignoring secrets via code comments](/code-scanning/scanning-practices/ignoring-secrets-via-code-comments.md)
* [Ignoring files using .aikido file](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md)
