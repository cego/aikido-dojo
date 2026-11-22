# Node.js (Restify)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

[https://github.com/AikidoSec/firewall-node](https://github.com/AikidoSec/firewall-node)

## Requirements

* Node.js 16+.
* Restify 8.x to 11.x.
* [Aikido account](https://help.aikido.dev/getting-started/setting-up-your-account) & [Zen Firewall token](https://help.aikido.dev/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token)

## Installation & Configuration

{% stepper %}
{% step %}
**Install Zen Firewall by Aikido**

Install Zen in your project:

{% tabs %}
{% tab title="npm" %}
```bash
npm install --save-exact @aikidosec/firewall
```
{% endtab %}

{% tab title="Yarn" %}
```bash
yarn add --exact @aikidosec/firewall
```
{% endtab %}

{% tab title="pnpm" %}
```bash
pnpm add --save-exact @aikidosec/firewall
```
{% endtab %}
{% endtabs %}

{% tabs %}
{% tab title="require Syntax" %}
Load Zen **before** any other imports:

```code-highlight language="javascript"
+require("@aikidosec/firewall");

const restify = require("restify");
const server = restify.createServer();
```
{% endtab %}

{% tab title="import Syntax" %}
Import Zen **before** any other library:

```code-highlight language="typescript"
+import "@aikidosec/firewall";

import restify from "restify";
const server = restify.createServer();
```

{% hint style="warning" %}
Many TypeScript projects use `import` syntax but still compile to CommonJS. In that case, the setup above works as-is. If your app runs as **native ESM** (e.g. `"type": "module"` in package.json), see [**ESM setup**](https://help.aikido.dev/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript/node.js-esm) for additional steps.
{% endhint %}
{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}
**Enable Request Blocking and User Identification**

Use this middleware to enable rate limiting, user identification, and blocking features.

Zen Firewall does not require this middleware to block attacks. Core attack protection works without it. The middleware provides additional request context that Zen uses for protections such as [rate limiting](https://help.aikido.dev/zen-firewall/zen-features/setting-up-rate-limiting-for-routes), [user blocking](https://help.aikido.dev/zen-firewall/zen-features/blocking-users-with-zen-firewall), [bot blocking](https://help.aikido.dev/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall), [country blocking](https://help.aikido.dev/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall), and [threat actor](https://help.aikido.dev/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall) blocking.

Adapt the example to fit how your application identifies users and handles requests.

```code-highlight language="javascript"
+const Zen = require("@aikidosec/firewall");
+
+// Optional middleware to set user for rate limiting or blocking specific users
+server.use((req, res, next) => {
+ Zen.setUser({ id: "123", name: "John Doe" });
+ return next();
+});
+
+// Call this after auth middleware, as early as possible in the middleware stack
+Zen.addRestifyMiddleware(server);
```
{% endstep %}

{% step %}
**Optionally Enable IDOR protection**

Stop users from accessing records that are not theirs. IDOR protection helps catch unauthorized access to resources like invoices, tickets, or user profiles with minimal setup.

{% content-ref url="../../zen-firewall/zen-features/idor-protection-with-zen-firewall.md" %}
[idor-protection-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/idor-protection-with-zen-firewall)
{% endcontent-ref %}
{% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

```bash
AIKIDO_BLOCK=false AIKIDO_TOKEN=AIK_RUNTIME_ node app.js
```

Set the token as an environment variable so the Aikido Zen agent can pick it up. If you don't have a token yet, follow [instructions here](https://help.aikido.dev/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token).

```bash
AIKIDO_TOKEN=AIK_RUNTIME_
```

We recommend to start your app in dry mode  to ensure it works as expected without blocking any requests. We advise running Zen Firewall in staging for two weeks to avoid false positives.

```bash
AIKIDO_BLOCK=false
```

{% hint style="info" %}
You can use `AIKIDO_DEBUG=true` to enable debug mode for more detailed information about what the agent is doing. For more information about your environment variables: [configuration-via-environment-variables.md](https://help.aikido.dev/zen-firewall/zen-installation-instructions/configuration-via-environment-variables "mention")
{% endhint %}
{% endstep %}

{% step %}
**Test your app**

Browse to your application and perform a couple of actions or open a couple of pages. Zen will automatically discover the routes in your application.&#x20;

{% hint style="info" %}
Zen sends data back to Aikido every 10 minutes
{% endhint %}

You can verify a working agent by looking at the following pages of your Zen application:

* **Events**: Should show an "Application started" event.
* **Routes**: After some time your application routes will start showing here with the method, route and requests.
* **Instances**: Should show the number of active instances for your application where Zen is installed.

{% endstep %}

{% step %}
**Setup rate limiting in the dashboard**

When you've added the Zen Firewall middleware you can test protecting a route from brute force attacks, you do this by setting up rate limit in the Aikido Dashboard:

1. Click on the created app.
2. Go to the **Routes** tab.
3. Find the route you would like to limit and click **Setup rate limiting**.
4. Follow the instructions to configure the rate limit (e.g., 5 requests per minute).

[image: API route management interface showing authentication routes with protection and rate limiting options.]

[image: Set rate limiting for POST /auth/login to 5 requests per minute.]

**Verify Rate Limiting**

Start your app and try to access the route you've rate limited 5 times within a minute. After the fifth attempt, you should receive a rate limit error:

```
You are rate limited by Aikido firewall. (Your IP: 1.2.3.4)
```
{% endstep %}

{% step %}
**Next steps**

Congrats you've successfully installed Zen Firewall. If you encountered any problems, have concerns or feature requests, don't hesitate to reach out to support.

You can now go and explore the many features that Zen Firewall provides:

* [blocking-bot-traffic-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall "mention")
* [blocking-tor-traffic-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-tor-traffic-with-zen-firewall "mention")
* [blocking-users-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-users-with-zen-firewall "mention")
* [blocking-known-threat-actors-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall "mention")
* [blocking-traffic-by-country-with-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall "mention")
* [setting-up-rate-limiting-for-routes.md](https://help.aikido.dev/zen-firewall/zen-features/setting-up-rate-limiting-for-routes "mention")
* [monitor-outbound-domains.md](https://help.aikido.dev/zen-firewall/zen-features/monitor-outbound-domains "mention")

Additional information:

* [how-zen-works-performance-reliability.md](https://help.aikido.dev/zen-firewall/miscellaneous/how-zen-works-performance-reliability "mention")
* [blocking-vs-detection-mode-in-zen-firewall.md](https://help.aikido.dev/zen-firewall/zen-features/blocking-vs-detection-mode-in-zen-firewall "mention")
* [understanding-your-zen-statistics.md](https://help.aikido.dev/zen-firewall/zen-features/understanding-your-zen-statistics "mention")
{% endstep %}
{% endstepper %}
