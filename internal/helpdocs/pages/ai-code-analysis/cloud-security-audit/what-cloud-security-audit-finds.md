# What Cloud Security Audit Finds

What Cloud Security Audit checks for in a connected cloud account and optional IaC.

A Cloud Security Audit reasons about misconfigurations in the clouds you select and chains them into attack paths. When you add infrastructure as code (IaC) repositories, it can correlate a live setting with the template that defines it. Each finding includes evidence and remediation. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md) for how to read it.

For how this differs from rule-based checks, see [Cloud Security Audit](/ai-code-analysis/cloud-security-audit.md). To run an audit, see [Run a Cloud Security Audit](/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit.md).

## Coverage at a glance

Cloud Security Audit looks for:

* **Misconfigurations** in the connected account's inventory.
* **Attack paths** formed by chaining those misconfigurations. For example, a publicly reachable resource and a broad role can be low risk alone and serious together.
* **IaC correlation**, when you include IaC repositories. A finding can point at the template or module that defines the live setting, so the fix can land in code as well as in the account.

**Configure** on the **Pricing** step changes how much of the selection gets investigated. It does not change which classes are in scope.

## Limits

Cloud Security Audit reads the inventory Aikido already collected. It does not change the account, and it does not inventory accounts you have not connected.

It does not replace [Cloud Scanning](/cloud-scanning/cloud-scanning-overview.md), which checks known misconfiguration rules on every cloud scan. It also does not replace a [Code Security Audit](/ai-code-analysis/ai-code-audit-overview/what-ai-code-audit-finds.md), which reasons about application source. IaC here is context for the live account.

You can retest a single Cloud Security Audit finding after you fix it, without running the whole audit again. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md).

## Related

{% content-ref url="/pages/k4yIjSNiww67YBu4pXls" %}
[Cloud Security Audit](/ai-code-analysis/cloud-security-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/tweiYLnrxouYwMmBjEiE" %}
[Cloud Scanning Overview](/cloud-scanning/cloud-scanning-overview.md)
{% endcontent-ref %}
