# Deploy Device Protection with Fleet

Use Fleet to deploy Aikido Device Protection across your managed Windows fleet with the required permissions in place.

## Installation

{% stepper %}
{% step %}
**What you'll need**

Before starting, make sure you have the following from the Aikido Device Protection dashboard:

* Your **PowerShell install command** (copied from the user group selector on the dashboard, includes token + installer download)

If you're missing it, go back to the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices) and click **Connect Device**.
{% endstep %}

{% step %}
**Run the install command on managed devices**

Deploy the PowerShell install command as a one-shot script running in the system context. Fleet runs PowerShell scripts on Windows hosts as the local SYSTEM account through `fleetd`.

1. In the Aikido Device Protection dashboard, copy the PowerShell install command for the target user group.
2. Save the command as a `.ps1` file, for example `aikido-install.ps1`.
3. In the Fleet console, go to **Controls** → **Scripts**.
4. Select your target team from the team dropdown in the upper-left corner.
5. Click **Upload** and add the `aikido-install.ps1` file.
6. For fleet-wide deployment, run: `fleetctl run-script --script-path aikido-install.ps1 --hosts <hostnames>` (or use a label to target all Windows hosts). For individual hosts, open the host's detail page and click **Actions** → **Run script**, then pick `aikido-install.ps1` and click **Run**.
7. Watch each host's **Activity** feed to confirm the script exited successfully.
   {% endstep %}

{% step %}
**Apply the lockdown policy**

Block uninstallation by pushing a configuration profile via Fleet's MDM profile management.

**Block uninstall via AppLocker**

1. Click **Add profile** again and create a second Windows profile with the following OMA-URI:
   * **OMA-URI:** `./Vendor/MSFT/AppLocker/ApplicationLaunchRestrictions/Aikido/MSI/Policy`
   * **Value (String):** paste the contents of [aikido-applocker-deny-mdm.xml](https://raw.githubusercontent.com/AikidoSec/safechain-internals/refs/heads/main/docs/aikido-applocker-deny-mdm.xml) (open the link, copy the full XML text, and paste it here)
2. Save and deploy the profile to your Windows hosts.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Fleet does not provide a team-wide reboot button in the UI. Use the Windows `RebootNow` MDM command through `fleetctl`, or push a one-line PowerShell restart script.

1. Save the following as `windows-restart-device.xml`:

   ```xml
   <Exec>
     <Item>
       <Target>
         <LocURI>./Device/Vendor/MSFT/Reboot/RebootNow</LocURI>
       </Target>
       <Meta>
         <Format xmlns="syncml:metinf">null</Format>
         <Type>text/plain</Type>
       </Meta>
       <Data></Data>
     </Item>
   </Exec>
   ```
2. Run `fleetctl mdm run-command --payload=windows-restart-device.xml --hosts=<host1>,<host2>` against each target host.
3. Alternatively, upload an `aikido-reboot.ps1` script containing `Restart-Computer -Force` through **Controls** → **Scripts** and run it against the target hosts.

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
