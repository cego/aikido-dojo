# Information for Support

When contacting Aikido Support about Zen Firewall, please work through the checks below first and include the requested information in your ticket. This helps us diagnose issues faster and avoid unnecessary back and forth.

### 1. Check your token

Zen needs a valid token to communicate with Aikido. If your app does not appear in the [Zen dashboard](https://app.aikido.dev/runtime/services), the token is the first thing to verify.

* Confirm `AIKIDO_TOKEN` is set in the environment where your app actually runs (not just your shell).
* Make sure the token has not been wrapped in extra quotes or whitespace.
* Verify the token belongs to the correct app. See [Creating an Aikido Zen Firewall Token](/zen-firewall/zen-installation-instructions/creating-an-aikido-zen-firewall-token.md).
* Enable [`AIKIDO_DEBUG=true`](/zen-firewall/zen-installation-instructions/configuration-via-environment-variables.md#debugging) and confirm the startup log reports that a token was detected.

### 2. Check connection to Aikido

Zen requires outbound HTTPS access to Aikido's cloud. If the token is set but the app still does not appear in the dashboard, network restrictions are the most likely cause.

See [Outbound Network Connections for Zen](/zen-firewall/miscellaneous/outbound-network-connections-for-zen.md) for the full list of domains (EU, US, ME) and connectivity test commands. Run those checks from the same host (and same container, if applicable) where your app runs, and share the output if any of them fail.

### 3. Information for support

When opening a ticket, please share:

**App & environment**

* Link to the app in the Aikido [Zen dashboard](https://app.aikido.dev/runtime/services)
* Programming language and runtime version (e.g. Node.js 20, Python 3.12, PHP 8.3)
* Framework and version (e.g. Express, Django, Spring MVC, Laravel, ASP.NET Core)
* Zen package / library version
* Operating system, container base image, or hosting platform (e.g. Ubuntu 24.04, `php:8.3-fpm`, AWS Lambda, Kubernetes)
* Whether the app runs behind a proxy, load balancer, CDN, or service mesh

**Configuration**

* Mode: blocking (`AIKIDO_BLOCK=true`) or detection only
* Any other `AIKIDO_*` environment variables you have set. See [Configuration via Environment Variables](/zen-firewall/zen-installation-instructions/configuration-via-environment-variables.md)

**Issue reproducibility**

* Is the issue consistently reproducible?
* Does it affect a single endpoint or multiple endpoints?
* A sample request that triggers the issue (with secrets and PII redacted)
* What you expected to happen vs. what actually happened

**Debug logs**

1. Set `AIKIDO_DEBUG=true` in the environment where your app runs.
2. Restart the app and reproduce the issue.
3. Capture the relevant section of the logs and attach it to the ticket. If you are not sure which lines are Zen's, search for `Aikido` or `Zen` in the output.

Where Zen's output ends up depends on your runtime:

* **Node.js**: Check the stdout or stderr of the Node process. Look wherever your process manager (PM2, systemd, Docker, Kubernetes, `npm start` console, etc.) collects them.
* **Python**: Check the stdout or stderr of the Python process. With WSGI / ASGI servers like Gunicorn or Uvicorn, check the error log destination configured for the workers (e.g. Gunicorn's `--error-logfile`, or stdout when run in the foreground).
* **PHP**: The output location depends on the SAPI:
  * Apache (`mod_php`): Apache's error log (the path configured in the `ErrorLog` directive, often `/var/log/apache2/error.log`).
  * PHP-FPM: the `error_log` configured in your FPM pool, plus the upstream web server's error log.
  * FrankenPHP: Caddy / FrankenPHP stdout, or the Caddy log file if you've configured one.
  * CLI workers (queue, cron, artisan): stderr of the command.
* **Java (Spring, Javalin, etc.)**: Check your application's configured logging framework (SLF4J / Logback / Log4j2). Zen messages appear in the same destination as the rest of your app logs (console, log file, JSON to stdout in containers).
* **.NET (ASP.NET Core / Framework)**: Check your application's configured `ILogger` sinks. With the default Console logger this is stdout (visible via `dotnet run`, `journalctl`, or container logs). IIS-hosted apps surface stdout via log files only when `stdoutLogEnabled="true"` is set in `web.config`.
* **Ruby on Rails**: Check the Rails logger (`log/development.log`, `log/production.log`, or whatever you've configured). Standalone Ruby scripts log to stderr.
* **Go**: Check the stdout or stderr of the Go binary.

If you run in a container or serverless platform (Docker, Kubernetes, AWS Lambda / CloudWatch, Cloud Run, App Service, Heroku, Fly.io, etc.), the process still emits the lines above, but the platform forwards them to its log collector. Pull the logs from there.

Providing all of the above allows us to identify framework, runtime, or environment specific issues much more efficiently.
