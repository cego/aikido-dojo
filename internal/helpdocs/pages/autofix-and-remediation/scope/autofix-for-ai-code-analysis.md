# AutoFix for AI Code Analysis

AutoFix writes code fixes for security issues found by [AI Code Analysis](https://help.aikido.dev/ai-code-analysis/ai-code-analysis), and delivers them as pull requests (PRs). All programming languages are supported.

You can get a fix in 2 ways:

* **From the AutoFix page:** preview a fix, refine it, and open a PR when you're ready.
* **Automatically:** Aikido opens PRs in the daily static application security testing (SAST) AutoFix run, with the same settings and limits as SAST.

## Prerequisites

* [AutoFix](https://help.aikido.dev/autofix-and-remediation/overview-aikido-autofix) is set up with write access to the repository.
* The finding points to a source code file in a connected repository.

## How it works

AI Code Analysis produces a root-cause analysis for each finding. AutoFix uses that analysis to change the code that causes the issue, not only the line where it was reported.

| | AutoFix page | Automated PRs |
| --- | --- | --- |
| **Trigger** | You click **Create Preview** or **Create PR** | The daily SAST AutoFix run |
| **Base branch** | Default branch, or a branch you pick | Default branch |
| **Where you review** | Preview in the Aikido app, then the PR | The PR |

## Fix findings from the AutoFix page

{% stepper %}
{% step %}
### Open the findings

Go to [**AutoFix** → **Code**](https://app.aikido.dev/issues/fix/code-audit?regionRedirect=1) and select **AI Code Analysis**.
{% endstep %}

{% step %}
### Preview the fix

Click **Create Preview** on a finding. The preview shows the code changes, and the result of running your tests against the fix when it's available.
{% endstep %}

{% step %}
### Adjust the fix

If the fix isn't quite right, click **Refine with Aikido AI** and describe what to change. See [Refine AutoFixes with Aikido AI](https://help.aikido.dev/autofix-and-remediation/automation-and-merging/refine-autofixes-with-aikido-ai). Click **Regenerate** to start over, or **Mark as Ignored** if the finding doesn't need a fix.
{% endstep %}

{% step %}
### Open the pull request

Click **Create PR**. The dropdown next to it has more options:

* **Create PR to Specific Branch:** open the PR against a branch other than the default branch, for example a release branch.
* **Create PR and Ship it:** open the PR and add your repository's deploy label to it.
{% endstep %}
{% endstepper %}

## Create fixes automatically

AI Code Analysis findings are fixed in the daily SAST AutoFix run. They follow the same settings and limits as SAST:

* **Settings:** go to [**AutoFix** → **Settings**](https://app.aikido.dev/issues/fix/settings?regionRedirect=1), turn on **Automatically fix SAST/AI Code Analysis & IaC issues**, and select the repositories, severity levels, and **Only private repos** there.
* **Limits:** SAST and AI Code Analysis share the limits of the daily run. See [Automatic Creation of Pull Requests](https://help.aikido.dev/autofix-and-remediation/automation-and-merging/automatic-creation-of-pull-requests).
* **Grouping:** findings from the same scan, with the same issue type and file, are fixed in 1 PR, with up to 5 findings per PR. The PR description explains which change fixes which finding.

Findings without a file location aren't fixed automatically.

{% hint style="info" %}
Automated PRs follow your [PR configuration](https://help.aikido.dev/autofix-and-remediation/configure/autofix-pr-configuration), such as title prefixes, branch names, and labels.
{% endhint %}

### Override the settings per repository

Use the `sast` section under `autofix.automated` in the repository's `.aikido` file, not `iac`. AI Code Analysis fixes count as high confidence, so the `confidence` setting doesn't skip them.

{% code title=".aikido" %}
```yaml
autofix:
  automated:
    sast:
      update_policy: critical_and_high_only # or all, critical_issues_only
```
{% endcode %}

Setting `disabled: true` in the `sast` section stops both SAST and AI Code Analysis PRs for that repository. See [Repository-specific configuration](https://help.aikido.dev/autofix-and-remediation/configure/autofix-pr-configuration#repository-specific-configuration) for all `.aikido` options.

## What happens after a fix is merged

When all issues linked to an AutoFix PR are resolved, for example because the fix was merged through another PR, Aikido closes the PR and deletes its branch. If someone else has pushed commits to the PR, Aikido leaves it open.

## Language support

AutoFix works with all programming languages. For JavaScript codebases, it also writes new unit tests and updates existing ones. It applies linting rules to the generated code for JavaScript, Python, and .NET.

## FAQ

<details>

<summary>Why didn't Aikido open an automated PR?</summary>

Check that the repository is selected in the SAST automated PR settings, the finding's severity matches your policy, and the `autofix.automated.sast` section of the repository's `.aikido` file doesn't disable fixes. If **Only private repos** is on, public repositories are skipped. Aikido also skips a finding that already had a PR, and waits for the next daily run when the repository has reached the PR limits.

</details>

## Related

{% content-ref url="../automation-and-merging/automatic-creation-of-pull-requests.md" %}
[automatic-creation-of-pull-requests.md](https://help.aikido.dev/autofix-and-remediation/automation-and-merging/automatic-creation-of-pull-requests)
{% endcontent-ref %}

{% content-ref url="../automation-and-merging/refine-autofixes-with-aikido-ai.md" %}
[refine-autofixes-with-aikido-ai.md](https://help.aikido.dev/autofix-and-remediation/automation-and-merging/refine-autofixes-with-aikido-ai)
{% endcontent-ref %}

{% content-ref url="../configure/autofix-pr-configuration.md" %}
[autofix-pr-configuration.md](https://help.aikido.dev/autofix-and-remediation/configure/autofix-pr-configuration)
{% endcontent-ref %}

{% content-ref url="../../ai-code-analysis/ai-code-analysis.md" %}
[ai-code-analysis.md](https://help.aikido.dev/ai-code-analysis/ai-code-analysis)
{% endcontent-ref %}
