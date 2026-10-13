# SBOM and VEX Export

Export a CycloneDX, SPDX, or CSV software bill of materials from Aikido, and add VEX exploitability statements when a security team needs them.

If an auditor, a customer, or a procurement team asks you for a machine-readable inventory of your software, you need the format the file is in and a clear boundary of what the export covers. This page gives you both. Field-level detail lives in [SBOM](/compliance-and-reporting/licenses-and-sbom-overview/sbom.md).

## What you can export today

| Artifact | Formats                      | Where                                              |
| -------- | ---------------------------- | -------------------------------------------------- |
| **SBOM** | CycloneDX 1.6, SPDX 2.3, CSV | Dashboard and API                                  |
| **VEX**  | CycloneDX 1.6                | Dashboard, as an option on a CycloneDX SBOM export |

An SBOM (Software Bill of Materials) lists every open-source component in your software: what it is, which version you use, who made it, and which license it's under. A VEX (Vulnerability Exploitability eXchange) document goes one step further and states, per vulnerability, whether you're affected and why. Auditors typically ask for the SBOM. Security teams downstream of you typically ask for both.

Aikido generates both on demand from your current scan data. The export is immediate, and nothing is queued.

## The three CycloneDX options, and when to use them

When you choose CycloneDX in the export dialog, three toggles appear:

* **Include VEX Analysis** adds the exploitability statement per vulnerability. Turn this on when the recipient is a security team, or when you're demonstrating active vulnerability management alongside the inventory.
* **Include Package Hashes** adds a hash per package so the recipient can verify the components are byte-for-byte the ones you listed. Turn this on for M\&A, procurement, and anything where tamper evidence matters.
* **Include Dependencies** adds the relationship graph: which packages are direct, and which are pulled in transitively.

If your SBOM is purely for license tracking, leave all three off. [SBOM](/compliance-and-reporting/licenses-and-sbom-overview/sbom.md) describes the same options field by field.

For the CycloneDX VEX, Aikido documents exploitability via the *analysis.state* field, while recommended actions (e.g. remediation advice) are provided in the *analysis.response* field.

## Export paths

### From the dashboard

{% stepper %}
{% step %}

### Open Licenses & SBOM

