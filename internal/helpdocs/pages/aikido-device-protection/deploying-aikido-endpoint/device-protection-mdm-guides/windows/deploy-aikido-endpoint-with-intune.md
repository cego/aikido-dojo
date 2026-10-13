# Deploy Device Protection with Intune

Use Microsoft Intune to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Microsoft Intune Windows management before you begin.
{% endhint %}

## Installation

{% stepper %}
{% step %}
**What you'll need**

Before starting, make sure you have the following from the Aikido Device Protection dashboard:

* Your **PowerShell install command** (copied from the user group selector on the dashboard, includes token + installer download)

If you're missing it, go back to the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device**.
{% endstep %}

{% step %}
**Run the install command on managed devices**

Deploy the dashboard command as a one-shot platform script that runs in the system context.

1. In the Intune admin center, go to **Devices** → **Scripts and remediations** → **Platform scripts** → **Add** → **Windows 10 and later**.
2. On the **Basics** tab, enter a name (e.g. "Aikido Device Protection install") and an optional description, then select **Next**.
3. On the **Script settings** tab, save the **PowerShell install command** from the Aikido dashboard as a `.ps1` file (the file must be under 200 KB and saved as ASCII), then upload it using the **Script file** browse button.
4. Set **Run this script using the logged on credentials** to **No** so the script runs as SYSTEM.
5. Set **Enforce script signature check** to **No**.
6. Set **Run script in 64-bit PowerShell host** to **Yes**.
7. Select **Next**, add any scope tags, then select **Next**.
8. On **Assignments**, select the target Windows device group, then select **Next**.
9. On **Review + add**, select **Add**.

{% hint style="info" %}
The Intune Management Extension checks in every 8 hours and after each reboot. To force a faster run on a test device, restart the **Microsoft Intune Management Extension** service on that device.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Block uninstallation of the Aikido Device Protection application.

**Block uninstall of the Aikido Device Protection application**

Intune has no built-in setting that simply hides "Uninstall" for a specific Win32 application. To prevent uninstall, deploy an AppLocker policy that denies execution of the Aikido uninstaller.

1. In the Intune admin center, go to **Devices** → **Manage devices** → **Configuration** → **Create** → **New policy**.
2. Set **Platform** to **Windows 10 and later** and **Profile type** to **Templates**, then select **Custom**. Select **Create**.
3. Give the profile a name (e.g. "Aikido Device Protection uninstall block").
4. Add an OMA-URI setting that pushes an AppLocker XML policy denying the Aikido uninstaller:
   * **Name**: `Aikido AppLocker uninstall block`
   * **OMA-URI**: `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Data type**: **String**
   * **Value**: paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
5. Assign the profile to the same Windows device group.

{% hint style="info" %}
If end users are local administrators, AppLocker rules can be bypassed. Combine this policy with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Intune admin center, go to **Devices** → **All devices** → **Bulk device actions**.
2. On the **Basics** tab, set **OS** to **Windows** and **Device action** to **Restart**, then select **Next**.
3. On the **Devices** tab, select the target devices (up to 100 per bulk action), then select **Next**. For larger groups, repeat the bulk action in batches, or assign a **Restart** compliance action to the device group via **Endpoint security** → **Device compliance** → **Policies** instead.
4. On **Review + create**, select **Create**.

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
