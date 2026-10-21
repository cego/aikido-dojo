# On-Prem Registry Proxy

## Overview

The Aikido On-Prem Registry Proxy can be deployed as a **transparent package firewall**, allowing organizations to enforce package security policies without requiring developers or CI systems to change their package manager configuration.

Instead of configuring clients to use private registry URLs or special proxy paths, organizations redirect supported public package registries to the Aikido Registry Proxy using corporate DNS. Developers continue using the default package manager configuration, while the proxy transparently evaluates requests before retrieving content from the public registries.

## Architecture

## Running the Proxy

```bash
docker run -p 8080:8080 -e DEPLOYMENT_MODE=onprem public.ecr.aws/aikido-security/aikido-registry-proxy:latest
```

The proxy runs on port `8080` and expects TLS to be terminated at the ingress layer. The ingress must forward the original `Host` header (or set `X-Forwarded-Host`) so the proxy can identify the ecosystem.

A health endpoint is available at `GET /health` (no authentication required), suitable for Kubernetes liveness and readiness probes.

## Configuration

<table data-search="false"><thead><tr><th>Variable</th><th>Default</th><th>Description</th></tr></thead><tbody><tr><td><code>DEPLOYMENT_MODE</code></td><td>Required to override default and set to <code>onprem</code></td><td>Required to setup the DNS Proxy</td></tr><tr><td><code>UPSTREAM_DNS</code></td><td><code>1.1.1.1</code></td><td>Public DNS resolver used for the proxy's own outbound calls. Must point to a real public resolver, not the intercepting corporate DNS. Accepts comma-separated values for multiple resolvers.</td></tr><tr><td><code>MINIMUM_PACKAGE_AGE_HOURS</code></td><td><code>48</code></td><td>Block packages published more recently than this threshold.</td></tr><tr><td><code>NPM_UPSTREAM_URL</code></td><td><code>https://registry.npmjs.org</code></td><td>Override to point at an internal npm mirror.</td></tr><tr><td><code>PYPI_UPSTREAM_URL</code></td><td><code>https://pypi.org</code></td><td>Override to point at an internal PyPI mirror.</td></tr><tr><td><code>PYPI_FILES_UPSTREAM_URL</code></td><td><code>https://files.pythonhosted.org</code></td><td>Override to point at an internal PyPI files mirror.</td></tr><tr><td><code>MAVEN_UPSTREAM_URL</code></td><td><code>https://repo1.maven.org/maven2</code></td><td>Override to point at an internal Maven mirror.</td></tr></tbody></table>

## Supported Host Mapping

The Registry Proxy routes requests to the appropriate ecosystem based on the incoming HTTP Host header.

| Public Hostname        | Ecosystem     |
| ---------------------- | ------------- |
| registry.npmjs.org     | npm           |
| pypi.org               | PyPI          |
| files.pythonhosted.org | PyPI          |
| repo1.maven.org        | Maven Central |

## Registry Compatibility

Metadata responses do not rewrite URLs to proxy paths. Clients receive URLs in their original public registry format and are never redirected to non-intercepted hosts.

### npm

* Packument responses remain compatible with `registry.npmjs.org`.
* Tarball downloads use the standard npm workflow.

### PyPI

* Metadata requests are served from `pypi.org`.
* Package downloads are served from `files.pythonhosted.org`.

### Maven

* Requests to Maven Central use the standard `/maven2/...` path.

## Egress IP and Rate Limiting

Because all proxied traffic egresses from a single IP (or two IPs with two replicas), Maven Central and other public registries may rate-limit or block requests that originate from that IP. This is a known limitation of the transparent proxy model.

To mitigate this, consider routing egress traffic through a NAT gateway with multiple IPs, or pointing the upstream URL overrides at an internal mirror that is not subject to public registry rate limits.

## Preventing Upstream Routing Loops

The proxy resolves all of its own outbound calls using the `UPSTREAM_DNS` resolver (default `1.1.1.1`), bypassing the container's system DNS entirely. This means even if the pod inherits the intercepting corporate DNS, the proxy's upstream fetches always reach the real public registries and never loop back.

No split-horizon DNS configuration is required on the proxy host itself. As an additional safeguard, you can configure the pod's DNS resolver explicitly via Kubernetes `dnsConfig`.

## Request Flow

1. A developer or CI system installs a package using the standard package manager.
2. Corporate DNS resolves the public registry hostname to the internal Registry Proxy.
3. The Registry Proxy receives the request using the original HTTP Host header.
4. The proxy determines the target ecosystem based on the requested hostname.
5. Malware detection and minimum package age policies are evaluated.
6. If the request is allowed, the proxy retrieves the package or metadata from the corresponding public registry.
7. The response is returned transparently to the client.
