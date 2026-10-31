# Fetch Pull Request Comments

Install the aikido-mcp-comments skill so your AI agent can read comments already posted on a pull request, including Aikido's scan results.

The `aikido-mcp-comments` skill lets your AI agent read comments that are already posted on a pull request (PR) or GitLab merge request: CI bot results, review feedback, and Aikido's own scan-result comment. The agent fetches them with your source control CLI (`gh`, `glab` or `az`), not with the Aikido MCP server.

{% hint style="info" %}
The Aikido plugins for Claude Code, Cursor and Codex already include this skill. Use this page for any other agent.
{% endhint %}

## Install with the skills CLI

The skill is listed on [skills.sh](https://skills.sh/aikidosec/aikido-skills/aikido-mcp-comments). The `skills` CLI asks which agents to install it for and puts it in the right folder for each one. It needs Node.js.

```shell
npx skills add AikidoSec/aikido-skills --skill aikido-mcp-comments -g
```

`-g` installs it for your user, so it works in every project. Leave it out to install it only in the current project. To update the skill later, run `npx skills update`.

## Install manually

{% stepper %}
{% step %}

### Create the skill folder

```shell
mkdir -p ~/.agents/skills/aikido-mcp-comments
```

{% endstep %}

{% step %}

### Download the skill files

The skill has two files: `SKILL.md` with the instructions and `reference.md` with the commands for each platform.

```shell
curl -fsSL "https://raw.githubusercontent.com/AikidoSec/aikido-skills/main/skills/aikido-mcp-comments/SKILL.md" \
  -o ~/.agents/skills/aikido-mcp-comments/SKILL.md
curl -fsSL "https://raw.githubusercontent.com/AikidoSec/aikido-skills/main/skills/aikido-mcp-comments/reference.md" \
  -o ~/.agents/skills/aikido-mcp-comments/reference.md
```

{% endstep %}

{% step %}

### Restart your agent

Restart the agent so it loads the new skill. To update the skill later, run the download commands again.
{% endstep %}
{% endstepper %}

## Example prompts

* "What did Aikido say on this PR?"
* "Summarize the comments on `org/repo#123`"
* "What did CI report on my merge request?"

## Related

{% content-ref url="/pages/iyRExFqpXza0GqBsgwjz" %}
[Using Aikido MCP](/ai-and-dev-tools/aikido-mcp/using-aikido-mcp.md)
{% endcontent-ref %}

{% content-ref url="/pages/XBHelJe7KcUn7BUyGhF0" %}
[AI Coding Assistants (MCP)](/ai-and-dev-tools/aikido-mcp.md)
{% endcontent-ref %}
