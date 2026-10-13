# Exploit Intelligence with KEV and EPSS

Learn how Aikido uses EPSS and the CISA KEV catalog to prioritize dependency CVEs, and how to configure optional EPSS-based noise reduction.

Aikido uses exploit intelligence from two sources to help you focus on CVEs that matter: the **CISA Known Exploited Vulnerabilities (KEV)** catalog and the **Exploit Prediction Scoring System (EPSS)**. KEV always influences severity for actively exploited CVEs. EPSS-based prioritization is optional and can further **auto-ignore** or **downgrade** low-risk vulnerabilities based on predicted exploit likelihood.

## How exploit intelligence works in Aikido

Aikido refreshes both data sources daily and applies them when calculating each issue's severity score.

### CISA KEV (always on)

The [CISA KEV catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog) lists CVEs with confirmed active exploitation in the wild. Aikido syncs this list every day.

When a CVE appears on the KEV catalog, Aikido **raises its severity by 5 points** (up to a maximum score of 100). You'll see **Actively exploited vulnerability** in the issue's severity breakdown, with the description *It is actively being exploited in the wild.*

KEV-based scoring is always enabled. You don't need to turn anything on.

### EPSS (optional)

[EPSS](https://www.first.org/epss/) estimates the probability that a CVE will be exploited in the wild within the next 30 days. A score of 1% means there's roughly a 1-in-100 chance of exploitation in that window.

Aikido syncs EPSS scores daily. When you enable EPSS-based prioritization, Aikido can use those scores to ignore or downgrade findings that are unlikely to be exploited soon.

EPSS is **off by default**. Workspace admins enable it in [Triage Settings](https://app.aikido.dev/settings/advanced).

### How KEV and EPSS work together

KEV and EPSS serve different purposes. KEV reflects **confirmed** exploitation. EPSS predicts **future** exploitation likelihood.

Aikido applies them in a specific order during severity scoring:

1. **KEV**: If the CVE is on the CISA catalog, severity goes up by 5 points.
2. **Public proof-of-concept (PoC)**: If exploit code for the CVE is available on GitHub, severity goes up by 1 point.
3. **EPSS rules**: Only run when the CVE is **not** on KEV and has **no** public PoC.

{% hint style="info" %}
CVEs that are actively exploited or have a public PoC are never auto-ignored or downgraded by EPSS rules. This keeps real-world threat signals from being filtered out by low EPSS scores.
{% endhint %}

| Signal         | Source              | Default behavior    | Configurable?       |
| -------------- | ------------------- | ------------------- | ------------------- |
| **KEV**        | CISA catalog        | Severity +5         | No, always on       |
| **Public PoC** | GitHub exploit code | Severity +1         | No, always on       |
| **Low EPSS**   | FIRST.org           | Ignore or downgrade | Yes, optional rules |

Click an issue's severity score to see which rules applied, including KEV, PoC, and EPSS adjustments.

## Configure EPSS-based prioritization

{% stepper %}
{% step %}

### Open Triage settings

Go to **Settings** → **Triage** and find the **EPSS-Based Prioritization** section.
{% endstep %}

{% step %}

### Manage EPSS rules

Click **Manage** in the EPSS-based prioritization section.

[image: Settings for scan frequency and advanced scan configuration in a security monitoring dashboard.]
{% endstep %}

{% step %}

### Choose your rules

Select which EPSS rules to apply in your workspace:

* **Auto-ignore issues with an EPSS below 1%**: Auto-ignores CVEs with a very low predicted chance of exploitation.
* **Lower severity with 10 points for EPSS below 5%**: Downgrades severity by 10 points. For example, a high-severity issue with score 60 becomes 50.
* **Lower severity with 5 points for EPSS between 5% and 10%**: Downgrades severity by 5 points.

[image: EPSS-based vulnerability prioritization options for automated issue ignoring and severity adjustment.]
{% endstep %}

{% step %}

### Save and rescan

Click **Save EPSS-based prioritization** to apply your rules.

Manually trigger a rescan to apply the new prioritization immediately, or wait until the next scheduled scan for changes to take effect.
{% endstep %}
{% endstepper %}

## Important notes

* Aikido checks EPSS and KEV scores daily. When an EPSS value crosses a threshold, an issue may be **ignored** or **unignored** on the next scan.
* EPSS-based prioritization is available everywhere in Aikido: your IDE, PR gating in CI, and the feed.
* When you disable EPSS-based prioritization, issues that were previously auto-ignored by EPSS rules are unignored again. KEV and PoC severity adjustments are unaffected.
* EPSS rules never override KEV or public PoC signals. Actively exploited CVEs stay prioritized even when their EPSS score is low.

## Related features

* [**How is severity score calculated?**](/getting-started/core-functionalities/how-is-severity-score-calculated.md): Overview of Aikido's 0-100 severity model and contextual adjustments.
* [**CVE Exploitability Analysis**](/aikido-agent/cve-exploitability-analysis.md): Agent triage for dependency CVEs in your repositories and containers.
* [**Reachability analysis**](/getting-started/reachability-analysis/introduction-to-reachability-analysis.md): Filters findings with no execution path to vulnerable code.
