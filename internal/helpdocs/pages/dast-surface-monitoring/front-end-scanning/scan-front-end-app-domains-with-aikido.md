# Scan Front-End App Domains with Aikido

The **Quick Scan** is a fast, non-intrusive check of your front-end app for security best practices: CSP headers, cookies, JWT issues, exposed libraries with known CVEs and more. It's the default starting point for any public-facing domain.

Quick Scan falls under Dynamic Application Security Testing (DAST): it inspects the externally-facing components of your software (web pages, headers, cookies, transferred data, third-party libraries) without ever touching your source code.

{% hint style="info" %}
Looking for deeper testing of APIs, business logic or authenticated flows? Use **Agentic Scan** instead. See [Aikido Pentest](/pentests/aikido-pentest.md).
{% endhint %}

## Use Cases

* 🔍 **Security Best Practices**: Verify your web app respects HTTPS, HSTS, CSP, X-Frame-Options and other web security baselines.
* 🔒 **Authentication Hygiene**: Catch JWT issues, weak cookie attributes (Secure, HttpOnly, SameSite) and other auth misconfigurations before they're exploited.
* 📦 **Vulnerable Library Detection**: Identify front-end libraries loaded on your pages that have known CVEs.
* 🔁 **Continuous Coverage**: Re-scan automatically so regressions introduced by new deploys are caught quickly.
* ✅ **Compliance Verification**: Demonstrate that your external-facing apps meet basic web-security hygiene required by standards and customers.

## What Quick Scan Checks

### Overview of checks performed

To see the checks performed by the Quick Scan, visit our [checks overview page](https://app.aikido.dev/domains/checks). Here, you'll find a detailed list of all the checks performed during the scan. Aikido will only perform safe, non-destructive automated tests (eg no automated SQL injection attempts,..).

### Overview of libraries detected

To see a list of all libraries detected by the Quick Scan, visit the **Resources** tab for a specific domain. You will find the library, version, date and time of detection and url for each detected resource. When CVEs are found for actively used resources, they will be added to your feed like other vulnerabilities.

## How to Set-up <a href="#add-a-domain-to-be-scanned-with-aikido" id="add-a-domain-to-be-scanned-with-aikido"></a>

**Step 1:** Navigate to the [Domains Overview Page](https://app.aikido.dev/domains) or [Domains Settings](https://app.aikido.dev/settings/domains), click **Add Domain** and select **Quick Scan**

[image: Choose scan type modal with Quick Scan, Attack Surface Scan and Agentic Scan options.]

**Step 2:** Fill in the service URL for the repositories which have public-facing domains by filling out the configuration form. You can specify full paths.

**Step 3: Optional**: link your domain to a repository or container

**Step 4: Optional:** set the sensitivity of the data

[image: Form to add domain details, link assets, and classify data sensitivity for issue scoring.]

**Step 5: Optional:** set up authentication if your domain requires login. Open the action menu for the domain and select **Authenticate Domain**. For the full list of supported authentication methods, see [Authenticated Scanning for Front-End Apps](/dast-surface-monitoring/front-end-scanning/authenticated-scanning-for-front-end-apps.md).
