# Install Device Protection for Linux

Install the agent on a single machine, either your own or one you're testing on, with the command from the Aikido dashboard. To cover a fleet, pick a method from [Deploying on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux.md) instead.

{% hint style="info" %}
Linux devices need systemd and root access. See [Supported Ecosystems](/aikido-device-protection/miscellaneous-aikido-endpoint/supported-ecosystems.md) for the distributions and versions Device Protection supports.
{% endhint %}

## Install with one command

{% stepper %}
{% step %}
**Get the install command**

1. Go to [Device Protection and open the Devices](https://app.aikido.dev/endpoint-protection/devices) tab.
2. Click **Connect Device** in the top right.
3. Choose **Manual install**.
4. Select the [user group](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md) the device should join, then copy the **Linux install command**.
   {% endstep %}

{% step %}
**Run it on the device**

Run the command as root. It picks the right package for the distribution and architecture, installs it, and registers the device against the token in the command.
{% endstep %}
{% endstepper %}

## Choose the install flavor

The dashboard command installs the desktop flavor. A machine with no desktop session takes one extra flag, and the rest of the command stays the same:

| Machine                                                                               | Flag          |
| ------------------------------------------------------------------------------------- | ------------- |
| A laptop or workstation somebody works on                                             | none          |
| A server or VM you keep: build server, CI runner, Docker or Kubernetes host, cloud VM | `--headless`  |
| Inside a container image, or a hosted CI runner that is thrown away after the job     | `--container` |

{% code overflow="wrap" %}

```bash
curl -fsSL https://github.com/AikidoSec/safe-chain/raw/refs/heads/main/install-scripts/install-endpoint-linux.sh | sudo sh -s -- --token "$AIKIDO_TOKEN" --headless
```

{% endcode %}

`--headless` leaves out the tray icon and the desktop packages, and protection is the same. `--container` does that too, and switches the agent to a local proxy instead of the transparent datapath, which also means it needs no reboot. Container images need a few more pieces, such as the entrypoint that starts the agent, so see [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md). Hosted CI runners are VMs, but a job cannot reboot one, so they take `--container` too and start the agent explicitly. The CI guides under [Deploying on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux.md) show that. A Docker, Kubernetes, or Jenkins host that builds and runs containers is a server, so use `--headless` there.

## Manual install

Use the package directly when you want the file itself, to stage it on an internal mirror or to install on a host that can't reach the download URL. Run `uname -m` first: `x86_64` means amd64, `aarch64` means arm64.

{% tabs %}
{% tab title="Ubuntu / Debian" %}

| Architecture      | Package                                                                                                                                |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| amd64 (`x86_64`)  | [EndpointProtection-amd64.deb](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb) |
| arm64 (`aarch64`) | [EndpointProtection-arm64.deb](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-arm64.deb) |

{% code overflow="wrap" %}

```bash
curl -fLO https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
sudo AIKIDO_TOKEN=<your-token> apt install ./EndpointProtection-amd64.deb
```

{% endcode %}
{% endtab %}

{% tab title="RHEL / Rocky Linux / CentOS" %}
Use the `el9` build on version 9 and the `el10` build on version 10. Amazon Linux 2023 also takes the `el9` build.

| Version | Architecture      | Package                                                                                                                                          |
| ------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| 9       | amd64 (`x86_64`)  | [EndpointProtection-amd64.el9.rpm](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.el9.rpm)   |
| 9       | arm64 (`aarch64`) | [EndpointProtection-arm64.el9.rpm](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-arm64.el9.rpm)   |
| 10      | amd64 (`x86_64`)  | [EndpointProtection-amd64.el10.rpm](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.el10.rpm) |
| 10      | arm64 (`aarch64`) | [EndpointProtection-arm64.el10.rpm](https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-arm64.el10.rpm) |

{% code overflow="wrap" %}

```bash
curl -fLO https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.el9.rpm
sudo AIKIDO_TOKEN=<your-token> dnf install ./EndpointProtection-amd64.el9.rpm
```

{% endcode %}
{% endtab %}
{% endtabs %}

The package reads `AIKIDO_TOKEN` while it installs, then enables and starts the service itself, so there's nothing else to run.

There are no flags here, so a machine without a desktop takes its flavor from environment variables set next to the token: `AIKIDO_HEADLESS=1` for a server or VM, or `AIKIDO_CONTAINER=1` with `AIKIDO_CI_CD=1` for a container image. Add `--no-install-recommends` on the Debian family, or `--setopt=install_weak_deps=False` on the Red Hat family, to keep the desktop packages off.

{% content-ref url="/pages/7N5hgY4HXXjWX1I6GnVP" %}
[Tray Icon Support on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux/tray-icon-support-on-linux.md)
{% endcontent-ref %}

## Verify the install

{% code overflow="wrap" %}

```bash
systemctl status aikido-endpoint-protection
aikidoctl version
```

{% endcode %}

The device then appears in your [device list](https://app.aikido.dev/endpoint-protection/devices) with an **Active** status. To confirm that blocking actually works, follow [How to test if Device Protection is working](/aikido-device-protection/miscellaneous-aikido-endpoint/how-to-test-device-protection.md).

{% hint style="info" %}
Installing points Node.js and uv at the Aikido CA through system-wide environment variables. Shells that were already running won't see them, so open a new login shell or reboot before you test package installs.

`aikidoctl` gives you the same inspect, pause, and diagnostics commands as on macOS and Windows. See [The aikidoctl CLI](/aikido-device-protection/miscellaneous-aikido-endpoint/aikidoctl-cli.md).

A `--headless` or `--container` install has no **Request Access** button. When a package needs admin approval, run `aikidoctl request` instead. See [Request a Blocked Install from the Command Line](/aikido-device-protection/miscellaneous-aikido-endpoint/aikidoctl-cli/request-a-blocked-install.md).
{% endhint %}

## Upgrading

Install the newer package over the existing one with the same command. The device keeps its registration and configuration, so you don't need the token again.

## Troubleshooting

| Problem                                                                                  | Fix                                                                                                                                                                                                                                        |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| The package installs but the device never appears in the dashboard                       | `AIKIDO_TOKEN` was not set for the install command. Configuration management tools do not forward your local environment to the host, so set the variable on the task itself, then install the package again                               |
| `apt` or `dnf` looks for the package in your repositories instead of installing the file | Pass a path, not a name: `apt install ./EndpointProtection-amd64.deb`. The leading `./` is what makes the package manager treat it as a local file                                                                                         |
| The install fails on a Red Hat-family host                                               | Match the build to the major version: use the `el9` package on version 9 and the `el10` package on version 10                                                                                                                              |
| The service is not running                                                               | Run `systemctl status aikido-endpoint-protection`, then `sudo aikidoctl diagnostics` to send us the details                                                                                                                                |
| Node.js or uv still reject the Aikido certificate                                        | Installing sets `NODE_EXTRA_CA_CERTS` and `UV_SYSTEM_CERTS` system-wide, and running shells do not pick them up. Open a new login shell or reboot the device                                                                               |
| The token shows up in run output or logs                                                 | Use your tool's redaction: `no_log` in Ansible, `sensitive true` in Chef, `Sensitive()` in Puppet                                                                                                                                          |
| Several machines share one entry in the device list                                      | They booted with the same `/etc/machine-id`. Device identity on Linux follows that file, so clear it in the image you clone from and let systemd write a fresh one on first boot                                                           |
| No tray icon appears                                                                     | Expected on GNOME outside Ubuntu, and cosmetic. See [Tray Icon Support on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux/tray-icon-support-on-linux.md) |
