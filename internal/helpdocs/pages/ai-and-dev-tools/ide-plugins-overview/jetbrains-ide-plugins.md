# JetBrains IDE

Aikido integrates with the majority of Jetbrains IDEs and scans your codebases for **secrets, API keys SAST** code issues and **SCA**. It runs scans whenever you open or save a file.

Every time you make and save changes in a file, a scan runs. If any issues are detected, they are highlighted in the editor and also displayed in the Problems panel. When you hover over a detected SAST issue, additional context about the problem is provided.

[image: Code editor highlighting security warnings in an Express.js application, including missing security headers.]

## Supported IDEs <a href="#supported-ides" id="supported-ides"></a>

We support the majority of Jetbrain IDEs

* IntelliJ IDEA (Java/Kotlin/Spring)
* GoLand (Go/JS/TS)
* PhpStorm (PHP/Laravel/Symfony)
* PyCharm (Python/Django)
* Rider (C#/.NET/ASP.NET)
* WebStorm (JS/TS/React)
* RubyMine (Ruby/Rails)
* RustRover (Rust)
* Android Studio (Android)
* CLion (C and C++)

## How to Install and Test <a href="#how-to-install-and-test" id="how-to-install-and-test"></a>

> This plugin is only available on paid plans.

**Step 1.** Head over to the [JetBrains Marketplace](https://plugins.jetbrains.com/plugin/24993-aikido-security) and click **Get.** After installation, you will be asked to add your personal access token (step 2).

**Step 2.** In Aikido, go to the [JetBrains IDE Integration Screen](https://app.aikido.dev/settings/integrations/ide/jetbrains) and create your token.

[image: JetBrains IDE Plugin user access and activity log with token generation options.]

**Step 3.** Check out the examples in our docs on the [JetBrains Marketplace](https://plugins.jetbrains.com/plugin/24993-aikido-security) to test whether everything works well.

## Troubleshooting

{% content-ref url="/pages/c6y7HYvl2fY0rM3QArnn" %}
[JetBrains Plugin - Access Token Not Stored](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/jetbrains-plugin-access-token-not-stored.md)
{% endcontent-ref %}

{% content-ref url="/pages/NAH7p1hgS3B4GR19LQ48" %}
[Jetbrains - Information for Support](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/jetbrains-information-for-support.md)
{% endcontent-ref %}
