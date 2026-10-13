# AWS CodeArtifact - Private Packages

You can allow Aikido's AutoFix to connect to registries hosting private packages on AWS CodeArtifact.

## Configuration in Aikido

Aikido provides a CloudFormation template that creates a new role in your AWS account. The permissions for this role are limited to read-only CodeArtifact accesss. This works by giving the Aikido AWS account a trust relationship with the newly created role in your account.

After creation of the role, Aikido only needs the specific ARN to get started, no access keys nor secrets.

In the [Aikido UI](https://app.aikido.dev/issues/fix/settings), go to AutoFix > Settings > Connect Registry,

follow the instructions in the modal:

* Create the IAM Role & Policy
* Paste the Role ARN from the CloudFormation output
* Fill in the domain of the CodeArtifact repositories as seen in the AWS console (eg. mydomain)
* Select the AWS region where the CodeArtifact domain is hosted (eg. us-east-1)

Aikido AutoFix will set a `CODEARTIFACT_AUTH_TOKEN` environment variable during the process. This environment variable can then be used by the package manager to authenticate with the CodeArtifact registry. See below for registry specific config which is required.

## NPM & PNPM <a href="#npm--pnpm" id="npm--pnpm"></a>

A basic `.npmrc` configuration must be present in the repository (or Aikido configuration) where the private package is being installed. It should contain the following information, with the placeholders replaced with your information:

```
@pied-piper:registry=https://mydomain-123456789.d.codeartifact.us-east-1.amazonaws.com/npm/pied-piper
//mydomain-123456789.d.codeartifact.us-east-1.amazonaws.com/npm/pied-piper/:_authToken=${CODEARTIFACT_AUTH_TOKEN}
```

During the AutoFix, the `${CODEARTIFACT_AUTH_TOKEN}` will be replaced by an actual auth token.

## Poetry

You need to set an additional environment variable to tell Poetry to use the CodeArtifact Auth Token: set `POETRY_HTTP_BASIC_[SOURCE]_USERNAME` to `aws` . Make sure to replace the `[SOURCE]` with the source name (in upper case) from your `pyproject.toml` , eg:

```toml
[[tool.poetry.source]]
name = "artifact"
url = "https://repo-1234567890.d.codeartifact.eu-west-1.amazonaws.com/pypi/poetry/simple/"
```

The name of the source in the example is `artifact`. So the `POETRY_HTTP_BASIC_ARTIFACT_USERNAME` environment variable needs to be set to `aws`.

By setting `POETRY_HTTP_BASIC_[SOURCE]_USERNAME` to `aws` , AutoFix will automatically populate `POETRY_HTTP_BASIC_[SOURCE]_PASSWORD` which will provide Poetry with access to the AWS CodeArtifact repository.

## UV

You need to set an additional environment variable to tell UV to use the CodeArtifact Auth Token: set `UV_INDEX_[INDEX]_USERNAME` to `aws` . Make sure to replace the `[INDEX]` with the index name (in upper case) from your `pyproject.toml` , eg:

```toml
[[tool.uv.index]]
name = "artifact"
url = "https://repo-1234567890.d.codeartifact.eu-west-1.amazonaws.com/pypi/uv-example/simple/"
```

The name of the index in the example is `artifact`. So the `UV_INDEX_ARTIFACT_USERNAME` environment variable needs to be set to `aws`.

By setting `UV_INDEX_[INDEX]_USERNAME` to `aws` , AutoFix will automatically populate `UV_INDEX_[INDEX]_PASSWORD` which will provide UV with access to the AWS CodeArtifact repository.

## Maven

After setting up the AWS CodeArtifact connection, you can use the `CODEARTIFACT_AUTH_TOKEN` environment variable in `settings.xml` . The settings.xml file can be set in the Maven Private Registry configuration, eg:

{% code title="settings.xml" %}

```
<servers>
    <server>
        <id>codeartifact</id>
        <username>aws</username>
        <password>${env.CODEARTIFACT_AUTH_TOKEN}</password>
    </server>
</servers>
```

{% endcode %}

## How to check your configuration

When creating a dependency AutoFix, you should see the following messages in the progress modal:

* Set 3 (or more) environment variables
* AWS CodeArtifact authentication token successfully set

If your AWS credentials are incorrect or are lacking permissions you will see a 'Failed to authenticate to AWS CodeArtifact' warning message:

If there is no message about AWS CodeArtifact, this means the configuration is not complete.

## Using AWS access keys (not recommended)

To enable the connection with AWS CodeArtifact, you need to provide the following 4 environment variables:

* `AWS_CODE_ARTIFACT_ACCESS_KEY`: an access key for a user or role which has the `AWSCodeArtifactReadOnlyAccess` permission
* `AWS_CODE_ARTIFACT_SECRET_KEY`: a secret key for a user or role which has the `AWSCodeArtifactReadOnlyAccess` permission
* `AWS_CODE_ARTIFACT_DOMAIN`: the domain of the CodeArtifact repositories as seen in the AWS console (eg. mydomain)
* `AWS_CODE_ARTIFACT_REGION`: the AWS region where the CodeArtifact domain is hosted (eg. us-east-1)

When these 4 environment variables are set, Aikido AutoFix will set a `CODEARTIFACT_AUTH_TOKEN` environment variable during the process. This environment variable can then be used by the package manager to authenticate with the CodeArtifact registry. See below for registry specific config which is required.

These environment variables can be set in Autofix > Settings > Connect Registry > [Set Environment variables](/autofix-and-remediation/connect-private-packages/custom-config-private-packages.md).
