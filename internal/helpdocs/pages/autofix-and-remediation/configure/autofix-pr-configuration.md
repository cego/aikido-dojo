# AutoFix PR Configuration

## Pull Request Configuration <a href="#pull-request-configuration" id="pull-request-configuration"></a>

You can configure the AutoFix pull requests to your own needs with the ["PR Configuration" settings](https://app.aikido.dev/issues/fix/settings). It's possible to change the following fields:

* **Title Prefix:** Customize the title to match your internal processes.
* **Labels:** Add specific labels to the PR for easier sorting and tracking, labels are comma-separated.
* **Branch name:** Customize the prefix of the branch Autofix creates.
* **Commit Message:** Customize the prefix of the commit message in the PR.
* **Summary:** Optionally include an additional message for each PR. For example add custom instructions for engineers or additional context.
* **Create Draft PR:** create draft PRs instead of ready for review.
* **Auto-Merge PR:** Aikido relies on built-in Github, Gitlab and Azure DevOps auto-merge features for this, learn more [here](/autofix-and-remediation/automation-and-merging/auto-merge-autofix-pull-requests.md).
* **Create Task for PR:** important for compliance reasons. This will automatically create a new task in your connected task tracker, linking to the created PR.

The following **variables** can be used in PR titles, labels, branch names and commit messages:

* `$TASK_REF` is the reference (eg. JIRA-1234) of the tasks linked to the issues solved in the PR and can be used to [link tasks with the PR](/autofix-and-remediation/configure/linking-tasks-and-prs.md)
* `$SEVERITY` is the highest severity of the issues solved in the PR (eg. `critical`, `high`, `medium`, `low`)
* `$AUTOFIX_TYPE` is replace by `dependency`, `sast`, `iac`, `quality`, `container`, `pentest`, `audit`, `gradle`, `dotnet` or `zen`

{% hint style="info" %}
Aikido doesn't support configuring default reviewers for AutoFix PRs. To make sure the right people review them, use CODEOWNERS and required reviews on the repository. See [Assigning Default Reviewers for AutoFix PRs](/autofix-and-remediation/faq/assigning-default-reviewers-for-autofix-prs.md).
{% endhint %}

## Repository-specific configuration

Global AutoFix settings in the Aikido UI apply to all repositories within the Aikido workspace. To configure settings specific to a repository, use the `.aikido` file located at the root of the repo. The configurations specified in this file will override the global settings.

For example, this configuation:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    branch_prefix: "aikido-"
    title_prefix: "Aikido "
    commit_prefix: "[fix] "
    labels: "security"
    draft_pr: true
    auto_merge: true
    create_task_for_pr: true
    task_title_prefix: "[AIK] "
    disable_codeowners_as_reviewers: true
    self_approve_pr: true
  automated:
    disabled: false
    dependency:
      update_policy: minor_and_patch_versions_only
      single_pr_for_repo: false
      lockfile_groups:
        - group_name: service1
          lockfiles:
            - 'service1/*'
  base_branch: develop
```

{% endcode %}

* under `pr_config`:
  * the branch name starts with `aikido-` [\[+\]](#branch-prefixes)
  * the PR title starts with `Aikido` [\[+\]](#pull-request-title-prefixes)
  * the commit starts with `[fix]` [\[+\]](#commit-prefixes)
  * the PR has a label `security` [\[+\]](#labels)
  * create draft PRs [\[+\]](#draft-pr)
  * the PR will be auto-merged by the SCM (eg Github) when the required approvals are present [\[+\]](#auto-merge)
  * a task (eg. Jira) is created for the PR [\[+\]](#create-task-for-pr)
  * the title of the task (eg. Jira) starts with `[AIK]` [\[+\]](#task-title-prefix)
  * *Gitlab specific:* the codeowners will not be added as reviewers on the PR [\[+\]](#disabling-codeowners-as-reviewers)
  * *Bitbucket specific:* AutoFix will approve it's own PR [\[+\]](#self-approve-bitbucket-prs)
* under `automated`:
  * automated autofixes are enabled for the repo [\[+\]](#disabling-automated-autofixes-at-the-repo-level)
  * under `dependency`:
    * only minor and patch versions upgrades are performed [\[+\]](#setting-the-update-policy)
    * creates multiple PRs for the lockfiles in the repo [\[+\]](#grouping-autofixes-per-lockfile-or-per-repo)
    * under `lockfile_groups`
      * a single PR will be created for all lockfiles in the `service1` directory, separate PRs will be created for the other lockfiles in the repo [\[+\]](#lockfile-groups)
* global:
  * creates the PR starting from the `develop` branch [\[+\]](#set-a-custom-base-branch)

The sections below explain each configuration parameter in detail.

### **Branch prefixes**

Branch prefixes can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    branch_prefix: "prefix-for-branch-"
```

{% endcode %}

AutoFix will use the `branch_prefix` value when creating a branch for that specific repository. For this example AutoFix will create branches like `prefix-for-branch-update-dependency-xyz`.

The branch prefix defined in the `.aikido` file takes priority over the branch prefix configured in the UI.

### Branch template

AutoFix creates a unique branch name by adding a random identifier to the branch name. This random identifier might not match your branch naming rules. In this case a branch template can be used to fully define the branch name format.

In many organizations these branch naming rules are used to enforce linking between tasks and branches. The following configuration makes AutoFix create specific branch names:

{% code title=".aikido" %}

```
autofix:
  pr_config:
    branch_template: fix/$TASK_REF
```

{% endcode %}

The `$TASK_REF` is replaced by the tasks linked with the issue fixed in the Pull Request. For example: `fix/AIK-1234` for a single task or `fix/AIK-1234-AIK-4567` for multiple tasks.

If you have `Create task for Pull Request` enabled, the automatically created task identifier will be using in the branch name.

In case of a strict branch format, you can use the following variables:

* `$TASK_REF_LOWER`: the task identifiers in lower case format
* `$TASK_REF_0` : the first task identifier
* `$TASK_REF_0_LOWER` : the first task identifier in lower case format

### Pull request title prefixes

Pull request title prefixes can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    title_prefix: "Title Prefix"
```

{% endcode %}

AutoFix will start the pull request title with `Title Prefix` . The title prefix defined in the `.aikido` file takes priority over the branch prefix configured in the UI.

### **Pull request title template**

Pull request titles can be customized further by adding a pull request title template in a `.aikido` file at the root of the repo with the following contents:

{% code title=".aikido" %}

```json
autofix:
  pr_config:
    title_template: "Aikido: $MSG - security"
```

{% endcode %}

AutoFix will populate `$MSG` when the pull request is created. This allows you to set a prefix, suffix or static title for Aikido AutoFix PRs. If both the `title_prefix` and `title_template` are provided, the template takes precidence.

### **Commit prefixes**

Commit prefixes can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    commit_prefix: "fix: "
```

{% endcode %}

AutoFix will use the `commit_prefix` value when creating a commit for that specific repository. For this example AutoFix will create commits like `fix: updating package uuid`

The commit prefix defined in the `.aikido` file takes priority over the commit prefix configured in the UI.

### Commit template

Commit messages can be customized further by adding a commit message template in a `.aikido` file at the root of the repo with the following contents:

{% code title=".aikido" %}

```json
autofix:
  pr_config:
    commit_template: "fix: $MSG - security"
```

{% endcode %}

AutoFix will populate `$MSG` when the commit message is created. This allows you to set a prefix, suffix or static commit message. If both the `commit_prefix` and `commit_template` are provided, the template takes precidence.

### Labels

Labels can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    labels: "aikido,security,$AUTOFIX_TYPE,$SEVERITY"
```

{% endcode %}

AutoFix will use the `labels` value when creating a Pull Request for that specific repository. For this example AutoFix add both the `aikido` and `security` label on the pull requests.

**Placeholders**

* `$AUTOFIX_TYPE`: the AutoFix type will be added as a label in when the PR is created. Possible values: `dependency`, `sast`, `iac`, `quality`, `container`, `pentest`, `audit`, `gradle`, `dotnet`, `zen`
* `$SEVERITY`: the maximum severity of the issues fixed in the Pull request will be added as a label when the PR is created. Possible values: `high`, `medium`, `low` .

The labels defined in the `.aikido` file takes priority over the labels configured in the UI. When defined in the `.aikido` file, the labels configured in the UI are ommitted.

### Draft PR

Draft PR can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    draft_pr: true
```

{% endcode %}

The example above creates Pull Requests as Draft for the specific repo.

The draft PR configuration in the `.aikido` file takes priority over the draft PR configuration in the UI.

### Auto-merge

Auto-merging can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    auto_merge: true
```

{% endcode %}

The example above enables auto-merging for all PRs in the specific repo.

The auto-merge configuration in the `.aikido` file takes priority over the configuration set in the UI.

#### Conditional auto-merge (by type and severity)

Auto-merge can be enabled conditionally based on the AutoFix type and the severity of the issue being fixed. The example below enables auto-merge for pull requests that remediate critical SAST issues.

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    auto_merge:
      sast:
        enabled: true
        severity: critical
```

{% endcode %}

**Supported configuration values**

* AutoFix types: `dependency`, `sast`, `iac`, `container`, `code-quality` and `pentest`
* Severity levels: `critical`, `high`, `medium` and `low`

**Severity logic**

If the severity is set to `medium`, all pull requests addressing critical, high and medium issues will be auto-merged.

If the pull request addresses multiple issues of varying severities, it will auto-merge if the highest severity issue meets the configured threshold. Example: if you configure auto-merge to critical severity, a pull request that contains both a critical and a high issue will still be auto-merged.

#### Instruction based auto-merge

For even more fine grained control, you can use natural language instructions on when to enable auto-merge.

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    auto_merge:
      instructions: >
        Auto-merge dependency updates when every dependency is bumped by a patch
        or minor version only. Never auto-merge a major version bump.
        Auto-merge SAST fixes only when the fix confidence is high and the
        change touches a single file.
        Do not auto-merge anything else.
```

{% endcode %}

The context provided for these instructions contains:

* AutoFix type
* Severity levels
* Confidence (for SAST, IaC & Code Quality issues)
* PR title and description
* List of changed files in the PR
* Diffs for the changed files in the PR

### Create task for PR

AutoFix automatically creates a task in your issue tracker when a pull request is opened in\
Aikido. This is useful for maintaining remediation traceability required by compliance\
standards.

Creating a task can be enabled or disabled for a specific repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    create_task_for_pr: true
```

{% endcode %}

The example above enables creating a task for the specific repo. The configuration in the `.aikido` file takes priority over the configuration in the UI.

### Task title prefix

When "create task for PR" is enabled, a prefix for the task title can be configured per repo by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    task_title_prefix: "[AIK] "
```

{% endcode %}

The example makes sure the task created for the PR starts with `[AIK]` in the specific repo.

The configuration in the `.aikido` file takes priority over the configuration in the UI.

### Disabling CODEOWNERS as reviewers

For Gitlab, Aikido AutoFix will add the `CODEOWNERS` as reviewers on Merge Requests (if there are less than 10 CODEOWNERS for the lockfile). This behavior can be disabled by adding a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    disable_codeowners_as_reviewers: true
```

{% endcode %}

### Self-approve BitBucket PRs

For BitBucket, Aikido AutoFix can approve the PRs it created. Enable this setting if your workflow requires the author of the PR to approve the PR.

{% code title=".aikido" %}

```yaml
autofix:
  pr_config:
    self_approve_pr: true
```

{% endcode %}

### Disabling Automated Autofixes at the repo level

Automated AutoFixes must first be enabled in the Aikido UI. To exclude a specific repository that is already enabled in the UI, add the following configuration to the `.aikido` file at the root level of your repository:

{% code title=".aikido" %}

```yaml
autofix:
  automated:
    disabled: true
```

{% endcode %}

**Note:** Setting `disabled: false` will **not** automatically enable Automated AutoFixes for a repository. This setting is strictly used to opt a repository out when the feature is globally turned on in the Aikido UI.

#### Disabling Automated Autofixes per AutoFix type

Automated AutoFixes can be disabled for specific AutoFix types. For example, you can disable automated SAST fixes while keeping automated PRs enabled for the other AutoFix types. The example below show how to disable every supported type (`sast`, `iac`, `dependency`, `container`, `pentest`, and `code_quality`) use the fine-grained configuration:

{% code title=".aikido" %}

```yaml
autofix:
  automated:
    sast:
      disabled: true
    iac:
      disabled: true
    dependency:
      disabled: true
    container:
      disabled: true
    pentest:
      disabled: true
    code_quality:
      disabled: true
```

{% endcode %}

### Setting the update policy

The update policy can set separately for each type of AutoFix. The settings in the `.aikido` file in the root of the repository overwrite the setting provided in the UI.

{% code title=".aikido" %}

```yaml
autofix:
  automated:
    dependency:
      update_policy: upgrade_all_packages | minor_and_patch_versions_only | critical_and_high_only | critical_issues_only
    sast:
      update_policy: all | critical_and_high_only | critical_issues_only
      confidence: high | medium
    iac:
      update_policy: all | critical_and_high_only | critical_issues_only
      confidence: high | medium
    pentest:
      update_policy: all | critical_and_high_only
```

{% endcode %}

For `dependency` AutoFixes the following update policies can be set:

* `upgrade_all_packages` : upgrade all packages regardless of the versions
* `minor_and_patch_versions_only` : only upgrade to newer minor and patch versions, don't perform major version updates
* `critical_and_high_only` : only upgrade versions for citical and high severity issues
* `critical_issues_only` : only upgrade versions for critical severity issues

For `sast` AutoFixes the following update policies can be set:

* `all` : perform SAST fixes regardless of the severity of the issue
* `critical_and_high_only` : perform SAST fixes only for critical and high severity issues
* `critical_issues_only` : perform SAST fixes only for critical severity issues
* Additionally a **minimum confidence level** can be set in the `confidence` field. This is set to `medium` by default. By setting the confidence to `high`, AutoFix will only create automated PRs for high confidence fixes.

For `iac` AutoFixes the following update policies can be set:

* `all` : perform IaC fixes regardless of the severity of the issue
* `critical_and_high_only` : perform IaC fixes only for critical and high severity issues
* `critical_issues_only` : perform IaC fixes only for critical severity issues
* Additionally a **minimum confidence level** can be set in the `confidence` field. This is set to `medium` by default. By setting the confidence to `high`, AutoFix will only create automated PRs for high confidence fixes.

For `pentest` AutoFixes the following update policies can be set:

* `all` : perform pentest fixes regardless of the severity of the issue
* `critical_and_high_only` : perform pentest fixes only for critical and high severity issues

### Grouping AutoFixes per lockfile or per repo

By default, Aikido AutoFix creates a separate Pull Request **per lockfile** in your repository. This keeps changes isolated and easy to review for standard projects.

For repos with many lockfiles owned by a single team, receiving many individual PRs can be noisy. Enabling a single Pull Request per repository reduces the administrative overhead and CI/CD usage.

To change this behavior, add a `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  automated:
    dependency:
      single_pr_for_repo: true
```

{% endcode %}

### Lockfile groups

By default, Aikido AutoFix creates a separate Pull Request per lockfile in your respository. This keeps changes isolated and easy to review for standard projects.

For repos with interdependent lockfiles, you can create a lockfile group to group all dependency AutoFixes in those lockfiles into a single Pull Request. This avoid duplicate pull requests in case of shared version resolution and makes it easier to review the changes.

Lockfile groups are configured by giving them a name and a list of lockfile patterns. For example by adding the following contents in a `.aikido` file at the root level of the repo, AutoFix will create a single automated dependency update Pull Request for all lockfiles in the `service1` folder and it's subfolders.

{% code title=".aikido" %}

```yaml
autofix:
  automated:
    dependency:
      lockfile_groups:
        - group_name: service1
          lockfiles:
            - 'service1/*'
```

{% endcode %}

### Set a custom base branch

By default, Aikido AutoFix creates Pull Requests based on the scanned branch. If you have a workflow with a different branching strategy, you can configure a custom base branch in the `.aikido` file at the root level of the repo with the following contents:

{% code title=".aikido" %}

```yaml
autofix:
  base_branch: develop
```

{% endcode %}

Make sure the scanned branch and base branch configured in the `.aikido` file are in sync. For example: if you perform a dependency AutoFix and your `package.json` has different content on the base branch compared to the scanned branch, the pull request will be created from the base branch and might not contain the results you would expect when starting from the scanned branch.
