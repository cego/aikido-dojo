# GitHub Issues

This one-time setup *per workspace* allows everyone in your Aikido organization to create issues directly in GitHub Issues

Following use cases are supported :

* **Manual Ticket Addition**: Manually sync security issue tickets to GitHub Issues.
* **Automated Ticket Creation**: Automatically create and push tickets to specified GitHub Issue projects/repos.

{% hint style="warning" %}
All issues will be automatically created in the respective repo. If you wish to disable this functionality, contact us.
{% endhint %}

{% hint style="info" %}
This integration is for GitHub.com. If you are on a GitHub Enterprise Cloud with data residency tenant (`SUBDOMAIN.ghe.com`), use the [GitHub Enterprise Cloud Issues](/getting-started/task-management-systems/all-supported-task-trackers/github-enterprise-cloud-issues.md) integration instead.
{% endhint %}

### Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* GitHub account needs to be an organisation account
* Issues need to be enabled inside GitHub

### Connecting the Aikido App to GitHub Issues <a href="#connecting-the-aikido-app-to-github-issues" id="connecting-the-aikido-app-to-github-issues"></a>

1. Navigate to [Integration Settings](https://app.aikido.dev/settings/integrations) within the Aikido app.
2. In the 'Task Trackers' section, select 'GitHub Issues'
3. **Install the Issues App inside your organisation** ([Install Link](https://github.com/apps/aikido-issues/installations/new)). This is needed in order to select your organisation.\
   ​

   [image: GitHub issue-tracking integration setup requires authorization and organization configuration.]
4. Select your organisation and repos\
   ​

   [image: Aikido Issues installation prompt with user selection interface.]
5. When installed succesfully, you will get a notification on top of the page in GitHub. **Return to Aikido.**\
   ​

   [image: Notification confirming "Aikido Issues" was updated for the specified user account.]
6. **Select your organisation** in the modal.\
   ​

   [image: Dropdown menu for selecting an organization in settings.]
7. **Click Save.** The status will now change to **Connected.**

   [image: GitHub issue-tracking integration successfully connected and ready for use.]
8. Close the modal & open the GitHub Issues Integration page. By default all issues will be synced to the 'Default Repo'. **Contact us** if you'd like to have autosyncing issues to the respective repo.

You can set a Default Label that will be send for all Aikido Issues. These will be synced and appear in GitHub.​\
​

[image: Open GitHub issue for minor aws/aws-sdk-php upgrade, labeled low priority security fix.]

#### Options for Task Creation in GitHub Issues via Aikido <a href="#options-for-task-creation-in-github-issues-via-aikido" id="options-for-task-creation-in-github-issues-via-aikido"></a>

There are two different options to create new tasks from Aikido into GitHub Issues

1. Manually create tasks from the Aikido interface
2. Automatically create new tasks via Aikido's auto creation functionality.

### Manual Task Creation <a href="#manual-task-creation" id="manual-task-creation"></a>

1. Hover over any issue in your feed and click the ***+*** in the **assignee column.**\
   ​

   [image: Task management interface showing "Assignee" column with add and options buttons.]

   Alternatively, you can click the triple dots in the last columns to open up the action menu. If you have grouped issues, the triple dot action menu is available on every subissue.\
   ​

   [image: Dropdown menu with task management actions: create, snooze, ignore, copy link, adjust severity.]
2. Fill in the required details in the popup modal.\
   ​

   [image: Form for creating a new GitHub issue task with project, assignee, and description fields.]

   ​
3. The newly created task in GitHub Issues will be linked in the Aikido Issue Detail under the 'Tasks' tab (sidepanel).

[image: Task management dashboard showing a low-priority AWS SDK PHP upgrade task.]

### Automated Task Creation <a href="#automated-task-creation" id="automated-task-creation"></a>

{% hint style="info" %}
Aikido will automatically create tasks **every hour in bulk.** There is at the moment no option to trigger this manually.
{% endhint %}

1. Go to the [GitHub Integration](https://app.aikido.dev/settings/integrations/tasktracker) settings page
2. Make sure to enable '**Autocreation**' by clicking the toggle to **On.**
3. Define the criteria for automatic task creation.

You can configure all of these settings in the autocreation modal:

* **On/Off:** Enable or disable autocreation. When off, the threshold is effectively set to `None`.
* **Severity threshold:** Create tickets for issues at `Critical`, `High`, `Medium`, or `Low` and above. For example, `High` creates tickets for `High` and `Critical`.
* **Issue types:** Create tickets for all issue types, or limit them to specific types like `SCA`, `SAST`, `Secrets`, `Cloud`, `IaC`, `Malware`, or `License`.
* **Daily limit:** Set the maximum number of tickets created per day. The default is `25`. Some plans support higher limits.
* **Ticket creation mode:** Create one ticket per issue group, or create separate tickets per location or scope. Per-location mode can create more tickets.
* **Scope:** Available only in per-location mode. Target all repos or clouds, only mapped repos or teams, or specific repos or clouds with include and exclude lists.

4. Aikido will then autonomously generate GitHub Issues based on these settings 🚀

### Automatically Close Issues in GitHub

Aikido can automatically close linked GitHub issues once the related vulnerability is resolved in Aikido, for example after a successful scan no longer reports it, or when you manually mark it solved. The GitHub issue is only closed when there are no remaining open or snoozed Aikido issues for that link (ignored issues do not block closing). You need a GitHub issue linked to the Aikido finding for this to apply.

***
