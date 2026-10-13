# Connecting and Managing Devices

A device is any computer connected to Device Protection through an installed agent. Once connected, Aikido automatically detects and monitors the packages installed on that device.

## Connecting a device

1. Go to [Device Protection and open the Devices](https://app.aikido.dev/endpoint-protection/devices) tab.
2. Click **Connect Device** in the top right.
3. Follow the setup instructions to install the agent on the device.

Once the agent is installed, the [device will appear in your device list](https://app.aikido.dev/endpoint-protection/devices) with an Active status.

## Device list

The device list shows every connected device. Each row shows:

* OS: Type of operating system
* Name: The device name and owner (e.g. "Kate Davis' Windows PC")
* User: The name of the account on the device
* User Group: Which group the device belongs to (e.g. Marketing, Product, Sales)
* Last seen: How recently the device checked in
* Agent Status: Whether the agent is currently active or offline and if Device Protection is up to date

### Device status

Aikido sets a device's status based on how long it has been since the agent's last heartbeat. The agent sends a heartbeat every 10 minutes.

| Status                             | Rule                                                                                                                         |
| ---------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Archived                           | No heartbeat for 30 days or more.                                                                                            |
| Offline                            | No heartbeat for more than 60 minutes.                                                                                       |
| Not protected / Update / Protected | Heartbeat within the last 60 minutes. Which one shows depends on whether setup is complete and the agent version is current. |

## User groups

{% content-ref url="/pages/l6Jr0yFrUmsVezoR1yoV" %}
[User Groups for Device Protection](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md)
{% endcontent-ref %}

## Device actions

Click the `⋮` menu next to a device to view device actions.

### Pausing protection on a device

Pause Protection temporarily turns off Device Protection on a specific device, without uninstalling the agent. The device stays connected, and you turn protection back on from the same menu when you're done.

Use it as an exception, not a routine step. For anything that needs to happen regularly, adjust your [installation policies](/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies.md) or [user groups](/aikido-device-protection/deploying-aikido-endpoint/user-groups-for-aikido-endpoint.md) instead.

Reasons to pause protection:

* Troubleshooting an issue that might be caused by Device Protection itself
* Urgently installing software that's blocked by a policy, while you set up a longer-term rule
* IT or support walking a user through a fix that requires protection to be off temporarily
* Checking whether Device Protection is conflicting with another security or monitoring tool on the device

{% hint style="warning" %}
While paused, the device isn't protected against unauthorized or malicious package installs. Turn protection back on as soon as you're done.
{% endhint %}

### Removing a device

Remove Device disconnects the device from Device Protection and takes it off your device list. This does not uninstall the agent from the device itself.

Remove a device when:

* It's been decommissioned, wiped, or replaced
* An employee has offboarded and the device no longer needs to be monitored
* It was connected by mistake or for testing and doesn't need to be tracked

{% hint style="warning" %}
Removing a device only removes it from Aikido. To fully stop Device Protection from running, [uninstall the agent](/aikido-device-protection/deploying-aikido-endpoint/uninstall-aikido-endpoint.md) from the device.
{% endhint %}
