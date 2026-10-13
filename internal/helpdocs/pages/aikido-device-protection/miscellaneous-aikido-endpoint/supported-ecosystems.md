# Supported Operating Systems & Ecosystems

## Operating systems

We support every version that is still supported by its vendor. Once a version reaches end of life, it drops off this list.

| Operating system         | Supported versions                                                                            |
| ------------------------ | --------------------------------------------------------------------------------------------- |
| macOS                    | 15 Sequoia, 26 Tahoe                                                                          |
| Windows                  | Windows 11 version 24H2 and later. Version 23H2 on Enterprise and Education until 10 Nov 2026 |
| Ubuntu                   | 22.04 LTS, 24.04 LTS, 26.04 LTS                                                               |
| Debian                   | 12 (Bookworm), 13 (Trixie)                                                                    |
| Red Hat Enterprise Linux | 9, 10                                                                                         |
| Rocky Linux              | 9, 10                                                                                         |
| AlmaLinux                | 9, 10                                                                                         |
| Oracle Linux             | 9, 10                                                                                         |
| CentOS                   | Stream 9, Stream 10                                                                           |
| Amazon Linux             | 2023                                                                                          |

### Linux specifics

| Requirement  | Details                                                                  |
| ------------ | ------------------------------------------------------------------------ |
| Architecture | amd64 (`x86_64`) or arm64 (`aarch64`). Run `uname -m` to check           |
| Kernel       | 5.11 or later, with BTF and cgroup v2 available. Run `uname -r` to check |
| Init system  | systemd, which runs the `aikido-endpoint-protection` service             |

The kernel requirements come from the eBPF datapath the agent uses to see package installs. A container image installed with `--container` uses an explicit proxy instead, so it runs on older kernels; see [Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md).

A desktop install also wants a tray that speaks StatusNotifierItem. That one is cosmetic, and protection runs either way. See [Tray Icon Support on Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux/tray-icon-support-on-linux.md).

## Package sources

Device Protection monitors packages from the following sources:

* NPM (npm, npx, yarn, pnpm, pnpx, bun, bunx, rush, rushx, ..)
* PyPI (pip, pip3, pipx, pdm, poetry, uv, uvx, ..)
* VS Code
* Open VSX (Cursor, Windsurf, Kiro, ..)
* Maven (Gradle, sbt, ..)
* NuGet
* Chrome extensions
* Firefox extensions
* Go
* Skills.sh
* Ruby 3+ (bundler, gem)
* Rust
* PHP (composer)
* Homebrew

{% hint style="info" %}
Skills are hosted on GitHub rather than a dedicated registry, so the **Block All Installs** and **Force Requests** policies do not apply to them. See [Configuring Installation Policies](/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies.md).
{% endhint %}
