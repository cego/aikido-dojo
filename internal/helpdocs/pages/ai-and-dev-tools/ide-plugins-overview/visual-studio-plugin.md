# Visual Studio IDE

Aikido automatically scans your projects for hardcoded secrets (API keys, tokens) and insecure code patterns (SQL injections, path traversal, ..) so you can catch issues early and keep your codebase safe.

Scans run automatically whenever you open a file or save changes, making it easy to catch issues early in development.

When security issues are found, they're highlighted directly in your code and listed in the Aikido window.

[image: C# code editor showing SQL query function, warnings, and security/code scan panels.]

## How to Install and Use <a href="#how-to-install-and-use" id="how-to-install-and-use"></a>

**Step 1.** Head over to the [Visual Studio Marketplace](https://marketplace.visualstudio.com/items?itemName=AikidoSecurity.aikido-visualstudio) and click **Install.** After installation, you will be asked to add your personal access token (step 2).

**Step 2.** In Aikido, go to the [Visual Studio Integration Screen](https://app.aikido.dev/settings/integrations/ide/visualstudio) and create your token.

[image: Visual Studio IDE Plugin page showing no registered tokens and option to generate a new PAT.]

**Step 3.** Check out the examples in our docs on the [Visual Studio Marketplace to test](https://marketplace.visualstudio.com/items?itemName=AikidoSecurity.aikido-visualstudio)[ ](https://marketplace.visualstudio.com/items?itemName=AikidoSecurity.aikido)whether everything works well.

## Supported Visual Studio Versions <a href="#supported-visual-studio-versions" id="supported-visual-studio-versions"></a>

The Aikido Extension can be installed on Visual Studio Community Edition, Visual Studio Professional and Visual Studio Enterprise. It adheres to [Visual Studio’s support lifecycle](https://learn.microsoft.com/en-us/visualstudio/productinfo/vs-servicing), remaining compatible with officially supported versions.

## Troubleshooting

{% content-ref url="/pages/jLKPaQyzd9YFiTyIsyzg" %}
[Visual Studio - Information for Support](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/visual-studio-information-for-support.md)
{% endcontent-ref %}
