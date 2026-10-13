# Fixing SSL/TLS Certificate Errors

Device Protection proxies your traffic through a local CA certificate. The following tools are configured automatically at install:

* Node.js, npm
* Python (pip, Poetry, uv)
* Ruby (Bundler, RubyGems)
* Java (all JDKs discovered at install time)
* git, curl
* Homebrew, Bazel, Nix (macOS)
* JetBrains IDEs, Android Studio, DBeaver
* Firefox
* macOS System Keychain / Windows Certificate Store

{% hint style="info" %}
Don't see your tool? [Contact support](mailto:support@aikido.dev), we'll look into adding automatic support for it.
{% endhint %}

If you need to fix it yourself in the meantime, use the information below.

## Unsupported tools

The CA certificate is stored at:

{% tabs %}
{% tab title="macOS" %}

```
/Library/Application Support/AikidoSecurity/EndpointProtection/run/endpoint-protection-proxy-ca-crt.pem
```

{% endtab %}

{% tab title="Windows" %}

```
C:\ProgramData\AikidoSecurity\EndpointProtection\run\endpoint-protection-proxy-ca-crt.pem
```

{% endtab %}
{% endtabs %}

You can also copy your workspace CA certificate in PEM format from the [Certificates](https://app.aikido.dev/endpoint-protection/certificates) page in Aikido, under **Root CA Certificate**.

Follow your tool's documentation for adding a custom CA certificate.

## Your organization's own CA

If the failing certificate is your organization's, not Aikido's, for example a corporate proxy that inspects TLS traffic or a private registry with an internal certificate, add it to the bundles Aikido builds instead. See [Adding a Custom CA Certificate](/aikido-device-protection/miscellaneous-aikido-endpoint/custom-ca-certificates.md).
