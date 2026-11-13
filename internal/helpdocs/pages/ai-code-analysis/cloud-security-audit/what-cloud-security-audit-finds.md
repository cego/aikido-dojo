# What Cloud Security Audit Finds

A Cloud Security Audit reasons about misconfigurations in the clouds you select and chains them into attack paths. When you add infrastructure as code (IaC) repositories, it can correlate a live setting with the template that defines it. Each finding includes evidence and remediation. See [Track Progress and Findings](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings) for how to read it.

For how this differs from rule-based checks, see [Cloud Security Audit](https://help.aikido.dev/ai-code-analysis/cloud-security-audit). To run an audit, see [Run a Cloud Security Audit](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit).

## Coverage at a glance

Cloud Security Audit looks for:

* **Misconfigurations** in the connected account's inventory.
* **Attack paths** formed by chaining those misconfigurations. For example, a publicly reachable resource and a broad role can be low risk alone and serious together.
* **IaC correlation**, when you include IaC repositories. A finding can point at the template or module that defines the live setting, so the fix can land in code as well as in the account.

**Configure** on the **Pricing** step changes how much of the selection gets investigated. It does not change which classes are in scope.

## Limits

Cloud Security Audit reads the inventory Aikido already collected. It does not change the account, and it does not inventory accounts you have not connected.

It does not replace [Cloud Scanning](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview), which checks known misconfiguration rules on every cloud scan. It also does not replace a [Code Security Audit](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/what-ai-code-audit-finds), which reasons about application source. IaC here is context for the live account.

You can retest a single Cloud Security Audit finding after you fix it, without running the whole audit again. See [Track Progress and Findings](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings).

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-code-analysis/cloud-security-audit)
{% endcontent-ref %}

{% content-ref url="../../cloud-scanning/cloud-scanning-overview.md" %}
[cloud-scanning-overview.md](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview)
{% endcontent-ref %}
