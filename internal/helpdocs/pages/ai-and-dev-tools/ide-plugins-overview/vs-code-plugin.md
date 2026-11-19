# VS Code IDE

Aikido automatically scans your projects for hardcoded secrets (API keys, tokens) and insecure code patterns (SQL injections, path traversal, ..) so you can catch issues early and keep your codebase safe.

Scans run automatically whenever you open a file or save changes, making it easy to catch issues early in development.

When security issues are found, they're highlighted directly in your code and listed in the Aikido window.

{% stepper %}
{% step %}
## Open Extensions and Install "Aikido Security"

Alternatively use these links to go to the Marketplaces

* [VS Code](https://marketplace.visualstudio.com/items?itemName=AikidoSecurity.aikido)
* [Windsurf](https://open-vsx.org/extension/AikidoSecurity/aikido) / [Cursor](https://open-vsx.org/extension/AikidoSecurity/aikido) / [Kiro](https://open-vsx.org/extension/AikidoSecurity/aikido) / [Google Antigravity](https://open-vsx.org/extension/AikidoSecurity/aikido)
{% endstep %}

{% step %}
## Authenticate with Aikido

 Open the Aikido plugin by clicking on the sidebar icon and click on "Connect to Aikido" to authenticate with Aikido platform.

Alternatively you can open up the Command Palette and run `Aikido: Log In`

If the automated authentication does not work you can manually create a personal access within Aikido by going to the [Integrations page and following the instructions](https://app.aikido.dev/settings/integrations?section=ide&regionRedirect=1).
{% endstep %}

{% step %}
## Try out our examples

Below you can find an example `index.js` file that can be used to verify if the extension is working correctly, it should detect one SAST issue (SQL injection) and one exposed secret (SQL Server connection string).

```javascript
const app = {}

app.get("/user", (req, res) => {
    const connStr = "Server=tcp:myserver.database.windows.net,1433;Initial Catalog=mydb;Persist Security Info=False;User ID=myuser;Password=$uperSecret123!@#";
    const username = req.query.username
    const unsafeQuery = `SELECT * FROM users WHERE username = '${username}'`
    sql.connect(connStr).query(unsafeQuery, (err, result) => {
        res.status(200).send(result)
    })
})
```
{% endstep %}

{% step %}
## Turn on Additional Security Tooling

Extend Aikido in your IDE with Expansion Packs like [MCP for AI agents](https://help.aikido.dev/ai-and-dev-tools/aikido-mcp), [pre-commit hooks](https://help.aikido.dev/ai-and-dev-tools/aikido-secrets-pre-commit-hook), and [Safe Chain](https://help.aikido.dev/code-scanning/aikido-malware-scanning). For more details, see the documentation below.

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs.md" %}
[aikido-expansion-packs.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-expansion-packs)
{% endcontent-ref %}
{% endstep %}
{% endstepper %}

Now that the plugin is installed, you can dive into the features that help you spot security issues while you work:

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/real-time-code-scanning-in-ide.md" %}
[real-time-code-scanning-in-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/real-time-code-scanning-in-ide)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/code-quality-in-ide.md" %}
[code-quality-in-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/code-quality-in-ide)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/open-source-dependency-scanning-sca-in-ide.md" %}
[open-source-dependency-scanning-sca-in-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/open-source-dependency-scanning-sca-in-ide)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/full-workspace-scan-in-ide.md" %}
[full-workspace-scan-in-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/full-workspace-scan-in-ide)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide.md" %}
[aikido-ai-in-ide.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/features/aikido-ai-in-ide)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/troubleshooting/vs-code-extension-keeps-disconnecting.md" %}
[vs-code-extension-keeps-disconnecting.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/troubleshooting/vs-code-extension-keeps-disconnecting)
{% endcontent-ref %}

{% content-ref url="../../ai-and-dev-tools/ide-plugins-overview/troubleshooting/vs-code-information-for-support.md" %}
[vs-code-information-for-support.md](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/troubleshooting/vs-code-information-for-support)
{% endcontent-ref %}
