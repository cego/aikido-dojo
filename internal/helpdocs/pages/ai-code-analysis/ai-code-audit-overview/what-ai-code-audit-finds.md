# What Code Security Audits Find

What a Code Security Audit checks for in your source code.

Code Security Audit reasons about your source code using the agentic engine [Aikido Pentest](/pentests/aikido-pentest.md) is built on. It covers the vulnerability classes Aikido Pentest validates, plus issues a pentest can't safely or reliably reach against a running app. The trade-off is that purely runtime concerns still need a live target. Each finding ships with code-based **Evidence**. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md) for how to read it.

## Coverage at a glance

Code Security Audit reasons about the issue classes Aikido Pentest validates, including:

* **Access control:** IDORs / BOLA, broken access control, and privilege escalation across tenants and roles.
* **Injection flaws:** SQL, NoSQL, LDAP, XPath, command injection, RCE, and server-side template injection (SSTI).
* **Cross-site and client-side attacks:** Stored and reflected XSS, CSRF, open redirects, and unsafe client-side data handling.
* **Authentication and session logic:** Weak flows, missing checks, session-handling errors, and hard-coded credentials in auth code.
* **SSRF, deserialization, and file handling:** Unsafe URL fetches, dangerous deserialization, path traversal, and unrestricted upload patterns.
* **Cryptographic failures:** Broken algorithms, weak signatures (e.g., JWT), and sensitive data exposure in code paths.
* **Business logic flaws:** Workflow bypasses, payment and checkout shortcuts, and multi-step exploit chains.
* **Agentic application risks:** Prompt injection, excessive agency, insecure tool use, and the rest of the [OWASP Top 10 for Agentic Applications](/pentests/coverage-and-findings/what-issues-can-aikido-pentest-find/owasp-top-10-for-agentic-applications.md).

For the full vulnerability catalog, see [What Issues Can Aikido Pentest Find?](/pentests/coverage-and-findings/what-issues-can-aikido-pentest-find.md). Every class listed there is in scope for Code Security Audit unless it appears in **What still needs a live target** below.

## Beyond what a live pentest reaches

Some issues are easier to identify from source, and others are only visible there. Code Security Audit reasons about them directly, even when a live pentest couldn't safely or feasibly reach them:

* **Code paths gated by credentials the pentest doesn't have.** If you don't provide an admin or manager account to a pentest, the admin panel won't be exercised. Code Security Audit still reads that code and reports issues like SQL injection in admin-only routes.
* **Denial-of-service and resource-exhaustion patterns.** Live pentests can't safely run these against your app. Code Security Audit flags algorithmic complexity attacks like ReDoS, hash-collision DoS, and exponential parsers directly from code.
* **Code paths that aren't externally reachable yet.** Newly merged routes that haven't deployed, feature-flagged code, and dormant branches are visible to source-level reasoning even when nothing at runtime exposes them.

## What still needs a live target

Some checks only make sense against a running application. Code Security Audit can flag risky configuration in source, but it can't confirm what your environment actually does on the wire. For these, use [Aikido Pentest](/pentests/aikido-pentest.md):

* **TLS and transport posture:** Cipher suites, protocol versions, certificate chains, and HSTS as served.
* **HTTP security headers as delivered:** Code Security Audit can see whether headers are set in code, but not what your CDN, WAF, or edge layer actually sends to a client.
* **Rate limiting and abuse protection in practice:** Code can show that limits exist; only a live run proves they hold up.
* **CORS behavior at the edge:** Configuration is visible, but real responses are not.
* **Web cache poisoning:** Requires interacting with a live cache.
* **Credential stuffing, brute force, and CAPTCHA bypass:** Live-only by definition.
* **DOM-based XSS as actually rendered:** Code Security Audit reasons about source; runtime rendering needs a browser.
* **Infrastructure on the host:** Confirmed against the running host, not from source.

## Application types and languages

Code Security Audit is especially valuable for codebases that are **hard or impractical to set up for live testing**. If wiring up a staging environment, simulator, or deployed instance is the blocker, source-level reasoning gets you pentest-grade findings without any of that overhead:

* **Mobile apps:** iOS and Android codebases that would otherwise need devices, simulators, or instrumentation to test.
* **Smart contracts:** Solidity and other on-chain code, where setting up a realistic exploit environment is non-trivial. For a review focused only on the contracts, use [Smart Contract Audit](/ai-code-analysis/smart-contract-audit.md).
* **Desktop apps:** Installed binaries with no public URL to point a pentest at.
* **Embedded software, libraries, and SDKs:** Code that ships into someone else's runtime, with no live target of your own.
* **Internal tools and back-office systems:** Code paths behind VPNs, SSO, or admin-only access that a live pentest typically can't reach.

Agents work across mainstream application languages, configuration, and infrastructure as code the same way a developer would. The analysis treats the whole repository as in scope, including monorepos with multiple services or packages, and follows references across files and modules to surface multi-step issues.

## FAQ

### Can you confirm coverage before I start?

Yes. [Contact support](/miscellaneous-info/contact-support.md) to check a specific framework, application type, or finding class before you spend credits.

## Related

{% content-ref url="/pages/kZRWK6NswIlL0dAmeNea" %}
[Code Security Audit](/ai-code-analysis/ai-code-audit-overview.md)
{% endcontent-ref %}
