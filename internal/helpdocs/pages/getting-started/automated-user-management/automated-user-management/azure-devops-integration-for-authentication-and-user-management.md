# Azure DevOps: Authentication and User Management

If your organization uses Azure DevOps, users can login with Google and Microsoft accounts. To allow auto-onboarding of users in your workspace: configure **Trusted Domains (see below for instructions)**

### Understanding Aikido's Azure DevOps Integration <a href="#understanding-aikidos-azure-devops-integration" id="understanding-aikidos-azure-devops-integration"></a>

* **Manual Onboarding:** invite users manually via email
* **Auto-onboard via Trusted Domains**: Users can automatically join the Azure DevOps workspace if their login email is part of a trusted domain that you can specify on workspace level. Aikido will verify this user has access to your Azure Devops organization. **Note:** the user needs to be a member on the organisation level in Azure, otherwise they will not be recognised during team sync.
* **Synchronization of Teams and Repositories:** Aikido replicates your Azure DevOps team and project structure. All users will have access to their repos, in line with the permissions set in Azure DevOps. By default, all users will have the **Team Only** role.

{% hint style="info" %}
If users, roles, or repository access change in Azure DevOps, Aikido syncs those updates automatically. Changes are synced nightly, or you can trigger a manual sync via the Teams page in settings.
{% endhint %}

### Onboarding of Users with Trusted Domains <a href="#onboarding-of-users-with-trusted-domains" id="onboarding-of-users-with-trusted-domains"></a>

{% hint style="info" %}
If you have multiple workspaces, you need to setup Trusted Domains in each workspace.
{% endhint %}

1. Go to [General Settings](https://app.aikido.dev/settings/account) in your workspace
2. In workspace info, click 'Add Trusted Domain'

   [image: Azure DevOps Server: Update token and add trusted domains for security.]
3. Fill in the trusted domain in the modal

   [image: Add a trusted domain for Azure DevOps user auto-enrollment and verification.]

{% hint style="info" %}
For security reasons, Aikido only allows you to add trusted domains that are the same as the current logged in user. This means that <user@aikido.dev> can only add [aikido.dev](http://aikido.dev) as trusted domain.
{% endhint %}

### Manually Inviting Users <a href="#manually-inviting-users" id="manually-inviting-users"></a>

**Manually invite via email**: You can invite users via the Aikido platform on the specified email, and the user will be able to access Aikido directly.
