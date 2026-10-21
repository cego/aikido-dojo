# Nodejs (Next.js standalone)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

<https://github.com/AikidoSec/firewall-node>

## Requirements

* Node.js 16+.
* Next.js 12.x, 13.x, or 14.x (standalone output).
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

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

Enable standalone output in `next.config.js`:

{% @aikido-custom-code/code-highlight language="javascript" content="const nextConfig = {

* output: "standalone",
  };" %}

Copy the firewall package to the standalone output after build:

```bash
mkdir -p .next/standalone/node_modules/@aikidosec
cp -r node_modules/@aikidosec/firewall .next/standalone/node_modules/@aikidosec
```

Run the standalone server with Zen preloaded:

```bash
node -r @aikidosec/firewall/instrument .next/standalone/server.js
```

{% hint style="info" %}
Please note that this requires a recent version of Node.js, we recommend to use Node.js 24 or higher. If you are using an older version of Node.js and you are still using Next.js 14, you can start the server with the following command: `node -r @aikidosec/firewall .next/standalone/server.js`. This method is no longer recommended for newer versions of Next.js, so please update to the latest version of Node.js.
{% endhint %}
{% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

```bash
AIKIDO_BLOCK=false AIKIDO_TOKEN=AIK_RUNTIME_ node -r @aikidosec/firewall/instrument .next/standalone/server.js
```

Set the token as an environment variable so the Aikido Zen agent can pick it up. If you don't have a token yet, follow [instructions here](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md).

```bash
AIKIDO_TOKEN=AIK_RUNTIME_
```

We recommend to start your app in dry mode to ensure it works as expected without blocking any requests. We advise running Zen Firewall in staging for two weeks to avoid false positives.

```bash
AIKIDO_BLOCK=false
```

{% hint style="info" %}
You can use `AIKIDO_DEBUG=true` to enable debug mode for more detailed information about what the agent is doing. For more information about your environment variables: [Configuration via Environment Variables](/zen-firewall/zen-installation-instructions/configuration-via-environment-variables.md)
{% endhint %}
{% endstep %}

{% step %}
**Optionally Enable IDOR protection**

Stop users from accessing records that are not theirs. IDOR protection helps catch unauthorized access to resources like invoices, tickets, or user profiles with minimal setup.

{% content-ref url="/pages/FSSao8N6iNSgll2UKqK2" %}
[IDOR Protection with Zen Firewall](/zen-firewall/zen-features/idor-protection-with-zen-firewall.md)
{% endcontent-ref %}
{% endstep %}

{% step %}
**Test your app**

Browse to your application and perform a couple of actions or open a couple of pages. Zen will automatically discover the routes in your application.

{% hint style="info" %}
Zen sends data back to Aikido every 10 minutes
{% endhint %}

You can verify a working agent by looking at the following pages of your Zen application:

* **Events**: Should show an "Application started" event.
* **Routes**: After some time your application routes will start showing here with the method, route and requests.
* **Instances**: Should show the number of active instances for your application where Zen is installed.

{% endstep %}

{% step %}
**Framework caveats**

* User-based rate limiting and manual user blocking are not supported because `setUser` is not supported in this setup.
* Detection coverage depends on [modules being externalized by Next.js](https://nextjs.org/docs/app/api-reference/config/next-config-js/serverExternalPackages) bundling behavior.
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

* [Blocking Bot traffic with Zen Firewall](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md)
* [Blocking or Monitoring Tor traffic with Zen Firewall](/zen-firewall/zen-features/blocking-tor-traffic-with-zen-firewall.md)
* [Tracking Users with Zen Firewall](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md)
* [Blocking Known Threat Actors with Zen Firewall](/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall.md)
* [Blocking Traffic by Country with Zen Firewall](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md)
* [Setting Up Rate Limiting for Routes](/zen-firewall/zen-features/setting-up-rate-limiting-for-routes.md)
* [Monitor Outbound Domains](/zen-firewall/zen-features/monitor-outbound-domains.md)

Additional information:

* [Zen Performance & Reliability](/zen-firewall/miscellaneous/how-zen-works-performance-reliability.md)
* [Blocking vs Detection Mode in Zen Firewall](/zen-firewall/zen-features/blocking-vs-detection-mode-in-zen-firewall.md)
* [Understanding Your Zen Statistics](/zen-firewall/zen-features/understanding-your-zen-statistics.md)
  {% endstep %}
  {% endstepper %}
