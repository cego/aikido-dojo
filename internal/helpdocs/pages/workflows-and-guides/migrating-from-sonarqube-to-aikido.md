# Migrating from SonarQube to Aikido

Learn how to migrate from SonarQube or SonarCloud to Aikido.

Aikido helps you migrate from SonarQube and SonarCloud by bringing security findings, code quality checks, and PR gating into one workflow.

Use this guide to decide how to map your Sonar workflows to Aikido and make the switch without surprising your developers.

## Choose your migration path

You don't need the same migration path for every repository. Start with the path that matches how heavily your team depends on Sonar today.

| Path                      | Use this when                                                                   | What to do                                                                                                                                           |
| ------------------------- | ------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Move to Aikido directly   | Your Sonar setup mostly covers workflows that Aikido handles out of the box.    | Connect repositories to Aikido, configure gating and custom rules, then retire the matching Sonar checks.                                            |
| Run both temporarily      | You want to compare findings before changing developer workflows.               | Enable Aikido first, keep Sonar in CI for a short validation window, then remove overlapping Sonar checks once teams are comfortable.                |
| Review specialized checks | You rely on a Sonar-specific check or metric that's important to your workflow. | Move your Sonar workflows to Aikido, then share the specialized check with us so we can prioritize a replacement or suggest an alternative workflow. |

## Map Sonar concepts to Aikido

Use this map to decide where each Sonar workflow belongs in Aikido.

| Sonar concept                       | Aikido equivalent                                                                                                                                                                                                                                                     | Notes                                                                                                                                                |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Security rules and vulnerabilities  | [SAST and IaC scanning](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md)                                                                                                                                                   | Aikido focuses on security findings and removes non-security styling noise from SAST.                                                                |
| Dependency findings                 | [Open-source dependency scanning](/code-scanning/scanning-practices/support-for-dependency-scanning-by-language.md)                                                                                                                                                   | Aikido also uses reachability and exploit intelligence to help you prioritize what matters.                                                          |
| Secrets                             | [Live Secret Detection](/code-scanning/scanning-practices/live-secret-detection.md)                                                                                                                                                                                   | Aikido checks whether leaked secrets are live and adjusts severity based on risk.                                                                    |
| IaC rules                           | [SAST and IaC scanning](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md)                                                                                                                                                   | Aikido scans Terraform, CloudFormation, Docker, Pulumi, and other IaC files.                                                                         |
| Quality Gates                       | [PR Gating](/pr-and-release-gating/aikido-ci-gating-functionality.md) and [Release Gating](/pr-and-release-gating/cli-for-pr-and-release-gating/aikido-cli-release-and-pr-gating.md)                                                                                  | Aikido can fail pull requests or release checks based on new issues and severity thresholds.                                                         |
| Code coverage                       | [Code Coverage](/code-quality/code-coverage.md) and [PR Gating on Code Coverage](/code-quality/code-coverage/pr-gating-on-code-coverage.md)                                                                                                                           | Upload coverage reports from CI and review coverage in Aikido. Fail the pull request when coverage is below your threshold.                          |
| Pull request decoration             | [Inline PR comments and AutoFix suggestions](/pr-and-release-gating/aikido-ci-gating-functionality/autofix-suggestions-and-inline-commenting-for-pr-checks.md)                                                                                                        | Aikido comments on pull requests and links developers to the related scan results.                                                                   |
| Quality Profiles                    | [Custom SAST and IaC rules](/code-scanning/scanning-practices/add-custom-sast-iac-rules.md) and [Custom Code Rules](/code-quality/add-custom-code-rules.md)                                                                                                           | Move only rules that are still useful. Aikido already filters out many low-signal checks.                                                            |
| Code smells and maintainability     | [Code Quality](/code-quality/code-quality-overview.md)                                                                                                                                                                                                                | Aikido Code Quality covers PR checks, repository scans, and custom team standards.                                                                   |
| False Positive, Won't Fix, Accepted | [Ignore issues](/getting-started/core-functionalities/ignore-issues-to-remove-issues-from-main-feed.md) and [Snooze issues](/getting-started/core-functionalities/snooze-issues-for-later.md)                                                                         | Use ignore for accepted risk or false positives. Use snooze when you want to revisit an issue later.                                                 |
| `//NOSONAR` suppressions            | [`NOAIKIDO` comments](/code-scanning/scanning-practices/excluding-sast-findings-using-comments.md) and [`.aikido` files](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md)                                                                      | Prefer UI ignores for reviewable decisions. Use `NOAIKIDO` when the exception belongs in the repository.                                             |
| Security Hotspots                   | [SAST](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md), [SAST AutoTriage](/aikido-agent/sast-autotriage.md), and [Reachability Analysis](/getting-started/reachability-analysis/introduction-to-reachability-analysis.md) | This isn't a one-to-one mapping. Aikido focuses on actionable security findings and uses triage signals to reduce manual review of low-risk results. |
| Reports and dashboards              | [Security Audit Report](/compliance-and-reporting/reports/security-audit-report.md) and [Trends Over Time Report](/compliance-and-reporting/reports/trends-over-time-report.md)                                                                                       | Use Aikido reports for security reviews, compliance evidence, and trend tracking.                                                                    |
| Remediation guidance                | [AutoFix](/autofix-and-remediation/overview-aikido-autofix.md)                                                                                                                                                                                                        | Aikido can create fixes for open-source dependencies, SAST, IaC, and containers.                                                                     |

