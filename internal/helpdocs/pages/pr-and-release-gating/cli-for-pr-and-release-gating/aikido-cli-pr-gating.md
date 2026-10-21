# Aikido CLI: PR Gating

The Aikido Security CLI client allows you to integrate Aikido Security scans into custom CI pipelines such as Jenkins or CircleCI. It helps ensure that security scans are part of your build process.

### Use Cases <a href="#use-cases" id="use-cases"></a>

* **Pull Request Checks**: Automatically add pass/fail checks to pull requests.

### Step-by-Step Guide <a href="#step-by-step-guide" id="step-by-step-guide"></a>

**Step 1**: Install the Aikido CI client globally. See our [CI API client](https://www.npmjs.com/package/@aikidosec/ci-api-client).

```sh
# For npm users
$ npm install -g @aikidosec/ci-api-client

# For yarn users
$ yarn global add @aikidosec/ci-api-client
```

**Step 2**: Configure your Aikido API key

1. Obtain the API key from the [Continuous Integration Settings page](https://app.aikido.dev/settings/integrations/continuous-integration).
2. Set the API key globally:

```
$ aikido-api-client apikey <your-api-key-here>
```

**Step 3**: Run a pull request scan\
The `scan` command is used for **pull request gating**.

```
$ aikido-api-client scan <repository_id or repository_name> <base_commit_id> <head_commit_id>
```

**Options**:

* `--no-fail-on-dependency-scan`: Skip failing the build on new dependency vulnerabilities.
* `--fail-on-sast-scan`: Fail the scan if new SAST issues are detected.
* `--fail-on-iac-scan`: Fail the scan if new infrastructure-as-code issues are found.
* `--fail-on-secrets-scan`: Fail the scan if exposed secrets are detected.
* `--minimum-severity-level="CRITICAL"`: Minimum severity level for the scan to fail/block\
  Options: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL` (default: `CRITICAL`)
* `--poll-interval [interval]`: Set the poll interval for checking updated scan results (default: 10).

**Step 4**: Upload custom test results (only supported for the `scan` command)\
You can optionally upload custom test results using the following command:

```
$ aikido-api-client upload --repository-id <repository_id> --type checkov --file <path_to_payload_file>
```

### For more detailed options <a href="#for-more-detailed-options" id="for-more-detailed-options"></a>

```
$ aikido-api-client help
```
