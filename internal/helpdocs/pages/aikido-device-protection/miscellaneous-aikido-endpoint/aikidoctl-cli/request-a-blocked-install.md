# Request a Blocked Install from the Command Line

When a package is blocked because it needs admin approval, the desktop shows **Request Access**. On a machine without that desktop, `aikidoctl request` sends the same request. On Windows, run `aikidoctl.exe`. You do not need administrator rights.

The block label is **Needs Admin Approval**, from **Force Requests for New Packages**. See [Configuring Installation Policies](/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies.md). Admins review the request in the [Inbox](/aikido-device-protection/using-aikido-endpoint/reviewing-install-requests.md).

The list is one row per package name, newest first. Several versions of the same package are one request. A row already sent shows **Request Sent**, **Request Approved**, or **Request Rejected**, and is not sent again. Malware, **Block All Installs**, and minimum package age blocks are not listed.

{% hint style="warning" %}
Device Protection keeps this list in memory, and only the newest 500 blocks. Restarting Device Protection, or moving the device to a different permission group, clears the list. Retry the install, then run the command again.
{% endhint %}

{% tabs %}
{% tab title="Terminal" %}

```shell
aikidoctl request
```

**Space** selects or clears the highlighted package that is still waiting. **a** selects every waiting package, or clears them if they are all selected. **Enter** sends the selected packages. If none are selected, **Enter** sends the highlighted package, as long as that row is still waiting. **q** or **Esc** leaves without sending.
{% endtab %}

{% tab title="Script" %}
Without a terminal, the command prints the list and does not send anything.

```shell
aikidoctl request --yes
aikidoctl request --package lodash,chalk
aikidoctl request --event <event-id>
```

`--yes` (or `-y`) requests every package still waiting. `--package` takes one name, repeated or comma-separated. `--event` takes the **ID** column, or a prefix that matches only one row. A name or id that is not waiting is refused, and nothing is sent. `aikidoctl request lodash` is refused too. Use `--package` or `--event`.
{% endtab %}

{% tab title="JSON" %}

```shell
aikidoctl request --json
```

Prints one JSON object per line and does not send. Each line has `product`, `package`, `versions`, `status`, `blocked_at`, and `event_id`. `status` is `blocked` while the package is waiting, then `request_pending`, `request_approved`, or `request_rejected`.

Add `--yes`, `--package`, or `--event` to send, and to print the packages that were requested.
{% endtab %}
{% endtabs %}

If the command prints `No packages are waiting for approval.`, retry the install. If it prints `daemon not reachable`, Device Protection is not running.
