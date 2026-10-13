# RunReveal

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to RunReveal for centralized monitoring, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses RunReveal as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from RunReveal:

* **Webhook ID** - from a Generic Webhook or Structured Webhook source in RunReveal.
* **Bearer Token** - required if you're using a Structured Webhook source, or if auth is enabled on a Generic Webhook source.

### Connecting RunReveal

{% stepper %}
{% step %}

#### Get your RunReveal webhook details

1. In the RunReveal portal, go to **Sources > New Source** and choose **Generic Webhook** or **Structured Webhook**.

{% hint style="warning" %}
If you choose **Structured Webhook**, every record Aikido sends must match RunReveal's [Structured Webhooks API schema](https://docs.runreveal.com/sources/source-types/structured-webhooks#structured-webhooks-api). Records that don't match are rejected on RunReveal's side. If you're not sure, use **Generic Webhook**, it accepts any JSON payload.
{% endhint %}

2. Give the source a name (for example, `Aikido`) and save it.
3. Copy the **Webhook ID**, only the ID (the last part of the URL RunReveal shows), not the full URL.
4. Copy the **Bearer Token**. This is required for Structured Webhooks, and optional for Generic Webhooks unless you've enabled auth on the source.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **RunReveal** and fill in:

**Settings**

* **Webhook ID** - the RunReveal webhook ID. Paste only the ID, not the full URL shown in the RunReveal UI.

**Options** *(optional, sensible defaults are pre-filled)*

* **Batch Data Size** - maximum total size in bytes of records in a batch before sending, default `1048576` (1 MiB).
* **Batch Record Count** - maximum number of records to accumulate before sending, default `100`.
* **Maximum Flush Interval** - maximum time in seconds to hold records before sending, default `5`.

Whichever of these three limits is hit first triggers the send.

**Authentication**

* **Bearer Token** - the bearer token issued by RunReveal for this webhook source. Leave empty if the source has auth disabled.

Click **Test connection**, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the RunReveal integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to RunReveal (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

{% hint style="warning" %}
RunReveal's webhook endpoint returns a success response for almost any request, including a wrong or nonexistent Webhook ID, an invalid token, or an empty payload. **Test connection passing, or issues appearing to sync, does not confirm data actually landed in RunReveal.** Always verify inside the RunReveal portal itself.
{% endhint %}

In the RunReveal portal, query the `logs` table filtering on `sourceType = 'webhook'` and your source name to confirm records are landing. To check for ingestion errors (including Structured Webhook schema mismatches), go to RunReveal's [source errors page](https://app.runreveal.com/dash/sources/errors) and select your source type and name.

### Managing the integration

From **Settings > Integrations > RunReveal** in Aikido, you can:

* **Manage Integration** - update the Webhook ID, batch settings, or bearer token.
* **Remove Integration** - disconnect RunReveal. This stops all future syncs; it doesn't delete data already sent to RunReveal.
