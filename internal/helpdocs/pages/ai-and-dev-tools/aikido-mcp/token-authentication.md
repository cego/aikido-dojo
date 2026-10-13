# Alternative Token Authentication

Token authentication is the alternative to the default browser-based login flow. Use it when:

* You are running Aikido MCP in a CI/CD pipeline or other headless environment
* Your setup cannot open a browser (e.g. remote servers, Docker containers)
* You prefer to manage credentials explicitly via environment variables

For interactive use on a local machine, the browser-based login flow is simpler and recommended.

## Creating a token

In Aikido, go to [Settings → Integrations → IDE → MCP](https://app.aikido.dev/settings/integrations/ide/mcp) and create a Personal Access Token (PAT). Copy it to a secure location.

## Setup by agent

### Claude Code

Pass the token as an argument to the setup skill:

```
/aikido:setup your-PAT-here
```

This saves the token to your Claude Code user settings and registers the MCP server automatically.

Alternatively, set `AIKIDO_API_KEY` in your environment before starting Claude Code:

{% tabs %}
{% tab title="macOS / Linux" %}

```shellscript
export AIKIDO_API_KEY=your-token-here
claude
```

{% endtab %}

{% tab title="Windows PowerShell" %}

```powershell
$env:AIKIDO_API_KEY="your-token-here"
claude
```

{% endtab %}
{% endtabs %}

### Cursor

When adding Aikido from the [Cursor Marketplace](https://cursor.com/marketplace/aikido), set your token as the `AIKIDO_API_KEY` environment variable when prompted during installation.

For a manual setup, add the token to your Cursor MCP configuration:

```json
{
  "mcpServers": {
    "aikido": {
      "command": "npx",
      "args": ["-y", "@aikidosec/mcp@latest"],
      "env": {
        "AIKIDO_API_KEY": "your-token-here"
      }
    }
  }
}
```

### OpenAI Codex CLI

Pass the token via the `--env` flag when registering the MCP server:

```shellscript
codex mcp add aikido \
  --env AIKIDO_API_KEY=your-token-here \
  -- npx -y @aikidosec/mcp@latest
```

### Gemini CLI

```shellscript
gemini mcp add aikido \
  --env AIKIDO_API_KEY=your-token-here \
  npx -y @aikidosec/mcp@latest
```

### JetBrains

Add the token to the `environment` block in the MCP server JSON config (Settings | Tools | AI Assistant | Model Context Protocol):

```json
{
  "mcpServers": {
    "aikido": {
      "command": "npx -y @aikidosec/mcp@latest",
      "args": [],
      "environment": {
        "AIKIDO_API_KEY": "your-token-here"
      }
    }
  }
}
```

### GitHub Copilot in IDE

Add the token to the `env` block in your `mcp.json`:

```json
{
  "servers": {
    "aikido": {
      "command": "npx",
      "args": ["-y", "@aikidosec/mcp@latest"],
      "env": {
        "AIKIDO_API_KEY": "your-token-here"
      }
    }
  }
}
```

### GitHub Copilot Cloud Agents

GitHub Copilot only exposes secrets to MCP config when their names start with `COPILOT_MCP_`. Create a repository secret:

`COPILOT_MCP_AIKIDO_API_KEY` = your Aikido token

Then use that variable name in your MCP config under Settings → Copilot → Coding agent:

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

### Mistral Vibe

Add the token as an environment variable in `~/.vibe/config.toml`:

```toml
[[mcp_servers]]
name = "aikido"
transport = "stdio"
command = "sh"
args = ["-c", "AIKIDO_API_KEY='your-token-here' npx -y @aikidosec/mcp@latest"]
```

### OpenCode

Add the token to the `environment` block in `~/.config/opencode/opencode.json`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "aikido": {
      "type": "local",
      "command": ["npx", "-y", "@aikidosec/mcp@latest"],
      "enabled": true,
      "environment": {
        "AIKIDO_API_KEY": "your-token-here"
      }
    }
  }
}
```
