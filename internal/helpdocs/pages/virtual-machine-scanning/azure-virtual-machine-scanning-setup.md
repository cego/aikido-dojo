# Azure Virtual Machine Scanning Setup

{% hint style="info" %}
This functionality is available for **Pro** and **Advanced** plans only. [Contact us](https://www.aikido.dev/contact) for more information.
{% endhint %}

### Why should I scan my virtual machines? <a href="#why-should-i-scan-my-virtual-machines" id="why-should-i-scan-my-virtual-machines"></a>

With virtual machine scanning, Aikido scans the disks of your virtual machines for vulnerable packages and outdated runtimes.

### Prerequisites

You'll need an Entra ID app registration with a client secret. If you've already connected your Azure cloud to Aikido for cloud scanning, you can reuse that same app registration, no need to create a new one. Otherwise, set that up first via [Connect your Azure cloud](https://help.aikido.dev/cloud-scanning/connect-your-cloud/azure/connect-azure-cloud-account-to-aikido).

### Create a custom role <a href="#getting-started" id="getting-started"></a>

Aikido uses a custom Azure role scoped to the minimum permissions needed to scan your virtual machines.

{% stepper %}
{% step %}

### Go to Access control (IAM)

Navigate to **Subscriptions** (or your **management group**), find the relevant scope for your virtual machines, and click **Access control (IAM)**.
{% endstep %}

{% step %}

### Add a custom role

Click **Add**, then **Add custom role**.
{% endstep %}

{% step %}

### Choose your scope

Decide whether the role should be assignable at **subscription level** or **management group level**:

* **Single subscription:** covers VMs in one subscription only.
* **Management group:** covers all subscriptions underneath it, useful if your VMs are spread across multiple subscriptions.

{% hint style="info" %}
For management group scope, you don't need to look up a separate ID. Use your **Directory (tenant) ID**, the same one from your app registration overview page. This maps to your organization's root management group.
{% endhint %}
{% endstep %}

{% step %}

### Paste the role JSON

Go to the **JSON** tab and click **Edit**. Paste the role definition below.

Set `assignableScopes` based on the scope you chose:

For a single subscription:

```json
"assignableScopes": [
  "/subscriptions/<your-subscription-id>"
]
```

For a management group:

```json
"assignableScopes": [
  "/providers/Microsoft.Management/managementGroups/<your-directory-tenant-id>"
]
```

Full role definition:

```json
{
  "properties": {
    "roleName": "Aikido VM Scanner",
    "description": "",
    "assignableScopes": [
      "<insert scope from above>"
    ],
    "permissions": [
      {
        "actions": [
          "Microsoft.Compute/virtualMachines/read",
          "Microsoft.Compute/disks/beginGetAccess/action",
          "Microsoft.Compute/disks/endGetAccess/action",
          "Microsoft.Compute/disks/read",
          "Microsoft.Compute/snapshots/read",
          "Microsoft.Compute/snapshots/write",
          "Microsoft.Compute/snapshots/delete",
          "Microsoft.Compute/snapshots/beginGetAccess/action",
          "Microsoft.Compute/snapshots/endGetAccess/action",
          "Microsoft.Authorization/roleAssignments/read"
        ],
        "notActions": [],
        "dataActions": [],
        "notDataActions": []
      }
    ]
  }
}
```

Click **Save**.
{% endstep %}

{% step %}

### Review and create

At the bottom, click **Review + assign**, then **Create**.
{% endstep %}
{% endstepper %}

### Assign the role to your app registration

{% stepper %}
{% step %}

### Go back to Access control (IAM)

Navigate to the same scope (subscription or management group) and click **Access control (IAM)**.
{% endstep %}

{% step %}

### Add a role assignment

Go to the **Role assignments** tab, click **Add**, then **Add role assignment**.
{% endstep %}

{% step %}

### Select the custom role

In the **Role** tab, search for and select the **Aikido VM Scanner** role, then click **Next**.
{% endstep %}

{% step %}

### Assign it to your app registration

Leave **Assign access to** on its default value. Click **Select members**, search for your app registration name (e.g. "AikidoSecurity"), and select it. Click **Select**, then **Review + assign** twice.
{% endstep %}
{% endstepper %}

{% hint style="info" %}
Your app registration now has the permissions it needs to scan your Azure virtual machines.
{% endhint %}

### VM grouping

To optimize scanning efficiency, Aikido groups certain Azure VMs and scans only one instance from each group. Grouping works as follows:

* **Virtual Machine Scale Sets (VMSS)**: All VMs that belong to the same VMSS are shown as a single VM group in Aikido. The VM group's name matches the VMSS identifier.
* **No grouping**: VMs that aren't part of a VMSS are treated as standalone VMs and scanned individually.
