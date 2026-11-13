# Run a Smart Contract Audit

A Smart Contract Audit reasons about on-chain contract source inside [AI Code Analysis](https://help.aikido.dev/ai-code-analysis/ai-code-analysis). Connect the repositories, confirm the credit cost, and start. You do not need a deployed network or a fork.

## Prerequisites

* You need the **Manage Pentests and AI Code Analysis** [permission](https://help.aikido.dev/getting-started/automated-user-management/setting-roles-and-permissions) to create an audit. You can still open AI Code Analysis and existing assessments without it.
* At least one repository is connected. [Self-scanned repositories](https://help.aikido.dev/code-scanning/local-code-scanning) are excluded from the picker. There is a limit on how many repositories you can pick. [Contact support](https://help.aikido.dev/miscellaneous-info/contact-support) if you need more.
* Your wallet has enough credits for the quoted cost. See [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits).

See [Aikido Never Stores Your Code](https://help.aikido.dev/getting-started/setting-up-your-account/aikido-never-stores-your-code) for how source is handled during an analysis.

## Set up the audit

Go to [**AI Code Analysis**](https://app.aikido.dev/agentic-review?regionRedirect=1), then click **Create Audit**.

{% stepper %}
{% step %}
### Select Smart Contract Audit

On **Select Audit Type**, select **Smart Contract Audit**, then click **Continue**.

**Cancel** returns to the AI Code Analysis home.
{% endstep %}

{% step %}
### Select repositories

1. Open **Repositories** and select the connected repositories for this audit. **Continue** stays disabled until at least one repository is selected.
2. To keep a monorepo focused, click **Exclude Paths** next to a selected repository and add directories that are not part of this audit, such as an off-chain service in the same repository.
3. Under **Scope**, leave **Audit all selected repositories** selected, or choose **Audit specific parts only** and describe the contracts to review. For example: "Audit the new staking contract and the contracts it calls."

{% hint style="info" %}
Leave documentation and READMEs in the repository. Agents can read that markdown for how the contracts are meant to work. Exclude unrelated applications or services.
{% endhint %}

**Additional Notes** are optional and apply to this audit only. Add them when context is missing from the source. Put standing guidance, such as the trust model, in a markdown file in the repository. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files). Useful notes include:

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

Aikido opens the assessment detail page so you can follow progress and read findings. See [Track Progress and Findings](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings).

## FAQ

### Why is there a repository limit?

The picker limits how many repositories you can include so the audit stays focused. If you need more in one audit, [contact support](https://help.aikido.dev/miscellaneous-info/contact-support).

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-code-analysis/smart-contract-audit)
{% endcontent-ref %}

{% content-ref url="what-smart-contract-audit-finds.md" %}
[what-smart-contract-audit-finds.md](https://help.aikido.dev/ai-code-analysis/smart-contract-audit/what-smart-contract-audit-finds)
{% endcontent-ref %}
