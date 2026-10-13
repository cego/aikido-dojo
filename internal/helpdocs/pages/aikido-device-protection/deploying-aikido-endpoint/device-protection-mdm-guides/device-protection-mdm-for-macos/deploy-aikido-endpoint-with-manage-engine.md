# Deploy Device Protection with ManageEngine

Use ManageEngine Endpoint Central to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In Endpoint Central, go to **Configurations** → **Mac** → **Custom Configuration**.
2. Click **Upload Configuration** and select the `.mobileconfig` file.
3. Give the profile a name (e.g. "Aikido Device Protection") and click **Save and Publish**.
4. Associate the profile with the target devices or device groups.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Create the Aikido Device Protection package**

1. In Endpoint Central, go to **Software Deployment** → **Package Creation** → **Packages** → **Add Package** → **Mac**.
2. Upload the Aikido Device Protection `.pkg` file and fill in the package details.
3. Save the package.
   {% endstep %}

{% step %}
**Create the install deployment**

1. In Endpoint Central, go to **Software Deployment** → **Install/Uninstall Software**.
2. Create a new deployment task, select **Mac**, and choose the Aikido Device Protection package you just created.
3. Set the action to **Install** and scope it to the same target devices or groups.
4. Under **Execution Settings**, leave **Apply this configuration during every Refresh Cycle** unchecked so the package deploys once per device.
5. Configure the schedule and save the task.

{% hint style="warning" %}
Enabling **Apply this configuration during every Refresh Cycle** makes Endpoint Central redeploy the package every 90-minute refresh cycle, which can reinstall an outdated build. The package should only be installed once per device.
{% endhint %}
{% endstep %}

{% step %}
**Deploy to your devices**

1. Verify the **Aikido Device Protection** configuration profile has reached all target devices. Check status under **Configurations** → **Configuration Summary**.
2. Deploy the **Aikido Device Protection** install task.
   {% endstep %}

{% step %}
**Reboot devices after installation**

1. In Endpoint Central, go to **Tools** → **Remote Shutdown**.
2. Select the target devices from the list.
3. Click **More Actions** and choose **Restart**.

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
