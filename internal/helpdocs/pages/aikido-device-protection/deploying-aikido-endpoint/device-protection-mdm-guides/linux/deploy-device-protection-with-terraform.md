# Deploy Device Protection with Terraform

Terraform creates machines, it does not keep software converged on machines that already exist. So Terraform's job here is to make sure every Linux instance it creates comes up with Device Protection already installed.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

{% hint style="info" %}
This is not the [Aikido Terraform provider](/miscellaneous-integrations/aikido-terraform-provider.md). That provider manages Aikido workspace configuration, such as repository activation and AutoFix settings. It has nothing to do with installing the agent on a machine.
{% endhint %}

## Install on new instances

Pass a cloud-init document as `user_data`. Keep the document in its own file so it stays readable and can be reused outside Terraform. See [cloud-init](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-cloud-init.md) for the document itself and its options.

{% code title="main.tf" %}

```terraform
variable "aikido_token" {
  type      = string
  sensitive = true
}

resource "aws_instance" "workstation" {
  ami           = var.workstation_ami
  instance_type = "m7i.large"

  user_data = templatefile("${path.module}/aikido-user-data.yaml.tftpl", {
    package      = "EndpointProtection-amd64.deb"
    aikido_token = var.aikido_token
  })

  tags = {
    Name = "workstation"
  }
}
```

{% endcode %}

{% code title="aikido-user-data.yaml.tftpl" overflow="wrap" %}

```yaml
#cloud-config
runcmd:
  - curl -fsSLo /tmp/${package} https://github.com/AikidoSec/safechain-internals/releases/latest/download/${package}
  - AIKIDO_TOKEN=${aikido_token} apt-get install -y /tmp/${package}
  - rm -f /tmp/${package}
```

{% endcode %}

{% hint style="info" %}
Changing `user_data` does not reinstall anything on instances that are already running, because cloud-init's `runcmd` only fires on first boot. A change only affects instances created after it.
{% endhint %}

### Keep the token out of state and metadata

Anything Terraform renders into `user_data` is stored in state, and `user_data` is readable from the instance metadata service by anything running on the machine. Pass the name of a secret instead of the token, and let the boot script fetch the value:

{% code title="main.tf" %}

```terraform
  user_data = templatefile("${path.module}/aikido-user-data.yaml.tftpl", {
    package         = "EndpointProtection-amd64.deb"
    token_parameter = "/aikido/device-protection-token"
  })
```

{% endcode %}

Grant the instance profile read access to that parameter, and use the [boot-time lookup](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#fetch-the-token-at-boot-instead) in `aikido-user-data.yaml.tftpl`, with `${token_parameter}` and `${package}` in place of the hardcoded parameter name and package file. The `aikido_token` variable then drops out entirely, so nothing sensitive reaches Terraform state.

## Install on instances Terraform already manages

For a one-off push at create time, use a `remote-exec` provisioner:

{% code overflow="wrap" %}

```terraform
resource "aws_instance" "workstation" {
  # ...

  connection {
    type        = "ssh"
    host        = self.public_ip
    user        = "ubuntu"
    private_key = file(var.ssh_private_key)
  }

  provisioner "remote-exec" {
    inline = [
      "curl -fsSLo /tmp/EndpointProtection-amd64.deb https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection-amd64.deb",
      "sudo AIKIDO_TOKEN=${var.aikido_token} apt-get install -y /tmp/EndpointProtection-amd64.deb",
      "rm -f /tmp/EndpointProtection-amd64.deb",
    ]
  }
}
```

{% endcode %}

Provisioners come with real limits, which is why HashiCorp calls them a last resort:

* They run once, at create time. They never run again, so drift is never corrected.
* A failed provisioner marks the resource tainted, and the next apply destroys and recreates the instance.
* Terraform needs SSH reachability from wherever you run it, including from CI.

{% hint style="info" %}
For a fleet that already exists, use [Ansible](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-ansible.md), [Chef](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-chef.md), or [Puppet](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/deploy-device-protection-with-puppet.md). On AWS, SSM Run Command and State Manager also fit better than Terraform, because they are built to act on running instances.
{% endhint %}

## Verify the rollout

On a new instance:

{% code overflow="wrap" %}

```bash
cloud-init status --wait
systemctl is-active aikido-endpoint-protection
aikidoctl version
```

{% endcode %}

The device then appears in your [device list](https://app.aikido.dev/endpoint-protection/devices) with an **Active** status.

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
