# Deploy Device Protection with SCCM

Deploy Aikido Device Protection to managed Windows devices with a Microsoft Configuration Manager (SCCM) package that runs the install command.

Use Microsoft Configuration Manager (SCCM) to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

## Installation

{% hint style="info" %}
All devices must run the Configuration Manager client and be members of a device collection before you begin.
{% endhint %}

{% stepper %}
{% step %}

### Prepare the package source

1. Copy the **PowerShell install command** from the [Aikido app](https://app.aikido.dev/endpoint-protection/devices) (click **Connect Device**).
2. Save it as `aikido-install.ps1` in a folder on a network share the site server can read, for example `\\<server>\<share>\Aikido`.
3. In the same folder, save the following script as `aikido-applocker.ps1`. It turns on the Application Identity service that AppLocker needs to enforce rules, then adds an AppLocker rule that stops users from launching the Aikido uninstaller.

{% code title="aikido-applocker.ps1" %}

```powershell
sc.exe config appidsvc start= auto
Start-Service -Name AppIDSvc
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

{% endcode %}
{% endstep %}

{% step %}

### Create the package and install program

1. In the Configuration Manager console, go to **Software Library** → **Application Management** → **Packages** and click **Create Package**.
2. Enter a name (e.g. "Aikido Device Protection"), select **This package contains source files**, browse to the UNC path of the folder, then click **Next**.
3. Select **Standard program**, then click **Next**.
4. Set the following, then finish the wizard:
   * **Name**: `Install`
   * **Command line**: `%windir%\Sysnative\WindowsPowerShell\v1.0\powershell.exe -NoProfile -ExecutionPolicy Bypass -File aikido-install.ps1`
   * **Run**: **Hidden**
   * **Program can run**: **Whether or not a user is logged on** (runs as SYSTEM)
     {% endstep %}

{% step %}

### Add the lockdown program

1. Select the **Aikido Device Protection** package, open the **Programs** tab, and click **Create Program**.
2. Select **Standard program** and set the following:
   * **Name**: `Block uninstall`
   * **Command line**: `%windir%\Sysnative\WindowsPowerShell\v1.0\powershell.exe -NoProfile -ExecutionPolicy Bypass -File aikido-applocker.ps1`
   * **Run**: **Hidden**
   * **Program can run**: **Whether or not a user is logged on**
3. On the **Requirements** page, check **Run another program first**, select the **Aikido Device Protection** package and the **Install** program, then finish the wizard.
4. On the **Programs** tab, right-click **Block uninstall** and select **Properties**.
5. On **General**, set **After running** to **Configuration Manager restarts computer**, then click **OK**.

The agent fully activates on the next boot, so the restart after this program completes the installation.

{% hint style="info" %}
If end users are local administrators, AppLocker rules can be bypassed. Combine this policy with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}

### Distribute and deploy the package

1. Select the **Aikido Device Protection** package, click **Distribute Content**, and add your distribution points.
2. Select the **Block uninstall** program and click **Deploy**. Configuration Manager runs **Install** first.
3. On **General**, select the target device collection.
4. On **Deployment Settings**, set **Purpose** to **Required**.
5. On **Scheduling**, add an assignment schedule of **As soon as possible**, then finish the wizard.

{% hint style="info" %}
Clients pick up new deployments at the next machine policy polling interval (60 minutes by default). To force a faster run on a test device, right-click it in the console and select **Client Notification** → **Download Computer Policy**.
{% endhint %}
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
