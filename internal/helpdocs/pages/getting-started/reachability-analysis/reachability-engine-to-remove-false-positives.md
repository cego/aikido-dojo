# Reachability Analysis Examples in Aikido

Aikido combines 100+ custom rules to reduce false positives and irrelevant alerts. Most of those ground rules are powered by our own reachability analysis engine.

When Aikido detects a known vulnerability in your dependencies or first-party code, it checks whether your application can actually call into the vulnerable code path. In practice, we build a lightweight call/dependency graph and trace references from your code to the functions/classes known to be affected. If the vulnerable code isn’t reachable (or is only used in non-production contexts like tests or tooling), the alert is downgraded or suppressed.

This results in significant noise reduction compared to legacy scanners that report every vulnerable package regardless of usage.

{% hint style="info" %}
Note: We continuously extend this mechanism across languages and ecosystems to further improve noise suppression.
{% endhint %}

#### **Example 1:** **You're not using the affected function**

Aikido’s internal knowledge base indicates that **CVE-2020-7774** only affects the `setLocale` function. If our analysis sees that your codebase never references `setLocale`, you’re not affected. Aikido will downgrade the issue and continue to monitor future changes in case the function starts being used.

[image: Security analysis found no usage of CVE-2020-7774 vulnerability in repository.]

#### **Example 2:** **Vulnerable package used only in tooling**

The JS package `minimatch` is vulnerable, but we detect it’s only pulled in by `eslint` and `mochajs`. Because these are linting/testing dependencies and are not shipped to production, your end-users are not exposed. Aikido downgrades the issue accordingly.

#### **Example 3:** Dead dependency path

A package (`path-parse`) is known to have a CVE and is required by `pug`. Our trace shows `pug` is no longer used by your application (no imports/references, or it’s excluded from your build). In this case, the CVE can be safely discarded.

[image: Dependency reachability analysis for CVE-2021-23343 showing vulnerable package path in repository.]

***

### How Aikido makes the call

* **Call/dependency graphing:** We map how your code imports and calls into packages, and then follow edges toward functions tied to a CVE.
* **Context awareness:** We distinguish production runtime code from dev and test tools (such as linting, test, or build code), so tooling-only usage won’t trigger production-impacting alerts
* **Guarded downgrades:** Findings are downgraded or suppressed only when the vulnerable symbol is provably unreachable. If evidence changes (such as with a new import), the severity is automatically reevaluted.

### Additional things to keep in mind

* **Heuristics and build setups:** Highly dynamic patterns, custom loaders, or unusual build steps can obscure call graphs. Some issues may be kept at a higher severity if usage is unclear.
* **Transitive changes:** Updating indirect dependencies can make previously unreachable code reachable. Aikido continually rechecks after dependency updates and code changes
