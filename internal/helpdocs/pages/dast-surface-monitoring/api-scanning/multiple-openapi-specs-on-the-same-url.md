# Multiple OpenAPI Specs on the Same URL

Each REST API domain in Aikido holds one OpenAPI spec, and each domain URL must be unique. If several services sit behind the same base URL, each with its own spec, add the same URL once per spec and make each entry unique with a named URL fragment, such as `#orders` or `#billing`.

## How it works

A URL fragment is the part after `#`. Browsers and HTTP clients never send it to the server. Aikido treats `https://staging.example.com/api#orders` and `https://staging.example.com/api#billing` as 2 separate domains, but every scan request from both goes to `https://staging.example.com/api`.

| Domain in Aikido                          | OpenAPI spec           | Requests go to                    |
| ----------------------------------------- | ---------------------- | --------------------------------- |
| `https://staging.example.com/api#orders`  | `orders-openapi.json`  | `https://staging.example.com/api` |
| `https://staging.example.com/api#billing` | `billing-openapi.json` | `https://staging.example.com/api` |

## Add a spec for the same URL

1. Go to [Domains](https://app.aikido.dev/domains) and click **Add Domain**.
2. Expand **Advanced** and select **REST API**.
3. Enter the base URL with a fragment named after the service, e.g. `https://staging.example.com/api#orders`.
4. Add the OpenAPI spec for this service. See [REST API & Web App Scanning](/dast-surface-monitoring/api-scanning/rest-api-scanning.md) for the spec options.
5. Repeat with a new fragment name, e.g. `#billing`, for each additional spec.

{% hint style="info" %}
Authentication and [additional headers](/dast-surface-monitoring/api-scanning/add-additional-headers-in-api-scanning.md) are set per domain. Configure them on every fragment domain that needs them.
{% endhint %}

## Related

{% content-ref url="/pages/omv5AOwmFHmKiGovLhvo" %}
[REST API & Web App Scanning](/dast-surface-monitoring/api-scanning/rest-api-scanning.md)
{% endcontent-ref %}

{% content-ref url="/pages/2zRQBTOpriFlX6nj9JgW" %}
[Authenticated API Scanning for REST/GraphQL](/dast-surface-monitoring/api-scanning/authenticated-api-scanning-for-rest-graphql.md)
{% endcontent-ref %}
