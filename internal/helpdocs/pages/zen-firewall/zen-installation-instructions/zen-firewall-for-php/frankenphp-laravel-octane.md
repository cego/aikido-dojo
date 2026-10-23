# FrankenPHP (Laravel Octane)

This guide will walk you through installing and setting up Zen Firewall by Aikido for your application. Follow the steps below to protect your application.

If you encounter any issues or problems, don't hesitate reach out on support chat or Github issues

<https://github.com/AikidoSec/firewall-php>

## Requirements

* FrankenPHP 8.2+ (tested up to FrankenPHP 8.5).
* No framework / custom middleware.
* [Aikido account](/getting-started/setting-up-your-account.md) & [Zen Firewall token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md)

## Installation & Configuration

{% stepper %}
{% step %}
**Install Zen Firewall by Aikido**

Zen for FrankenPHP is installed as a system package.

Prerequisites:

* Ensure you have `sudo` privileges.
* Use a supported FrankenPHP version (8.2+).

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

{% endstep %}

{% step %}
**Start Zen Firewall in dry / detection-only mode**

**Via Laravel Sail**

**Download Aikido installation package**

In `docker/version/Dockerfile`:

```docker
ARG AIKIDO_VERSION=1.5.2

RUN curl -L -o /tmp/aikido-php-firewall.deb \
    "https://github.com/AikidoSec/firewall-php/releases/download/v${AIKIDO_VERSION}/aikido-php-firewall.$(uname -m).deb"
```

**Pass the Aikido environment variables to FrankenPHP in Octane Caddyfile and configure the worker**

In `vendor/laravel/octane/src/Commands/stubs/Caddyfile`:

