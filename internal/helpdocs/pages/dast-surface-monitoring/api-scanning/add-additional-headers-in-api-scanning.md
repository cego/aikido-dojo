# Add Additional Headers in API Scanning

Add custom headers that Aikido will include in every request during an API scan. This is useful for auth tokens, tenant IDs or any header your service needs before it accepts traffic.

Aikido sends each header you define with all scan requests. The scanner does not try to detect missing headers on its own, so make sure you include everything your API requires. Unlike our [Authentication](/dast-surface-monitoring/api-scanning/authenticated-api-scanning-for-rest-graphql.md) these headers are included during the scan.

### Typical use cases

* Static auth tokens for staging or test environments
* X API Key headers
* Version or tenant selectors
* Headers needed to bypass a gateway or WAF in front of your API

## Add a header

1. Click Edit headers in the API Settings page.

1. Enter the header name.
2. Enter its value.
3. Add more headers if needed.
4. Save changes.
