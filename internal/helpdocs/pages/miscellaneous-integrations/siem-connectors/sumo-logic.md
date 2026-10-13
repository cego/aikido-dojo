# Sumo Logic

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Sumo Logic for centralized monitoring, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses Sumo Logic as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, you need one of the following from Sumo Logic:

* **An HTTP Collector URL**, set up manually in Sumo Logic, or
* **Access ID and Access Key** with permissions to create collectors and sources, so Aikido can set up the collector for you automatically.

### Connecting Sumo Logic

{% stepper %}
{% step %}

#### Get your Sumo Logic credentials

Pick one of the two connection methods:

**Option A: HTTP Collector URL**

1. In Sumo Logic, go to [Collection](https://service.sumologic.com/collection) and click **Add Collector**.
2. Once the collector is created, add an **HTTP Logs & Metrics** source to it.
3. Copy the collector URL shown, this is your **URL**.
4. Edit the HTTP source and, under **Advanced Options for Logs (Optional)**, make sure **One Message Per Request** is unchecked.
5. If your data may contain newline characters within a JSON message, also enable **Multiline Processing** set to **Detect Automatically**, so messages are indexed correctly.

**Option B: Access ID and Access Key**

1. In Sumo Logic, [generate an Access ID and Access Key](https://help.sumologic.com/docs/manage/security/access-keys/) with permissions to create collectors and sources.
2. Aikido will automatically create and manage a hosted collector and HTTP source for you the first time it connects, and reuse them on subsequent connections.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Sumo Logic** and fill in either:

* **URL** - your Sumo Logic HTTP Collector URL, or
* **Access ID** and **Access Key** - so Aikido can create and manage the collector for you.

Click **Test connection** to confirm Aikido can reach Sumo Logic, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Sumo Logic integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Sumo Logic (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Sumo Logic, search the collector's source. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up in your search results.

{% hint style="info" %}
If Aikido created the collector for you, it's named in the format `Monad-<pipeline_id>-<node_id>`, making it easy to spot among your existing collectors.
{% endhint %}

### Managing the integration

From **Settings > Integrations > Sumo Logic** in Aikido, you can:

* **Manage Integration** - update the URL, or Access ID and Access Key.
* **Remove Integration** - disconnect Sumo Logic. This stops all future syncs; it doesn't delete data already sent to Sumo Logic, or remove any collector Aikido created.
