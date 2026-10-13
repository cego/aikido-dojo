# Code Quality in PRs

PR checks focus on **new code**, meaning the lines you add or change in a pull request. Aikido scans those changes automatically and leaves inline comments (or fails the gate) so quality issues don't make it into your main branch.

This guide walks you through setting up Code Quality in your PR / MR Checks. To scan **existing code** across a full repository and see category scores, open findings, and hotspots, see [Repository Scans](/code-quality/repository-scans.md).

## Prerequisites

Before you begin, ensure you have:

* An account connected to GitHub, Bitbucket, GitLab, or Azure DevOps Git
* Admin access to your Aikido account
* [Connected repositories to Aikido](/code-scanning/connect-your-source-code.md)

Code Quality isn't available for Local Scanning Accounts or Team Foundation Version Control (TFVC) repositories.

## Getting started <a href="#getting-started" id="getting-started"></a>

There are 3 steps in total to setup and run your first code quality steps. The Aikido UI shows clearly which steps to follow on the Code Quality page (accessible from sidebar). Below is a more extensive description.

{% stepper %}
{% step %}
**Enable Code Quality for repositories**

Choose which repositories should have Code Quality checks enabled:

1. Click **Configure Repositories** on the Code Quality Page to open the PR checks page. If you haven't configured the CI checks yet, you will be taken through the installation flow.
2. Select the repositories where you want to run Code Quality checks and enable **Code Quality Scanning**.
3. Select the lowest severity for Code Quality comments. Aikido adds comments for findings at that severity and above.
4. Enable **Fail check for open quality comments** to fail the check when Aikido finds Code Quality issues at or above the selected severity.

<div data-with-frame="true"></div>
{% endstep %}

{% step %}
**Select quality checks**

Configure which code quality rules to enforce on the Code Quality Checks page

1. Browse through the available checks organized by:
   * **Default checks** - Pre-configured best practices for multiple languages
   * **Custom checks** - Rules specific to your team (you can add these later)
2. Toggle checks on/off based on your team’s needs

<div data-with-frame="true"></div>
{% endstep %}

{% step %}
**Run your first scan**

Once you’ve configured repositories and checks:

1. Create a new pull request, or push a new commit to an existing pull request in one of the configured repositories. This will automatically trigger your first Code Quality scan.
2. Review the results in the **Activity** tab.
   {% endstep %}
   {% endstepper %}

When you fix a Code Quality issue, Aikido automatically resolves the comment. [Learn how Aikido manages comments after an issue is fixed](/pr-and-release-gating/aikido-ci-gating-functionality/autofix-suggestions-and-inline-commenting-for-pr-checks.md#what-happens-to-comments-after-you-fix-an-issue).

To fail the pull request when uploaded test coverage is below a percentage, see [PR Gating on Code Coverage](/code-quality/code-coverage/pr-gating-on-code-coverage.md). That setting is in the same **Code Quality Scanning** section.

## How severity is determined

Every Code Quality rule has a base severity, shown on the [**Checks** page](https://app.aikido.dev/code-quality/checks). When Aikido finds a violation, it starts with that base severity and adjusts it for the specific context. The resulting finding severity can be higher or lower than the rule's base severity.
