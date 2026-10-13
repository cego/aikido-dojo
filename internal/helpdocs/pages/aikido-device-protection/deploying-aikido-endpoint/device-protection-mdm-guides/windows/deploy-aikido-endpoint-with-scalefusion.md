# Deploy Device Protection with Scalefusion

Use Scalefusion to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Scalefusion Windows management before you begin.
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

Deploy the dashboard command as a PowerShell script that runs at the device level.

1. Save the **PowerShell install command** from the Aikido dashboard as a `.ps1` file.
2. In the Scalefusion Dashboard, go to **Application Management** → **Enterprise Store** → **Upload New App**.
3. Choose **PowerShell (PS1)** from the drop-down, then upload the `.ps1` file.
4. Enter a script name (e.g. "Aikido Device Protection install") to identify it in Device Profiles, and set the **Execution Level** to **Device**, then click **SAVE**.
5. When Scalefusion finishes processing, click **PUBLISH**, set the schedule to **Run Once on Publish**, and select the target Windows Device Profile.

{% hint style="info" %}
The **Execution Level** must be **Device**. A user-level script does not have the privileges required to register the Aikido Device Protection service.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Block uninstallation of the Aikido Device Protection application by deploying an AppLocker policy that denies execution of the Aikido uninstaller.

1. In the Scalefusion Dashboard, go to **Device Profiles & Policies** → **Device Profiles** and edit the target Windows Device Profile.
2. Go to **Settings** → **Advanced Configurations** to open the payload editor.
3. Enter a SyncML payload that pushes an AppLocker XML policy denying the Aikido uninstaller:

   ```xml
   <Add>
       <CmdID>1</CmdID>
       <Item>
           <Target>
               <LocURI>./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy</LocURI>
           </Target>
           <Meta>
               <Format xmlns="syncml:metinf">chr</Format>
           </Meta>
           <Data>ESCAPED_APPLOCKER_XML_HERE</Data>
       </Item>
   </Add>
   ```
4. Replace `ESCAPED_APPLOCKER_XML_HERE` with the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link and copy the full XML text), with the XML special characters escaped so the policy is valid inside the `<Data>` element: replace every `&` with `&amp;`, every `<` with `&lt;`, and every `>` with `&gt;`.
5. Click **VALIDATE** and confirm the green checkmark.
6. Click **SAVE** to apply the payload to the same Windows Device Profile.

{% hint style="info" %}
If end users are local administrators, AppLocker rules can be bypassed. Combine this policy with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Scalefusion Dashboard, go to **Devices** and select each target device.
2. Click **Actions** and select **Restart**.

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
