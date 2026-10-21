# Custom Reports Creation

Create custom security reports with Aikido Agent using natural language, then save and share them across your workspace.

**Custom Reports** lets you ask [Aikido Agent](/aikido-agent/aikido-agent.md) questions about your security data in plain language. The Agent builds charts and tables from your findings, and you can refine the report in chat before saving it for your workspace.

{% hint style="info" %}
Custom Reports is available on paid plans and during a pro trial. It is not available to team-only users.
{% endhint %}

## How it works

1. You describe the insights you want in natural language.
2. Aikido Agent generates one or more widgets (tables, bar charts, or line charts) from your workspace data.
3. You refine the report in chat (for example, change the breakdown, add a filter, or adjust the date range).
4. When you're happy with the result, you save it to **Workspace Reports** so the rest of your workspace can view it.<br>

Saved reports keep a cached snapshot of their data. You can refresh that data or change the date range later without rebuilding the report from scratch.

{% hint style="warning" %}
Custom Reports only use real data from your workspace. If Aikido Agent cannot access the data needed for your request, it tells you instead of filling gaps with invented numbers or charts.
{% endhint %}

## Create a custom report

{% stepper %}
{% step %}

### Open Custom Reports

Go to [**Reports**](https://app.aikido.dev/reports) and open **Custom Reports**, or go directly to [Custom Reports](https://app.aikido.dev/reports/custom).
{% endstep %}

{% step %}

### Describe what you need

Click **Create Custom Report**. Enter a description of the insights you want, or pick one of the suggested prompts.
{% endstep %}

{% step %}

### Review and refine

The Agent builds widgets based on your request. Use the chat panel to iterate. For example, ask to group by team, filter to a severity, or switch from a table to a chart.

Your previous reports appear in the **Previous Chats** sidebar so you can reopen one that is still private to you.\

{% endstep %}

{% step %}

### Save to your workspace

When the report looks right, click **Add to Workspace Reports**. Saved reports are shared with your workspace. Until you save, the report stays private to you.
{% endstep %}
{% endstepper %}

## Save to Workspace Reports

New custom reports start private to you. They show up in **Previous Chats**, and other people in your workspace can't see them until you save.

When you're ready to share the report:

1. Open the report and confirm the report looks right.
2. Click **Add to Workspace Reports**.
3. Aikido publishes the report to the shared **Custom Reports** list for your workspace.

Admins and default users in your workspace can then open the report from [Custom Reports](https://app.aikido.dev/reports/custom). Only admins can create reports or save them to the workspace.

## Naming Reports

When Aikido Agent builds your report, it also assigns a **report name** based on what you asked for. That name shows up in **Previous Chats** (while the report is private to you) and in the shared Custom Reports list (after you save).

To rename a report saved to your workspace:

1. Go to [Custom Reports](https://app.aikido.dev/reports/custom).
2. Open the **Actions** menu on the report row.
3. Click **Rename report**.
4. Enter the new **Report Name** and click **Save**.
