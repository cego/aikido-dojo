# Allowing IP Addresses for Cloud Scanning

Aikido’s Cloud Scanning checks your cloud environments for misconfigurations (CSPM).

These scans are performed from a fixed set of Aikido IP addresses. To allow scanning, you may need to allowlist these IPs in your cloud provider.

{% hint style="warning" %}
This applies only to cloud misconfiguration scanning (CSPM).

Other Aikido features, such as [container scanning](/code-scanning/miscellaneous/allowing-ip-addresses-for-code-container-scanning.md), [code scanning](/code-scanning/miscellaneous/allowing-ip-addresses-for-code-container-scanning.md), or [api scanning](/dast-surface-monitoring/miscellaneous/allowing-ip-addresses-for-dast-surface-monitoring.md), do not use these IPs.
{% endhint %}

## Why allowlist our IPs?

Some cloud providers or firewall setups block requests from unknown sources. Adding Aikido’s IPs ensures our scanner can reach your environment and report on misconfigurations.

Common use cases in which you need to allowlist IP addresses include:

* Access to specific Azure resources, such as [Key Vaults](https://learn.microsoft.com/en-us/azure/key-vault/general/network-security#key-vault-firewall-enabled-ipv4-addresses-and-ranges---static-ips).
* 3rd-party monitoring systems detecting cloud API requests.

## Which IPs should you allowlist?

Cloud scanning traffic comes from the region where **your Aikido workspace** is hosted, not from the region where your cloud resources run. Check the URL you use to log in to Aikido, then allowlist that region's IPs:

| Aikido URL                                     | Region | IP list               |
| ---------------------------------------------- | ------ | --------------------- |
| [app.aikido.dev](https://app.aikido.dev)       | EU     | EU-based IP addresses |
| [app.us.aikido.dev](https://app.us.aikido.dev) | US     | US-based IP addresses |
| [app.au.aikido.dev](https://app.au.aikido.dev) | AU     | Coming soon           |
| [app.me.aikido.dev](https://app.me.aikido.dev) | ME     | ME-based IP addresses |

{% hint style="warning" %}
If you log in at `app.aikido.dev`, scanning uses **EU** IPs, even when your cloud resources are hosted in the US. Allowlisting only US IPs in that case will cause scans to fail.
{% endhint %}

## Current IP addresses

### EU-based IP addresses (`app.aikido.dev`)

* 34.240.182.181
* 52.31.217.247

### US-based IP addresses (`app.us.aikido.dev`)

* 98.80.32.100
* 52.54.124.156

### ME-based IP addresses (`app.me.aikido.dev`)

* 40.172.244.108
* 3.29.255.25
