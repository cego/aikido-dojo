# DAST Overview

Dynamic Application Security Testing (DAST) probes your running applications and infrastructure from the outside, the same way an attacker would. Aikido's DAST suite watches the domains you operate to surface misconfigurations, exposed assets and exploitable vulnerabilities before someone else finds them.

When you add a domain in Aikido, you choose the scan type that matches your use-case:

[**Quick Scan**](/dast-surface-monitoring/front-end-scanning/scan-front-end-app-domains-with-aikido.md): fast, non-intrusive checks for security best practices (CSP headers, cookies, JWT issues and more). A good starting point for any public-facing app.

[**Attack Surface Scan**](/dast-surface-monitoring/attack-surface-scanning.md): extensive external attack surface mapping, covering subdomains, leaked credentials, exposed servers, SSL configuration and more. Use this to understand everything you have exposed to the internet.

[**Agentic Scan**](/pentests/aikido-pentest.md): AI agents that simulate real attacks on APIs and web apps. Covers OWASP risks like SQL injection plus logic flaws, IDOR and cross-tenant data leaks. See [Aikido Pentest](/pentests/aikido-pentest.md) for the full setup.
