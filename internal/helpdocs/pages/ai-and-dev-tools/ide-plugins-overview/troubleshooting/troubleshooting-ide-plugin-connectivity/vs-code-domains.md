# VS Code - domains

To function correctly, the Local Scanner must be able to access the following domains over HTTPS (port 443):

* `*.aikido.dev`
* &#x20;Depending on `Download All Dependencies From Aikido Domain` in Settings
  * Disabled (default): `github.com`
  * Enabled: `static-binaries.aikido.dev`

\
**Optional domains**

For Aikido Expansion Pack usage

* `device-protection.aikido.help`
* `raw.githubusercontent.com`

For dependency scanning

* Depending on `Download All Dependencies From Aikido Domain` in Settings
  * Disabled (default): `mirror.gcr.io` / `ghcr.io`
  * Enabled: `static-cve.aikido.dev`

For Autofix usage (depending on ecosystem):

* `registry.npmjs.org`
* `repo1.maven.org`
* `api.nuget.org`

Links to help docs:

* `help.aikido.dev`

Analytics

* `product.aikido.dev`

**If you cannot use wildcards:**

Where `*.aikido.dev` is not accepted, the wildcard expands to a small, fixed set of hosts (additional to the ones mentioned above):

| Host                                  | When                |
| ------------------------------------- | ------------------- |
| `ide.aikido.dev` and `app.aikido.dev` | EU region (default) |
| `app.us.aikido.dev`                   | US region           |
| `app.au.aikido.dev`                   | Australia region    |
| `app.me.aikido.dev`                   | Middle East region  |
