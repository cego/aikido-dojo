# Devin MCP

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

For the full list of tools the MCP server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools).

## Installation

{% stepper %}
{% step %}
**Install the Aikido MCP server**

```shellscript
devin mcp add -s user aikido -- npx -y @aikidosec/mcp@latest
```

Or add it by hand to `~/.config/devin/mcp_config.json` (`%APPDATA%\devin\mcp_config.json` on Windows).

```json
{
  "mcpServers": {
    "aikido": {
      "command": "npx",
      "args": ["-y", "@aikidosec/mcp@latest"]
    }
  }
}
```

{% hint style="info" %}
Use `-s project` instead of `-s user` to write the server to `.devin/mcp_config.json` and share it with your team.
{% endhint %}
{% endstep %}

{% step %}
**Authenticate**

On first use, the MCP server will detect no authentication is set. Ask the agent to run `aikido_login`. It will return a browser link. Follow the link to log in to Aikido. Your token is stored securely once login is complete.
{% endstep %}

{% step %}
**Add the Aikido rule to your global AGENTS file**

Create the Devin config directory if it doesn't exist yet.

```shellscript
mkdir -p ~/.config/devin/
```

Append the Aikido rule to `~/.config/devin/AGENTS.md`. Rules in this file are always on.

```shellscript
curl -fsSL "https://gist.githubusercontent.com/kidk/aa48cad6db80ba4a38493016aae67712/raw/3644397b7df43423e3da06434491b40bbb79dd47/aikido-rule.txt" \
  >> ~/.config/devin/AGENTS.md
```

{% hint style="info" %}
To scope the rule to one repository instead, add the same content to the repository's `AGENTS.md`, or to `.devin/rules/aikido.md` with `trigger: always_on` frontmatter.
{% endhint %}
{% endstep %}

{% step %}
**Finished**

Aikido MCP is now available in Devin.

{% hint style="info" %}
Restart Devin Desktop if it was open.
{% endhint %}
{% endstep %}
{% endstepper %}

Devin namespaces MCP tools as `mcp__aikido__<tool>` and puts them behind its permission system, so approve the Aikido tools the first time the agent calls them.

For ways to ask your agent to use Aikido, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).
