# AI Coding Assistants (MCP)

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

## Demo

Demo of the Aikido MCP server working with an agent rule to scan and fix vulnerabilities in AI generated code:

{% embed url="<https://www.youtube.com/watch?v=D0ltRTSuKmk>" %}

## Installation

### AI Platforms

* [Anthropic Claude Code MCP](/ai-and-dev-tools/aikido-mcp/anthropic-claude-code-mcp.md)
* [Cursor MCP](/ai-and-dev-tools/aikido-mcp/cursor-mcp.md)
* [OpenAI Codex CLI MCP](/ai-and-dev-tools/aikido-mcp/openai-codex-cli-mcp.md)
* [Gemini CLI MCP](/ai-and-dev-tools/aikido-mcp/gemini-cli-mcp.md)
* [Jetbrains AI](/ai-and-dev-tools/aikido-mcp/jetbrains-ai.md)
* [Github Copilot](/ai-and-dev-tools/aikido-mcp/github-copilot.md)
* [Mistral Vibe MCP](/ai-and-dev-tools/aikido-mcp/mistral-vibe-mcp.md)
* [OpenCode MCP](/ai-and-dev-tools/aikido-mcp/opencode-mcp.md)
* [Devin MCP](/ai-and-dev-tools/aikido-mcp/devin-mcp.md)

### Via Aikido IDE plugins

When the [Aikido IDE plugin](/ai-and-dev-tools/ide-plugins-overview.md) is installed you can use the Aikido Expansion Packs to install the Aikido MCP server with one click. [Learn more in the Expansion Packs docs.](/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs.md)

{% hint style="warning" %}
Currently available for Jetbrains IDE's, VS Code and variants (e.g., Windsurf, Cursor, Kiro, and AntiGravity).
{% endhint %}

* [Cursor IDE](/ai-and-dev-tools/ide-plugins-overview/cursor-ide.md)
* [Google Antigravity IDE](/ai-and-dev-tools/ide-plugins-overview/google-antigravity.md)
* [JetBrains IDE](/ai-and-dev-tools/ide-plugins-overview/jetbrains-ide-plugins.md)
* [Kiro IDE](/ai-and-dev-tools/ide-plugins-overview/kiro-ide.md)
* [VS Code IDE](/ai-and-dev-tools/ide-plugins-overview/vs-code-plugin.md)
* [Windsurf IDE](/ai-and-dev-tools/ide-plugins-overview/windsurf-ide.md)

### Manual installation for other platforms

For any other AI platform or custom MCP setup, refer to the [npm package page for detailed manual installation instructions](https://www.npmjs.com/package/@aikidosec/mcp).

## Using the MCP server

For the tools the server exposes, see [Available tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools). For prompts that put them to work, see [Example prompts](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#example-prompts).

## Pull request comments

Use the MCP of your source control system to fetch pull request comments and to react to them. Aikido recommends this for every provider. Have the agent reply and interact on those comments the same way a person would.

See [Fetch Pull Request Comments](/ai-and-dev-tools/aikido-mcp/fetch-pull-request-comments.md) for more information.
