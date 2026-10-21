# IBM QRadar

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to IBM QRadar for centralized monitoring, so your security findings are searchable alongside the rest of your logs.

This is useful if your team already uses QRadar as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, make sure you have:

* A QRadar deployment reachable from Aikido on the port its HTTP Receiver listens on.
* Administrative access in QRadar to create or edit a log source.
* If using HTTPS with a self-signed certificate, the **CA certificate** that signed the receiver's server certificate.
* If using HTTPS with mutual TLS, a **client certificate** and matching **private key**.

### Connecting QRadar

{% stepper %}
{% step %}

#### Set up the HTTP Receiver log source in QRadar

1. In QRadar, go to **Admin > Data Sources > Log Sources > Add**.
2. Set **Log Source Type** to whatever fits your events (or a generic type like `Universal DSM`), and set **Protocol Configuration** to `HTTP Receiver`.
3. Configure the receiver:
   * **Listen Port** - any free port on the QRadar host (default `12469`). You'll enter the same port in Aikido.
   * **Communication Type** - `HTTP`, `HTTPS`, or `HTTPS and Client Authentication (mTLS)`. This must match what you select in Aikido.
   * **Event Parsing Method** - `Event Per Line` for newline-delimited JSON, or `Event Per JSON Array` for one JSON array per request. This must also match what you select in Aikido.
   * *(Optional)* Turn on **Use HTTP Authentication Token Header** and set a header name and value, if you want to require an auth header on every request.
4. Save and deploy the log source.
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **IBM QRadar** and fill in:

**Settings**

* **Host** - the hostname or IP of your QRadar HTTP Receiver. Don't include `http://` or `https://`.
* **Port** - the port the HTTP Receiver listens on (QRadar's default is `12469`).
* **Communication Type** - must match the log source's Communication Type:
  * **HTTP** - plaintext. Only recommended on trusted networks.
  * **HTTPS** - TLS with server certificate verification. Provide a **CA Certificate**, or enable **Insecure Skip Verify** for testing (you can't use both together).
  * **HTTPS with mTLS** - TLS plus a client certificate. Provide the same **CA Certificate** / **Insecure Skip Verify** choice as HTTPS, plus a **Client Certificate** and **Client Key**.
* **Event Parsing Method** - must match the log source's Event Parsing Method:
  * **Event Per Line** - newline-delimited JSON, one event per line.
  * **Event Per JSON Array** - a single JSON array of events per request.
* **Authentication** *(optional)* - leave unset for an unauthenticated log source, or choose:
  * **HTTP Basic** - enter a **Username** and **Password**.
  * **Custom Header** - enter a **Header Name** and **Header Value** (use this if QRadar's **Use HTTP Authentication Token Header** option is on).

Click **Test connection** to confirm Aikido can reach QRadar, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the QRadar integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to QRadar (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In QRadar, check the log source's event count, or search for recently indexed events. Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up.

### Troubleshooting

* **Connection refused or timeouts** - double check **Host** and **Port**, and confirm firewalls allow traffic between Aikido and QRadar on that port.
* **TLS verification failures** - make sure the **CA Certificate** matches the CA that signed QRadar's server certificate. You can temporarily enable **Insecure Skip Verify** to check whether the issue is verification-specific, but don't leave that on in production.
* **mTLS handshake failures** - confirm the client certificate is trusted by the log source's trust store, and that **Client Certificate** and **Client Key** are a valid matching pair.
* **Events not appearing** - confirm **Event Parsing Method** in Aikido matches the log source's setting. A mismatch can cause QRadar to misparse or drop events.
* **Authentication failures (401)** - if QRadar's log source has **Use HTTP Authentication Token Header** on, use **Custom Header** in Aikido with matching name/value. For HTTP Basic, double check the credentials on the log source.

### Managing the integration

From **Settings > Integrations > IBM QRadar** in Aikido, you can:

* **Manage Integration** - update the host, port, communication settings, or authentication.
* **Remove Integration** - disconnect QRadar. This stops all future syncs; it doesn't delete events already sent to QRadar.
