# Deploy Device Protection with Pulseway

Use Pulseway to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must have the Pulseway agent installed and reporting as online before you begin.
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
**Create the install script**

1. In the Pulseway WebApp, go to **Automation** → **Scripts**.
2. Click **Create Category** and name it (e.g. "Aikido"), or open an existing category of your own.
3. Click **Create Script** and name it "Aikido Device Protection install".
4. Open the **Windows** tab and set it to **Enabled**.
5. Set the script type to **PowerShell**.
6. Paste the following guard, then the PowerShell install command copied from the Aikido Device Protection dashboard, into the editor:

```powershell
if (Get-Service -Name EndpointProtection -ErrorAction SilentlyContinue) { exit 0 }
# paste the install command from the dashboard below this line
```

7. Click **Save Script**.

{% hint style="info" %}
The Pulseway agent service runs as Local System, so scripts you run through it already have the privileges the installer needs to write to Program Files and register the service. There is no separate run-as setting.
{% endhint %}
{% endstep %}

{% step %}
**Run the install script on managed devices**

To deploy across a group of devices, use a task:

1. Go to **Automation** → **Tasks** and click **Create Task**.
2. Fill in a name and description (e.g. "Aikido Device Protection install").
3. Under **Scope**, select the systems or system group you want to target.
4. Open the **Scripts** tab, pick your category from the dropdown, and click the plus icon next to **Aikido Device Protection install** to add it to the task.
5. Click **Save**, then open the task and click **Run** to execute it immediately.

To install on a single device, go to **Systems**, select the device, open **Scripts**, choose your category, and run **Aikido Device Protection install**.

{% hint style="info" %}
Tick **Enable Scheduling** on the task and set it to repeat daily to cover devices that are offline now or enroll later. The guard at the top of the install script exits early on devices that already have the agent, so the installer does not run again.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Block uninstallation of the Aikido Device Protection application and keep the service running.

**Block uninstallation via AppLocker**

AppLocker is the built-in Windows feature that controls which applications users are allowed to run. This policy stops users from launching the Aikido uninstaller.

1. Go to **Automation** → **Scripts** and click **Create Script** in your category.
2. Name it "Aikido AppLocker uninstall block", open the **Windows** tab, set it to **Enabled**, and select **PowerShell**.
3. Paste the following script and click **Save Script**:

```powershell
$tmp = "$env:TEMP\aikido-applocker.xml"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml" -OutFile $tmp
Set-AppLockerPolicy -XmlPolicy $tmp -Merge
Remove-Item $tmp -Force
```

4. Add the script to a task scoped to your target devices, following the steps above, and run it.

**Monitor and restart the Aikido service**

1. Go to **Server Admin** → **Policies** and open (or create) the policy assigned to your Windows devices.
2. Open the **Monitored Services** category and add the `EndpointProtection` service so Pulseway raises a notification if it stops.
3. Save the policy.

Pulseway no longer runs actions directly from monitoring rules, so pair the notification with a scheduled task that starts the service again:

1. Create a script named "Aikido service restart" with the following body:

```powershell
if (Get-Service -Name EndpointProtection -ErrorAction SilentlyContinue) { Start-Service EndpointProtection }
```

2. Add it to a task scoped to your Windows devices, tick **Enable Scheduling**, and set it to repeat hourly.

{% hint style="info" %}
The Aikido Device Protection installer hardens the `EndpointProtection` service ACL at the OS level, so standard users cannot stop or disable it. If your users run as local administrators, AppLocker rules can be bypassed, so combine the uninstall block above with removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Pulseway WebApp, go to **Systems** and select the target Windows devices.
2. Run the **Restart** system command.

You can also schedule the restart by adding a script with `Restart-Computer -Force` to a task and setting it to run outside working hours.

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
