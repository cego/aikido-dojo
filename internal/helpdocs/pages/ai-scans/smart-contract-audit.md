# AI Smart Contract scan

AI Smart Contract scan is the blockchain scan type in [AI Scans](https://help.aikido.dev/ai-scans/ai-code-analysis). Agents reason about on-chain logic, access control, accounting, and cross-contract flaws that involve libraries and integrations. They simulate transactions across multiple parties to find maximal extractable value (MEV) flaws, including ones that depend on an oracle price update.

An [AI Code scan](https://help.aikido.dev/ai-scans/ai-code-audit-overview) reviews application code, and it can include contract code in that broader review. An AI Smart Contract scan keeps the agents on the contracts.

## Use cases

* **Review contracts before a deployment:** Find exploitable flaws with less time and cost than a manual audit
* **Scope a release:** Limit the scan to the contracts you just changed
* **Record accepted risk:** Add known issues with **Add Notes**, or upload a document under **Additional Information**, so agents can prioritize what is still unexpected

## AI Code scan vs AI Smart Contract scan

|                   | AI Smart Contract scan                      | AI Code scan                                             |
| ----------------- | ------------------------------------------- | -------------------------------------------------------- |
| **Starts from**   | Contract repositories                       | Application repositories                                 |
| **Core question** | Which on-chain flaws chain into an exploit? | Which application flaws are exploitable from source?     |
| **Best fit**      | Solidity contracts (officially supported)   | Applications, including ones that also contain contracts |

## How it works

You pick the repositories that contain the contracts. There is no staging URL, chain RPC, or deployed address required. Agents read the contract source and follow calls across files.

Findings land in the **Issues** tab and in the Feed, with code-based evidence and remediation. An AI Smart Contract scan includes a **Threat Model** tab. After you push a fix, use **Re-test** on the finding.

To run one, see [Run an AI Smart Contract scan](https://help.aikido.dev/ai-scans/smart-contract-audit/run-a-smart-contract-audit). For what is in scope, see [What an AI Smart Contract scan finds](https://help.aikido.dev/ai-scans/smart-contract-audit/what-smart-contract-audit-finds).

## Related

{% content-ref url="run-a-smart-contract-audit.md" %}
[run-a-smart-contract-audit.md](https://help.aikido.dev/ai-scans/smart-contract-audit/run-a-smart-contract-audit)
{% endcontent-ref %}

{% content-ref url="../ai-code-audit-overview/what-ai-code-audit-finds.md" %}
[what-ai-code-audit-finds.md](https://help.aikido.dev/ai-scans/ai-code-audit-overview/what-ai-code-audit-finds)
{% endcontent-ref %}
