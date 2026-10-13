# Bitbucket PR Gating

Aikido Bitbucket PR gating runs security scan on each pull request (PRs). It focuses on **newly introduced** issues in the PR and shows which existing issues the PR solves.

Use it to catch risky changes before they land on your main branch. You decide what should fail the gate.

### Pick the right guide

Start with the page that matches your Bitbucket plan. Then configure repositories in the dashboard.

* **Bitbucket Free & Standard**
  * Follow: [Bitbucket Free & Standard: PR Scans Setup](/pr-and-release-gating/bitbucket-pr-gating/bitbucket-ci-pr-gating-via-aikido-dashboard.md)
* **Bitbucket Premium**
  * Use a dedicated Workspace Access Token
  * Follow: [Bitbucket Premium: PR Gating Setup](/pr-and-release-gating/bitbucket-pr-gating/bitbucket-ci-pr-gating-via-aikido-dashboard-with-a-workspace-access-token.md)
