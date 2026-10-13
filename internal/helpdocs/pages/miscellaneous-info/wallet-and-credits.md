# Wallet & Credits

## **What are Aikido credits?**

Aikido Credits are used to access powerful, resource-intensive features on a pay-as-you-go basis. This system gives you direct control over your consumption of these advanced tools, ensuring you only pay for what you use.

All your credits are managed in your **Aikido Wallet**, where you can check your balance, purchase more credits, view your transaction history, and set spending limits.

## **What can I use credits for?**

Aikido Credits unlock specific, high-value actions in Aikido:

* [**Aikido Libraries**](/autofix-and-remediation/aikido-libraries-overview.md): Keep every single open-source dependency in your repo secured automatically, for as long as the repo is in production. No manual upgrades, no breaking changes. Zero-CVE repos, forever.
* [**Deep Review**](/deep-review/how-deep-review-works.md): Get an agentic review of your pull requests, with relevant repository context beyond the diff. Costs **1 credit per pull request**. Additional reviews on the same PR don't cost extra.
* [**AI Code Analysis**](/ai-code-analysis/ai-code-analysis.md): Run a Code Security, Data Exposure, Cloud Security, or Smart Contract audit. Agents reason across your repositories or, for Cloud Security, your connected clouds.
* [**Aikido Pentest**](/pentests/aikido-pentest.md): Run advanced, AI-powered penetration tests on your applications, including compliance reports.
* [**Code Quality**](/code-quality/repository-scans/start-a-repository-scan.md): Enable full-repository code quality scanning for existing code.
* [**CVE Exploitability Analysis**](/aikido-agent/cve-exploitability-analysis.md): Have Aikido Agent assess open dependency CVEs in your repositories and containers and triage what's truly exploitable.

## Spending limits

Workspace admins can cap how many credits each feature spends from [Wallet settings](https://app.aikido.dev/settings/wallet/settings). Limits are set per feature, and the reset period depends on the feature.

1. Open **Organization Settings** and click **Manage Wallet**.
2. On the **Settings** tab, click **Configure Limits**.
3. Turn on a feature, enter a credit cap, and click **Save**.

When a feature reaches its cap, Aikido skips actions that would charge a credit until the limit resets. Turn a feature's toggle off if you don't want a cap.

## Monthly complimentary credits

Active workspaces on eligible plans can receive complimentary credits each calendar month:

| Plan      | Monthly complimentary credits |
| --------- | ----------------------------: |
| Developer |                  Not included |
| Basic     |                  Not included |
| Pro       |                           100 |
| Advanced  |                           200 |

One credit equals $1 USD of wallet balance. Aikido adds complimentary credits to the same wallet you use for paid features. They aren't a separate allowance or a monthly reset.

Your credits don't expire. Any credits left in your wallet stay there, including credits you bought and credits Aikido added as a monthly top-up.

### When you'll receive a top-up

Monthly top-ups are for active workspaces. Use your wallet credits on eligible features during the month to keep receiving complimentary credits.

If you don't use your wallet credits, Aikido may not add a complimentary top-up that month. Your existing wallet balance is never removed.

## **How to purchase credits?**

Follow these steps to add credits to your wallet:

1. Open **Organization Settings** and click **Manage Wallet**, then click **Add Credits**. You can also open [Wallet settings](https://app.aikido.dev/settings/wallet/settings) directly.
2. **Enter Purchase Details.**

* **Amount**: Enter the number of credits you wish to buy.
* **Minimum Purchase**: Please note there is a minimum purchase of 500 credits per transaction.
* **Total**: The total cost is calculated automatically for your review.

3. **Select a Payment Method.**

* Choose between **Credit Card** and **Invoice**. These options are explained in the next section.

4. **Confirm Your Purchase.**

* Click the **Confirm** button to proceed with the payment.

## **Payment methods**

You can purchase credits using one of two methods:

* **Credit Card**
  * This is the standard and fastest payment option. Selecting it will take you to a secure checkout flow to enter your card details.
* **Invoice**
  * This option is available **for paying customers only**. If you choose Invoice, we will generate and send one to your billing contact. Credits are added to your wallet upon successful payment.

## Automatic top-ups

Workspace admins can have Aikido buy more credits when the wallet drops below a balance you choose. Aikido emails an invoice to workspace admins, and the credits are added after that invoice is paid.

Set this up from the **Settings** tab in [Wallet settings](https://app.aikido.dev/settings/wallet/settings). See [Automatic Top-Ups](/miscellaneous-info/wallet-and-credits/automatic-top-ups.md) for the trigger, amount, monthly cap, and when a top-up is skipped.

## **Checking your credit history**

To see a full list of your credit purchases and usage, go to **Organization Settings** and click [**View History**](https://app.aikido.dev/settings/wallet/history) in the Wallet section.

Your credit history table includes:

* **Timestamp**: When the transaction happened.
* **Description**: What the transaction was for (e.g., "Credit Purchase," "AI Pentest").
* **Status**: The state of the transaction, such as "Completed," "Awaiting Payment," or "Failed."
