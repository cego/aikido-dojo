# SBOM

A **Software Bill of Materials (SBOM)** is a structured inventory of every open-source component in your software, what it is, what version you use, who made it, and what license it's under. Customers, auditors, and regulators increasingly ask for one, and it's the foundation for tracking vulnerabilities and license obligations across your supply chain.

For formats, export paths, VEX, and what the file covers, see [SBOM and VEX export](/compliance-and-reporting/licenses-and-sbom-overview/sbom-and-vex-export.md).

## Downloading an SBOM

Click **Download SBOM** in the page header to open the export dialog. The export always reflects your **current filters**, so if you only want an SBOM for one repository or one container image, filter to it first.

### Choose a format

You can export in one of two formats:

* **CSV**: A flat comma-separated list of every package, its version, license, and location. Useful for pulling into Excel for ad-hoc analysis.
* **SPDX**: Software Package Data Exchange (SPDX) is an open SBOM standard mostly associated with open-source license compliance and copyright.
* **CycloneDX**: The industry-standard format for SBOMs, recognised by virtually every SBOM consumer, scanner, and compliance platform.

### CycloneDX options

When you pick CycloneDX, three additional toggles appear:

* **Include VEX Analysis**: VEX stands for *Vulnerability Exploitability eXchange*. An SBOM alone tells you *what packages you use*, VEX tells you *which packages have known vulnerabilities* **and** *if you're affected*, along with the reason.
* **Include Package Hashes**: Adds a hash of each package, which uniquely identifies the exact contents of that package version. The hashes allow anyone to verify that the packages they have are byte-for-byte the same as the ones you listed and haven't been tampered with.
* **Include Dependencies**: This option adds the relationships between packages. Which packages are direct dependencies of your project, and which are pulled in transitively by other packages.

If your SBOM is purely for license tracking you can safely ignore these options.

## Monitoring your self-generated SBOMs

To learn more about uploading your self-generated SBOMs, see [uploading your SBOM](/code-scanning/miscellaneous/generate-sbom-based-on-open-source-packages.md).
