# Bundle Aikido in Your Own Claude Code Plugin

Bundle the Aikido MCP server and your own skills into an internal Claude Code plugin, so developers get a fixed setup with nothing to configure.

The [official Aikido plugin](/ai-and-dev-tools/aikido-mcp/anthropic-claude-code-mcp.md) always ships with its own skills. If you want to decide exactly which skills your developers get, build your own internal plugin that bundles the Aikido MCP server together with the skills you pick, then distribute it from your own marketplace. A plugin is a folder of configuration that Claude Code loads as one unit: skills, agents, hooks, and MCP servers.

Use this when:

* You want a fixed set of skills: your own, a subset of Aikido's, or both
* Developers should not have to install anything or run a setup step
* You want the rollout versioned and reviewed in your own repository

{% hint style="info" %}
Don't install the official Aikido plugin alongside your own. Both would register the Aikido MCP server and you would get two sets of skills.
{% endhint %}

## Build the plugin

{% stepper %}
{% step %}
**Create the plugin folder**

```
aikido-internal/
├── .claude-plugin/
│   └── plugin.json
├── .mcp.json
└── skills/
    └── scan/
        └── SKILL.md
```

Only `plugin.json` goes inside `.claude-plugin/`. Everything else sits at the plugin root.

`.claude-plugin/plugin.json`:

```json
{
  "name": "aikido-internal",
  "displayName": "Aikido (internal)",
  "description": "Aikido security scanning with our own skills",
  "version": "1.0.0"
}
```

The `name` is the namespace for your skills, so a skill folder called `scan/` becomes `/aikido-internal:scan`.
{% endstep %}

{% step %}
**Add the Aikido MCP server**

`.mcp.json` at the plugin root:

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

The server starts automatically when the plugin is enabled. The first time someone uses an Aikido tool, the server asks them to sign in through the browser.

For headless environments, or to skip the browser sign-in, use a Personal Access Token (PAT). The server picks up `AIKIDO_API_KEY` from the environment it inherits, so set it in the developer's shell profile or through your device management tool and leave the config above unchanged. That keeps the token out of the plugin repository. See [Token Authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md) for how to create a PAT.

If you would rather have Claude Code ask each developer for a token once and store it securely, declare it in `plugin.json` under `userConfig` and reference it as `${user_config.aikido_token}` in the `env` block:

```json
{
  "userConfig": {
    "aikido_token": {
      "type": "string",
      "title": "Aikido token",
      "description": "Your Aikido Personal Access Token",
      "sensitive": true
    }
  }
}
```

Then pass it to the server in `.mcp.json`:

```json
"env": {
  "AIKIDO_API_KEY": "${user_config.aikido_token}"
}
```

{% endstep %}

{% step %}
**Add the skills you want**

Each skill is a folder under `skills/` with a `SKILL.md` file. The `description` tells Claude when to use it.

`skills/scan/SKILL.md`:

```markdown
---
description: Scan changed code with Aikido for vulnerabilities and hardcoded secrets. Use after writing or editing code, and before committing.
---

Run the Aikido `aikido_scan_paths` tool on the files changed in the working tree.

Report every finding with its severity and file location, then fix the ones
you introduced. Do not commit while a critical finding is open.
```

Add as many as you need, for example a triage skill that uses `aikido_issues_list` to pull open issues for the current repository. See the [available MCP tools](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md#available-tools) for what the server exposes.
{% endstep %}

{% step %}
**Test it locally**

```shellscript
claude --plugin-dir ./aikido-internal
```

Then run `/aikido-internal:scan` to check the skill, and `/plugin` to confirm the Aikido MCP server started. Run `/reload-plugins` to pick up edits without restarting.

Validate the structure before you ship:

```shellscript
claude plugin validate ./aikido-internal
```

{% endstep %}
{% endstepper %}

## Distribute it to your team

Put the plugin in a repository your developers can read, with a marketplace file at `.claude-plugin/marketplace.json` in the repository root:

```json
{
  "name": "acme-tools",
  "owner": {
    "name": "Platform Team",
    "email": "platform@acme.com"
  },
  "plugins": [
    {
      "name": "aikido-internal",
      "source": "./aikido-internal",
      "description": "Aikido security scanning with our own skills"
    }
  ]
}
```

A private repository works. Claude Code uses each developer's existing Git credentials to clone it.

Then commit this to your project's `.claude/settings.json`:

```json
{
  "extraKnownMarketplaces": {
    "acme-tools": {
      "source": {
        "source": "github",
        "repo": "acme/claude-plugins"
      }
    }
  },
  "enabledPlugins": {
    "aikido-internal@acme-tools": true
  }
}
```

Once a developer trusts the repository folder, Claude Code adds the marketplace and enables the plugin with no further prompt.

{% hint style="warning" %}
Keep the plugin folder inside the marketplace repository and reference it with a relative `./` path, as above. A plugin that points at an external source, such as another GitHub repository, is not installed by the settings file alone: each developer has to run `claude plugin install` first.
{% endhint %}

## Control what the tools can do

Bundling your own plugin controls which skills exist on the developer's machine. It does not change what the Aikido tools are allowed to do. That stays a workspace-level setting in Aikido, see [MCP Permissions](/ai-and-dev-tools/aikido-mcp/mcp-permissions.md).
