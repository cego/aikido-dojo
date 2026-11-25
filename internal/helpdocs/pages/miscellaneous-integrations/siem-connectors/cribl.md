# Cribl

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Cribl for centralized monitoring, and route, enrich, or filter them before they reach your SIEM or data lake.

This is useful if your team already uses Cribl Stream as a telemetry pipeline and wants Aikido findings to flow through it alongside the rest of your security data.

[image: Cribl card in the Aikido app with the description Send Aikido issues to Cribl for centralized monitoring]

## Prerequisites

Before connecting, get these from Cribl:

* **Cribl HTTP source** - a Cribl HTTP source that is enabled and reachable from Aikido. [Allow Aikido's egress IP addresses](https://help.aikido.dev/miscellaneous-integrations/siem-connectors/siem-egress-ip-addresses) if your Cribl instance restricts inbound traffic.
* **Ingress address and port** - the hostname of your Cribl instance and the port the HTTP source listens on (default `10080`).
* **API path** - the HTTP event API path set on the source.
* **Authentication token** - a token configured on the HTTP source.

## Set up Cribl

{% stepper %}
{% step %}
### Create an HTTP source in Cribl

1. In Cribl Stream, go to **Data** → **Sources** and click **Add Source** next to **HTTP**.
2. Enter an **Input ID** (for example, `aikido`).
3. Set the **Port**. The default is `10080`.
4. Under **Auth tokens**, click **Add Token** and enter a secret value. Copy it, you'll paste it into Aikido later.
5. Note the **Cribl HTTP event API** path. The default is `/cribl`.
6. Click **Save**, then commit and deploy your changes.
{% endstep %}

{% step %}
### Connect a SIEM

Go to **Settings** → **Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Cribl** and fill in:

**Settings**

* **Ingress Address** - the hostname where your Cribl instance is accessible, for example `example.cribl.cloud`. You can find it in your group information panel in Cribl. Don't include the port or path.
* **Port** - the port your HTTP source listens on. Pre-filled with `10080`.
* **API Path** - the path set on your HTTP source's HTTP event API, for example `/cribl`. Don't add `_bulk`, Aikido adds it.

**Authentication**

* **Authentication Token** - the token you created on the Cribl HTTP source. Aikido stores it encrypted and never returns it after you save.

Click **Test connection**, then **Connect**.
{% endstep %}

{% step %}
### Configure which issues sync

On the Cribl integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Cribl (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

## Verify events are arriving

In Cribl Stream, open your HTTP source and go to the **Live Data** tab to capture incoming events. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm events arrive.

Events are sent as newline-delimited JSON (NDJSON). For the fields in each event, see the [SIEM field reference](https://help.aikido.dev/miscellaneous-integrations/siem-connectors/siem-field-reference).

## Troubleshooting

* **Connection errors** - check the ingress address, port, and API path, and confirm your Cribl instance is reachable from Aikido. If you restrict inbound traffic, allow Aikido's [egress IP addresses](https://help.aikido.dev/miscellaneous-integrations/siem-connectors/siem-egress-ip-addresses).
* **Authentication errors** - confirm the **Authentication Token** matches a token on the Cribl HTTP source, and that the token is enabled.
* **No events in Live Data** - confirm you committed and deployed the source in Cribl, and that the source isn't disabled.
* **Timeouts** - check that your Cribl instance is healthy and not overloaded.

If the problem persists, [contact support](https://help.aikido.dev/miscellaneous-info/contact-support).

## Manage the integration

From **Settings** → **Integrations** → **Cribl** in the Aikido app, you can:

* **Manage Integration** - update the ingress address, port, API path, or authentication token.
* **Remove Integration** - disconnect Cribl. This stops all future syncs. It doesn't delete events already sent to Cribl.