Go to **Reports >** [**Licenses & SBOM**](https://app.aikido.dev/licenses).
{% endstep %}

{% step %}

### Filter to the scope you want

The export always reflects your current filters, so filter first and export second. You can filter by license type, language, risk level, and container settings, and by team from the dropdown at the top of the page.
{% endstep %}

{% step %}

### Download the SBOM

Click **Download SBOM** in the page header and choose your format.
{% endstep %}
{% endstepper %}

{% hint style="info" %}
If you use multi-branch scanning, select the branch repository in the dropdown before you export. Each branch produces its own SBOM.
{% endhint %}

### From a single container

Open the container, then choose **Actions > Export RAW SBOM**. This returns that one container image on its own, regardless of the report-level filters. See [Export RAW SBOM of Your Containers](/container-image-scanning/configuration/export-raw-sbom-of-your-containers.md).

### From the API

Use the Export SBOM endpoint to generate and download SBOMs programmatically. Code repositories and containers have separate endpoints, matching the split in the dashboard:

* [Export a source code SBOM](https://apidocs.aikido.dev/reference/exportcoderepolicenses)
* [Export a container SBOM](https://apidocs.aikido.dev/reference/exportcontainerrepolicenses)

See the [API authorization docs](https://apidocs.aikido.dev/reference/authorization) for your API key.

## One document or several?

You can export one combined file or one file per container, depending on the path you take.

* **One document.** An unfiltered export from **Reports > Licenses & SBOM** covers your code repositories, your container images, and your self-reported SBOM uploads in a single file. Self-reported uploads are included because they are stored as container entries.
* **Separate documents.** An export from an individual container page, or from the per-container API endpoint, returns that container alone.

Use one document when someone asks for your SBOM. Use a separate document when a customer asks about one shipped image.

{% hint style="info" %}
You can combine multiple repositories (such as backend and frontend) into a single product SBOM. This is particularly useful for regulations like the Cyber Resilience Act (CRA), under which an SBOM becomes a mandatory document for virtually all digital products sold in Europe starting December 11th, 2027.
{% endhint %}

Self-generated SBOMs from embedded builds, such as Buildroot or Yocto, can be uploaded through the [Upload Container SBOM API](https://apidocs.aikido.dev/reference/uploadcontainersbom). They then appear under **Containers** with the registry label **Self-reported SBOM**. See [Generate SBOM Based on Open-Source Packages](/code-scanning/miscellaneous/generate-sbom-based-on-open-source-packages.md) for the build commands.

## Check the format yourself

Export CycloneDX 1.6 from your own account and run the file through your pipeline before you commit to a format. A full SBOM lists components, versions, suppliers, and licenses. Turn on **Include VEX Analysis** to add the exploitability statement per vulnerability. That file is the same output you'll get for any project in your account.

If your pipeline expects SPDX 2.3 or CSV, export those directly from **Reports > Licenses & SBOM**. CycloneDX carries the most detail: VEX, package hashes, and the dependency graph are CycloneDX only. See [Scope and limits](#scope-and-limits).

## Scope and limits

These are the points auditors come back to.

* The export reflects the current dependency graph at export time. Aikido generates it live from your latest scan data.
* Aikido holds the present state. To prove what your software contained on a given date, keep your own snapshots. See the pattern below.
* Container coverage is at the layer level. Packages inside scanned container layers are included. For firmware images and executables without a manifest, generate an SBOM at build time and upload it as a self-reported SBOM.
* VEX is CycloneDX only. SPDX and CSV exports list the inventory.

{% hint style="warning" %}
You cannot retrieve the SBOM as it looked on a past date, including 90 days ago. Record keeping stays on your side.
{% endhint %}

## Recommended pattern: point-in-time evidence

If you need to prove what your software contained on a given date (a release audit, a Cyber Resilience Act (CRA) or Medical Device Regulation (MDR) obligation, or a customer contract with a retention clause), then build an SBOM archive on your side. It takes one scheduled job to export the file and store it in your technical documentation.

{% stepper %}
{% step %}

### Export on a schedule

Export manually from **Reports > Licenses & SBOM**, or automate it through the [Export SBOM endpoint](https://apidocs.aikido.dev/reference/exportcoderepolicenses).
{% endstep %}

{% step %}

### Tie the schedule to something meaningful

Export once per release, or on a fixed interval such as every 90 days. Per release is the stronger evidence, because the artifact matches a version number your customer can name.
{% endstep %}

{% step %}

### Store the snapshots

Keep them in your own object storage or artifact repository. Put the export date and the release tag in the filename, and apply your own retention policy.
{% endstep %}

{% step %}

### Include the VEX in the SBOM

Archiving the VEX ensures you capture the exact vulnerability status and known security risks at the time of generation. This can prevent future questions from auditors, when new vulnerabilities are discovered after the release.
{% endstep %}
{% endstepper %}

Aikido gives you the current state on demand, in the format your auditor asks for. You keep the timeline, and a scheduled export is all it takes.

## Related

* [SBOM](/compliance-and-reporting/licenses-and-sbom-overview/sbom.md): the export dialog, field by field
* [CRA Compliance](https://github.com/AikidoSec/docs/tree/main/compliance-and-reporting/compliance-and-reporting/how-aikido-enables-compliance/cra-compliance.md): How Aikido enables CRA compliance
* [Open-source licenses](/compliance-and-reporting/licenses-and-sbom-overview/licenses.md): license risk per package and license groups
* [Generate SBOM Based on Open-Source Packages](/code-scanning/miscellaneous/generate-sbom-based-on-open-source-packages.md): Buildroot, Yocto, and self-reported uploads
* [Security Audit Report](/compliance-and-reporting/reports/security-audit-report.md): the shareable report auditors ask for alongside the SBOM
* [Compliance reporting](/compliance-and-reporting/reports/compliance-reporting.md): ISO 27001, SOC 2, and the other framework reports

## Need help?

Open the **Intercom chat** in the bottom right corner. Our team is here to help!
