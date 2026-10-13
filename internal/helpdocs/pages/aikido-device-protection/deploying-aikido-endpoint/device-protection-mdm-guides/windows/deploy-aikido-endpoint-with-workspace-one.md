# Deploy Device Protection with Omnissa Workspace ONE

Use Omnissa Workspace ONE to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Omnissa Workspace ONE Windows management before you begin.
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

1. In the Workspace ONE UEM console, go to **Resources** → **Scripts**.
2. Click **Add** and select **Windows** from the drop-down.
3. On the **General** tab, give the script a name (e.g. "Aikido Device Protection Install").
4. On the **Details** tab, set **Language** to **PowerShell** and **Execution Context and Privileges** to **System Context**.
5. Paste the PowerShell install command you copied from the Aikido dashboard into the script body.
6. On the **Assignment** tab, assign the script to the target Windows smart group and set the trigger to run on the next sample.
7. Click **Save**.
8. From the script's row in **Scripts**, click **More** → **Run Now** on a test device first and confirm a successful execution status before rolling out.

{% hint style="info" %}
Scripts require Intelligent Hub on the device. Check deployment status under **Devices** → **List View** → device → **More** → **Scripts**.
{% endhint %}
{% endstep %}

{% step %}
**Apply the lockdown policy**

Workspace ONE does not expose a single toggle that blocks stopping a specific Windows service. Combine the two controls below to harden the agent against tampering.

Create two separate Workspace ONE profiles: one to block uninstallation using AppLocker, and one to restrict access to the Windows Services console.

**Block uninstall of the Aikido Device Protection application**

1. In the Workspace ONE UEM console, go to **Resources** → **Profiles & Baselines** → **Profiles**.
2. Click **Add** → **Add Profile**, then select **Windows** → **Windows Desktop** → **Device Profile**.
3. Select the **Application Control** payload. Application Control uses Microsoft AppLocker to prevent users from adding or removing applications.
4. Paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) into the Application Control payload (open the link, copy the full XML text, and paste it here).
5. Assign the profile to the same smart group and click **Save & Publish**.

**Restrict access to the Services console**

1. Create a second profile: go to **Resources** → **Profiles & Baselines** → **Profiles**, click **Add** → **Add Profile**, then select **Windows** → **Windows Desktop** → **Device Profile**.
2. Select the **Custom Settings** payload.
3. Paste the following SyncML block into the Custom Settings payload:

   ```xml
   <SyncML>
     <SyncBody>
       <Replace>
         <CmdID>1</CmdID>
         <Item>
           <Target>
             <LocURI>./Device/Vendor/MSFT/Policy/Config/ADMX_MMCSnapins/MMC_Services</LocURI>
           </Target>
           <Meta><Format xmlns="syncml:metinf">chr</Format></Meta>
           <Data><disabled/></Data>
         </Item>
       </Replace>
     </SyncBody>
   </SyncML>
   ```
4. Assign the profile to the same smart group and click **Save & Publish**.

{% hint style="info" %}
The Aikido Device Protection installer also applies a service ACL that prevents standard users from stopping the `EndpointProtection` service directly. The Workspace ONE controls above add a second layer by stopping users from reaching the uninstaller or the Services console.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Workspace ONE UEM console, go to **Devices** → **List View**.
2. Filter to the target smart group, then select devices using Shift-click. (Avoid the global checkbox, bulk actions may not trigger correctly with it.)
3. Click **More Actions** → **Restart Device** and confirm.

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
