# Google SecOps

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Google SecOps (Chronicle), so your security findings are searchable alongside the rest of your security telemetry.

This is useful if your team already uses Google SecOps as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Google Cloud and Google SecOps:

* **GCP Project ID** - the Google Cloud project that hosts your SecOps instance.
* **Region** - the regional endpoint your instance is provisioned in (for example `us`, `europe`, or `asia-southeast1`). This determines which endpoint your data is sent to, so it must match your instance.
* **Instance ID** - the Customer ID (GUID) of your SecOps instance. Find it in SecOps under **Settings > SIEM Settings > Instance ID**.
* **Log Type** - the Chronicle log type the data is ingested as. It must be a type Chronicle supports, and should match the parser you want to apply.
* **Service account JSON key** - for a service account with permission to import logs.

### Connecting Google SecOps

{% stepper %}
{% step %}

#### Create a service account

1. In the Google Cloud Console, go to **IAM & Admin > Service Accounts** and create a service account (for example, "aikido").
2. Grant it a custom role that includes the `chronicle.logs.import` permission on the project that hosts your SecOps instance. Google recommends this least-privilege approach. The predefined **Chronicle API Editor** role also works, but it grants broader access than log import needs.
3. Open the service account's **Keys** tab, click **Add key > Create new key**, and select **JSON**.
4. Save the downloaded key file securely. You'll paste its contents into Aikido next. Aikido accepts the key as raw JSON or base64-encoded JSON.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to [**Settings > Integrations**](https://app.aikido.dev/settings/integrations/) in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Google SecOps** and fill in:

**Settings**

* **GCP Project ID** - the project that hosts your SecOps instance.
* **Region** - your SecOps regional endpoint. Defaults to `us`.
* **Instance ID** - your SecOps instance GUID.
* **Log Type** - the Chronicle log type sent on every batch, for example `WINDOWS_DHCP`, `OKTA`, or `CORELIGHT`.
* **Log Entry Time Field** - a JSON path pointing at the event timestamp (Chronicle's `logEntryTime`), for example `first_detected_at` for when the issue was first created. See the [SIEM Field Reference](/miscellaneous-integrations/siem-connectors/siem-field-reference.md) for the fields in each issue. Falls back to the current time if the field is missing or can't be parsed.
* **Collection Time Field** - a JSON path pointing at the collection timestamp (Chronicle's `collectionTime`). The Aikido issue payload has no collection timestamp field, so any path you enter falls back to the current time. The field is required, so you can enter the placeholder path `metadata.ingestionTime`.

**Authentication**

* **Service Account Credentials JSON** - the JSON key you created in Google Cloud. Aikido stores it encrypted and never returns it once saved.

**Options** *(optional)*

* **Environment Namespace** - a static tag used to partition your data in SecOps, for example `aikido`.
* **Forwarder Resource Name** - the resource name of an existing forwarder, in the format `projects/<p>/locations/<r>/instances/<i>/forwarders/<f>`.
* **Endpoint Override** - the full URL to send to, bypassing the region-based default. Use it for Private Service Connect or test environments.
* **Compress Request** - gzip-compress the request body. Reduces egress at the cost of CPU. Off by default.
* **Batch Config** - controls how data is batched before sending:
  * **Batch Record Count** - maximum number of records to accumulate before a batch is sent. Defaults to 500.
  * **Batch Data Size** - maximum total size in bytes of a batch before it's sent. Defaults to 2097152 (2 MiB).
  * **Maximum Flush Interval** - maximum time in seconds to hold records before sending, regardless of size or count. Defaults to 30.

Click **Test connection** to confirm Aikido can reach Google SecOps, then **Connect**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Google SecOps integration page, under **Issue input settings**, configure which Aikido issues are synced:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Google SecOps (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Google SecOps, check for the ingested logs using **Search**, a **Raw Log Scan**, or the **Data Feeds** page under **Settings > SIEM Settings**. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up.

### Limitations

* Batches over 4 MB are rejected by Google SecOps. If you increase **Batch Data Size**, keep it at or below this limit.
* **Log Type** is fixed for the whole connection, it can't vary per issue.
* Your Google SecOps network must allow access from Aikido. See [SIEM Egress IP Addresses](/miscellaneous-integrations/siem-connectors/siem-egress-ip-addresses.md) if you restrict inbound traffic.

### Troubleshooting

* **`401` errors** - the service account key is invalid or expired. Generate a new JSON key and update it in Aikido.
* **`403` errors** - the service account is missing the `chronicle.logs.import` permission on your instance.
* **`404` errors** - double check **GCP Project ID**, **Region**, **Instance ID**, and **Log Type**. A region that doesn't match where your instance is provisioned is a common cause.
* **`400` errors** - the **Log Type** isn't supported by Chronicle. Use a supported log type.
* **`413 Payload Too Large`** - reduce **Batch Data Size** or **Batch Record Count** so batches stay under 4 MB.
* **Events show the wrong time** - set **Log Entry Time Field** to the field that carries the event time.

### Managing the integration

From **Settings > Integrations > Google SecOps** in Aikido, you can:

* **Manage Integration** - update the project, region, instance, log type, batching, or credentials.
* **Remove Integration** - disconnect Google SecOps. This stops all future syncs; it doesn't delete data already sent to Google SecOps.
