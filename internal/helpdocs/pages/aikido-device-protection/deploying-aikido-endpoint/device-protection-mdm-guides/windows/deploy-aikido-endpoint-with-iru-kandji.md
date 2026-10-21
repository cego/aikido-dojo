# Deploy Device Protection with Iru / Kandji

Use Iru / Kandji to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Iru Windows management before you begin. Iru manages Windows 11 (24H2 or 25H2) on the Pro, Pro Education, Enterprise, and Education editions.
{% endhint %}

## Installation

{% stepper %}
{% step %}
**What you'll need**

Before starting, make sure you have the following from the Aikido Device Protection dashboard:

* Your **PowerShell install command** (copied from the user group selector on the dashboard, includes token + installer download)

If you're missing it, go back to the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device**.

{% hint style="info" %}
The PowerShell command bundles your Device Protection token, downloads the latest installer, and installs silently. Unlike the macOS rollout, there is no configuration profile or installer package to upload, so you don't need a Windows Custom App Library Item.
{% endhint %}
{% endstep %}

{% step %}
**Run the install command on managed devices**

Deploy the dashboard command as a Windows Custom Script. The Iru Agent runs Windows custom scripts once per device.

1. In Iru, go to **Library** → **Add Library Item** and select **Windows Custom Script**.
2. Enter a **Name** (e.g. "Install Aikido Device Protection") and confirm the status is **Active**.
3. Leave **Execute in** set to **64-bit**.
4. Paste the following into the **Audit script** field:

```powershell
if (Get-Service -Name EndpointProtection -ErrorAction SilentlyContinue) { exit 0 }
exit 1
```

5. Add a **Remediation script** and paste the PowerShell install command copied from the Aikido Device Protection dashboard, then add this restart line as the last line of that same script (the **Reboot devices after installation** step below explains why):

```powershell
# last line of the remediation script, after the install command
shutdown.exe /r /t 300 /c "Restarting to activate Aikido Device Protection"
```

6. Click **Save**, then assign the Library Item to the Blueprint targeting your Windows devices through its Assignment Map.

The audit script reports the device as non-compliant while the `EndpointProtection` service is missing, which is what triggers the remediation script to install the agent. Devices that enroll later run the same check once, on their first check-in.
{% endstep %}

{% step %}
**Apply the lockdown policy**

Iru has no custom profile or OMA-URI payload for Windows, so deliver the AppLocker policy that blocks the Aikido uninstaller as a second Windows Custom Script.

1. In Iru, go to **Library** → **Add Library Item** and select **Windows Custom Script**.
2. Enter a **Name** (e.g. "Block Aikido uninstall") and confirm the status is **Active**.
3. Paste the following into the **Audit script** field:

```powershell
if ((Get-AppLockerPolicy -Local -Xml) -like "*c0a1b2c3-d4e5-46f7-8a9b-0c1d2e3f4a5b*") { exit 0 }
exit 1
```

4. Add a **Remediation script** and paste:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

5. Click **Save** and assign the Library Item to the same Blueprint.

The audit script looks for the Aikido deny rule in the local AppLocker policy, so remediation only runs on devices that don't have the rule yet.

{% hint style="info" %}
The Aikido Device Protection installer hardens the `EndpointProtection` service ACL at the OS level, so standard users cannot stop or disable it. If your users run as local administrators, AppLocker rules can be bypassed, so combine the uninstall block above with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

Iru's Windows device actions don't include a restart command, which is why the install remediation script in step 2 ends with:

```powershell
shutdown.exe /r /t 300 /c "Restarting to activate Aikido Device Protection"
```

The remediation script runs once per device, so the restart is scheduled only on the run that installs the agent. If you'd rather not restart devices from a script, leave the line out of step 2 and ask users to reboot instead.

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
