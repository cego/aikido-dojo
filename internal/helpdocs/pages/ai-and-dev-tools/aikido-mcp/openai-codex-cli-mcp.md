# OpenAI Codex CLI MCP

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

For the full list of tools the MCP server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools).

## Installation

Add the Aikido MCP and rules via our [Aikido Plugin](https://chatgpt.com/plugins/plugins_6abf6f017e3481919e681ac223b8424e).\
\
Or install manually using the steps below:

{% stepper %}
{% step %}
**Install the Aikido MCP server**

```shellscript
codex mcp add aikido -- npx -y @aikidosec/mcp@latest
```

{% endstep %}

{% step %}
**Authenticate**

On first use, the MCP server will detect no authentication is set. Ask Codex to run `aikido_login`. It will return a browser link. Follow the link to log in to Aikido. Your token is stored securely once login is complete.
{% endstep %}

{% step %}
**Add the Aikido rule to Global AGENTS file**

Create the codex directory if it doesn't exist yet.

```shellscript
mkdir -p ~/.codex/skills
```

Download the Aikido rule and add it to `~/.codex/skills/aikido-rule.txt`.

```shellscript
curl -fsSL "https://gist.githubusercontent.com/kidk/aa48cad6db80ba4a38493016aae67712/raw/3644397b7df43423e3da06434491b40bbb79dd47/aikido-rule.txt" \
  -o ~/.codex/skills/aikido-rule.txt
```

{% endstep %}

{% step %}
**Finished**

Aikido MCP is now available in Codex CLI.

{% hint style="info" %}
Restart Codex CLI if it was open.
{% endhint %}
{% endstep %}
{% endstepper %}

{% hint style="info" %}
For CI/CD or headless setups where browser login is not possible, see [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md).
{% endhint %}

For ways to ask your agent to use Aikido, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).
