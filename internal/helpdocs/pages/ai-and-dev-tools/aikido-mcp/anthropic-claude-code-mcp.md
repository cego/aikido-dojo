# Anthropic Claude Code Plugin

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

## Installation for Claude Code CLI

{% stepper %}
{% step %}
**Install the Aikido Plugin**

Installing the Aikido Claude Code plugins can be done via Claude Code CLI.

* Open up the terminal and enter `claude`.
* Now, you can install the plugin from the Claude Plugins Marketplace
  * `/plugin install aikido@claude-plugins-official`
* After install, run `/reload-plugins` to activate the plugin.
* A restart of Claude Code CLI is needed to load the skills. You can exit Claude Code CLI by typing `exit`.
  {% endstep %}

{% step %}
**Authenticate**

After you've started Claude Code again, run the setup skill with no arguments:

```
/aikido:setup
```

This opens a browser window where you can log in to Aikido. Once authenticated, the token is stored securely in your Claude Code settings and the MCP server is registered automatically.
{% endstep %}

{% step %}
**Finished**

The Aikido Claude Plugin is now available in Claude Code.
{% endstep %}
{% endstepper %}

{% hint style="info" %}
For CI/CD or headless setups where browser login is not possible, see [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md).
{% endhint %}

The Aikido plugin comes with its own skills. To decide exactly which skills your developers get, and roll them out with nothing for them to configure, see [Bundle Aikido in your own Claude Code plugin](/ai-and-dev-tools/aikido-mcp/anthropic-claude-code-mcp/claude-code-internal-plugin.md).

For the full list of tools the MCP server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools). For ways to ask your agent to use Aikido, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).

### Troubleshooting

If the setup skill fails to start the MCP server, check the [troubleshooting guide](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/mcp-troubleshooting).
