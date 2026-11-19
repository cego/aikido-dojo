# Deploy Device Protection with Ivanti

Use Ivanti Neurons for MDM to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All Macs must be enrolled in Ivanti Neurons for mobile device management (MDM) before you begin. Custom configurations need an Ivanti Silver license or higher.
{% endhint %}

## Prerequisites

Download the following from the [Aikido app](https://app.aikido.dev/endpoint-protection/devices?regionRedirect=1):

* The **Aikido Device Protection `.mobileconfig`**
* The **Aikido Device Protection installer** (`.pkg` file)

If you're missing any of these, go to the [Aikido app](https://app.aikido.dev/endpoint-protection/devices?regionRedirect=1), click **Connect Device**, and complete the pre-flight steps.

## Deploy Device Protection

{% stepper %}
{% step %}
### Add the configuration profile

1. In Ivanti Neurons for MDM, go to **Configurations** and click **+ Add**.
2. Type "custom" in the search field, then click **Custom**.
3. Enter a name, for example "Aikido Device Protection profile".
4. Under **Choose OS**, click the macOS icon.
5. Drag the `.mobileconfig` file into the drop box, or click **Choose File** to select it.
6. Click **Next**.
7. Set **Who does this configuration apply to** to **Device Wide**.
8. Select a distribution with your test Mac, then click **Done**.

{% hint style="warning" %}
Upload the `.mobileconfig` unchanged. Don't rebuild it with Ivanti's built-in macOS configurations, and keep **Device Wide** selected.
{% endhint %}

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in **System Settings** → **Login Items**.
{% endstep %}

{% step %}
### Add the installer

Deploy the profile first so the system extension is already approved when the installer runs.

1. Go to **Apps** → **App Catalog** and click **Add**.
2. Under **Select the platforms to upload file**, select **Apple Apps**, then select **macOS**.
3. Drag the Aikido Device Protection `.pkg` into the drop box, or click **Choose File** to select it, then click **Confirm**.
4. Click **Next** and step through the description pages without changes.
5. Select a distribution with your test Mac.
6. Expand **Advanced Options & App Configuration** and enable **Install on Device**.
7. Click **Done**.

{% hint style="info" %}
Upload the signed `.pkg` from the Aikido app as-is. Don't repackage it with MobileIron Packager.
{% endhint %}
{% endstep %}

{% step %}
### Restart the Mac

The agent fully activates on the next boot.

1. Go to **Devices** and click the test Mac in the **Name** column.
2. Click **Actions** → **Restart/Shutdown device**.
3. Select **Restart Device**, then click **Send to Device**.

To restart several Macs at once, select them in the **Devices** list and click **Actions** → **Restart/Shutdown Device**.
{% endstep %}

{% step %}
### Verify the deployment

On the test Mac, confirm:

1.  The system extension is activated:

    ```bash
    systemextensionsctl list | grep aikido
    ```

    Expect to see the extension marked `[activated enabled]`.
2. Open **System Settings** → **General** → **Login Items & Extensions** and confirm the Aikido Device Protection entries can't be toggled off.

Once the test Mac passes, change the distribution of the configuration and the app to the rest of your fleet.
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
