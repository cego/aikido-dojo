# Deploy Device Protection with Rippling

Use Rippling to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Rippling MDM before you begin.
{% endhint %}

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
**Upload the Aikido Device Protection configuration profile**

1. In Rippling, go to **IT** → **Device Management** → **Configurations** and select the **macOS** tab.
2. Click **Upload**.
3. Enter a **Policy name** (e.g. "Aikido Device Protection").
4. Set **Platform** to **macOS**.
5. Drop or select the `.mobileconfig` file.
6. Click **Save & continue**.

{% hint style="info" %}
Rippling requires `.mobileconfig` uploads to be done from a Mac. Uploading from another operating system may silently fail.
{% endhint %}

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Deploy the configuration profile**

1. In **Configurations**, switch to the **Everything Else** tab.
2. Find the profile you just created.
3. Click the **three-dot menu** on the right and select **Deploy**.
4. Select the target employees or devices and click **Save**.
   {% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In Rippling, go to **IT** → **Software**.
2. Click **Upload Software**.
3. Fill in the following fields:
   * **Name:** Aikido Device Protection
   * **Operating System:** macOS
   * **Description:** Aikido Device Protection
4. Under **Upload Installer File**, drop or select the Aikido Device Protection `.pkg` file.
5. Click **Submit**.
6. Click **Add** on the newly created software item, then click **Finished Selecting**.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Confirm the **Aikido Device Protection** configuration profile is installed on a test device.
2. Find the **Aikido Device Protection** software item under **IT** → **Software**, click **Edit**, select the target employees or devices, and click **Save**.

{% hint style="info" %}
Newly created software items in Rippling may stay in a "Pending" status for a few minutes before they can be deployed.
{% endhint %}
{% endstep %}

{% step %}
**Reboot devices after installation**

1. In Rippling, go to **IT** → **Devices**.
2. Select the target devices and choose **Restart** from the bulk actions menu (or open an individual device and select **Restart** from device actions).

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
