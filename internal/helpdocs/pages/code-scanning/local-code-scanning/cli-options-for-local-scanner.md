# CLI Options for local code scanning

The Aikido Security Local Scanner is a tool that enables you to perform Aikido Security scans within your environment, ensuring your code never leaves your premises. Below you can find the options that can be passed when running the local scanner.

### Repository scanning <a href="#repository-scanning" id="repository-scanning"></a>

```markdown
Usage: aikido-local-scanner scan [options] <path>

Run a scan.

Arguments:
  path                                     The path you want to scan.

Options:
  --apikey <apikey>                        Apikey to send scanning results to Aikdo.  (env: AIKIDO_API_KEY)
  --repositoryname <repositoryname>        Repo name to create or send results to.
  --branchname <branchname>                Branch name that is being scanned.
  --tmpdirectory <tmpdirectory>            Temporary directory to use during scanning. (default: "./.aikidotmp")
  --debug                                  Add additional debug information to command output.
  --disable-artifact-scanning              Disable artifact scanning. Use to speed up scanning at the cost of not scanning artifacts such as .jar files.
  --secrets-scanning-full-git-history      Enable scanning the full Git history for secrets.
  --scan-types [types...]                  Specify which types of scans should be executed. This will overwrite the --scanners flag (cf below). (choices: "code", "dependencies", "iac", "secrets", default: [])
  --exclude <exclude_path>                 Specify a file or folder path that should be excluded from the scan. This option may be specified multiple times. (default: [])
  --fail-on <severity>                     Runs scanner in gating mode and fails on the given severity or higher. (choices: "low", "medium", "high", "critical")
  --gating-mode <mode>                     Indicate whether the scanner should run in release or PR gating mode. Release gating mode scans your main branch and waits to see if there are any issues that should prevent release. Pull request mode, Aikido seeks ONLY new vulnerabilities introduced in a branch (scans the diff). You must supply a base and head commit for the comparison to work. Should be combined with the --fail-on flag (choices: "release", "pr", default: "release")
  --sla-mode                               In release gating mode, only fail the gate on issues that are out of SLA.
  --base-commit-id <commit-id>             Base commit id, this is the commit that Aikido will compare against to determine if a finding is new. Only used for PR gating mode.
  --head-commit-id <commit-id>             Head commit id, the commit that is being scanned. Required for PR gating mode.
  --gating-result-output <output>          JSON file to write issues to when running in gating mode only
  --no-fail-on-timeout                     Do not fail the process in case the scan result polling times out (gating mode only)
  --max-polling-attempts <amount>          Amount of times to poll for scan results, increase this if the default value of 20 is not enough (gating mode only)
  --linked-team-name <name>                Team name to link the repository to. Specify this option multiple times to link to multiple teams. (default: [])
  --no-snippets                            Use this mode to not share any code snippets with Aikido.
  --checkov-skip-extension <extension>     Specify an extension to be skipped by Checkov scans specifically. This option may be specified multiple times. Example: If you want to ignore JSON files, pass .json as the value for this option. (default: [])
  --no-lockfiles-cache                     Do not allow caching of dependency- and lockfiles to enable automated rescans
  --scan-timeout <timeout>                 Timeout in milliseconds for each scan (defaults to 900000ms).
  --force-create-repository-for-branch     Create a new repository in Aikido per branch.
  --enable-proxy                           Uses HTTPS_PROXY environment variable to proxy requests.
  --ca-bundle <path>                       Path to a PEM file containing custom root CA(s) to trust when using --enable-proxy. (env: AIKIDO_CA_BUNDLE)
  --include-dev-deps                       Enable scanning of development dependencies (e.g., devDependencies in package.json, etc.)
  -h, --help                               display help for command
```

For image scanning CLI options, see [here](/container-image-scanning/local-image-scanning/cli-options-for-local-image-scanning.md).

## Common use cases <a href="#common-use-cases" id="common-use-cases"></a>

Below you can find some example configurations of common use cases.

**I want to run a scan on a repository**

> The api key can also be passed as an environment variable *AIKIDO\_API\_KEY*

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main
```

**I want to run only SCA scanning on my repository**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--scan-types dependencies
```

**I want to run only secrets scanning on my repository**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--scan-types secrets
```

**I want to run only code scanning on my repository**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--scan-types code
```

**I want to run only IaC scanning on my repository**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--scan-types iac
```

**I want to exclude my /staging and /development folder from being scanned**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--exclude staging 
--exclude development
```

**I want to run only secrets and code scanning and I want to exclude my /staging folder.**

```shellscript
aikido-local-scanner scan ./ 
--apikey AIK_CI_xxx 
--repositoryname DemoApp 
--branchname main 
--scan-types secrets code
--exclude staging
```
