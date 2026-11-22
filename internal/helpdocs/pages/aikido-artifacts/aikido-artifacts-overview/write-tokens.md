# Write Tokens

{% hint style="info" %}
Only workspace admins can create, view, and delete write tokens.
{% endhint %}

A write token lets a package manager or CI job publish private packages to Aikido Artifacts. Create one token per publisher, such as one per CI pipeline, so you can revoke each one on its own.

## Create a write token

1. Go to **Artifacts** → **Write Tokens**.
2. Click **Add Token**, enter a name, and click **Add Token**. The name must be unique in your workspace.
3. Click **Copy** and store the token in a secure place, such as your CI secret store.

{% hint style="warning" %}
Aikido shows the token once. You can't view it again after you close the dialog. If you lose it, create a new token.
{% endhint %}

## Use a write token

Use the token as `<write-token>` in the publishing steps in [Set Up Aikido Artifacts](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/set-up-aikido-artifacts#publish-private-packages).

## Delete a write token

In the token's row, click the three-dot menu, click **Delete token**, then click **Delete** to confirm.

{% hint style="danger" %}
Deleting a token revokes access for every package manager or CI job that uses it. Those publishes fail until you give them a new token.
{% endhint %}

## Token permissions

A write token can publish and read private packages in your workspace. It doesn't expire, and you can't limit it to an ecosystem or a package. Delete it when you no longer need it.

## Related

{% content-ref url="set-up-aikido-artifacts.md" %}
[set-up-aikido-artifacts.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/set-up-aikido-artifacts)
{% endcontent-ref %}

{% content-ref url="manage-your-packages.md" %}
[manage-your-packages.md](https://help.aikido.dev/aikido-artifacts/aikido-artifacts-overview/manage-your-packages)
{% endcontent-ref %}
