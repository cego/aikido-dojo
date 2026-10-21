# Composer - plugin-api-version changes

By default AutoFix is using Composer 2.6, when you're using a different composer version this can lead to unwanted updates of the `plugin-api-version` field in the lockfile:

To avoid these changes, you can specify the composer version that should be used by AutoFix. The composer version can be configured in `composer.json` by adding the following configuration:

{% code title="composer.json" %}

```
{
    "config": {
        "platform": {
            "composer": "2.9.0"
        }
    }
}
```

{% endcode %}
