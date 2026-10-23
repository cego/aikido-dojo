# Code Security Audit

Pentest-grade reasoning on your source code, no environment required.

Code Security Audit is the security-focused audit type in [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md). It uses pentest-grade reasoning to identify vulnerabilities, logic flaws, and exploit chains in your source code.

You don't need a staging URL, a crawl, or pentest scope setup. Connect a repository, confirm the price in credits, and start the audit.

<table data-view="cards"><thead><tr><th>Section</th><th>What you will find</th><th data-hidden data-card-target data-type="content-ref">Target</th></tr></thead><tbody><tr><td><strong>Create a Code Security Audit</strong></td><td>Select the audit type, choose repositories, and confirm the credit cost.</td><td><a href="/ai-code-analysis/ai-code-audit-overview/create-an-ai-code-audit.md">Create a Code Security Audit</a></td></tr><tr><td><strong>Track progress and findings</strong></td><td>Status vocabulary, detail tabs, and how to read findings.</td><td><a href="/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md">Track Progress and Findings</a></td></tr><tr><td><strong>What Code Security Audits Find</strong></td><td>Vulnerability classes covered, plus what still needs a live target.</td><td><a href="/ai-code-analysis/ai-code-audit-overview/what-ai-code-audit-finds.md">What Code Security Audits Find</a></td></tr></tbody></table>

## How it works

### Connect repositories

You pick the repositories in the create flow. Code Security Audit only needs source code, with no live target, domains, or auth setup. See [Aikido Never Stores Your Code](/getting-started/setting-up-your-account/aikido-never-stores-your-code.md) for how your source is handled during an analysis.

### Agents reason about your code

Multiple security agents work through the connected codebases together, chaining context across files to surface architectural and logic flaws that traditional static scanners can't see. Commit `AGENTS.md` in the repository when you want that reasoning to follow your architecture and accepted risks. See [Guide Agents with Repository Files](/ai-and-dev-tools/guide-an-audit-with-repository-files.md).

### Read the findings

Findings land in the **Issues** tab on the assessment and in the Feed, alongside everything else Aikido detects. Each finding includes a summary, root cause, remediation, and code-based **Evidence**. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md) for the full read-out.

## Pricing and credits

Code Security Audit is paid for with [Aikido credits](/miscellaneous-info/wallet-and-credits.md).

* The **Pricing** step shows the exact credit total before you commit.
* Cost depends on the size and complexity of the codebases you select. Larger or more sprawling selections cost more.
* Credits are charged when you click **Start Audit**.

For how these fit with SAST, Deep Review, dynamic application security testing (DAST), and Attack Surface Monitoring, see [Which Aikido product should I use?](/getting-started/which-aikido-product-should-i-use.md).

## Code Security Audit vs Aikido Pentest

Both run the same agentic engine, but they answer different questions. Code Security Audit reasons about your source code. Aikido Pentest exercises your running application.

### Use Code Security Audit when:

* You want deep code reasoning on logic and architectural flaws, including broken access control, privilege escalation, and multi-step exploit chains, without configuring a live environment.
* You don't have a stable staging or QA target, or auth flows aren't ready for live testing.
* Your codebase is hard to set up for live testing, such as a mobile app, desktop app, embedded system, or library. For a review focused on on-chain contracts, use [Smart Contract Audit](/ai-code-analysis/smart-contract-audit.md).
* You need a fast turnaround with minimal setup: connect a repository, confirm credits, start.
* You want to validate changes in source before they ship to a live deployment.

### Use Aikido Pentest when:

* You have a live target and want to validate real exploitability with real traffic.
* You want runtime evidence, including reproduction requests, attack-surface mapping, and live agent activity.
* Your scope includes domains, authenticated user roles, and crawl-discovered endpoints beyond what's visible in source.
* You need **Exploit Further** on validated findings to escalate exploitation against the running app and prove real impact, such as whether sensitive data can actually be reached.
* You're satisfying compliance frameworks like SOC 2 or ISO 27001 that expect a live pentest.

## Code Security Audit vs SAST

Both look at your source code, but they work differently. [Static application security testing (SAST)](/code-scanning/scanning-practices/sast-by-aikido-supported-languages-and-security-focus.md) flags patterns, such as a tainted parameter, a risky API call, or a missing check. Code Security Audit reasons about intent and context across files, so it can identify issues that need an attacker's perspective. The two are complementary.

### Use Code Security Audit when:

* You want to catch logic and architectural flaws that pattern-matching can't see, including broken access control, privilege escalation, business logic bypasses, and multi-step exploit chains.
* You need cross-file reasoning that follows references through services, modules, and helpers.
* You're analyzing a high-stakes change, release, or codebase and want attacker-style review on top of your continuous scans.
* You want deeper context on a specific finding, including the root cause, exploit path, and remediation, instead of just a pattern hit.

### Use SAST when:

* You want fast, every-commit feedback on common vulnerability patterns, including injection sinks, taint flows, hard-coded secrets, dangerous APIs, and infrastructure as code (IaC) misconfigurations.
* You need broad, continuous coverage across every pull request and branch without spending credits per scan.
* You're enforcing pull-request gates in CI/CD on known-bad patterns.

## FAQ

### Why doesn't Code Security Audit need a live URL?

Code Security Audit reads and reasons about your source code directly. There's no crawl phase, traffic replay, or live exploitation, so there's no environment to point at. For live testing against a deployed target, use [Aikido Pentest](/pentests/aikido-pentest.md).

### Can I retest a Code Security Audit finding?

Yes. Once you've pushed a fix, use **Re-test** on the finding to have Aikido re-check it against your latest code. **Exploit Further** is an Aikido Pentest action. It needs a running app. Retesting for the other audit types is covered in [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md).

## Related

{% content-ref url="/pages/X4yc23Inhvx1gPKx2AwE" %}
[Create a Code Security Audit](/ai-code-analysis/ai-code-audit-overview/create-an-ai-code-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/X3BH6OfadNRDh0qSJiwm" %}
[Guide Agents with Repository Files](/ai-and-dev-tools/guide-an-audit-with-repository-files.md)
{% endcontent-ref %}

{% content-ref url="/pages/64GGpVtiCGGThzz1y8B7" %}
[AI Code Analysis](/ai-code-analysis/ai-code-analysis.md)
{% endcontent-ref %}
