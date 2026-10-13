# SIEM Egress IP Addresses

When Aikido sends issue data to your connected SIEM, that traffic originates from a fixed set of Aikido-owned IP addresses.

If your SIEM destination sits behind a firewall or IP allowlist, add these addresses so Aikido can connect.

### Egress IPs

* 35.81.245.33
* 35.83.26.83
* 35.84.110.88

{% hint style="info" %}
This page is the canonical list of Aikido's SIEM egress IPs. Confirm against it before configuring an allowlist, as the set can change over time.
{% endhint %}

### When you need them

Add the addresses above to your SIEM's allowlist if it restricts inbound access by IP, for example a self-managed Splunk HEC, Elasticsearch cluster, or QRadar HTTP Receiver sitting behind a firewall.

See [SIEM Connectors](/miscellaneous-integrations/siem-connectors.md) for setup steps for each SIEM.
