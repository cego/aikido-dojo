# Linking Tasks and PRs

## Use cases

* Improve visibility of security fixes in your task tracker by auto-linking PRs.
* Ensure each autofix PR is associated with the correct task.
* Enable clickable references between PRs and tasks

## What is it?

In order to keep your task tracker in sync with the open Pull Requests, AutoFix can link the Pull Requests it creates with tasks in your task tracker. AutoFix has 2 mechanisms to achieve this:

1. Linking tasks created for issues with Pull Requests
2. Creating a new task for each opened Pull Request

For 1, you need to setup [automatic task creation](https://help.aikido.dev/getting-started/task-management-systems) for issues. For 2, you need to enable `Create task for PR` in the [AutoFix PR configuration](/autofix-and-remediation/configure/autofix-pr-configuration.md). Depending on your workflow you can enable 1 and/or 2.

The task reference placeholder (`$TASK_REF` ) can be used in the following [Pull Request Configuration](/autofix-and-remediation/configure/autofix-pr-configuration.md#pull-request-configuration) fields: branch prefix, commit message and PR title.

The `$TASK_REF` will be populated with task linked to the issues solved in the PR and the new task that was created for the PR if the `Create task for PR` option is enabled.

## Creating a new task for each opened Pull Request

Creating a task for each AutoFix PR can be used to assign people to review and merge the AutoFix PRs. Enable "*Create task for PR*" in your [AutoFix configuration](/autofix-and-remediation/configure/autofix-pr-configuration.md). A task will be created for each AutoFix PR whether it is linked with existing tasks or not.

Supported task trackers:

* Jira
* Jira Datacenter
* Linear
* Asana
* Azure Boards

## Linking Pull Requests with existing Issue Tasks

Task identifiers for issues fixed in the PR are included in the branch name by default. If your workspace is integrated with a task tracker, the PR will automatically be linked to the corresponding task by the task tracker.

## How to integrate Jira and your SCM

As described above, AutoFix can add `$TASK_REF` to your PRs. Configure an integration between your source control management (SCM) platform and task tracker. The tracker then detects the task reference and links the PR to its task.

#### **GitHub → Jira**

**Step 1:** Install [**GitHub for Jira**](https://github.com/marketplace/jira-software-github) from the GitHub Marketplace.

**Step 2:** Grant access to the repositories you want linked.

**Step 3:** In Jira Cloud:

* Go to **Apps → Manage apps → GitHub for Jira → Configure**.
* Connect to your GitHub organization.

**Step 4:** In Aikido, use `$TASK_REF` in your PR title config.

**Step 5:** Jira ticket keys:

* Show as clickable links in Jira.
* Appear in GitHub’s UI (clickable when viewed in Jira).

***

#### **GitLab → Jira**

**Step 1:** In GitLab:

* Go to your project → **Settings → Integrations → Jira**.
* Fill in:
  * **Jira URL**: `https://yourcompany.atlassian.net`
  * **Username**: Your Atlassian email
  * **Password/API Token**: Generate from your Atlassian account

**Step 2:** Save the integration.

**Step 3:** In Aikido, use `$TASK_REF` in your PR title config.

**Step 4:** Jira ticket keys:

* Auto-link in GitLab UI.
* Update Jira’s development panel.

***

#### **Bitbucket → Jira (Cloud)**

**Step 1:** In Bitbucket Cloud:

* Go to **Settings → Linked repositories → Link Jira project**.
* Sign in to your Atlassian account and select the Jira project.

**Step 2:** In Aikido, use `$TASK_REF` in your PR title config.

**Step 3:** Jira ticket keys:

* Auto-link in Bitbucket UI.
* Appear in Jira’s development info panel.

***

#### **Azure DevOps → Jira**

**Step 1:** Install a connector between Azure DevOps and Jira, such as:

* TFS4JIRA
* [Git Integration for Jira](https://marketplace.atlassian.com/apps/4984/git-integration-for-jira-azure-devops-github-gitlab)
* Exalate

**Step 2:** Authorize both Azure DevOps and Jira accounts in the chosen connector.

**Step 3:** Configure the connector to recognize ticket patterns (e.g. `ABC-\d+`).

**Step 4:** In Aikido, use `$TASK_REF` in your PR title config.

**Step 5:** Jira ticket keys:

* Link PRs to Jira and appear in the Jira development panel.

## Closing tasks

Aikido can auto-close tasks when issues are resolved or ignored for [Jira](https://help.aikido.dev/getting-started/task-management-systems/advanced-functionalities/auto-close-jira-tasks-when-aikido-issues-are-resolved) and [Linear](https://help.aikido.dev/getting-started/task-management-systems/advanced-functionalities/auto-close-linear-tasks-when-aikido-issues-are-resolved). If auto-close tasks is enabled, AutoFix will automatically close the tasks it created when the PRs are closed.
