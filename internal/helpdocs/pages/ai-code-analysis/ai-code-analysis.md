# AI Code Analysis

Choose a Code Security, Data Exposure, Cloud Security, or Smart Contract audit.

AI Code Analysis uses agents to find security, data-exposure, cloud, and smart-contract risks that need context across files, services, repositories, or a connected cloud account.

Choose the audit type that matches your goal:

<table data-view="cards"><thead><tr><th>Audit type</th><th>What it helps you find</th><th data-hidden data-card-target data-type="content-ref">Target</th></tr></thead><tbody><tr><td><strong>Code Security Audit</strong></td><td>AI SAST: Find real attack paths and validated security risks in code with AI.</td><td><a href="/ai-code-analysis/ai-code-audit-overview.md">Code Security Audit</a></td></tr><tr><td><strong>Data Exposure Audit</strong></td><td>DSPM: Find sensitive data flows and data exposure risks in code with AI.</td><td><a href="/ai-code-analysis/what-is-dspm.md">Data Exposure Audit (DSPM)</a></td></tr><tr><td><strong>Cloud Security Audit</strong></td><td>CSPM: Audit connected clouds for misconfigurations and correlate them with optional IaC.</td><td><a href="/ai-code-analysis/cloud-security-audit.md">Cloud Security Audit</a></td></tr><tr><td><strong>Smart Contract Audit</strong></td><td>Blockchain: Find exploitable vulnerabilities in on-chain smart contracts with AI.</td><td><a href="/ai-code-analysis/smart-contract-audit.md">Smart Contract Audit</a></td></tr></tbody></table>

## How it works

Code Security, Data Exposure, and Smart Contract audits read connected repositories. They don't need a staging environment or a live target. Data Exposure Audit doesn't access your data stores.

Cloud Security Audit reads the inventory of clouds you already connected. You can add infrastructure as code (IaC) repositories so a finding can point at the template that defines a live setting.

Findings include evidence and remediation guidance.

* **Agentic analysis:** Aikido uses multiple agents to investigate the selected code or cloud inventory, follow relevant context, and validate potential findings.
* **Models and harnesses matched to the task:** Aikido selects the models and analysis harness that best fit each stage of the analysis, balancing the depth of analysis with the resources it uses.
* **Confirmed cost before you start:** The pricing step shows the credit total before you launch. **Configure** can change that total before you start. Once you click **Start Audit**, the total is fixed for the analysis.

For testing a running application, use [Aikido Pentest](/pentests/aikido-pentest.md). For fast, continuous checks of known vulnerability patterns, use [SAST](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md). For continuous rule-based cloud checks, use [Cloud Scanning](/cloud-scanning/cloud-scanning-overview.md).
