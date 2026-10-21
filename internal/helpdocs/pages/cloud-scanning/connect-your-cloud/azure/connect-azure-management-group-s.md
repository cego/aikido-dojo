# Connect Azure Management Group(s)

This page describes Aikido's support for automatically and continuously connecting Azure subscriptions within the same Azure management group(s).

{% hint style="info" %}
This functionality is available only for **Pro** and **Advanced** plans. **Contact us** via chat for more information.
{% endhint %}

When onboarding your Azure subscription, you have the option to connect entire Azure management groups automatically.

## Why Connect Azure Management Group(s)?

Connecting entire Azure Management Groups provides the following benefits:

* **Faster setup**: You only need to manually connect one subscription.
* **Automatic subscription discovery**: Azure subscriptions are automatically added to Aikido, including future ones.
* **Automatic ACR scanning setup**: Aikido will automatically discover and scan all your Azure Container Registries.
* **Automatic Azure VM scanning setup**: If enabled, Aikido will scan all your Azure VMs using our agentless technology.

## Prerequisites

* You are on the Pro, Advanced, or Enterprise plan.
* You have access to the management group(s) in your Azure environment.

## Getting Started

* Follow the standard flow of onboarding an Azure subscription.
* On step 3, select the "Management Groups" option.
* Provide the ID of the subscription that will act as a "management cloud" in Aikido. This is simply the subscription holding the connection details. There are no setup differences between this and the rest of the subscriptions. We recommend choosing one that you intend to keep for the foreseeable future.
* The other fields are optional. **If left empty, all Azure subscriptions from your tenant will be automatically connected.**
* Next, you will need to configure RBAC/IAM access to your Azure environment.

{% hint style="info" %}
You must grant access at a management group level. Granting access at the tenant root management group will automatically propagate to all management groups and subscriptions.

Alternatively, you can grant access at a specific (sub-)management group to onboard only a subset of your environment.
{% endhint %}

* For the CSPM components, Aikido requires Azure role `Security Reader` and `Log Analytics Reader`.
* For ACR scanning, Aikido requires `AcrPull` and `Container Registry Repository Catalog Lister`.
* For VM scanning, Aikido requires a role with the following permissions:

```
Microsoft.Compute/virtualMachines/read,
Microsoft.Compute/disks/beginGetAccess/action,
Microsoft.Compute/disks/endGetAccess/action,
Microsoft.Compute/disks/read,
Microsoft.Compute/snapshots/read,
Microsoft.Compute/snapshots/write,
Microsoft.Compute/snapshots/delete,
Microsoft.Compute/snapshots/beginGetAccess/action,
Microsoft.Compute/snapshots/endGetAccess/action,
Microsoft.Authorization/roleAssignments/read,
```

* Once you complete the setup, the other subscriptions will start to appear within a few minutes.

#### Cloud Purpose Determination <a href="#cloud-purpose-determination" id="cloud-purpose-determination"></a>

The purpose/environment of your Azure subscriptions is automatically determined from their names. Aikido looks for terms such as "production", "staging", "uat", etc., and sets the cloud purpose accordingly. If it doesn't find any match, the purpose will be "mixed". You can manually update the purpose of each cloud connection using the "Configure" button.
