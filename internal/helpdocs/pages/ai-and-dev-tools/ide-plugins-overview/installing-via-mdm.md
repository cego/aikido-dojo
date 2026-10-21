# Installing IDE Plugins via MDM

MDM (mobile device management) tools let IT teams pre-install software on developer machines at scale. This page covers deploying the Aikido plugin to JetBrains IDEs and VS Code-based editors using the most common MDM platforms.

## JetBrains IDEs

**Plugin ID:** `dev.aikido.aikido-intellij-plugin`

JetBrains IDEs accept an `installPlugins` flag on the command line. The IDE must be fully closed before running the command. The syntax is the same across all JetBrains products; only the executable name changes (e.g., `idea`, `webstorm`, `phpstorm`, `rider`). See the full list in the [JetBrains docs](https://www.jetbrains.com/help/idea/install-plugins-from-the-command-line.html).

### macOS

Deploy the script below as a Jamf policy script or a Kandji custom script. Replace `PhpStorm` with the IDE your team uses.

```bash
#!/bin/bash
PLUGIN_ID="dev.aikido.aikido-intellij-plugin"

osascript -e 'quit app "PhpStorm"'
sleep 5
open -na "PhpStorm.app" --args installPlugins "$PLUGIN_ID"
```

{% hint style="warning" %}
The IDE must have been launched at least once on the machine before this script will work.
{% endhint %}

### Windows

Deploy as a PowerShell script or package it as a Win32 app. Replace the path and executable with the IDE installed on your machines.

```powershell
& "C:\Program Files\JetBrains\IntelliJ IDEA\bin\idea64.exe" installPlugins dev.aikido.aikido-intellij-plugin
```

## VS Code, Cursor, Windsurf and other variants

**Extension ID:** `AikidoSecurity.aikido`

All VSCode editors support installing extensions from the command line. Use the command for whichever editors your team has installed.

### macOS

Deploy as a Jamf policy script or Kandji custom script.

```bash
#!/bin/bash
code --install-extension AikidoSecurity.aikido      # VS Code
cursor --install-extension AikidoSecurity.aikido    # Cursor
windsurf --install-extension AikidoSecurity.aikido  # Windsurf
```

### Windows

Deploy as a PowerShell remediation script.

```powershell
code --install-extension AikidoSecurity.aikido      # VS Code
cursor --install-extension AikidoSecurity.aikido    # Cursor
windsurf --install-extension AikidoSecurity.aikido  # Windsurf
```

For managing VS Code extensions at scale with Intune, see the [VS Code enterprise documentation](https://code.visualstudio.com/docs/enterprise/overview).
