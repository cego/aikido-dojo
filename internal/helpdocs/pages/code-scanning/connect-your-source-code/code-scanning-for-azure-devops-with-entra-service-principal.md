# Code Scanning for Azure DevOps with an Entra Service Principal

Let Aikido scan your Azure DevOps organization as a Microsoft Entra service principal instead of a personal access token.

## Code Scanning for Azure DevOps with an Entra Service Principal

Instead of a personal access token (PAT), you can let Aikido scan your Azure DevOps organization as a Microsoft Entra service principal that you own. No token or secret is created or stored: for every scan, Aikido proves its identity to Microsoft Entra, which issues a short-lived access token for your service principal. Nothing expires and nothing needs rotating.

{% hint style="info" %}
Your workspace must already be connected to Azure DevOps Services (dev.azure.com) with a PAT, and your Azure DevOps organization must be connected to a Microsoft Entra tenant (**Organization settings** → **Microsoft Entra**). You need to be able to register applications in that tenant and to add users to your Azure DevOps organization. Azure DevOps Server, the Aikido broker, and the local scanner keep using a PAT.
{% endhint %}

**First**, go to [General Settings](https://app.aikido.dev/settings/account), click **Update Personal Access Token** and, on the **Code Scanning** tab, click **Use An Entra Service Principal Instead**.

[image: Code Scanning tab in Aikido with the personal access token card and the button to use an Entra service principal instead]

You end up on the service principal management page, which can also be found at [Manage Azure DevOps Scanning service principal](https://app.aikido.dev/settings/integrations/azure-devops/scanning/service-principal). Keep this page open: it shows the **Issuer**, **Subject identifier**, and **Audience** you need in the next section.

[image: Aikido page for the scanning service principal showing the issuer, subject identifier, and audience values]

{% hint style="info" %}
Create a separate application for code scanning, even if you already set up [AutoFix with an Entra service principal](/autofix-and-remediation/configure/autofix-for-azure-devops-with-entra-service-principal.md). The Issuer and Audience are the same, but the Subject identifier is different because scanning uses its own Aikido identity, and the AutoFix application has write access to your projects.
{% endhint %}

**Next**, follow the instructions below to create the service principal in Microsoft Entra and give it read access in Azure DevOps.

### Register the application in Microsoft Entra

1. Open the [Microsoft Entra admin center](https://entra.microsoft.com), go to **Entra ID** → **App registrations** and click **New registration**.
2. Give the application a name, for example **Aikido Code Scanning**, keep **Single tenant only** selected, leave the redirect URI empty and click **Register**.

[image: Register an application form in Microsoft Entra with the name Aikido Code Scanning and single tenant selected]

1. On the **Overview** page of the application, copy the **Application (client) ID** and the **Directory (tenant) ID**. You insert both in Aikido later.
2. Go to **Certificates & secrets** → **Federated credentials** and click **Add credential**.

[image: Federated credentials tab of the Entra application with the Add credential button]

1. Select **Other issuer** and fill in the values exactly as shown on the service principal management page in Aikido:

* **Issuer**: the Issuer shown in Aikido
* **Type**: Explicit subject identifier
* **Value**: the Subject identifier shown in Aikido
* **Audience**: click **Edit** and replace the default with the Audience shown in Aikido. It contains your workspace id, so only tokens Aikido mints for your workspace are accepted. Give the credential a name, for example **aikido-code-scanning**, and click **Add**. You don't need a client secret or certificate.

[image: Add a credential form in Microsoft Entra filled in with the issuer, subject identifier, and audience from Aikido]

### Grant read access in Azure DevOps

1. In Azure DevOps, go to **Organization settings** → **Users** and click **Add users**.
2. Search for the application name (**Aikido Code Scanning**), set the **Access level** to **Basic**, and add it to the projects Aikido should scan. If the application doesn't show up yet, wait a minute and search again.
3. Under **Azure DevOps Groups**, select **Project Readers**. Then click **Add**.

[image: Add new users panel in Azure DevOps with the service principal, Basic access level, a project, and the Project Readers group selected]

{% hint style="warning" %}
The default group is **Project Contributors**, which grants write access. Select **Project Readers** so the service principal can only read your code.
{% endhint %}

### Connect in Aikido

1. Go back to the service principal management page in Aikido.
2. Insert the **Directory (tenant) ID** and the **Application (client) ID** you copied earlier and click **Save Service Principal**. Aikido requests a token from Microsoft Entra and verifies that it can access the repositories of your Azure DevOps organization before saving.

You are all set. Scanning now authenticates as the service principal, and the **Code Scanning** tab shows which application is in use.

[image: Code Scanning tab in Aikido showing the configured Entra service principal with Manage and Remove buttons]

{% hint style="info" %}
Your PAT stays stored but is no longer used. To switch back, click **Remove** on the **Code Scanning** tab: scanning then uses the stored PAT again, so update it first if it has expired. Deleting the federated credential in Microsoft Entra or removing the service principal from your Azure DevOps organization revokes the service principal only.
{% endhint %}
