# Deploy Device Protection with ManageEngine

Use ManageEngine Endpoint Central to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in ManageEngine Endpoint Central Windows management before you begin.
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

1. In Endpoint Central, go to **Configurations** → **Add Configurations** → **Windows** → **Scripts** → **Custom Script**.
2. Give the script a name (e.g. "Aikido Device Protection install").
3. Set **Run as** to **System User** so the installer has the privileges to write to Program Files and register the service.
4. Set the execution frequency to **Once** so the script runs a single time on each target device.
5. Paste the PowerShell install command copied from the Aikido Device Protection dashboard into the script body.
6. Click **Save and Publish**, then assign the configuration to the target devices or device groups.

{% hint style="info" %}
The PowerShell command bundles your Device Protection token, downloads the latest installer, and installs silently. There is no extra MSI to upload.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Block uninstallation of the Aikido Device Protection application.

**Block uninstallation via AppLocker**

1. Create a new Custom Script configuration using the same steps above.
2. Give it a name (e.g. "Aikido - AppLocker uninstall block").
3. Paste the following PowerShell script into the script body:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

4. Click **Save and Publish** and assign it to the target devices or device groups.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In Endpoint Central, go to **Tools** → **Shutdown/Restart**.
2. Select the target Windows devices from the list.
3. Choose **Restart** and click **Proceed**.

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
