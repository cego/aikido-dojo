# Cyber Resilience Act (CRA) Compliance

An overview of how Aikido helps you meet the Cyber Resilience Act (CRA) requirements.

The European Cyber Resilience Act (CRA) defines the cybersecurity requirements for all products with digital elements sold in the EU. One of the main requirements is reporting vulnerabilities and incidents to ENISA, Europe's cybersecurity agency.

## What do I have to report under the CRA?

You only have a reporting obligation once you become **aware** that a vulnerability in your product is being actively exploited, or that a severe incident has occurred. Becoming aware can take many forms:

* A customer or user notifying you that they were attacked or compromised through your product.
* You notice an indicator of compromise (IoC) during debugging or servicing
* Your monitoring or logging flagging malicious behavior.
* [Zen](/zen-firewall/getting-started-with-zen-firewall.md) flagging suspicious activity, for example blocked outbound requests could indicate a system has been compromised.

{% hint style="info" %}
Note: Becoming aware of **exploitation with malicious intent** is the trigger, not just exploitability. Finding a flaw is not the same as it being exploited with malicious intent.
{% endhint %}

## Incident reporting obligations

From 11 September 2026, Article 14 of the CRA requires manufacturers to report to ENISA and their national CSIRT as soon as they become aware of either:

* An **actively exploited vulnerability** for which there is reliable evidence that it has been exploited by a malicious actor, or
* A **severe incident** affecting that product's security posture.

Reporting happens in three stages:

* **Early warning, within 24 hours**: for an *actively exploited vulnerability*, which EU member states the product is available in. For a *severe incident*, whether it was caused by malicious acts and which member states are affected.
* **Notification, within 72 hours**: provide information about the product, the nature of the exploit or incident, any mitigating measures available, and guidance users can act on.
* **Final report**: for an *actively exploited vulnerability*, no later than 14 days after a fix or mitigation becomes available. For a *severe incident*, within one month of the 72-hour notification.

All three reports go through the ENISA Single Reporting Platform, which routes them to your coordinating CSIRT and ENISA at the same time. Aikido helps you identify and prioritize what needs to be reported. Filing the report itself is still your responsibility.

## Common misconceptions

* **"Every vulnerability with a known exploit must be reported."** A vulnerability having a public exploit does not automatically mean that it has been exploited on your product. It's important to prioritize remediating vulnerabilities with known exploits, but having an public exploit or high EPSS score is not a trigger for the reporting requirements.
* **"I should wait until I've fully investigated before reporting."** The 24-hour early warning starts the moment you become aware, even if you haven't fully investigated the situation. You can add more details as your investigation progresses, in the 72-hour notification and the final report.
* **"I only have to notify ENISA & the national CSIRT via the single reporting platform (SRP)"** You must also provide actionable risk mitigations to all impacted users and, where appropriate, all users.
* **"A pentest or bug bounty finding is reportable."** These typically do not qualify as malicious exploitation. The CRA only requires a report once there's "reliable evidence" the vulnerability "has been exploited by a malicious actor", and a pentester or bug bounty researcher finding a flaw doesn't qualify due to their non-malicious intent. You only need to report if you can confirm the vulnerability was actually exploited by a malicious actor.
* **"The final report for a vulnerability is due 14 days after becoming aware of it."** The 14-day clock doesn't start at awareness. It starts once a corrective or mitigating measure is available. If you don't have a fix yet, the final report isn't due, though your 24-hour early warning and 72-hour notification obligations still apply from the moment you become aware.

## How Aikido helps you avoid reportable events

There's no silver-bullet solution that can give you a list of vulnerabilities you're required to report. That list depends on evidence of real-world malicious exploitation of your specific product, which lives in your logs, your customers' reports, your support team's inbox, or in your firewall alerts, not in scan data.

What Aikido can do is show you which vulnerabilities are most likely to be exploited, so you can fix or mitigate them before they turn into an incident you have to report to ENISA and your users.

* Aikido monitors the CISA Known Exploited Vulnerabilities (KEV). Use the *"Exploit Status"* filter in the main feed to get a list of vulnerabilities that have a known exploit and should be prioritised. These are the most likely to cause an severe incident or become actively exploited, ultimately triggering the reporting obligations.
* Open an issue's severity score to confirm whether **Actively being exploited in the wild** is listed as a contributing factor. See [Exploit Intelligence with KEV and EPSS](/code-scanning/miscellaneous/use-epss-values-to-further-reduce-noise.md) for how this scoring works.
* Use [CVE Exploitability Analysis](/aikido-agent/cve-exploitability-analysis.md) to understand how a CVE could be exploited in your specific codebase. This is useful for prioritization, and, if you do end up filing a report, for describing the technical nature of the exploit.
* If you run [Zen](/zen-firewall/getting-started-with-zen-firewall.md) in front of your application, check for attacks it has already detected or blocked, one of the triggers of the reporting obligation. For example, *blocked* outbound traffic can be an indicator of system compromise or an ongoing incident.

{% hint style="info" %}
Aikido helps you prioritize and fix vulnerabilities before they're exploited. It cannot tell you which vulnerabilities you're legally required to report, since that depends on evidence of actual malicious exploitation that only you can confirm, and it does not submit reports to ENISA or your CSIRT on your behalf.
{% endhint %}

## How Aikido helps you meet CRA Annex I requirements

| CRA Requirements                          | Aikido Features                                           |
| ----------------------------------------- | --------------------------------------------------------- |
| No known exploitable vulnerabilities      | SCA, Container scanning, SAST/DAST, CSPM, Secret scanning |
| Protect against unauthorized access       | CSPM, IaC Scanning, DAST, AI Pentest, DSPM, Zen Firewall  |
| Secure-by-default                         | DAST (EASM), CSPM, DSPM, Aikido Images, Zen Firewall      |
| Secure-by-design                          | CSPM, Aikido Images, Zen Firewall                         |
| Limit attack surface                      | DAST (EASM), CSPM, Aikido Images, Zen Firewall            |
| Data encryption & integrity               | CSPM, DSPM, Secret scanning, DAST, AI Pentest             |
| Data minimisation                         | DSPM                                                      |
| Software Bill of Materials (SBOM)         | SBOM Export (CycloneDX), SCA, Container scanning          |
| Remediate vulnerabilities without delay   | AutoFix, AutoShip, SLA Reporting, Aikido Libraries        |
| Regular testing and security reviews      | SAST/DAST, AI Pentest / Code Analysis, PR checks          |
| Logging & monitoring of internal activity | Device Protection, Zen Firewall                           |
| Resilience against Denial of Service      | Zen Firewall                                              |
| Minimise negative impact on other devices | Zen Firewall (Outbound)                                   |
| Third-party and open-source due diligence | SBOM & Licenses, Malware, Intel, Device Protection        |

For more information, see [this CRA blogpost](https://www.aikido.dev/blog/cyber-resilience-act-compliance#how-aikido-security-simplifies-cra-compliance-requirements) on our website.

## Need help?

If you're not sure whether a finding qualifies as a reportable event under the CRA, open the **chat** in the bottom right corner. Our team is here to help!
