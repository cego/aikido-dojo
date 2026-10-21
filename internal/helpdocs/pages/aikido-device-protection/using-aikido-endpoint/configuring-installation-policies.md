# Configuring Installation Policies

You can set installation policies for each ecosystem independently. This lets you control what happens when someone on your team tries to install a new package.

## Accessing policy settings

Go to [Device Protection and open the Settings tab](https://app.aikido.dev/endpoint-protection/settings). You will see a list of all supported ecosystems with their current policy status. Click on any ecosystem to configure its settings.

## How Aikido decides to block or allow a package

When a package is installed, Aikido checks it against a set of rules in order. The first rule that matches decides the outcome.

<table data-full-width="true"><thead><tr><th width="102">Priority</th><th>Rule</th><th>Outcome</th></tr></thead><tbody><tr><td>1</td><td>Package is identified as <strong>malware</strong></td><td>Always <mark style="color:$danger;">blocked</mark>. No exceptions.</td></tr><tr><td>2</td><td>Your group has an exception that <strong>blocks</strong> this package</td><td><mark style="color:$danger;">Blocked</mark></td></tr><tr><td>3</td><td>Your group has an exception that <strong>allows</strong> this package, or an admin previously approved it</td><td><mark style="color:$success;">Allowed</mark></td></tr><tr><td>4</td><td><strong>Block All Installs</strong> is enabled for this ecosystem</td><td><mark style="color:$danger;">Blocked</mark></td></tr><tr><td>5</td><td><strong>Force Requests</strong> is enabled and the package has not been requested yet (or was previously denied)</td><td>Held for <mark style="color:orange;">approval</mark>. The person on the device requests access, and that request goes to the Inbox.</td></tr><tr><td>6</td><td>No rules match</td><td><mark style="color:$success;">Allowed</mark></td></tr></tbody></table>

{% hint style="info" %}
Group exceptions always take priority over default settings. If you're unsure why an install was blocked or allowed, check your group exceptions first.
{% endhint %}

## Available controls

{% hint style="warning" %}
It can take up to 10 minutes for changes to sync to all devices
{% endhint %}

### Block All Installs

When enabled, all installations from that ecosystem are blocked on connected devices. No new packages can be installed until you turn this off or create an exception.

Use this when you want to fully lock down an ecosystem for your organization.

### Force Requests for New Packages

When enabled, team members cannot install new packages on their own. Instead, their install attempt is sent to the **Inbox** for admin approval. The admin can then approve or reject it.

This is a good middle ground when you want oversight without completely blocking installs.

{% hint style="info" %}
**Skills are an exception.** Block All Installs and Force Requests do not apply to Skills. Skills are hosted on GitHub rather than a dedicated registry Aikido can isolate, so there is no way to block or hold skill installs without blocking or force-requesting all GitHub traffic.
{% endhint %}

### Minimum Package Age

Sets how long a package must have existed before it can be installed. For example, if set to 24 hours, any package published less than 24 hours ago will be blocked.

This helps protect against supply chain attacks, where an attacker publishes malicious code to a public registry and tries to get people to install it before the community can flag it.

Versions that fix a **critical** vulnerability are allowed through even when they are too new, provided Aikido scanned them and found no malware. See [Minimum Package Age & Version Fallback](/aikido-device-protection/using-aikido-endpoint/minimum-package-age-fallback.md#critical-security-fixes) for what happens when a version is too new.

## Exceptions

Exceptions let you override your default policies for specific packages or user groups. This is useful when your default policy is strict but certain teams need access to specific tools.

#### Adding an exception

1. Go to **Device Protection** and open the **Settings** tab.
2. Click on the ecosystem you want to configure.
3. Scroll down to the **Exceptions** section.
4. Click **Add Exception**.

In the **Add Package** dialog you set:

* **Package Name**: the package this exception applies to. Use `*` to match any substring, so `@yourorg/*` matches every package under the `@yourorg` scope.
* **Status**: whether the package is **Blocked** or **Allowed**.
* Age limit: **Standard age limit (48 hours)**, **Allow all versions (no age limit)**, or **Custom age limit** to set your own.
* **Permission Group**: the group this exception applies to.

{% hint style="warning" %}
Package identifiers vary by ecosystem. [Double-check the correct identifier before adding an exception](/aikido-device-protection/miscellaneous-aikido-endpoint/package-identifiers.md).
{% endhint %}

### Use cases

* **Removing the minimum package age**: add a new exception and set the age limit to **Allow all versions (no age limit)**
