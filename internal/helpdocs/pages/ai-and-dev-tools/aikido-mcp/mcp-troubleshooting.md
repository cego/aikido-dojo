# MCP Troubleshooting

When running into issues with the Aikido MCP server, work through the steps below before reaching out to support. If the problem persists, share the output from these steps with the Aikido team.

## Check your Node.js version

The Aikido MCP server requires Node.js 18.19.0 or higher.

Run this in your terminal to verify:

```bash
node -v
```

If your version is below 18.19.0, update Node.js before continuing.

## Run the MCP server directly

Most AI coding tools don't make it easy to get detailed logs for plugin-integrated MCP servers. If the steps above didn't help, try running the MCP server on its own to isolate the issue.

Run the following command in your terminal:

```bash
LOG_LEVEL=debug npx -y @aikidosec/mcp@latest
```

This starts the MCP server outside of any AI tool and triggers the browser login flow if you are not yet authenticated. It outputs debug-level logs directly in your terminal.

If you are using [token authentication](/ai-and-dev-tools/aikido-mcp/token-authentication.md), pass the token directly instead:

```bash
AIKIDO_API_KEY=your_api_key_here LOG_LEVEL=debug npx -y @aikidosec/mcp@latest
```

Look for errors related to authentication, network connectivity, or missing configuration.

## Collect Claude debug logs

If you are using Claude, start it in debug mode:

```bash
claude --debug
```

This writes startup logs (including MCP startup logs) to a file. These logs can contain useful information about connection or configuration problems, but they won't always surface the root cause.
