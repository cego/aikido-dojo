# Linux Rollout Reference

Every method under [Deploying on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux.md) does the same thing: fetch the package for the host, install it with your user group token in the environment, and get users onto a fresh login shell afterwards. This page explains those shared steps once. Each guide shows how its own tool does them and links here for the why.

## How the token reaches the package

The `.deb` and `.rpm` packages read `AIKIDO_TOKEN` while they install, and register the device against it. The variable has to be in the environment of the process that runs `apt`, `dnf`, or `dpkg`, which is why every guide sets it on the install step itself rather than anywhere else.

Setting `AIKIDO_TOKEN` in your own shell first does nothing. Ansible, Chef, and Puppet do not forward the environment of the machine you run them from, and a Docker build cannot see it either.

One token covers one Aikido [user group](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md), because the token is what selects the group a device joins. Machines in different groups need a token each, kept in a separate inventory group, Hiera layer, data bag item, or module instance.

## Keep the token out of logs and state

The token is only needed on first install, but most tools record what they ran. Every guide uses its tool's redaction, so what is left to check is where the value still lands:

| Tool       | Where the token can still leak                                                              |
| ---------- | ------------------------------------------------------------------------------------------- |
| Ansible    | Plain-text group vars. Keep it in Ansible Vault                                             |
| Chef       | Plain-text data bags. Use an encrypted data bag or Chef Vault                               |
| Puppet     | The install script on disk, and the catalog cached under `/opt/puppetlabs/puppet/cache`     |
| Terraform  | State, for anything rendered into `user_data`. `TF_LOG=DEBUG` output is not redacted either |
| Packer     | A `.pkrvars.hcl` file. Pass the token in from your secret store instead                     |
| cloud-init | The instance metadata service, and `/var/log/cloud-init-output.log`                         |
| Docker     | The image and its build history, when the token is passed with `--build-arg`                |

{% hint style="warning" %}
Redaction also hides the errors you need when an install fails. Turn it off on a single test host rather than on the whole fleet.
{% endhint %}

### Fetch the token at boot instead

On cloud instances, the cleanest fix for user data and state is to not put the token there. Give the machine an identity that can read your secret store and let the install fetch it:

{% code title="user-data" overflow="wrap" %}

```yaml
#cloud-config
runcmd:
  - curl -fsSLo /tmp/aikido.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb
  - |
    AIKIDO_TOKEN=$(aws ssm get-parameter \
      --name /aikido/device-protection-token \
      --with-decryption \
      --query Parameter.Value \
      --output text)
    export AIKIDO_TOKEN
    apt-get install -y /tmp/aikido.deb
  - rm -f /tmp/aikido.deb
```

{% endcode %}

Swap the lookup for your own cloud: `az keyvault secret show` with a managed identity, or `gcloud secrets versions access` with a service account.

## When protection becomes active

Installing points Node.js and uv at the Aikido CA through system-wide environment variables (`NODE_EXTRA_CA_CERTS` and `UV_SYSTEM_CERTS`). Shells and services that were already running never see them, so protection is only complete after a new login shell or a reboot.

Don't restart machines people are working on from an automation run. Ask those users to log out and back in, or reboot in a maintenance window. For unattended machines, reboot at the end of the install, which each guide shows for its own tool.

## Make repeat runs cheap

Ansible, Chef, and Puppet run again on every play or converge, so guard the download and the install with a check for the agent. A guarded run does nothing on a machine that already has it, and never touches the token:

| Platform                  | Check                                                                                     |
| ------------------------- | ----------------------------------------------------------------------------------------- |
| Debian and Ubuntu         | `dpkg-query -W -f='${db:Status-Status}' aikido-endpoint-protection \| grep -qx installed` |
| RHEL, Rocky Linux, CentOS | `rpm -q aikido-endpoint-protection`                                                       |
| Either                    | `aikidoctl --simple`                                                                      |

Use `dpkg-query` and not `dpkg -s` on the Debian family. `dpkg -s` also succeeds for a package that was removed but not purged, so a machine where someone ran `apt remove aikido-endpoint-protection` looks installed to your tool and never gets converged back.

## How devices are identified

Device identity on Linux follows `/etc/machine-id`. Machines that boot with the same file share a single entry in your device list, so you lose visibility on all but one of them. Truncate the file as the last step of any image build or VM template, and systemd writes a fresh ID on each machine's first boot.

Container images are the exception: every container from one image reports as the same device, which is usually what you want in CI.

## Versions and upgrades

The `releases/latest/download` URLs always resolve to the newest release. To move the fleet on your own schedule instead, mirror the package internally and point your tool at the mirror.

Upgrades keep the device's registration and need no token, so install the newer package over the existing one. If your automation skips the install once the agent is present, make that check version-aware or temporarily remove it for the upgrade run. See [Upgrading](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md#upgrading).

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
