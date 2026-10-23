# Ephemeral Code Scanning

To run a one-time stand-alone SAST & dependencies scan that writes results to a file run a command like:

```
./aikido-local-scan ephemeral-scan path_to_your_repo --apikey AIK_CI_xxx
```

This will write findings to an `aikido-scan-results.json` file.\
\
All CLI options:

```
Usage: aikido-local-scanner ephemeral-scan [options] <path>

Run a scan that outputs results to a JSON file. BETA: supports SAST and dependency scanning.

Arguments:
  path                               The path you want to scan.

Options:
  --apikey <apikey>                  Apikey to process scan results. (env: AIKIDO_API_KEY)
  --repositoryname <repositoryname>  Aikido repository name to validate the scan results against. (default: "")
  --tmpdirectory <tmpdirectory>      Temporary directory to use during scanning. (default: "./.aikidotmp-1791275799502")
  --debug                            Add additional debug information to command output.
  --exclude <exclude_path>           Specify a file or folder path that should be excluded from the scan. This option may be specified multiple times. (default: [])
  --result-output <output>           JSON file to write issues to. (default: "./aikido-scan-results.json")
  --fail-on <severity>               Exit with a non-zero status code when the scan produces findings at or above the given severity. (choices: "low", "medium", "high", "critical")
  --scan-types [types...]            Specify which types of scans should be executed. (choices: "code", "dependencies", default: ["code","dependencies"])
  --scan-timeout <timeout>           Timeout in milliseconds for each scan (defaults to 900000ms).
  --enable-proxy                     Uses HTTPS_PROXY environment variable to proxy requests
  --ca-bundle <path>                 Path to a PEM file containing custom root CA(s) to trust when using --enable-proxy. (env: AIKIDO_CA_BUNDLE)
  --disable-telemetry                Disable all error reporting (telemetry) to Aikido. (env: AIKIDO_DISABLE_TELEMETRY)
  --include-dev-deps                 Enable scanning of development dependencies (e.g., devDependencies in package.json, etc.)
  -h, --help                         display help for command

```

## Common use cases <a href="#common-use-cases" id="common-use-cases"></a>

Below you can find some example configurations of common use cases.\
\
I want to run only dependency scanning

```
aikido-local-scanner ephemeral-scan ./ 
--apikey AIK_CI_xxx 
--scan-types dependencies
```

I want to run only code scanning

```
aikido-local-scanner ephemeral-scan ./ 
--apikey AIK_CI_xxx 
--scan-types code
```
