# PHP on nginx or Apache with PHP-FPM

This guide installs Zen Firewall for a PHP app served by nginx or Apache through PHP-FPM. Zen runs as a PHP extension inside the PHP-FPM workers, so you configure it in the PHP-FPM pool. nginx and Apache need no changes.

<https://github.com/AikidoSec/firewall-php>

{% hint style="info" %}
If Apache runs PHP through `mod_php` instead of PHP-FPM, follow [PHP on Apache (mod\_php)](/zen-firewall/zen-installation-instructions/zen-firewall-for-php/apache-mod-php.md).
{% endhint %}

## Prerequisites

* PHP 7.2 or later (tested up to PHP 8.5).
* nginx or Apache that passes PHP requests to PHP-FPM (`fastcgi_pass` in nginx, `proxy_fcgi` in Apache).
* `sudo` access on the server.
* An [Aikido account](/getting-started/setting-up-your-account.md) and a [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md).

## Set up Zen Firewall

{% stepper %}
{% step %}

### Install the Zen package

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

The package enables the extension for every PHP version and runtime it finds on the server, including PHP-FPM.
{% endstep %}

{% step %}

### Add the token to your PHP-FPM pool

PHP-FPM clears the environment of its workers by default (`clear_env = yes`), so variables you export in your shell or service manager don't reach Zen. Set them with `env[...]` in the pool that serves your app instead.

Open the pool file. The default pool is called `www`:

{% tabs %}
{% tab title="Debian and Ubuntu" %}

```shell
sudo nano /etc/php/<php-version>/fpm/pool.d/www.conf
```

{% endtab %}

{% tab title="RHEL, CentOS, and Fedora" %}

```shell
sudo nano /etc/php-fpm.d/www.conf
```

{% endtab %}
{% endtabs %}

Add these lines at the end of the pool:

{% code title="[www.conf](http://www.conf)" %}

```ini
env[AIKIDO_TOKEN] = AIK_RUNTIME_...
env[AIKIDO_BLOCK] = false
```

{% endcode %}

`AIKIDO_BLOCK=false` starts Zen in detection-only mode: it reports attacks but doesn't block them. Run it like this in staging for about 2 weeks to rule out false positives. See [Configuration via environment variables](/zen-firewall/zen-installation-instructions/configuration-via-environment-variables.md) for the other options.

{% hint style="warning" %}
If your app uses more than 1 pool, add the lines to every pool that serves it. A pool without `AIKIDO_TOKEN` doesn't report to Aikido.
{% endhint %}
{% endstep %}

{% step %}

### Restart PHP-FPM

Restart PHP-FPM to load the extension and the new pool settings. Reloading nginx or Apache is not enough.

{% tabs %}
{% tab title="Debian and Ubuntu" %}

```shell
sudo systemctl restart php<php-version>-fpm
```

{% endtab %}

{% tab title="RHEL, CentOS, and Fedora" %}

```shell
sudo systemctl restart php-fpm
```

{% endtab %}
{% endtabs %}
{% endstep %}

{% step %}

### Enable request blocking and user identification

Use this middleware to enable rate limiting, user identification, and blocking features.

Zen Firewall does not require this middleware to block attacks. Core attack protection works without it. The middleware provides additional request context that Zen uses for protections such as [rate limiting](/zen-firewall/zen-features/setting-up-rate-limiting-for-routes.md), [user blocking](/zen-firewall/zen-features/blocking-users-with-zen-firewall.md), [bot blocking](/zen-firewall/zen-features/blocking-bot-traffic-with-zen-firewall.md), [country blocking](/zen-firewall/zen-features/blocking-traffic-by-country-with-zen-firewall.md), and [threat actor](/zen-firewall/zen-features/blocking-known-threat-actors-with-zen-firewall.md) blocking.

Adapt the example to fit how your application identifies users and handles requests.

Call `set_user` and `should_block_request` early in every request, for example in your front controller (`index.php`) or a shared bootstrap file:

```php
if (extension_loaded('aikido')) {
    \aikido\set_user('123', 'John Doe');

    $decision = \aikido\should_block_request();

    if ($decision->block) {
        if ($decision->type === 'ratelimited') {
            http_response_code(429);
            header('Retry-After: ' . $decision->retry_after);
            exit('You are rate limited by Zen.');
        }
        http_response_code(403);
        exit('You are blocked by Zen.');
    }
}
```

{% endstep %}

{% step %}

### Test your app

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

### Set up rate limiting

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

### Explore next steps

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
