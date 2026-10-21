# Deploy Device Protection with Primo

Use Primo to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Primo Windows management before you begin.
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

Primo runs PowerShell scripts on Windows hosts through the Controls library. Upload the dashboard command as a `.ps1` script, then run it against your target devices.

1. In the Primo admin console, go to **MDM** → **Controls** → **Scripts**.
2. Copy the PowerShell install command from the Aikido Device Protection dashboard. Save it as a `.ps1` file (e.g. `aikido-install.ps1`), then click **Add script**, upload the file, and give it a name (e.g. "Aikido Device Protection install").
3. Assign the script to the team or device group that should receive Aikido Device Protection.
4. Go to **Hosts**, filter to your target Windows devices, select them, open the **Actions** menu, and choose **Run script**. Select the Aikido install script.
5. Wait for the run to finish. Open the script run results and confirm every targeted device reports success before continuing.
   {% endstep %}

{% step %}
**Apply the lockdown policy**

End users must not be able to stop the Aikido service or uninstall the Aikido Device Protection app.

**Block uninstallation of Aikido Device Protection**

Run a PowerShell script that downloads and applies the AppLocker policy.

1. In the Primo admin console, go to **MDM** → **Controls** → **Scripts**.
2. Click **Add script**, name it (e.g. "Aikido uninstall block"), and paste the following:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

3. Apply it to your target Windows device group.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Add a reboot step to the install script: open the Aikido install script in **MDM** → **Controls** → **Scripts** (or add a follow-up script) and append `Restart-Computer -Force` so the device reboots at the end of the run.

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
