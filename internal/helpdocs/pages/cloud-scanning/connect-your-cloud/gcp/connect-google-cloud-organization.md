# Connect Google Cloud Organization

Connect your entire Google Cloud Organization to Aikido by onboarding one of your GCP projects. Automatically discover and connect all other projects, simplify setup, and ensure complete coverage.

If you have multiple GCP projects part of the same Google Cloud organization, you can connect one of them, grant permission at the organization level, and let Aikido automatically connect the rest of them, including projects you will create in the future.

### Why Connect Your Google Cloud Organization? <a href="#why-connect-aws-organization" id="why-connect-aws-organization"></a>

By onboarding at the organization level, you benefit from:

* **Faster setup**: You only need to manually connect one project.
* **Automatic project discovery**: Projects are automatically added to Aikido, including the ones you will create in the future.
* **Automatic container scanning setup**: Aikido will automatically scan your containers from Artifact Registry from all your projects.

### Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* You are on the Pro, Advanced, or Enterprise plan.
* You have access to the Google Cloud organization, including the ability to grant permissions at the organization level.
* You designate one of your GCP projects to be your "management" project. You will create the service account/workload identity federation pool in this project, and onboard it to Aikido.

### Getting Started <a href="#getting-started" id="getting-started"></a>

To connect your Google Cloud organization, select the 'Full GCP Organization' option in the [GCP onboarding wizard](https://app.aikido.dev/clouds/add/gcp).

After you provide the ID and number of your "management" project, you will need to provide the following:

* **Organization ID**: It should look like '783352941112'.
* **Included Folder IDs (optional)**: This option enables you to onboard a subset of your organization - e.g., a folder and all its descendants.
* **Excluded Folder IDs (optional)**: This option enables you to exclude entire folders from the automatic onboarding. This also works in conjunction with the **Included Folders** option.
* **Excluded Project IDs**: Optionally, you can exclude specific GCP projects from being added to Aikido.
* You can also choose to automatically configure container scanning across your projects by enabling **Artifact Registry Scanning**.

You can obtain this information from your [Google Cloud Resource Manager page](https://console.cloud.google.com/cloud-resource-manager) (select the organization from the project selector).

{% hint style="info" %}
If the credential verification fails during onboarding, IAM permissions may still be propagating in your Google Cloud organization. Wait a few minutes and retry. Propagation time can vary depending on the size of your organization.
{% endhint %}

#### Editing an existing organization connection <a href="#editing-an-existing-connection" id="editing-an-existing-connection"></a>

Once your organization is connected, you can change its scope from **Clouds** → your GCP org → **Actions** → **Configure**. This opens the **Edit organization** modal, which presents the same include/exclude settings as the onboarding wizard, but combined into fewer fields:

* **Organization or Folder IDs**: accepts either your organization ID, or one or more folder IDs to scope the sync to a subset of your organization. Folder IDs here use the prefixed format.
* **Excluded Folder IDs**: excludes specific folders and everything beneath them. Folder IDs here use the plain format (no prefix).
* **Excluded Project IDs**: excludes specific projects.

{% hint style="info" %}
Folder ID format differs between the two fields above:

* **Organization or Folder IDs**: `folders/123456, folders/789012`
* **Excluded Folder IDs**: `123456, 987654`
  {% endhint %}

{% hint style="warning" %}
Narrowing the scope does not remove projects that were already connected. Any project that falls outside the new scope stays in Aikido until you manually delete it from the [clouds page](https://app.aikido.dev/clouds).
{% endhint %}

#### Terraform Module <a href="#terraform-module" id="terraform-module"></a>

If you prefer Infrastructure as Code, you can use Aikido's Terraform module to set up the required Google Cloud IAM and Workload Identity Federation configuration for organization-level onboarding. The module supports both single-project and full-organization setups, optional Artifact Registry scanning, and optional VM scanning permissions.

<https://github.com/AikidoSec/gcp-onboarding-terraform-module>

Use the modules/org module for organization-level onboarding, then upload the generated `credential_config_json` to Aikido during the cloud connection flow.

If your Aikido account is not on the EU instance (app.aikido.dev), set `aikido_region` to `us`, `me`, or `au` accordingly.

{% hint style="info" %}
**Important**

The Terraform module only sets up the required Google Cloud IAM configuration. You must still complete the onboarding flow in Aikido to connect the cloud account by uploading the generated `credential_config_json` and finishing the remaining in-app steps.
{% endhint %}

#### Cloud Purpose Determination <a href="#cloud-purpose-determination" id="cloud-purpose-determination"></a>

The purpose/environment of your GCP projects is automatically determined based on the name of the project. Aikido looks for terms such as "production", "staging", "uat", etc., and sets the cloud purpose accordingly. If it doesn't find any match, the purpose will be "mixed". You can manually update the purpose of each cloud connection using the "Configure" button.

### FAQs <a href="#faqs" id="faqs"></a>

* **Is it secure?**

Yes. Connecting your GCP Organization relies on the same mechanism we use for connecting individual GCP projects. [Workload Identity Federation](https://help.aikido.dev/cloud-scanning/connect-your-cloud/gcp/google-cloud-workload-identity-federation-setup) is also supported and recommended.

* **If I add a GCP project to my organization, will it appear in Aikido?**

Yes. Aikido scans your GCP organization every time it scans your "management" project, and automatically connects new GCP projects. There might be some delay until you see your project in Aikido, especially if your GCP organization has thousands of projects (we respect [the Google Cloud API limits](https://cloud.google.com/resource-manager/docs/limits)).

* **I just added a new GCP project to my org, and it did not show up in Aikido.**

If Aikido scanned your GCP "management" project (you can manually scan it) and the new project does not appear after 15 minutes, you may have reached the limit of clouds for your plan. Contact us to increase your limit.

* **What happens if I shut down a GCP project or move it to another org?**

Aikido will detect that the project is no longer active or part of the organization and will mark the corresponding cloud as "not reachable". You will see this on the [clouds page](https://app.aikido.dev/clouds).
