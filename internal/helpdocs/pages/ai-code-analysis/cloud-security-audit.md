# Cloud Security Audit

What Cloud Security Audit is, and how it differs from rule-based cloud scanning.

Cloud Security Audit is the cloud-focused audit type in [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md). It is Aikido's agentic approach to cloud security posture management (CSPM). Agents reason about a connected cloud account, and optional infrastructure as code (IaC), then chain misconfigurations into attack paths.

[Cloud Scanning](/cloud-scanning/cloud-scanning-overview.md) flags individual resources that fail a known check. Cloud Security Audit asks which of those settings, taken together, become an attack path.

## Use cases

* **Chain misconfigurations:** See which settings combine into an attack path
* **Correlate with IaC:** Add the repositories that define the account so a finding can point at the template
* **Review before a change:** Run an on-demand audit of clouds you already connected

## Cloud Scanning vs Cloud Security Audit

Cloud Scanning checks connected accounts continuously against known misconfiguration rules. Cloud Security Audit runs on demand. It reads the same inventory, then reasons across resources and optional IaC.

|                   | Cloud Security Audit                                 | Cloud Scanning                      |
| ----------------- | ---------------------------------------------------- | ----------------------------------- |
| **Starts from**   | Cloud inventory, plus optional IaC                   | Connected cloud accounts            |
| **Core question** | Which misconfigurations combine into an attack path? | Which resources fail a known check? |
| **When it runs**  | On demand, paid in credits                           | Continuously, on each cloud scan    |

Cloud Scanning flags a broken check. Cloud Security Audit shows how checks combine.

## How it works

You pick one or more connected clouds. IaC repositories are optional. Agents read the inventory Aikido already collected. They do not change the account. When you include IaC, a finding can point at the file that defines a live setting.

Findings land in the **Issues** tab and in the Feed. To run one, see [Run a Cloud Security Audit](/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit.md). For how to read findings, including retesting, see [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md). For what is in scope, see [What Cloud Security Audit Finds](/ai-code-analysis/cloud-security-audit/what-cloud-security-audit-finds.md).

## Related

{% content-ref url="/pages/sU44cPzbmeqn9z1n9Y1a" %}
[Run a Cloud Security Audit](/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/tweiYLnrxouYwMmBjEiE" %}
[Cloud Scanning Overview](/cloud-scanning/cloud-scanning-overview.md)
{% endcontent-ref %}
