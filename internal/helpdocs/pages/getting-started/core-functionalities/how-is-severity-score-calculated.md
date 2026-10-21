# How is Severity Score Calculated

Aikido provides a contextual, risk-based severity score from 0 to 100, offering 10× more granularity than traditional CVSS scoring (which is on a 0–10 scale). This allows for better prioritization and filtering.

<table><thead><tr><th width="143.943603515625">Severity</th><th width="155.927001953125">Score</th></tr></thead><tbody><tr><td>Critical</td><td>90 - 100</td></tr><tr><td>High</td><td>70 - 89</td></tr><tr><td>Medium</td><td>40 - 69</td></tr><tr><td>Low</td><td>1 - 39</td></tr></tbody></table>

### 1. Multiple Vulnerability Data Sources

We continuously monitor a variety of vulnerability feeds and databases, that provide baseline severity information and help establish initial severity scores. Databases include:

* Public vulnerability databases (e.g., NVD, GHSA,... | [See full list](https://help.aikido.dev/code-scanning/scanning-practices/external-vulnerability-databases-used-in-our-sca-engine))
* Operating system and vendor-specific advisories
* Our own Aikido Intel: <https://intel.aikido.dev/>

### 2. Contextual Severity Adjustments

To reflect actual risk more accurately, Aikido layers in additional context such as exploitability, environment, threat intelligence, and custom rules.<br>

**Exploitability & Threat Intelligence:**

Severity can increase when there’s evidence of real-world risk:

* The vulnerability is actively exploited or appears on the CISA KEV list
* A public PoC exploit is available (e.g., on GitHub)

#### Business Context

Severity is adjusted based on the importance of the affected asset. Some example are:

* Production vs test environments
* Backend vs frontend code
* Whether the vulnerable code is reachable or executed

#### Customer Rules

You can further refine issue scoring by adding contextual information to your project

* Learn how you can improve the risk score for repositories and containers [here](/code-scanning/miscellaneous/improve-risk-scoring-for-repositories-and-containers.md)

#### Exploit Prediction (EPSS)

Aikido also supports EPSS-based prioritization to automatically downgrade or ignore vulnerabilities that are unlikely to be exploited in the next 30 days. This is optional and turned off by default, more info here: [Exploit Intelligence with KEV and EPSS](/code-scanning/miscellaneous/use-epss-values-to-further-reduce-noise.md)

{% hint style="success" %}
You can click the score to view a detailed breakdown of why this issue received this severity rating
{% endhint %}
