# Ephemeral Image Scanning

To run a one-time stand-alone image scan that writes results to a file run a command like:

```
./aikido-local-scan ephemeral-image-scan your-image-name --apikey AIK_CI_xxx
```

This will write findings to an `aikido-image-scan-results.json` file.\
\
All CLI options:

```
Usage: aikido-local-scanner ephemeral-image-scan [options] <image>

Run an ephemeral image scan that outputs results to a JSON file without persisting findings in Aikido.

Arguments:
  image                     The image you want to scan.

Options:
  --apikey <apikey>                 Apikey to send scanning results to Aikido. (env: AIKIDO_API_KEY)
  --platform <platform>             Set platform (to pull arm64 image on a amd64 system for example)
  --image-name <name>               Specify a name for the scanned image. This overwrites the default behaviour of deducting the image name from the <image> argument.
  --result-output <output>          JSON file to write issues to. (default: "./aikido-image-scan-results.json")
  --output-cyclonedx-json <output>  JSON file to write CycloneDX JSON SBOM to.
  --fail-on <severity>              Exit with a non-zero status code when the scan produces findings at or above the given severity. (choices: "low", "medium", "high", "critical")
  --include-unfixed                 Include CVEs that have no fix available yet (by default only fixable findings are reported).
  --debug                           Add additional debug information to command output.
  --scan-timeout <timeout>          Timeout in milliseconds for each scan (defaults to 900000ms).
  --enable-proxy                    Uses HTTPS_PROXY environment variable to proxy requests
  --ca-bundle <path>                Path to a PEM file containing custom root CA(s) to trust when using --enable-proxy. (env: AIKIDO_CA_BUNDLE)
  -h, --help                        display help for command
```
