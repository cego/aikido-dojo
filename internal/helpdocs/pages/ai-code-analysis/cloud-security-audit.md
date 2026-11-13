# Cloud Security Audit

Cloud Security Audit is the cloud-focused audit type in [AI Code Analysis](https://help.aikido.dev/ai-code-analysis/ai-code-analysis). It is Aikido's agentic approach to cloud security posture management (CSPM). Agents reason about a connected cloud account, and optional infrastructure as code (IaC), then chain misconfigurations into attack paths.

[Cloud Scanning](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview) flags individual resources that fail a known check. Cloud Security Audit asks which of those settings, taken together, become an attack path.

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

Findings land in the **Issues** tab and in the Feed. To run one, see [Run a Cloud Security Audit](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit). For how to read findings, including retesting, see [Track Progress and Findings](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings). For what is in scope, see [What Cloud Security Audit Finds](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/what-cloud-security-audit-finds).

## Related

{% content-ref url="run-a-cloud-security-audit.md" %}
[run-a-cloud-security-audit.md](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/run-a-cloud-security-audit)
{% endcontent-ref %}

{% content-ref url="../../cloud-scanning/cloud-scanning-overview.md" %}
[cloud-scanning-overview.md](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview)
{% endcontent-ref %}
