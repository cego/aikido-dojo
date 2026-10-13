# Deploy Device Protection with Jamf

Use Jamf Pro to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
These steps use Jamf Pro. If you use Jamf School, the menus and payload names can differ.
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
**Sign the `.mobileconfig` before uploading**

{% hint style="warning" %}
You **must** sign the `.mobileconfig` before uploading it to Jamf Pro. Jamf Pro does not support some of the newer Apple payloads used by Aikido Device Protection, and an unsigned profile will fail to deploy correctly.
{% endhint %}

Sign the downloaded `.mobileconfig` using a signing certificate from Jamf Pro's built-in CA. Follow Jamf's guide: [Creating a Signing Certificate Using Jamf Pro's Built-in CA](https://learn.jamf.com/r/en-US/technical-articles/Creating_a_Signing_Certificate_Using_Jamf_Pros_Built-in_CA_to_Use_for_Signing_Configuration_Profiles_and_Packages).
{% endstep %}

{% step %}
**Add the Aikido Device Protection configuration profile**

1. In Jamf Pro, use the **upload** to add it as **Configuration Profile** for macOS.
2. Upload the signed `.mobileconfig` file.
3. Scope it to the target devices.

{% hint style="info" %}
If Jamf Pro shows an error in the **VPN** payload asking for a **Server** hostname or IP address, enter `127.0.0.1`. This is a known Jamf Pro limitation: it doesn't fully support the proxy payload used by Aikido Device Protection.
{% endhint %}

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In Jamf Pro, go to **Settings** → **Packages**.
2. Upload the Aikido Device Protection `.pkg`.
3. Confirm the package finishes processing before you use it in a policy.
   {% endstep %}

{% step %}
**Create the install policy**

1. In Jamf Pro, go to **Computers** → **Policies**.
2. Create a new policy scoped to the same devices.
3. Add the uploaded Aikido Device Protection package to the policy.
4. Set the policy's **Execution Frequency** to **Once per computer**.

{% hint style="warning" %}
Leaving the **Execution Frequency** at **Ongoing** makes the policy reinstall the package on every recurring check-in, which can redeploy an outdated build. The package should only be installed once per device.
{% endhint %}
{% endstep %}

{% step %}
**Deploy to your devices**

1. Deploy the **Aikido Device Protection** configuration profile.
2. Verify it has reached all target devices. Check profile status in the Jamf Pro console.
3. Deploy the **Aikido Device Protection** install policy.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In Jamf Pro, go to **Computers** → **Policies** and create a new policy.
2. Add a **Restart Options** payload and set it to **Restart**.
3. Scope it to the same target devices and trigger it after the agent package has installed.

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
