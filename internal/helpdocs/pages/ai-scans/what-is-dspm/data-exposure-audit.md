# Run an AI Data Exposure scan

An AI Data Exposure scan is Aikido's code-based data security posture management (DSPM). It runs inside [AI Scans](https://help.aikido.dev/ai-scans/ai-code-analysis) and reasons over your repositories to find how sensitive data is handled. It does not connect to your databases.

## Prerequisites

* You need the **Manage Pentests and AI Scans** [permission](https://help.aikido.dev/getting-started/automated-user-management/setting-roles-and-permissions) to create a scan. You can still open AI Scans and existing assessments without it.
* At least one repository is connected. [Self-scanned repositories](https://help.aikido.dev/code-scanning/local-code-scanning) are excluded from the picker. There is a limit on how many repositories you can pick. [Contact support](https://help.aikido.dev/miscellaneous-info/contact-support) if you need more.
* Your wallet has enough credits for the quoted cost. See [Wallet & Credits](https://help.aikido.dev/miscellaneous-info/wallet-and-credits).

## Set up the scan

Go to [**AI Scans**](https://app.aikido.dev/agentic-review?regionRedirect=1), then open **New AI Scan**.

{% stepper %}
{% step %}
### Select AI Data Exposure scan

On **Select AI scan type**, select **AI Data Exposure scan**, then click **Continue**.

**Cancel** returns to the AI Scans home.

{% endstep %}

{% step %}
### Select repositories

1. Open **Repositories** and select the connected repositories for this scan. **Continue** stays disabled until at least one repository is selected.
2. To keep a monorepo focused, click **Exclude Paths** next to a selected repository and add directories that are not part of this scan.
3. Under **Scope**, leave **Scan all selected repositories** selected, or choose **Scan specific parts only** and describe the data handling to review.

[image: Exclude Paths button next to a selected repository]

{% hint style="info" %}
Leave documentation and READMEs in the repository. Agents can read that markdown for data flows and how sensitive data is handled. Exclude unrelated applications or services.
{% endhint %}

**Additional Information** is optional and applies to this scan only. Click **Add Additional Information**, then **Upload File** or **Add Notes**, when context is missing from the source, such as intentionally exposed private data, known issues to ignore, or proxy-level protection. Put guidance that should apply to every scan in a markdown file in the repository. See [Guide Agents with Repository Files](https://help.aikido.dev/ai-and-dev-tools/guide-an-audit-with-repository-files).

**Cancel** returns to **Select AI scan type**.

{% endstep %}

{% step %}
### Confirm pricing and start

Wait until the **Pricing** step shows the credit total. Your balance is in the footer. Click **Add Credits** if the wallet is short.

The card title depends on scope:

* **Full AI Scan** when you scan all selected repositories.
* **Scoped AI Scan** when you chose **Scan specific parts only**. The step also shows the scope you wrote.

When **Configure** appears, use it to set the depth of the scan:

* **High risk only** prioritizes the areas most likely to expose serious vulnerabilities. Some areas get less coverage, or none.
* **Recommended coverage** investigates the security-sensitive parts of your selected scan scope.
* **Deeper coverage** builds on recommended coverage with additional investigation of lower-priority areas and less common code paths or cloud resources.

The credit total updates before you start.

Click **Start Audit** to launch. Credits are charged when you click **Start Audit**.
{% endstep %}
{% endstepper %}

Aikido opens the assessment detail page. See [Track Progress and Findings](https://help.aikido.dev/ai-scans/ai-code-audit-overview/track-progress-and-findings).

For the data-handling classes in scope, see [What an AI Data Exposure scan finds](https://help.aikido.dev/ai-scans/what-is-dspm/what-data-exposure-audit-finds).

## FAQ

### Why is there a repository limit?

The picker limits how many repositories you can include so the scan stays focused. If you need more in one scan, [contact support](https://help.aikido.dev/miscellaneous-info/contact-support).

## Related

{% content-ref url="./" %}
[.](https://help.aikido.dev/ai-scans/what-is-dspm)
{% endcontent-ref %}

{% content-ref url="what-data-exposure-audit-finds.md" %}
[what-data-exposure-audit-finds.md](https://help.aikido.dev/ai-scans/what-is-dspm/what-data-exposure-audit-finds)
{% endcontent-ref %}
