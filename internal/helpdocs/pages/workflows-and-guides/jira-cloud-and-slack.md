# Jira Cloud and Slack

Configure Aikido to send Slack messages and create Jira tickets when new vulnerabilities are found.

Linking your ticket management system and chat solution to Aikido will increase the visibility of new vulnerabilities, allow your team to bring vulnerability management into their development cycle, and increase the transparency for the entire organisation.

## Workflow

**Trigger**: A new vulnerability is detected by Aikido

**Actions**:

* New ticket created in Jira
* Message on Slack

## How to

* [ ] [Set-up Slack integration](/getting-started/chat-and-alerts/slack-notifications.md)
  * [ ] **Optionally**: Send notifications to different channels based on a number of parameters (severity, team): [Send alerts to multiple Slack channels](/getting-started/chat-and-alerts/send-alerts-to-multiple-slack-channels.md)
* [ ] [Set-up Jira integration](/getting-started/task-management-systems/all-supported-task-trackers/jira-cloud.md)
  * [ ] [Enable "Automated Task Creation"](/getting-started/task-management-systems/all-supported-task-trackers/jira-cloud.md#id-2-automated-task-creation)
