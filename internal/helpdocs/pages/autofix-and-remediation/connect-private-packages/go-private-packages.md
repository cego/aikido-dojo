# Go - Private Packages

For Aikido to update dependencies that include private packages, it needs access to your private registries so it can generate accurate lockfile updates. Follow the steps below to configure private Go packages hosted on GitHub or GitLab.

1. Navigate to the [autofix settings](https://app.aikido.dev/issues/fix/settings) page and click on "*Connect Registry*"
2. Select "*Go*" as your Package Manager

3. In the modal below, enter the private registry host, this is the hostname with the path to the repo where the package is being hosted. Aikido will expose this as the `GOPRIVATE` env variable during the update process.
4. eg for GitHub: `github.com/[ORG]/[REPO]`
5. eg for Gitlab: `gitlab.com/[ORG]/[REPO]`

Aikido uses this value to set the `GOPRIVATE` environment variable, so it can be a list of comma-separated glob patterns of module path prefixes that should be considered private (eg. `github.com/myorg/*` or `gitlab.com/myorg/repo1,gitlab.com/myorg/repo2`.

1. Next create a personal access token in either GitHub or Gitlab with read-only access for that repository and it to the token field

5. Click on "*Connect Registry*"
