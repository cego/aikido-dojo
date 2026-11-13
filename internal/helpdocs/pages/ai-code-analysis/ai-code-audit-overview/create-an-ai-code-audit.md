# Create a Code Security Audit

The create flow is three steps: choose **Code Security Audit**, pick the repositories, then confirm the credit cost and start.

## Prerequisites

* You need the **Manage Pentests and AI Code Analysis** [permission](https://help.aikido.dev/getting-started/automated-user-management/setting-roles-and-permissions) to create an audit. You can still open AI Code Analysis and existing assessments without it.
* At least one repository is connected. [Self-scanned repositories](https://help.aikido.dev/code-scanning/local-code-scanning) are excluded from the picker. There is a limit on how many repositories you can pick. [Contact support](https://help.aikido.dev/miscellaneous-info/contact-support) if you need more.
* Your wallet has enough credits for the quoted cost. You can top up from inside the create flow. See [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits).

## Set up the audit

Go to **AI Code Analysis**, then click **Create Audit**.

{% stepper %}
{% step %}
### Select Code Security Audit

On **Select Audit Type**, select **Code Security Audit**, then click **Continue**.

**Cancel** returns to the AI Code Analysis home.
{% endstep %}

{% step %}
### Select repositories

1. Open **Repositories** and select the connected repositories for this audit. **Continue** stays disabled until at least one repository is selected.
2. To keep a monorepo focused, click **Exclude Paths** next to a selected repository and add directories that are not part of this audit.
3. Under **Scope**, leave **Audit all selected repositories** selected, or choose **Audit specific parts only** and describe what to review. For example: "Audit the new billing flow and the services it calls."
4. **Additional Notes** are optional and apply to this audit only. Add them when context is missing from the source, such as runtime configuration, known architecture quirks, or firewall and proxy settings. Put guidance that should apply to every audit in a markdown file in the repository. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files).

[image: Exclude Paths button next to a selected repository]

{% hint style="info" %}
Leave documentation and READMEs in the repository. Agents can read that markdown as context. Exclude unrelated applications or services.
{% endhint %}

**Cancel** returns to **Audit Type**.
{% endstep %}

{% step %}
### Confirm pricing and start

Wait until the **Pricing** step shows the credit total. Your balance is in the footer. Click **Add Credits** if the wallet is short.

The card title depends on scope:

* **Full AI Code Analysis** when you audit all selected repositories.
* **Scoped AI Code Analysis** when you chose **Audit specific parts only**. The step also shows the scope you wrote.

When **Configure** appears, use it to set the depth of the audit:

* **High risk only** prioritizes the areas most likely to expose serious vulnerabilities. Some areas get less coverage, or none.
* **Recommended coverage** investigates the security-sensitive parts of your selected audit scope.
* **Deeper coverage** builds on recommended coverage with additional investigation of lower-priority areas and less common code paths or cloud resources.

The credit total updates before you start.

Click **Start Audit** to launch. Credits are charged when you click **Start Audit**.
{% endstep %}
{% endstepper %}

Aikido opens the assessment detail page. See [Track Progress and Findings](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings).

## FAQ

### Why is there a repository limit?

The picker limits how many repositories you can include so the audit stays focused. If you need more in one audit, [contact support](https://help.aikido.dev/miscellaneous-info/contact-support).

### The credit estimate looks wrong, or a repository is missing

[Contact support](https://help.aikido.dev/miscellaneous-info/contact-support). Self-scanned repositories are excluded from the picker.

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview)
{% endcontent-ref %}

{% content-ref url="track-progress-and-findings.md" %}
[track-progress-and-findings.md](https://help.aikido.dev/ai-code-analysis/ai-code-audit-overview/track-progress-and-findings)
{% endcontent-ref %}
