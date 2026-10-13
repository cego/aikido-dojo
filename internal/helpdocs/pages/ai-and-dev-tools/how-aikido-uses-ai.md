# How Aikido Uses AI

Aikido uses AI throughout the full software development lifecycle. It starts in the IDE, continues in pull requests and scans, helps teams fix and enforce security standards, and extends into runtime protection and pentesting.

For which product to turn on (SAST, Code Quality, Deep Review, AI Code Analysis, DAST, and Pentest), see [Which Aikido product should I use?](/getting-started/which-aikido-product-should-i-use.md).

At the center of that intelligence is [Aikido Agent](/aikido-agent/aikido-agent.md). It understands your security data and powers [Ask Aikido](/aikido-agent/ask-aikido.md), [Custom Reports](/aikido-agent/custom-reports.md), [SAST AutoTriage](/aikido-agent/sast-autotriage.md), [Secrets AutoTriage](/aikido-agent/secrets-autotriage.md), and [CVE Exploitability Analysis](/aikido-agent/cve-exploitability-analysis.md).

{% hint style="info" %}
We never use, store, or train on any customer data. Small, anonymized code fragments may be used to guarantee the accuracy of triaging and fixing vulnerabilities.

All AI operations are inference-only, meaning data are processed transiently in memory and never retained or reused.

For a detailed overview of our controls, environments, and compliance measures, please refer to our AI Policy, available upon request through our [Trust Center](https://www.aikido.dev/trust-center).
{% endhint %}

### AI across the SDLC

From first code to production, Aikido uses AI to:

* prioritize issues early
* interrogate findings with Ask Aikido
* analyze dependency CVEs against real code and container usage
* review pull requests with Deep Review
* reason about existing source with AI Code Analysis
* generate and refine fixes
* generate API specs
* enforce custom code and cloud rules
* learn repo-specific standards through extra code context
* track AI usage in production
* validate security with AI-powered pentests

### 1. Start in the IDE

Aikido brings AI directly into the IDE. Developers can prioritize findings and apply fixes before code reaches a pull request.

Read [Aikido AI in IDE](/ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide.md).

### 2. Scans and review

As code moves through scans and pull requests, [Aikido Agent](/aikido-agent/aikido-agent.md) cuts noise and surfaces what matters first. It checks exploitability, reads real code context, and reprioritizes issues based on likely impact. Check out [SAST AutoTriage](/aikido-agent/sast-autotriage.md).

Aikido applies the same approach to secret findings. [Secrets AutoTriage](/aikido-agent/secrets-autotriage.md) filters out false positives and prioritizes exposed secrets by whether they are still active and how much access they grant.

For dependency findings, Aikido also runs [CVE Exploitability Analysis](/aikido-agent/cve-exploitability-analysis.md). The Agent inspects how a vulnerable package is used in your repository or container and decides whether the CVE is actually exploitable in your environment. Based on your settings, Aikido can downgrade, upgrade, snooze, or ignore the finding automatically.

When you want to inspect one finding in more detail, use [Ask Aikido](/aikido-agent/ask-aikido.md). It helps you validate severity, understand realistic attack paths, and ask whether an issue is actually reachable or exploitable in your app. For dependency CVEs, it can also explain the reachability path and the full impact of the vulnerable package in plain language.

To turn security data into reusable workspace reports, use [Custom Reports](/aikido-agent/custom-reports.md). Describe the insights you need in natural language, refine the charts and tables in chat, and save the result for your team.

[Code Quality](/code-quality/code-quality-overview.md) is also AI-powered. It reviews newly introduced pull request changes and helps enforce engineering standards across many languages.

[Deep Review](/deep-review/how-deep-review-works.md) is the agentic pull request reviewer. It reads the change with context beyond the diff, and it covers both **security and quality**.

[AI Code Analysis](/ai-code-analysis/ai-code-analysis.md) reasons across existing source when you want depth without a live target: logic flaws, authorization gaps, and attack paths that pattern checks will miss.

You can also add extra code context. This gives Aikido AI more signal and less noise. Use it to explain accepted exceptions, architectural choices, and repo-specific standards. Aikido then uses that context to make Code Quality comments more relevant. Check out how to [Add Extra Code Context](/code-quality/add-extra-code-context.md).

### 3. Fix issues with AutoFix

Aikido generates reviewable fixes for code, dependencies, infrastructure, and containers. You can apply fixes in the IDE or open pull requests for review.

You can also refine generated fixes with follow-up instructions, so the patch better matches your codebase and standards.

Read [AutoFix Overview](/autofix-and-remediation/overview-aikido-autofix.md), [AutoFix for SAST and IaC Issues](/autofix-and-remediation/scope/ai-autofix-for-sast-and-iac-issues.md), [AutoFix for Open Source Dependencies](/autofix-and-remediation/scope/autofix-for-open-source-dependencies.md), [AutoFix for Containers](/autofix-and-remediation/scope/ai-autofix-for-containers.md), and [Refine AutoFixes with Aikido AI](/autofix-and-remediation/automation-and-merging/refine-autofixes-with-aikido-ai.md).

### 4. Generate specs and enforce policies

Aikido can generate an OpenAPI specification directly from backend code. That helps teams start API scanning without maintaining the spec by hand.

Read [Autogenerate OpenAPI via Aikido AI (Code2Swagger)](/dast-surface-monitoring/api-scanning/autogenerate-openapi-via-aikido-ai-code2swagger.md).

Aikido also uses AI for custom rule creation. Teams can define custom code checks in natural language and generate custom cloud misconfiguration checks for their environment.

Read [Add Custom Code Rules](/code-quality/add-custom-code-rules.md) and [Custom CSPM Rules](/cloud-scanning/custom-cspm-rules.md).

### 5. Track runtime AI usage

Zen Firewall tracks LLM provider usage, model activity, token consumption, and estimated cost. That gives teams visibility into how AI is used in running applications.

Read [Tracking AI / LLM usage with Zen Firewall](/zen-firewall/zen-features/tracking-ai-llm-usage-with-zen-firewall.md).

### 6. Validate security with Pentests

Aikido Pentest uses autonomous agents to discover, exploit, and validate vulnerabilities across applications, APIs, and infrastructure.

Read [Pentest Overview](/pentests/aikido-pentest.md).
