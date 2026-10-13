# Deploy Device Protection with Puppet

Use a Puppet class to install Device Protection on the Linux nodes your Puppet agent already manages. The class below picks the right package per platform, installs it with your user group token, and keeps the service running on every run.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

## Module layout

```
modules/aikido_device_protection/
└── manifests/init.pp
data/
└── groups/developers.yaml
```

## Set up the class

{% stepper %}
{% step %}
**Store the token in Hiera eyaml**

Encrypt the token for the group of nodes that maps to your Aikido user group:

{% code overflow="wrap" %}

```bash
eyaml encrypt -l 'aikido_device_protection::token' -s '<your-token>'
```

{% endcode %}

Paste the output into the matching Hiera layer:

{% code title="data/groups/developers.yaml" %}

```yaml
---
aikido_device_protection::token: >
  ENC[PKCS7,MIIBiQYJKoZIhvcNAQcDoIIBejCCAXYCAQAxggEhMIIBHQIBADAFMAACAQEw...]
```

{% endcode %}

The class parameter is typed `Sensitive`, and Puppet wraps values coming from Hiera automatically, so nothing else is needed.
{% endstep %}

{% step %}
**Add the class**

{% code title="modules/aikido\_device\_protection/manifests/init.pp" %}

```puppet
class aikido_device_protection (
  Sensitive[String[1]] $token,
  String[1] $release_url = 'https://github.com/AikidoSec/safechain-internals/releases/latest/download',
) {
  $arch = $facts['os']['architecture'] ? {
    /^(aarch64|arm64)$/ => 'arm64',
    default             => 'amd64',
  }

  case $facts['os']['family'] {
    'Debian': {
      $package_file    = "EndpointProtection-${arch}.deb"
      $install_command = 'apt-get install -y'
      $installed_check = "dpkg-query -W -f='\${db:Status-Status}' aikido-endpoint-protection 2>/dev/null | grep -qx installed"
    }
    'RedHat': {
      $package_file    = "EndpointProtection-${arch}.el${facts['os']['release']['major']}.rpm"
      $install_command = 'dnf install -y'
      $installed_check = 'rpm -q aikido-endpoint-protection'
    }
    default: {
      fail("Aikido Device Protection does not support ${facts['os']['family']}")
    }
  }

  $package_path = "/var/cache/aikido/${package_file}"

  file { '/var/cache/aikido':
    ensure => directory,
    owner  => 'root',
    group  => 'root',
    mode   => '0700',
  }

  exec { 'download aikido-endpoint-protection':
    command  => "curl -fsSLo ${package_path} ${release_url}/${package_file}",
    provider => shell,
    path     => ['/usr/bin', '/bin'],
    creates  => $package_path,
    unless   => $installed_check,
    require  => File['/var/cache/aikido'],
  }

  $install_script = @("SCRIPT")
    #!/bin/sh
    set -e
    export AIKIDO_TOKEN='${$token.unwrap}'
    exec ${install_command} ${package_path}
    | SCRIPT

  file { '/usr/local/sbin/aikido-install':
    ensure  => file,
    owner   => 'root',
    group   => 'root',
    mode    => '0700',
    content => Sensitive($install_script),
  }

  exec { 'install aikido-endpoint-protection':
    command  => '/usr/local/sbin/aikido-install',
    provider => shell,
    path     => ['/usr/bin', '/bin', '/usr/sbin', '/sbin'],
    unless   => $installed_check,
    require  => [
      Exec['download aikido-endpoint-protection'],
      File['/usr/local/sbin/aikido-install'],
    ],
  }

  service { 'aikido-endpoint-protection':
    ensure  => running,
    enable  => true,
    require => Exec['install aikido-endpoint-protection'],
  }
}
```

{% endcode %}

Four things in that class are worth knowing about:

* The install runs from a root-only script wrapped in `Sensitive()`, because the package needs [`AIKIDO_TOKEN` in its environment](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#how-the-token-reaches-the-package) and Puppet does not redact an `exec`'s `environment` attribute.
* [`unless => $installed_check`](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#make-repeat-runs-cheap) guards the download and the install, so nodes that already run the agent do nothing.
* `provider => shell` is there because the Debian check is a pipeline, and `exec` does not use a shell by default.
* The download uses `curl`, which minimal images sometimes leave out. If yours does, install it alongside this class.

{% hint style="warning" %}
The token lands on disk in two root-only places: `/usr/local/sbin/aikido-install`, which has no guard and stays in place after the install, and the cached catalog under `/opt/puppetlabs/puppet/cache`, which the agent needs to apply the resource. Treat both as sensitive and keep report access restricted.

Once a group of nodes is installed, switching the `file` resource to `ensure => absent` clears the script from all of them. Put it back before the next new node joins that group.
{% endhint %}
{% endstep %}

{% step %}
**Classify your nodes**

Declare the class where you classify nodes, through roles and profiles, the Puppet Enterprise console, or `site.pp`:

```puppet
include aikido_device_protection
```

The next agent run installs the agent. Classify one group first, then widen.

Puppet has no built-in reboot resource, so protection becomes active once users log out and back in. To reboot unattended nodes instead, notify [puppetlabs-reboot](https://forge.puppet.com/modules/puppetlabs/reboot) from the install `exec`. See [When protection becomes active](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#when-protection-becomes-active).
{% endstep %}
{% endstepper %}

## Verify the rollout

On a node that has run:

{% code overflow="wrap" %}

```bash
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
