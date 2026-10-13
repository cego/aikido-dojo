# Microsoft Sentinel

{% hint style="info" %}
This SIEM connector is currently available by request. [Reach out to support](https://www.aikido.dev/contact) to have it enabled for your workspace.
{% endhint %}

Send your Aikido issues to Microsoft Sentinel so your security findings live alongside the rest of your Azure monitoring and observability data.

This is useful if your team already uses Sentinel as a SIEM and wants a single place to correlate Aikido findings with infrastructure, application, and audit signals.

### Requirements

Before connecting, get these from Azure:

* **Endpoint URL** - the Logs Ingestion URI of your Data Collection Endpoint.
* **Stream Name** - the name of the data stream defined on your Data Collection Rule (DCR), for example `Custom-Aikido_CL`.
* **DCR ID** - either a Rule ID (manual mode) or a DCR Resource ID (managed mode). See DCR Configuration below.
* **Tenant ID, Client ID, Client Secret** - from a Microsoft Entra application registered in your tenant.

{% hint style="info" %}
An onboarded Log Analytics workspace with at least contributor rights is required, and your Entra application needs the **Monitoring Metrics Publisher** role on the DCR to send data. For managed mode, the registered application also needs permission to read, create, and delete Data Collection Rules under the workspace's subscription.
{% endhint %}

### Connecting Sentinel

{% stepper %}
{% step %}

#### Set up your Azure prerequisites

1. Create a Microsoft Entra application and note the Tenant ID, Client ID, and a Client Secret.
2. Create a Data Collection Endpoint and copy its Logs Ingestion URI, this is your Endpoint URL.
3. Create a table in your Log Analytics workspace to receive the ingested data. Upload a sample JSON payload matching your issue data when creating the table so Azure can derive its schema.
4. Create a Data Collection Rule (DCR) under **Monitor > Data Collection Rules > Create**, using the Data Collection Endpoint from step 2. Add a data source of type **Custom**, upload the same sample JSON payload to generate its stream declaration, and add a data flow that routes that stream to the table from step 3. This is the DCR you'll reference below, either directly (manual mode) or as the template Aikido clones (managed mode).
5. On the DCR's **Access control (IAM)** page, assign the **Monitoring Metrics Publisher** role to your Entra application, scoped to the DCR. This data-plane role is required for the Logs Ingestion API to accept data from the application; workspace Contributor rights alone aren't enough and won't grant it.
6. Decide on a DCR mode:
   * **Manual** - you manage a single DCR yourself. Find its Rule ID under **Monitor > Data Collection Rules**, select your DCR, then **JSON View > immutableId**.
   * **Managed** - Aikido scales DCRs for you automatically for higher throughput. Find the DCR Resource ID the same way, using the top-level `id` field instead.
7. Find your Stream Name under the DCR's `properties.streamDeclarations` (it follows the pattern `Custom-[TableName]_CL`).
   {% endstep %}

{% step %}

#### Connect a SIEM

Go to **Settings > Integrations** in the Aikido app and add a new SIEM integration. In the **Connect a SIEM** dialog, select **Microsoft Azure Sentinel** and fill in:

**Settings**

* **Endpoint URL** - the Azure Monitor DCR ingestion endpoint URL.
* **Stream Name** - the name of the data stream defined in the DCR.
* **DCR Configuration** - choose the DCR Mode:
  * **Manual** - ingest into a single DCR you provide.
  * **Managed** - Aikido clones your DCR to scale ingestion for higher throughput.
* **DCR Resource ID** *(managed mode)* - the full ARM resource ID of your DCR, in the form `/subscriptions/<sub>/resourceGroups/<rg>/providers/Microsoft.Insights/dataCollectionRules/<name>`. Aikido clones this to create a pod-local DCR.
* **Rule ID** *(manual mode)* - the immutable ID of your DCR.

**Authentication**

* **Tenant ID** - the Azure Active Directory tenant (directory) ID.
* **Client ID** - the application (client) ID registered in Azure Active Directory.
* **Client Secret** - the client secret associated with the registered application.

Click **Test connection** to confirm Aikido can reach Sentinel, then **Save changes**.
{% endstep %}

{% step %}

#### Configure which issues sync

On the Sentinel integration page, under **Issue input settings**, configure which Aikido issues are synced to your SIEM:

* **Issue status** - for example, sync open issues only, or include resolved and ignored issues too.
* **Minimum severity** - the lowest severity level to sync.
* **Issue type** - limit to specific issue types, or sync all of them.
* **Language** - limit to specific languages, or sync all of them.
* **Sync frequency** - how often Aikido pushes matching issues to Sentinel (for example, every hour).

Click **Save Settings**.
{% endstep %}
{% endstepper %}

### Verifying events are arriving

In Microsoft Sentinel, go to **Logs** and query your custom table (for example `Aikido_CL`). Wait for your configured sync frequency to pass, or trigger a test issue in Aikido, and confirm it shows up in the table.

### Troubleshooting

* **Records rejected for size** - individual records above 1 MB are rejected. This is an Azure Monitor Logs Ingestion API limit.
* **Connection succeeds but no data arrives (403 on ingestion)** - the Entra application authenticated but isn't authorized to publish to the DCR. Confirm it has the **Monitoring Metrics Publisher** role assigned directly on the DCR; workspace Contributor rights don't grant this.
* **Managed mode fails to set up a DCR** - confirm the registered application can read, create, and delete Data Collection Rules under the target subscription, and that the DCR Resource ID is valid.
* **Data not appearing** - verify the Stream Name matches a stream declared on the DCR, and that the table schema matches the issue data Aikido sends.

### Managing the integration

From **Settings > Integrations > Microsoft Sentinel** in Aikido, you can:

* **Manage Integration** - update the endpoint URL, stream name, DCR settings, or credentials.
* **Remove Integration** - disconnect Sentinel. This stops all future syncs; it doesn't delete events already sent to Sentinel.
