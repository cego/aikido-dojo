# Aikido Images: Supported Versions & Release Lines

Aikido Images provides drop-in replacements for popular base images with Critical and High vulnerabilities fixed. The full catalog, spanning thousands of images across language runtimes, infrastructure tools, and Kubernetes operators, is available at [images.aikido.dev](https://images.aikido.dev).

Aikido builds a new Aikido Image whenever a new upstream release is published.

## How version tracking works

A **release line** is a version pattern (for example, Node 24.x.x or Python 3.13.x) that Aikido monitors on upstream registries for new releases. When a matching tag is published, Aikido automatically builds an Aikido Image for it and makes it available in the catalog. You don't need to request it.

New releases within a tracked line are always picked up automatically. What's not automatic is new major or minor versions. When Python 3.15 or Node 28 are released upstream, Aikido needs to add a new release line before they are tracked. New lines are added as versions reach general availability.

{% hint style="info" %}
Rolling tags like `node:24` or `python:3.13` are supported too. They always resolve to the latest Aikido Image release within that line, so your Dockerfile doesn't need to pin a specific version.
{% endhint %}

## Patch cadence and SLA

Aikido scans every tracked release line daily. Aikido patches vulnerabilities in the CISA Known Exploited Vulnerabilities (KEV) catalog within 48 hours, and other Critical and High CVEs within 7 days.

## End-of-life versions

Release lines that upstream no longer maintains aren't tracked by default. Aikido can add extended support for them under the same SLA, for example `php:5.6.40`.

## Variants

Aikido Images are available in the following variants:

* **Standard images** with the regular image tag.
* **FIPS images** for workloads that require FIPS-validated cryptography.
* **Hardened images** with a reduced attack surface.
* **Runtime variants** that are smaller and leaner.
* **Dev variants** that include build tooling.

Not every variant exists for every image. Browse the public catalog at [images.aikido.dev](https://images.aikido.dev) to see the images and variants available.

## Request a version or image

If the version or image you need isn't covered by an existing release line, see [Requesting Aikido Images](https://help.aikido.dev/autofix-and-remediation/aikido-images/requesting-aikido-images).
