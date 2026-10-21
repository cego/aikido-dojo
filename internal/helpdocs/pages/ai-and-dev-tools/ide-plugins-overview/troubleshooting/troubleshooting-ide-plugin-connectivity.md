# IDE Plugin Connectivity

## How the plugin uses the Network

The IDE plugin processes code **locally on your device.** Your source code is never uploaded to Aikido servers. Cloud connectivity is used for the following:

| Direction   | What is transferred                                                                                         |
| ----------- | ----------------------------------------------------------------------------------------------------------- |
| Cloud → IDE | Scanning rules and configuration (e.g. SAST/secrets rulesets, opengrep rules)                               |
| IDE ↔ Cloud | Scan findings (issues) for scoring, deduplication, and applying your organisation’s ignore and triage rules |
| IDE ↔ Cloud | AutoFix and AutoTriage requests are sent to Aikido’s backend for processing and returned to the IDE         |

This applies to SAST, secrets, SCA, and IaC scans. The same data-handling policies that apply to your connected Aikido workspace apply to the IDE plugin.

## Required Domains

To function correctly, the plugin must be able to access the following domains over HTTPS (port 443):

```
*.aikido.dev
or
ide.aikido.dev
```

## Troubleshooting

If the Aikido IDE plugin cannot connect or login is not persistent, it’s likely due to a network restriction or firewall rule.

### 1. Check connectivity from your terminal

**macOS / Linux (bash)**

```
# Test DNS resolution
nslookup ide.aikido.dev

# Test HTTPS connectivity
curl -v https://ide.aikido.dev/ping

# Or check if the port is open
nc -vz ide.aikido.dev 443
```

**Windows (PowerShell)**

```
# Test DNS resolution
Resolve-DnsName ide.aikido.dev

# Test HTTPS connectivity
Invoke-WebRequest -Uri https://ide.aikido.dev/ping -UseBasicParsing

# Test port connectivity
Test-NetConnection ide.aikido.dev -Port 443
```

### 2. Verify proxy or firewall settings

* Ensure your proxy or corporate firewall allows outbound HTTPS connections to \*.aikido.dev.
* If you’re using a VPN, try disabling it temporarily to see if it affects connectivity.
* Confirm no SSL interception or certificate pinning issues are blocking access.

[Follow this guide for more information about configuring proxy and CA within your IDE.](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/ide-plugin-proxy-and-ca-settings.md)

### 3. Still having issues?

If all tests fail or you receive timeouts, please share the command output with the Aikido team so we can help diagnose further.

{% content-ref url="/pages/jLKPaQyzd9YFiTyIsyzg" %}
[Visual Studio - Information for Support](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/visual-studio-information-for-support.md)
{% endcontent-ref %}

{% content-ref url="/pages/NAH7p1hgS3B4GR19LQ48" %}
[Jetbrains - Information for Support](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/jetbrains-information-for-support.md)
{% endcontent-ref %}
