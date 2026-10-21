# Outbound Network Connections for Zen

Zen requires outbound internet access to function correctly.

During operation, Zen [sends data](/zen-firewall/miscellaneous/data-sent-by-aikido-zen.md) to the following Aikido services:

**EU-based domains:**

* guard.aikido.dev
* runtime.aikido.dev

**US-based domains:**

* guard.us.aikido.dev
* runtime.us.aikido.dev

**ME-based domains:**

* guard.me.aikido.dev
* runtime.me.aikido.dev

{% hint style="warning" %}
Both domains use dynamic IPs managed by Aikido’s cloud infrastructure.
{% endhint %}

If your environment restricts outbound traffic, allow these domains by hostname rather than IP address.

## Testing connectivity

Run these checks from the same host (and same container, if applicable) where your app runs. Substitute the regional domain that matches your account if you are not on the EU cluster.

**Linux / macOS**

```bash
# DNS resolution
nslookup guard.aikido.dev
nslookup runtime.aikido.dev

# HTTPS connectivity
curl -v https://guard.aikido.dev
curl -v https://runtime.aikido.dev
```

**Windows (PowerShell)**

```powershell
Resolve-DnsName guard.aikido.dev
Test-NetConnection guard.aikido.dev -Port 443
```

If you use a proxy, corporate firewall, VPN, or SSL interception, verify that outbound traffic to `*.aikido.dev` on port 443 is allowed and that Aikido's certificates are not being rewritten. Proxy configuration is covered in [Proxy & Load Balancer Settings](/zen-firewall/zen-installation-instructions/proxy-load-balancer-settings.md).
