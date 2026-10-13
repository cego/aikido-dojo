# Device Protection in WSL2

Package installs inside Windows Subsystem for Linux (WSL) run in their own Linux environment, outside the reach of the Windows agent. With **Protect WSL2** turned on, Device Protection on Windows finds the running WSL2 distros on each device and installs the Linux agent into them, so `npm`, `pip`, and other package managers in WSL get inspected the same way they do on Windows.

{% hint style="info" %}
WSL2 protection is available on request. [Contact Support](/miscellaneous-info/contact-support.md) to turn it on for your workspace.
{% endhint %}

## Prerequisites

* WSL2. WSL1 distros are skipped.
* `curl` or `wget` inside each distro, and outbound HTTPS to `github.com` to download the installer.

## Turn on WSL2 protection

1. Go to [**Device Protection** → **Settings**](https://app.aikido.dev/endpoint-protection/settings).
2. Under **Device Settings**, turn on **Protect WSL2**.

[image: Device Settings card with the Protect Containers and Protect WSL2 toggles]

The setting applies to every Windows device in your workspace. Devices pick it up on their next check-in, within 3 minutes, then wait a random delay of up to 10 minutes before the first install, so a large fleet doesn't install everywhere at once.

## How it works

The Windows agent watches the WSL2 distros of the user signed in at the device, and checks each running distro every 10 minutes. It only does this while protection is on and the agent is running normally on Windows.

* **Only running distros are touched.** Aikido never starts a stopped distro. A distro you start later is picked up within a few minutes.
* **The Linux agent is installed when it is missing or broken.** Aikido runs the standard Linux install script as root inside the distro with `--headless`, using the Windows device's token. A distro that already has a healthy package is left alone.
* **One distro at a time.** An install can take a few minutes, so each check installs into at most 1 distro and moves on to the next one on the following check.
* **The distro keeps itself up to date.** After the install, the Linux agent updates itself inside the distro, the same as on any other Linux device.
* **Each distro is its own device.** A protected distro has its own machine ID, so it reports to Aikido as a separate Linux device.

These distros are skipped:

| Distro                                                      | Why                                                                                              |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| WSL1 distros                                                | Not supported                                                                                    |
| `docker-desktop`, `rancher-desktop`, and `podman-machine-*` | Vendor images, not your workspace. To protect containers, turn on **Protect Containers** instead |
| Distros that are being installed, converted, or exported    | Aikido waits until they are in a normal state                                                    |

Turning **Protect WSL2** off stops new installs. It doesn't remove the agent from distros that already have it. To remove it, follow the Linux steps in [Uninstall Device Protection](/aikido-device-protection/deploying-aikido-endpoint/uninstall-aikido-endpoint.md) inside the distro.

## Verify it works

Inside the distro, run:

```shell
sudo aikidoctl version
```

It prints the installed version, the service status, and whether the distro is enrolled. If the command is not found, the install hasn't happened yet: check that the distro is running and wait for the next check.

## Troubleshooting

| Problem                                                     | Fix                                                                                                                                                                                                                                                                                                                                   |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| A distro is never protected                                 | Check that it is WSL2 (`wsl -l -v` in PowerShell shows `VERSION 2`), that it is running, and that it isn't on the skipped list above                                                                                                                                                                                                  |
| The install fails with `neither curl nor wget is installed` | Install `curl` in the distro, e.g. `sudo apt-get install curl`                                                                                                                                                                                                                                                                        |
| The install keeps failing                                   | Aikido retries with an increasing delay, up to 30 minutes between attempts. The Windows agent log, `%ProgramData%\AikidoSecurity\EndpointProtection\logs\EndpointProtection.log`, has the installer output for each attempt on lines starting with `WSL:`. Send that log to [Contact Support](/miscellaneous-info/contact-support.md) |
| Nothing happens while nobody is signed in                   | Expected. Aikido only acts on the distros of the user signed in at the device                                                                                                                                                                                                                                                         |

## Related

{% content-ref url="/pages/762qjKwRIkF3Q2lgx56s" %}
[Install Device Protection for Linux](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/install-device-protection-for-linux.md)
{% endcontent-ref %}

{% content-ref url="/pages/U1zOPz2xS3bI8ZXDYgin" %}
[Device Protection in Containers](/aikido-device-protection/deploying-aikido-endpoint/device-protection-mdm-guides/linux/device-protection-in-containers.md)
{% endcontent-ref %}
