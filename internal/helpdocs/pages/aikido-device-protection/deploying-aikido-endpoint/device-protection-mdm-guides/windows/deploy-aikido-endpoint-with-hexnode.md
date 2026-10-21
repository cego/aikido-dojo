# Deploy Device Protection with Hexnode

Use Hexnode UEM to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Hexnode Windows management before you begin.
{% endhint %}

## Installation

{% stepper %}
{% step %}
**What you'll need**

Before starting, make sure you have the following from the Aikido Device Protection dashboard:

* Your **PowerShell install command** (copied from the user group selector on the dashboard)

If you're missing it, go back to the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device**.
{% endstep %}

{% step %}
**Run the install command on managed devices**

Push the PowerShell command from the Aikido dashboard to your Windows fleet through Hexnode's Execute Custom Script action.

1. In your Hexnode portal, go to **Manage** → **Devices** and select the target Windows devices, or open a device group.
2. Open **Actions** → **Others** → **Execute Custom Script**.
3. Select **Windows** as the platform.
4. Choose **PowerShell script** as the script type.
5. Paste the install command copied from the Aikido Device Protection dashboard into the script body.
6. Set the timeout to at least **30 minutes** to allow for the installer download.
7. Click **Execute** to push the command. Hexnode runs custom scripts on Windows in the device (SYSTEM) context by default, which is what the Aikido installer requires.

You can monitor the execution status of each device under **Manage** → **Action History**.
{% endstep %}

{% step %}
**Apply the lockdown policy**

Use a Custom Configuration policy to block uninstallation of the Aikido Device Protection app.

1. In your Hexnode portal, go to **Policies** and create a new blank policy for Windows.
2. Under the **Windows** tab, go to **Configurations** and click **Configure** next to **Deploy Custom Configuration**.
3. Click **Add Payload** and enter the following:
   * **Name**: `Block Aikido uninstall`
   * **OMA-URI**: `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Data Type**: `String`
   * **Value**: paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
4. Go to **Policy Targets**, assign the policy to the same target devices or device groups, and save.

As a second layer, add the Aikido Device Protection installer to **Mandatory Apps** so that Hexnode automatically reinstalls it if it's ever removed:

1. In your Hexnode portal, go to **Policies** and edit the same Windows policy.
2. Under the **Windows** tab, go to **App Management** and click **Configure** next to **Mandatory Apps**.
3. Add the Aikido Device Protection app (uploaded to your Hexnode enterprise app catalog) and save.

{% hint style="warning" %}
The Aikido Device Protection service is hardened at the OS level by the installer to resist user tampering. No additional Hexnode policy is needed for service protection.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In your Hexnode portal, go to **Manage** → **Devices** and select the target Windows devices, or open a device group.
2. Choose **Actions** → **Device Control** → **Restart Device**.
3. Click **Restart** in the confirmation dialog.

The agent fully activates on the next boot.
{% endstep %}
{% endstepper %}

## Troubleshooting

| Problem                                                      | Fix                                                                                                                                                                             |
| ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Device doesn't appear in the dashboard after the script runs | Reboot the device; the agent fully registers on the next boot                                                                                                                   |
| Script exits with an error in the MDM console                | Run the script manually as SYSTEM on a test host to see the full error; confirm the device has internet access to download the installer                                        |
| Agent service is not running                                 | Reboot the device; if the service is still absent, re-run the install script                                                                                                    |
| Script is blocked by execution policy                        | Scripts run via MDM execute as SYSTEM and bypass user-level execution policy restrictions; confirm the MDM is not applying an additional policy that restricts script execution |
| Device shows as inactive after installation                  | The agent needs a reboot to fully activate; check back after rebooting                                                                                                          |
