# Aikido Artifacts Overview

Aikido Artifacts is a single package registry for your private packages and your public dependencies. You publish internal packages to it and install everything through it, while Aikido screens public packages for malware and holds back brand-new releases.

{% hint style="info" %}
Aikido Artifacts is in **closed beta** and supports npm, PyPI, Maven, and Go. To enable it for your workspace, [contact Aikido](https://help.aikido.dev/miscellaneous-info/contact-support).
{% endhint %}

## Why use Aikido Artifacts

* **Malware screening.** Public packages are checked against Aikido Intel before they reach a laptop or CI runner.
* **Minimum package age.** By default, releases newer than 48 hours aren't served, which gives malware detection time to catch up with fresh publishes. You can change this for your workspace.
* **Patched packages.** Existing [Aikido Libraries](https://help.aikido.dev/autofix-and-remediation/aikido-libraries-overview) patches are served from the same registry.
* **One registry for everything.** Host internal packages and install public dependencies from the same URL.
* **Nothing to host.** Aikido runs the registry. You change a registry URL and keep working.

## How it works

Point your package manager at Aikido Artifacts. For every install, the registry looks for the package in this order:

1. **Your private registry.** Aikido looks for the package among the packages you published to your own registry. If it's there, Aikido returns it.
2. **The Aikido package registry.** If the package isn't in your private registry, Aikido looks in its own package registry. It holds all Aikido Libraries packages and cached copies of upstream packages. Every package in it has passed the malware and minimum age checks.
3. **The upstream registry.** If the package isn't in the Aikido package registry yet, Aikido fetches it from the upstream registry (e.g. npmjs.com or PyPI). The same malware and minimum age checks apply before Aikido returns it and caches it for the next install.

{% hint style="warning" %}
Because private packages are looked up first, a private package takes priority over a public package with the same name.
{% endhint %}

## Malware and minimum age checks

Public packages (steps 2 and 3) are screened before Aikido returns them:

* **Malware.** Aikido doesn't serve a version that Aikido Intel flags as malware.
* **Minimum package age.** Aikido doesn't serve a version published less than 48 hours ago (to change this for your workspace, see [Configuring Installation Policies](https://help.aikido.dev/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies#minimum-package-age)).

Publishing a private package to Aikido Artifacts doesn't trigger a scan. If the package's source code is in a repository you've connected to Aikido, Aikido already scans it as part of your code.

## Get started

{% content-ref url="set-up-aikido-artifacts.md" %}
[set-up-aikido-artifacts.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/set-up-aikido-artifacts)
{% endcontent-ref %}

{% content-ref url="manage-your-packages.md" %}
[manage-your-packages.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/manage-your-packages)
{% endcontent-ref %}

{% content-ref url="write-tokens.md" %}
[write-tokens.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/write-tokens)
{% endcontent-ref %}

## Related

{% content-ref url="../autofix-and-remediation/aikido-libraries-overview/" %}
[aikido-libraries-overview](https://help.aikido.dev/autofix-and-remediation/aikido-libraries-overview)
{% endcontent-ref %}
