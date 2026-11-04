# Aikido Libraries

Secured, drop-in replacements for vulnerable packages. No upgrades. No breaking changes.

Upgrading a vulnerable dependency means using a version you haven't tested, risking breaking changes, and coordinating deploying that upgrade across your team. Transitioning to Aikido Libraries removes that tradeoff. It patches the exact version already in use - same package, same version, same API, same behavior - so you get the security fix without touching your application code.

Aikido has **11,500+** secured libraries across JavaScript, Java, Python, PHP, Go, .NET, Ruby, and Rust. Aikido Libraries is most useful when:

* A CVE has no safe upstream fix yet
* The upstream fix requires a major version bump you cannot take
* A package is no longer maintained

### Roll out Aikido Libraries to one or more repositories

1. In the left sidebar, go to **AutoFix** and click **Aikido Libraries**.

2. Open the **Bulk Rollout** tab and click **Select Repos**.
3. Click **Add Repo** and select the repositories you want to protect.
4. Review the credit breakdown in the **Credits** panel on the right.
5. Click **Enable Bulk Rollout** to confirm.

Once enabled, Aikido opens a daily PR that pins your packages to secured Aikido variants. New CVEs are picked up automatically and included in subsequent PRs - no manual intervention needed.

<details>

<summary>Pricing</summary>

Repository protection is a credit-based 1-year subscription, renewed annually. The cost is based on the number of packages and open CVEs in your repository at the time of enrollment. The exact cost is shown in the Credits panel before you confirm.

</details>

<details>

<summary>What gets sent to Aikido</summary>

When a repository is protected, Aikido sends this information for each package so it can match and build the secured variants:

* CVE identifier
* Package name
* Package version

No source code from your repository is sent.

</details>

### Fix a single dependency

When an Aikido Libraries variant is available for a vulnerable dependency, Aikido surfaces it as a fix option in two places. No credits required for one-off fixes.

* **From the Feed:** open the issue detail and select the Aikido Libraries option when creating the AutoFix PR.
* **From AutoFix > Dependencies:** click **View Fix** on the dependency and choose the Aikido Libraries variant in the modal.<br>

<details>

<summary>Example: log4j@1.2.16</summary>

Although released more than 20 years ago, version 1.2.16 of log4j is still used by many teams because upgrading to log4j 2.x is a breaking change.

Version `1.2.16` contains these critical issues:

* [CVE-2019-17571](https://nvd.nist.gov/vuln/detail/cve-2019-17571) - remote code execution
* [CVE-2020-9493](https://nvd.nist.gov/vuln/detail/cve-2020-9493) - malicious code execution
* [CVE-2022-23305](https://nvd.nist.gov/vuln/detail/cve-2022-23305) - SQL injection

The Aikido Libraries variant of `log4j 1.2.16` fixes these CVEs without forcing a migration to `2.x`.

You can safely transition from your existing log4j\@1.2.16 library to the respective Aikido version of it using AutoFix.

AutoFix proposes a transition to an Aikido Libraries variant when one is available. You can review this in the [AutoFix dependency overview](https://app.aikido.dev/issues/fix).

To transition to the Aikido Libraries variant, select this option in the AutoFix creation modal.

</details>
