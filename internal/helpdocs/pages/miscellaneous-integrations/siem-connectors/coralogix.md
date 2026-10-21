# Coralogix

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Coralogix for centralized monitoring, so your security findings are searchable alongside the rest of your observability data.

This is useful if your team already uses Coralogix and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Coralogix:

* **Region** - the Coralogix region your account is provisioned in (`us1`, `us2`, `us3`, `eu1`, `eu2`, `ap1`, `ap2`, or `ap3`). This determines which regional endpoint your data is sent to, so it must match your account.
* **Send-Your-Data API Key** - from your Coralogix account.
* **Application Name and Subsystem Name** - labels Coralogix uses to organize incoming data.

### Connecting Coralogix

{% stepper %}
{% step %}

#### Get your Coralogix credentials

1. Sign in to Coralogix and check your account URL against [Coralogix's endpoints reference](https://coralogix.com/docs/integrations/coralogix-endpoints/) to confirm your **Region**.
2. In Coralogix, go to **Settings > API Keys** and click **Add Key**.
3. Select **Send Your Data Key**, give it a name (for example, "Aikido"), and create it.
4. Copy the key immediately, Coralogix only shows it once.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Coralogix** and fill in:

**Settings**

* **Region** - the Coralogix region your account is provisioned in. Defaults to `us1`.
* **Application Name** - stamped on every record. Coralogix uses this to separate environments or tenants, for example `aikido`.
* **Subsystem Name** - stamped on every record. Coralogix uses this to separate components within an application, for example `prod`.
* **Communication Type** - the transport used to reach Coralogix. Currently only **REST (Send-Your-Data)** is supported.
* **Send-Your-Data API Key** - the key you generated in Coralogix.

**Options**

* **Timestamp Record Location** *(optional)* - a JSONPath or gjson path pointing at a timestamp field in the incoming data (for example `$.timestamp`). Leave empty to stamp every record with the current time when it's sent.
* **Batch Config** - controls how data is batched before sending:
  * **Batch Data Size** - maximum total size in bytes of a batch before it's sent, defaults to under Coralogix's 2 MB per-request limit.
  * **Batch Record Count** - maximum number of records to accumulate before a batch is sent automatically.
  * **Maximum Flush Interval** - maximum time in seconds to hold records before sending, regardless of size or count.

Click **Test connection** to confirm Aikido can reach Coralogix, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Coralogix integration page, under **Issue input settings**, configure which Aikido issues are synced:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Coralogix (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Coralogix, search for records under the **Application Name** and **Subsystem Name** you configured. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up.

### Limitations

* Requests over 2 MB are rejected by Coralogix. The default batch size stays under this, but if you increase it, you may hit `413` errors.
* **Application Name** and **Subsystem Name** are fixed for the whole connection, they can't vary per issue.
* Only the REST (Send-Your-Data) transport is supported.
* Requests that don't complete within 30 seconds are treated as failures and retried.

### Troubleshooting

* **Connection refused or timeouts** - confirm **Region** matches the region your Coralogix account is actually provisioned in. The wrong region can produce timeouts or `401`/`403` errors.
* **`401` / `403` errors** - regenerate the Send-Your-Data API key and confirm it has the `SendData` permission. Keys are region-scoped, a key from a different Coralogix account or region won't work.
* **`413 Payload Too Large`** - reduce **Batch Data Size** so batches stay under Coralogix's 2 MB request cap.
* **`429 Too Many Requests`** - Coralogix is rate-limiting your account. Lower the sync frequency, or contact Coralogix to raise your ingest quota.
* **Records show the wrong time** - set **Timestamp Record Location** to the field that carries the event time.
* **Records show up under the wrong app or subsystem** - double check **Application Name** and **Subsystem Name**, both are stamped verbatim onto every record.

### Managing the integration

From **Settings > Integrations > Coralogix** in Aikido, you can:

* **Manage Integration** - update the region, application/subsystem names, batching, or API key.
* **Remove Integration** - disconnect Coralogix. This stops all future syncs; it doesn't delete data already sent to Coralogix.
