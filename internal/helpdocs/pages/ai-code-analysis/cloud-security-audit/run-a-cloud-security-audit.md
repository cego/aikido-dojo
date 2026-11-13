# Run a Cloud Security Audit

A Cloud Security Audit reasons about connected cloud accounts inside [AI Code Analysis](https://help.aikido.dev/ai-code-analysis/ai-code-analysis). It uses the inventory Aikido already has. You can add infrastructure as code (IaC) repositories so findings can point back at the templates that define the account.

## Prerequisites

* You need the **Manage Pentests and AI Code Analysis** [permission](https://help.aikido.dev/getting-started/automated-user-management/setting-roles-and-permissions) to create an audit. You can still open AI Code Analysis and existing assessments without it.
* At least one cloud account is [connected](https://help.aikido.dev/cloud-scanning/connect-your-cloud). Supabase projects do not appear in the picker.
* Your wallet has enough credits for the quoted cost. See [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits).

IaC repositories are optional. [Self-scanned repositories](https://help.aikido.dev/code-scanning/local-code-scanning) are excluded from the picker. There is a limit on how many repositories you can add. [Contact support](https://help.aikido.dev/miscellaneous-info/contact-support) if you need more.

## Set up the audit

Go to [**AI Code Analysis**](https://app.aikido.dev/agentic-review?regionRedirect=1), then click **Create Audit**.

{% stepper %}
{% step %}
### Select Cloud Security Audit

On **Select Audit Type**, select **Cloud Security Audit**, then click **Continue**.

**Cancel** returns to the AI Code Analysis home.
{% endstep %}

{% step %}
### Select clouds and IaC

1. Open **Clouds** and select the connected accounts to audit. **Continue** stays disabled until at least one cloud is selected.
2. To tie findings back to code, open **IaC repositories (optional)** and select the repositories that define those accounts. You can leave this empty.
3. To keep a repository focused, click **Exclude Paths** next to a selected repository and add directories that are not part of this audit.

{% hint style="info" %}
Leave documentation and READMEs in the repository. Agents can read that markdown for how the account is defined. Exclude unrelated applications or services.
{% endhint %}

**Additional Notes** are optional and apply to this audit only. Add them when context is missing from the inventory and the IaC, such as runtime configuration, known architecture quirks, or firewall and proxy settings. Put guidance that should apply to every audit in a markdown file in the repository. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files).

**Cancel** returns to **Audit Type**.
{% endstep %}

{% step %}
### Confirm pricing and start

Wait until the **Pricing** step shows the credit total next to **AI Cloud Security Audit**. Your balance is in the footer. Click **Add Credits** if the wallet is short.

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

### Why is a cloud missing from the picker?

Connect the account first. See [Connect Your Cloud](https://help.aikido.dev/cloud-scanning/connect-your-cloud). Supabase projects are not included in Cloud Security Audit.

### Do I have to add IaC repositories?

No. The audit runs on the cloud inventory alone. Add IaC when you want a finding to point at the template or module that defines the live setting.

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-code-analysis/cloud-security-audit)
{% endcontent-ref %}

{% content-ref url="what-cloud-security-audit-finds.md" %}
[what-cloud-security-audit-finds.md](https://help.aikido.dev/ai-code-analysis/cloud-security-audit/what-cloud-security-audit-finds)
{% endcontent-ref %}
