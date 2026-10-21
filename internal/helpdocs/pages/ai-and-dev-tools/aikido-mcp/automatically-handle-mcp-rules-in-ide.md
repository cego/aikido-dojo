# Automatically handle MCP rules in IDE

When this option is enabled, Aikido automatically manages the MCP rule required for LLM and AI usage.

This works together with [Aikido MCP](https://help.aikido.dev/ide-plugins/aikido-mcp) to ensure AI generated code is scanned by Aikido when using MCP compatible tools.

#### What this does

Aikido automatically adds and maintains one MCP rule that enforces security scanning for AI generated, added, and modified code. No manual MCP rule setup is required.

#### The MCP rule Aikido manages

When enabled, Aikido manages the following rule:

```
---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "**"
description: Aikido MCP Security Rules
---

# Security best practices

- Always run aikido_scan_paths on generated, added and modified first-party code, unless the prompt explicitly says not to.
- Pass the file paths to the scanner. Only use aikido_full_scan when the code is not saved to disk yet.
- If the scan identifies any security issues in newly generated or modified code, fix the issues using the scan result and remediation provided by Aikido.
- After applying fixes, run aikido_scan_paths again to verify that the issues were successfully resolved, and no new issues were introduced.
- Repeat the fix-and-rescan cycle until the code passes with zero remaining or newly introduced security issues.
```

#### Gitignore

To apply this rule, Aikido may create or update an MCP rule file in your repository. This file is not always added to .gitignore.

Add it to .gitignore if you want Aikido to fully manage the rule without committing it. Keep it out of .gitignore if you want to review or version the rule in Git.

#### When to disable this option

Disable this option if you want to manage MCP rules manually or do not want Aikido to modify rule files automatically.
