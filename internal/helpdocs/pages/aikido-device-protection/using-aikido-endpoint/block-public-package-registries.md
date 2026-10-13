# Block Public Package Registries

When you route package traffic through a private registry, developers can still bypass the setup by pointing their package manager directly at the public registry. Device Protection enforces that all traffic goes through your approved private registry instead.

## How it works

Device Protection monitors network traffic on each device. When a request goes directly to a public package registry, it is blocked at the network level before the package is downloaded.

## Prerequisites

* Aikido Device Protection deployed on developer machines. See [Deploying Device Protection](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides.md).
* Your private registry added under [Custom Registries](/aikido-device-protection/using-aikido-endpoint/custom-registries.md), so installs through it keep working once the public registry is blocked.

## Configure the block

1. Go to [Device Protection and open the Settings tab](https://app.aikido.dev/endpoint-protection/settings).
2. Click the ecosystem you want to block.
3. Turn on **Block All Installs**.

Repeat for every ecosystem you route through a private registry.

{% hint style="warning" %}
It can take up to 10 minutes for changes to sync to all devices.
{% endhint %}

Need to let a specific package or team through? Add an exception under the same ecosystem. See [Configuring Installation Policies](/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies.md#exceptions).
