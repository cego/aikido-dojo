# Device Protection & Bitdefender for Windows

Bitdefender Content Control rules, such as blocking YouTube, stop working on Windows devices that run Device Protection. To block websites on these devices, use the Aikido Device Firewall instead.

## How it works

On Windows, Device Protection relays the device's web traffic through its local proxy. Bitdefender then sees every connection as coming from the Aikido proxy instead of from the browser. Bitdefender Content Control doesn't support traffic that arrives through a proxy, so its rules no longer apply. When you pause Device Protection, the traffic stops going through the proxy and Content Control works again.

{% hint style="warning" %}
Adding the domain as an exception in Device Protection doesn't fix this. An exception only skips inspection: the traffic still goes through the proxy, so Bitdefender still can't see it.
{% endhint %}

## Block websites with the Device Firewall

Block the domains in the Aikido Device Firewall instead of in Bitdefender Content Control. For example, to block YouTube:

1. Go to [Device Protection and open the Firewall tab](https://app.aikido.dev/endpoint-protection/firewall?regionRedirect=1).
2. Open the **Exceptions** tab and click **Add Exception**.
3. Enter `youtube.com`, set the status to **Blocked**, choose the **User Groups**, and save.
4. Repeat for `*.youtube.com` to cover every subdomain.

For all options, see [Add an exception](https://help.aikido.dev/aikido-device-protection/using-aikido-endpoint/firewall#add-an-exception) in the Device Firewall guide.
