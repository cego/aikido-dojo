# GitHub PR Gating

You can easily configure GitHub PR Gating via the Aikido interface. This functionality lets you block **newly introduced** issues for a threshold you decide and see which existing issues a pull request solves. These checks run **everytime** changes are committed to the feature branch. This doc focusses on managing configurations in bulk - multiple repos at the same time - without code.

### Use Cases <a href="#use-cases" id="use-cases"></a>

* **Bulk Repository Management**: Easily specify and manage configurations for multiple repositories at once through the Aikido interface.
* **Zero Code Integration**: Install the Aikido app on GitHub to manage checks without embedding any code, simplifying the setup process.
* **Cost Efficiency**: By managing GitHub Checks through Aikido, avoid using CI minutes on GitHub, leading to significant cost savings.

### Setting up GitHub CI <a href="#setting-up-github-ci" id="setting-up-github-ci"></a>

**Step 1.** Go to the [Integrations Page](https://app.aikido.dev/settings/integrations) and select GitHub in the CI gating section.

<div data-with-frame="true">[image: Connect CI tools (GitHub, GitLab, Bitbucket, Azure) to block risky code merges.]</div>

**Step 2.** Install the **Aikido PR Checks app** in GitHub. Make sure that you select the GitHub organisation that is currently being used in your workspace. Choose which repos that Aikido is allowed to access. We recommend giving access to all repos so these can easily be managed from within Aikido.

<div data-with-frame="true">[image: Aikido PR Checks installation prompt with logo and installation location query.]</div>

**Step 3.** Aikido redirects you to the [GitHub CI page](https://app.aikido.dev/settings/integrations/github/checks) with an overview of your repos. You can start configuring your repos. We recommend starting out with 1 repo to make sure everything works well.

<div data-with-frame="true"></div>

**Step 4. Select repos in bulk** and click **Setup PR Scans** in the floating bar on the bottom

<div data-with-frame="true"></div>

**Step 5.** This opens the configuration modal, where you choose the failure threshold and scans for the selected repositories.

To run the enabled scans without failing the GitHub PR check, select **Always Pass**.

<div data-with-frame="true">[image: Security scanning configuration settings for the "about-github" repository.]</div>

{% hint style="info" %}
If you've added new repositories after the initial setup, you'll need to configure those repos as well.

To apply a default configuration automatically, see [Default PR/MR gating configuration for new repositories](/pr-and-release-gating/aikido-ci-gating-functionality/default-pr-mr-gating-configuration-for-new-repositories.md).
{% endhint %}

### Ignore issues directly from PR comments <a href="#ignore-issues-directly-from-pr-comments" id="ignore-issues-directly-from-pr-comments"></a>

When Aikido posts an inline PR comment for a finding, you can ignore that issue directly from GitHub by replying to the comment with:

`@AikidoSec ignore: [your reason to ignore]`

Example:

`@AikidoSec ignore: This secret is used for tests only.`

This performs the same action as manually ignoring the issue in the Aikido platform:

* The issue is marked as ignored in Aikido.
* The ignore reason is stored.
* PR gating can turn green once all blocking issues are resolved or ignored.

{% hint style="info" %}
This reply ignores the finding only when you have the **Can ignore issues** permission. Without that permission, Aikido replies that the ignore is not allowed, and the finding stays open.
{% endhint %}

### Make failed PR checks block merges <a href="#make-failed-pr-checks-block-merges" id="make-failed-pr-checks-block-merges"></a>

Aikido can fail the PR check, but GitHub decides whether that failing check blocks a merge. Configure branch protection with required status checks:

* [GitHub Docs: About protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches#require-status-checks-before-merging)
* [GitHub Docs: Managing a branch protection rule](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/managing-a-branch-protection-rule)

### Adding Exceptions for specific repos <a href="#adding-exceptions-for-specific-repos" id="adding-exceptions-for-specific-repos"></a>

You might want to have 1 specific repo where the configuration slightly differs. You can easily add exceptions by clicking the triple dots on a repo item or just select 1 or more items and go through the **Configure Scans** process again.

<div data-with-frame="true">[image: Critical vulnerability status dashboard with configuration management options for projects.]</div>
