# PHP Framework Laravel

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

<https://github.com/AikidoSec/firewall-php>

## Requirements

* PHP 7.2+ (tested up to PHP 8.5).
* Laravel (supported by Zen Firewall).
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

## Installation & Configuration

#### Install Zen Firewall by Aikido

Zen for PHP is installed as a system package.

Prerequisites:

* Ensure you have `sudo` privileges.
* Use a supported PHP version (7.2+).

For Red Hat-based systems (RHEL, CentOS, Fedora):

```bash
rpm -Uvh --oldpackage https://github.com/AikidoSec/firewall-php/releases/latest/download/aikido-php-firewall.x86_64.rpm
```

For Red Hat-based arm64 / aarch64:

```bash
rpm -Uvh --oldpackage https://github.com/AikidoSec/firewall-php/releases/latest/download/aikido-php-firewall.aarch64.rpm
```

For Debian-based systems (Debian, Ubuntu):

```bash
curl -L -O https://github.com/AikidoSec/firewall-php/releases/latest/download/aikido-php-firewall.x86_64.deb
dpkg -i -E ./aikido-php-firewall.x86_64.deb
```

For Debian-based arm64 / aarch64:

```bash
curl -L -O https://github.com/AikidoSec/firewall-php/releases/latest/download/aikido-php-firewall.aarch64.deb
dpkg -i -E ./aikido-php-firewall.aarch64.deb
```

#### Enable Request Blocking and User Identification

Use this middleware to enable rate limiting, user identification, and blocking features.

Zen Firewall does not require this middleware to block attacks. Core attack protection works without it. The middleware provides additional request context that Zen uses for protections such as [rate limiting](/zen-firewall/zen-features/setting-up-rate-limiting-for-routes.md), [user blocking](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md), [bot blocking](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md), [country blocking](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md), and [threat actor](/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall.md) blocking.

Adapt the example to fit how your application identifies users and handles requests.

Create `app/Http/Middleware/AikidoMiddleware.php`:

{% @aikido-custom-code/code-highlight language="php" content="! app/Http/Middleware/AikidoMiddleware.php
+namespace App\Http\Middleware;
\+
+use Closure;
+use Illuminate\Support\Facades\Auth;
\+
+class AikidoMiddleware
+{

* public function handle($request, Closure $next)
* {
* ```
     if (!extension_loaded('aikido')) {
  ```
* ```
         return $next($request);
  ```
* ```
     }
  ```
*
* ```
     $userId = Auth::id();
  ```
* ```
     if ($userId) {
  ```
* ```
         \aikido\set_user($userId);
  ```
* ```
     }
  ```
*
* ```
     $decision = \aikido\should_block_request();
  ```
* ```
     if ($decision->block) {
  ```
* ```
         if ($decision->type == 'ratelimited') {
  ```
* ```
             abort(429, 'You are rate limited by Zen.', ['Retry-After' => (string) $decision->retry_after]);
  ```
* ```
         }
  ```
* ```
         abort(403, 'You are blocked by Zen.');
  ```
* ```
     }
  ```
*
* ```
     return $next($request);
  ```
* }
  +}" %}

Register middleware in `bootstrap/app.php`.

#### Start Zen Firewall in dry / detection-only mode

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

For example: Set environment variables in `.env`:

```env
AIKIDO_TOKEN=AIK_RUNTIME_...
AIKIDO_BLOCK=false
```

Restart PHP-FPM / web server after configuration changes.

#### Test your app

Browse to your application and perform a couple of actions or open a couple of pages. Zen will automatically discover the routes in your application.

{% hint style="info" %}
Zen sends data back to Aikido every 10 minutes
{% endhint %}

You can verify a working agent by looking at the following pages of your Zen application:

* **Events**: Should show an "Application started" event.
* **Routes**: After some time your application routes will start showing here with the method, route and requests.
* **Instances**: Should show the number of active instances for your application where Zen is installed.

#### Setup rate limiting in the dashboard

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

#### Next steps

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
