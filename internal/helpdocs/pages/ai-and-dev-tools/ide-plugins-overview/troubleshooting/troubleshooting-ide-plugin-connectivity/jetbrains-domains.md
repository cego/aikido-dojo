# JetBrains - domains

To function correctly, the IDE must be able to access the following domains over HTTPS (port 443):

* `*.aikido.dev`

\
**Optional domains**

For IaC scanning:

* `github.com`

For Aikido Expansion Pack usage

* `device-protection.aikido.help`
* `raw.githubusercontent.com`

For dependency scanning

* `static-binaries.aikido.dev`
* `mirror.gcr.io` / `ghcr.io` or `static-cve.aikido.dev` when `Download From Aikido Domain` is enabled in the IDE settings

Links to help docs:

* `help.aikido.dev`

Analytics

* `product.aikido.dev`. Can be disabled by disabling `Aikido Telemetry > Send Usage Statistics` in the IDE settings.

**If you cannot use wildcards:**

Where `*.aikido.dev` is not accepted, the wildcard expands to a small, fixed set of hosts (additional to the ones mentioned above):

| Host                                  | When                |
| ------------------------------------- | ------------------- |
| `ide.aikido.dev` and `app.aikido.dev` | EU region (default) |
| `app.us.aikido.dev`                   | US region           |
| `app.au.aikido.dev`                   | Australia region    |
| `app.me.aikido.dev`                   | Middle East region  |
