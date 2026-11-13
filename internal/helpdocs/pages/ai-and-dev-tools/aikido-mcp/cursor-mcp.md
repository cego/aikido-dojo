# Cursor MCP

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

For the full list of tools the MCP server exposes, see [Available tools](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/using-aikido-mcp#available-tools).

## Installation

### Via Cursor Marketplace  (recommended)

You can also install Aikido MCP directly from the [Cursor Marketplace](https://cursor.com/marketplace/aikido).

{% stepper %}
{% step %}
**Add via Cursor Marketplace**

Open [Aikido on the Cursor Marketplace](https://cursor.com/marketplace/aikido) and click **Add to Cursor**.
{% endstep %}

{% step %}
**Authenticate**

On first use, the MCP server will detect no authentication is set and prompt you to log in. Ask the agent to run `aikido_login`. It will return a browser link. Follow the link to authenticate with Aikido. Your token is stored securely once login is complete.
{% endstep %}

{% step %}
**Add custom instructions**

Add an instruction to your `.cursorrules` file or Cursor custom instructions to ensure the agent uses the MCP server.

```
Scan any new or modified code with `aikido_scan_paths` before finalizing changes. This security check is required.
```
{% endstep %}
{% endstepper %}

{% hint style="info" %}
For CI/CD or headless setups where browser login is not possible, see [Token Authentication](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/token-authentication).
{% endhint %}

### Via Aikido IDE Expansion Packs

The easiest way to add the Aikido MCP to Cursor is through the Aikido IDE plugin's Expansion Packs. This handles configuration automatically without any manual setup.

{% content-ref url="../ide-plugins-overview/cursor-ide.md" %}
[cursor-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/cursor-ide)
{% endcontent-ref %}

Once the plugin is installed, open the Aikido sidebar and go to **Aikido Cursor Plugin** to enable the Aikido MCP server.

[Learn more about Expansion Packs](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs)

For ways to ask your agent to use Aikido, see [Example prompts](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/using-aikido-mcp#example-prompts).

### Manual installation

For any custom setup, refer to the [npm package page](https://www.npmjs.com/package/@aikidosec/mcp) for detailed manual installation instructions.
