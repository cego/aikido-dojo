# Nodejs (Google Cloud Pub/Sub)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

<https://github.com/AikidoSec/firewall-node>

## Requirements

* Node.js 16+.
* Google Cloud Pub/Sub (`@google-cloud/pubsub` 4.x or 5.x).
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

Load Zen before creating Pub/Sub clients:

{% @aikido-custom-code/code-highlight language="javascript" content="+require("@aikidosec/firewall");

const { PubSub } = require("@google-cloud/pubsub");
const client = new PubSub();" %}
{% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

Set runtime variables:

```bash
AIKIDO_TOKEN=AIK_RUNTIME_...
AIKIDO_BLOCK=false
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
**Test your app**

Publish and consume a few test messages.

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
