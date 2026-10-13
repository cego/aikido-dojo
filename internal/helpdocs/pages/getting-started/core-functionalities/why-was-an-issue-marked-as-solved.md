# Why Was an Issue Marked as Solved

In Aikido, issues automatically move from **Open** to **Solved** when they are no longer detected in the latest scan.

There are several common reasons why this might happen:

***

### ✅ The Issue Was Fixed <a href="#the-issue-was-fixed" id="the-issue-was-fixed"></a>

Someone in your team resolved the problem, so it's no longer showing up in the scan because it has been **completely removed or fixed in code**.

***

### 🌿 Branch Change <a href="#branch-change" id="branch-change"></a>

If you're scanning a **different branch**, Aikido will:

* Close issues tied to the old branch
* Open any issues found on the new one
* 💡You can verify which branch is being scanned via the tag on the [repository page](https://app.aikido.dev/repositories)

  [image: Highlighted "master" branch in the about-github repository navigation menu.]

***

### 📁 File Relocation <a href="#file-relocation" id="file-relocation"></a>

If the affected file was **moved or renamed**, Aikido will:

* Mark the original issue as solved (since the file path changed)
* Potentially open a new issue under the new file location

***

### 🔐 CVE Database Update <a href="#cve-database-update" id="cve-database-update"></a>

Sometimes CVEs get updated to narrow down the versions affected.

**Example**:

* Aikido flagged a vulnerability in log4j:2.17.1 based on CVE-2021-44832.
* Later, the CVE was updated to say only versions ≤2.17.0 are affected.
* Aikido picked up the change and automatically removed the issue for 2.17.1.

***

### 📦 Dependency Moved to Dev-Dependency <a href="#dependency-moved-to-dev-dependency" id="dependency-moved-to-dev-dependency"></a>

If a vulnerable package is moved from a production dependency to a dev-dependency, Aikido will:

* Mark the original issue as solved (since the dependency type changed)
* Not longer report the issue as dev-dependencies are **by default** not scanned.
* 💡 You can enable scanning of dev-dependencies on a **per-repository** basis. [More information in this doc](/code-scanning/scanning-practices/scanning-dev-dependencies-for-cves.md)

***

### 🧠 SAST Rule Improvements <a href="#sast-rule-improvements" id="sast-rule-improvements"></a>

Static analysis rules are continuously improved to reduce false positives. As these rules become more accurate, Aikido may:

* **Mark previously flagged issues as solved**
* Prevent future false alarms for similar patterns

***

### 🗂️ Local Scanning – Scan Location Changed <a href="#local-scanning-scan-location-changed" id="local-scanning-scan-location-changed"></a>

When scanning locally, changing the **target directory** will influence what gets scanned.

* Issues in unscanned paths will be marked as solved
* New issues may show up in the new location

***

### Still Not Sure? <a href="#still-not-sure" id="still-not-sure"></a>

This list covers the most frequent causes, but there may be other edge cases.

👉 If you're unsure why something was marked as solved, reach out through the in-app chat. We're happy to help.
