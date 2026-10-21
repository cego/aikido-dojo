# Aikido Agent

Aikido Agent understands your security data. It can answer questions, investigate attack paths, triage findings, and adapt fixes.

**Aikido Agent** is Aikido's security intelligence. It understands your findings, code context, and triage signals, so it can answer questions, investigate attack paths, adapt fixes, and power automated triage.

{% hint style="info" %}
We never use, store, or train on any customer data. Small, anonymized code fragments may be used to guarantee the accuracy of triaging and fixing vulnerabilities. All AI operations are inference-only. Data is processed transiently in memory and never retained or reused. For a detailed overview of our controls, see our AI Policy via the [Trust Center](https://www.aikido.dev/trust-center).
{% endhint %}

## What you can do with Aikido Agent

<table data-view="cards"><thead><tr><th>Capability</th><th>What it does</th><th data-hidden data-card-target data-type="content-ref">Target</th></tr></thead><tbody><tr><td><strong>Ask Aikido</strong></td><td>Chat with the Agent on a finding to validate risk, inspect attack paths, and adapt fixes.</td><td><a href="/aikido-agent/ask-aikido.md">Ask Aikido</a></td></tr><tr><td><strong>SAST AutoTriage</strong></td><td>Filter and prioritize static analysis findings based on exploitability and code context.</td><td><a href="/aikido-agent/sast-autotriage.md">SAST AutoTriage</a></td></tr><tr><td><strong>Secrets AutoTriage</strong></td><td>Filter false-positive secrets and rank real exposures by risk.</td><td><a href="/aikido-agent/secrets-autotriage.md">Secrets AutoTriage</a></td></tr><tr><td><strong>CVE Exploitability Analysis</strong></td><td>Assess dependency CVEs in your repositories and containers and triage what's truly exploitable.</td><td><a href="/aikido-agent/cve-exploitability-analysis.md">CVE Exploitability Analysis</a></td></tr><tr><td><strong>Custom Reports</strong></td><td>Ask natural-language questions about your security data and save the results as workspace reports.</td><td><a href="/aikido-agent/custom-reports.md">Custom Reports Creation</a></td></tr></tbody></table>

## Related

Related products that use AI but live elsewhere in the docs: [AutoFix](/autofix-and-remediation/overview-aikido-autofix.md), [Code Analysis](/ai-code-analysis/ai-code-audit-overview.md), and [Aikido Pentest](/pentests/aikido-pentest.md). For the full map of AI across the SDLC, see [How Aikido Uses AI](/ai-and-dev-tools/how-aikido-uses-ai.md).
