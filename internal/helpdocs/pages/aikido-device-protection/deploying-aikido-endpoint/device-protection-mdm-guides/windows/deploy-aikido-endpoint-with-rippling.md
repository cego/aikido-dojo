# Deploy Device Protection with Rippling

Use Rippling to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Rippling Device Management for Windows before you begin.
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

1. In Rippling, go to **IT** → **Device Management** → **Scripts** (the PowerShell script library for Windows).
2. Click **New script** and select **PowerShell** as the script type.
3. Set a **Name** (e.g. "Install Aikido Device Protection").
4. Set the **Run as** context to **SYSTEM** so the installer has the privileges to write to Program Files and register the service.
5. Set the **Schedule** to **Run once** so the script fires on each targeted device a single time.
6. Paste the PowerShell install command copied from the Aikido Device Protection dashboard into the script body.
7. Click **Save**, then click **Deploy** and select the target employees or device groups.

{% hint style="info" %}
The PowerShell command bundles your Device Protection token, downloads the latest installer, and installs silently. There is no extra MSI to upload.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

1. In Rippling, go to **IT** → **Device Management** → **Scripts** and click **New script**.
2. Select **PowerShell**, set **Run as** to **SYSTEM**, and set **Schedule** to **Run once**.
3. Paste the following script into the body:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

4. Click **Save**, then **Deploy** to the same device group.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Rippling does not expose a one-click restart action for Windows in the device list. Trigger the reboot with a short script:

1. In Rippling, go to **IT** → **Device Management** → **Scripts** and click **New script**.
2. Select **PowerShell**, set **Run as** to **SYSTEM**, and set the **Schedule** to **Run once**.
3. Paste the following into the script body:

```powershell
shutdown.exe /r /t 60 /c "Restarting to activate Aikido Device Protection"
```

4. Click **Save**, then **Deploy** to the same device group that received the installer.

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
