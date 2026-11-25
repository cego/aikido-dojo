# Custom Registries

If your team installs packages through a private registry like AWS CodeArtifact, JFrog Artifactory, or GCP Artifact Registry, add its domain in Aikido. Device Protection only inspects installs on domains it knows about.

Packages served through a registered domain are checked for malware exactly like packages from the public registry, and your existing installation policies apply to them too.

## How it works

When you register a custom registry, Device Protection recognizes its domain as a proxy for a public ecosystem. Package installs that flow through that domain are inspected on the device and checked against Aikido's malware signatures, just like installs from the public registry itself. Your existing [installation policies](https://help.aikido.dev/aikido-device-protection/using-aikido-endpoint/configuring-installation-policies) apply to these installs too.

## Add a custom registry

1. Go to [Custom Registries](https://app.aikido.dev/endpoint-protection/custom-registries?regionRedirect=1) in your Aikido account.
2. Click **Add Custom Registry**.
3. Select the registry **Type**: AWS CodeArtifact, JFrog Artifactory, or GCP Artifact Registry.
4. Enter the registry **Domain**, for example `my-org.d.codeartifact.us-east-1.amazonaws.com`, `my-company.jfrog.io`, or `europe-west1-npm.pkg.dev`.
5. For JFrog Artifactory, select the **Ecosystem** that the repository proxies (for example Npm) and enter its **Repository Key** (for example `npm-remote`).
6. For GCP Artifact Registry, select the **Ecosystem** (for example npm), then enter the **Project ID** (for example `my-gcp-project`) and the **Repository** (for example `npm-remote`).
7. Save.

A JFrog repository key or a GCP repository maps to a single ecosystem, so add a separate entry for each repository you use.

The registry now appears in the list with its type and domain. Use the **⋮** menu on a row to edit or remove it.
