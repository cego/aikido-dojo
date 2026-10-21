# Asana

This one-time setup *per workspace* allows everyone in your Aikido organization to create issues directly in Asana.

Following use cases are supported :

* **Automated Ticket Creation**: Automatically create and push tickets to specified Asana projects for seamless tracking of security issues.
* **Manual Ticket Addition**: Manually add security issue tickets to Asana, ensuring targeted attention for critical vulnerabilities.

## Connecting the Aikido App to Asana <a href="#connecting-the-aikido-app-to-asana" id="connecting-the-aikido-app-to-asana"></a>

1. Navigate to [Integration Settings](https://app.aikido.dev/settings/integrations) within the Aikido app.
2. In the 'Task Trackers' section, select 'Asana'
3. A prompt will request authorization for Asana.

   [image: Asana integration unauthorized; prompt to authorize access for issue-tracking.]

   ​
4. Login into your Asana account
5. Grant Aikido permission to access your Asana workspace

   [image: Asana permission request for Aikido Security to access and modify your account data.]
6. Once authorized, Aikido is successfully connected to Asana, enhancing your task management capabilities 🚀

[image: Configure Asana integration to map repositories and enable automatic task creation.]

## ​Options for Task Creation in Asana via Aikido <a href="#options-for-task-creation-in-asana-via-aikido" id="options-for-task-creation-in-asana-via-aikido"></a>

There are two different options to create new tasks from Aikido into Asana.

1. Manually create tasks from the Aikido interface
2. Automatically create new tasks via Aikido's auto creation functionality.

## 1. Manual Task Creation <a href="#id-1-manual-task-creation" id="id-1-manual-task-creation"></a>

1. Hover over any issue in your feed and click the ***+*** in the **assignee column.**\
   ​

   [image: Task assignee section with add and options buttons visible.]

   Alternatively, you can click the triple dots in the last columns to open up the action menu. If you have grouped issues, the triple dot action menu is available on every subissue.

   [image: Action menu with options to manage, snooze, ignore, share, or adjust task severity.]

   ​
2. Fill in the required details in the popup modal.\
   ​

   [image: Create a ClickUp task reporting a high-severity XSS vulnerability with details and link.]
3. The newly created task in Asana will be linked in the Aikido Issue Detail under the 'Tasks' tab (sidepanel).

[image: High-priority task warning: Unescaped input may cause XSS attacks.]

## 2. Automated Task Creation <a href="#id-2-automated-task-creation" id="id-2-automated-task-creation"></a>

{% hint style="info" %}
Aikido will automatically create tasks **every hour in bulk.** There is at the moment no option to trigger this manually.
{% endhint %}

1. Go to the [Asana Integration](https://app.aikido.dev/settings/integrations/tasktracker) settings
2. Make sure to enable '**Autocreation**' by clicking the toggle to **On.**
3. Define the criteria for automatic task creation.

You can configure all of these settings in the autocreation modal:

* **On/Off:** Enable or disable autocreation. When off, the threshold is effectively set to `None`.
* **Severity threshold:** Create tickets for issues at `Critical`, `High`, `Medium`, or `Low` and above. For example, `High` creates tickets for `High` and `Critical`.
* **Issue types:** Create tickets for all issue types, or limit them to specific types like `SCA`, `SAST`, `Secrets`, `Cloud`, `IaC`, `Malware`, or `License`.
* **Daily limit:** Set the maximum number of tickets created per day. The default is `25`. Some plans support higher limits.
* **Ticket creation mode:** Create one ticket per issue group, or create separate tickets per location or scope. Per-location mode can create more tickets.
* **Scope:** Available only in per-location mode. Target all repos or clouds, only mapped repos or teams, or specific repos or clouds with include and exclude lists.

4. Aikido will then autonomously generate Asana tickets based on these settings 🚀
