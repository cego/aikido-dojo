# AutoFix Overview

While most security programs are optimized for finding vulnerabilities instead of fixing them, **AutoFix** treats remediation as a core workflow in the Aikido platform. AutoFix generates concrete, reviewable patches for open-source dependencies, application code, infrastructure as code, containers, and even pentest findings, and delivers them through the same channels engineers already use: PRs, CI gating, and IDEs.

AutoFix is available for the following use cases

* [Open-source dependencies](/autofix-and-remediation/scope/autofix-for-open-source-dependencies.md) (including [Aikido Libraries](/autofix-and-remediation/aikido-libraries-overview.md) drop-in replacements)
* [SAST and IaC issues](/autofix-and-remediation/scope/ai-autofix-for-sast-and-iac-issues.md)
* [Containers](/autofix-and-remediation/scope/ai-autofix-for-containers.md) (including [Aikido Images](/autofix-and-remediation/aikido-images/autofix-for-containers-using-aikido-images.md) zero-CVE base images)
* [Code Quality](/code-quality/code-quality-overview.md)
* Pentest Issues
* AI Code Audit issues

***

### How AutoFix fits into your workflows

AutoFix is designed to show up where you already work:

* **IDE integration.** Aikido’s IDE plugins surface SAST, IaC, and dependency issues directly in the editor. When an issue is detected, AutoFix offers an inline diff that you can apply to your working copy, preventing the vulnerability from ever reaching a shared branch.
* **CI and PR gating.** When PR and release gating is enabled, Aikido scans pipeline branches and pull requests for new SAST and IaC issues. If new problems are introduced, AutoFix offers remediation suggestions that can be applied with a click.
* **Manual and automatic PR creation.** AutoFix PRs can be created manually from the Aikido UI, or configured to be generated daily for selected repositories and issue types. Aikido allows you to fine-tune PR metadata so that AutoFix PRs align with existing contribution guidelines and task tracking practices.

***

### What AutoFix covers

#### Open-source dependencies

For open-source dependencies, AutoFix proposes changes that remove vulnerabilities through package upgrades or other compatible changes. In many cases, a single AutoFix run can remove an entire class of vulnerabilities instead of addressing them one by one.

By default, AutoFix suggests the minimum version required to remediate a vulnerability; it prefers minor and patch bumps over major upgrades. If a major upgrade is proposed, it is because no lower version would fix the issue.

Under the hood, AutoFix analyzes the dependency tree to find the optimal point to apply a change. This is often by upgrading a direct dependency so that multiple vulnerable transitive dependencies are updated in one step.

**Aikido Libraries**

Vulnerable dependencies can be hard to remediate when the upstream fix requires a breaking major version change, or when the maintainer hasn't shipped a fix yet. Aikido addresses this with **Aikido Libraries**: drop-in replacements for vulnerable open-source packages that keep the same package name and base version while backporting security patches.

Aikido Libraries variants are designed to be 100% compatible with the original version while removing CVEs, and can serve either as a long-term solution or as a temporary stop-gap until a full major upgrade is feasible. AutoFix surfaces Aikido Libraries upgrades directly in the dependency AutoFix overview when available. Aikido Libraries variants are available for JavaScript, Java, Python, PHP, .NET, Go, Ruby, and Rust.

#### SAST and IaC issues

For application code, AutoFix generates rule-specific patches across a broad set of languages: JavaScript, Java, .NET, PHP, Python, Ruby, Go, Elixir, Rust, C, Kotlin, Scala, and Swift. Each SAST AutoFix is tied to a particular rule and remediation pattern. The instructions for that rule are tuned to keep changes minimal, preserve behavior, and address the root cause rather than only the immediate symptom.

AutoFix uses an agentic approach for these patches: it gathers the relevant code context, formulates a remediation plan, and then applies a minimal fix that is designed to be easy to read and review. For every SAST AutoFix, Aikido assigns a confidence level (High, Medium, Low) that reflects how strongly benchmarking suggests the generated change both fixes the vulnerability and preserves correct behavior.

Infrastructure as Code (IaC) AutoFixes follow the same pattern but operate on configuration rather than application logic. They target insecure infrastructure definitions, such as publicly exposed resources, weak TLS termination, or missing encryption, across AWS, GCP, Azure, and Kubernetes.

#### Code Quality

AutoFix can fix [Code Quality](/code-quality/code-quality-overview.md) findings, including default checks and [custom rules](/code-quality/add-custom-code-rules.md). For a custom rule, add an AutoFix prompt when you create or edit the rule and validate it against a non-compliant example. Aikido suggests a fix on pull request comments and can open a pull request from a repository scan finding.

Turn on [automatic pull request creation](/autofix-and-remediation/automation-and-merging/automatic-creation-of-pull-requests.md) to open a pull request each day for critical and high severity findings on default checks.

#### Containers and base images

Container AutoFix combines information from image scans and Dockerfiles to propose concrete base-image updates and OS-level remediations.

The primary workflow focuses on base-image upgrades. When vulnerabilities are found in a container’s base image, Aikido suggests multiple update options: patch, minor, or major. For each option, AutoFix surfaces which vulnerabilities would be remediated and whether any new ones would be introduced. If the automatically selected tags are not ideal for a given stack, users can add additional tags to the evaluation set while keeping configuration overhead low.

Because scanning all relevant tags and evaluating their vulnerability profiles is non-trivial, container AutoFix runs a more expensive analysis to compute suggestions, but returns results that capture both remediation and residual risk. It supports both public base images and private registries, as long as those private images are scanned by Aikido.

**Aikido Images**

When moving to a newer base image would be disruptive, AutoFix integrates **Aikido Images**: a registry of zero-CVE base images that patch CRITICAL and HIGH severity vulnerabilities across both maintained and end-of-life base images. Aikido continuously rebuilds these images as new CVEs land. AutoFix proposes them when available, and applies them by switching the base image in the Dockerfile to the Aikido Images variant hosted on `docker.aikido.io`.

#### Pentest findings

AutoFix is also wired into Aikido’s AI pentesting. When a pentest identifies an issue, Aikido generates a detailed root-cause analysis that spans code, containers, and cloud configuration. AutoFix can then use this analysis to synthesize a code fix that addresses the underlying defect rather than just the observed symptom.

***

### Security and privacy considerations

#### SAST, IaC, and Container AutoFix

When generating fixes for SAST, IaC, or container issues, Aikido sends only the minimal required code snippets to AI models hosted on AWS Bedrock over encrypted channels. Neither Aikido nor AWS Bedrock uses this code for training or fine-tuning models.

#### Dependency AutoFix

On the supply-chain side, AutoFix assists with several structural hardening tasks:

* Pinning GitHub Actions to specific commits
* Pinning base images to digests
* Installing missing security libraries

AutoFix is ultimately intended to handle the large volume of “mechanical” security work (dependency upgrades, small code fixes, base image updates) so that human effort can focus on architecture and ambiguous trade-offs. Its defaults are designed to keep changes safe and understandable while still reducing risk.
