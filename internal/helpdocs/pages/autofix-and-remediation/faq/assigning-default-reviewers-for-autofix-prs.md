# Assigning Default Reviewers for AutoFix PRs

Aikido doesn't currently support configuring default reviewers for AutoFix pull requests. AutoFix opens the PR in your SCM, but it doesn't assign reviewers from an Aikido setting.

The practical path today is to enforce reviewers on the repository itself with **CODEOWNERS** and **required reviews**.

## Workaround for GitHub

Use GitHub's native CODEOWNERS support and branch protection so the right people are requested (and required) whenever AutoFix changes files they own.

{% stepper %}
{% step %}

### Add a CODEOWNERS file

Create a `CODEOWNERS` file in `.github/`, the repository root, or `docs/`. Map the paths AutoFix typically changes (for example lockfiles and package manifests) to the users or teams that should review those PRs:

```
# Dependency files commonly changed by AutoFix
package.json @your-org/security-team
package-lock.json @your-org/security-team
yarn.lock @your-org/security-team
pnpm-lock.yaml @your-org/security-team
requirements.txt @your-org/security-team
poetry.lock @your-org/security-team
Cargo.lock @your-org/security-team
go.sum @your-org/security-team
```

Adjust the paths and owners to match your repo layout. See [About code owners](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) for syntax details.
{% endstep %}

{% step %}

### Require review from code owners

In your repository's branch protection rules (or rulesets) for the target branch:

1. Require a pull request before merging
2. Enable **Require review from Code Owners**

With this enabled, GitHub requests reviews from the matching CODEOWNERS and blocks merges until those reviews are satisfied. AutoFix PRs follow the same rules as any other PR.
{% endstep %}
{% endstepper %}

{% hint style="info" %}
GitHub can also automatically request reviews from CODEOWNERS when a PR touches owned files, even before you require code owner reviews for merging. Required code owner reviews is what makes approval mandatory.
{% endhint %}

## GitLab

On GitLab, Aikido AutoFix adds CODEOWNERS as reviewers on merge requests when there are fewer than 10 CODEOWNERS for the changed files. You can turn that off with `disable_codeowners_as_reviewers` in your repo's `.aikido` file. See [Disabling CODEOWNERS as reviewers](/autofix-and-remediation/configure/autofix-pr-configuration.md#disabling-codeowners-as-reviewers).

To require those reviews before merge, configure approval rules or CODEOWNER approvals in GitLab itself.

## BitBucket

AutoFix adds CODEOWNERS as reviewers on BitBucket based on the `.bitbucket/CODEOWNERS` file when there are fewer than 10 CODEOWNERS for the changed files.

This requires a Workspace Access Token with `account:read` permissions.

Supported owner identifiers include email, team and username. Workspace groups are currently not supported, BitBucket advises to use teams instead.

You can turn this off with `disable_codeowners_as_reviewers` in your repo's `.aikido` file. See [Disabling CODEOWNERS as reviewers](/autofix-and-remediation/configure/autofix-pr-configuration.md#disabling-codeowners-as-reviewers).
