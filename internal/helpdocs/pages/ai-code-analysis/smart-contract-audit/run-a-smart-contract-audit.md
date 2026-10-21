# Run a Smart Contract Audit

Select contract repositories, add chain context, and start a Smart Contract Audit.

A Smart Contract Audit reasons about on-chain contract source inside [AI Code Analysis](/ai-code-analysis/ai-code-analysis.md). Connect the repositories, confirm the credit cost, and start. You do not need a deployed network or a fork.

## Prerequisites

* You need the **Manage Pentests and AI Code Analysis** [permission](/getting-started/automated-user-management/setting-roles-and-permissions.md) to create an audit. You can still open AI Code Analysis and existing assessments without it.
* At least one repository is connected. [Self-scanned repositories](/code-scanning/local-code-scanning.md) are excluded from the picker. You can pick up to 20 repositories.
* Your wallet has enough credits for the quoted cost. See [Wallet & Credits](/miscellaneous-info/wallet-and-credits.md).

See [Aikido Never Stores Your Code](/getting-started/setting-up-your-account/aikido-never-stores-your-code.md) for how source is handled during an analysis.

## Set up the audit

Go to [**AI Code Analysis**](https://app.aikido.dev/agentic-review), then click **Create Audit**.

{% stepper %}
{% step %}

### Select Smart Contract Audit

On **Select Audit Type**, select **Smart Contract Audit**, then click **Continue**.

**Cancel** returns to the AI Code Analysis home.
{% endstep %}

{% step %}

### Select repositories

1. Open **Repositories** and select up to 20 connected repositories. **Continue** stays disabled until at least one repository is selected.
2. To keep a monorepo focused, click **Exclude Paths** next to a selected repository and add directories that are not part of this audit, such as an off-chain service in the same repository.
3. Under **Scope**, leave **Audit all selected repositories** selected, or choose **Audit specific parts only** and describe the contracts to review. For example: "Audit the new staking contract and the contracts it calls."

{% hint style="info" %}
Exclude unrelated applications or services, not documentation, READMEs, or other non-code assets. Aikido can use those files as context for how the contracts are meant to work.
{% endhint %}

**Additional Notes** are optional. Add them when context is missing from the source. Useful notes include:

* Target chains
* Trusted admin or oracle roles
* Deployed contract addresses
* Known accepted risks

**Cancel** returns to **Audit Type**.
{% endstep %}

{% step %}

### Confirm pricing and start

Wait until the **Pricing** step shows the credit total next to **AI Smart Contract Audit**. If you limited the scope, the step also shows the scope you wrote. Your balance is in the footer. Click **Add Credits** if the wallet is short.

When **Configure** appears, use it to set the depth of the audit:

* **High risk only** prioritizes the areas most likely to expose serious vulnerabilities. Some areas get less coverage, or none.
* **Recommended coverage** investigates the security-sensitive parts of your selected audit scope.
* **Deeper coverage** builds on recommended coverage with additional investigation of lower-priority areas and less common code paths or cloud resources.

The credit total updates before you start.

Click **Start Audit** to launch. Credits are charged when you click **Start Audit**.
{% endstep %}
{% endstepper %}

Aikido opens the assessment detail page so you can follow progress and read findings. See [Track Progress and Findings](/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings.md).

## FAQ

### Why is there a 20 repository cap?

The picker allows 20 repositories so the audit stays focused. If you need more in one audit, [contact support](/miscellaneous-info/contact-support.md).

## Related

{% content-ref url="/pages/kLIRtXQIow79GwwO54tE" %}
[Smart Contract Audit](/ai-code-analysis/smart-contract-audit.md)
{% endcontent-ref %}

{% content-ref url="/pages/sDmfg35GWHVySIbqxISE" %}
[What Smart Contract Audit Finds](/ai-code-analysis/smart-contract-audit/what-smart-contract-audit-finds.md)
{% endcontent-ref %}
