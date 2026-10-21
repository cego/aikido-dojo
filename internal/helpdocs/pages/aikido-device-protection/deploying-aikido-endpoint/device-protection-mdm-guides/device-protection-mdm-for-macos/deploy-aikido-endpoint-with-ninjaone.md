# Deploy Device Protection with NinjaOne

Use NinjaOne MDM to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in NinjaOne MDM (via manual profile install or Apple Business Manager / Apple School Manager Automated Device Enrollment) before you begin.
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

1. Open the `.mobileconfig` file in a text editor to copy its XML contents.
2. In the NinjaOne Admin Portal, click **Administration** in the left navigation pane, open the **Policies** drop-down, and select **MDM Policies**.
3. Click **Create New Policy** (or edit an existing macOS policy) and select **Apple** as the platform with **macOS** as the device type.
4. Give the policy a name (e.g. "Aikido Device Protection").
5. Open the **Custom payload** section and click **+ Add payload**.
6. In the **Add payload** dialog, enter a name (e.g. "Aikido Device Protection") and paste the full XML contents of the `.mobileconfig` file into the **Content** field.
7. Click **Update** to save the payload, then save the policy.
8. Assign the policy to the target macOS device group.

{% hint style="info" %}
NinjaOne deploys each custom payload as a separate MDM profile to the device. The Aikido `.mobileconfig` already contains all required payloads, so a single **Add payload** entry is enough.
{% endhint %}

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In the NinjaOne Admin Portal, go to **Administration** → **Library** → **Automation**.
2. Click **+ Add** and select **Installation**.
3. In the **Install Application** dialog, fill in the following:
   * **Name:** Aikido Device Protection
   * **Description:** Aikido Device Protection
   * **Operating system:** macOS
   * **Installer:** Upload the Aikido Device Protection `.pkg` file.
   * **Run as:** System
4. Save the installation package.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Assign the **Aikido Device Protection** MDM policy (containing the configuration profile) to the target device group.
2. Verify it is installed on a test device.
3. Deploy the **Aikido Device Protection** installation package to the same devices using an on-demand run or a **Run once** scheduled task.

{% hint style="warning" %}
When you schedule the deployment, use the **Run once** option (or set **Ends** to **After** and **Occurrences** to **1**). A recurring schedule reinstalls the package on every run, which can redeploy an outdated build. The package should only be installed once per device.
{% endhint %}
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
