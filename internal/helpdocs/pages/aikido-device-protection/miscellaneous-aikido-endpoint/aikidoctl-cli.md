# The aikidoctl CLI

`aikidoctl` is a command-line help tool that ships with Aikido Device Protection. It provides utilities to inspect, control, and troubleshoot the agent when something goes wrong.

## Usage

{% tabs %}
{% tab title="macOS" %}

```
aikidoctl [global options] command [command options]
```

{% endtab %}

{% tab title="Windows" %}

```
aikidoctl.exe [global options] command [command options]
```

{% endtab %}

{% tab title="Linux" %}

```
aikidoctl [global options] command [command options]
```

{% endtab %}
{% endtabs %}

Run `aikidoctl help` (or `aikidoctl <command> --help`) at any time to list the available commands and options. On Windows, use `aikidoctl.exe`.

## Commands

Not every command is available on every platform. The **Not available on** column lists where a command is missing. An empty cell means it works everywhere.

| Command           | Not available on | Description                                                                                                                                                                                                                                                                                                                                                                                  |
| ----------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `certs`           | Windows, Linux   | Manage CA certificates.                                                                                                                                                                                                                                                                                                                                                                      |
| `certconfig`      |                  | Manage per-tool certificate trust configuration (node, git, java, etc.).                                                                                                                                                                                                                                                                                                                     |
| `pause`           |                  | Pause Aikido Device Protection for a specified duration. **This fully disables all protection on the device while paused.**                                                                                                                                                                                                                                                                  |
| `unpause`         |                  | Resume Aikido Device Protection if it is currently paused.                                                                                                                                                                                                                                                                                                                                   |
| `sbom`            | Linux            | Inspect the software bill of materials (SBOM) collection.                                                                                                                                                                                                                                                                                                                                    |
| `secretscanning`  | Windows, Linux   | Inspect on-device secret scanning.                                                                                                                                                                                                                                                                                                                                                           |
| `device-identity` | macOS, Linux     | Show or reset the identity the device registers with. Run `device-identity reset` on a machine created from a VM clone or snapshot. See [Cloned VMs and images](#cloned-vms-and-images).                                                                                                                                                                                                     |
| `diagnostics`     |                  | <p>Collect and submit diagnostics to Aikido.<br><br>The <code>diagnostics</code> command is the command-line equivalent of the in-app <strong>Report A Problem</strong> flow. See <a href="/aikido-device-protection/miscellaneous-aikido-endpoint/endpoint-information-for-support.md">Report a Problem / Device Protection Support</a> for when and how to send diagnostics to Aikido.</p> |
| `request`         |                  | List packages blocked pending admin approval, and send that request with `aikidoctl request`. See [Request a Blocked Install from the Command Line](/aikido-device-protection/miscellaneous-aikido-endpoint/aikidoctl-cli/request-a-blocked-install.md).                                                                                                                                     |
| `restart-daemon`  | Linux            | Restart the Aikido Device Protection daemon.                                                                                                                                                                                                                                                                                                                                                 |
| `version`         |                  | Show the installed Device Protection version, build details, and service state.                                                                                                                                                                                                                                                                                                              |
| `help`, `h`       |                  | Show a list of commands or help for one command.                                                                                                                                                                                                                                                                                                                                             |

To add your organization's own CA to the certificate bundles Aikido builds, see [Adding a Custom CA Certificate](/aikido-device-protection/miscellaneous-aikido-endpoint/custom-ca-certificates.md).

## Cloned VMs and images

A device's identity is created when the agent first registers. If you install Device Protection on a VM and then clone or snapshot that VM, every machine created from it carries the same identity, so they all collapse into a single entry in your [device list](https://app.aikido.dev/endpoint-protection/devices) and you lose visibility on all but one.

On Windows, after deploying a clone or snapshot:

1. Give the machine its own hostname.
2. Reset the device identity:

{% code overflow="wrap" %}

```powershell
aikidoctl.exe device-identity reset
```

{% endcode %}

The machine then registers as its own device on the next check-in. Remove the stale shared entry from your device list if one is already there.

On Linux, device identity follows `/etc/machine-id` instead. Truncate that file as the last step of your image build so each machine generates its own ID on first boot.

## Global options

| Option            | Description                                                      |
| ----------------- | ---------------------------------------------------------------- |
| `--simple`        | Print only the installed Device Protection version, for scripts. |
| `--help`, `-h`    | Show help.                                                       |
| `--version`, `-v` | Print the version.                                               |
