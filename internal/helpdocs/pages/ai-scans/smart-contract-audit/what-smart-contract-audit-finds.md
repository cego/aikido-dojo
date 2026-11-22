# What an AI Smart Contract scan finds

An AI Smart Contract scan reasons about the contract source in the repositories you select. Agents assess weaknesses both individually and in combination, tracing how they can form an exploit chain across contracts and transactions. Each finding includes code-based **Evidence**. See [Track Progress and Findings](https://help.aikido.dev/ai-scans/ai-code-audit-overview/track-progress-and-findings) for how to read it.

For how this differs from an AI Code scan, see [AI Smart Contract scan](https://help.aikido.dev/ai-scans/smart-contract-audit). To run a scan, see [Run an AI Smart Contract scan](https://help.aikido.dev/ai-scans/smart-contract-audit/run-a-smart-contract-audit).

## Coverage at a glance

AI Smart Contract scan focuses on these classes:

* **Protocol logic bugs:** Can valid sequences of actions bypass intended rules, pay claims twice, or leave withdrawals stuck?
* **Access control:** Across owner, admin, operator, and other roles, can any actor perform a privileged action or change permissions without authorization?
* **Accounting:** Can asset-to-share conversions, rounding, or fee and debt calculations create unbacked claims or incorrect payouts?
* **Oracle and external data:** Can stale, manipulated, or incorrectly scaled prices and other feeds distort protocol decisions?
* **Economic and transaction-ordering risks, including maximal extractable value (MEV):** Can inadequate slippage protection, first-mover advantages, or front-running let actors extract value from users?
* **External calls:** Can interactions with other contracts exploit the order of state updates?
* **Availability and gas limits:** Can unbounded loops, queue growth, or reverting dependencies block withdrawals, liquidations, or other critical actions?

Agents follow contract calls across files, libraries, and other protocols, so a finding can span more than one function or contract.

**Configure** on the **Pricing** step changes how much of the selection gets investigated. It does not change which classes are in scope. **Scan specific parts only** and **Additional Information** are set in the create flow. Use **Add Notes** or **Upload File** to record target chains, trusted admin or oracle roles, deployed contract addresses, and known accepted risks when that context is not in the repository. Commit standing guidance as markdown in the repository. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files).

## Limits

AI Smart Contract scan reviews the source code in the selected repositories and the context you provide. It does not send transactions, run a fork, or verify deployed bytecode, configuration, or current on-chain state. A deployed address in your notes or an uploaded file provides context for that source review. It does not make the scan a live-chain assessment.

AI Smart Contract scan checks how contracts handle inputs from off-chain systems, such as oracle data, signed messages, and keeper or relayer calls, including when those inputs are delayed, missing, or invalid. Use an [AI Code scan](https://help.aikido.dev/ai-scans/ai-code-audit-overview) for the off-chain software itself, including frontends, APIs, keepers, and relayers.

After you push a fix, use **Re-test** on the finding. See [Track Progress and Findings](https://help.aikido.dev/ai-scans/ai-code-audit-overview/track-progress-and-findings).

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-scans/smart-contract-audit)
{% endcontent-ref %}

{% content-ref url="../ai-code-audit-overview/" %}
[ai-code-audit-overview](https://help.aikido.dev/ai-scans/ai-code-audit-overview)
{% endcontent-ref %}

{% content-ref url="../ai-code-audit-overview/track-progress-and-findings.md" %}
[track-progress-and-findings.md](https://help.aikido.dev/ai-scans/ai-code-audit-overview/track-progress-and-findings)
{% endcontent-ref %}

{% content-ref url="run-a-smart-contract-audit.md" %}
[run-a-smart-contract-audit.md](https://help.aikido.dev/ai-scans/smart-contract-audit/run-a-smart-contract-audit)
{% endcontent-ref %}
