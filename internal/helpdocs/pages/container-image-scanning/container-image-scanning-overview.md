# Container Image Scanning Overview

<table data-card-size="large" data-view="cards"><thead><tr><th></th><th></th><th data-hidden data-card-target data-type="content-ref"></th></tr></thead><tbody><tr><td><a href="/container-image-scanning/cloud-provider-registries.md"><strong>Connect Your Container Registry</strong></a></td><td></td><td><a href="/container-image-scanning/cloud-provider-registries.md">Cloud Provider Registries</a></td></tr><tr><td><a href="/container-image-scanning/standalone-registries.md"><strong>Connect Standalone Registry</strong></a></td><td></td><td><a href="/container-image-scanning/standalone-registries.md">Standalone Registries</a></td></tr></tbody></table>

Container image scanning in Aikido is designed to treat your images as first-class security assets. It connects three things that are usually siloed: registries (where images live), code (where you can actually fix issues), and runtimes (where risk materializes).

{% hint style="info" %}
Aikido scans both static images and containers deployed to your clusters. Learn more about [in-cluster image scanning](https://help.aikido.dev/cloud-scanning/kubernetes-cluster-scanning/kubernetes-in-cluster-image-scanning).
{% endhint %}

### Defining container image scanning

A container image bundles your application code, OS packages, language runtimes, web servers, and supporting tools into a single immutable artifact. Aikido’s container image scanning focuses on four broad categories of risk in those artifacts.

1. **Known vulnerabilities (CVEs)** in OS packages and other components.
2. **Outdated or end-of-life (EOL) runtimes** (e.g. language runtimes and web servers).
3. **Software composition and licensing**, via a detailed software bill of materials (SBOM).
4. **Malware** in application packages installed in the image (npm, PyPI, Maven, NuGet, Go, Cargo, RubyGems, and Packagist).

When you enable scanning for a registry like AWS ECR, Aikido inspects images for vulnerabilities, malware, licenses, and EOL runtimes, and rescans them regularly so that new CVEs are caught even if the image hasn’t been rebuilt or pushed recently.

Aikido also tracks the status of key runtimes discovered inside images (such as Python, Node.js, PHP, Nginx) and classifies them as **Outdated**, **Almost Out-of-Date**, **Up-to-Date**, or **Not Found**, based on vendor support timelines. As runtimes approach or pass EOL, Aikido raises alerts and escalates severity as the deprecation date nears or passes. This is so container owners get an early signal that security updates will stop, even if everything still appears to work.

Finally, for images scanned by Aikido’s own scanner, you can export a raw Software Bill of Materials. The SBOM includes file system locations of discovered components and information about their origin. Together, vulnerabilities, runtime status, and SBOM data form the basis of Aikido’s container security model.

### Where Aikido gets image data

Aikido treats your container registry (or build pipeline) as the source of truth for which images exist, and then standardizes the findings across different environments.

#### Cloud registries

For most cloud registries, you can either use the cloud provider's scanner (e.g., AWS Inspector for AWS ECS) or the Aikido scanner. When the Aikido Scanner is used, Aikido adds:

* License and EOL runtime scanning.
* Tag-level targeting (i.e., scan only images matching specific tags)
* Continuous scanning on a daily schedule, even for images that haven’t been pushed recently, so newly disclosed CVEs are picked up over time.

#### Standalone registries

Aikido also connects directly to standalone registries like Docker Hub. After you provide a read-only access token, Aikido enumerates repositories you can access, optionally links them to code repositories for better deduplication, and scans them for vulnerabilities with results flowing into your Aikido feed.

#### Local image scanning and gating

Not all images ever leave your continuous integration system. For those, Aikido provides a local image scanner that you run in your own environment. This scanner can:

* Scan images locally and report issues back into Aikido.
* Act as a release gate, failing a pipeline when issues at or above a chosen severity are present before pushing to a registry.
* Act as a pull request gate, comparing the security posture of an image built from a pull request against a base commit, and failing when new issues above a threshold are introduced.

This lets you treat container security as a guardrail around build and release instead of just a periodic check after deployment.

### Linking images to code and cloud context

A vulnerability in an image is only useful if it can be mapped to the team and code that can fix it. Aikido’s model is explicitly built around linking containers to repositories and cloud assets.

You can link container images to their corresponding code repositories manually or via smart suggestions. Once linked:

* Container issues appear directly within the relevant repository
* Container AutoFix becomes available since Aikido can reliably locate the Dockerfile associated with an image (see below)

### SBOMs and software supply chain visibility

A container's Software Bill of Materials helps you reason concretely about what is in the image and where it came from. Once an image has been scanned, you can export a raw SBOM from the container detail page. The export is designed to support:

* Inspecting the **f**ilesystem location**s** where specific components (and their vulnerabilities) were found.
* Tracking the origin of installed components across layers.

SBOM generation can also be automated via Aikido’s public API, which exposes endpoints for [generating](https://apidocs.aikido.dev/reference/generatecontainersbom) and [downloading](https://apidocs.aikido.dev/reference/exportcontainerrepolicenses) SBOMs for scanned images.

### Monitoring runtimes and end-of-life risk

Many high-impact issues in containers arise not from a single CVE, but from running on a runtime or web server that has fallen out of vendor support. After you connect a registry, Aikido monitors key runtimes in your images, including Python and Node.js runtimes and common web servers like Apache and Nginx.

For each monitored runtime, Aikido assigns statuses like "Outdated", "Up-to-Date", or "Almost Out-of-Date". Aikido also creates alerts in the feed when a specific package version is no longer maintained by the vendor, starting before deprecation and increasing issue severity as the date approaches and passes.

### AI AutoFix for containers

Container scanning identifies where vulnerabilities live. AutoFix for Containers helps you change the image to remove the vulnerabilities, with minimal manual steps.

#### Base-image-centric remediation

When Aikido finds vulnerabilities in a container’s base image, AutoFix suggests Dockerfile changes that update that base image. It does this by proposing multiple update options and generating 3–5 Dockerfile variants, each tied to a specific base image. For each variant, AutoFix shows which vulnerabilities would be fixed and whether any new ones would be introduced.

#### Aikido Images

When upgrading to a newer OS release or base image is expensive or risky, Aikido offers [**Aikido Images**](/autofix-and-remediation/aikido-images/autofix-for-containers-using-aikido-images.md), a registry of zero-CVE base images for containers. These are:

* Available for both maintained and end-of-life base images
* Provided as drop-in replacements that are intended to be compatible while removing known severe (critical and high) CVEs
* [More than 2000 base images](https://images.aikido.dev/images) are available, usable either as a long-term replacement or as an interim step until a full upgrade is feasible

#### Integration with repositories and registries

Container AutoFix relies on the earlier linking steps:

* Linking containers to repositories enables AutoFix for those containers and ensures Dockerfiles are located correctly
* AutoFix can work with both public and private base images, as long as private base images are scanned by Aikido.

When you select an option, AutoFix generates a pull request in your source control system containing the updated Dockerfile and associated changes. It can be run manually per image or configured to create PRs on a regular cadence for selected images.
