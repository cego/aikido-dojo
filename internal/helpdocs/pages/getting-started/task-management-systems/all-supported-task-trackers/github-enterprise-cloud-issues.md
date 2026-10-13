# GitHub Enterprise Cloud Issues

This one-time setup *per workspace* allows everyone in your Aikido organization to create issues directly in a GitHub Enterprise Cloud **data residency** tenant (`SUBDOMAIN.ghe.com`).

Following use cases are supported:

* **Manual Ticket Addition**: Manually sync security issue tickets to GitHub Issues.
* **Automated Ticket Creation**: Automatically create and push tickets to specified GitHub Issue repos.

{% hint style="info" %}
This integration is for **GitHub Enterprise Cloud with data residency** tenants (`SUBDOMAIN.ghe.com`). If you are on GitHub.com, use the [GitHub Issues](/getting-started/task-management-systems/all-supported-task-trackers/github-issues.md) integration instead.
{% endhint %}

Unlike GitHub.com, GitHub Enterprise Cloud data-residency tenants require you to register your own OAuth app directly on your tenant. Setup therefore has two parts: first register the OAuth app on your tenant, then connect it in Aikido.

{% embed url="<https://www.youtube.com/watch?v=6lvyiBRTctU>" %}

### Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* A GitHub Enterprise Cloud data-residency tenant (`SUBDOMAIN.ghe.com`).
* You must be an **owner** of the GitHub organization you want to connect. Only one organization is supported per connection.
* Permission to register an OAuth app on your tenant.
* Issues need to be enabled inside GitHub.

### Step 1: Register an OAuth App on Your Tenant

1. Log in to your GitHub Enterprise Cloud tenant at `https://SUBDOMAIN.ghe.com` (e.g. `https://octocorp.ghe.com`).
2. Click your **profile picture** in the top-right corner and select **Settings**.
3. In the left sidebar, scroll to **Developer settings**.
4. Click **OAuth Apps**, then click **New OAuth App**.
5. Fill in the application details:
   * **Application name**: a meaningful name (e.g. `Aikido integration`).
   * **Homepage URL**: your company's website or product URL (e.g. `https://aikido.dev`).
   * **Authorization callback URL**: set this to `https://unify.apideck.com/vault/callback`.
6. Click **Register application**.

### Step 2: Copy the Client ID and Client Secret

1. On the app's settings page, copy the **Client ID**.
2. Click **Generate a new client secret** and copy the value immediately.

{% hint style="warning" %}
Copy the client secret right away. GitHub only displays it once and you won't be able to see it again.
{% endhint %}

### Step 3: Connect the Aikido App to GitHub Enterprise Cloud

1. Navigate to [Integration Settings](https://app.aikido.dev/settings/integrations) within the Aikido app.
2. In the 'Task Trackers' section, select **GitHub Enterprise Cloud issues**.
3. In the modal, fill in:
   * **Tenant Subdomain**: the `SUBDOMAIN` part of your `SUBDOMAIN.ghe.com` tenant.
   * **Client ID**: the Client ID from Step 2.
   * **Client Secret**: the Client Secret from Step 2.
4. Click **Save**. The status will show as **Unauthorized**.
5. Click **Authorize**. You will be redirected to a GitHub authorization screen. Review the access and click **Authorize**.
6. Back in the modal, open the **Organization** dropdown and select your organization. The dropdown is populated after authorizing, and you must be an **owner** of the organization.
7. Click **Save**. The status will now change to **Connected**.
8. Close the modal & open the GitHub Enterprise Cloud Issues Integration page. By default all issues will be synced to the 'Default Repo'. Use **Map To GitHub Repos** to map individual repositories.

You can set a Default Label that will be sent for all Aikido Issues. These will be synced and appear in GitHub.

#### Options for Task Creation in GitHub Enterprise Cloud via Aikido

There are two different options to create new tasks from Aikido into GitHub Issues

1. Manually create tasks from the Aikido interface
2. Automatically create new tasks via Aikido's auto creation functionality.

### Manual Task Creation

1. Hover over any issue in your feed and click the ***+*** in the **assignee column.** Alternatively, click the triple dots in the last column to open the action menu. If you have grouped issues, the triple dot action menu is available on every subissue.
2. Fill in the required details in the popup modal.
3. The newly created task in GitHub Issues will be linked in the Aikido Issue Detail under the 'Tasks' tab (sidepanel).

### Automated Task Creation

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
