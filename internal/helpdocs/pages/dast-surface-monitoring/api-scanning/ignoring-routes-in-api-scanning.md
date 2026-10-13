# Ignoring Routes in API Scanning

Our API scanning solution helps protect your endpoints by identifying potential security vulnerabilities. However, there may be times when you need to exclude specific routes from being scanned. This document explains when and how to ignore routes in the scanning process.

### How to Ignore a Route

Excluding a route from scanning is straightforward:

* Navigate to the Routes page of the domain you would like to modify
* Locate the specific route you wish to exclude
* Hover over the route to reveal the action menu on the right side
* Click on the action menu and select "Exclude from scan"

Once ignored, the route will be excluded from future security scans until you choose to re-enable it.

{% hint style="warning" %}
While ignoring routes can be necessary, remember that each excluded endpoint represents a potential security gap. Only ignore routes when there’s a clear justification, and regularly review your ignored routes list.
{% endhint %}

### Use Cases

There are several legitimate reasons to exclude certain routes from API scanning:

1. **Test or Development Routes**: Endpoints that are only used during development and aren’t exposed in production.
2. **High-Volume Endpoints**: Routes that cannot handle large amounts of traffic where scanning might impact performance.
3. **False Positives**: Routes that consistently trigger false security alerts due to their unique functionality.
4. **Internal Tools**: Admin or debugging endpoints that use different security models.
5. **Third-Party Integrations**: Routes that interface with external systems that have their own security measures.
