# Deploy Device Protection with Ivanti

Use Ivanti Neurons for MDM to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Ivanti Neurons for mobile device management (MDM) before you begin. Running PowerShell scripts needs Ivanti Bridge, which is part of the Ivanti Secure UEM package (legacy Gold). Ivanti Bridge doesn't support Windows 10 Home.
{% endhint %}

## Prerequisites

* The **PowerShell install command** from the [Aikido app](https://app.aikido.dev/endpoint-protection/devices?regionRedirect=1). Click **Connect Device** to copy it.
* Ivanti Bridge installed on the target devices. Ivanti adds the Bridge app to the **App Catalog** and distributes it to all devices by default. Check that it's in **Apps** → **App Catalog** and distributed to your devices.

## Deploy Device Protection

{% stepper %}
{% step %}
### Save the install script

1. Save the PowerShell install command from the Aikido app as `aikido-install.ps1`.
2. In Ivanti Neurons for MDM, go to **Admin** → **Scripts** → **All Scripts** and click **+ Add**.
3. Enter a name, for example "Aikido Device Protection install".
4. Under **Platform type**, select **Windows**, then set **Script Type** to **ps1**.
5. Click **Import code from a Script**, select `aikido-install.ps1`, and click **Import**.
6. Click **Save**.
{% endstep %}

{% step %}
### Run the script with Ivanti Bridge

1. Go to **Configurations** and click **+ Add**.
2. Select **Ivanti Bridge for Windows Scripts and Actions Configuration**.
3. Enter a name, for example "Aikido Device Protection install".
4. In the **Configure Ivanti Bridge** section, choose **Select scripts from Scripts Library** and select the "Aikido Device Protection install" script.
5. Click **Next**.
6. Select a distribution with your test device, then click **Done**.

The script downloads the installer, registers the device against your user group token, and installs silently. Check the result on the device's details page, on the **Logs** tab.
{% endstep %}

{% step %}
### Block uninstall with AppLocker

1. Go to **Configurations**, click **+ Add**, type "custom" in the search field, and click **Custom**.
2. Enter a name, for example "Aikido Device Protection uninstall block".
3. Under **Choose OS**, click the Windows icon, then select **Custom CSP OMA-URI Schema Node**.
4. Under **ACTIONS**, click **+** and set the following:
   * **Description**: `Aikido AppLocker uninstall block`
   * **OMA-URI**: `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Data type**: **CHARACTER**
   * **Value**: paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
   * **Access Type**: **Add**
5. Click **Save & Close**, then click **Next**.
6. Select the same distribution as the install script, then click **Done**.

{% hint style="info" %}
If end users are local administrators, they can bypass AppLocker rules. Combine this policy with App Control for Business (WDAC) and removal of local admin rights for stronger enforcement.
{% endhint %}
{% endstep %}

{% step %}
### Restart the device

The agent fully activates on the next boot.

1. Go to **Devices** and click the test device in the **Name** column.
2. Click **Actions** → **Restart/Shutdown device**.
3. Select **Restart Device**, then click **Send to Device**.

To restart several devices at once, select them in the **Devices** list and click **Actions** → **Restart/Shutdown Device**.

Once the test device passes, change the distribution of both configurations to the rest of your fleet.
{% endstep %}
{% endstepper %}

## Troubleshooting

| Problem | Fix |
| ------- | --- |
| Device doesn't appear in the dashboard after the script runs | Reboot the device; the agent fully registers on the next boot |
| Script exits with an error in the MDM console | Run the script manually as SYSTEM on a test host to see the full error; confirm the device has internet access to download the installer |
| Agent service is not running | Reboot the device; if the service is still absent, re-run the install script |
| Script is blocked by execution policy | Scripts run via MDM execute as SYSTEM and bypass user-level execution policy restrictions; confirm the MDM is not applying an additional policy that restricts script execution |
| Device shows as inactive after installation | The agent needs a reboot to fully activate; check back after rebooting |
