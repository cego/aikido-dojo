# Disable Attack Protection for Specific Path or URL

Disable Attack Protection lets you mark a route as safe if Zen is incorrectly blocking legitimate traffic. This is the recommended way to handle false positives. Once disabled, Zen stops blocking requests on that route (or group of routes)

## Use Cases

* Database tools, for example PHPMyAdmin: Routes where raw SQL queries are expected. Zen would normally block these, but here they’re legitimate.
* Search APIs with custom query syntax: For example, an API that accepts filter expressions (field:value) that look like injection attempts.
* Developer / debug endpoints: Internal tools that accept large or unusual payloads.
* Custom applications: Endpoints intentionally designed to accept input patterns that resemble attacks.

{% hint style="warning" %}
Only disable protection if you are sure the behavior is expected and safe. For real attacks, Zen should stay enabled.
{% endhint %}

## How to Disable Protection

### Disable for a Single Route

1. Open the Routes view in your Zen dashboard.
2. Locate the route that triggered the false positive.
3. Click the action menu (⋮).
4. Select **Disable Protection** for this route.
5. Confirm your choice.

Zen will now no longer block attacks on this specific route.

### Disable for Multiple Routes (Wildcard)

1. Open the Routes view in your Zen dashboard.
2. Add a new route with wildcard `*`<br>

   
3. Click the action menu (⋮)
4. Select Disable protection for route wildcard.
5. Confirm your choice.

All routes that match the wildcard will bypass attack blocking.
