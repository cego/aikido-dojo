# PR Gating for Code Using Local Scanner

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

For implementation details see:

* [GitLab Self Managed](/code-scanning/local-code-scanning/pr-gating/gitlab-pr-gating-using-local-scanner.md)
* [GitHub](/code-scanning/local-code-scanning/pr-gating/github-pr-gating-using-local-scanner.md)
* [BitBucket](/code-scanning/local-code-scanning/pr-gating/bitbucket-pr-gating-using-local-scanner.md)
* [Azure DevOps Server](/code-scanning/local-code-scanning/pr-gating/azure-devops-pr-gating-using-local-scanner.md)
* [Jenkins](/code-scanning/local-code-scanning/pr-gating/jenkins-pr-gating-using-local-scanner.md)
