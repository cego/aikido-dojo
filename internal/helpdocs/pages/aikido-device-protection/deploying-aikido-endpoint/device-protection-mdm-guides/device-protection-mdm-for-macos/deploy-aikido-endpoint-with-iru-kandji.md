# Deploy Device Protection with Iru / Kandji

Use Iru / Kandji to roll out Aikido Device Protection across your organization and protect managed macOS devices with a consistent setup.

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
**Upload the Aikido Device Protection configuration profile to Iru**

1. In Iru, go to the Library.
2. Upload the `.mobileconfig` file as a **Custom Profile**.
3. Assign it to the same Blueprint as the app.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In Iru, add the Aikido Device Protection `.pkg` as a custom app.
2. Set the custom app's install type to **Install once per device**.
3. Assign the app to the same Blueprint.

{% hint style="warning" %}
Choosing **Audit and enforce** makes Iru reinstall the package on every check-in, which can redeploy an outdated build. The package should only be installed once per device.
{% endhint %}
{% endstep %}

{% step %}
**Deploy to your devices**

1. Deploy the **Aikido Device Protection** configuration profile
2. Verify it has reached all target devices. Check status in the Iru console.
3. Deploy the **Aikido Device Protection** app
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
