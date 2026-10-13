# Splunk

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Splunk so your security findings live alongside the rest of your monitoring and observability data.

This is useful if your team already uses Splunk Cloud / Splunk Enterprise as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Splunk:

* **HEC URL** - the base URL of your Splunk HTTP Event Collector (HEC) endpoint.
* **HEC Port** - the port your Splunk HEC listens on.
* **HEC Token** - the token you generate in Splunk for authentication (see below).
* **Username and password** *(optional)* - only needed if you want Aikido to automatically create the Splunk index for you.

### Connecting Splunk

{% stepper %}
{% step %}

#### Get your Splunk HEC token

If you're on **Splunk Enterprise**, first enable HEC: go to **Settings > Data Inputs > HTTP Event Collector > Global Settings**, and set **All Tokens** to **Enabled**.

Then, on both Splunk Cloud and Splunk Enterprise, create a token:

1. Go to **Settings > Add Data**, then click **Monitor**, then **HTTP Event Collector**.
2. Enter a **Name** for the token (for example, "Aikido").
3. Optionally set a source name, description, and index for the token.
4. Click **Next**, review the settings, and click **Submit**.
5. Copy the token value shown, you'll paste it into Aikido in the next step.

Note the HEC URL and port for your Splunk instance too; you'll need both.
{% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Splunk** and fill in:

**Settings**

* **Splunk HEC URL** - the base URL of your Splunk HTTP Event Collector.
* **Splunk HEC Port** - the HEC port.

**Options** *(all optional)*

* **Allow Insecure Connection** - allows insecure connections (for example, self-signed certificates). Not recommended for production.
* **Splunk Index** - the index to send data to. Leave empty to use the default index.
* **Username** - an administrative username, only needed if you want Aikido to create the index for you.
* **To Create Index** - check this if you want Aikido to create the index automatically. Leave unchecked if you're using an existing index. If checked, make sure your HEC token isn't already scoped to specific indices in Splunk, or the connection will error.

**Authentication**

* **Splunk HEC Token** - the token you generated in Splunk.
* **Password** - only needed if **To Create Index** is checked.

Click **Test connection** to confirm Aikido can reach Splunk, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Splunk integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Splunk (for example, every hour).

Click **Save Settings**.\
\

{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Splunk, search your configured index (or `source=aikido` if using the default). Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up in Splunk.

### Managing the integration

From **Settings > Integrations > Splunk** in Aikido, you can:

* **Manage Integration** - update the HEC URL, port, options, or token.
* **Remove Integration** - disconnect Splunk. This stops all future syncs; it doesn't delete events already sent to Splunk.
