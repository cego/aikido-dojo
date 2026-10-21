# Attack Surface Monitoring

Attack Surface Monitoring (ASM) provides visibility in your exposure to external threats from the internet. Most ASM solutions focus on purely infrastructure related exposures but most often breaches happen due to identity related attacks. "Attackers don't 'hack' in, they log in". That's why the Aikido platform includes both infrastructure and identity related threats to cover the whole attack surface.

## Use Cases

* 🔍 **Attack Surface Monitoring**: Detect forgotten subdomains, exposed databases, outdated software, and more before attackers do.
* 🕵️ **Shadow IT Discovery**: Uncover unauthorised or forgotten assets that may not be under proper security management
* :tophat: **Darknet monitoring:** Identify already compromised credentials on the darknet.
* :office: **Boost Third Party Risk Management (TPRM) ratings:** Address key factors TPRM platforms use to assess your security posture and boost your chances of closing new deals.
* ✅ **Compliance Verification**: Validate that your external-facing systems meet industry standards and regulatory requirements.

## Infrastructure Attack Surface Monitoring

infrastructure monitoring provides visibility into your external security posture by continuously monitoring your internet-facing assets for vulnerabilities and exposures.

This proactive approach helps you identify and remediate security risks before they can be exploited.

### **Asset Inventory Management**

* Performs subdomain enumeration to map systems that are publicly accessible
* Resolves IP addresses and verifies reachability from the internet

{% hint style="info" %}
If AWS is connected in Aikido CSPM, the subdomains for a configured root domain are automatically enriched with Route 53 records and continuously monitored.\
\
**Note:** Root domains like insecure.org are not automatically added to the Attack Surface module to avoid unwanted traffic to your domains. The subdomains like web.insecure.org are imported.
{% endhint %}

### **Vulnerability Scanning**

* Leverages a vast library of security checks to identify known vulnerabilities
* Detects misconfigurations in web servers, applications, and infrastructure
* Identifies outdated software versions with known security issues

### **Subdomain Takeover Identification**

Detects misconfigurations in DNS records that could lead to subdomain takeovers

* Identifies vulnerable subdomains pointing to deprovisioned or unclaimed services
* Continuously monitors for dangling DNS records (Domain Dangling) that could be exploited by attackers
* Detects Route 53 `A` records pointing to released AWS Elastic IPs (EIPs) that an attacker could re-acquire

### Email Spoofing Protection

Identifies missing or misconfigured mail security records like SPF, DKIM and DMARC to validate email security settings and prevent spoofing of your domains in phishing attacks.

### **DNS Zone Transfer Attacks**

Identifies misconfigurations in the DNS nameservers that would allow for zone transfer attacks.

### **Port Exposure Detection**

Identifies open ports that are exposed to the internet and highlights management, file sharing ports or database ports.

### TLS Security Hardening

Identifies insecure or weak TLS configurations

* Detects weak or insecure TLS protocols in use like SSL 3.0, TLS 1.0 and TLS 1.1
* Discovers weak or insecure cipher suites that are configured to establish the encrypted connection
* Assesses the Post-Quantum readiness of your TLS configuration
* Highlights issues in certificates that could provide problems in verifying and establishing a chain of trust
* Monitors the health of the certificates used by your servers and alerts these are about to expire

### Supabase Potential Data Breach Validation

Identifies Supabase instances for which a misconfigured Row Level Security (RLS) policy results in a potential data breach. It works in the following steps:

1. **Key discovery**: Crawl deployed frontend JavaScript assets to extract the Supabase client anon key.
2. **Project identification**: Determine the Supabase project endpoint (e.g., https\://\<project>.[supabase.co](http://supabase.co/)) from the JWT token or code references.
3. **Credential validation**: Establish a connection to the project to verify that the extracted anon key is valid.
4. **Surface enumeration**: Enumerate tables within the public schema that are addressable through the client context.
5. **Permission testing**: Attempt read and write operations on each table using only the anon key to identify gaps caused by missing or misconfigured RLS.
6. **Sensitivity assessment**: For any table that is readable, analyze returned records for PII and other sensitive fields to determine exposure severity and remediation priority.

## Identity Exposure Monitoring

To defend against identity-related attacks, the platform detects exposed credentials on the dark web and safely tests for default, weak, and leaked passwords across your organisation’s login portals.

### **Darknet Monitoring**

* Discovers leaked usernames and passwords on the darknet
* Monitors code repositories, paste sites, and data breach collections
* Helps prevent account takeovers and unauthorized access

### Weak and Default Credentials Usage Detection (COMING SOON!)

Identifies the exposure of weak and default credentials on login portals related to your organisation.

### Leaked Credential Validation (COMING SOON!)

Validates the leaked credentials identified in previous steps to verify whether these still work on exposed login portals and whether these have been reused accross applications.

## How to Set-up

**Step 1:** Click **Add Domain** in the [Domain Overview](https://app.aikido.dev/domains) and select **Attack Surface Scan**

[image: Choose scan type modal with Quick Scan, Attack Surface Scan and Agentic Scan options.]

**Step 2.** Enter the domain name of your **environment**. Ensure this is a root domain (e.g., `example.com`). In case you only want to scan a subset of domains that are f.e. connected to a specific department or geography, you can also use `us.example.com` to only scan subdomains from `*.us.example.com`.

[image: Attack surface domain name input field]

**Step 3.** To prevent abuse, we require you to prove that you are the owner of the domain you'd like to scan. Configure the CNAME or TXT record in your DNS registrar and once done, click "Verify Domain".

{% hint style="info" %}
Domains are automatically marked as verified when the Attack Surface scan is configured by a user logged in with Google or Microsoft whose email address matches the target domain in scope.
{% endhint %}

**Step 4.** A scan will automatically start after completing the set-up.

## Identifying Traffic

All requests coming from Aikido Attack surface scans will have:

* the `User-Agent` set to `aikido-surface-agent/1.0`
* the following header `aikido-surface-test` set to value `1` in the request
* [will come from the IP's documented here](/dast-surface-monitoring/miscellaneous/allowing-ip-addresses-for-dast-surface-monitoring.md)
