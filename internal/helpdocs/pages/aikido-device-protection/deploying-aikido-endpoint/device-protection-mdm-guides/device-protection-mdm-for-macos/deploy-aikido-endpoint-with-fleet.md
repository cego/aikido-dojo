# Deploy Device Protection with Fleet

Use Fleet to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In the Fleet console, go to **Controls** → **OS settings** → **Configuration profiles**.
2. Select your target fleet from the dropdown in the upper-left corner.
3. Click **Add profile** and upload the `.mobileconfig` file. Fleet signs it automatically.
4. Monitor deployment status by hovering over the profile row and clicking the **(i)** icon. Profiles move from **Verifying** to **Verified** within approximately one hour.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload and configure the Aikido Device Protection installer**

1. In the Fleet console, go to **Software** in the main navigation.
2. Select your target fleet from the dropdown in the upper-left corner.
3. Click **Add software** → **Custom package** and upload the Aikido Device Protection `.pkg`.
4. Click **Add software**.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Confirm the **Aikido Device Protection** configuration profile shows **Verified** for all target hosts. Check **Controls** → **OS settings** → **Configuration profiles** → **(i)**.
2. Install the **Aikido Device Protection** package: open the target host's detail page, go to **Software** → **Library**, find the package, and click **Install**.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Restart your target devices after the agent installation completes. The agent fully activates on the next boot.
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
