# Tray Icon Support on Linux

Tray icons on Linux use the [StatusNotifierItem](https://www.freedesktop.org/wiki/Specifications/StatusNotifierItem/) protocol. Most desktops implement it in their panel, but GNOME Shell needs the AppIndicator extension. Ubuntu ships and enables that extension by default, so on Debian and Red Hat-family distributions you install it yourself.

A missing tray icon is cosmetic. Protection runs in the `aikido-endpoint-protection` service and keeps blocking packages either way, but users lose the local dashboard described in [Desktop Experience](/aikido-device-protection/using-aikido-endpoint/end-user-experience.md).

| Desktop environment                                   | What you need                                                                                              |
| ----------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| GNOME on Ubuntu                                       | Nothing, the extension is enabled by default                                                               |
| GNOME on Debian                                       | Install and enable `gnome-shell-extension-appindicator`                                                    |
| GNOME on RHEL / Rocky Linux / CentOS                  | Enable EPEL, then install and enable the same package                                                      |
| KDE Plasma, Xfce 4.16 and later, Cinnamon, MATE, LXQt | Nothing, their panels handle status notifier items                                                         |
| Headless servers and containers                       | No tray. Use [the aikidoctl CLI](/aikido-device-protection/miscellaneous-aikido-endpoint/aikidoctl-cli.md) |

## Enable the extension on GNOME

{% tabs %}
{% tab title="Debian" %}
{% code overflow="wrap" %}

```bash
sudo apt install gnome-shell-extension-appindicator
gnome-extensions enable appindicatorsupport@rgcjonas.gmail.com
```

{% endcode %}
{% endtab %}

{% tab title="RHEL / Rocky Linux / CentOS" %}
The extension comes from [EPEL](https://docs.fedoraproject.org/en-US/epel/) rather than the base repositories.

{% code overflow="wrap" %}

```bash
sudo dnf install epel-release
sudo dnf install gnome-shell-extension-appindicator
gnome-extensions enable appindicatorsupport@rgcjonas.gmail.com
```

{% endcode %}
{% endtab %}
{% endtabs %}

Then restart the shell so the extension loads: press `Alt` + `F2`, type `r`, and press Enter on X11, or log out and back in on Wayland.

{% hint style="warning" %}
Run `gnome-extensions enable` as the desktop user, without `sudo`. Extensions are enabled per user, so under `sudo` it enables the extension for root and the icon still won't appear.
{% endhint %}

If the icon is still missing, check that `gnome-extensions list --enabled` includes the ID, and on Xfce that the panel has its **Status Tray** item. Otherwise [send us diagnostics](/aikido-device-protection/miscellaneous-aikido-endpoint/endpoint-information-for-support.md).
