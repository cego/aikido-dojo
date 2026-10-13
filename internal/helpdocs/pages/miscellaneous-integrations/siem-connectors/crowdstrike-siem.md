# CrowdStrike Falcon Next-Gen SIEM

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to CrowdStrike Falcon Next-Gen SIEM, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses CrowdStrike as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, you need:

* An active CrowdStrike Falcon subscription with Next-Gen SIEM access.
* A **Next-Gen SIEM endpoint** URL and **auth token**, generated from a Monad connector set up inside CrowdStrike (see below).

### Connecting CrowdStrike

{% stepper %}
{% step %}

#### Set up the Monad connector in CrowdStrike

1. Log in to your CrowdStrike console, open the main menu, and go to **Next-Gen SIEM > Data onboarding**.
2. Find the **Monad** connector on the Data Sources page and click **Configure**.
3. Fill in the connector details:
   * **Data Source** - select your specific data source type.
   * **Data Type** - select **JSON** (recommended for best compatibility).
   * **Connector name** - a descriptive name (for example, "Aikido").
   * **Description** - optional.
4. Under **Parser details**, select an existing parser or create a new one that matches the data format Aikido sends.
5. Accept the terms and click **Save**.
6. Once the connector finishes initializing, look for the **Generate API key** option in the banner and click it.
7. Copy and securely store both the generated **API key** and the **API URL**, you'll paste these into Aikido next.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **CrowdStrike Falcon Next-Gen SIEM** and fill in:

**Settings**

* **Next-Gen SIEM endpoint** - the CrowdStrike Monad connector ingestion URL, in the format `https://<connection-id>.ingest.<your-region>.crowdstrike.com/services/collector`.

**Authentication**

* **Next-Gen SIEM auth token** - the CrowdStrike Monad connector API token you generated.

Click **Test connection** to confirm Aikido can reach CrowdStrike, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the CrowdStrike integration page, under **Issue input settings**, configure which Aikido issues are synced:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to CrowdStrike (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In CrowdStrike's Next-Gen SIEM, search for events from your configured data source. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up.

### Limitations

* Maximum 100 events, or 16 MB, per batch.
* Events older than 7 days may be rejected by CrowdStrike.
* API rate limits apply based on your CrowdStrike subscription.

### Troubleshooting

* **Authentication failures** - confirm the auth token hasn't expired and has the right permissions, and double check the endpoint URL.
* **Connection timeouts** - check network connectivity and any firewall rules between Aikido and CrowdStrike, and validate proxy settings if you use one.
* **Data rejected** - confirm your CrowdStrike parser matches the data format Aikido sends, and that events aren't exceeding the batch or size limits above.

### Managing the integration

From **Settings > Integrations > CrowdStrike Falcon Next-Gen SIEM** in Aikido, you can:

* **Manage Integration** - update the endpoint or auth token.
* **Remove Integration** - disconnect CrowdStrike. This stops all future syncs; it doesn't delete data already sent to CrowdStrike.
