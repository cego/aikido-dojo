# Securing AI-Generated Code

**Enhance Your AI Development Workflow with Aikido**

AI coding agents like Claude can speed up development by writing and committing code autonomously. They also introduce new risks. They can leak secrets, install fake packages, miss IDE feedback, or ship code without enough security checks.

**How Aikido Helps**

* **Pre-commit controls:** Catch mistakes before code lands in the repo.
* **During code generation:** Keep security checks close to the agent workflow.
* **Pre-merge checks:** Stop risky changes before they hit your main branch.

Whether your team uses Claude or another AI tool, Aikido adds guardrails across the full workflow.

### Main Risks

* Hallucinated package names that point to malware
* Secrets or tokens committed by mistake
* Security findings missed because the agent runs outside the IDE
* No security checks between code generation and merge

### Recommended Controls

#### Safe Chain: Block Malicious and Hallucinated Packages

[Safe Chain](https://help.aikido.dev/code-scanning/aikido-malware-scanning) validates packages against trusted registries before installation. It blocks fake or hallucinated package names before they reach your environment. It also scans nested dependencies for malicious behavior such as obfuscated code, data exfiltration, install scripts, and crypto miners. New package versions published less than 24 hours ago are blocked by default, and Safe Chain falls back to the latest older safe version so builds keep moving.

#### Secrets Pre-Commit Hook: Prevent Secrets from Being Committed

The [Aikido Secrets Pre-Commit Hook](https://help.aikido.dev/ai-and-dev-tools/aikido-secrets-pre-commit-hook) scans staged changes and blocks commits that contain secrets, passwords, or API keys. It works for both developer-written and agent-generated changes. Because it runs locally before code leaves the machine, it reduces the chance of accidental exposure in Git history, CI logs, or pull requests. Teams using the IDE plugins can also install it directly through [Aikido Expansion Packs](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs) (currently available for VS Code and related IDEs).

#### Aikido MCP: Scan Code During Generation

The [AI Coding Assistants (MCP)](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp) setup connects Aikido directly to your agent environment. It runs secret scanning and SAST while code is being written, and can help fix issues before commit. When used with [Automatically handle MCP rules in IDE](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/automatically-handle-mcp-rules-in-ide), Aikido can automatically enforce a scan-fix-rescan loop for generated, added, and modified code. That keeps security checks inside the AI workflow instead of waiting for a later scan.

#### Aikido AI in IDE and Expansion Packs: Keep Findings Visible

[Aikido AI in IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide) keeps findings and fixes visible while code is being written. It can use AI AutoTriage to prioritize the issues that matter most and AI AutoFix to generate reviewable fixes with a diff preview and automatic rescan. If you want faster rollout, [Aikido Expansion Packs](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs) can install related tooling like the pre-commit hook, Safe Chain, and MCP directly from the IDE. This keeps the full developer and agent workflow in one place.

#### PR Gating: Enforce Checks Before Merge

[PR Gating](https://help.aikido.dev/pr-and-release-gating/aikido-ci-gating-functionality) is the final control point for agent-generated changes before they reach your main branch. It scans the branch diff for new issues across SCA, IaC, Secrets, SAST, malware, license risks, and code quality. It also shows which existing issues the branch solves. You can set severity thresholds, choose which scans run, and decide whether draft pull requests should be checked. That makes it a strong backstop even if earlier controls were skipped or bypassed.

### Suggested Setup

1. Enable [Safe Chain](https://help.aikido.dev/code-scanning/aikido-malware-scanning) for package installs
2. Install the [Aikido Secrets Pre-Commit Hook](https://help.aikido.dev/ai-and-dev-tools/aikido-secrets-pre-commit-hook)
3. Connect [AI Coding Assistants (MCP)](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp) to your agent environment
4. Keep [PR Gating](https://help.aikido.dev/pr-and-release-gating/aikido-ci-gating-functionality) enabled as a final check
