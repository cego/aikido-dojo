# Wazuh

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to the Wazuh indexer for centralized monitoring, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses Wazuh and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

{% hint style="warning" %}
This connector writes directly to the **Wazuh indexer**, not the Wazuh server/manager. That means Aikido issues do **not** pass through Wazuh's decoders and ruleset: they won't trigger Wazuh rules or alerts, won't get rule levels, descriptions, or MITRE ATT\&CK enrichment, and won't appear in the Wazuh dashboard's Security Events views. They are fully searchable in Discover and usable in dashboards through an index pattern, just not part of Wazuh's alerting pipeline.
{% endhint %}

### Requirements

Before connecting, make sure you have:

* **Wazuh indexer URL** - the indexer API endpoint (port `9200` by default), reachable from Aikido. This targets the indexer, not the Wazuh server/manager.
* **Indexer username and password** - a user with permission to create and write to the target index, including bulk indexing.
* **A dedicated index name** - don't use Wazuh's own `wazuh-alerts-*` indices, they're managed by the Wazuh server and expect a schema this connector doesn't produce.
* **CA Certificate** - required unless your indexer presents a publicly-trusted certificate. A default Wazuh install uses a self-signed certificate.

### Connecting Wazuh

{% stepper %}
{% step %}

#### Get your Wazuh indexer details

1. Confirm your Wazuh indexer API is reachable from Aikido. If it's behind a firewall or IP allowlist, you'll need to allow Aikido's egress IPs.
2. Create or choose an indexer user with write and bulk-indexing permissions on the index you plan to use.
3. Pick a dedicated index name for Aikido data (don't reuse `wazuh-alerts-*`).
4. If your indexer uses a self-signed or internally-issued certificate (the Wazuh default), get the CA certificate in PEM format. On a default Wazuh Docker deployment it's at `config/wazuh_indexer_ssl_certs/root-ca.pem`.

{% hint style="warning" %}
The hostname in your indexer URL must appear as a Subject Alternative Name (SAN) on the indexer's certificate. A default Wazuh certificate is only issued for `wazuh.indexer`, so pointing Aikido at a different hostname, even a public one, means reissuing the certificate to include that name. A mismatch here fails with a hostname-related error that can look like an unrelated problem. You can check what a certificate covers with:

```
openssl x509 -in config/wazuh_indexer_ssl_certs/wazuh.indexer.pem -noout -ext subjectAltName
```

{% endhint %}
{% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Wazuh** and fill in:

**Settings**

* **URL** - the URL of the Wazuh indexer API, must start with `https`.
* **Username** - a Wazuh indexer user with write permissions on the target index.
* **Index** - the index to write to. Use a dedicated index, not one of Wazuh's own `wazuh-alerts-*` indices.

**Options**

* **Insecure Skip Verify** - skips TLS certificate verification. Not recommended for production, supply a **CA Certificate** instead to verify a self-signed indexer. Can't be combined with **CA Certificate**.

**Authentication**

* **Password** - the password for the Wazuh indexer user.
* **CA Certificate** - the CA certificate (PEM format) used to verify the indexer's certificate. Required unless the indexer has a publicly-trusted certificate.

Click **Test connection** to confirm Aikido can reach your Wazuh indexer, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Wazuh integration page, under **Issue input settings**, configure which Aikido issues are synced:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Wazuh (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In the Wazuh dashboard, go to **Discover**, create an index pattern for your dedicated index if you haven't already, and search it. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up. Remember: it won't appear under Security Events, since it bypasses Wazuh's ruleset.

### Managing the integration

From **Settings > Integrations > Wazuh** in Aikido, you can:

* **Manage Integration** - update the URL, index, or credentials.
* **Remove Integration** - disconnect Wazuh. This stops all future syncs; it doesn't delete data already sent to your indexer.
