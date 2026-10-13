# Deploy Device Protection with Absolute

Use Absolute to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must have the Absolute agent installed and reporting to the Absolute console before you begin.
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

Push the PowerShell command from the Aikido dashboard to your Windows fleet using Absolute's Device Reach feature.

1. In the Absolute console, go to **Device Reach** in the left navigation.
2. Select **Custom Scripts**, then click **Create Script**.
3. Set the platform to **Windows** and the script type to **PowerShell**.
4. Paste the install command copied from the Aikido Device Protection dashboard into the script editor.
5. Set the timeout to at least **30 minutes** to allow for the installer download.
6. Assign the script to the target device group and click **Run**.

{% hint style="info" %}
After the script runs, check the Device Reach results panel to confirm successful execution on each device.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Use Absolute's capabilities to ensure Aikido Device Protection stays installed and to harden the service against tampering.

**Block uninstallation of the Aikido Device Protection app**

Use Absolute's Device Reach to deliver an AppLocker policy that prevents users from running the Aikido uninstaller.

1. In the Absolute console, go to **Device Reach** in the left navigation.
2. Select **Custom Scripts**, then click **Create Script**.
3. Set the platform to **Windows** and the script type to **PowerShell**.
4. Paste the following PowerShell script into the script body:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

5. Assign the script to the target device group and click **Run**.

**Prevent removal via Application Persistence**

Use Absolute's Application Persistence feature to detect and heal removal of the Aikido Device Protection application.

1. In the Absolute console, go to **Resilience** in the left navigation.
2. Select **Application Persistence**, then click **Add application**.
3. Search for or enter the Aikido Device Protection application name.
4. Configure the persistence action (reinstall or alert) and assign it to the target device group.
5. Save the configuration.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Absolute console, go to **Devices** and select the target Windows devices.
2. Use **Device Reach** to run a PowerShell reboot command:

```powershell
Restart-Computer -Force
```

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
