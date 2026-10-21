# Allowing IP Addresses for Code and Container Scanning

If your Git provider only allows specific static IP addresses to access your code, you will have to whitelist Aikido's static IP addresses for code scanning before you can start scanning.

## Which IPs should you allowlist?

Code and container scanning traffic comes from the region where **your Aikido workspace** is hosted, not from the region where your Git provider or registry runs. Check the URL you use to log in to Aikido, then allowlist that region's IPs:

| Aikido URL                                     | Region | IP list               |
| ---------------------------------------------- | ------ | --------------------- |
| [app.aikido.dev](https://app.aikido.dev)       | EU     | EU-based IP addresses |
| [app.us.aikido.dev](https://app.us.aikido.dev) | US     | US-based IP addresses |
| [app.au.aikido.dev](https://app.au.aikido.dev) | AU     | AU-based IP addresses |
| [app.me.aikido.dev](https://app.me.aikido.dev) | ME     | ME-based IP addresses |

{% hint style="warning" %}
If you log in at `app.aikido.dev`, scanning uses **EU** IPs, even when your Git provider or registry is hosted in the US. Allowlisting only US IPs in that case will cause scans to fail.
{% endhint %}

## IP Addresses

### EU-based IP addresses (`app.aikido.dev`)

* 52.214.244.18
* 18.202.209.180
* 52.50.198.227
* 52.51.98.186

### US-based IP addresses (`app.us.aikido.dev`)

* 3.211.221.73
* 54.163.131.24
* 54.225.143.47

### ME-based IP addresses (`app.me.aikido.dev`)

* 3.29.13.194
* 40.172.18.244
* 40.172.67.79

### AU-based IP addresses (`app.au.aikido.dev`)

* 3.24.60.82
* 52.65.15.49

The ports required to be opened are at least port **443** for HTTPS. For Docker container registries, additional ports might be required. For example, the Gitlab Container Registry requires port **4567** to be open.

After adding the IPs, rescan your repositories to confirm connectivity.

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
