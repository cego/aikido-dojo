# Deploy Device Protection with SimpleMDM

Use SimpleMDM to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In SimpleMDM, go to **Profiles** and click **New Profile**.
2. Select **Custom Configuration Profile** as the profile type.
3. Give it a name (e.g. "Aikido Device Protection").
4. Upload the `.mobileconfig` file, or open it in a text editor and paste its XML contents.
5. Click **Create Profile**.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In SimpleMDM, go to **Apps** and click **Add App**.
2. Select **macOS Package**.
3. Upload the Aikido Device Protection `.pkg` file.
4. Click **Save App**.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Go to **Devices** → **Groups** and select your target device group.
2. On the **Profiles** tab, assign the **Aikido Device Protection** configuration profile.
3. Verify it is installed on all target devices.
4. Back in the device group, open the **Apps** tab and assign the **Aikido Device Protection** package.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In SimpleMDM, go to **Devices** and select each target device.
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
