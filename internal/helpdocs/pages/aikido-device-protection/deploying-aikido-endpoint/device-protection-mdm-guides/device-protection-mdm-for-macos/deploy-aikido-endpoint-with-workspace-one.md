# Deploy Device Protection with Omnissa Workspace ONE

Use Omnissa Workspace ONE UEM to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in Workspace ONE UEM before you begin.
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
**Add the Aikido Device Protection configuration profile**

1. In the Workspace ONE UEM console, go to **Resources** → **Profiles and Baselines**.
2. Click **Add**, then **Upload Profile**.
3. Select **macOS** and upload the `.mobileconfig` file.
4. Click **Save and Continue**.
5. Assign the profile to the target smart group.
6. Click **Save & Publish**.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. Generate a `.plist` metadata file for the package using the Omnissa Admin Assistant Tool: drag and drop the `.pkg` onto the tool.
2. In the Workspace ONE UEM console, go to **Resources** → **Apps** → **Native** → **Internal**.
3. From the **Add** dropdown, select **Application File**.
4. Upload the Aikido Device Protection `.pkg` and click **Save**, then **Continue**.
5. Upload the generated `.plist` metadata file and click **Save**, then **Continue**.
6. Assign the app to the same smart group and publish.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Save and publish the **Aikido Device Protection** configuration profile.
2. Verify it is installed on a test device.
3. Publish the **Aikido Device Protection** app.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the Workspace ONE UEM console, go to **Devices** → **List View**.
2. Filter to the target smart group, select all devices, and click **More Actions** → **Restart**.

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
