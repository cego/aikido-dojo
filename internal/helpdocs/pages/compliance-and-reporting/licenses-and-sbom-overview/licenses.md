# OSS Licenses

The **Licenses & SBOM** page gives you a complete picture of the open-source licenses used across your code and containers, the legal risk each license carries, and lets you export a Software-Bill-of-Materials (SBOM) for compliance and audit purposes.

You can find it in the main navigation under **Reports > Licenses & SBOM**.

## Overview

The overview provides a searchable list of every package we found and the license it ships under.

For each package you'll see:

* **Package name**: with an icon showing its language or ecosystem.
* **License type**: the license(s) we detected.
* **Legal risk**: how risky this license is for your project.
* **Location**: where the package was found. If it appears in several places, you'll see the first location and a "\[X] others" tooltip to view the rest.

### Acting on a package

Each package has an *actions* menu (the **⋯** icon) with the following options:

| Action                       | What it does                                                                                                                                                                                      |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Mark as Internal Package** | Tells Aikido this is your own code, not a third-party dependency. The package is excluded from license risk reporting. Useful for internal libraries that picked up an unexpected license signal. |
| **Change Risk Level**        | Override the legal risk for this specific package. Use this when the default risk doesn't match how you actually use the package.                                                                 |
| **Overwrite License**        | Manually correct the detected license. For example, if the auto-detection got it wrong or the package's license file is missing.                                                                  |

You can also click the info icon next to any license type for a simple explanation of what that license allows and restricts, or the locations link to see every place a package shows up.

### Bulk actions

Tick the checkboxes on multiple packages to perform *bulk actions* on them:

* **Mark as Internal License**: mark all selected packages as internal in a single step.
* **Change Risk Level**: apply a new risk level to every selected package.

## License groups

The **License Groups** tab is the same data viewed by license instead of by package. Each row represents one license type and shows:

* The license name (e.g. MIT, Apache-2.0, GPL-3.0).
* How many of your dependencies use that license.
* The legal risk level assigned to that license.

Click a license to jump back to the *Overview tab* filtered to just that license. From the row menu you can **Adjust risk** to change the risk level for the license itself, this affects every package using it, which is useful when your legal team has approved or rejected a license.
