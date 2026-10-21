# Deploy Device Protection with Miradore

Use Miradore to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Miradore Windows management using the **Full enrollment** method. Software deployment and remote reboot are not supported on Windows 10/11 Home.
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

Wrap the dashboard command in a `.ps1` file and deploy it as a Windows application via a Business Policy.

1. Save the PowerShell command from the dashboard as a file named `install-aikido.ps1`.
2. In the Miradore admin console, go to **Management** → **Applications**.
3. Click **Add** → **Windows application** → **Advanced (file upload)** → **Next**.
4. Fill in the required fields:
   * **Application name** (e.g. "Aikido Device Protection installer")
   * **File**: select `install-aikido.ps1`
   * **Command type**: **Execute PowerShell**
   * **Install command**: `powershell.exe -ExecutionPolicy Bypass -File .\install-aikido.ps1`
5. Click **Create**, then **Close**.
6. Go to **Management** → **Business policies** and click **Add**.
7. Give the policy a name, add the Windows application you just created, assign your target devices using tags, and set the status to **Enabled**.

The Miradore Client executes the script on each in-scope device the next time it checks in. The Miradore Client runs as a Windows service under the SYSTEM account, so the installer has the privileges it needs to register the agent and write the token.
{% endstep %}

{% step %}
**Apply the lockdown policy**

Deploy a Custom Policy using the AppLocker OMA-URI to block uninstallation of Aikido Device Protection.

1. Go to **Management** → **Configuration profiles**, click **Add** → **Windows** → **Custom policy**.
2. Add a payload with the following values:
   * **OMA-URI**: `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Data type**: String
   * **Value**: paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
3. Click **Add**, then attach the profile to the same Business Policy you used for the installer.

As a supplementary control, you can also block MDM unenrollment: add a second payload with **OMA-URI** `./Vendor/MSFT/Policy/Config/Experience/AllowManualMDMUnenrollment`, **Data type** Integer, **Value** `0`.

{% hint style="warning" %}
The Aikido Device Protection installer hardens the EndpointProtection service ACL at the OS level, so standard users cannot stop or disable it. If your users run as local administrators, combine the AppLocker uninstall block above with removal of local admin rights for stronger enforcement.
{% endhint %}

{% hint style="info" %}
To restrict local administrator rights on managed devices, use your existing Windows Group Policy or Active Directory configuration. Miradore does not enforce account privilege restrictions directly.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

Trigger a reboot on each target device from the Miradore admin console after the installer reports success.

1. Go to **Management** → **Devices** and open a target device.
2. From the page toolbar, click **Actions** → **Reboot device**.
3. Confirm in the wizard by clicking **Reboot**.
4. Repeat for each managed device.

The device restarts about five minutes after the action reaches it. The agent fully activates on the next boot.
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
