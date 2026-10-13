# Supported Versions & Release Lines

Aikido Images provides drop-in replacements for popular base images with Critical and High OS vulnerabilities fixed. The full catalog - spanning thousands of images across language runtimes, infrastructure tools, and Kubernetes operators - is available at [images.aikido.dev](https://images.aikido.dev).

This page covers a specific part of the catalog: the release lines that Aikido actively tracks, where a new Aikido Image is built automatically whenever a new upstream release is published.

## How version tracking works

A **release line** is a version pattern - for example, Node 24.x.x or Python 3.13.x - that Aikido monitors on upstream registries for new releases. When a matching tag is published, Aikido automatically builds an Aikido Image for it and makes it available in the catalog. You don't need to request it.

New releases within a tracked line are always picked up automatically. What's not automatic is new major or minor versions - when Python 3.15 or Node 28 are released upstream, Aikido needs to add a new release line before they are tracked. New lines are added as versions reach general availability, and you can also [request one](/autofix-and-remediation/aikido-images/requesting-aikido-images.md).

{% hint style="info" %}
Rolling tags like `node:24` or `python:3.13` are supported too - they always resolve to the latest Aikido Image release within that line, so your Dockerfile doesn't need to pin a specific version.
{% endhint %}

## Currently tracked release lines

<details>

<summary>Python</summary>

Versions tracked: **3.11, 3.12, 3.13, 3.14**

Each version is available in the following variants:

| Variant          | Example tag                    |
| ---------------- | ------------------------------ |
| Default (Debian) | `python:3.13.15`               |
| Slim             | `python:3.13.15-slim`          |
| Bookworm         | `python:3.13.15-bookworm`      |
| Slim Bookworm    | `python:3.13.15-slim-bookworm` |
| Trixie           | `python:3.13.15-trixie`        |
| Slim Trixie      | `python:3.13.15-slim-trixie`   |
| Alpine 3.23      | `python:3.13.15-alpine3.23`    |
| Alpine 3.24      | `python:3.13.15-alpine3.24`    |

</details>

<details>

<summary>Node.js</summary>

Versions tracked: **22, 24, 26**

Each version is available in the following variants:

| Variant       | Example tag                  |
| ------------- | ---------------------------- |
| Default       | `node:26.10.0`               |
| Bookworm      | `node:26.10.0-bookworm`      |
| Bookworm Slim | `node:26.10.0-bookworm-slim` |
| Trixie        | `node:26.10.0-trixie`        |
| Trixie Slim   | `node:26.10.0-trixie-slim`   |
| Alpine 3.23   | `node:26.10.0-alpine3.23`    |
| Alpine 3.24   | `node:26.10.0-alpine3.24`    |

</details>

<details>

<summary>Go</summary>

Versions tracked: **1.26, 1.27**

Each version is available in the following variants:

| Variant          | Example tag                |
| ---------------- | -------------------------- |
| Default (Debian) | `golang:1.26.8`            |
| Bookworm         | `golang:1.26.8-bookworm`   |
| Trixie           | `golang:1.26.8-trixie`     |
| Alpine 3.23      | `golang:1.26.8-alpine3.23` |
| Alpine 3.24      | `golang:1.26.8-alpine3.24` |

</details>

<details>

<summary>Ruby</summary>

Versions tracked: **3.3, 3.4, 4.0**

Each version is available in the following variants:

| Variant          | Example tag                 |
| ---------------- | --------------------------- |
| Default (Debian) | `ruby:3.4.10`               |
| Bookworm         | `ruby:3.4.10-bookworm`      |
| Slim Bookworm    | `ruby:3.4.10-slim-bookworm` |
| Trixie           | `ruby:3.4.10-trixie`        |
| Slim Trixie      | `ruby:3.4.10-slim-trixie`   |
| Alpine 3.23      | `ruby:3.4.10-alpine3.23`    |
| Alpine 3.24      | `ruby:3.4.10-alpine3.24`    |

</details>

<details>

<summary>Alpine</summary>

Minor versions tracked: **3.20, 3.21, 3.22, 3.23, 3.24**

New patch releases within each tracked minor version are picked up automatically (e.g. `alpine:3.23.7` when it is published).

</details>

<details>

<summary>Infrastructure &#x26; applications</summary>

New releases for these images are picked up automatically as they are published upstream.

| Image                                  | Example tag                    | Notes                                                            |
| -------------------------------------- | ------------------------------ | ---------------------------------------------------------------- |
| `quay.io/minio/minio`                  | `RELEASE.2025-09-07T16-13-09Z` | All stable releases in the `RELEASE.YYYY-MM-DDTHH-MM-SSZ` format |
| `n8nio/n8n`                            | `2.40.2`                       | Current release line + the 2.39.x LTS line                       |
| `registry.datadoghq.com/agent`         | `7.83.2`                       | Datadog Agent                                                    |
| `registry.datadoghq.com/cluster-agent` | `7.83.2`                       | Datadog Cluster Agent                                            |
| `mongodb/atlas`                        | `v1.58.3`                      | MongoDB Atlas Kubernetes Operator                                |
| `mongo`                                | `8.0.32`                       | MongoDB                                                          |
| `rabbitmq`                             | `4.2.9`                        | OpenSource Message Broker                                        |
| `kibana`                               | `9.5.3`                        | OpenSource Analytics and Visualisation Platform                  |
| `elasticsearch`                        | `9.5.3`                        | Search and Analytics Platform                                    |
| `nginx`                                | `1.31.6-alpine3.24`            | HTTP Web Server                                                  |
| `localstack/localstack`                | `2026.8`                       | Local Cloud Development Tool                                     |

</details>

<details>

<summary>Kubernetes operators &#x26; platform tools</summary>

| Image                                                                      | Example tag |
| -------------------------------------------------------------------------- | ----------- |
| `ghcr.io/external-secrets/kubernetes-external-secrets`                     | `8.5.5`     |
| `ghcr.io/open-telemetry/opentelemetry-operator/opentelemetry-operator`     | `0.159.0`   |
| `ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-java`   | `2.31.1`    |
| `ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-python` | `0.66b0`    |
| `quay.io/argoproj/kubectl-argo-rollouts`                                   | `v1.10.0`   |
| `quay.io/argoproj/argo-rollouts`                                           | `v1.10.0`   |
| `quay.io/brancz/kube-rbac-proxy`                                           | `v0.22.1`   |
| `gcr.io/knative-releases/knative.dev/serving/cmd/activator`                | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler`               | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/serving/cmd/autoscaler-hpa`           | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/serving/cmd/controller`               | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/serving/cmd/webhook`                  | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/operator/cmd/operator`                | `v1.23.1`   |
| `gcr.io/knative-releases/knative.dev/operator/cmd/webhook`                 | `v1.23.1`   |
| `gcr.io/knative-releases/knative.dev/net-gateway-api/cmd/controller`       | `v1.23.0`   |
| `gcr.io/knative-releases/knative.dev/net-gateway-api/cmd/webhook`          | `v1.23.0`   |
| `danielqsj/kafka-exporter`                                                 | `v1.10.0`   |

</details>

## Request a version or image

If the version or image you need isn't covered by an existing release line, see [Requesting Aikido Images](/autofix-and-remediation/aikido-images/requesting-aikido-images.md).
