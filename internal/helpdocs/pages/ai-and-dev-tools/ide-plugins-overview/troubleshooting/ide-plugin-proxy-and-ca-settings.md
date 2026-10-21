# IDE Plugin Proxy and CA Settings

If your organization routes traffic through a corporate proxy or uses a custom Certificate Authority (CA) for SSL/TLS inspection, the Aikido IDE plugin may not be able to connect. You will need to configure your IDE's proxy and/or CA settings so the plugin can reach `*.aikido.dev`.

Verify you can access `ide.aikido.dev` via proxy before continuing with next steps. You can [find test commands on our troubleshooting connectivity page.](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/troubleshooting-ide-plugin-connectivity.md)

### VS Code (and variants like Cursor, Windsurf)

{% content-ref url="/pages/yS43oHZnXnEoeOXnQlN3" %}
[VS Code Proxy Setup](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/vs-code-proxy-setup.md)
{% endcontent-ref %}

### JetBrains (IntelliJ IDEA, WebStorm, PyCharm, etc.)

{% content-ref url="/pages/95U6ZucURUXNVpoJvpYL" %}
[Jetbrains Proxy Setup](/ai-and-dev-tools/ide-plugins-overview/troubleshooting/jetbrains-proxy-setup.md)
{% endcontent-ref %}

### Visual Studio

1. Configure your proxy using the built-in Proxy Configuration page or your Windows system proxy. See [Configure proxy settings in Visual Studio](https://learn.microsoft.com/en-us/visualstudio/ide/reference/proxy-settings) and [Troubleshoot proxy and firewall issues](https://learn.microsoft.com/en-us/visualstudio/ide/reference/proxy-firewall-troubleshoot).
2. Add your custom CA certificate to the Windows Trusted Root Certification Authorities store. See [Troubleshoot proxy and firewall issues - SSL certificate errors](https://learn.microsoft.com/en-us/visualstudio/ide/reference/proxy-firewall-troubleshoot).

### Eclipse

1. Configure your proxy under Window > Preferences > General > Network Connections. Set the Active Provider to Manual and configure the HTTP and HTTPS entries. See [Eclipse Network Connections](https://help.eclipse.org/latest/topic/org.eclipse.platform.doc.user/reference/ref-net-preferences.htm).
2. Add your custom CA certificate to your OS trust store. If Eclipse still shows certificate errors, import the certificate into Eclipse's JRE trust store using `keytool`. See [Eclipse Network Connections](https://help.eclipse.org/latest/topic/org.eclipse.platform.doc.user/reference/ref-net-preferences.htm).

### Still having issues?

Check the connectivity steps in [Troubleshooting IDE Plugin Connectivity](https://help.aikido.dev/ai-and-dev-tools/ide-plugins-overview/troubleshooting/troubleshooting-ide-plugin-connectivity). If the issue persists, contact the Aikido support team with your debug logs and connectivity check output.
