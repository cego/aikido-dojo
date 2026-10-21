# Bitbucket PR gating using Local Scanner

The Aikido Local Scanner can be used to enforce security gates in your CI pipeline.

{% hint style="info" %}
You can also use PR gating with your cloud-connected workspaces (GitHub, BitBucket, etc). Check out our CI [Integration section](https://help.aikido.dev/section/ci-integrations/sg3q6UrIf4qE).
{% endhint %}

PR gating ensures that new code meets your security standards before it is merged into the default branch. It scans only the changes introduced in the pull request. If the scan detects new issues that meet or exceed your configured severity threshold, the CI pipeline fails.

This prevents new vulnerabilities from being introduced while allowing existing findings to be addressed separately.

{% hint style="info" %}
We also support [release gating](/container-image-scanning/local-image-scanning/release-gating-for-container-images-using-local-image-scanner.md) for teams that want to enforce checks at release time.
{% endhint %}

To enable PR gating, add the `--fail-on <severity>` option to select your preferred severity level. Then, add the `--gating-mode pr` option to signify that you wish to perform PR gating. You must also specify the base (`--base-commit-id <commit-id>`) and head commit (`--head-commit-id <commit-id>`). If there a scan was previously performed on the base commit id, the scan results will be compared to those. If not, the results will be compared to the most recent scan on your default branch.

For general information about setting up the Local Scanner in a Bitbucket environment, check out [this article](/code-scanning/local-code-scanning/bitbucket-pipeline-setup-for-local-code-scanning.md).

Example `bitbucket-pipelines.yml` for PR gating:

```yaml
image: aikidosecurity/local-scanner:latest

pipelines:
  pull-requests:
    "**":   # this applies the pipeline to all pull requests
      - step:
          name: Aikido PR Gating Scan
          script:
            - >- 
              aikido-local-scanner scan ./
              --apikey $AIKIDO_API_KEY
              --repositoryname $BITBUCKET_REPO_SLUG
              --branchname $BITBUCKET_BRANCH
              --gating-mode pr
              --fail-on critical
              --base-commit-id $BITBUCKET_PR_DESTINATION_COMMIT
              --head-commit-id $BITBUCKET_COMMIT
```
