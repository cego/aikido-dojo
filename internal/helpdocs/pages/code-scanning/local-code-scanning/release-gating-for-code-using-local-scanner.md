# Release Gating for Code Using Local Scanner

The Aikido Local Scanner can enforce security gates at different stages of your CI pipeline, including pull requests and releases.

{% hint style="info" %}
You can also use PR gating with your cloud-connected workspaces (GitHub, BitBucket, etc). Check out our CI [Integration section](https://help.aikido.dev/section/ci-integrations/sg3q6UrIf4qE).
{% endhint %}

Release gating scans your code before it is published or deployed. Unlike [PR gating](/container-image-scanning/local-image-scanning/pr-gating-for-container-images-using-local-image-scanner.md), which scans only the changes introduced in a pull request, release gating evaluates the entire build. If the scan finds issues that meet or exceed your configured severity threshold, the CI pipeline fails and the release is blocked.

This ensures that the final code meets your security standards before it reaches production.

To enable release gating add the `--fail-on <severity>` option to the scan command.

`aikido-local-scanner scan ./ --apikey $AIKIDO_LOCAL_SCANNER_TOKEN --repositoryname $CI_PROJECT_NAME --branchname main --fail-on critical`

## Examples <a href="#examples" id="examples"></a>

### GitLab Self Managed <a href="#gitlab-self-managed" id="gitlab-self-managed"></a>

For general information about setting up the Local Scanner in a GitLab environment, check out [this article](/code-scanning/local-code-scanning/gitlab-self-managed-setup-for-local-code-scanning.md).

Example `.gitlab-ci.yml` for release gating:

```yaml
default:
  image:
    name: aikidosecurity/local-scanner:latest
    entrypoint: [""]

run_aikido_selfscanner:
  script:
  - aikido-local-scanner scan ./ --apikey $AIKIDO_LOCAL_SCANNER_TOKEN --repositoryname $CI_PROJECT_NAME --branchname main --fail-on critical
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
```

### GitHub <a href="#github" id="github"></a>

For general information about setting up the Local Scanner in a GitHub environment, check out [this article](/code-scanning/local-code-scanning/github-action-setup-for-local-code-scanning.md).\
Example `.github/workflows/aikido-scan.yml` for release gating:

```yaml
on:
  push:
    branches:
      - main

name: Aikido Scan
jobs:
  aikido-local-scan-repo:
    runs-on: ubuntu-latest
    container:
      image: aikidosecurity/local-scanner:latest
    steps: 
      - uses: actions/checkout@v4 
        with: 
          token: ${{ secrets.GITHUB_TOKEN }} 
          path: my-repo 
      - name: Run scan
        run: aikido-local-scanner scan my-repo --apikey ${{ secrets.AIKIDO_API_KEY }} --repositoryname MyRepo --branchname main --fail-on critical
```
