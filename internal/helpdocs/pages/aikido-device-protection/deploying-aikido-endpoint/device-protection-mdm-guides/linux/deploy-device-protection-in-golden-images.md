# Deploy Device Protection in Golden Images

If you build your own Linux images with Packer, or any other image pipeline, install Device Protection during the build. Every machine created from the image is then protected from its first boot, with no install step at runtime.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

## Install during the build

{% code title="workstation.pkr.hcl" overflow="wrap" %}

```hcl
variable "aikido_token" {
  type      = string
  sensitive = true
}

build {
  sources = ["source.amazon-ebs.workstation"]

  provisioner "shell" {
    environment_vars = [
      "AIKIDO_TOKEN=${var.aikido_token}",
    ]
    inline = [
      "curl -fsSLo /tmp/aikido.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb",
      "sudo AIKIDO_TOKEN=\"$AIKIDO_TOKEN\" apt-get install -y /tmp/aikido.deb",
      "rm -f /tmp/aikido.deb",
      "systemctl is-enabled aikido-endpoint-protection",
      "sudo truncate -s 0 /etc/machine-id",
    ]
  }
}
```

{% endcode %}

Nothing has to be deferred to first boot: the package registers the device during the build. Three things in that build are worth knowing about:

* Pass the token in from your secret store rather than a `.pkrvars.hcl` file, for example `PKR_VAR_aikido_token=$(...) packer build .`. Marking the variable `sensitive` keeps it out of Packer's output. See [Keep the token out of logs and state](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#keep-the-token-out-of-logs-and-state).
* `systemctl is-enabled` fails the build if the service did not get enabled, which is the check you want. A broken image is much more expensive to find later.
* `truncate -s 0 /etc/machine-id` comes last, so systemd writes a fresh ID on each machine's first boot. This one is not optional: machines that share a machine ID share a single entry in your device list. See [How devices are identified](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#how-devices-are-identified).

{% hint style="info" %}
The instance Packer builds on registers itself during the install, and it disappears when the build finishes. Delete that entry so it does not sit in your device list as an inactive device. See [Removing a device](/aikido-device-protection/deploying-aikido-endpoint/connecting-and-managing-devices.md#removing-a-device).
{% endhint %}

## One image, one user group

The token selects the Aikido [user group](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md) the device joins, so a baked image belongs to one user group. If one image has to serve several groups, keep the token out of it and stage only the package file:

{% code title="workstation.pkr.hcl" overflow="wrap" %}

```hcl
provisioner "shell" {
  inline = [
    "sudo mkdir -p /var/cache/aikido",
    "sudo curl -fsSLo /var/cache/aikido/EndpointProtection-amd64.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb",
  ]
}
```

{% endcode %}

Then install from the staged file on first boot, with the token of the group that machine belongs to:

{% code title="user-data" %}

```yaml
#cloud-config
runcmd:
  - AIKIDO_TOKEN=<your-token> apt-get install -y /var/cache/aikido/EndpointProtection-amd64.deb
```

{% endcode %}

First boot then needs no network access to GitHub. See [cloud-init](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-cloud-init.md) for the rest of that document.

## Verify the rollout

Boot one machine from the finished image and check:

{% code overflow="wrap" %}

```bash
cat /etc/machine-id
systemctl is-active aikido-endpoint-protection
aikidoctl version
```

{% endcode %}

Boot a second machine from the same image and confirm the two machine IDs differ, and that both appear in your [device list](https://app.aikido.dev/endpoint-protection/devices) with an **Active** status.

Images pin whatever version was current at build time, so rebuild on your normal image cadence and upgrade machines already running from an older image by [installing the newer package over it](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#versions-and-upgrades).

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
