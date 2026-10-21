# Deploy Device Protection with Ansible

Use Ansible to install Device Protection across a Linux fleet you push to over SSH. The role below picks the right package for each host, installs it with your user group token, and leaves the service enabled.

Read [Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md) first for the token, the packages, and the install options that every method uses. The [Linux Rollout Reference](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md) covers the rest: token handling, reboots, repeat runs, and device identity.

## Role layout

```
roles/aikido_device_protection/
├── defaults/main.yml
├── handlers/main.yml
└── tasks/main.yml
group_vars/
└── developer_workstations.yml
site.yml
```

## Set up the role

{% stepper %}
{% step %}
**Add the defaults**

{% code title="roles/aikido\_device\_protection/defaults/main.yml" %}

```yaml
---
aikido_release_url: https://github.com/AikidoSec/safechain-internals/releases/latest/download

# Reboot after installing. Off by default, see step 3.
aikido_reboot: false
```

{% endcode %}
{% endstep %}

{% step %}
**Add the tasks**

The first task decides whether there is anything to do, so repeat runs are cheap and never touch the token.

{% code title="roles/aikido\_device\_protection/tasks/main.yml" %}

```yaml
---
- name: Check whether Device Protection is already installed
  ansible.builtin.command: aikidoctl --simple
  register: aikido_installed
  changed_when: false
  failed_when: false

- name: Install Device Protection
  when: aikido_installed.rc != 0
  block:
    - name: Work out the package architecture
      ansible.builtin.set_fact:
        aikido_arch: "{{ 'arm64' if ansible_facts['architecture'] == 'aarch64' else 'amd64' }}"

    - name: Work out the package name
      ansible.builtin.set_fact:
        aikido_package: >-
          {{
            'EndpointProtection-' ~ aikido_arch ~ '.deb'
            if ansible_facts['os_family'] == 'Debian'
            else 'EndpointProtection-' ~ aikido_arch ~ '.el'
                 ~ ansible_facts['distribution_major_version'] ~ '.rpm'
          }}

    - name: Download the package
      ansible.builtin.get_url:
        url: "{{ aikido_release_url }}/{{ aikido_package }}"
        dest: "/tmp/{{ aikido_package }}"
        mode: "0600"

    - name: Install the package on Debian and Ubuntu
      ansible.builtin.apt:
        deb: "/tmp/{{ aikido_package }}"
      environment:
        AIKIDO_TOKEN: "{{ aikido_token }}"
      no_log: true
      when: ansible_facts['os_family'] == 'Debian'
      notify: Reboot to activate Device Protection

    - name: Install the package on RHEL, Rocky Linux and CentOS
      ansible.builtin.dnf:
        name: "/tmp/{{ aikido_package }}"
        state: present
      environment:
        AIKIDO_TOKEN: "{{ aikido_token }}"
      no_log: true
      when: ansible_facts['os_family'] == 'RedHat'
      notify: Reboot to activate Device Protection

    - name: Remove the downloaded package
      ansible.builtin.file:
        path: "/tmp/{{ aikido_package }}"
        state: absent

- name: Make sure the service is enabled and running
  ansible.builtin.systemd_service:
    name: aikido-endpoint-protection
    state: started
    enabled: true
```

{% endcode %}

Leave `gather_facts` on, the role needs it. `distribution_major_version` is what selects the `el9` or `el10` build. Amazon Linux 2023 reports `2023`, so map it to `el9` if you run it.

The `environment` keyword is what [hands the token to the install](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#how-the-token-reaches-the-package), and `no_log: true` keeps it out of the play output.
{% endstep %}

{% step %}
**Add the reboot handler**

{% code title="roles/aikido\_device\_protection/handlers/main.yml" %}

```yaml
---
- name: Reboot to activate Device Protection
  ansible.builtin.reboot:
    reboot_timeout: 600
    msg: Rebooting to activate Aikido Device Protection
  when: aikido_reboot | bool
```

{% endcode %}

Set `aikido_reboot: true` for hosts nobody is working on. On developer machines, leave it off and ask people to log out and back in. See [When protection becomes active](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/linux-rollout-reference.md#when-protection-becomes-active).
{% endstep %}

{% step %}
**Store the token in Ansible Vault**

Encrypt the token for the inventory group that maps to your Aikido user group:

{% code overflow="wrap" %}

```bash
ansible-vault encrypt_string --name aikido_token '<your-token>'
```

{% endcode %}

Paste the output into the group's variables:

{% code title="group\_vars/developer\_workstations.yml" %}

```yaml
---
aikido_token: !vault |
  $ANSIBLE_VAULT;1.1;AES256
  62313365396662343061393464336163383764373764613633653634306231386433626436623361
  ...
```

{% endcode %}

Repeat per group if you use more than one user group. Each group gets its own token and its own `group_vars` file.

{% hint style="info" %}
If the token already lives in your CI secrets, read it from the controller's environment instead of committing a vaulted copy: `aikido_token: "{{ lookup('env', 'AIKIDO_TOKEN') }}"`.
{% endhint %}
{% endstep %}

{% step %}
**Write the playbook and run it**

{% code title="site.yml" %}

```yaml
---
- name: Deploy Aikido Device Protection
  hosts: developer_workstations
  become: true
  serial: 20%
  roles:
    - aikido_device_protection
```

{% endcode %}

Start with a few machines, then widen:

{% code overflow="wrap" %}

```bash
ansible-playbook site.yml --ask-vault-pass --limit pilot
ansible-playbook site.yml --ask-vault-pass
```

{% endcode %}
{% endstep %}
{% endstepper %}

## Verify the rollout

{% code overflow="wrap" %}

```bash
ansible developer_workstations -m command -a "aikidoctl --simple"
ansible developer_workstations -b -m command -a "systemctl is-active aikido-endpoint-protection"
```

{% endcode %}

The devices then appear in your [device list](https://app.aikido.dev/endpoint-protection/devices) with an **Active** status.

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
