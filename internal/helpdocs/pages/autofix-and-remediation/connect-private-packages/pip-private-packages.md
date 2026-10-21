# PIP - Private Packages

For Aikido to update dependencies that include private packages, it needs access to your private packages so it can generate accurate requirement file updates. You can provide private git repository configuration in Aikido for this.

## Private PyPI registry <a href="#same-scm-organization" id="same-scm-organization"></a>

You can configure Aikido to authenticate with your private PyPI registry when updating the PIP dependencies by following the steps below:

1. Go to your account's settings page for the autofixer in Aikido, [here](https://app.aikido.dev/issues/fix/settings).
2. Click on "Connect Registry", the private registry modal will now be shown

3. Select PIP and provide the index URL to the PyPI registry including the username and password in the following format:

```
https://<username>:<password>@domain.com/repo
```

4. Click "Connect Registry" to save the configuration.

If the private registry is hosted on **GitHub**, **GitLab**, **BitBucket** or **Azure DevOps**, consult the sections below.

## Private Git repository <a href="#configuration-in-aikido---private-rubygems-repository" id="configuration-in-aikido---private-rubygems-repository"></a>

### Same SCM organization <a href="#same-scm-organization" id="same-scm-organization"></a>

If the private Git repository for the pip is in the same organization as the code being autofixed, make sure your [AutoFix access token](https://app.aikido.dev/issues/fix/settings) has access to the private git repository.

### Different SCM organization <a href="#different-scm-organization" id="different-scm-organization"></a>

You can configure Aikido to authenticate with your private repositories when updating the dependencies by following the steps below:

1. Go to your account's settings page for the autofixer in Aikido, [here](https://app.aikido.dev/settings/integrations/autofix).
2. Click on "Connect Registry", the configuration modal will now be shown

3. Select "Environment Variables"

4. In the example above, we show a possible setup for a private GitHub repository. The environment variable key is `PIP_GITHUB__COM`, and the value is a GitHub access token with repository access.

| Provider     | Environment Variable Key | Value Format        |
| ------------ | ------------------------ | ------------------- |
| GitHub       | `PIP_GITHUB__COM`        | `your-access-token` |
| GitLab       | `PIP_GITLAB__COM`        | `your-access-token` |
| Bitbucket    | `PIP_BITBUCKET__ORG`     | `your-access-token` |
| Azure DevOps | `PIP_DEV__AZURE__COM`    | `your-access-token` |

5. Click "Connect Registry" to save the configuration.
