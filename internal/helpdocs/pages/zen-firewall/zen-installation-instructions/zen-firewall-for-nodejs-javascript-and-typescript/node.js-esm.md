# Nodejs (ESM)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

{% hint style="warning" %}
Please also [check the documentation on how to integrate Zen with your used web framework](/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript.md). This page only includes additional steps for apps running as native ESM.
{% endhint %}

## Requirements

* We recommend Node.js 24+
* AWS Lambda handler integration.
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

## Installation & Configuration

{% hint style="info" %}
If you haven't already [follow the specific guide for your framework](/zen-firewall/zen-installation-instructions/zen-firewall-for-nodejs-javascript-and-typescript.md) before continuing here
{% endhint %}

Modify the start command of your application to include the Zen firewall:

```bash
node -r @aikidosec/firewall/instrument your-app.js
```

Alternatively, you can set the `NODE_OPTIONS` environment variable to include the Zen firewall:

```bash
export NODE_OPTIONS='-r @aikidosec/firewall/instrument'
```

## Loading environment variables

When using `--require` / `-r` to preload the Zen firewall, the instrumentation hook runs before your application code. This means environment variables loaded by packages like `dotenv` will not be available when Zen starts.

To ensure `AIKIDO_TOKEN` and other environment variables are available during instrumentation, use Node.js's native `--env-file` flag:

```bash
node --env-file=.env -r @aikidosec/firewall/instrument your-app.js
```

## Use Zen together with Sentry

You need to use Node.js v24.11.1 / v25.1.0 or later to use Zen together with Sentry in an ESM application. Follow the [Sentry instructions for ESM](https://docs.sentry.io/platforms/javascript/guides/node/install/esm/) to set up Sentry. After that, make sure to preload Zen using `--require` / `-r` before loading Sentry:

```bash
node -r @aikidosec/firewall/instrument --import ./instrument.mjs your-app.js
```

## Known issues

Zen can not protect ESM sub-dependencies of an ESM package. For example if an ESM package `foo` imports a sub-dependency `bar` that is also an ESM package, Zen will not be able to protect the code in `bar`. This is because the V8 engine does not allow Node.js to observe the evaluation of inner ESM packages (yet). Open issue: [Adding an evaluation hook for v8::Module](https://issues.chromium.org/u/1/issues/384413088).
