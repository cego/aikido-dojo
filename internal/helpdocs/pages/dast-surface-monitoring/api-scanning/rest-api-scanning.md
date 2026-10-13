# REST API & Web App Scanning

{% hint style="warning" %}
**Agentic Scan is the recommended option for testing REST APIs and web apps.** It covers the same OWASP risks as the REST scanner (SQL injection, path traversal, shell injection) plus business-logic flaws, IDOR and cross-tenant data leaks that fuzzing cannot find. See [Aikido Pentest](/pentests/aikido-pentest.md).
{% endhint %}

Aikido can scan your REST API and Web App endpoints to uncover critical endpoint vulnerabilities, such as SQL injection or path traversal. Aikido uses API fuzzing, which essentially includes spamming dangerous payloads to each field in your API.

{% hint style="danger" %}
NEVER run on a production environment, always on staging to avoid potential downtime
{% endhint %}

### Main Use Cases <a href="#main-use-cases" id="main-use-cases"></a>

**Critical Vulnerability Detection**:

* SQL injection
* NoSQL injection
* Path traversal
* Shell injection

You can see all checks in the [Aikido app here](https://app.aikido.dev/domains/checks?scanner=rest).

### Setting up REST API Scanning <a href="#setting-up-rest-api-scanning" id="setting-up-rest-api-scanning"></a>

**Step 1:** Click **Add Domain** in the [Domain Overview](https://app.aikido.dev/domains), expand **Advanced** in the scan-type chooser and select **REST API**

[image: Choose scan type modal with the Advanced section expanded showing REST API and GraphQL options.]

**Step 2.** Enter the domain name of your **staging environment**. Ensure this is the base URL for your REST APIs (e.g., `https://example.io/api`)

[image: API domain name input field for application configuration.]

**Step 3:** Add your OpenAPI specification using one of these options:

* [**Connect to Zen App (recommended)**](/zen-firewall/getting-started-with-zen-firewall.md): Integrate with Zen to automatically discover and update API endpoints for continuous scanning. More info about Zen [can be found here](https://help.aikido.dev/section/zen-by-aikido/sgIt4HRxlrFr). No manual work or maintenance needed!
* [**Generate via Aikido AI (in beta)**](/dast-surface-monitoring/api-scanning/autogenerate-openapi-via-aikido-ai-code2swagger.md)**:** Using your codebase Aikido will generate an OpenAPI spec. No manual work is needed. Regular rescans keep everything current.
* **Fetch from URL:** Provide a URL that has the latest version of your OpenAPI spec. Aikido will fetch the spec before each scan.
* **Manual Upload**: Upload a OpenAPI file to define your API endpoints. You will be required to manually update and upload your spec each time new API endpoints are added or changed.

[image: Add or upload an OpenAPI spec to scan for security risks; staging use only.]

**Step 4:** Add authorization information to your API to make sure Aikido can access endpoints that require login. You can do this by clicking the triple dots action menu on the domain, and then '**Authenticate Domain**'

[image: Domain action menu with options: Start scan, Edit, Authenticate, and Delete domain.]

This will trigger the modal where you can fill in the authentication details. For the full list of supported authentication methods, see [Authenticated API Scanning for REST/GraphQL](/dast-surface-monitoring/api-scanning/authenticated-api-scanning-for-rest-graphql.md).

[image: Domain authentication setup form with input fields for API URL, email, and password.]

### Identifying Aikido traffic

All requests coming from Aikido REST and GraphQL scans will have:

* the `User-Agent` set to `aikido-scan-agent/1.0`
* the following header `aikido-api-test` set to value `1` in the request
* [will come from the IP's documented here](/dast-surface-monitoring/miscellaneous/allowing-ip-addresses-for-dast-surface-monitoring.md)
