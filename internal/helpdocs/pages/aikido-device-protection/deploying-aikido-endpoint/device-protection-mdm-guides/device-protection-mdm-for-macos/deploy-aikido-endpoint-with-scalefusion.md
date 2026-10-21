# Deploy Device Protection with Scalefusion

Use Scalefusion to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

## Installation

{% stepper %}
{% step %}
**What you'll need**

Before starting, make sure you have the following from the Aikido Device Protection dashboard:

* The **Aikido Device Protection `.mobileconfig`** (downloaded from the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices))
* The **Aikido Device Protection installer** (`.pkg` file, downloaded from the dashboard)

If you're missing any of these, go back to the [Aikido Device Protection dashboard](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device**, and complete the pre-flight steps.
{% endstep %}

{% step %}
**Add the Aikido Device Protection configuration profile**

1. In the Scalefusion Dashboard, go to **Device Profiles & Policies** → **Device Profiles** and edit the macOS Device Profile you want to target.
2. In the profile, open the **Advanced Configurations** section to access the payload editor.
3. Paste the contents of the `.mobileconfig` into the payload editor, or use **IMPORT FILES** to upload the `.mobileconfig` directly.
4. Click **VALIDATE** and confirm the green checkmark.
5. Click **SAVE** to push the payload to devices in the profile.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In the Scalefusion Dashboard, go to **Enterprise** → **My Apps** → **Enterprise Store**.
2. Click **Upload New App** → **Upload macOS App**.
3. Upload the Aikido Device Protection `.pkg` file (or provide a link to it hosted in your cloud storage).
4. Confirm the **App Name**, **Bundle ID**, and **App Version**, then click **SAVE**.

{% hint style="warning" %}
Scalefusion only distributes **signed** PKG files. The Aikido Device Protection installer downloaded from the dashboard is signed, so upload it as-is; do not repackage or re-sign it.
{% endhint %}
{% endstep %}

{% step %}
**Deploy to your devices**

1. Confirm the **Advanced Configurations** payload from the earlier step is applied to all devices in the target macOS Device Profile.
2. Back in **Enterprise Store**, select the uploaded Aikido Device Protection app.
3. Click **PUBLISH** and choose the same Device Profile.
4. Click **PUBLISH** again to confirm.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Scalefusion Dashboard, go to **Devices** and select each target device.
2. Click **Actions** and select **Restart**.

{% hint style="info" %}
The Restart action requires devices to be supervised. If your devices are not supervised, trigger a reboot via a script job that runs `sudo reboot`.
{% endhint %}

The agent fully activates on the next boot.
{% endstep %}

{% step %}
**Verify the deployment**

On a test device, confirm:

1. The system extension is activated:

   ```bash
   systemextensionsctl list | grep aikido
   ```

   Expect to see the extension marked `[activated enabled]`.
2. Open **System Settings** → **General** → **Login Items & Extensions** and confirm the Aikido Device Protection entries cannot be toggled off.
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
