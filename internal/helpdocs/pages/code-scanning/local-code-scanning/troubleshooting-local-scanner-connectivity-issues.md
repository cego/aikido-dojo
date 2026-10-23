# Troubleshooting Local Scanner connectivity issues

This guide helps customers run the Aikido Local Scanner behind corporate proxies, custom SSL inspection, or locked-down networks. It covers when and how to use `--enable-proxy` and `--ca-bundle`, which domains to whitelist, and how to fix common certificate errors.

To function correctly, the Local Scanner must be able to access the following domains over HTTPS (port 443):

* `*.aikido.dev`
* `*.aikido.io`

Optional

* `o1405212.ingest.sentry.io` for error reporting (can be disabled with `--disable-telemetry`)
* `mirror.gcr.io` / `ghcr.io` for dependency scanning vulnerability database
  * Not neccessary when using `--enable-aikido-hosted-dependencies-db` and `--enable-aikido-maven-cache`
* `api.nuget.org` for csproj dependency resolution
* `raw.githubusercontent.com` for vcpkg dependency resolution

**If you cannot use wildcards:**

Where `*.aikido.dev` and `*.aikido.io` are not accepted, each wildcard expands to a small, fixed set of hosts:

<table><thead><tr><th width="264.83984375">Host</th><th>When</th></tr></thead><tbody><tr><td><code>app.aikido.dev</code></td><td>EU region (default)</td></tr><tr><td><code>app.us.aikido.dev</code></td><td>US region</td></tr><tr><td><code>app.au.aikido.dev</code></td><td>Australia</td></tr><tr><td><code>app.me.aikido.dev</code></td><td>Middle East</td></tr><tr><td><code>static-cve.aikido.dev</code></td><td>With <code>--enable-aikido-hosted-dependencies-db</code></td></tr><tr><td><code>static-maven.aikido.io</code></td><td>With <code>--enable-aikido-maven-cache</code></td></tr></tbody></table>

#### Using `--enable-proxy`

The scanner does **not** use the system proxy by default. To send all HTTPS traffic through your corporate proxy, you must:

1. **Enable proxy support** with the `--enable-proxy` flag.
2. **Set the proxy URL** via the `HTTPS_PROXY` (or `https_proxy`) environment variable.

#### Using `--ca-bundle` (custom root CAs)

When your proxy or firewall performs **SSL inspection**, it typically replaces the server’s certificate with one signed by an **internal/corporate CA**.

Use **`--ca-bundle`** to point the scanner at a **PEM file** that contains your corporate root CA certificate(s) (and optionally intermediate CAs). The scanner will then trust those CAs when connecting through the proxy.

1. **Get your corporate CA certificate** (e.g. from your IT team or export from the machine’s trust store). It must be in **PEM** format.
2. **Pass it to the scanner** via `--ca-bundle` or the `AIKIDO_CA_BUNDLE` environment variable. Additionally it might be required to pass this value as `SSL_CERT_FILE` and `NODE_EXTRA_CA_CERTS` environment variables.
