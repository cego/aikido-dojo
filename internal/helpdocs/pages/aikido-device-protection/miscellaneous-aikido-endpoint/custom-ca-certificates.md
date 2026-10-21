# Adding a Custom CA Certificate

To inspect package manager traffic, Device Protection uses a local Certificate Authority (CA) and builds CA bundles for the tools it configures at install. Those tools are pointed at Aikido's bundle instead of their own default trust list, so any CA that isn't in the bundle is no longer trusted by them.

That matters if your organization runs its own CA, for example a corporate proxy that inspects TLS traffic (Zscaler, Netskope) or a private registry served with an internal certificate. Adopting your CA adds it to every bundle Aikido builds, and keeps it there whenever Aikido rebuilds or repairs those bundles.

## When you need this

Adopt a CA when, after Device Protection is installed, your tools fail to reach an internal host or a registry behind TLS inspection. The errors typically mention a self-signed certificate, an unknown or untrusted issuer, or a certificate that could not be verified.

If the failing certificate is Aikido's own CA rather than your organization's, see [Fixing SSL/TLS Certificate Errors](/aikido-device-protection/miscellaneous-aikido-endpoint/ssl-certificate-errors.md) instead.

## Before you start

* Your CA certificate as a PEM file on the device, for example `corporate-ca.pem`.
* Administrator access. The command must run as root on macOS and Linux, or from an elevated prompt on Windows.

## Adopt the CA

{% tabs %}
{% tab title="macOS / Linux" %}

```bash
sudo aikidoctl certconfig adopt /path/to/corporate-ca.pem
```

{% endtab %}

{% tab title="Windows" %}
Run from an elevated Command Prompt or PowerShell:

```
aikidoctl.exe certconfig adopt C:\path\to\corporate-ca.pem
```

{% endtab %}
{% endtabs %}

Aikido adds the CA to the bundles it builds and reinstalls the trust configuration for each tool.

{% hint style="info" %}
Shells and services that were already running keep their old environment. Open a new login shell before you test, or restart the affected service.
{% endhint %}

## What adopting does and doesn't change

* It covers every CA bundle Aikido builds. See [Fixing SSL/TLS Certificate Errors](/aikido-device-protection/miscellaneous-aikido-endpoint/ssl-certificate-errors.md) for the tools Device Protection configures at install.
* The CA survives Aikido's automatic drift repair, so it stays trusted when a rule is reinstalled.
* It does not change the macOS System Keychain or the Windows Certificate Store. Tools that read the operating system trust store keep using whatever your MDM installed there.
* It does not register a registry domain for inspection. To have Device Protection scan packages coming from your own registry, add it under [Custom Registries](/aikido-device-protection/using-aikido-endpoint/custom-registries.md).

## Related pages

{% content-ref url="/pages/aqSzOzkyVOAcDF4MxAnb" %}
[Fixing SSL/TLS Certificate Errors](/aikido-device-protection/miscellaneous-aikido-endpoint/ssl-certificate-errors.md)
{% endcontent-ref %}

{% content-ref url="/pages/w1tOir0a3TK3dYnQtikz" %}
[The aikidoctl CLI](/aikido-device-protection/miscellaneous-aikido-endpoint/aikidoctl-cli.md)
{% endcontent-ref %}
