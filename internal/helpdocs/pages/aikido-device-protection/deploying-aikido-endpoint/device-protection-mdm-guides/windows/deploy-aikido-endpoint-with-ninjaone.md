# Deploy Device Protection with NinjaOne

Use NinjaOne to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in NinjaOne before you begin.
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

1. In the NinjaOne Admin Portal, go to **Administration** → **Library** → **Automation**.
2. Click **Add automation** → **Add new script**.
3. Configure the script:
   * **Name:** Aikido Device Protection installer
   * **Language:** PowerShell
   * **Operating System:** Windows
   * **Architecture:** All
   * **Run as:** System
4. Paste the PowerShell install command from the Aikido Device Protection dashboard into the editor and click **Save**.
5. In the left navigation, go to **Devices**, then filter by OS to show only Windows devices.
6. Select the target devices, then click **Run** in the toolbar, choose **Script**, and select **Aikido Device Protection installer**.
7. Confirm **Run As** is set to **System**, then click **Run**.

{% hint style="info" %}
You can also schedule this script under **Policies** → **Scheduled Automations** to cover devices that come online later.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Protect the Aikido service against tampering using NinjaOne's policy conditions. For uninstall prevention, use AppLocker outside NinjaOne.

1. In the NinjaOne Admin Portal, go to **Administration** → **Policies**.
2. Select (or create) the Windows policy assigned to your target devices.
3. Open the **Conditions** tab and click **Add a condition**.
4. Choose **Windows Service (Windows)** and configure:
   * **Service:** `EndpointProtection`
   * **State:** Stopped
5. Under the condition's **Automation**, add a script action that restarts `EndpointProtection` (for example, a PowerShell one-liner `Start-Service EndpointProtection` set to **Run as: System**).
6. Save the policy.

#### Block uninstallation of the Aikido Device Protection app

Create a NinjaOne automation script that applies the AppLocker policy via PowerShell:

1. In the NinjaOne Admin Portal, go to **Administration** → **Library** → **Automation** and click **Add automation** → **Add new script**.
2. Set **Language** to **PowerShell**, **Operating System** to **Windows**, and **Run as** to **System**.
3. Paste the following script and click **Save**:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

4. Run the script against your target Windows devices from the **Devices** tab.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the NinjaOne Admin Portal, open **Devices** and select your target Windows devices.
2. Click **Run** in the toolbar and choose **Reboot**. Select **Later** to schedule a restart, or **Now** for an immediate restart.
3. Pick a date, time, and time zone, then click **Set**.

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
