# AI Scans

AI Scans uses agents to find security, data-exposure, cloud, and smart-contract risks that need context across files, services, repositories, or a connected cloud account.

Choose the scan type that matches your goal:

<table data-view="cards"><thead><tr><th>Scan type</th><th>What it helps you find</th><th data-hidden data-card-target data-type="content-ref">Target</th></tr></thead><tbody><tr><td><strong>AI Code scan</strong></td><td>AI code scans are complementary to normal code scanning and find novel attack paths and security risks.</td><td><a href="https://help.aikido.dev/ai-scans/ai-code-audit-overview">ai-code-audit-overview</a></td></tr><tr><td><strong>AI Data Exposure scan</strong></td><td>DSPM: Find sensitive data flows and data exposure risks in code with AI.</td><td><a href="https://help.aikido.dev/ai-scans/what-is-dspm">what-is-dspm</a></td></tr><tr><td><strong>AI Cloud scan</strong></td><td>CSPM: Audit connected clouds for misconfigurations and correlate them with optional IaC.</td><td><a href="https://help.aikido.dev/ai-scans/cloud-security-audit">cloud-security-audit</a></td></tr><tr><td><strong>AI Smart Contract scan</strong></td><td>Blockchain: Find exploitable vulnerabilities in on-chain smart contracts with AI. Currently supports Solidity.</td><td><a href="https://help.aikido.dev/ai-scans/smart-contract-audit">smart-contract-audit</a></td></tr></tbody></table>

## How it works

AI Code, AI Data Exposure, and AI Smart Contract scans read connected repositories. They don't need a staging environment or a live target. An AI Data Exposure scan doesn't access your data stores.

An AI Cloud scan reads the inventory of clouds you already connected. You can add infrastructure as code (IaC) repositories so a finding can point at the template that defines a live setting.

Findings include evidence and remediation guidance.

* **Agentic analysis:** Aikido uses multiple agents to investigate the selected code or cloud inventory, follow relevant context, and validate potential findings.
* **Repository instructions:** Commit `AGENTS.md` at the repository root for architecture, trust boundaries, and accepted risks. Aikido also reads `CLAUDE.md`, `.claude/`, and `.cursor/`. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files).
* **Models and harnesses matched to the task:** Aikido selects the models and analysis harness that best fit each stage of the analysis, balancing the depth of analysis with the resources it uses.
* **Confirmed cost before you start:** The pricing step shows the credit total before you launch. **Configure** can change that total before you start. Once you click **Start Audit**, the total is fixed for the analysis.

For testing a running application, use [Aikido Pentest](https://help.aikido.dev/pentests/aikido-pentest). For fast, continuous checks of known vulnerability patterns, use [SAST](https://help.aikido.dev/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus). For continuous rule-based cloud checks, use [Cloud Scanning](https://help.aikido.dev/cloud-scanning/cloud-scanning-overview).