## Migration checklist

{% stepper %}
{% step %}

### Connect your repositories

Start by connecting the repositories you want Aikido to scan through your source code management (SCM) integration. You usually don't need to add a scanner config file to each repository first. Once Aikido is connected to your SCM, Aikido can scan on push, daily rescans, and pull request activity.

Use the setup guide for your SCM:

* [GitHub](/code-scanning/connect-your-source-code/connect-github-account-to-aikido.md)
* [GitLab](/code-scanning/connect-your-source-code/connect-gitlab-account-to-aikido.md)
* [Bitbucket](/code-scanning/connect-your-source-code/connect-bitbucket-account-to-aikido.md)
* [Azure DevOps](/code-scanning/connect-your-source-code/connect-azure-devops-projects-to-aikido.md)

See [Connect Your Source Code](/code-scanning/connect-your-source-code.md) for all setup options.
{% endstep %}

{% step %}

### Match your Sonar workflows

Enable the Aikido workflows that match what Sonar does for the repository today:

* SAST and IaC scanning for code and infrastructure security
* Open-source dependency scanning for CVEs and license risk
* Secrets scanning with liveness checks
* Code Quality for maintainability and team-specific standards
* Code Coverage for test coverage reports from CI, including a check that fails below your coverage threshold
* PR gating for pull request blocking

If you use Sonar in the IDE, review [Aikido IDE Plugins](/ai-and-dev-tools/ide-plugins-overview.md) so developers can catch issues before they commit.
{% endstep %}

{% step %}

### Move useful custom rules

Don't copy every Sonar Quality Profile rule by default. First, check whether Aikido already covers the risk or pattern. Aikido often has broader coverage out of the box, so many Sonar rules won't need a custom replacement.

Move only the rules that still add value, such as rules that represent environment-specific security, architecture, or team-standard decisions.

Use [Custom SAST and IaC Rules](/code-scanning/scanning-practices/add-custom-sast-iac-rules.md) for security patterns and infrastructure checks. Use [Custom Code Rules](/code-quality/add-custom-code-rules.md) for team-specific code quality expectations.
{% endstep %}

{% step %}

### Move accepted-risk decisions

Review Sonar issues marked False Positive, Won't Fix, or Accepted. Recreate the decisions you still agree with in Aikido.

Use [Ignore issues](/getting-started/core-functionalities/ignore-issues-to-remove-issues-from-main-feed.md) for UI-based decisions. Use [`.aikido` files](/code-scanning/scanning-practices/ignore-via-code-with-aikido-files.md) or [`NOAIKIDO` comments](/code-scanning/scanning-practices/excluding-sast-findings-using-comments.md) when the exception should live in the repository.
{% endstep %}

{% step %}

### Configure PR and release gating

Replace Sonar Quality Gate checks with Aikido PR gating or release gating for new security and code-quality issues.

In Aikido, you can set severity thresholds, choose which scan types can fail a check, and start with an always-green setup if you want visibility before blocking merges. See [PR Gating Overview](/pr-and-release-gating/aikido-ci-gating-functionality.md) and [Aikido CLI: Release Gating](/pr-and-release-gating/cli-for-pr-and-release-gating/aikido-cli-release-and-pr-gating.md).
{% endstep %}

{% step %}

### Update CI and branch protection

After Aikido is configured for the repository, remove the matching Sonar scanner job from CI and update branch protection rules to require the Aikido check instead.

If a repository depends on a specialized check outside your Aikido setup, share it with us so we can prioritize a replacement or suggest an alternative workflow.
{% endstep %}

{% step %}

### Review reporting and ownership

Make sure security, engineering, and compliance teams know where to find their new views in Aikido.

Use the main feed for day-to-day triage, PR scan results for developer workflows, and reports for audit or trend reviews. If you route findings into Jira, Linear, Slack, or Microsoft Teams, connect those workflows before removing old Sonar notifications.
{% endstep %}
{% endstepper %}

## Cutover checklist

Before you remove Sonar from a repository, confirm that:

* Aikido scans the repository and the right branch.
* Coverage reports upload to Aikido if you used Sonar coverage, and the coverage check uses the threshold you want.
* PR gating is configured with the severity threshold your team wants.
* Required branch protection checks point to Aikido.
* Important custom rules have been recreated or intentionally dropped.
* Accepted-risk decisions have been reviewed in Aikido.
* Dedicated CI checks remain in place if your team still needs them.
* Developers know where to view findings, request ignores, and use AutoFix.

## Need help?

Migrating from Sonar can touch security, engineering, and compliance workflows. If you're not sure how to map a Sonar rule, gate, or report to Aikido, open the **Intercom chat** in the bottom right corner. Our team is here to help.
