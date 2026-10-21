# Allowing IP Addresses for Issue/Task Tracker Integrations

If you only allows specific IP addresses to access your Issue/Task Management Systems (eg Jira, Linear,..), you will have to allowlist Aikido's IP addresses so tasks can be created and managed.

Add the following IPs:

* **18.197.244.247**
* **18.156.9.3**
* **3.65.139.215**

The ports required to be opened are at least port **443** for HTTPS.

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
