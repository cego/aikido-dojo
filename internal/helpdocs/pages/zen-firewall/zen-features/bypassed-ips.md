# Bypassed IPs

Bypassed IPs are trusted IP addresses or CIDR ranges that Aikido Zen explicitly ignores during request processing. Any request coming from a bypassed IP is fully excluded from Zen inspection and enforcement.

This means Zen will not analyze the request, generate findings, or apply blocking rules for that traffic.

## Why use bypassed IPs

Not all traffic represents real attacker behavior. Some requests come from systems you trust and control.

Common examples include internal infrastructure such as load balancers or reverse proxies, monitoring and uptime services, internal testing environments, approved security tools, or Aikido scanning traffic.

Bypassing these IPs helps reduce noise, avoids false positives, and keeps alerts focused on real threats.

When running **AI pentesting or Domain and API security** scans, bypassing IPs can be useful depending on the type of scan you want to perform.

In some cases, you may want Zen fully enabled to observe detections and blocking behavior. In other cases, temporarily bypassing scan traffic prevents interference and allows uninterrupted testing.

## How to configure bypassed IPs

{% hint style="danger" %}
All Zen security features will be disabled for the configured IPs, so only add sources you fully trust.
{% endhint %}

1. Go to the Firewall tab in the Aikido dashboard
2. Scroll to the Danger zone section at the bottom of the page<br>

   
3. Select Manage Bypass List<br>

   
4. Add one or more IP addresses or CIDR ranges
5. Save your changes
