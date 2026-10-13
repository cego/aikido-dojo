# Datadog SIEM

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Datadog Logs so your security findings live alongside the rest of your observability and monitoring stack.

This is useful if your team already uses Datadog as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Datadog:

* **API key -** from **Organization Settings > API Keys** in Datadog.
* **Domain URL** – the base domain for your Datadog account's region, for example `datadoghq.com` (default), `us3.datadoghq.com`, `us5.datadoghq.com`, `datadoghq.eu`, `ap1.datadoghq.com`, `ap2.datadoghq.com`, or `ddog-gov.com`. You can confirm your region from the URL you use to log in to Datadog.

### Connecting Datadog

{% stepper %}
{% step %}

#### Get your Datadog API key

In Datadog, go to **Organization Settings > API Keys**, click **New Key**, and give it a descriptive name (for example, "Aikido"). Copy the key - you'll paste it into Aikido in the next step.
{% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Datadog Logs** and fill in:

**Settings**

* **Domain URL** – select the base domain that matches your Datadog account's region (defaults to `datadoghq.com`).

**Options** *(all optional)*

* **Source** – the integration name Datadog associates with these logs. Defaults to `aikido`, so Datadog automatically installs the matching parsers and facets.
* **Tags** – custom tags (for example `team:security`, `env:production`) for filtering and grouping in Datadog.
* **Hostname** – the originating host of the log, if relevant to your setup.
* **Service** – lets you correlate Aikido logs with an existing APM service. Defaults to `aikido`.

**Authentication**

* **API Key** – paste the key you generated in Datadog.

Click **Test connection** to confirm Aikido can reach Datadog, then **Save changes**.

{% endstep %}

{% step %}

#### Configure which issues sync

On the Datadog integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** – for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** – the lowest severity level to sync.
* **Issue type** – limit to specific issue types, or sync all of them.
* **Language** – limit to specific languages, or sync all of them.
* **Sync frequency** – how often Aikido pushes matching issues to Datadog (for example, every hour).

Click **Save Settings**.

{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Datadog, go to **Logs > Search** and filter on `source:aikido` (or the custom source you configured). Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up in Datadog.

### Managing the integration

From **Settings > Integrations > Datadog** in Aikido, you can:

* **Manage Integration** – update the domain URL, options, or API key.
* **Remove Integration** – disconnect Datadog. This stops all future syncs; it doesn't delete logs already sent to Datadog.
