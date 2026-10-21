# Deploy Device Protection with Mosyle

Use Mosyle Business to deploy Aikido Device Protection across your managed macOS fleet with the required permissions in place.

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

1. In the Mosyle admin console, go to **Management** → **Manage** **macOS** → **Management Profiles (sidebar)**.<br>
2. Click "**Certificates / Custom Profiles"** and give it a name (e.g. "Aikido Device Protection").
3. Upload the `.mobileconfig` file via Select File.
4. Under **Add Assignment**, select your target device group with **System** scope.
5. Click **Save**.

This profile allows the Aikido Device Protection system extension to load silently, enables the network content filter, and stops users from disabling background services in System Settings → Login Items.
{% endstep %}

{% step %}
**Upload and deploy the Aikido Device Protection installer**

1. In the same view as Step 3. In sidebar under **Management Profiles** click **Install PKG**
2. Select **"Already have a .PKG"**
3. On "Host your .PKG file" select "**Automatically set App info**"
4. Use this public URL "<https://github.com/AikidoSec/safechain-internals/releases/latest/download/EndpointProtection.pkg>" and **enable "This app is Signed"** and **"Install as Managed"**. Continue with "Add Enterprise App"<br>
5. Give the profile a name (e.g. "Aikido Device Protection - Pkg").
6. Under **Add Assignment**, select the same target device group.
7. Click **Save**.
   {% endstep %}

{% step %}
**Deploy to your devices**

1. Confirm the **Aikido Device Protection** configuration profile is installed: open the profile and click **View Results**, all target devices should show green.
2. Save and assign the **Install PKG** deployment profile.
   {% endstep %}

{% step %}
**Reboot devices after installation**

Restart your target devices after confirming the agent package has installed. Check **View Results** on the Install PKG profile first. The agent fully activates on the next boot.
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
