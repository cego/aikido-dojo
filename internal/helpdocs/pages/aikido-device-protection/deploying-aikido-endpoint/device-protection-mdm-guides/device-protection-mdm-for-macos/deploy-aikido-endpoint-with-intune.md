# Deploy Device Protection with Intune

Use Microsoft Intune to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In the Intune admin center, go to **Devices** → **macOS** → **Configuration** → **Create** → **New Policy**.
2. Select **macOS** as the platform. Under **Profile type**, select **Templates**, then choose **Custom**.
3. Give the profile a name (e.g. "Aikido Device Protection").
4. Under **Custom configuration profile name**, enter a name and browse to the `.mobileconfig` file to upload it.
5. Assign the profile to the target device group.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In the Intune admin center, go to **Apps** → **macOS** → **Add**.
2. Select **macOS app (PKG)** as the app type.
3. Upload the Aikido Device Protection `.pkg` file.
4. Complete the remaining tabs (**Detection rules**, **Assignments**) and assign the app as **Required** for the same device group.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Assign and sync the **Aikido Device Protection** configuration profile.
2. Verify it appears on a test device.
3. Assign the **Aikido Device Protection** app.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Intune admin center, go to **Devices** → **macOS**.
2. Filter to the target device group, select all devices, and click **Restart** from the bulk-action menu.

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
