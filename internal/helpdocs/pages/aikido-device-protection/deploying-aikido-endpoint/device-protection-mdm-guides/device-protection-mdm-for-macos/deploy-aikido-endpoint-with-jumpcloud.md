# Deploy Device Protection with JumpCloud

Use JumpCloud to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
All devices must be enrolled in [JumpCloud Apple MDM](https://jumpcloud.com/support/set-up-apple-mdm) before you begin.
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

1. In the JumpCloud Admin Portal, go to **Device Management** → **Policy Management**.
2. Click **(+)**, select the **Mac** tab, and choose **MDM Custom Configuration Profile**. Click **configure**.
3. Give the policy a name (e.g. "Aikido Device Protection").
4. Under **Settings**, click **upload file** and select the `.mobileconfig` file.
5. Select the **Device Groups** tab and assign the policy to the target device group.
6. Click **save**.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In the JumpCloud Admin Portal, go to **Device Management** → **Software Management** → **Apple**.
2. Click **(+Add New)** and select **JumpCloud Private Repo**.
3. Enter an **Application Name** (e.g. "Aikido Device Protection").
4. Under **Upload File**, upload the Aikido Device Protection `.pkg` file and click **Upload**.
5. Once processing completes, select the **Device Groups** tab and assign the app to the same device group.
6. Click **save**, then confirm the install.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Assign and save the **Aikido Device Protection** configuration policy.
2. Verify it is installed on a test device.
3. Deploy the **Aikido Device Protection** app.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In the JumpCloud Admin Portal, go to **Device Management** → **Commands**.
2. Click **+ Command**, select **MDM Command**, and choose **Restart Device**.
3. Assign it to the same device group and click **Run Now**.

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
