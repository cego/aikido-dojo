# AI Coding Assistants (MCP)

The Aikido MCP Plugin connects Aikido’s security engine to AI coding tools. It automatically scans AI generated code for vulnerabilities and hardcoded secrets as soon as it is created.

AI assistants can review their own output, but that review is not perfect. Aikido adds a reliable and consistent security layer that checks every generated snippet with proven scanning rules.

**Why use Aikido MCP**

* Deterministic, independent security checks on every AI generated snippet before it is committed
* Immediate detection and remediation of vulnerabilities and hardcoded secrets in AI assisted workflows
* Real time feedback, making AI driven development safer by default

## Demo

Demo of the Aikido MCP server working with an agent rule to scan and fix vulnerabilities in AI generated code:

{% embed url="https://www.youtube.com/watch?v=D0ltRTSuKmk" %}

## Installation

### AI Platforms

* [Anthropic Claude Code MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/anthropic-claude-code-mcp)
* [Cursor MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/cursor-mcp)
* [OpenAI Codex CLI MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/openai-codex-cli-mcp)
* [Gemini CLI MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/gemini-cli-mcp)
* [Jetbrains AI](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/jetbrains-ai)
* [Github Copilot](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/github-copilot)
* [Mistral Vibe MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/mistral-vibe-mcp)
* [OpenCode MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/opencode-mcp)
* [Devin MCP](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/devin-mcp)

### Via Aikido IDE plugins

When the [Aikido IDE plugin](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview) is installed you can use the Aikido Expansion Packs to install the Aikido MCP server with one click. [Learn more in the Expansion Packs docs.](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs)

{% hint style="warning" %}
Currently available for Jetbrains IDE's, VS Code and variants (e.g., Windsurf, Cursor, Kiro, and AntiGravity).
{% endhint %}

* [Cursor IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/cursor-ide)
* [Google Antigravity IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/google-antigravity)
* [jetbrains-ide-plugins.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/jetbrains-ide-plugins "mention")
* [Kiro IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/kiro-ide)
* [VS Code IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/vs-code-plugin)
* [Windsurf IDE](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/windsurf-ide)

### Manual installation for other platforms

For any other AI platform or custom MCP setup, refer to the [npm package page for detailed manual installation instructions](https://www.npmjs.com/package/@aikidosec/mcp).

## Using the MCP server

For the tools the server exposes, see [Available tools](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/using-aikido-mcp#available-tools). For prompts that put them to work, see [Example prompts](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/using-aikido-mcp#example-prompts).

## Pull request comments

Use the MCP of your source control system to fetch pull request comments and to react to them. Aikido recommends this for every provider. Have the agent reply and interact on those comments the same way a person would.

See [Fetch Pull Request Comments](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp/fetch-pull-request-comments) for more information.
