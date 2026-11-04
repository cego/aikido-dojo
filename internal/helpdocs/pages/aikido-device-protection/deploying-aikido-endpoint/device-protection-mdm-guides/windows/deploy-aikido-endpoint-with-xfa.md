# Deploy Device Protection with XFA

Use XFA to deploy Aikido Device Protection to every Windows device your team uses for work. That includes personal (BYOD) and contractor devices that aren't enrolled in a mobile device management (MDM) solution.

{% hint style="info" %}
Your organization must have XFA verification at sign-in set up for its identity provider before you begin.
{% endhint %}

## How it works

1. The user signs in to a work app with your organization's identity provider.
2. XFA checks the device against your policy, including whether Device Protection is installed and running.
3. If the check fails, XFA shows the user what to fix and starts the Aikido installation with your organization's enrollment token.
4. Once Device Protection is active, the check passes and the user continues signing in. The device appears in the Aikido app.

## Prerequisites

* An XFA organization with verification at sign-in set up for your identity provider. See the [XFA docs](https://docs.xfa.tech).
* The XFA app on the devices you want to cover. Users who don't have it yet are prompted to install it at their first sign-in.
* Your Device Protection enrollment token. In the [Device Protection devices page](https://app.aikido.dev/endpoint-protection/devices), click **Connect Device** and copy the token from the installation instructions. It starts with `AIK_`.

## Set up Device Protection with XFA

{% stepper %}
{% step %}

### Add your Aikido token to XFA

1. In XFA, go to **Integrations** → **Aikido**.
2. Paste your Device Protection enrollment token.
3. Click **Save**.

XFA uses this token to enroll each device in your Aikido workspace during installation.
{% endstep %}

{% step %}

### Add the Aikido check to your policy

1. In XFA, go to **Policies** and open the policy that applies to your target users or groups.
2. Enable the **Aikido Endpoint** check for Windows.
3. Set the enforcement level and save the policy.

Enforcement levels:

* **Warn**: users are notified and guided to install Aikido, but can still sign in.
* **Block**: users can't sign in until Device Protection is installed and running.

{% hint style="info" %}
Start with **Warn** for a few days so users can install Aikido at a convenient moment, then switch to **Block**. You can also apply the policy to a small group first.
{% endhint %}

If you also use [Safe Chain](/code-scanning/aikido-malware-scanning.md), enable the **Aikido Safe Chain** check in the same policy.
{% endstep %}

{% step %}

### Let users install Device Protection at sign-in

No further action is needed from you. The next time a user in scope signs in:

1. XFA shows that Device Protection is missing and offers to install it.
2. The user confirms. XFA asks for administrator rights only for this installation, then runs the Aikido installer with your token.
3. XFA guides the user through the approval prompts the Aikido installer shows on Windows.
4. XFA rechecks the device. Once Aikido is active, the user continues to their app.

XFA requests administrator rights on demand, only when an installation needs them. After installation, XFA keeps checking that Aikido stays installed and running.
{% endstep %}

{% step %}

### Verify the deployment

On a test device, confirm:

1. The Aikido check shows as passing for the device in XFA under **Devices**.
2. The device appears with an *Active* status on the [Device Protection devices page](https://app.aikido.dev/endpoint-protection/devices).
3. Opening [device-protection.aikido.help](https://device-protection.aikido.help/) on the device shows that Device Protection is active.
   {% endstep %}
   {% endstepper %}

## Troubleshooting

| Problem                                                             | Fix                                                                                                |
| ------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| Aikido is installed but the device doesn't appear in the Aikido app | Check that the token saved in XFA matches the one in your Aikido workspace and hasn't been revoked |
| The XFA check still fails after installation                        | Restart the device so the agent fully activates, then sign in again                                |
| You want users to install Aikido without being blocked yet          | Set the Aikido check to **Warn** in the XFA policy                                                 |
