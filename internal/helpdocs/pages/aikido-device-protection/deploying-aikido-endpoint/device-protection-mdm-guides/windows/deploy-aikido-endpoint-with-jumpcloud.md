# Deploy Device Protection with JumpCloud

Use JumpCloud to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in [JumpCloud Windows MDM](https://jumpcloud.com/support/get-started-windows-mdm) before you begin.
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

1. In the JumpCloud Admin Portal, go to **Device Management** → **Commands**.
2. Click **+ Command** and select **Command**.
3. Give it a name (e.g. "Install Aikido Device Protection").
4. Set **Type** to **Windows PowerShell**. JumpCloud runs Windows commands as the `LocalSystem` (NT Authority\System) account, so no extra user context is needed.
5. Paste the PowerShell install command you copied from the Aikido Device Protection dashboard into the command body.
6. Leave **Launch Type** as **One-Time**.
7. Select the **Devices** tab and assign the command to your target Windows device group.
8. Click **save**, then run the command using **Run Now**.

The command pulls the installer, registers the device against your user group token, and installs silently.
{% endstep %}

{% step %}
**Apply the lockdown policy**

JumpCloud doesn't offer a single built-in policy that both locks a Windows service and blocks an app uninstaller, so apply these two policies together.

**Block uninstallation of the Aikido Device Protection app**

1. In the JumpCloud Admin Portal, go to **Device Management** → **Policy Management**.
2. Click **+ Add New**, select **Device Policy**, then select the **Windows** tab.
3. Search for and select **Custom MDM (OMA URI) Policy**, then click **Configure**.
4. Give the policy a name (e.g. "Aikido Device Protection uninstall block").
5. In the **Settings** section, enter the following values:
   * **OMA URI**: `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Format**: **string**
   * **Value**: paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
6. Select the **Device Groups** tab, select the same Windows device group, and click **Create Policy**.

**Prevent users from stopping or disabling the Aikido service**

The Aikido Device Protection installer hardens its own service ACL so non-admin users cannot stop or disable it. No additional JumpCloud policy is required.

{% hint style="warning" %}
JumpCloud's Custom Registry Keys policy cannot modify ACLs under the `Security\` registry path. Service permission hardening is therefore enforced by the Aikido installer itself, not by JumpCloud.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In the JumpCloud Admin Portal, go to **Device Management** → **Policy Management**.
2. Click the **(+)** icon, select the **Windows** tab, choose **Schedule Restart**, and click **configure**.
3. Set the frequency to **Single** and pick a date and time shortly after your install command runs.
4. Select the **Device Groups** tab and assign the policy to the same Windows device group.
5. Click **save**.

For an immediate one-off restart of a single device, go to **Device Management** → **Devices**, click the three dots in the **Actions** column, and choose **Restart Device**.

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