{% @aikido-custom-code/code-highlight content="! vendor/laravel/octane/src/Commands/stubs/Caddyfile
{$CADDY\_SERVER\_SERVER\_NAME} {
route {
php\_server {

* ```
        env AIKIDO_TOKEN "AIK_RUNTIME_...."
  ```
* ```
        env AIKIDO_BLOCK "True"
    }
  ```
  }
  }" %}

**Call \aikido\worker\_rinit() and \aikido\worker\_rshutdown() in Octane worker script**

Wrap your request handler with these calls to ensure Aikido processes each request.

In `vendor/laravel/octane/bin/frankenphp-worker.php`:

{% @aikido-custom-code/code-highlight language="php" content="! vendor/laravel/octane/bin/frankenphp-worker.php

<?php

// ... Octane bootstrap ...

$handleRequest = static function () use ($worker, $frankenPhpClient, $debugMode) {
+    \aikido\worker_rinit();
    // Octane request handling logic
+    \aikido\worker_rshutdown();
};

while ($requestCount < $maxRequests && frankenphp_handle_request($handleRequest)) {
    $requestCount++;
}" %}

In `docker/version/start-container`, add the following installation steps towards the end of the script:

```bash
ls -sf /var/www/html/frankenphp /usr/bin/frankenphp
dpkg -i -E ./aikido-php-firewall.$(uname -m).deb
```

Now you can start Sail:

```bash
./vendor/bin/sail up --build
```

**Via Docker Image**

**Install Aikido in your Dockerfile**

{% @aikido-custom-code/code-highlight content="! Dockerfile
FROM dunglas/frankenphp:php${PHP_VERSION}-bookworm

+ARG AIKIDO_VERSION=1.5.2

+RUN apt-get update && apt-get install -y curl \
+   && curl -L -o /tmp/aikido-php-firewall.deb \
+   "https://github.com/AikidoSec/firewall-php/releases/download/v${AIKIDO_VERSION}/aikido-php-firewall.$(uname -m).deb" \
+   && dpkg -i /tmp/aikido-php-firewall.deb \
+   && rm /tmp/aikido-php-firewall.deb \
+   && apt-get clean

COPY . /app

ENTRYPOINT ["php", "artisan", "octane:frankenphp"]" %}

**Pass the Aikido environment variables to FrankenPHP in Octane Caddyfile and configure the worker**

In `vendor/laravel/octane/src/Commands/stubs/Caddyfile`:

{% @aikido-custom-code/code-highlight content="vendor/laravel/octane/src/Commands/stubs/Caddyfile
{$CADDY_SERVER_SERVER_NAME} {
    route {
        php_server {
+            env AIKIDO_TOKEN "AIK_RUNTIME_...."
+            env AIKIDO_BLOCK "True"
        }
    }
}" %}

**Call \aikido\worker\_rinit() and \aikido\worker\_rshutdown() in Octane worker script**

Wrap your request handler with these calls to ensure Aikido processes each request.

In `vendor/laravel/octane/bin/frankenphp-worker.php`:

{% @aikido-custom-code/code-highlight language="php" content="! vendor/laravel/octane/bin/frankenphp-worker.php
<?php

// ... Octane bootstrap ...

$handleRequest = static function () use ($worker, $frankenPhpClient, $debugMode) {
+    \aikido\worker_rinit();
    // Octane request handling logic
+    \aikido\worker_rshutdown();
};

while ($requestCount < $maxRequests && frankenphp_handle_request($handleRequest)) {
    $requestCount++;
}" %}

**Start the container**

```
docker run -e AIKIDO_TOKEN="AIK_RUNTIME_..." your-app
```

We recommend to start your app in dry mode to ensure it works as expected without blocking any requests. We advise running Zen Firewall in staging for two weeks to avoid false positives.

```
env AIKIDO_BLOCK "False"
```

{% hint style="info" %}
You can define an environment variable `AIKIDO_DEBUG=true` to enable debug mode for more detailed information about what the agent is doing. For more information about your environment variables: [/pages/cS1d1R2qLfIJt7OJ7dW5](/pages/cS1d1R2qLfIJt7OJ7dW5 "mention")
{% endhint %}

For example

```env
AIKIDO_TOKEN=AIK_RUNTIME_...
AIKIDO_BLOCK=false
```

Restart your FrankenPHP process manager/web server after changes.
{% endstep %}

{% step %}
**Enable Rate limiting and User blocking**

Use `set_user` and `should_block_request` in your request pipeline:

{% @aikido-custom-code/code-highlight language="php" content="+if (extension_loaded('aikido')) {
+   \aikido\set_user('123', 'John Doe');
+
+   $decision = \aikido\should_block_request();
+
+   if ($decision->block) {
+       if ($decision->type === 'ratelimited') {
+           http_response_code(429);
+           header('Retry-After: ' . $decision->retry_after);
+           exit('You are rate limited by Zen.');
+       }
+       http_response_code(403);
+       exit('You are blocked by Zen.');
+   }
+}" %}
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

* [/pages/f45DP3VAPtfnxJHfRRQO](/pages/f45DP3VAPtfnxJHfRRQO "mention")
* [/pages/iP6p6jOdj6CHR4SrBjeM](/pages/iP6p6jOdj6CHR4SrBjeM "mention")
* [/pages/DwUeaYMkfLPRxTzU0H7N](/pages/DwUeaYMkfLPRxTzU0H7N "mention")
* [/pages/Cuhw6Bwome1nyBc0Wq51](/pages/Cuhw6Bwome1nyBc0Wq51 "mention")
* [/pages/aTqnMKVKyl8OLchBr12Y](/pages/aTqnMKVKyl8OLchBr12Y "mention")
* [/pages/X0AXHHmUDrRPMsgyA71F](/pages/X0AXHHmUDrRPMsgyA71F "mention")
* [/pages/F2jaW8oZmBttuMA7fANi](/pages/F2jaW8oZmBttuMA7fANi "mention")

Additional information:

* [/pages/9pwdJjtFQRt3r4jy7U2q](/pages/9pwdJjtFQRt3r4jy7U2q "mention")
* [/pages/7fhmf4YIX1YczgqFA1rE](/pages/7fhmf4YIX1YczgqFA1rE "mention")
* [/pages/VL9A5lxpqzryy4ThsGRJ](/pages/VL9A5lxpqzryy4ThsGRJ "mention")
{% endstep %}
{% endstepper %}
