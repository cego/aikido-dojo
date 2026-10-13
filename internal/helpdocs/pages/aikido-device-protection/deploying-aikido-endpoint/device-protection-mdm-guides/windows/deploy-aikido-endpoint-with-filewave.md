# Deploy Device Protection with FileWave

Deploy Aikido Device Protection to managed Windows devices with a FileWave fileset that runs the install command.

Use FileWave to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must have the FileWave client installed and reporting to your FileWave server before you begin.
{% endhint %}

## Prerequisites

Copy your **PowerShell install command** from the user group selector in the [Aikido app](https://app.aikido.dev/endpoint-protection/devices). It includes your token and the installer download.

If you're missing it, go to the [Aikido app](https://app.aikido.dev/endpoint-protection/devices) and click **Connect Device**.

## Deploy Device Protection

{% stepper %}
{% step %}

### Create the install fileset

1. In FileWave Central, go to **Filesets** and click **New Desktop Fileset** → **Empty**.
2. Name the fileset "Aikido Device Protection".
3. Select the fileset and click **Scripts** in the toolbar.
4. In the **Activation** section, create a new script named `aikido-install.ps1`.
5. Paste the PowerShell install command copied from the Aikido app into the script and save it.

{% hint style="info" %}
FileWave runs fileset scripts as the local SYSTEM account, so the installer has the privileges it needs to write to Program Files and register the service. There is no separate run-as setting.
{% endhint %}
{% endstep %}

{% step %}

### Add the lockdown script

AppLocker is the built-in Windows feature that controls which applications users are allowed to run. This script adds an AppLocker rule that stops users from launching the Aikido uninstaller.

1. In the same **Scripts** dialog, create a second **Activation** script named `aikido-applocker.ps1`.
2. Paste the following script and save it:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

3. Confirm `aikido-applocker.ps1` is listed below `aikido-install.ps1`. Scripts in the same section run from top to bottom, so drag them into this order if needed.

{% hint style="info" %}
The Aikido Device Protection installer hardens the `EndpointProtection` service at the OS level, so standard users can't stop or disable it. If your users run as local administrators, they can bypass AppLocker rules, so combine the uninstall block with removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}

### Force a reboot after installation

1. Select the **Aikido Device Protection** fileset and click **Properties**.
2. Check **Force reboot** and save.

FileWave restarts the device 2 minutes after the scripts run. The agent fully activates on the next boot.

To restart devices manually instead, right-click them in FileWave Central and select **Restart Windows Devices**.
{% endstep %}

{% step %}

### Deploy to your devices

1. In FileWave Central, go to **Deployments** and create a new deployment named "Aikido Device Protection".
2. Click **+** beside **Clients** and select the target Windows devices or group.
3. Click **+** beside **Filesets** and add the **Aikido Device Protection** fileset.
4. Click **OK**, then click **Update Model** in the toolbar to send the deployment to devices.

Start with a single test device. Open its **Client Info** → **Filesets Status** to confirm the fileset activated, and check that the device appears in the [Aikido app](https://app.aikido.dev/endpoint-protection/devices) after the reboot. Then add the rest of your fleet to the deployment and click **Update Model** again.
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
