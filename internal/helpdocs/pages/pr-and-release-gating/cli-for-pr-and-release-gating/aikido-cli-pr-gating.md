# Aikido CLI: PR Gating

The Aikido Security CLI client allows you to integrate Aikido Security scans into custom CI pipelines such as Jenkins or CircleCI. It helps ensure that security scans are part of your build process.

### Use Cases <a href="#use-cases" id="use-cases"></a>

* **Pull Request Checks**: Automatically add pass/fail checks to pull requests.

### Step-by-Step Guide <a href="#step-by-step-guide" id="step-by-step-guide"></a>

**Step 1**: Install the [`@aikidosec/ci-api-client`](https://www.npmjs.com/package/@aikidosec/ci-api-client) package globally.

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

Flags for this command are documented on the [`@aikidosec/ci-api-client`](https://www.npmjs.com/package/@aikidosec/ci-api-client) package. Print them locally with:

```shell
aikido-api-client help scan
```

**Step 4**: Upload custom test results (only supported for the `scan` command)\
You can optionally upload custom test results using the following command:

```
$ aikido-api-client upload --repository-id <repository_id> --type checkov --file <path_to_payload_file>
```
