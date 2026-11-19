# Set Up Aikido Artifacts

Configure your package manager to install through Aikido Artifacts, and to publish your private packages to it.

## Prerequisites

* Aikido Artifacts enabled for your workspace. It's in beta, so [contact Aikido](https://help.aikido.dev/miscellaneous-info/contact-support) to turn it on.

## Open the instructions

Go to **Artifacts** in the Aikido app, open the **Instructions** tab, and choose **NPM**, **Maven**, **PyPi**, or **Go**. The snippets below the selector update to match.

The tab has two sections: **Publishing Instructions** for your private packages, and **Install Instructions** for installing through the registry.

## Install packages through the registry

Add the install configuration from the **Install Instructions** section to your project or user configuration. Copy it from the app: the registry URL contains a token that is specific to your workspace.

| Package manager | Config file          | What it sets                                |
| --------------- | -------------------- | ------------------------------------------- |
| npm             | `.npmrc`             | `registry=<registry-url>`                   |
| pip             | `pip.conf`           | `index-url = <registry-url>`                |
| Maven           | `settings.xml`       | A profile with the registry as a repository |
| Go              | Environment variable | `GOPROXY=<registry-url>`                    |

{% hint style="danger" %}
The registry URL contains a token that identifies your workspace. Anyone with the URL can install through your registry. Treat it like a secret and don't commit it to a public repository.
{% endhint %}

## Publish private packages

Publishing needs a [write token](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/write-tokens). Create one first, then follow the **Publishing Instructions** for your package manager in the app. They already contain your registry URL, so you only need to replace `<write-token>` with your token.

For example, to publish a package and install it again:

{% tabs %}
{% tab title="npm" %}
In your project root, create a `.npmrc` file and publish. `<registry-host>` is the registry URL without `https://`.

{% code title=".npmrc" %}
```ini
//<registry-host>/:_authToken=<write-token>
registry=<registry-url>/
```
{% endcode %}

```shell
npm publish
```

To install it in another project, add the registry to its `.npmrc`, then run `npm install`:

{% code title=".npmrc" %}
```ini
registry=<registry-url>
```
{% endcode %}

```shell
npm install <package-name>
```
{% endtab %}

{% tab title="PyPI" %}
Create a `.pypirc` file, then build and upload with twine:

{% code title=".pypirc" %}
```ini
[distutils]
index-servers =
    aikido

[aikido]
repository = <registry-url>
username = __token__
password = <write-token>
```
{% endcode %}

```shell
python -m build
python -m twine upload --config-file .pypirc --repository aikido dist/*
```

To install it, add the registry to your `pip.conf`, then run `pip install`:

{% code title="pip.conf" %}
```ini
[global]
index-url = <registry-url>
extra-index-url = https://pypi.org/simple
```
{% endcode %}

```shell
pip install <package-name>
```
{% endtab %}
{% endtabs %}

After you publish, the version appears in [My packages](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/manage-your-packages).

## Related

{% content-ref url="write-tokens.md" %}
[write-tokens.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/write-tokens)
{% endcontent-ref %}

{% content-ref url="manage-your-packages.md" %}
[manage-your-packages.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/manage-your-packages)
{% endcontent-ref %}
