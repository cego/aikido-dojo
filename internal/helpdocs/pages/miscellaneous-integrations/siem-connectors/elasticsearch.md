# Elasticsearch

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Elasticsearch for centralized monitoring, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses Elasticsearch (or Elastic Cloud) as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Elastic:

* **Username and Index Name** - always required.
* **Authentication** - either an **API Key**, or a **Password** for the username above.
* **Connection details** - either a direct **URL** to your Elasticsearch cluster, or a **Cloud ID** if you're on Elastic Cloud.

You'll pick one option from each of the last two groups when connecting, giving four possible combinations (for example, URL + API Key, or Cloud ID + Password).

### Connecting Elasticsearch

{% stepper %}
{% step %}

#### Get your Elasticsearch credentials

Log in to your [Elastic Cloud console](https://cloud.elastic.co/home) and gather the details for the connection type you plan to use:

* **URL** - find your Elasticsearch endpoint URL from your deployment overview.
* **Cloud ID** - find this on your deployment's overview page, if you're connecting via Elastic Cloud.
* **API Key** - generate one from **Stack Management > API Keys**.
* **Username and Password** - use an existing user, or create one with permissions to write to your index.
* **Index Name** - the index you want Aikido to write to. If it doesn't already exist, Elasticsearch creates it automatically the first time data is sent.
  {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Elasticsearch** and fill in:

**Settings**

* **Username** - the username for authenticating with your Elasticsearch cluster.
* **Index Name** - the index to write data to. Created automatically if it doesn't exist.
* **Authentication Type** - choose how to authenticate:
  * **API Key** - enter your **API Key**.
  * **Password** - enter the **Password** for the username above.
* **Connection Type** - choose how to connect:
  * **URL** - enter the direct **URL** of your Elasticsearch cluster.
  * **Cloud ID** - enter your Elastic Cloud deployment's **Cloud ID**.

**Options**

* **Insecure Skip Verify** - skips verification of the server's TLS certificate. Insecure, use for testing only.

Click **Test connection** to confirm Aikido can reach Elasticsearch, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Elasticsearch integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Elasticsearch (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Kibana (or by querying Elasticsearch directly), search your configured index. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up as a new document in the index.

### Managing the integration

From **Settings > Integrations > Elasticsearch** in Aikido, you can:

* **Manage Integration** - update the username, index, authentication, or connection details.
* **Remove Integration** - disconnect Elasticsearch. This stops all future syncs; it doesn't delete documents already sent to Elasticsearch.
