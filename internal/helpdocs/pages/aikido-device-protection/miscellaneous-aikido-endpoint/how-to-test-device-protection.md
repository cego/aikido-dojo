# How to test if Device Protection is working

Use one of the methods below to verify that Aikido Device Protection is active on your machine.

## Quick check (recommended)

Open the link below in your browser. The page will immediately show whether your device is protected or not.

{% embed url="<https://device-protection.aikido.help/>" %}

## Package-based test

Install a safe test package in your preferred language. These packages are harmless and exist solely to trigger a detection.

{% tabs %}
{% tab title="npm" %}

```bash
npm install aikido-endpoint-test
```

{% endtab %}

{% tab title="pip" %}

```bash
pip install safe-chain-pi-test
```

{% endtab %}

{% tab title="NuGet" %}
<https://www.nuget.org/packages/AikidoEndpointTest>
{% endtab %}

{% tab title="Go" %}

```
go install github.com/Aikido-demo-apps/endpoint-golang-test@latest
```

{% endtab %}
{% endtabs %}

## IDE extension test

Install the test extension in your editor. It will trigger a detection if Device Protection is running.

* **VS Code:** [Aikido Endpoint Test](https://marketplace.visualstudio.com/items?itemName=AikidoSecurity.aikido-endpoint-test)
* **Open VSX (e.g. Cursor, Windsurf):** [Aikido Endpoint Test](https://open-vsx.org/extension/AikidoSecurity/aikido-endpoint-test)

Test package or extension not blocked? See [Why is a Package Not Blocked?](/aikido-device-protection/using-aikido-endpoint/why-is-a-package-not-blocked.md)

## Automated test

```shellscript
#!/bin/bash
curl -s https://device-protection.aikido.help/ | grep -q "Device Protection Active"
```

This script will return exit code `0` when Device Protection is active, and exit code `1` when it's not.
