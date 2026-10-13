# IP Addresses for Domain Scanning

Aikido uses dedicated IP addresses to perform scanning of your domains (DAST). To prevent connectivity issues, rate limiting, or security blocks, add these IPs to your firewall’s allowlist or other security software. After this, rescan your domains to confirm connectivity.

Looking for Aikido Pentest IP addresses? They are listed on [IP Addresses for Pentest](/pentests/prepare-a-pentest/ip-addresses-for-pentest.md).

{% hint style="warning" %}
To use 'Fetch OpenAPI by URL', you must also add the [Code & Container scanning IP addresses.](/code-scanning/miscellaneous/allowing-ip-addresses-for-code-container-scanning.md)
{% endhint %}

## Which IPs should you allowlist?

Domain scanning traffic comes from the region where **your Aikido workspace** is hosted, not from the region where your application runs. Check the URL you use to log in to Aikido, then allowlist that region's IPs:

| Aikido URL                                     | Region | IP list               |
| ---------------------------------------------- | ------ | --------------------- |
| [app.aikido.dev](https://app.aikido.dev)       | EU     | EU-based IP addresses |
| [app.us.aikido.dev](https://app.us.aikido.dev) | US     | US-based IP addresses |
| [app.au.aikido.dev](https://app.au.aikido.dev) | AU     | AU-based IP addresses |
| [app.me.aikido.dev](https://app.me.aikido.dev) | ME     | ME-based IP addresses |

{% hint style="warning" %}
If you log in at `app.aikido.dev`, scanning uses **EU** IPs, even when your endpoints are hosted in the US. Allowlisting only US IPs in that case will cause scans to fail.
{% endhint %}

## IP Addresses

### EU-based IP addresses (`app.aikido.dev`)

* 3.248.4.169
* 54.76.211.68
* 54.228.156.63
* 54.247.155.164
* 18.200.152.99
* 18.202.99.112

### US-based IP addresses (`app.us.aikido.dev`)

* 98.85.190.95
* 52.204.144.1
* 44.209.56.130
* 18.210.114.117
* 35.168.38.209
* 35.173.56.162
* 54.227.161.94
* 44.209.154.183

### ME-based IP addresses (`app.me.aikido.dev`)

* 158.252.118.40
* 158.252.52.197
* 40.172.160.56

### AU-based IP addresses (`app.au.aikido.dev`)

* 54.66.103.236

### Optional IP addresses

Used for troubleshooting with support.

* 79.127.239.171

## Request Headers

All HTTP requests originating from Aikido scans (except for Front-end scans) include a specific User-Agent header. This can be used to identify the DAST traffic or for allowlisting purposes:

* `User-Agent: aikido-scan-agent/1.0`

For instructions on whitelisting IP addresses with third-party providers, refer to the following resources:

* [Cloudflare WAF](https://developers.cloudflare.com/waf/custom-rules/use-cases/allow-traffic-from-ips-in-allowlist/)
  * Cloudflare Turnstile does not support allowlisting specific client IP addresses. If you need to [bypass Turnstile for Aikido scanning traffic, you must do it in your application code.](https://developers.cloudflare.com/turnstile/tutorials/conditionally-enforcing-turnstile/) We recommend bypassing only when both conditions are true:
    1. The request originates from an Aikido IP range
    2. The request includes the `aikido` User Agent in headers as described above
* [Azure WAF](https://learn.microsoft.com/en-us/azure/web-application-firewall/ag/custom-waf-rules-overview)
* [AWS WAF](https://docs.aws.amazon.com/waf/latest/developerguide/waf-rule-statement-type-ipset-match.html).
  * For WAFs behind Application Load Balancers or CloudFront, your [WAF should check the last IP address in the `X-Forwarded-For` header](https://docs.aws.amazon.com/waf/latest/developerguide/waf-rule-statement-forwarded-ip-address.html).
* [Vercel WAF](https://vercel.com/docs/vercel-firewall/vercel-waf/custom-rules)
  * Use the ["bypass" action](https://vercel.com/docs/vercel-firewall/firewall-concepts#bypass) for trusted IPs

{% hint style="info" %}
[The IP address lists are also available as JSON arrays](https://aikido.help/ips/)
{% endhint %}

***
