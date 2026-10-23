# CLI Options for Local image scanning

Image Scanning CLI options:

```
Usage: aikido-local-scanner image-scan [options] <image>

Run an image scan.

Arguments:
  image                            The image you want to scan.

Options:
  --apikey <apikey>                Apikey to send scanning results to Aikido.  (env: AIKIDO_API_KEY)
  --platform <platform>            Set platform (to pull arm64 image on a amd64 system for example)
  --image-name <name>              Specify a name for the scanned image. This overwrites the default behaviour of deducting the image name from the <image> argument. Can be used to specify a descriptive name if the image name is not.
  --multi-arch <arch>              Run a separate scan for each given architecture (e.g. linux/amd64). Specify this option multiple times to scan multiple architectures. Each architecture is reported under a unique image name (<image-name>-<arch>), and when combined with --output-cyclonedx-json each architecture is written to its own file (<name>-<arch>.json). Not supported in gating mode and cannot be combined with --platform.
  --tag <tag>                      Specify a tag for the scanned image for reporting purposes. This overwrites the default behaviour of deducting the tag from the <image> or <image-name> argument.
  --debug                          Add additional debug information to command output.
  --fail-on <severity>             Runs scanner in gating mode and fails on the given severity or higher. (choices: "low", "medium", "high", "critical")
  --gating-mode <mode>             Indicate whether the scanner should run in release or PR gating mode. Release gating mode scans your main branch and waits to see if there are any issues that should prevent release, in pull request mode, Aikido seeks ONLY new vulnerabilities introduced in a branch. You must supply a base and head commit for the comparison to work. Should be combined with the --fail-on flag (choices: "release", "pr", default: "release")
  --sla-mode                       In release gating mode, only fail the gate on issues that are out of SLA.
  --base-commit-id <commit-id>     Base commit id, this is the commit that Aikido will compare against to determine if a finding is new. Only used for PR gating mode.
  --head-commit-id <commit-id>     Head commit id, the commit that is being scanned. Only used for PR gating mode.
  --gating-result-output <output>  JSON file to write issues to (when running in gating mode only)
  --no-fail-on-timeout             Do not fail the process in case the scan result polling times out (release gating mode only)
  --max-polling-attempts <amount>  Amount of times to poll for scan results, increase this if the default value of 20 is not enough (release gating mode only)
  --linked-team-name <name>        Team name to link the image to.
  --force-create-image-for-tag     Create a new image in Aikido with the tag instead of updating the tag on existing image.
  --scan-timeout <timeout>         Timeout in milliseconds for each scan (defaults to 900000ms).
  --output-cyclonedx-json <output> JSON file to write CycloneDX JSON SBOM to.
  --scan-javascript-lockfiles.     Also detect JavaScript packages from lockfiles (e.g. package-lock.json, yarn.lock) found inside the image, in addition to installed packages.
  --enable-proxy                   Uses HTTPS_PROXY environment variable to proxy requests.
  --ca-bundle <path>               Path to a PEM file containing custom root CA(s) to trust when using --enable-proxy. (env: AIKIDO_CA_BUNDLE)
  --disable-telemetry              Disable all error reporting (telemetry) to Aikido. (env: AIKIDO_DISABLE_TELEMETRY)
  --linked-repository-name <name>  Repository name to link the image to.
  -h, --help                       display help for command
```

For repository scanning CLI options, see [here](/code-scanning/local-code-scanning/cli-options-for-local-scanner.md).

## Common use cases <a href="#common-use-cases" id="common-use-cases"></a>

Below you can find some example configurations of common use cases.<br>

**I want to run an image scan**

```shellscript
aikido-local-scanner image-scan my-image
--apikey AIK_CI_xxx 
```

**I want to run an image scan and link the image to team "Cloud Team"**

```shellscript
aikido-local-scanner image-scan my-image
--apikey AIK_CI_xxx 
--linked-team-name "Cloud Team"
```

**I want to run a image scan and link the image to teams "Dev Team 1" and "Dev Team 2"**

```shellscript
aikido-local-scanner image-scan my-image 
--apikey AIK_CI_xxx 
--linked-team-name "Dev Team 1"
--linked-team-name "Dev Team 2"
```

**I want to run an image scan and output the SBOM to a JSON file**

```shellscript
aikido-local-scanner image-scan my-image
--apikey AIK_CI_xxx 
--output-cyclonedx-json sbom.json
```

**I want to run an image scan on a registry image for multiple architectures in one command. This will create a separate image in feed with results per architecture.**

```shellscript
aikido-local-scanner image-scan my-image
--apikey AIK_CI_xxx 
--multi-arch linux/amd64
--multi-arch linux/arm64
```
