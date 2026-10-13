# Deploy Device Protection with Miradore

Use Miradore to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

{% hint style="info" %}
Deploying packages and scripts requires a Miradore **Premium** or **Premium+** subscription.
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

1. In the Miradore admin console, go to **Management** → **Configuration profiles**.
2. Click **Add** → **macOS** → **Advanced (custom)**.
3. Upload the `.mobileconfig` file.
4. Click **Save**.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload the Aikido Device Protection installer**

1. In the Miradore admin console, go to **Management** → **Applications**.
2. Click **Add** → **macOS application** → **PKG (Uploaded)** → **Next**.
3. Click **Select file** and upload the Aikido Device Protection `.pkg`.
4. Fill in the required fields:
   * **Application name** (e.g. "Aikido Device Protection")
   * **Bundle identifier** (e.g. `dev.aikido.endpoint`)
   * **Version**
5. Click **Create** and wait for the package to finish processing before continuing.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Restart your target devices after the agent installation completes. The agent fully activates on the next boot.
{% endstep %}

{% step %}
**Deploy using a Business Policy**

Use a Business Policy to deploy all three items.

1. In the Miradore admin console, go to **Management** → **Business policies**.
2. Click **Add**, give the policy a name, and set its status to **Disabled**.
3. Add all three items to the policy: the configuration profile, the `.pkg`, and the restart script.
4. Select the restart script, click **Actions** → **Edit item dependency**, and set the `.pkg` as its dependency so devices only restart once the agent has installed.
5. Assign the policy to your target devices using tags.
6. Set the policy status to **Enabled**.
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
3. To review deployment status, go to **Management** → **Action log** and filter by "Business policy" in the **Sender** field.
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
