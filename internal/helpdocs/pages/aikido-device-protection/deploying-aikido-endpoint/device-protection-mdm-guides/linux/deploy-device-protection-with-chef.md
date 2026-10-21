# Deploy Device Protection with Chef

Use a Chef cookbook to install Device Protection on the Linux machines your Chef Infra Client already manages. The recipe below picks the right package per platform, installs it with your user group token, and keeps the service running on every converge.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

## Cookbook layout

```
cookbooks/aikido_device_protection/
├── attributes/default.rb
├── metadata.rb
└── recipes/default.rb
```

## Set up the cookbook

{% stepper %}
{% step %}
**Add the attributes**

{% code title="cookbooks/aikido\_device\_protection/attributes/default.rb" %}

```ruby
default['aikido']['release_url'] =
  'https://github.com/AikidoSec/safechain-internals/releases/latest/download'

# Reboot after installing. Off by default, see step 3.
default['aikido']['reboot'] = false
```

{% endcode %}
{% endstep %}

{% step %}
**Store the token in an encrypted data bag**

Create a data bag item holding the token for the user group these nodes belong to:

{% code overflow="wrap" %}

```bash
knife data bag create aikido device_protection --secret-file /path/to/encrypted_data_bag_secret
```

{% endcode %}

`data_bag_item` decrypts it on the node using the secret at `Chef::Config[:encrypted_data_bag_secret]`. If you use [Chef Vault](https://github.com/chef/chef-vault) instead, swap the lookup for `chef_vault_item('aikido', 'device_protection')['token']`.
{% endstep %}

{% step %}
**Add the recipe**

{% code title="cookbooks/aikido\_device\_protection/recipes/default.rb" %}

```ruby
arch = node['kernel']['machine'] == 'aarch64' ? 'arm64' : 'amd64'

case node['platform_family']
when 'debian'
  package_file    = "EndpointProtection-#{arch}.deb"
  install_command = 'apt-get install -y'
  installed_check = "dpkg-query -W -f='${db:Status-Status}' aikido-endpoint-protection 2>/dev/null | grep -qx installed"
when 'rhel'
  package_file    = "EndpointProtection-#{arch}.el#{node['platform_version'].to_i}.rpm"
  install_command = 'dnf install -y'
  installed_check = 'rpm -q aikido-endpoint-protection'
else
  raise "Aikido Device Protection does not support #{node['platform_family']}"
end

package_path = "/var/cache/aikido/#{package_file}"

directory '/var/cache/aikido' do
  owner 'root'
  group 'root'
  mode '0700'
end

remote_file package_path do
  source "#{node['aikido']['release_url']}/#{package_file}"
  owner 'root'
  mode '0600'
  not_if installed_check
end

execute 'install aikido-endpoint-protection' do
  command "#{install_command} #{package_path}"
  environment 'AIKIDO_TOKEN' => data_bag_item('aikido', 'device_protection')['token']
  sensitive true
  not_if installed_check
  notifies :request_reboot, 'reboot[activate aikido-endpoint-protection]', :delayed
end

service 'aikido-endpoint-protection' do
  action [:enable, :start]
end

reboot 'activate aikido-endpoint-protection' do
  action :nothing
  reason 'Activate Aikido Device Protection'
  delay_mins 5
  only_if { node['aikido']['reboot'] }
end
```

{% endcode %}

Three things in that recipe are worth knowing about:

* The install is an `execute` resource because Chef's `package` resources cannot pass the [`AIKIDO_TOKEN` the install needs](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#how-the-token-reaches-the-package).
* [`not_if installed_check`](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#make-repeat-runs-cheap) guards the download and the install, which also means it blocks upgrades. Remove it, or point `release_url` at a mirrored package, when you want the fleet to move.
* `node['aikido']['reboot']` is off by default so a converge never restarts a machine someone is working on. Turn it on for unattended nodes. See [When protection becomes active](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#when-protection-becomes-active).
  {% endstep %}

{% step %}
**Add it to the run list**

Put the recipe in the run list of the nodes you want protected, through a role, an environment, or a Policyfile:

```ruby
run_list 'recipe[aikido_device_protection]'
```

The next converge on each node installs the agent. Roll out gradually by adding the recipe to one role first.
{% endstep %}
{% endstepper %}

## Verify the rollout

On a converged node:

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
