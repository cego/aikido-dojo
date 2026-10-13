# What Smart Contract Audit Finds

What Smart Contract Audit checks for in contract code.

A Smart Contract Audit reasons about the contract source in the repositories you select. Agents assess weaknesses both individually and in combination, tracing how they can form an exploit chain across contracts and transactions. Each finding includes code-based **Evidence**. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md) for how to read it.

For how this differs from a Code Security Audit, see [Smart Contract Audit](/ai-code-analysis/smart-contract-audit.md). To run an audit, see [Run a Smart Contract Audit](/ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit.md).

## Coverage at a glance

Smart Contract Audit focuses on these classes:

* **Protocol logic bugs:** Can valid sequences of actions bypass intended rules, pay claims twice, or leave withdrawals stuck?
* **Access control:** Across owner, admin, operator, and other roles, can any actor perform a privileged action or change permissions without authorization?
* **Accounting:** Can asset-to-share conversions, rounding, or fee and debt calculations create unbacked claims or incorrect payouts?
* **Oracle and external data:** Can stale, manipulated, or incorrectly scaled prices and other feeds distort protocol decisions?
* **Economic and transaction-ordering risks, including maximal extractable value (MEV):** Can inadequate slippage protection, first-mover advantages, or front-running let actors extract value from users?
* **External calls:** Can interactions with other contracts exploit the order of state updates?
* **Availability and gas limits:** Can unbounded loops, queue growth, or reverting dependencies block withdrawals, liquidations, or other critical actions?

Agents follow contract calls across files, libraries, and other protocols, so a finding can span more than one function or contract.

**Configure** on the **Pricing** step changes how much of the selection gets investigated. It does not change which classes are in scope. **Audit specific parts only** and **Additional Notes** are set in the create flow. Notes are where you record target chains, trusted admin or oracle roles, deployed contract addresses, and known accepted risks when that context is not in the repository.

## Limits

Smart Contract Audit reviews the source code in the selected repositories and the context you provide. It does not send transactions, run a fork, or verify deployed bytecode, configuration, or current on-chain state. A deployed address in your notes provides context for that source review. It does not make the audit a live-chain assessment.

Smart Contract Audit checks how contracts handle inputs from off-chain systems, such as oracle data, signed messages, and keeper or relayer calls, including when those inputs are delayed, missing, or invalid. Use a [Code Security Audit](/ai-code-analysis/ai-code-audit-overview.md) for the off-chain software itself, including frontends, APIs, keepers, and relayers.

After you push a fix, use **Re-test** on the finding. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md).

## Related

{% content-ref url="/pages/kLIRtXQIow79GwwO54tE" %}
[Smart Contract Audit](/ai-code-analysis/smart-contract-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/kZRWK6NswIlL0dAmeNea" %}
[Code Security Audit](/ai-code-analysis/ai-code-audit-overview.md)
{% endcontent-ref %}

{% content-ref url="/pages/P12625XkjR5KcYlJmjTa" %}
[Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md)
{% endcontent-ref %}

{% content-ref url="/pages/rpzhJqLOk9BfP8JQGOSq" %}
[Run a Smart Contract Audit](/ai-code-analysis/smart-contract-audit/run-a-smart-contract-audit.md)
{% endcontent-ref %}
