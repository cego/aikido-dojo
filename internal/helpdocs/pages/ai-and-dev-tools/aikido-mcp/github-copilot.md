# Github Copilot

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

For the full list of tools the MCP server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools).

## Installation

### IDE

The Aikido VSCode IDE plugin uses [Expansion Packs](/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs.md) to provide additional features. The Aikido MCP for Github Copilot is one of these Expansion Packs, making installation simple and fast without a separate setup process. See the linked page below for instructions on how to enable it.

{% content-ref url="/pages/mgd9l65M1Oa5zH7GUy7e" %}
[VS Code IDE](/ai-and-dev-tools/ide-plugins-overview/vs-code-plugin.md)
{% endcontent-ref %}

### Integration with Copilot IDE plugins

{% stepper %}
{% step %}
**Add the MCP configuration in your IDE settings**

Go to Github Copilot plugin settings and open the `mcp.json` file. Add the Aikido server:

```json
{
  "servers": {
    "aikido": {
      "command": "npx",
      "args": ["-y", "@aikidosec/mcp@latest"]
    }
  }
}
```

{% endstep %}

{% step %}
**Authenticate**

On first use, ask Copilot to run `aikido_login`. It will return a browser link. Follow it to log in to Aikido. Your token is stored securely once login is complete.
{% endstep %}

{% step %}
**Add custom prompt instructions**

Add an instruction to your `AGENTS.md` file or custom instructions to tell agent to use the new tooling.

```
Scan any new or modified code with `aikido_scan_paths` before finalizing changes. This security check is required.
```

{% endstep %}
{% endstepper %}

{% hint style="info" %}
For CI/CD or headless setups where browser login is not possible, see [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md).
{% endhint %}

### Integration with Copilot Cloud Agents

GitHub Copilot coding agent can call tools from MCP servers while it runs. Follow the steps below to set up Aikido MCP for your repository. Configuration is done per repository.

{% hint style="warning" %}
Copilot Cloud Agents run in a headless environment and require a token because browser-based login is not available. You will need an Aikido Personal Access Token. See [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md) for how to create one.
{% endhint %}

If you already use Aikido MCP in VS Code, [you can reuse that configuration](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/extend-coding-agent-with-mcp#reusing-your-mcp-configuration-from-visual-studio-code) and only add the token in your repository settings.

#### Set-up Instructions

{% hint style="info" %}
For the full GitHub flow (where to paste JSON, how validation works, and how Copilot environments expose secrets), [check out GitHub’s guide](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/extend-coding-agent-with-mcp).
{% endhint %}

{% stepper %}
{% step %}
**Add the MCP configuration in your repo settings**

In your GitHub repository, go to Settings → Copilot → Coding agent and paste an MCP JSON config.

Example configuration for Aikido MCP (local/stdio via npx):

{% tabs %}
{% tab title="Copilot Cloud Agent" %}

```json
{
  "mcpServers": {
    "aikido": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@aikidosec/mcp@latest"],
      "tools": ["aikido_scan_paths", "aikido_full_scan"],
      "env": {
        "AIKIDO_API_KEY": "COPILOT_MCP_AIKIDO_API_KEY"
      }
    }
  }
}
```

{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}
**Add the token to your Copilot environment**

Github Copilot only exposes secrets/variables to MCP config if their names start with `COPILOT_MCP_`

Create a secret like:

`COPILOT_MCP_AIKIDO_API_KEY` = your Aikido token
{% endstep %}

{% step %}
**Add custom prompt instructions**

To ensure GitHub Copilot uses the MCP server, configure repository instructions as [described in the GitHub guide](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions).

Add an instruction such as:

```
Scan any new or modified code with `aikido_scan_paths` before finalizing changes. This security check is required.
```

{% endstep %}

{% step %}
**Validate**

You can now [verify your configuration](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/extend-coding-agent-with-mcp#validating-your-mcp-configuration) by asking Copilot to security scan one of your files.
{% endstep %}
{% endstepper %}

For ways to ask your agent to use Aikido, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).
