# Deploy Device Protection with Hexnode

Use Hexnode UEM to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In your Hexnode portal, go to **Policies** and create a new blank policy for macOS.
2. Under the **macOS** tab, go to **Configurations** and click **Configure** next to **Deploy Custom Configuration**.
3. Click **Upload** and select the `.mobileconfig` file.
4. Go to **Policy Targets**, assign the policy to the target devices or device groups, and save.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In your Hexnode portal, go to **Apps** → **+Add Apps** → **Enterprise App**.
2. Select **macOS** as the platform and upload the Aikido Device Protection `.pkg` file.
3. Complete the app details and save to finish adding it to your app inventory.
   {% endstep %}

{% step %}
**Create the install policy**

1. In your Hexnode portal, go to **Policies** and create a new blank policy for macOS.
2. Under the **macOS** tab, go to **App Management** and click **Configure** next to **Mandatory Apps**.
3. Add the uploaded Aikido Device Protection app.
4. Go to **Policy Targets**, assign it to the same target devices or device groups, and save.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Deploy the **Aikido Device Protection** configuration profile policy.
2. Verify it has reached all target devices. Check device compliance status in the Hexnode console.
3. Deploy the **Aikido Device Protection** install policy.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In your Hexnode portal, go to **Manage** → **Devices** and select the target devices.
2. Choose **Actions** → **Device Control** → **Restart Device**.

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
