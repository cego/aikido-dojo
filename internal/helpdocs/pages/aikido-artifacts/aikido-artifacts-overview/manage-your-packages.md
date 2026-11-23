# Manage Your Packages

The **My packages** tab lists every private package version you've published to Aikido Artifacts. Use it to check what's published and to remove versions.

## View your packages

Go to **Artifacts** → **My packages**. The tab lists only the private packages you published, not cached public packages.

| Column              | Description                                          |
| ------------------- | ---------------------------------------------------- |
| **Ecosystem**       | The package ecosystem, e.g. npm, pypi, maven, or go. |
| **Package name**    | The name of the package.                             |
| **Package version** | The published version. Each version is its own row.  |
| **Upload date**     | When the version was published.                      |

Use the search field to filter by package name or ecosystem. Use the **Ecosystem** dropdown to show one ecosystem.

## Add a package

Click **Add Package**. The dialog shows the publishing steps for each package manager. Follow them as described in [Set Up Aikido Artifacts](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/set-up-aikido-artifacts#publish-private-packages).

## Delete a package version

Only workspace admins can delete package versions. In the row of the version you want to remove, click the three-dot menu, click **Delete version**, then click **Delete** to confirm.

{% hint style="warning" %}
A deleted version is no longer resolvable from the registry, and any build that depends on it fails. Deleting removes one version, not the whole package. You can't publish a version number again after you delete it.
{% endhint %}

## Related

{% content-ref url="set-up-aikido-artifacts.md" %}
[set-up-aikido-artifacts.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/set-up-aikido-artifacts)
{% endcontent-ref %}

{% content-ref url="write-tokens.md" %}
[write-tokens.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/write-tokens)
{% endcontent-ref %}
