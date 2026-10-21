# Jetbrains AI / Junie MCP

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

## Installation

{% stepper %}
{% step %}
**Add Aikido MCP server to your Jetbrains IDE config**

In your IDE go to Settings | Tools | AI Assistant | Model Context Protocol (MCP).

{% hint style="info" %}
[Check out Jetbrains documentation for more information](https://www.jetbrains.com/help/ai-assistant/mcp.html)
{% endhint %}

On the Model Context Protocol (MCP) settings page, click  Add to add a new MCP server configuration.

In the New MCP Server dialog, select how you want to connect to the MCP server and provide a JSON configuration:

```shellscript
{
  "mcpServers": {
    "aikido": {
      "command": "npx -y @aikidosec/mcp@latest",
      "args": [],
      "environment": {}
    }
  }
}
```

{% endstep %}

{% step %}
**Authenticate**

On first use, ask the AI assistant to run `aikido_login`. It will return a browser link, follow it to log in to Aikido. Your token is stored securely once login is complete.
{% endstep %}

{% step %}
Adapt your agent instructions to use the Aikido MCP by adding the following skill or instructions.

```
---
name: aikido-security
description: Enforce mandatory security scanning for all code changes
license: MIT
compatibility: opencode
metadata:
  audience: developers
  security: aikido
---

## Purpose

Ensure that all first-party code is scanned for security issues whenever code is written or changed.

## Rules

- Always run **aikido_scan_paths** on all newly written, generated, added, or modified first-party code.
- Pass the file paths to the scanner. Only use **aikido_full_scan** when the code is not saved to disk yet.
- If any security issues are reported, apply the fixes using the remediation guidance from Aikido.
- Run **aikido_scan_paths** again after fixes to confirm that no issues remain and no new ones were introduced.
- Repeat the fix and rescan cycle until the scan passes with zero issues.
- If the Aikido MCP server is not installed or not running, inform the user to install it using the official setup guide at  
  https://help.aikido.dev/ide-plugins/aikido-mcp
```

{% endstep %}

{% step %}
**Finished**

Aikido MCP is now available in Jetbrains IDE
{% endstep %}
{% endstepper %}

{% hint style="info" %}
For CI/CD or headless setups where browser login is not possible, see [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md).
{% endhint %}

For the full list of tools the MCP server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools). For ways to ask your agent to use Aikido, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).
