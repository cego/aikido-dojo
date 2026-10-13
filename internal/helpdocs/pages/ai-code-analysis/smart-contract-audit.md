# Smart Contract Audit

What Smart Contract Audit is, and when to use it instead of a Code Security Audit.

Smart Contract Audit is the blockchain audit type in [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md). Agents reason about on-chain logic, access control, accounting, and cross-contract flaws that involve libraries and integrations. They simulate transactions across multiple parties to find maximal extractable value (MEV) flaws, including ones that depend on an oracle price update.

A [Code Security Audit](/ai-code-analysis/ai-code-audit-overview.md) reviews application code, and it can include contract code in that broader review. Smart Contract Audit keeps the agents on the contracts.

## Use cases

* **Review contracts before a deployment:** Find exploitable flaws with less time and cost than a manual audit
* **Scope a release:** Limit the audit to the contracts you just changed
* **Record accepted risk:** Note known issues in **Additional Notes** so agents can prioritize what is still unexpected

## Code Security Audit vs Smart Contract Audit

|                   | Smart Contract Audit                        | Code Security Audit                                      |
| ----------------- | ------------------------------------------- | -------------------------------------------------------- |
| **Starts from**   | Contract repositories                       | Application repositories                                 |
| **Core question** | Which on-chain flaws chain into an exploit? | Which application flaws are exploitable from source?     |
| **Best fit**      | Solidity contracts (officially supported)   | Applications, including ones that also contain contracts |

## How it works

You pick the repositories that contain the contracts. There is no staging URL, chain RPC, or deployed address required. Agents read the contract source and follow calls across files.

Findings land in the **Issues** tab and in the Feed, with code-based evidence and remediation. A Smart Contract Audit includes a **Threat Model** tab. After you push a fix, use **Re-test** on the finding.

To run one, see [Run a Smart Contract Audit](/ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit.md). For what is in scope, see [What Smart Contract Audit Finds](/ai-code-analysis/smart-contract-audit/what-smart-contract-audit-finds.md).

## Related

{% content-ref url="/pages/rpzhJqLOk9BfP8JQGOSq" %}
[Run a Smart Contract Audit](/ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/vx7C5jq5oGOUqGbKvpE2" %}
[What Code Security Audits Find](/ai-code-analysis/ai-code-audit-overview/what-ai-code-audit-finds.md)
{% endcontent-ref %}
