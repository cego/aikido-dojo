# Scanning Dev Dependencies for CVEs

{% hint style="warning" %}
Aikido **always** scans dev dependencies for malware. The following article relates to CVE scanning and how to enable.
{% endhint %}

**By default**, Aikido will **not** report vulnerabilities for dependencies that are only installed on the developer machine (devdeps). The assumption here is that they will not ship to production and won't impact the security of your live product.

However, there are some cases in which scanning for dev dependencies might be a valuable addition to the other scans:

* Compliance reasons, including the software that is only available on the developer machines
* SvelteKit: packages are often marked as dev dependency although they make it into production.

We support **JavaScript, Java and Python.**

### Enabling Dev Dependencies Per Repository <a href="#enabling-dev-dependencies-per-repository" id="enabling-dev-dependencies-per-repository"></a>

It is possible to enable dev dependency scanning on per repo basis.

**Step 1.** Contact Aikido to enable the functionality. Quickest way to do this is contact us via chat.

**Step 2.** Go to the detail page of a specific repository and click '**Configure**'.

[image: Security scan dashboard showing critical PHP issues with statuses and recommended fix times.]

**Step 3.** Scroll down in the modal to enable dev dependency scanning

[image: Option to enable or disable scanning of developer dependencies with warning about false positives.]

**Step 4.** Trigger a manual rescan to see results immediately, or wait for the next scheduled scan.

{% hint style="info" %}
Dev dependency findings will also surface on the next push to the default branch via PR gating.
{% endhint %}
