# Deploy Device Protection with cloud-init

cloud-init runs on the first boot of a new machine, which makes it the simplest way to have Device Protection installed before anyone logs in. It works anywhere cloud-init does: AWS, Azure, GCP, OpenStack, Proxmox, and plain libvirt images.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

## Set up the user data

The agent is not published in a distribution repository, so the `packages:` module cannot install it. Download the package and install it from `runcmd`.

{% tabs %}
{% tab title="Ubuntu / Debian" %}
{% code title="user-data" overflow="wrap" %}

```yaml
#cloud-config
runcmd:
  - curl -fsSLo /tmp/aikido.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
  - AIKIDO_TOKEN=<your-token> apt-get install -y /tmp/aikido.deb
  - rm -f /tmp/aikido.deb
```

{% endcode %}
{% endtab %}

{% tab title="RHEL / Rocky Linux / CentOS" %}
{% code title="user-data" overflow="wrap" %}

```yaml
#cloud-config
runcmd:
  - curl -fsSLo /tmp/aikido.rpm https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.el9.rpm
  - AIKIDO_TOKEN=<your-token> dnf install -y /tmp/aikido.rpm
  - rm -f /tmp/aikido.rpm
```

{% endcode %}

Use the `el10` package on version 10.
{% endtab %}
{% endtabs %}

The package enables and starts the `aikido-endpoint-protection` service itself, so there is nothing to add for that.

One document covers one image family, one architecture, and on Red Hat-family images one major version, because the package name differs for each. Keep a document per image rather than adding detection logic to `runcmd`.

{% hint style="warning" %}
A literal token in `runcmd` is not private. User data is readable from the instance metadata service by anything running on the machine, and cloud-init writes the commands it runs to `/var/log/cloud-init-output.log`. For anything past a quick test, [fetch the token from your secret store at boot](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#fetch-the-token-at-boot-instead) instead.
{% endhint %}

{% hint style="info" %}
`runcmd` only fires on the first boot, so changing the user data does nothing to machines that are already running. Treat it as image configuration rather than a way to manage a fleet. For machines that already exist, use [Ansible](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-ansible.md), [Chef](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-chef.md), or [Puppet](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-puppet.md).
{% endhint %}

## Reboot at the end of the first boot

`runcmd` runs late in the boot, which matters for the Aikido CA environment variables the install sets system-wide. Logins after boot pick them up. Services that started earlier in the same boot, such as a container runtime or a build agent, do not. If the machine runs those, add `power_state` so it reboots once the install is done:

{% code title="user-data" overflow="wrap" %}

```yaml
#cloud-config
runcmd:
  - curl -fsSLo /tmp/aikido.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
  - AIKIDO_TOKEN=<your-token> apt-get install -y /tmp/aikido.deb
  - rm -f /tmp/aikido.deb

power_state:
  mode: reboot
  message: Rebooting to activate Aikido Device Protection
  condition: true
```

{% endcode %}

## Verify the rollout

On a new machine:

{% code overflow="wrap" %}

```bash
cloud-init status --wait
systemctl is-active aikido-endpoint-protection
aikidoctl version
```

{% endcode %}

If the install failed, `/var/log/cloud-init-output.log` has the package manager output. The device appears in your [device list](https://app.aikido.dev/endpoint-protection/devices) with an **Active** status.

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
