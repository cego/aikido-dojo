# Deploy Device Protection with FileWave

Deploy Aikido Device Protection to managed Macs with FileWave filesets and a deployment.

Use FileWave to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All Macs must be enrolled in FileWave mobile device management (MDM) before you begin. FileWave delivers configuration profiles through MDM, so unenrolled Macs don't get the profile.
{% endhint %}

## Prerequisites

Download the following from the [Aikido app](https://app.aikido.dev/endpoint-protection/devices):

* The **Aikido Device Protection `.mobileconfig`**
* The **Aikido Device Protection installer** (`.pkg` file)

If you're missing any of these, go to the [Aikido app](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device**, and complete the pre-flight steps.

## Deploy Device Protection

{% stepper %}
{% step %}

### Add the configuration profile

1. In FileWave Central, go to **Filesets**.
2. Drag the `.mobileconfig` file from Finder into the **Filesets** view.
3. Rename the new Profile fileset to "Aikido Device Protection profile".

{% hint style="warning" %}
Don't recreate or edit the profile in FileWave's Profile editor. The editor only knows Apple's predefined payloads and can change profiles with keys it doesn't recognize. Dragging the file in keeps it unchanged.
{% endhint %}

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}

### Add the installer

1. In **Filesets**, drag the Aikido Device Protection `.pkg` into the view, or click **New Desktop Fileset** → **MSI / PKG** and select the file.
2. Rename the new PKG fileset to "Aikido Device Protection installer".
3. Wait for the upload to finish. FileWave shows the progress in the bottom bar.
   {% endstep %}

{% step %}

### Force a reboot after installation

1. Select the **Aikido Device Protection installer** fileset and click **Properties**.
2. Check **Force reboot** and save.

FileWave restarts the Mac 2 minutes after the installer runs. The agent fully activates on the next boot.
{% endstep %}

{% step %}

### Deploy to your devices

Deploy the profile first so the system extension is already approved when the installer runs.

1. In FileWave Central, go to **Deployments** and create a new deployment named "Aikido Device Protection".
2. Click **+** beside **Clients** and select the target Macs or group.
3. Click **+** beside **Filesets** and add the **Aikido Device Protection profile** fileset.
4. Click **OK**, then click **Update Model** in the toolbar to send the deployment to devices.
5. Open a target Mac's **Client Info** → **Filesets Status** and confirm the profile is installed.
6. Edit the deployment, click **+** beside **Filesets**, and add the **Aikido Device Protection installer** fileset.
7. Click **OK**, then click **Update Model** again.

Start with a single test Mac. Once it passes the checks in the next step, add the rest of your fleet to the deployment and click **Update Model**.
{% endstep %}

{% step %}

### Verify the deployment

In FileWave Central, open the test Mac's **Client Info** → **Filesets Status** and confirm both filesets are active.

On the test Mac, confirm:

1. The system extension is activated:

   ```bash
   systemextensionsctl list | grep aikido
   ```

   Expect to see the extension marked `[activated enabled]`.
2. Open **System Settings** → **General** → **Login Items & Extensions** and confirm the Aikido Device Protection entries can't be toggled off.
   {% endstep %}
   {% endstepper %}

## Troubleshooting

| Problem                                              | Fix                                                                                                  |
| ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| Users still see the popup                            | Make sure the configuration profile is scoped correctly and installed before the package policy runs |
| Extension is waiting for user approval               | Check the System Extensions payload and confirm the team ID and bundle ID match exactly              |
| The package installs but the device does not connect | Confirm the token script ran before the package install and that the token was copied correctly      |
| Duplicate extension entries appear                   | Reboot the device                                                                                    |
| The content filter is not approved silently          | Re-upload the `.mobileconfig` profile and verify it is installed on the device                       |
